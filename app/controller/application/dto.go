package application

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

type createOrUpdateDto struct {
	Tenant      *types.Tenant
	Project     *types.Project
	Environment *types.Environment
	Application *types.Application
	Actor       string
}

func (d *Controller) convertInputToDto(
	ctx context.Context,
	in *types.ApplicationInput,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	actor string,
) (*createOrUpdateDto, error) {
	application, err := d.specSvc.ToApplication(ctx, in, application)
	if err != nil {
		return nil, err
	}

	application.TenantID = tenant.ID
	application.ProjectID = project.ID
	application.EnvironmentID = environment.ID
	application.EnvironmentUID = environment.UID

	return &createOrUpdateDto{
		Tenant:      tenant,
		Project:     project,
		Environment: environment,
		Application: application,
		Actor:       actor,
	}, nil
}

func (c *Controller) convertSpecToDTO(
	ctx context.Context,
	spec *types.ApplicationSpec,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	actor string,
) (*createOrUpdateDto, error) {
	return c.convertInputToDto(ctx, spec.ToInput(), tenant, project, environment, nil, actor)
}
