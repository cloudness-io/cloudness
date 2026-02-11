package tenant

import (
	"context"
	"errors"
	"time"

	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
	"github.com/cloudness-io/cloudness/types/enum"
)

type CreateTenantInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *Controller) CreateNoAuth(ctx context.Context, userID int64, in *CreateTenantInput, role enum.TenantRole) (*types.Tenant, error) {
	if err := c.sanitizeCreateInput(in); err != nil {
		return nil, err
	}

	now := time.Now().UTC().UnixMilli()
	defaults := c.configSvc.GetTenantDefaults()
	tenant := &types.Tenant{
		UID:                helpers.GenerateUID(),
		Name:               in.Name,
		Slug:               helpers.Slugify("t", in.Name),
		Description:        in.Description,
		AllowAdminToModify: defaults.DefaultAllowAdminToModify,
		MaxProjects:        defaults.DefaultMaxProjectsPerTenant,
		MaxApps:            defaults.DefaultMaxApplicationsPerTenant,
		MaxInstances:       defaults.DefaultMaxInstancesPerApplication,
		MaxCPUPerApp:       defaults.DefaultMaxCPUPerApplication,
		MaxMemoryPerApp:    defaults.DefaultMaxMemoryPerApplication,
		MaxVolumes:         defaults.DefaultMaxVolumeCount,
		MinVolumeSize:      defaults.DefaultMinVolumeSize,
		MaxVolumeSize:      defaults.DefaultMaxVolumeSize,
		CreateBy:           userID,
		Created:            now,
		Updated:            now,
	}

	tenant, err := c.tenantStore.Create(ctx, tenant)
	if err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			return nil, tenantAlreadyTaken()
		}

		return nil, err
	}

	//Add membership for current user to tenant
	membership := &types.TenantMembership{
		TenantMembershipKey: &types.TenantMembershipKey{
			TenantID:    tenant.ID,
			PrincipalID: userID,
			Role:        role,
		},
		CreatedBy: userID,
		Created:   now,
		Updated:   now,
	}

	err = c.tenantMembershipStore.Create(ctx, membership)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func tenantAlreadyTaken() error {
	err := check.NewValidationErrors()
	err.AddValidationError("Name", errors.New("name is taken"))
	return err
}
