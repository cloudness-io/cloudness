package token

import (
	"context"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/cloudness-io/cloudness/app/jwt"
	"github.com/cloudness-io/cloudness/app/store"
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/gotidy/ptr"
)

const (
	// userSessionTokenLifeTime is the duration a login / register token is valid.
	// TODO: Users should be able to list / delete session tokens if they want to cleanup earlier.
	userSessionTokenLifeTime time.Duration = 30 * 24 * time.Hour // 30 days.
)

func generateIdentifier(prefix string) string {
	r := rand.IntN(0x10000)
	return fmt.Sprintf("%s-%08x-%04x", prefix, time.Now().Unix(), r)
}

func CreateUserSession(
	ctx context.Context,
	tokenStore store.TokenStore,
	user *types.User,
) (*types.Token, string, error) {
	principal := user.ToPrincipal()
	return create(
		ctx,
		tokenStore,
		enum.TokenTypeSession,
		principal,
		principal,
		generateIdentifier("login"),
		ptr.Duration(userSessionTokenLifeTime),
	)
}

func create(
	ctx context.Context,
	tokenStore store.TokenStore,
	tokenType enum.TokenType,
	createdBy *types.Principal,
	createdFor *types.Principal,
	identifier string,
	lifetime *time.Duration,
) (*types.Token, string, error) {
	issuedAt := time.Now().UTC()

	var expiresAt *int64
	if lifetime != nil {
		expiresAt = ptr.Int64(issuedAt.Add(*lifetime).UnixMilli())
	}

	// create db entry first so we get the id.
	token := types.Token{
		Type:        tokenType,
		Identifier:  identifier,
		PrincipalID: createdFor.ID,
		IssuedAt:    issuedAt.UnixMilli(),
		ExpiresAt:   expiresAt,
		CreatedBy:   createdBy.ID,
	}

	err := tokenStore.Create(ctx, &token)
	if err != nil {
		return nil, "", fmt.Errorf("failed to store token in db: %w", err)
	}

	// create jwt token.
	jwtToken, err := jwt.GenerateForToken(&token, createdFor.Salt)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create jwt token: %w", err)
	}

	return &token, jwtToken, nil
}
