package auth

import (
	"context"

	"github.com/cloudness-io/cloudness/app/token"
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *Controller) Login(ctx context.Context, instance *types.Instance, in *LoginInput) (*types.TokenResponse, error) {
	demoUser, _ := c.userCtrl.GetDemoUser(ctx)
	if instance.DemoUserEnabled == false && in.Email == demoUser.Email {
		return nil, usererror.Forbidden("Demo user is disabled")
	}
	user, err := c.userCtrl.FindUserByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(in.Password),
	)

	if err != nil {
		log.Debug().Err(err).
			Str("user_uid", user.UID).
			Msg("invalid password")

		return nil, errors.NotFound("User not found")
	}

	token, jwtToken, err := token.CreateUserSession(ctx, c.tokenStore, user)
	if err != nil {
		return nil, err
	}

	return &types.TokenResponse{Token: *token, AccessToken: jwtToken}, nil
}
