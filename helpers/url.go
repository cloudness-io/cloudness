package helpers

import (
	"net/url"
	"strings"
)

// ParseFQDN parses the FQDN into its components and returns schema, subdomain and domain
func ParseFQDN(fqdn string) (scheme string, subdomain string, domain string) {
	u, err := url.Parse(fqdn)
	if err != nil {
		return
	}

	scheme = u.Scheme
	host := u.Hostname()
	parts := strings.Split(host, ".")
	if len(parts) > 2 {
		subdomain = parts[0]
		domain = strings.Join(parts[1:], ".")
	} else {
		domain = host
	}
	if scheme == "" {
		scheme = "http"
	}
	return
}

// GenerateFQDN generates the FQDN from the provided schema, subdomain and domain
func GenerateFQDN(schema string, subdomain string, domain string) string {
	return schema + "://" + subdomain + "." + domain
}

// TrimProtocol removes http/https prefixes and any trailing slashes, preserving host and path.
func TrimProtocol(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ""
	}

	lower := strings.ToLower(trimmed)
	switch {
	case strings.HasPrefix(lower, "https://"):
		trimmed = trimmed[len("https://"):]
	case strings.HasPrefix(lower, "http://"):
		trimmed = trimmed[len("http://"):]
	default:
		if u, err := url.Parse(trimmed); err == nil && u.Scheme != "" {
			// Fallback for other schemes (keeps host/path, drops scheme)
			trimmed = strings.TrimPrefix(trimmed, u.Scheme+"://")
		}
	}

	return strings.TrimRight(trimmed, "/")
}
