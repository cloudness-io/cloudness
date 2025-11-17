package request

import (
	"context"

	"github.com/cloudness-io/cloudness/types/enum"
)

func IsAuthenticated(ctx context.Context) bool {
	_, ok := PrincipalFrom(ctx)
	return ok
}

func IsSuperAdmin(ctx context.Context) bool {
	instance, ok := InstanceSettingsFrom(ctx)
	session, sessionOk := AuthSessionFrom(ctx)
	return ok && sessionOk && instance.SuperAdmin != nil && session.Principal.ID == *instance.SuperAdmin
}

func IsTeamAdmin(ctx context.Context) bool {
	membership, ok := TenantMembershipFrom(ctx)
	return ok && membership.Role == enum.TenantRoleAdmin
}

func IsProjectOwner(ctx context.Context) bool {
	membership, ok := ProjectMembershipFrom(ctx)
	return IsTeamAdmin(ctx) || (ok && membership.Role == enum.ProjectRoleOwner)
}

func IsProjectContributor(ctx context.Context) bool {
	membership, ok := ProjectMembershipFrom(ctx)
	return IsProjectOwner(ctx) || (ok && membership.Role == enum.ProjectRoleContributor)
}

func IsProjectViewer(ctx context.Context) bool {
	membership, ok := ProjectMembershipFrom(ctx)
	return IsProjectContributor(ctx) || (ok && membership.Role == enum.ProjectRoleViewer)
}
