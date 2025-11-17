package auth

import (
	"context"

	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/token"
	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) Callback(ctx context.Context, authSetting *types.AuthSetting, code string) (*types.TokenResponse, error) {
	var userEmail string
	var err error
	switch authSetting.Provider {
	case enum.AuthProviderGithub:
		userEmail, err = c.githubCallback(ctx, authSetting, code)
	default:
		return nil, errors.BadRequest("Auth provider not implemented")
	}

	if err != nil {
		return nil, err
	}

	if userEmail == "" {
		return nil, errors.BadRequest("User email not found")
	}

	in := &user.CreateInput{
		Email:       userEmail,
		DisplayName: userEmail,
	}

	user, err := c.userCtrl.CreateNoAuth(ctx, in)
	if err != nil {
		return nil, err
	}

	token, jwtToken, err := token.CreateUserSession(ctx, c.tokenStore, user)
	if err != nil {
		return nil, err
	}

	return &types.TokenResponse{Token: *token, AccessToken: jwtToken}, nil
}
