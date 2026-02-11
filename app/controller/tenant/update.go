package tenant

import (
	"context"

	"github.com/cloudness-io/cloudness/app/request"
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"

	"github.com/rs/zerolog/log"
)

type TenantGeneralUpdateModel struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

func (c *Controller) UpdateGeneral(ctx context.Context, tenant *types.Tenant, in *TenantGeneralUpdateModel) (*types.Tenant, error) {
	if err := c.sanitizeUpdateInput(in); err != nil {
		return nil, err
	}

	tenant.Name = in.Name
	tenant.Description = in.Description

	return c.tenantStore.Update(ctx, tenant)
}

func (c *Controller) UpdateRestrictions(ctx context.Context, tenant *types.Tenant, in *types.TenantRestrictions) (*types.TenantRestrictions, error) {
	restrictions := c.GetRestrctions(ctx, tenant)
	if !restrictions.AllowAdminToModify && !request.IsSuperAdmin(ctx) {
		return nil, usererror.Forbidden("Not allowed")
	}

	//check if not super admin and allow admin flag modified
	if !request.IsSuperAdmin(ctx) && restrictions.AllowAdminToModify != in.AllowAdminToModify {
		return nil, usererror.Forbidden("Forbidden: Only super admin can modify allow admin flag")
	}

	if in.MaxVolumeSize < in.MinVolumeSize {
		log.Ctx(ctx).Debug().Any("in", in).Msg("Volume size")
		return nil, usererror.BadRequest("Max volume size must be greater than or equal to min volume size")
	}

	tenant.AllowAdminToModify = in.AllowAdminToModify
	tenant.MaxProjects = in.MaxProjects
	tenant.MaxApps = in.MaxApps
	tenant.MaxInstances = in.MaxInstances
	tenant.MaxCPUPerApp = in.MaxCPU
	tenant.MaxMemoryPerApp = in.MaxMemory
	tenant.MaxVolumes = in.MaxVolumes
	tenant.MinVolumeSize = in.MinVolumeSize
	tenant.MaxVolumeSize = in.MaxVolumeSize

	tenant, err := c.tenantStore.Update(ctx, tenant)
	if err != nil {
		return nil, err
	}
	return c.GetRestrctions(ctx, tenant), nil
}

func (c *Controller) sanitizeUpdateInput(in *TenantGeneralUpdateModel) error {
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
