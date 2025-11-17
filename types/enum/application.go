package enum

import (
	"strings"
)

// ApplicationAttr defines application attributes that can be used for sorting and filtering.
type ApplicationAttr int

// ApplicationAttr enumeration.
const (
	ApplicationAttrNone ApplicationAttr = iota
	ApplicationAttrCreated
	ApplicationAttrUpdated
	ApplicationAttrDeleted
)

// ParseApplicationAttr parses the application attribute string
// and returns the equivalent enumeration.
func ParseApplicationAttr(s string) ApplicationAttr {
	switch strings.ToLower(s) {
	case created, createdAt:
		return ApplicationAttrCreated
	case updated:
		return ApplicationAttrUpdated
	case deleted, deletedAt:
		return ApplicationAttrDeleted
	default:
		return ApplicationAttrNone
	}
}

// String returns the string representation of the attribute.
func (a ApplicationAttr) String() string {
	switch a {
	case ApplicationAttrCreated:
		return created
	case ApplicationAttrUpdated:
		return updated
	case ApplicationAttrDeleted:
		return deleted
	case ApplicationAttrNone:
		return ""
	default:
		return undefined
	}
}
