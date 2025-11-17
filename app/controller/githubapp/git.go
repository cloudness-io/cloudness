package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) ListRepos(ctx context.Context, ghApp *types.GithubApp) ([]*types.GitRepo, error) {
	repos, err := c.ghAppSvc.ListRepos(ctx, ghApp)
	if err != nil {
		return nil, err
	}

	dst := make([]*types.GitRepo, 0)
	for _, repo := range repos.Repositories {
		dst = append(dst, &types.GitRepo{
			Name:     repo.GetName(),
			FullName: repo.GetFullName(),
		})
	}
	return dst, nil
}

func (c *Controller) ListBranches(ctx context.Context, ghApp *types.GithubApp, owner string, repo string) ([]*types.GitBranch, error) {
	branches, err := c.ghAppSvc.ListBranches(ctx, ghApp, owner, repo)
	if err != nil {
		return nil, err
	}

	dst := make([]*types.GitBranch, 0)
	for _, branch := range branches {
		dst = append(dst, &types.GitBranch{
			Name: branch.GetName(),
		})
	}

	return dst, nil
}
