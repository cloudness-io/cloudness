package agent

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

func (a *Agent) runStatusMonitor(ctx context.Context, config types.RunnerConfig) error {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			status, err := a.serverManager.ListApplicationStatuses(ctx, nil)
			if err != nil {
				continue
			}

			if len(status) == 0 {
				log.Debug().Msg("agent: no application status found")
				continue
			}

			err = a.client.UploadAppStatus(ctx, status)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("agent: error updating application status in monitor")
				continue
			}
		}
	}
}
