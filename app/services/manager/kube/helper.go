package kube

import (
	"regexp"
	"strings"

	"github.com/cloudness-io/cloudness/app/usererror"
)

// Domain helpers
func (m *K8sManager) getKubeKeyForDomain(subdomain, hostname string) (string, error) {
	if hostname == "" {
		return "", usererror.BadRequest("Invalid hostname: hostname cannot be empty, check if you have added http/https scheme to your domain")
	}

	if subdomain == "*" {
		subdomain = "wildcard"
	}

	// 1. Convert to lowercase
	kubeName := strings.ToLower(subdomain + "-" + hostname)

	// 2. Replace underscores with hyphens (as underscores are not allowed)
	kubeName = strings.ReplaceAll(kubeName, "_", "-")

	// 3. Use a regex to keep only alphanumeric characters and hyphens
	// This removes any special characters that might be in a raw hostname
	reg := regexp.MustCompile("[^a-z0-9-]+")
	kubeName = reg.ReplaceAllString(kubeName, "")

	// 4. Ensure it starts and ends with an alphanumeric character.
	// Trim leading and trailing hyphens.
	kubeName = strings.Trim(kubeName, "-")

	// Optional: Limit length to 253 characters (standard DNS limit)
	if len(kubeName) > 253 {
		kubeName = kubeName[:253]
	}

	return kubeName, nil
}

func (m *K8sManager) certificateSecretName(key string) string {
	return key + "-secret"
}
