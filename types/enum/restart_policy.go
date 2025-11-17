package enum

type RestartPolicyType string

func GetAllRestartPolicyStr() []string { return builderTypesStr }

const (
	RestartPolicyTypeOnFailure RestartPolicyType = "ON_FAILURE"

	RestartPolicyTypeAlways RestartPolicyType = "ALWAYS"

	RestartPolicyTypeNever RestartPolicyType = "NEVER"
)

var restartPolicyTypesStr = sortEnum([]string{
	string(RestartPolicyTypeOnFailure),
	string(RestartPolicyTypeAlways),
	string(RestartPolicyTypeNever),
})
