package application

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
)

type UpdateDomainInput struct {
	Scheme    string `json:"scheme"`
	Domain    string `json:"domain"`
	Subdomain string `json:"subdomain"`
	Port      int    `json:"serviceDomainPort,string"`
}

func FromDomainSpec(spec *types.ApplicationSpec) *UpdateDomainInput {
	input := spec.ToServiceDomainInput()
	if input == nil {
		return nil
	}
	return &UpdateDomainInput{
		Scheme:    input.Scheme,
		Domain:    input.Domain,
		Subdomain: input.Subdomain,
		Port:      input.Port,
	}
}

func (c *Controller) GenerateDomain(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
) (*types.Application, error) {
	fqdn, err := c.SuggestFQDN(ctx, application)
	if err != nil {
		return nil, err
	}

	networkIn := application.Spec.ToNetorkInput()
	networkIn.ServiceDomainInput = &types.ServiceDomainInput{
		FQDN: fqdn,
		Port: application.Spec.Networking.ContainerPorts[0],
	}
	in := &types.ApplicationInput{
		NetworkInput: networkIn,
	}

	application.Domain = fqdn
	dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, application, session.Principal.DisplayName)
	if err != nil {
		return nil, err
	}

	return c.updateWithTx(ctx, dto)
}

func (c *Controller) DeleteDomain(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
) (*types.Application, error) {
	application.Domain = ""

	networkIn := application.Spec.ToNetorkInput()
	networkIn.ServiceDomainInput = nil
	in := &types.ApplicationInput{
		NetworkInput: networkIn,
	}

	dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, application, session.Principal.DisplayName)
	if err != nil {
		return nil, err
	}

	return c.updateWithTx(ctx, dto)
}

func (c *Controller) UpdateDomain(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	in *UpdateDomainInput,
) (*types.Application, error) {
	//TODO: validate input
	fqdn := helpers.GenerateFQDN(in.Scheme, helpers.Normalize(in.Subdomain), in.Domain)
	networkIn := application.Spec.ToNetorkInput()
	networkIn.ServiceDomainInput = &types.ServiceDomainInput{
		FQDN: fqdn,
		Port: in.Port,
	}
	appIn := &types.ApplicationInput{
		NetworkInput: networkIn,
	}

	application.Domain = fqdn
	dto, err := c.convertInputToDto(ctx, appIn, tenant, project, environment, application, session.Principal.DisplayName)
	if err != nil {
		return nil, err
	}

	return c.updateWithTx(ctx, dto)
}
