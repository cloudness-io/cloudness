package auth

import (
	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"
)

// AnonymousPrincipal is an in-memory principal for users with no auth data.
// Authorizer is in charge of handling anonymous access.
var AnonymousPrincipal = types.Principal{
	ID:   -1,
	UID:  types.AnonymousPrincipalUID,
	Type: enum.PrincipalTypeUser,
}

func IsAnonymousSession(session *Session) bool {
	return session != nil && session.Principal.UID == types.AnonymousPrincipalUID
}
