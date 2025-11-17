package types

import (
	"github.com/cloudness-io/cloudness/types/enum"
)

// Represents server side infos stored for tokens we distribute.
type Token struct {
	// TODO: int64 ID doesn't match DB
	ID          int64          `db:"token_id"                 json:"-"`
	PrincipalID int64          `db:"token_principal_id"       json:"principal_id"`
	Type        enum.TokenType `db:"token_type"               json:"type"`
	Identifier  string         `db:"token_uid"                json:"identifier"`
	// ExpiresAt is an optional unix time that if specified restricts the validity of a token.
	ExpiresAt *int64 `db:"token_expires_at"         json:"expires_at,omitempty"`
	// IssuedAt is the unix time at which the token was issued.
	IssuedAt  int64 `db:"token_issued_at"          json:"issued_at"`
	CreatedBy int64 `db:"token_created_by"         json:"created_by"`
}

// TokenResponse is returned as part of token creation for PAT / SAT / User Session.
type TokenResponse struct {
	AccessToken string `json:"access_token"`
	Token       Token  `json:"token"`
}
