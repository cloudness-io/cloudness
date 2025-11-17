package auth

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) Init(ctx context.Context) error {
	for _, authSetting := range defaultAuthSettings {
		if _, err := c.authSettingStore.Create(ctx, authSetting); err != nil {
			return err
		}
	}
	return nil
}

var defaultAuthSettings = []*types.AuthSetting{
	{
		Provider: enum.AuthProviderPassword,
		Enabled:  true,
	},
	{
		Provider: enum.AuthProviderGoogle,
		Enabled:  false,
	},
	{
		Provider: enum.AuthProviderGithub,
		Enabled:  false,
	},
	{
		Provider: enum.AuthProviderGitlab,
		Enabled:  false,
	},
}
