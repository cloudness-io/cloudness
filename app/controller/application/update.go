package application

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) Update(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	in *types.ApplicationInput,
) (*types.Application, error) {

	dto, err := c.convertInputToDto(ctx, in, tenant, project, environment, application, session.Principal.DisplayName)
	if err != nil {
		return nil, err
	}

	return c.updateWithTx(ctx, dto)
}
