package enum

import (
	"strings"
)

// VolumeAttr defines volume attributes that can be used for sorting and filtering.
type VolumeAttr int

// VolumeAttr enumeration.
const (
	VolumeAttrNone VolumeAttr = iota
	VolumeAttrCreated
	VolumeAttrUpdated
	VolumeAttrDeleted
)

// ParseVolumeAttr parses the volume attribute string
// and returns the equivalent enumeration.
func ParseVolumeAttr(s string) VolumeAttr {
	switch strings.ToLower(s) {
	case created, createdAt:
		return VolumeAttrCreated
	case updated:
		return VolumeAttrUpdated
	case deleted, deletedAt:
		return VolumeAttrDeleted
	default:
		return VolumeAttrNone
	}
}

// String returns the string representation of the attribute.
func (a VolumeAttr) String() string {
	switch a {
	case VolumeAttrCreated:
		return created
	case VolumeAttrUpdated:
		return updated
	case VolumeAttrDeleted:
		return deleted
	case VolumeAttrNone:
		return ""
	default:
		return undefined
	}
}
