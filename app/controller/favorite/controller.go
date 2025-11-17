package favorite

import "github.com/cloudness-io/cloudness/app/store"

type Controller struct {
	favStore store.FavoriteStore
}

func NewController(favStore store.FavoriteStore) *Controller {
	return &Controller{
		favStore: favStore,
	}
}
