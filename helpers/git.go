package helpers

import (
	"regexp"
	"strings"

	"github.com/cloudness-io/cloudness/errors"
)

func SplitGitRepoFullname(fullName string) (owner, repo string, err error) {
	parts := strings.Split(fullName, "/")
	if len(parts) != 2 {
		return "", "", errors.BadRequest("invalid Github repository format: %s", fullName)
	}

	return parts[0], parts[1], nil
}

func SplitGitRepoUrl(url string) (owner, repo string, err error) {
	httpsPattern := `^https:\/\/([\w.-]+)\/([\w-]+)\/([\w.-]+)\.git$`

	httpsRegex := regexp.MustCompile(httpsPattern)
	if matches := httpsRegex.FindStringSubmatch(url); matches != nil {
		return matches[2], matches[3], nil
	}

	return "", "", errors.BadRequest("invalid Github repository format: %s", url)
}

func GetGitHttpUrl(gitUrl string) string {
	return strings.TrimSuffix(gitUrl, ".git")
}

func SanitizeGitUrl(gitUrl string) string {
	if !strings.HasSuffix(gitUrl, ".git") {
		return gitUrl + ".git"
	}
	return gitUrl
}
