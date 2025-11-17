package auth

import "github.com/cloudness-io/cloudness/types/enum"

type Metadata interface {
	ImpactsAuthorization() bool
}

// EmptyMetadata represents the state when the auth session doesn't have any extra metadata.
type EmptyMetadata struct{}

func (m *EmptyMetadata) ImpactsAuthorization() bool {
	return false
}

// TokenMetadata contains information about the token that was used during auth.
type TokenMetadata struct {
	TokenType enum.TokenType
	TokenID   int64
}

func (m *TokenMetadata) ImpactsAuthorization() bool {
	return false
}
