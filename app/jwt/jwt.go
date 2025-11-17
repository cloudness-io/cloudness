package jwt

import (
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

const (
	issuer = "Cloudness"
)

// Claims defines Cloudness jwt claims.
type Claims struct {
	jwt.StandardClaims

	PrincipalID int64 `json:"pid,omitempty"`

	Token *SubClaimsToken `json:"tkn,omitempty"`
}

// SubClaimsToken contains information about the token the JWT was created for.
type SubClaimsToken struct {
	Type enum.TokenType `json:"typ,omitempty"`
	ID   int64          `json:"id,omitempty"`
}

// GenerateForToken generates a jwt for a given token.
func GenerateForToken(token *types.Token, secret string) (string, error) {
	var expiresAt int64
	if token.ExpiresAt != nil {
		expiresAt = *token.ExpiresAt
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer: issuer,
			// times required to be in sec not millisec
			IssuedAt:  token.IssuedAt / 1000,
			ExpiresAt: expiresAt / 1000,
		},
		PrincipalID: token.PrincipalID,
		Token: &SubClaimsToken{
			Type: token.Type,
			ID:   token.ID,
		},
	})

	res, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return "", errors.Wrap(err, "Failed to sign token")
	}

	return res, nil
}
