package variable

const (
	SystemVarTeamID           = "CLOUDNESS_TEAM_ID"
	SystemVarProjectID        = "CLOUDNESS_PROJECT_ID"
	SystemVarEnvironmentID    = "CLOUDNESS_ENVIRONMENT_ID"
	SystemVarAppName          = "CLOUDNESS_APP_NAME"
	SystemVarAppID            = "CLOUDNESS_APP_ID"
	SystemVarAppPrivateDomain = "CLOUDNESS_PRIVATE_DOMAIN"
	SystemVarAppPublicDomain  = "CLOUDNESS_PUBLIC_DOMAIN"
	SystemVarAppPublicURL     = "CLOUDNESS_PUBLIC_URL"
	SystemVarAppTCPPort       = "CLOUDNESS_TCP_APPLICATION_PORT"
	SystemVarServiceName      = "CLOUDNESS_SERVICE_NAME"
)

var systemVarMap = map[string]bool{
	SystemVarTeamID:           true,
	SystemVarProjectID:        true,
	SystemVarEnvironmentID:    true,
	SystemVarAppName:          true,
	SystemVarAppID:            true,
	SystemVarAppPrivateDomain: true,
	SystemVarAppPublicDomain:  true,
	SystemVarAppPublicURL:     true,
	SystemVarAppTCPPort:       true,
	SystemVarServiceName:      true,
}

func IsSystemVar(key string) bool {
	return systemVarMap[key]
}
