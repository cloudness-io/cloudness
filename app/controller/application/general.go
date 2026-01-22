package application

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

type UpdateIconInput struct {
	Icon string `json:"icon"`
}

type UpdateNameInput struct {
	Name string `json:"name"`
}

func (c *Controller) UpdateIcon(ctx context.Context, application *types.Application, in *UpdateIconInput) (*types.Application, error) {
	application.Spec.Icon = in.Icon
	err := application.UpdateSpecJSON()
	if err != nil {
		return nil, err
	}
	return c.applicationStore.UpdateSpec(ctx, application)
}

func (c *Controller) UpdateName(ctx context.Context, application *types.Application, in *UpdateNameInput) (*types.Application, error) {
	if err := c.sanitizeNameUpdateInput(in); err != nil {
		return nil, err
	}

	application.Name = in.Name
	application.Spec.Name = in.Name
	err := application.UpdateSpecJSON()
	if err != nil {
		return nil, err
	}
	return c.applicationStore.UpdateSpec(ctx, application)
}

func (c *Controller) sanitizeNameUpdateInput(in *UpdateNameInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.Name); err != nil {
		errors.AddValidationError("name", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}
