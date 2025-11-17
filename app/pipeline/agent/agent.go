package agent

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline/manager/client"
	"github.com/cloudness-io/cloudness/app/pipeline/runner"
	"github.com/cloudness-io/cloudness/app/pipeline/runner/engine"
	"github.com/cloudness-io/cloudness/app/pipeline/runner/engine/kubernetes"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

var engines = []engine.Engine{
	kubernetes.New(),
}

type Agent struct {
	serverUID     int64
	client        client.RunnerClient
	stopChan      <-chan bool
	configUpdated int64
}

func New(client client.RunnerClient) *Agent {
	return &Agent{
		client: client,
	}
}

func (a *Agent) Start(ctx context.Context) error {
	serverUID, err := a.client.GetPrimaryServerUID(ctx) //NOTE: temporary untill agent is split, serverUID will be passed an env or config to independent agent
	if err != nil {
		return err
	}
	a.serverUID = serverUID

	configChan := make(chan *types.RunnerConfig)
	defer close(configChan)
	defer func() {
		if r := recover(); r != nil {
			stack := string(debug.Stack())
			log.Ctx(ctx).Error().
				Str("panic", fmt.Sprintf("[%T] agent panic: %v", r, r)).
				Msg(stack)
		}
	}()

	//runner context and cancelFunc
	var rCtx context.Context
	var rCtxCancel context.CancelFunc
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case config := <-configChan:
				a.configUpdated = config.UpdatedAt
				if a.configUpdated == 0 { //adding this so the agent does not restart on config without changes
					a.configUpdated = time.Now().UTC().UnixMilli()
				}
				if rCtx != nil && rCtxCancel != nil { //avoid panic during first run
					rCtxCancel() // cancels previous run
				}
				rCtx, rCtxCancel = context.WithCancel(ctx)
				if config.EnableRunner {
					go func() {
						if err := a.run(rCtx, *config); err != nil {
							log.Ctx(ctx).Error().Err(err).Msg("agent: error running runner")
							return
						}
					}()
				}
			}
		}
	}()

	if err := a.checkConfigChange(ctx, configChan); err != nil {
		return err
	}

	return nil
}

func (a *Agent) checkConfigChange(ctx context.Context, configChan chan<- *types.RunnerConfig) error {
	ticker := time.NewTicker(time.Duration(10) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			config, err := a.client.Config(ctx, a.serverUID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("agent: error fetching runner config")
				return err
			}
			if a.configUpdated < config.UpdatedAt || a.configUpdated == 0 {
				log.Ctx(ctx).Trace().Msg("agent: config updated, sending to channel for restart")
				configChan <- config
			}
		}
	}
}

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
