package agent

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (a *Agent) runMetricsScrapper(ctx context.Context, config types.RunnerConfig) error {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			metrics, err := a.serverManager.ListMetrics(ctx, nil)
			if err != nil {
				continue
			}

			if len(metrics) == 0 {
				log.Debug().Msg("scrapper: no metrics found")
				continue
			}

			err = a.client.UploadMetrics(ctx, metrics)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("agent: error updating metrics in scrapper")
				continue
			}
		}
	}
}
