package auth

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/app/controller/tenant"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/app/token"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
	"k8s.io/apimachinery/pkg/api/errors"
)

type RegisterUserInput struct {
	DisplayName     string `json:"displayName"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

func (c *Controller) Register(ctx context.Context, in *RegisterUserInput) (*types.TokenResponse, error) {
	user, err := c.registerUser(ctx, in)
	if err != nil {
		return nil, err
	}

	token, jwtToken, err := token.CreateUserSession(ctx, c.tokenStore, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create token after successful user creation: %w", err)
	}
	return &types.TokenResponse{Token: *token, AccessToken: jwtToken}, nil
}

func (c *Controller) registerUser(ctx context.Context, in *RegisterUserInput) (*types.User, error) {
	allowed, err := c.instanceCtrl.IsUserSignupAllowed(ctx)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, errors.NewBadRequest("User registration is disabled")
	}

	if err := c.sanitizeRegisterInput(in); err != nil {
		return nil, err
	}

	createIn := &user.CreateInput{
		Email:       in.Email,
		DisplayName: in.DisplayName,
		Password:    in.Password,
	}

	user, err := c.userCtrl.CreateNoAuth(ctx, createIn)
	if err != nil {
		return nil, err
	}

	uCount, err := c.userCtrl.CountTotal(ctx)
	if err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("Error counting users")
		return user, nil //skipping tenant creation
	}

	if uCount == 1 {
		//create personal tenant for the user
		tenantIn := tenant.CreateTenantInput{
			Name:        "Root",
			Description: "Default",
		}
		_, err = c.tenantCtrl.CreateNoAuth(ctx, user.ID, &tenantIn, enum.TenantRoleAdmin)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("unable create tenant")
			return user, nil
		}

		// update the first user to super admin
		err = c.instanceCtrl.UpdateSuperAdmin(ctx, &user.ID)
		if err != nil {
			return nil, err
		}
	}

	return user, err
}

func (c *Controller) sanitizeRegisterInput(in *RegisterUserInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.DisplayName); err != nil {
		errors.AddValidationError("displayName", err)
	}
	if err := check.Email(in.Email); err != nil {
		errors.AddValidationError("email", err)
	}
	if err := check.Password(in.Password); err != nil {
		errors.AddValidationError("password", err)
	}

	if errors.HasError() {
		return errors
	}
	return nil

}
