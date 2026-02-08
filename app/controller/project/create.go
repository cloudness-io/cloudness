package project

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/app/controller/environment"
	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type CreateProjectInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Controller) Create(ctx context.Context, session *auth.Session, tenant *types.Tenant, in *CreateProjectInput) (*types.Project, error) {
	if err := c.sanitizeCreateInput(in); err != nil {
		log.Ctx(ctx).Err(err).Msg("Error validating project create input")
		return nil, err
	}

	restrictions := c.configSvc.GetTenantRestrictions(tenant)
	projectFilter := &types.ProjectFilter{
		TenantID: &tenant.ID,
	}
	pCount, err := c.projectStore.Count(ctx, projectFilter)
	if err != nil {
		return nil, err
	}

	if pCount >= restrictions.MaxProjects {
		return nil, errors.PreconditionFailed("Max project limit reached")
	}

	userID := session.Principal.ID

	now := time.Now().UTC().UnixMilli()
	project := &types.Project{
		UID:         helpers.GenerateUID(),
		TenantID:    tenant.ID,
		Name:        in.Name,
		Description: in.Description,
		CreateBy:    userID,
		Created:     now,
		Updated:     now,
	}

	c.tx.WithTx(ctx, func(ctx context.Context) error {
		project, err = c.projectStore.Create(ctx, project)
		if err != nil {
			return err
		}

		if err := c.AddMember(ctx, session, tenant.ID, project.ID, &ProjectMembershipAddModel{
			Email: session.Principal.Email,
			Role:  enum.ProjectRoleOwner,
		}); err != nil {
			return err
		}
		_, err = c.envCtrl.Create(ctx, session, tenant, project, &environment.CreateEnvironmentInput{
			Name: "Production",
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return project, err
}

func (c *Controller) sanitizeCreateInput(in *CreateProjectInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.Name); err != nil {
		errors.AddValidationError("name", err)
	}
	if err := check.Description(in.Description); err != nil {
		errors.AddValidationError("description", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}
