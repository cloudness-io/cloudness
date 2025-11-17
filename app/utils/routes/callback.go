package routes

import (
	"fmt"

	"github.com/cloudness-io/cloudness/types/enum"
)

const (
	CallbackOAuthGithub = "/auth/github/callback"
)

func GetOAuthRedirectUrl(provider enum.AuthProvider) string {
	return fmt.Sprintf("/auth/%s/redirect", provider)
}

func GetOAuthCallbackUrl(provider enum.AuthProvider) string {
	return fmt.Sprintf("/auth/%s/callback", provider)
}
