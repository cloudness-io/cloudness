package authn

import (
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/types"

	"github.com/google/wire"
)

// WireSet provides a wire set for this package.
var WireSet = wire.NewSet(
	ProvideAuthenticator,
)

func ProvideAuthenticator(
	config *types.Config,
	principalStore store.PrincipalStore,
	tokenStore store.TokenStore,
) Authenticator {
	return NewTokenAuthenticator(principalStore, tokenStore, config.Token.CookieName)
}
