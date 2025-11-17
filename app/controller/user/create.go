package user

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/store"
	"github.com/cloudness-io/cloudness/types"

	"github.com/dchest/uniuri"
	"golang.org/x/crypto/bcrypt"
)

// CreateAdminUserIntegration struct    is the input to create and integration along with user creation for SCM, should be
// populated only for admin
type CreateAdminUserIntegration struct {
	Token   string
	Expires int64
	Refresh string
}

// CreateInput is the input used for create operations.
// On purpose don't expose admin, has to be enabled explicitly.
type CreateInput struct {
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	Password    string `json:"-"`
	AvatarUrl   string `json:"avatar_url"`
}

/*
 * CreateNoAuth creates a new user without auth checks.
 * WARNING: Never call as part of user flow.
 *
 * Note: take admin separately to avoid potential vulnerabilities for user calls.
 */
func (c *Controller) CreateNoAuth(ctx context.Context, in *CreateInput) (*types.User, error) {
	currUser, err := c.principalStore.FindUserByEmail(ctx, in.Email)
	if err != nil && !errors.Is(err, store.ErrResourceNotFound) {
		return nil, err
	}

	if currUser != nil {
		return currUser, nil
	}

	var hash []byte
	if in.Password != "" {
		hash, err = bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to create hash: %w", err)
		}
	}
	user := &types.User{
		UID:         strings.ToLower(in.Email),
		DisplayName: in.DisplayName,
		AvatarUrl:   in.AvatarUrl,
		Email:       in.Email,
		Password:    string(hash),
		Salt:        uniuri.NewLen(uniuri.UUIDLen),
		Created:     time.Now().UTC().UnixMilli(),
		Updated:     time.Now().UTC().UnixMilli(),
	}
	return c.principalStore.CreateUser(ctx, user)
}
