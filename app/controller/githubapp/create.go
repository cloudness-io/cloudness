package githubapp

import (
	"context"
	"time"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
)

type CreateGithubAppInput struct {
	Name         string `json:"name"`
	Organization string `json:"organization"`
	IsTenantWide bool   `json:"isTenantWide,string"`
	ApiUrl       string `json:"apiUrl"`
	HtmlUrl      string `json:"htmlUrl"`
	CustomUser   string `json:"customUser"`
	CustomPort   int64  `json:"customPort"`
}

func (c *Controller) Create(ctx context.Context, session *auth.Session, tenantID, projectID int64, in *CreateGithubAppInput) (*types.GithubApp, error) {
	//TODO: permission check, only admin?
	if err := c.sanitizeCreateInput(in); err != nil {
		log.Ctx(ctx).Err(err).Msg("Error validating project create input")
		return nil, err
	}

	now := time.Now().UTC().UnixMilli()

	ghApp := &types.GithubApp{
		UID:          helpers.GenerateUID(),
		TenantID:     tenantID,
		Name:         in.Name,
		Organization: in.Organization,
		IsTenantWide: in.IsTenantWide,
		ApiUrl:       in.ApiUrl,
		HtmlUrl:      in.HtmlUrl,
		CustomUser:   in.CustomUser,
		CustomPort:   in.CustomPort,

		CreateBy: session.Principal.ID,
		Created:  now,
		Updated:  now,
	}

	if projectID > 0 && !in.IsTenantWide {
		ghApp.ProjectID = projectID
	}

	return c.ghAppSvc.Create(ctx, ghApp)
}
