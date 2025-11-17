package user

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) GetDemoUser(ctx context.Context) (*types.DemoUser, error) {
	return &types.DemoUser{
		Email:    demoUserEmail,
		Password: demoUserPassword,
	}, nil
}
