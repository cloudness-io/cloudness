package enum

import (
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

const (
	id            = "id"
	uid           = "uid"
	name          = "name"
	email         = "email"
	admin         = "admin"
	number        = "number"
	created       = "created"
	createdAt     = "created_at"
	createdBy     = "created_by"
	updated       = "updated"
	deleted       = "deleted"
	deletedAt     = "deleted_at"
	sequence      = "sequence"
	date          = "date"
	defaultString = "default"
	undefined     = "undefined"
	asc           = "asc"
	ascending     = "ascending"
	desc          = "desc"
	descending    = "descending"
)

func Sanitize[E constraints.Ordered](element E, all func() ([]E, E)) (E, bool) {
	allValues, defValue := all()
	var empty E
	if element == empty && defValue != empty {
		return defValue, true
	}
	idx, exists := slices.BinarySearch(allValues, element)
	if exists {
		return allValues[idx], true
	}
	return defValue, false
}

func toInterfaceSlice[T any](vals []T) []any {
	res := make([]any, len(vals))
	for i := range vals {
		res[i] = vals[i]
	}
	return res
}

func sortEnum[T constraints.Ordered](slice []T) []T {
	slices.Sort(slice)
	return slice
}
