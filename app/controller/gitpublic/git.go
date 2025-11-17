package gitpublic

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/check"
)

type RepoLoadModel struct {
	RepoURL string `json:"repoURL"`
}

func (c *Controller) FindRepo(ctx context.Context, repoURL string) (*types.GitRepo, error) {
	if err := c.sanitizeRepoURL(repoURL); err != nil {
		return nil, err
	}
	return c.gitPublicSvc.FindRepo(ctx, repoURL)
}

func (c *Controller) ListBranches(ctx context.Context, repoURL string) ([]*types.GitBranch, error) {
	if err := c.sanitizeRepoURL(repoURL); err != nil {
		return nil, err
	}

	return c.gitPublicSvc.ListBranches(ctx, repoURL)
}

func (c *Controller) GetLatestCommit(ctx context.Context, repoURL string, branch string) (*types.GitCommit, error) {
	return c.gitPublicSvc.GetLatestCommit(ctx, repoURL, branch)
}

func (c *Controller) GetCommitBySha(ctx context.Context, repoURL string, sha string) (*types.GitCommit, error) {
	return c.gitPublicSvc.GetCommitBySha(ctx, repoURL, sha)
}

func (c *Controller) sanitizeRepoURL(repoURL string) error {
	errors := check.NewValidationErrors()
	if err := check.GitRepo(repoURL); err != nil {
		errors.AddValidationError("repoURL", err)
	}
	if errors.HasError() {
		return errors
	}
	return nil
}
