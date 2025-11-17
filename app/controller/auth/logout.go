package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudness-io/cloudness/app/auth"
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/types/enum"
)

func (c *Controller) Logout(ctx context.Context, session *auth.Session) error {

	var (
		tokenID   int64
		tokenType enum.TokenType
	)

	if session == nil {
		return usererror.ErrUnauthorized
	}

	switch t := session.Metadata.(type) {
	case *auth.TokenMetadata:
		tokenID = t.TokenID
		tokenType = t.TokenType
	default:
		return errors.New("provided jwt doesn't support logout")
	}

	if tokenType != enum.TokenTypeSession {
		return usererror.BadRequestf("unsupported logout token type %v", tokenType)
	}

	err := c.tokenStore.Delete(ctx, tokenID)
	if err != nil {
		return fmt.Errorf("failed to delete token from store: %w", err)
	}

	return nil
}
