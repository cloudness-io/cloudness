package favorite

import (
	"github.com/cloudness-io/cloudness/app/store"

	"github.com/google/wire"
)

var WireSet = wire.NewSet(
	ProvideController,
)

func ProvideController(favStore store.FavoriteStore) *Controller {
	return NewController(favStore)
}
