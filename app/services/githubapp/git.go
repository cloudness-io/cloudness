package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/go-github/v69/github"
)

func (c *Service) ListRepos(ctx context.Context, ghApp *types.GithubApp) (*github.ListRepositories, error) {
	ghClient, err := c.getGithubClient(ctx, ghApp)
	if err != nil {
		return nil, err
	}

	repos, _, err := ghClient.Apps.ListRepos(ctx, &github.ListOptions{})
	if err != nil {
		return nil, err
	}
	return repos, nil
}

func (c *Service) ListBranches(ctx context.Context, ghApp *types.GithubApp, owner string, repo string) ([]*github.Branch, error) {
	ghClient, err := c.getGithubClient(ctx, ghApp)
	if err != nil {
		return nil, err
	}

	branches, _, err := ghClient.Repositories.ListBranches(ctx, owner, repo, &github.BranchListOptions{})
	if err != nil {
		return nil, err
	}
	return branches, nil
}

func (c *Service) GetLatestCommit(ctx context.Context, ghApp *types.GithubApp, repoURL string, branch string) (*github.RepositoryCommit, error) {
	owner, repo, err := helpers.SplitGitRepoUrl(repoURL)
	if err != nil {
		return nil, err
	}

	ghClient, err := c.getGithubClient(ctx, ghApp)
	if err != nil {
		return nil, err
	}

	ghBranch, _, err := ghClient.Repositories.GetBranch(ctx, owner, repo, branch, 1)
	if err != nil {
		return nil, err
	}
	return ghBranch.Commit, nil
}

func (c *Service) GetCommitBySha(ctx context.Context, ghApp *types.GithubApp, repoURL string, sha string) (*github.RepositoryCommit, error) {
	owner, repo, err := helpers.SplitGitRepoUrl(repoURL)
	if err != nil {
		return nil, err
	}

	ghClient, err := c.getGithubClient(ctx, ghApp)
	if err != nil {
		return nil, err
	}

	commit, _, err := ghClient.Repositories.GetCommit(ctx, owner, repo, sha, &github.ListOptions{})
	return commit, err
}
