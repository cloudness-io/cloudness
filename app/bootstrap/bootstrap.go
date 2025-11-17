package bootstrap

import (
	"context"

	"github.com/cloudness-io/cloudness/app/controller/auth"
	"github.com/cloudness-io/cloudness/app/controller/instance"
	"github.com/cloudness-io/cloudness/app/controller/server"
	"github.com/cloudness-io/cloudness/app/controller/template"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

// Bootstrap is an abstraction of a function that bootstraps a system.
type Bootstrap func(context.Context) error

func System(
	config *types.Config,
	instanceCtrl *instance.Controller,
	serverCtrl *server.Controller,
	authCtrl *auth.Controller,
	userCtrl *user.Controller,
	templateCtrl *template.Controller,
) func(context.Context) error {
	return func(ctx context.Context) error {
		server, err := serverCtrl.Init(ctx)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to initialize server controller")
			return err
		}

		if _, err := instanceCtrl.Init(ctx, server); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to initialize instance controller")
			return err
		}

		if err := authCtrl.Init(ctx); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to initialize auth controller")
			return err
		}

		if err := userCtrl.Init(ctx); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to initialize user controller")
			return err
		}

		if err := templateCtrl.Init(ctx); err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("failed to initialize template controller")
			return err
		}

		return nil
	}
}
