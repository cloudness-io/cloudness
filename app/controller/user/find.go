package user

import (
	"context"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/types"
)

func (c *Controller) FindUserByEmail(ctx context.Context, email string) (*types.User, error) {
	return c.principalStore.FindUserByEmail(ctx, email)
}

/*
 * Find tries to find the provided user.
 */
func (c *Controller) Find(ctx context.Context, session *auth.Session,
	userUID string) (*types.User, error) {
	user, err := c.FindNoAuth(ctx, userUID)
	if err != nil {
		return nil, err
	}

	// Ensure principal has required permissions on parent.
	// if err = apiauth.CheckUser(ctx, c.authorizer, session, user, enum.PermissionUserView); err != nil {
	// 	return nil, err
	// }

	return user, nil
}

/*
 * FindNoAuth finds a user without auth checks.
 * WARNING: Never call as part of user flow.
 */
func (c *Controller) FindNoAuth(ctx context.Context, userUID string) (*types.User, error) {
	return findUserFromUID(ctx, c.principalStore, userUID)
}
