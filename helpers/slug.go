package helpers

import "strings"

func GenerateSlug(len int) string {
	return RandomLower(len)
}

func Slugify(prefix string, name string) string {
	const suffixLen = 8
	const maxLen = 63
	reserve := suffixLen + 1

	base := Normalize(name)
	if base == "" {
		base = RandomLower(8)
	}

	pref := Normalize(prefix)
	// Reserve room: suffix + hyphen.
	if pref != "" {
		reserve += len(pref) + 1
	}

	if reserve >= maxLen {
		// are you really dumb?
		return RandomLower(maxLen)
	}

	// if the name is already greater than max length, reserving space for slug since k8s has a limit
	if len(base) > maxLen-reserve {
		base = base[:maxLen-reserve]
		base = strings.Trim(base, "-")
		if base == "" {
			base = RandomLower(8)
		}
	}

	slug := base + "-" + RandomLower(suffixLen)
	if pref != "" {
		slug = pref + "-" + slug
	}

	return slug
}
