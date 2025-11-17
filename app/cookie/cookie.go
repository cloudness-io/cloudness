package cookie

import (
	"errors"
	"net/http"
	"time"

	"github.com/cloudness-io/cloudness/types"
)

func IncludeTokenCookie(
	r *http.Request, w http.ResponseWriter,
	tokenResponse *types.TokenResponse,
	cookieName string,
) {
	cookie := newEmptyTokenCookie(r, cookieName)
	cookie.Value = tokenResponse.AccessToken
	if tokenResponse.Token.ExpiresAt != nil {
		cookie.Expires = time.UnixMilli(*tokenResponse.Token.ExpiresAt)
	}

	http.SetCookie(w, cookie)
}

func DeleteTokenCookieIfPresent(r *http.Request, w http.ResponseWriter, cookieName string) {
	// if no token is present in the cookies, nothing todo.
	// No other error type expected here - and even if there is, let's try best effort deletion.
	_, err := r.Cookie(cookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return
	}

	cookie := newEmptyTokenCookie(r, cookieName)
	cookie.Value = ""
	cookie.Expires = time.UnixMilli(0) // this effectively tells the browser to delete the cookie

	http.SetCookie(w, cookie)
}

func newEmptyTokenCookie(r *http.Request, cookieName string) *http.Cookie {
	return &http.Cookie{
		Name:     cookieName,
		SameSite: http.SameSiteLaxMode,
		HttpOnly: true,
		Path:     "/",
		Domain:   r.URL.Hostname(),
		Secure:   r.URL.Scheme == "https",
	}
}
