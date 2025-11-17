package triggerer

import (
	"context"
	"runtime/debug"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline/canceler"
	"github.com/cloudness-io/cloudness/app/pipeline/scheduler"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type TriggerHook struct {
	ApplicaitonID int64              `json:"application_id"`
	Triggerer     string             `json:"triggerer"`
	Title         string             `json:"title"`
	Action        enum.TriggerAction `josn:"action"`
}

type Triggerer interface {
	Trigger(ctx context.Context, hook *TriggerHook) (*types.Deployment, error)
}

type triggerer struct {
	tx              dbtx.Transactor
	appStore        store.ApplicationStore
	deploymentStore store.DeploymentStore
	scheduler       scheduler.Scheduler
	canceler        canceler.Canceler
}

func New(
	tx dbtx.Transactor,
	appStore store.ApplicationStore,
	deploymentStore store.DeploymentStore,
	scheduler scheduler.Scheduler,
	canceler canceler.Canceler,
) Triggerer {
	return &triggerer{
		tx:              tx,
		appStore:        appStore,
		deploymentStore: deploymentStore,
		scheduler:       scheduler,
		canceler:        canceler,
	}
}

var _ Triggerer = (*triggerer)(nil)

func (t *triggerer) Trigger(ctx context.Context, hook *TriggerHook) (*types.Deployment, error) {
	log := log.With().
		Int64("application.id", hook.ApplicaitonID).
		Logger()

	log.Debug().Msg("trigger: trigger received")
	defer func() {
		//in ideal case this should never happen
		if r := recover(); r != nil {
			log.Error().Msgf("trigger: unexpected panic: %s", r)
			debug.PrintStack()
		}
	}()

	now := time.Now().UTC().UnixMilli()
	application, err := t.appStore.Find(ctx, hook.ApplicaitonID)
	if err != nil {
		return nil, err
	}

	deployment := &types.Deployment{
		UID:           helpers.GenerateUID(),
		ApplicationID: application.ID,
		SpecJson:      application.SpecJSON,
		//TODO Determined needs build based on previous spec
		NeedsBuild: true,
		Triggerer:  hook.Triggerer,
		Title:      hook.Title,
		Action:     hook.Action,
		Status:     enum.DeploymentStatusPending,
		Error:      "",
		Version:    0,
		Created:    now,
		Updated:    now,
	}

	err = t.tx.WithTx(ctx, func(ctx context.Context) error {
		// cancelling other deployments before creating new
		if err := t.canceler.CancelIncompleteBuilds(ctx, application.TenantID, application.ProjectID, application.ID); err != nil {
			log.Warn().Err(err).Msg("trigger: failed to cancel incomplete builds")
		}

		deployment, err = t.deploymentStore.Create(ctx, deployment)
		if err != nil {
			return err
		}
		log = log.With().
			Int64("deployment.id", deployment.ID).
			Logger()

		log.Debug().Msg("trigger: deployment created")

		_, err = t.appStore.UpdateDeploymentTriggerTime(ctx, application)
		if err != nil {
			return err
		}

		err = t.scheduler.Schedule(ctx, deployment)
		if err != nil {
			//TODO: should we error this out?
			return nil
		}

		log.Debug().Msg("trigger: deployment scheduled")

		return nil
	})

	if err != nil {
		return nil, err
	}

	return deployment, nil
}
