package enum

type TenantRole string

const (
	TenantRoleAdmin  TenantRole = "admin"
	TenantRoleMember TenantRole = "member"
)

var TenantRoles = sortEnum([]TenantRole{
	TenantRoleAdmin,
	TenantRoleMember,
})

var TenantRolesStr = []string{
	string(TenantRoleMember),
	string(TenantRoleAdmin),
}

func TenantRoleFromString(s string) TenantRole {
	switch s {
	case string(TenantRoleAdmin):
		return TenantRoleAdmin
	case string(TenantRoleMember):
		return TenantRoleMember
	default:
		return ""
	}
}
