package githubapp

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
)

func (s Service) Find(ctx context.Context, tenantID, projectID, githubAppID int64) (*types.GithubApp, error) {
	githubApp, err := s.githubAppStore.Find(ctx, tenantID, projectID, githubAppID)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, err
	}
	return githubApp, nil
}

func (s *Service) FindByUID(ctx context.Context, tenantID, projectID, githubAppUID int64) (*types.GithubApp, error) {
	githubApp, err := s.githubAppStore.FindByUID(ctx, tenantID, projectID, githubAppUID)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, err
	}
	return githubApp, nil
}
