package database

import "strings"

func PartialMatch(column, value string) (string, string) {
	var (
		n       int
		escaped bool
	)

	if n, value = len(value), strings.ReplaceAll(value, `\`, `\\`); n < len(value) {
		escaped = true
	}
	if n, value = len(value), strings.ReplaceAll(value, "_", `\_`); n < len(value) {
		escaped = true
	}
	if n, value = len(value), strings.ReplaceAll(value, "%", `\%`); n < len(value) {
		escaped = true
	}

	sb := strings.Builder{}
	sb.WriteString("LOWER(")
	sb.WriteString(column)
	sb.WriteString(") LIKE '%' || LOWER(?) || '%'")
	if escaped {
		sb.WriteString(` ESCAPE '\'`)
	}

	return sb.String(), value
}
