package application

import (
	"context"

	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) CreateGithub(
	ctx context.Context,
	actor string,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	app *types.Application,
	ghApp *types.GithubApp,
	in *types.ApplicationInput,
) (*types.Application, error) {
	in.RepoURL = ghApp.GetHttpUrl(in.Repo)

	dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, app, actor)
	if err != nil {
		return nil, err
	}

	// Updating github app specific application context
	dto.Application.GithubAppID = &ghApp.ID

	return c.createWithTx(ctx, dto)
}

func (c *Controller) CreateRegistry(
	ctx context.Context,
	actor string,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	in *types.RegistryInput) (*types.Application, error) {
	//TODO: validation of image
	appIn := c.GetRegistryIn(in.Image, in.NetworkInput)
	dto, err := c.convertInputToDto(ctx, appIn, tenant, project, environment, nil, actor)
	if err != nil {
		return nil, err
	}

	return c.createWithTx(ctx, dto)
}

func (c *Controller) Create(
	ctx context.Context,
	actor string,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	in *types.ApplicationInput) (*types.Application, error) {
	dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, nil, actor)
	if err != nil {
		return nil, err
	}

	return c.createWithTx(ctx, dto)
}

func (c *Controller) CreateWithoutTx(ctx context.Context, dto *createOrUpdateDto) (*types.Application, error) {
	application := dto.Application
	if err := c.schemaSvc.ValidateApplication(ctx, application.Spec); err != nil {
		return nil, err
	}

	//Validate restrictions
	if err := c.validateRestrictions(dto); err != nil {
		return nil, err
	}

	// update application spec json
	if err := application.UpdateSpecJSON(); err != nil {
		return nil, err
	}

	restrictions := c.configSvc.GetTenantRestrictions(dto.Tenant)
	aCount, err := c.applicationStore.Count(ctx, &types.ApplicationFilter{
		TenantID: &application.TenantID,
	})
	if err != nil {
		return nil, err
	}

	if aCount >= restrictions.MaxApps {
		return nil, errors.PreconditionFailed("Max applications limit reached")
	}
	// update server
	server, err := c.serverCtrl.Get(ctx)
	if err != nil {
		return nil, err
	}
	privateDomain, err := c.suggestPrivateDomain(ctx, application)
	if err != nil {
		return nil, err
	}

	application.ServerID = server.ID
	application.DeploymentStatus = enum.ApplicationDeploymentStatusNeedsDeployment
	application.PrivateDomain = privateDomain
	application.Name = helpers.CapitalizeString(privateDomain)
	application.Slug = helpers.Slugify("", application.Name)
	application.ParentSlug = dto.Environment.Slug

	application, err = c.applicationStore.Create(ctx, application)
	if err != nil {
		return nil, err
	}
	err = c.varCtrl.UpdateDefaultVariables(ctx, server, dto.Tenant, dto.Project, dto.Environment, application)
	if err != nil {
		return nil, err
	}

	return application, nil
}

func (c *Controller) createWithTx(ctx context.Context, dto *createOrUpdateDto) (*types.Application, error) {
	var app *types.Application
	var err error
	err = c.tx.WithTx(ctx, func(ctx context.Context) error {
		app, err = c.CreateWithoutTx(ctx, dto)
		return err
	})
	if err != nil {
		return nil, err
	}
	return app, nil
}
