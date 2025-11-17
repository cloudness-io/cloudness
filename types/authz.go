package types

import "github.com/cloudness-io/cloudness/types/enum"

// PermissionCheck represents a permission check.
type PermissionCheck struct {
	Scope      Scope
	Resource   Resource
	Permission enum.Permission
}

// Resource represents the resource of a permission check.
// Note: Keep the name empty in case access is requested for all resources of that type.
type Resource struct {
	Type       enum.ResourceType
	Identifier string
}

// Scope represents the scope of a permission check
// Notes:
//   - In case the permission check is for resource App, keep app empty (app is resource, not scope)
//   - In case the permission check is for resource Tenant, Tenant is an ancestor of the tenant (tenant is
//     resource, not scope)
//   - App isn't use as of now (will be useful once we add access control for app child resources, e.g. branches).
type Scope struct {
	Tenant string
	App    string
}
