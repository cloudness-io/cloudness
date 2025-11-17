package tenant

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/app/controller/user"
	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"
)

type TenantMembershipModel struct {
	Email string          `json:"email"`
	Role  enum.TenantRole `json:"role"`
}

func (c *Controller) FindMembership(ctx context.Context, tenantID int64, principalID int64) (*types.TenantMembership, error) {
	membership, err := c.tenantMembershipStore.Find(ctx, tenantID, principalID)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, err
	}
	return membership, nil
}

func (c *Controller) CreateTenantMembership(ctx context.Context, tenant *types.Tenant, session *auth.Session, in *TenantMembershipModel) error {
	if err := c.sanitizeCreateMembershipInput(in); err != nil {
		return err
	}

	principal, err := c.userCtrl.FindUserByEmail(ctx, in.Email)
	if err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	if principal == nil {
		userIn := &user.CreateInput{
			DisplayName: in.Email,
			Email:       in.Email,
		}

		principal, err = c.userCtrl.CreateNoAuth(ctx, userIn)
		if err != nil {
			return err
		}
	}

	membership := &types.TenantMembership{
		TenantMembershipKey: &types.TenantMembershipKey{
			TenantID:    tenant.ID,
			PrincipalID: principal.ID,
			Role:        in.Role,
		},
		CreatedBy: session.Principal.ID,
	}

	return c.tenantMembershipStore.Create(ctx, membership)
}

func (c *Controller) sanitizeCreateMembershipInput(in *TenantMembershipModel) error {
	errors := check.NewValidationErrors()
	if err := check.Email(in.Email); err != nil {
		errors.AddValidationError("email", err)
	}

	role := enum.TenantRoleFromString(string(in.Role))
	if role == "" {
		errors.AddValidationError("role", check.NewValidationError("Invalid role"))
	}

	if errors.HasError() {
		return errors
	}
	return nil
}

func (c *Controller) UpdateMembership(ctx context.Context, tenantID int64, in *TenantMembershipModel) error {
	user, err := c.userCtrl.FindUserByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	role := enum.TenantRoleFromString(string(in.Role))
	if role == "" {
		return errors.BadRequest("Invalid role")
	}

	return c.tenantMembershipStore.Update(ctx, tenantID, user.ID, role)
}

func (c *Controller) DeleteMembership(ctx context.Context, tenantID int64, in *TenantMembershipModel) error {
	user, err := c.userCtrl.FindUserByEmail(ctx, in.Email)
	if err != nil {
		return nil
	}

	return c.tenantMembershipStore.Delete(ctx, tenantID, user.ID)
}
