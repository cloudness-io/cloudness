package authn

import (
	"errors"
	"net/http"

	"github.com/cloudness-io/cloudness/app/auth"
)

var (
	// ErrNoAuthData that is returned if the authorizer doesn't find any data in the request that can be used for auth.
	ErrNoAuthData = errors.New("the request doesn't contain any auth data that can be used by the Authorizer")
)

// Authenticator is an abstraction of an entity that's responsible for authenticating principals
// that are making calls via HTTP.
type Authenticator interface {
	/*
	 * Tries to authenticate the acting principal if credentials are available.
	 * Returns:
	 *		(session, nil) 		    - request contains auth data and principal was verified
	 *		(nil, ErrNoAuthData)	- request doesn't contain any auth data
	 *		(nil, err)  			- request contains auth data but verification failed
	 */
	Authenticate(r *http.Request) (*auth.Session, error)
}
