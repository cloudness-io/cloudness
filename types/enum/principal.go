package enum

// PrincipalType defines the supported types of principals.
type PrincipalType string

func (PrincipalType) Enum() []interface{}                    { return toInterfaceSlice(principalTypes) }
func (s PrincipalType) Sanitize() (PrincipalType, bool)      { return Sanitize(s, GetAllPrincipalTypes) }
func GetAllPrincipalTypes() ([]PrincipalType, PrincipalType) { return principalTypes, "" }

const (
	// PrincipalTypeUser represents a user.
	PrincipalTypeUser PrincipalType = "user"
	// PrincipalTypeServiceAccount represents a service account.
	PrincipalTypeServiceAccount PrincipalType = "serviceaccount"
	// PrincipalTypeService represents a service.
	PrincipalTypeService PrincipalType = "service"
)

var principalTypes = sortEnum([]PrincipalType{
	PrincipalTypeUser,
	PrincipalTypeServiceAccount,
	PrincipalTypeService,
})
