package enum

type AuthProvider string

const (
	AuthProviderPassword AuthProvider = "password"
	AuthProviderGithub   AuthProvider = "github"
	AuthProviderGitlab   AuthProvider = "gitlab"
	AuthProviderGoogle   AuthProvider = "google"
)

func ProviderFromString(s string) AuthProvider {
	switch s {
	case string(AuthProviderPassword):
		return AuthProviderPassword
	case string(AuthProviderGithub):
		return AuthProviderGithub
	case string(AuthProviderGitlab):
		return AuthProviderGitlab
	case string(AuthProviderGoogle):
		return AuthProviderGoogle
	default:
		return ""
	}
}
