package user

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

// UpdateInput store infos to update an existing user.
type UpdateInput struct {
	DisplayName string `json:"name"`
}

// Update updates the provided user.
func (c *Controller) Update(ctx context.Context, userID int64, in *UpdateInput) (*types.User, error) {
	user, err := findUserFromID(ctx, c.principalStore, userID)
	if err != nil {
		return nil, err
	}

	if err = c.sanitizeUpdateInput(in); err != nil {
		return nil, err
	}

	user.DisplayName = in.DisplayName

	err = c.principalStore.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Controller) sanitizeUpdateInput(in *UpdateInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.DisplayName); err != nil {
		errors.AddValidationError("name", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}
