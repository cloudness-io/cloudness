package githubapp

import (
	"github.com/cloudness-io/cloudness/app/services/githubapp"
	"github.com/cloudness-io/cloudness/types/check"
)

const (
	githubApiUrl  = "https://api.github.com"
	githubHtmlUrl = "https://github.com"
	customUser    = "git"
	customPort    = 22
)

type Controller struct {
	ghAppSvc *githubapp.Service
}

func NewController(ghAppSvc *githubapp.Service) *Controller {
	return &Controller{
		ghAppSvc: ghAppSvc,
	}
}

// func (c *Controller) findByUIDNoAuth(ctx context.Context, githubAppUID int64) (*types.GithubApp, error) {
// 	githubApp, err := c.githubAppStore.FindByUIDNoAuth(ctx, githubAppUID)
// 	if err != nil && !errors.Is(err, dbStore.ErrResourceNotFound) {
// 		return nil, err
// 	}
// 	return githubApp, nil
// }

func (c *Controller) sanitizeCreateInput(in *CreateGithubAppInput) error {
	errors := check.NewValidationErrors()
	if err := check.DisplayName(in.Name); err != nil {
		errors.AddValidationError("name", err)
	}
	if errors.HasError() {
		return errors
	}

	if in.ApiUrl == "" {
		in.ApiUrl = githubApiUrl
	}
	if in.HtmlUrl == "" {
		in.HtmlUrl = githubHtmlUrl
	}
	if in.CustomUser == "" {
		in.CustomUser = customUser
	}
	if in.CustomPort == 0 {
		in.CustomPort = customPort
	}

	return nil
}
