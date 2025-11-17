package types

import "github.com/cloudness-io/cloudness/types/enum"

type ProjectMembershipKey struct {
	TenantID           int64            `json:"-"`
	TenantMembershipID int64            `json:"-"`
	ProjectID          int64            `json:"-"`
	PrincipalID        int64            `json:"-"`
	Role               enum.ProjectRole `json:"role"`
}

type ProjectMembership struct {
	*ProjectMembershipKey
	CreatedBy int64 `json:"-"`

	Created int64 `json:"created"`
	Updated int64 `json:"updated"`
}

type ProjectMembershipUser struct {
	*ProjectMembership
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}
