package application

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
)

type PrivateNetworkInput struct {
	Subdomain string `json:"privateSubdomain"`
}

func PrivateNetworkFromApp(app *types.Application) *PrivateNetworkInput {
	return &PrivateNetworkInput{
		Subdomain: app.PrivateDomain,
	}
}

func (c *Controller) UpdatePrivateNetwork(
	ctx context.Context,
	session *auth.Session,
	tenant *types.Tenant,
	project *types.Project,
	environment *types.Environment,
	application *types.Application,
	in *PrivateNetworkInput,
) (*types.Application, error) {
	if err := c.sanitizePrivateNetworkInput(in); err != nil {
		return nil, err
	}
	application.PrivateDomain = helpers.Normalize(in.Subdomain)

	dto, err := c.convertInputToDto(ctx, application.Spec.ToInput(), tenant, project, environment, application, session.Principal.DisplayName)
	if err != nil {
		return nil, err
	}

	application, err = c.updateWithTx(ctx, dto)
	if err != nil {
		if errors.Is(err, store.ErrDuplicate) {
			return nil, errors.BadRequest("Private domain is already taken")
		}
		return nil, err
	}

	return application, nil
}

func (c *Controller) sanitizePrivateNetworkInput(in *PrivateNetworkInput) error {
	if in.Subdomain == "" {
		return errors.BadRequest("Private network subdomain is required")
	}

	return nil
}
