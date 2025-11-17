package types

// TokenResponse is returned as part of token creation for PAT / SAT / User Session.
type OauthLoginResponse struct {
	LoginUrl string `json:"login_url"`
}
