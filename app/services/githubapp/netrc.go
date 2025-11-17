package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/types"
)

func (s *Service) Netrc(ctx context.Context, tenantID, projectID, githubAppID int64, repoURL string) (*types.Netrc, error) {
	ghApp, err := s.Find(ctx, tenantID, projectID, githubAppID)
	if err != nil {
		return nil, err
	}
	token, err := s.generateAccessToken(ctx, ghApp)
	if err != nil {
		return nil, err
	}
	netrc := &types.Netrc{
		Password: token,
		Login:    "x-access-token",
	}

	if err := netrc.SetMachine(repoURL); err != nil {
		return nil, err
	}

	return netrc, nil
}
