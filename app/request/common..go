package request

import (
	"context"
	"strings"
)

const (
	HeaderAuthorization = "Authorization"
)

func IsLoginOrRegistrationPage(ctx context.Context) bool {
	currentUrl, ok := CurrentFullUrlFrom(ctx)
	if !ok {
		return false
	}
	if strings.HasPrefix(currentUrl, "/login") || strings.HasPrefix(currentUrl, "/register") {
		return true
	}
	return false
}
