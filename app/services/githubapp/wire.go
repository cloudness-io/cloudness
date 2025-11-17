package githubapp

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideService,
)

func ProvideService(tx dbtx.Transactor, ghAppStore store.GithubAppStore, privateKeyStore store.PrivateKeyStore) *Service {
	return New(tx, ghAppStore, privateKeyStore)
}
