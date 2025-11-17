package project

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"
)

type ProjectMembershipAddModel struct {
	Email string           `json:"new_email"`
	Role  enum.ProjectRole `json:"role"`
}

type ProjectMembershipUpdateModel struct {
	Email string           `json:"email"`
	Role  enum.ProjectRole `json:"role"`
}

type ProjectMembershipRemoveModel struct {
	Email string `json:"email"`
}

func (c *Controller) AddMember(ctx context.Context, session *auth.Session, tenantID, projectID int64, in *ProjectMembershipAddModel) error {
	if err := c.sanitizeCreateMembershipInput(in); err != nil {
		return err
	}

	user, err := c.userCtrl.FindUserByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	tenantMembership, err := c.tenantMembershipStore.Find(ctx, tenantID, user.ID)
	if err != nil {
		return err
	}

	membership := &types.ProjectMembership{
		ProjectMembershipKey: &types.ProjectMembershipKey{
			TenantID:           tenantID,
			TenantMembershipID: tenantMembership.ID,
			ProjectID:          projectID,
			PrincipalID:        user.ID,
			Role:               in.Role,
		},
		CreatedBy: session.Principal.ID,
	}
	return c.projectMembershipStore.Create(ctx, membership)
}

func (c *Controller) UpdateMember(ctx context.Context, tenantID, projetID int64, in *ProjectMembershipUpdateModel) error {
	if err := c.sanitizeUpdateMembershipInput(in); err != nil {
		return err
	}

	user, err := c.userCtrl.FindUserByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	return c.projectMembershipStore.Update(ctx, tenantID, projetID, user.ID, in.Role)
}

func (c *Controller) RemoveMember(ctx context.Context, tenantID, projectID int64, in *ProjectMembershipRemoveModel) error {
	if err := c.sanitizeRemoveMembershipInput(in); err != nil {
		return err
	}

	user, err := c.userCtrl.FindUserByEmail(ctx, in.Email)
	if err != nil {
		return err
	}

	return c.projectMembershipStore.Delete(ctx, tenantID, projectID, user.ID)
}

func (c *Controller) FindMembership(ctx context.Context, tenantID, projectID int64, principalID int64) (*types.ProjectMembership, error) {
	membership, err := c.projectMembershipStore.Find(ctx, tenantID, projectID, principalID)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, err
	}
	return membership, nil
}

func (c *Controller) ListMembers(ctx context.Context, tenantID, projectID int64) ([]*types.ProjectMembershipUser, error) {
	return c.projectMembershipStore.List(ctx, tenantID, projectID)
}

func (c *Controller) sanitizeCreateMembershipInput(in *ProjectMembershipAddModel) error {
	errors := check.NewValidationErrors()
	if err := check.Email(in.Email); err != nil {
		errors.AddValidationError("new_email", err)
	}

	role := enum.ProjectRoleFromString(string(in.Role))
	if role == "" {
		errors.AddValidationError("role", check.NewValidationError("Invalid role"))
	}

	if errors.HasError() {
		return errors
	}
	return nil
}

func (c *Controller) sanitizeUpdateMembershipInput(in *ProjectMembershipUpdateModel) error {
	errors := check.NewValidationErrors()
	if err := check.Email(in.Email); err != nil {
		errors.AddValidationError("email", err)
	}

	role := enum.ProjectRoleFromString(string(in.Role))
	if role == "" {
		errors.AddValidationError("role", check.NewValidationError("Invalid role"))
	}

	if errors.HasError() {
		return errors
	}
	return nil
}

func (c *Controller) sanitizeRemoveMembershipInput(in *ProjectMembershipRemoveModel) error {
	errors := check.NewValidationErrors()
	if err := check.Email(in.Email); err != nil {
		errors.AddValidationError("email", err)
	}

	if errors.HasError() {
		return errors
	}
	return nil
}
