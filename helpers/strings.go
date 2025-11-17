package helpers

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func CapitalizeString(s string) string {
	return cases.Title(language.English, cases.NoLower).String(s)
}

func CapitalizeSentance(sen string) string {
	ret := make([]string, 0)
	for word := range strings.SplitSeq(sen, " ") {
		ret = append(ret, CapitalizeString(word))
	}
	return strings.Join(ret, " ")
}

func Normalize(text string) string {
	normalized := regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(strings.ToLower(text), "-")
	return strings.Trim(normalized, "-")
}
