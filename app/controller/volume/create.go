package volume

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

func (c *Controller) Create(ctx context.Context, tenant *types.Tenant, project *types.Project, env *types.Environment, app *types.Application, in *types.VolumeCreateInput) (*types.Volume, error) {
	err := c.sanitizeCreateInput(in)
	if err != nil {
		return nil, err
	}

	restriction := c.configSvc.GetVolumeRestrictions(in.Server, tenant)

	//Validate restrictions
	if err := c.validateRestrictions(in, restriction); err != nil {
		return nil, err
	}

	now := time.Now().UTC().UnixMilli()
	volume := &types.Volume{
		Name:           in.Name,
		UID:            helpers.GenerateUID(),
		TenantID:       tenant.ID,
		ProjectID:      project.ID,
		EnvironmentID:  env.ID,
		EnvironmentUID: env.UID,
		ServerID:       in.Server.ID,
		ApplicaitonID:  &app.ID,
		MountPath:      in.MountPath,
		Size:           in.Size,
		Created:        now,
		Updated:        now,
	}

	return c.volumeStore.Create(ctx, volume)
}

func (c *Controller) sanitizeCreateInput(in *types.VolumeCreateInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.Name); err != nil {
		errors.AddValidationError("name", err)
	}
	if err := check.Directory(in.MountPath); err != nil {
		errors.AddValidationError("mountPath", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}

func (c *Controller) validateRestrictions(in *types.VolumeCreateInput, restrictions *types.VolumeRestriction) error {
	err := check.NewValidationErrors()
	if in.Size < restrictions.MinVolumeSize {
		err.AddValidationError("size", check.NewValidationErrorf("Volume is below min allowed limit of %dMiB", restrictions.MinVolumeSize))
	}
	if in.Size > restrictions.MaxVolumeSize {
		err.AddValidationError("size", check.NewValidationErrorf("Volume is above max allowed limit of %dMiB", restrictions.MaxVolumeSize))
	}

	if err.HasError() {
		return err
	}
	return nil
}
