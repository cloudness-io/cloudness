package auth

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) ListEnabled(ctx context.Context) (map[enum.AuthProvider]*types.AuthSetting, error) {
	auths, err := c.authSettingStore.List(ctx)
	if err != nil {
		return nil, err
	}

	enabled := map[enum.AuthProvider]*types.AuthSetting{}
	for _, a := range auths {
		if a.Enabled {
			enabled[a.Provider] = a
		}
	}

	return enabled, nil
}

func (c *Controller) ListAllMap(ctx context.Context) (map[enum.AuthProvider]*types.AuthSetting, error) {
	auths, err := c.authSettingStore.List(ctx)
	if err != nil {
		return nil, err
	}

	dst := map[enum.AuthProvider]*types.AuthSetting{}
	for _, a := range auths {
		dst[a.Provider] = a
	}
	return dst, nil
}
