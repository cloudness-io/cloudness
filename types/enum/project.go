package enum

import "strings"

// ProjectAttr defines project attibutes that can be used for sorting and filtering.
type ProjectAttr int

// Order enumeration.
const (
	ProjectAttrNone ProjectAttr = iota
	ProjectAttrUID
	ProjectAttrCreated
	ProjectAttrUpdated
	ProjectAttrDeleted
)

// ParseProjectAttr parses the space attribute string and returns the equivalent enumeration.
func ParseProjectAttr(s string) ProjectAttr {
	switch strings.ToLower(s) {
	case uid:
		return ProjectAttrUID
	case created, createdAt:
		return ProjectAttrCreated
	case updated:
		return ProjectAttrUpdated
	case deleted, deletedAt:
		return ProjectAttrDeleted
	default:
		return ProjectAttrNone
	}
}

// String returns the string representation of the attribute.
func (a ProjectAttr) String() string {
	switch a {
	case ProjectAttrUID:
		return uid
	case ProjectAttrCreated:
		return created
	case ProjectAttrUpdated:
		return updated
	case ProjectAttrDeleted:
		return deleted
	case ProjectAttrNone:
		return ""
	default:
		return undefined
	}
}
