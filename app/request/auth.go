package request

import (
	"net/http"
)

const (
	QueryParamAccessToken   = "access_token"
	QueryParamIncludeCookie = "include_cookie"
)

func GetAccessTokenFromQuery(r *http.Request) (string, bool) {
	return QueryParam(r, QueryParamAccessToken)
}

func GetIncludeCookieFromQueryOrDefault(r *http.Request, dflt bool) (bool, error) {
	return QueryParamAsBoolOrDefault(r, QueryParamIncludeCookie, dflt)
}

func GetTokenFromCookie(r *http.Request, cookieName string) (string, bool) {
	return GetCookie(r, cookieName)
}
