package gitpublic

import (
	"context"
	"strings"

	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/helpers"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"

	"github.com/drone/go-scm/scm"
	"github.com/drone/go-scm/scm/driver/github"
)

func (c *Service) getScmClient(repoURL string) (*scm.Client, error) {
	switch true {
	case strings.HasPrefix(repoURL, "https://github.com"):
		return github.NewDefault(), nil
		// case strings.HasPrefix(repoURL, "https://gitlab.com"):
		// 	return gitlab.NewDefault(), nil
		// case strings.HasPrefix(repoURL, "https://bitbucket.org"):
		// 	return bitbucket.NewDefault(), nil
	}

	return nil, errors.BadRequest("Scm client not configured for this public repo")
}

func (c *Service) FindRepo(ctx context.Context, repoURL string) (*types.GitRepo, error) {
	owner, repo, err := helpers.SplitGitRepoUrl(repoURL)
	if err != nil {
		return nil, err
	}

	client, err := c.getScmClient(repoURL)
	if err != nil {
		return nil, err
	}

	scmRepo, _, err := client.Repositories.Find(ctx, owner+"/"+repo)
	if err != nil {
		if errors.As(err, &scm.ErrNotFound) {
			return nil, check.NewValidationErrorsKey("repoURL", "Git repository not found")
		}
		return nil, err
	}
	return &types.GitRepo{
		Name:     scmRepo.Name,
		FullName: scmRepo.Namespace + "/" + scmRepo.Name,
	}, nil
}

func (c *Service) ListBranches(ctx context.Context, repoURL string) ([]*types.GitBranch, error) {
	owner, repo, err := helpers.SplitGitRepoUrl(repoURL)
	if err != nil {
		return nil, err
	}

	client, err := c.getScmClient(repoURL)
	if err != nil {
		return nil, err
	}

	branches := []*types.GitBranch{}
	for {
		page := 1
		size := 100
		gitBranches, response, err := client.Git.ListBranches(ctx, owner+"/"+repo, scm.ListOptions{Page: page, Size: size})
		if err != nil {
			if errors.As(err, &scm.ErrNotFound) {
				return nil, check.NewValidationErrorsKey("repoURL", "Git repository not found")
			}
			return nil, err
		}

		for _, branch := range gitBranches {
			branches = append(branches, &types.GitBranch{Name: branch.Name})
		}

		page = response.Page.Next
		if page == 0 && response.Page.NextURL == "" {
			break
		}
	}

	return branches, nil
}

func (c *Service) GetLatestCommit(ctx context.Context, repoURL string, branch string) (*types.GitCommit, error) {
	owner, repo, err := helpers.SplitGitRepoUrl(repoURL)
	if err != nil {
		return nil, err
	}

	client, err := c.getScmClient(repoURL)
	if err != nil {
		return nil, err
	}

	commit, _, err := client.Git.FindCommit(ctx, owner+"/"+repo, branch)
	if err != nil {
		return nil, err
	}

	return &types.GitCommit{
		Sha:     commit.Sha,
		Message: commit.Message,
	}, nil
}

func (c *Service) GetCommitBySha(ctx context.Context, repoURL string, sha string) (*types.GitCommit, error) {
	owner, repo, err := helpers.SplitGitRepoUrl(repoURL)
	if err != nil {
		return nil, err
	}

	client, err := c.getScmClient(repoURL)
	if err != nil {
		return nil, err
	}

	commit, _, err := client.Git.FindCommit(ctx, owner+"/"+repo, sha)
	if err != nil {
		return nil, err
	}

	return &types.GitCommit{
		Sha:     commit.Sha,
		Message: commit.Message,
	}, nil
}
