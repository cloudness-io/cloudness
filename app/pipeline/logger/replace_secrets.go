package logger

import "strings"

// newSecretReplacer returns a replacer that avoids replacing substrings inside unicode characters.
// It only replaces exact matches, not substrings within other runes.
func newSecretReplacer(secrets []string) *strings.Replacer {
	var oldNew []string
	for _, part := range secrets {
		if len(part) == 0 {
			continue
		}
		// Only replace if the secret is not a single character or
		// is not a unicode symbol
		// This is a simple heuristic; for more robust
		// handling, use regex or a custom replacer
		if len([]rune(part)) == 1 && part != " " {
			// skip single unicode
			// characters
			continue
		}
		oldNew = append(oldNew, part)
		oldNew = append(oldNew, "********")
	}
	return strings.NewReplacer(oldNew...)
}

func noOpSecretReplacer(secrets []string) *strings.Replacer {
	return strings.NewReplacer()
}
