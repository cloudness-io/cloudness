package application

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

type AppDeleteOption struct {
	Volume bool `json:"volume,string"`
}

func (c *Controller) softDeleteInternal(ctx context.Context, app *types.Application, now int64) error {
	if err := c.canceler.CancelIncompleteBuilds(ctx, app.TenantID, app.ProjectID, app.ID); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error cancelling incomplete builds, but ignoring this to continue deletion")
	}

	temp := app.GetIdentifierStr()
	app.Domain = temp
	app.CustomDomain = temp
	app.PrivateDomain = temp
	return c.applicationStore.SoftDelete(ctx, app, now)
}

func (c *Controller) SoftDelete(ctx context.Context, app *types.Application, opts *AppDeleteOption) error {
	now := time.Now().UTC().UnixMilli()
	server, err := c.serverCtrl.FindByID(ctx, app.ServerID)
	if err != nil {
		return err
	}

	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		volumes, err := c.volumeCtrl.ListForApp(ctx, app)
		if err != nil {
			return err
		}
		for _, v := range volumes {
			if opts.Volume {
				if err := c.volumeCtrl.SoftDelete(ctx, v); err != nil {
					return err
				}
			} else {
				v.ApplicaitonID = nil
				if _, err := c.volumeCtrl.Update(ctx, v); err != nil {
					return err
				}
			}
		}

		if err := c.softDeleteInternal(ctx, app, now); err != nil {
			return err
		}

		if err := c.clearResources(ctx, server, app); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting resources, supressing error for now, will be tried during cleanup job")
		}
		return nil
	})
}

func (c *Controller) SoftDeleteInEnvironment(ctx context.Context, envID, now int64) error {
	apps, err := c.applicationStore.List(ctx, &types.ApplicationFilter{
		EnvironmentID: &envID,
	})
	if err != nil {
		return err
	}

	serverMap := make(map[int64]*types.Server)

	for _, app := range apps {
		if err := c.softDeleteInternal(ctx, app, now); err != nil {
			return err
		}
		if err := c.volumeCtrl.SoftDeleteInApplication(ctx, app.ID, now); err != nil {
			return err
		}
		server := serverMap[app.ServerID]
		var err error
		if server == nil {
			server, err = c.serverCtrl.FindByID(ctx, app.ServerID)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("Error finding server during deletion for resource cleanup, skipping...")
				continue
			}
			serverMap[app.ServerID] = server
		}
		if err := c.clearResources(ctx, server, app); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("Error deleting resources, supressing error for now, will be tried during cleanup job")
		}
	}

	return nil
}

func (c *Controller) clearResources(ctx context.Context, server *types.Server, app *types.Application) error {
	manager, err := c.manager.GetServerManager(server)
	if err != nil {
		return err
	}

	return manager.DeleteResources(ctx, server, app.Namespace(), app.GetIdentifierStr())
}
