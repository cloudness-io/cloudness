package agent

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline/runner"
	"github.com/cloudness-io/cloudness/app/pipeline/runner/engine"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (a *Agent) run(ctx context.Context, config types.RunnerConfig) error {
	engine, err := engine.FindEngine(engines, config)
	if err != nil {
		log.Error().Err(err).Msg("agent: error finding engine")
		return err
	}

	_, err = engine.Load(ctx, config)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msgf("agent: error loading %s engine", engine.Type())
		return err
	}

	hostname := config.Hostname
	if len(hostname) == 0 {
		hostname, _ = os.Hostname()
	}

	log.Info().Msgf("agent: starting agent with %s engine", engine.Type())

	var wg sync.WaitGroup
	jobChan := make(chan *types.Deployment)
	defer close(jobChan)
	// number of runners based on job
	for i := range config.ParallelWorkers {
		wg.Go(func() {
			for {
				select {
				case <-ctx.Done():
					return
				case deployment := <-jobChan:
					log := log.With().
						Int64("agent.thread", i).
						Int64("deployment.id", deployment.ID).
						Int64("deployment.uid", deployment.UID).
						Logger()
					deployment.Machine = hostname
					if err := a.client.Accept(ctx, deployment); err != nil {
						log.Error().Err(err).Msg("agent: unable to accept deployment")
						continue
					}

					log.Trace().Msg("agent: deployment accepted for processing")

					runnerContext, err := a.client.Init(ctx, deployment.ID)
					if err != nil {
						log.Error().Err(err).Msg("agent: unable to initialize deployment")
						continue
					}
					log.Trace().Msg("agent: deployment initialized for processing")
					log = log.With().
						Int64("application.uid", runnerContext.ApplicationUID).
						Int64("deployment.uid", runnerContext.Deployment.UID).
						Str("runner.id", runnerContext.RunnerName).
						Logger()

					runnerCtx, runnerCancel := context.WithCancel(log.WithContext(ctx))
					runner := runner.NewRunner(a.client, config, engine)

					go func() {
						log.Debug().Msg("agent: watch for deployment cancellation")
						done, _ := a.client.Watch(ctx, deployment.ID)
						if done {
							log.Debug().Msg("agent: received cancellation event")
							runnerCancel()
						} else {
							log.Debug().Msg("agent: done listening for cancellation event")
						}
					}()

					runner.Start(runnerCtx, runnerContext)
				}
			}
		})
	}

	//listing incomplete
	runner := runner.NewRunner(a.client, config, engine)
	if incomplete, err := runner.ListPending(ctx); err != nil {
		log.Ctx(ctx).Error().Err(err).Any("Incomplete", incomplete).Msg("agent: error listing pending jobs, skipping...")
	} else {
		if len(incomplete) == 0 {
			log.Ctx(ctx).Info().Msg("agent: no pending jobs found")
		} else {
			log.Ctx(ctx).Info().Any("pending", len(incomplete)).Msg("agent: found pending jobs, TODO: yet to be done")
		}
	}

	a.startPolling(ctx, config, jobChan)
	wg.Wait()

	return nil
}

func (a *Agent) startPolling(ctx context.Context, config types.RunnerConfig, jobChan chan<- *types.Deployment) {
	ticker := time.NewTicker(time.Duration(config.PollingInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			deployment, err := a.client.Request(ctx)
			if err == context.Canceled || err == context.DeadlineExceeded {
				log.Trace().Err(err).Msg("agent: no deployment received")
				continue
			}
			if err != nil {
				log.Error().Err(err).Msg("agent: unable to request deployment")
				//TODO: max retries?
				continue
			}

			log.Trace().Msg("agent: deployment received")

			if deployment == nil || deployment.ID == 0 {
				continue
			}

			jobChan <- deployment
		}
	}
}
