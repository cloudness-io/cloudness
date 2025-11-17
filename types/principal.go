package types

import (
	"strings"

	"github.com/cloudness-io/cloudness/types/enum"
)

// AnonymousPrincipalUID is an internal UID for anonymous principals.
const AnonymousPrincipalUID = "anonymous"

// Principal represents the identity of an acting entity (User, ServiceAccount, Service).
type Principal struct {
	ID          int64              `db:"principal_id"           json:"-"`
	UID         string             `db:"principal_uid"          json:"uid"`
	Email       string             `db:"principal_email"        json:"email"`
	Type        enum.PrincipalType `db:"principal_type"         json:"type"`
	DisplayName string             `db:"principal_display_name" json:"display_name"`

	// Should be part of principal or not?
	Blocked  bool   `db:"principal_blocked"            json:"blocked"`
	Password string `db:"principal_user_password"      json:"-"`
	Salt     string `db:"principal_salt"               json:"-"`

	// Other info
	Created int64 `db:"principal_created"                json:"created"`
	Updated int64 `db:"principal_updated"                json:"updated"`
}

func (p *Principal) ToPrincipalInfo() *PrincipalInfo {
	return &PrincipalInfo{
		ID:          p.ID,
		UID:         p.UID,
		DisplayName: p.DisplayName,
		Email:       strings.ToLower(p.Email),
		Type:        p.Type,
		Created:     p.Created,
		Updated:     p.Updated,
	}
}

// PrincipalInfo is a compressed representation of a principal we return as part of non-principal APIs.
type PrincipalInfo struct {
	ID          int64              `json:"id"`
	UID         string             `json:"uid"`
	DisplayName string             `json:"display_name"`
	Email       string             `json:"email"`
	Type        enum.PrincipalType `json:"type"`
	Created     int64              `json:"created"`
	Updated     int64              `json:"updated"`
}

func (p *PrincipalInfo) Identifier() int64 {
	return p.ID
}

type PrincipalFilter struct {
	Page  int                  `json:"page"`
	Size  int                  `json:"size"`
	Query string               `json:"query"`
	Types []enum.PrincipalType `json:"types"`
}
