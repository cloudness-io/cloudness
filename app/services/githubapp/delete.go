package githubapp

import (
	"context"
	"errors"

	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"
)

func (s *Service) Delete(ctx context.Context, ghApp *types.GithubApp) error {
	return s.tx.WithTx(ctx, func(ctx context.Context) error {
		if ghApp.PrivateKeyID > 0 {
			err := s.privateKeyStore.Delete(ctx, ghApp.TenantID, ghApp.PrivateKeyID)
			if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
				return err
			}
		}

		return s.githubAppStore.Delete(ctx, ghApp)
	})

}
