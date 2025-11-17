package githubapp

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
)

func (s *Service) List(ctx context.Context, tenantID, projectID int64) ([]*types.GithubApp, error) {
	ghApps, err := s.githubAppStore.List(ctx, tenantID, projectID)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, err
	}
	return ghApps, nil
}
