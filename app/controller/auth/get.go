package auth

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) GetByProvider(ctx context.Context, provider enum.AuthProvider) (*types.AuthSetting, error) {
	authSetting, err := c.authSettingStore.FindByProvider(ctx, provider)
	if err != nil {
		if errors.Is(err, store.ErrResourceNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return authSetting, nil
}

func (c *Controller) GetDemoUser(ctx context.Context, instance *types.Instance) (*types.DemoUser, error) {
	demoUser := &types.DemoUser{}
	if instance.DemoUserEnabled {
		demoUser, _ = c.userCtrl.GetDemoUser(ctx)
	}
	return demoUser, nil
}
