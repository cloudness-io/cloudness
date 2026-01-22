package application

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

type UpdateIconInput struct {
	Icon string `json:"icon"`
}

func (c *Controller) UpdateIcon(ctx context.Context, application *types.Application, in *UpdateIconInput) (*types.Application, error) {
	application.Spec.Icon = in.Icon
	err := application.UpdateSpecJSON()
	if err != nil {
		return nil, err
	}
	return c.applicationStore.UpdateSpec(ctx, application)
}
