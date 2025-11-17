package application

import (
	"context"
	"maps"
	"slices"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) CreateFromTemplate(
	ctx context.Context,
	actor string,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	template *types.Template) ([]*types.Application, error) {

	restrictions := c.configSvc.GetTenantRestrictions(tenant)
	serviceAppMap := make(map[string]*types.Application)

	server, err := c.serverCtrl.Get(ctx)
	if err != nil {
		return nil, err
	}

	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		for _, service := range template.Spec.Services {
			spec := &types.ApplicationSpec{
				Icon:       service.Icon,
				Name:       service.Name,
				Build:      service.Build,
				Deploy:     service.Deploy,
				Networking: service.Networking,
				Volumes:    service.ToVolumeMounts(),
			}

			dto, err := c.convertSpecToDTO(ctx, spec, tenant, project, environment, actor)
			if err != nil {
				return err
			}

			if service.Networking.ServiceDomain != nil {
				fqdn, err := c.SuggestFQDN(ctx, dto.Application)
				if err != nil {
					return err
				}
				dto.Application.Domain = fqdn
				spec.Networking.ServiceDomain.Domain = fqdn
			}

			app, err := c.CreateWithoutTx(ctx, dto)
			if err != nil {
				return err
			}

			//volumes
			for _, volume := range service.Volumes {
				_, err := c.volumeCtrl.Create(ctx, tenant, project, environment, app, &types.VolumeCreateInput{
					Name:      volume.Name,
					MountPath: volume.MountPath,
					Size:      restrictions.MinVolumeSize,
					Server:    server,
				})
				if err != nil {
					return err
				}
			}
			serviceAppMap[service.Name] = app
		}

		return c.varCtrl.AddTemplateVariables(ctx, environment.ID, serviceAppMap, template)
	})
	if err != nil {
		return nil, err
	}

	return slices.Collect(maps.Values(serviceAppMap)), nil
}
