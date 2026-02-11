package environment

import (
	"context"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

type CreateEnvironmentInput struct {
	Name string `json:"name"`
}

func (c *Controller) Create(ctx context.Context, session *auth.Session, tenant *types.Tenant, project *types.Project, in *CreateEnvironmentInput) (*types.Environment, error) {
	//TODO: permission check, only admin?
	if err := c.sanitizeCreateInput(in); err != nil {
		log.Ctx(ctx).Err(err).Msg("Error validating project create input")
		return nil, err
	}

	envs, err := c.List(ctx, tenant.ID, project.ID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC().UnixMilli()
	env := &types.Environment{
		UID:       helpers.GenerateUID(),
		TenantID:  tenant.ID,
		ProjectID: project.ID,
		Seq:       int64(len(envs) + 1),
		Name:      in.Name,
		Slug:      helpers.Slugify("e-"+strings.TrimPrefix(project.Slug, "p-"), in.Name),
		CreateBy:  session.Principal.ID,
		Created:   now,
		Updated:   now,
	}

	return c.environmentStore.Create(ctx, env)
}
