package enum

import "strings"

// TenantAttr defines tenant attibutes that can be used for sorting and filtering.
type TenantAttr int

// Order enumeration.
const (
	TenantAttrNone TenantAttr = iota
	TenantAttrUID
	TenantAttrCreated
	TenantAttrUpdated
	TenantAttrDeleted
)

// ParseTenantAttr parses the space attribute string and returns the equivalent enumeration.
func ParseTenantAttr(s string) TenantAttr {
	switch strings.ToLower(s) {
	case uid:
		return TenantAttrUID
	case created, createdAt:
		return TenantAttrCreated
	case updated:
		return TenantAttrUpdated
	case deleted, deletedAt:
		return TenantAttrDeleted
	default:
		return TenantAttrNone
	}
}

// String returns the string representation of the attribute.
func (a TenantAttr) String() string {
	switch a {
	case TenantAttrUID:
		return uid
	case TenantAttrCreated:
		return created
	case TenantAttrUpdated:
		return updated
	case TenantAttrDeleted:
		return deleted
	case TenantAttrNone:
		return ""
	default:
		return undefined
	}
}
