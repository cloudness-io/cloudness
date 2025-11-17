package auth

import (
	"context"

	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/dchest/uniuri"
	"golang.org/x/oauth2"
)

func (c *Controller) GetRedirectUrl(ctx context.Context, authSetting *types.AuthSetting) (string, error) {
	switch authSetting.Provider {
	case enum.AuthProviderGithub:
		config := c.githubOAuth2Config(authSetting)
		return config.AuthCodeURL(uniuri.NewLen(16), oauth2.AccessTypeOnline), nil
	}
	return "", errors.BadRequest("Auth provider %s is not supported", authSetting.Provider)
}
