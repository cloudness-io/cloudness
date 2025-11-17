package application

import (
	"context"

	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
)

// Default in
func (c *Controller) GetDefaultGitIn() *types.ApplicationInput {
	return c.specSvc.DefaultGitIn()
}

func (c *Controller) GetGitpublicIn(ctx context.Context, repoURL string) (*types.ApplicationInput, error) {
	repoURL = helpers.SanitizeGitUrl(repoURL)
	repo, err := c.gitPublicCtrl.FindRepo(ctx, repoURL)
	if err != nil {
		return nil, err
	}

	appIn := c.GetDefaultGitIn()
	appIn.RepoURL = repoURL
	appIn.Repo = repo.FullName
	return appIn, nil
}

func (c *Controller) GetRegistryIn(image string, networkIn *types.NetworkInput) *types.ApplicationInput {
	appIn := &types.ApplicationInput{
		RegistryInput: &types.RegistryInput{
			Image: image,
		},
		NetworkInput: networkIn,
	}

	return appIn
}
