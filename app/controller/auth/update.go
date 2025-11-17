package auth

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

// TODO: looks clumsy, check if you can reformat
func (c *Controller) UpdatePasswordSettings(ctx context.Context, in *ChangePasswordSettings) error {
	return c.tx.WithTx(ctx, func(ctx context.Context) error {
		if err := c.instanceCtrl.UpdateEnableSignup(ctx, in.EnableRegistration); err != nil {
			return err
		}

		authSetting, err := c.authSettingStore.FindByProvider(ctx, enum.AuthProviderPassword)
		if err != nil {
			return err
		}
		if authSetting.Enabled != in.EnablePasswordLogin {
			authSetting.Enabled = in.EnablePasswordLogin
			if _, err := c.authSettingStore.Update(ctx, authSetting); err != nil {
				return err
			}
		}
		return nil
	})
}

func (c *Controller) UpdateDemoUserAuth(ctx context.Context, in *DemoUserSettings) error {
	return c.instanceCtrl.UpdateDemoUser(ctx, in.DemoUserEnabled)
}

func (c *Controller) UpdateOauthSettings(ctx context.Context, auth *types.AuthSetting, provider enum.AuthProvider) error {
	auth.Provider = provider
	log.Ctx(ctx).Debug().Any("auth", auth).Msg("Auth settings")
	_, err := c.authSettingStore.Update(ctx, auth)
	return err
}
