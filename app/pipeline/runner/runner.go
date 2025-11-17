package runner

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline"
	"github.com/cloudness-io/cloudness/app/pipeline/logger"
	"github.com/cloudness-io/cloudness/app/pipeline/manager/client"
	"github.com/cloudness-io/cloudness/app/pipeline/runner/engine"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog"
)

var noContext = context.Background()

type Runner struct {
	// client is the remote client to communicate with server manager
	client client.RunnerClient

	// engine is the backend engine to run the deployment
	engine engine.Engine

	// config is the configuration of the runner
	config types.RunnerConfig
}

func NewRunner(client client.RunnerClient, config types.RunnerConfig, engine engine.Engine) *Runner {
	return &Runner{
		client: client,
		engine: engine,
		config: config,
	}
}

// Pending returns the deployment id of pending deployments
func (r *Runner) ListPending(ctx context.Context) ([]int64, error) {
	return r.engine.ListIncomplete(ctx)
}

func (r *Runner) Start(ctx context.Context, runnerContext *pipeline.RunnerContext) error {
	log := zerolog.Ctx(ctx)
	defer func() {
		if r := recover(); r != nil {
			log.Error().Any("Panic", r).Msg("Panic occured during execution")
		}
		log.Trace().Msg("runner: destroying engine artifacts")
		time.Sleep(5 * time.Second) //waiting for upload to finish
		r.engine.Destroy(noContext, runnerContext)
	}()

	// Calcualte resource limits based on context and config
	if runnerContext.ResourcesLimit == nil {
		runnerContext.ResourcesLimit = &pipeline.ResourcesLimit{
			CPU:    r.config.MaxCPU,
			Memory: r.config.MaxMemory,
		}
	} else {
		runnerContext.ResourcesLimit.CPU = helpers.Min(runnerContext.ResourcesLimit.CPU, r.config.MaxCPU)
		runnerContext.ResourcesLimit.Memory = helpers.Min(runnerContext.ResourcesLimit.Memory, r.config.MaxMemory)
	}

	log.Trace().Msg("runner: engine artifacts setup started")
	if err := r.engine.Setup(ctx, runnerContext); err != nil {
		log.Error().Err(err).Msg("runner: error setting up engine artifacts")
		runnerContext.Deployment.Fail(err)
		return r.report(ctx, runnerContext)
	}
	log.Trace().Msg("runner: engine artifacts setup completed")

	log.Trace().Msg("runner: starting steps execution")
	//start steps in sequence
	log.Trace().Int("steps.count", len(runnerContext.Steps)).Msg("runner: steps counts")

	runnerContext.Deployment.Start()
	if err := r.report(ctx, runnerContext); err != nil {
		return err
	}

	for _, step := range runnerContext.Steps {
		log := zerolog.Ctx(ctx).With().
			Str("step.name", step.Name).
			Logger()

		secrets := make([]string, len(runnerContext.Secrets))
		for i, s := range runnerContext.Secrets {
			secrets[i] = s.Data
		}

		state, err := r.runStep(ctx, runnerContext, step, secrets...)
		if err != nil {
			log.Trace().Err(err).Msg("runner: step run failed")
			runnerContext.Deployment.Fail(err)
			return r.report(ctx, runnerContext)
		}

		if state.ExitCode != 0 {
			runnerContext.Deployment.Fail(fmt.Errorf("deployment exited with non-zero exit code: %d", state.ExitCode))
			return r.report(ctx, runnerContext)
		}

		log.Trace().Msg("runner: step completed")
	}

	runnerContext.Deployment.Finish()
	return r.report(ctx, runnerContext)
}

func (r *Runner) runStep(ctx context.Context, runnerContext *pipeline.RunnerContext, step *pipeline.Step, secrets ...string) (*engine.State, error) {
	log := zerolog.Ctx(ctx)
	stepLogger := logger.New(r.client, runnerContext.Deployment.ID, secrets)
	log.Trace().Str("step.name", step.Name).Msg("runner: starting step")
	if err := r.engine.StartStep(ctx, runnerContext, step); err != nil {
		log.Error().Err(err).Msg("runner: error starting step")
		return nil, err
	}

	rc, err := r.engine.TailStep(ctx, runnerContext, step)
	if err != nil {
		log.Error().Err(err).Msg("engine: error tailing step")
		return nil, err
	}

	if err := logger.Log(log, stepLogger, rc); err != nil {
		log.Error().Err(err).Msg("runner: error logging step")
	}

	return r.engine.WaitStep(ctx, runnerContext, step)
}

func (r *Runner) report(ctx context.Context, runnerCtx *pipeline.RunnerContext) error {
	log := zerolog.Ctx(ctx)
	if ctx.Err() != nil {
		log.Trace().Err(ctx.Err()).Msg("runner: context error skipping reporting")
		return nil
	}
	err := r.client.Update(noContext, runnerCtx.Deployment)
	if err != nil {
		log.Error().Err(err).Msg("runner: error updating deployment")
		return err
	}
	return nil
}
