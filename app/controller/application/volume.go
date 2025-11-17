package application

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

func (c *Controller) AddVolume(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	volumeIn *types.VolumeCreateInput,
) (*types.Application, error) {

	//check mount paths
	volumes := c.specSvc.GetVolumeMounts(application.Spec)
	for _, volume := range volumes {
		if volume.MountPath == volumeIn.MountPath {
			errors := check.NewValidationErrors()
			errors.AddValidationError("mountPath", check.NewValidationErrorf("Application already has the same mount path"))
			return nil, errors
		}
	}

	server, err := c.serverCtrl.FindByID(ctx, application.ServerID)
	if err != nil {
		return nil, err
	}

	volumeIn.Server = server

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		//craete volume entity
		_, err := c.volumeCtrl.Create(ctx, tenant, project, environment, application, volumeIn)
		if err != nil {
			return err
		}
		volumes = append(volumes, &types.VolumeMounts{
			VolumeName: volumeIn.Name,
			VolumeSize: volumeIn.Size,
			MountPath:  volumeIn.MountPath,
		})
		in := &types.ApplicationInput{
			VolumeInput: &types.VolumeInput{
				Volumes: volumes,
			},
		}

		dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, application, session.Principal.DisplayName)
		if err != nil {
			return err
		}

		_, err = c.updateWithoutTx(ctx, dto)
		return err
	})

	return application, err
}

func (c *Controller) UpdateVolume(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	volume *types.Volume,
) (*types.Application, error) {

	err := c.tx.WithTx(ctx, func(ctx context.Context) error {
		//craete volume entity
		volume.ApplicaitonID = &application.ID
		_, err := c.volumeCtrl.Update(ctx, volume)
		if err != nil {
			return err
		}

		volumes, err := c.volumeCtrl.ListForApp(ctx, application)
		if err != nil {
			return err
		}
		volumeMounts := []*types.VolumeMounts{}
		for _, v := range volumes {
			volumeMounts = append(volumeMounts, &types.VolumeMounts{
				VolumeName: v.Name,
				VolumeSize: v.Size,
				MountPath:  v.MountPath,
			})
		}

		in := &types.ApplicationInput{
			VolumeInput: &types.VolumeInput{
				Volumes: volumeMounts,
			},
		}

		dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, application, session.Principal.DisplayName)
		if err != nil {
			return err
		}

		_, err = c.updateWithoutTx(ctx, dto)
		return err
	})
	if err != nil {
		return nil, err
	}

	return application, err
}

func (c *Controller) DetachVolume(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	volume *types.Volume,
) (*types.Application, error) {
	err := c.tx.WithTx(ctx, func(ctx context.Context) error {
		_, err := c.volumeCtrl.Detach(ctx, volume)
		if err != nil {
			return err
		}
		volumes, err := c.volumeCtrl.ListForApp(ctx, application)
		if err != nil {
			return err
		}
		volumeMounts := []*types.VolumeMounts{}
		for _, v := range volumes {
			volumeMounts = append(volumeMounts, &types.VolumeMounts{
				VolumeName: v.Name,
				VolumeSize: v.Size,
				MountPath:  v.MountPath,
			})
		}

		in := &types.ApplicationInput{
			VolumeInput: &types.VolumeInput{
				Volumes: volumeMounts,
			},
		}

		dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, application, session.Principal.DisplayName)
		if err != nil {
			return err
		}

		_, err = c.updateWithoutTx(ctx, dto)
		return err

	})
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (c *Controller) DeleteVolume(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	volume *types.Volume,
) (*types.Application, error) {

	err := c.tx.WithTx(ctx, func(ctx context.Context) error {
		//craete volume entity
		err := c.volumeCtrl.SoftDelete(ctx, volume)
		if err != nil {
			return err
		}

		volumes, err := c.volumeCtrl.ListForApp(ctx, application)
		if err != nil {
			return err
		}
		volumeMounts := []*types.VolumeMounts{}
		for _, v := range volumes {
			volumeMounts = append(volumeMounts, &types.VolumeMounts{
				VolumeName: v.Name,
				VolumeSize: v.Size,
				MountPath:  v.MountPath,
			})
		}

		in := &types.ApplicationInput{
			VolumeInput: &types.VolumeInput{
				Volumes: volumeMounts,
			},
		}

		dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, application, session.Principal.DisplayName)
		if err != nil {
			return err
		}

		_, err = c.updateWithoutTx(ctx, dto)
		return err
	})
	return application, err
}
