package manager

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/cloudness-io/cloudness/app/controller/variable"
	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/app/pipeline/convert"
	"github.com/cloudness-io/cloudness/app/pipeline/scheduler"
	"github.com/cloudness-io/cloudness/app/services/config"
	"github.com/cloudness-io/cloudness/app/services/githubapp"
	specSvc "github.com/cloudness-io/cloudness/app/services/spec"
	"github.com/cloudness-io/cloudness/app/sse"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/logstream"
	basestore "github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

var _ RunnerManager = (*runnerManager)(nil)

type (
	RunnerManager interface {
		// GetPrimaryServerUID returns the primary server UID.
		GetPrimaryServerUID(ctx context.Context) (int64, error) //NOTE: temporary

		// Config returns the runner config for the server.
		Config(ctx context.Context, serverUID int64) (*types.RunnerConfig, error)

		// Request requests the next available deployment for execution.
		Request(ctx context.Context) (*types.Deployment, error)

		// Accept accepts a deployment for execution.
		Accept(ctx context.Context, deploymentID int64, machine string) (*types.Deployment, error)

		// Details returns all the information for the deployment to be executed.
		Details(ctx context.Context, deploymentID int64) (*pipeline.RunnerContext, error)

		// Watch watches for deployment cancellation events.
		Watch(ctx context.Context, deploymentID int64) (bool, error)

		// Update updates the deployment with new information.
		Update(ctx context.Context, deployment *types.Deployment) (*types.Deployment, error)

		// stream streams real time logs to backend
		Stream(ctx context.Context, deploymentID int64, logs []*types.LogLine) error

		// Upload uploads the full logs to server
		Upload(ctx context.Context, deploymentID int64, r io.Reader) error
	}
)

type runnerManager struct {
	scheduler        scheduler.Scheduler
	serverStore      store.ServerStore
	applicationStore store.ApplicationStore
	deploymentStore  store.DeploymentStore
	volumeStore      store.VolumeStore
	logStore         store.LogStore
	varCtrl          *variable.Controller
	configSvc        *config.Service
	ghAppSvc         *githubapp.Service
	sseStreamer      sse.Streamer
	logz             logstream.LogStream
}

func New(scheduler scheduler.Scheduler,
	serverStore store.ServerStore,
	applicationStore store.ApplicationStore,
	deploymentStore store.DeploymentStore,
	volumeStore store.VolumeStore,
	logStore store.LogStore,
	varCtrl *variable.Controller,
	configSvc *config.Service,
	ghAppSvc *githubapp.Service,
	sseStreamer sse.Streamer,
	logStream logstream.LogStream,
) RunnerManager {
	return &runnerManager{
		scheduler:        scheduler,
		serverStore:      serverStore,
		applicationStore: applicationStore,
		deploymentStore:  deploymentStore,
		volumeStore:      volumeStore,
		logStore:         logStore,
		varCtrl:          varCtrl,
		configSvc:        configSvc,
		ghAppSvc:         ghAppSvc,
		sseStreamer:      sseStreamer,
		logz:             logStream,
	}
}

func (m *runnerManager) GetPrimaryServerUID(ctx context.Context) (int64, error) {
	server, err := m.serverStore.Find(ctx, 1)
	if err != nil {
		return 0, err
	}

	return server.UID, nil
}

func (m *runnerManager) Config(ctx context.Context, serverUID int64) (*types.RunnerConfig, error) {
	server, err := m.serverStore.FindByUID(ctx, serverUID)
	if err != nil {
		return nil, err
	}
	config := &types.RunnerConfig{
		EnableRunner:    server.BuildEnabled,
		Hostname:        server.Name,
		ParallelWorkers: server.MaxConcurrentBuilds,
		PollingInterval: 1,
		MaxCPU:          server.MaxCPUPerBuild,
		MaxMemory:       server.MaxMemoryPerBuild,
		UpdatedAt:       server.Updated,
	}
	switch server.Type {
	case enum.ServerTypeK8s:
		config.Engine = "kubernetes"
		config.KubeNameSpace = "cloudness"
	}
	return config, nil
}

func (m *runnerManager) Request(ctx context.Context) (*types.Deployment, error) {
	log := log.Ctx(ctx).With().Logger()

	log.Trace().Msg("manager: request deployment from queue")

	deployment, err := m.scheduler.Request(ctx, scheduler.Filter{})
	if err != nil && ctx.Err() != nil {
		log.Trace().Err(err).Msg("manager: context canceled")
		return nil, err
	}
	if err != nil {
		log.Warn().Err(err).Msg("manager: error requesting queue item")
		return nil, err
	}

	log.Trace().Msg("manager: deployment request successful")
	return deployment, nil
}

