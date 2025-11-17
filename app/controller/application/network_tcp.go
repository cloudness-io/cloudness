package application

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/types"
)

type TCPProxyInput struct {
	Port int `json:"tcpProxyPort,string"`
}

func (c *Controller) AddTCPProxy(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	in *TCPProxyInput,
) (*types.Application, error) {
	networkIn := application.Spec.ToNetorkInput()
	networkIn.TCPProxiesInput = &types.TCPProxiesInput{
		Port: in.Port,
	}

	appIn := &types.ApplicationInput{
		NetworkInput: networkIn,
	}

	dto, err := c.convertInputToDto(ctx, appIn, tenant, project, environment, application, session.Principal.DisplayName)
	if err != nil {
		return nil, err
	}

	return c.updateWithTx(ctx, dto)
}

func (c *Controller) DeleteTCPProxy(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
) (*types.Application, error) {
	networkIn := application.Spec.ToNetorkInput()
	networkIn.TCPProxiesInput = nil

	appIn := &types.ApplicationInput{
		NetworkInput: networkIn,
	}

	dto, err := c.convertInputToDto(ctx, appIn, tenant, project, environment, application, session.Principal.DisplayName)
	if err != nil {
		return nil, err
	}

	return c.updateWithTx(ctx, dto)
}
