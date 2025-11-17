package githubapp

import (
	"context"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
)

type Service struct {
	tx              dbtx.Transactor
	githubAppStore  store.GithubAppStore
	privateKeyStore store.PrivateKeyStore
}

func New(
	tx dbtx.Transactor,
	githubAppStore store.GithubAppStore,
	privateKeyStore store.PrivateKeyStore,
) *Service {
	return &Service{
		tx:              tx,
		githubAppStore:  githubAppStore,
		privateKeyStore: privateKeyStore,
	}
}

func (s *Service) Create(ctx context.Context, ghApp *types.GithubApp) (*types.GithubApp, error) {
	return s.githubAppStore.Create(ctx, ghApp)
}
