package enum

import "strings"

// EnvironmentAttr defines project attibutes that can be used for sorting and filtering.
type EnvironmentAttr int

// Order enumeration.
const (
	EnvironmentAttrNone EnvironmentAttr = iota
	EnvironmentAttrUID
	EnvironmentAttrCreated
	EnvironmentAttrUpdated
	EnvironmentAttrDeleted
	EnvironmentAttrSequence
)

// ParseEnvironmentAttr parses the space attribute string and returns the equivalent enumeration.
func ParseEnvironmentAttr(s string) EnvironmentAttr {
	switch strings.ToLower(s) {
	case uid:
		return EnvironmentAttrUID
	case created, createdAt:
		return EnvironmentAttrCreated
	case updated:
		return EnvironmentAttrUpdated
	case deleted, deletedAt:
		return EnvironmentAttrDeleted
	case sequence:
		return EnvironmentAttrSequence
	default:
		return EnvironmentAttrNone
	}
}

// String returns the string representation of the attribute.
func (a EnvironmentAttr) String() string {
	switch a {
	case EnvironmentAttrUID:
		return uid
	case EnvironmentAttrCreated:
		return created
	case EnvironmentAttrUpdated:
		return updated
	case EnvironmentAttrDeleted:
		return deleted
	case EnvironmentAttrSequence:
		return sequence
	case EnvironmentAttrNone:
		return ""
	default:
		return undefined
	}
}
