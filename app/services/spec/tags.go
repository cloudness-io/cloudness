package spec

import (
	"context"

	"github.com/cloudness-io/cloudness/errors"
	"github.com/cloudness-io/cloudness/types"
)

func (s *Service) GetLatestTag(ctx context.Context, app *types.Application) (string, error) {
	switch true {
	// Git
	case app.Spec.IsGit():
		gitSpec := app.Spec.Build.Source.Git
		switch true {
		case app.GithubAppID != nil:
			// github app
			ghApp, err := s.ghAppSvc.Find(ctx, app.TenantID, app.ProjectID, *app.GithubAppID)
			if err != nil {
				return "", err
			}
			commit, err := s.ghAppSvc.GetLatestCommit(ctx, ghApp, gitSpec.RepoURL, gitSpec.Branch)
			if err != nil {
				return "", err
			}
			gitSpec.Commit = commit.GetSHA()
			return commit.Commit.GetMessage(), nil
		default:
			// git public
			commit, err := s.gitpublicSvc.GetLatestCommit(ctx, gitSpec.RepoURL, gitSpec.Branch)
			if err != nil {
				return "", err
			}
			gitSpec.Commit = commit.Sha
			return commit.Message, nil
		}
	// Registry
	case app.Spec.IsRegistry():
		return app.Spec.Build.Source.Registry.Image, nil
	}

	return "", nil
}

func (s *Service) GetCurrentTag(ctx context.Context, app *types.Application) (string, error) {
	switch true {
	case app.Spec.IsGit():
		gitSpec := app.Spec.Build.Source.Git
		if gitSpec.Commit == "" {
			return "", errors.BadRequest("No commit sha found")
		}

		switch true {
		case app.GithubAppID != nil:
			// github app
			ghApp, err := s.ghAppSvc.Find(ctx, app.TenantID, app.ProjectID, *app.GithubAppID)
			if err != nil {
				return "", err
			}
			if commit, err := s.ghAppSvc.GetCommitBySha(ctx, ghApp, gitSpec.RepoURL, gitSpec.Commit); err != nil {
				return "", err
			} else {
				return commit.Commit.GetMessage(), nil
			}
		default:
			// git public
			if commit, err := s.gitpublicSvc.GetCommitBySha(ctx, gitSpec.RepoURL, gitSpec.Commit); err != nil {
				return "", err
			} else {
				return commit.Message, nil
			}
		}

	case app.Spec.IsRegistry():
		return app.Spec.Build.Source.Registry.Image, nil
	}
	return "", errors.BadRequest("No type configured")
}
