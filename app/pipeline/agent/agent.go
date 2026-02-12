package agent

import (
	"context"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cloudness-io/cloudness/app/pipeline/manager/client"
	"github.com/cloudness-io/cloudness/app/pipeline/runner/engine"
	"github.com/cloudness-io/cloudness/app/pipeline/runner/engine/kubernetes"
	"github.com/cloudness-io/cloudness/app/services/manager"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

var engines = []engine.Engine{
	kubernetes.New(),
}

type Agent struct {
	serverUID     int64
	client        client.RunnerClient
	serverManager manager.ServerManager
	stopChan      <-chan bool
	configUpdated int64
}

func New(client client.RunnerClient, factory manager.ManagerFactory) *Agent {
	serverManager, err := factory.GetServerManagerByType(enum.ServerTypeK8s)
	if err != nil {
		panic(err)
	}

	return &Agent{
		client:        client,
		serverManager: serverManager,
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

				//Build Runner
				if config.EnableRunner {
					go func() {
						if err := a.run(rCtx, *config); err != nil {
							log.Ctx(ctx).Error().Err(err).Msg("agent: error running runner")
							return
						}
					}()
				}

				//Metrics Scrapper
				// if config.EnableMetricsScrapper{
				go func() {
					if err := a.runMetricsScrapper(rCtx, *config); err != nil {
						log.Ctx(ctx).Error().Err(err).Msg("agent: error running metrics scrapper")
						return
					}
				}()

				go func() {
					if err := a.runStatusMonitor(rCtx, *config); err != nil {
						log.Ctx(ctx).Error().Err(err).Msg("agent: error running metrics scrapper")
						return
					}
				}()
				// }
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
