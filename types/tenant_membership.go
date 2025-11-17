package types

import "github.com/cloudness-io/cloudness/types/enum"

type TenantMembershipKey struct {
	ID                int64           `json:"-"`
	TenantID          int64           `json:"-"`
	TenantUID         int64           `json:"tenant_uid"`
	PrincipalID       int64           `json:"-"`
	TenantName        string          `json:"tenant_name"`
	TenantDescription string          `json:"tenant_description"`
	Role              enum.TenantRole `json:"role"`
}

type TenantMembership struct {
	*TenantMembershipKey
	CreatedBy int64 `json:"-"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type TenantMembershipUser struct {
	*TenantMembership
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}