func (m *runnerManager) Accept(ctx context.Context, deploymentID int64, machine string) (*types.Deployment, error) {
	log := log.With().
		Int64("deployment.id", deploymentID).
		Str("machine", machine).
		Logger()

	log.Debug().Msg("manager: accept deployment request received")

	deployment, err := m.deploymentStore.Find(ctx, deploymentID)
	if err != nil {
		log.Warn().Err(err).Msg("manager: error finding deployment")
		return nil, err
	}

	if deployment.Machine != "" {
		log.Debug().Msg("manager: deployment already assigned")
		return nil, fmt.Errorf("deployment already assigned")
	}

	deployment.Machine = machine
	deployment.Status = enum.DeploymentStatusRunning

	err = m.deploymentStore.Update(ctx, deployment)
	switch {
	case errors.Is(err, basestore.ErrVersionConflict):
		log.Debug().Err(err).Msg("manager: deployment processed by another agent")
	case err != nil:
		log.Debug().Err(err).Msg("manager: cannot update deployment")
	default:
		log.Info().Msg("manager: stage accepted")
	}

	if err != nil {
		return nil, err
	}

	// create log stream for deployment
	if err := m.logz.Create(ctx, deploymentID); err != nil {
		log.Warn().Err(err).Msgf("manager: unable to create log stream for %d", deploymentID)
		return nil, err
	}

	return deployment, nil
}

func (m *runnerManager) Details(ctx context.Context, deploymentID int64) (*pipeline.RunnerContext, error) {
	deployment, err := m.deploymentStore.Find(ctx, deploymentID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("manager: error finding deployment")
		return nil, err
	}

	app, err := m.applicationStore.Find(ctx, deployment.ApplicationID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("manager: error getting application")
		return nil, err
	}

	var previousDeployment *types.Deployment
	if app.DeploymentID != nil {
		previousDeployment, err = m.deploymentStore.Find(ctx, *app.DeploymentID)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("manager: error finding previous deployment")
			return nil, err
		}
	}

	variables, err := m.varCtrl.ListWithUpdate(ctx, app.EnvironmentID, app.ID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("manager: error getting variables")
		return nil, err
	}

	volumes, err := m.volumeStore.List(ctx, &types.VolumeFilter{
		TenantID:      &app.TenantID,
		ProjectID:     &app.ProjectID,
		EnvironmentID: &app.EnvironmentID,
		ApplicationID: &app.ID,
	})
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("manager: error getting volumes")
		return nil, err
	}

	server, err := m.serverStore.Find(ctx, app.ServerID)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("manager: error getting server")
		return nil, err
	}

	config, err := m.configSvc.Pipeline(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("manager: error getting pipeline config")
	}
	serverRestrictions := m.configSvc.GetServerRestrictions(server)

	runnerCtxIn := &pipeline.RunnerContextInput{
		Application:        app,
		Variables:          variables,
		Deployment:         deployment,
		PreviousDeployment: previousDeployment,
		Volumes:            volumes,
		Config:             config,
		ServerRestctions:   serverRestrictions,
	}

	if app.GetGithubAppID() > 0 {
		netrc, err := m.ghAppSvc.Netrc(ctx, app.TenantID, app.ProjectID, app.GetGithubAppID(), specSvc.GetGitRepoUrl(app.Spec))
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("manager: error getting netrc for ghApp")
			return nil, err
		}
		runnerCtxIn.Netrc = netrc
	}

	runnerCtx, err := convert.ToRunnerContext(runnerCtxIn)
	if err != nil {
		deployment.Fail(err)
		m.Update(ctx, deployment)
		log.Ctx(ctx).Error().Err(err).Msg("manager: error converting context")
		return nil, err
	}

	return runnerCtx, nil
}

func (m *runnerManager) Update(ctx context.Context, deployment *types.Deployment) (*types.Deployment, error) {
	if err := m.deploymentStore.Update(ctx, deployment); err != nil {
		return nil, err
	}

	if deployment.Status == enum.DeploymentStatusRunning {
		if err := m.up(ctx, deployment); err != nil {
			return nil, err
		}
	} else {
		if err := m.down(ctx, deployment); err != nil {
			return nil, err
		}

		// destroy log streams
		if err := m.logz.Delete(ctx, deployment.ID); err != nil {
			log.Warn().Err(err).Msgf("manager: error deleting log stream for deployment %d", deployment.ID)
			return nil, err
		}
	}

	return deployment, nil
}

func (m *runnerManager) Watch(ctx context.Context, deploymentID int64) (bool, error) {
	cancel, err := m.scheduler.Cancelled(ctx, deploymentID)
	if err != nil {
		return cancel, err
	}

	// if no error received, check database
	deployment, err := m.deploymentStore.Find(ctx, deploymentID)
	if err != nil {
		log := log.With().Int64("deployment.id", deploymentID).Logger()
		log.Warn().Msg("manager: cannot find deployment")
		return cancel, fmt.Errorf("could not find deployment for cancellation: %w", err)
	}

	return deployment.IsDone(), nil
}

func (m *runnerManager) Stream(ctx context.Context, deploymentID int64, lines []*types.LogLine) error {
	log := log.With().Int64("deployment.id", deploymentID).Logger()
	for _, line := range lines {
		if err := m.logz.Write(ctx, deploymentID, line); err != nil {
			log.Trace().Err(err).Msg("manager: error writing log to stream")
			return err
		}
	}
	return nil
}

func (m *runnerManager) Upload(ctx context.Context, deploymentID int64, r io.Reader) error {
	if err := m.logStore.Create(ctx, deploymentID, r); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("manager: error uploading logs")
		return err
	}
	return nil
}
