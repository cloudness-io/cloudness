package user

import (
	"context"

	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/store/database/dbtx"
	"github.com/cloudness-io/cloudness/types"
)

type Controller struct {
	tx             dbtx.Transactor
	principalStore store.PrincipalStore
}

func NewController(
	tx dbtx.Transactor,
	userStore store.PrincipalStore,
) *Controller {
	return &Controller{
		tx:             tx,
		principalStore: userStore,
	}
}

func findUserFromUID(ctx context.Context, userStore store.PrincipalStore, userUID string) (*types.User, error) {
	return userStore.FindUserByUID(ctx, userUID)
}

func findUserFromID(ctx context.Context, userStore store.PrincipalStore, userID int64) (*types.User, error) {
	return userStore.FindUser(ctx, userID)
}

// func findUserFromEmail(ctx context.Context, userStore store.UserStore, email string) (*types.User, error) {
// 	return userStore.FindUserByEmail(ctx, email)
// }
