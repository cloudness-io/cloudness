package types

import "time"

type Certificate struct {
	Name        string
	Namespace   string
	DNSNames    []string
	IssuerRef   string
	SecretName  string
	NotBefore   time.Time
	NotAfter    time.Time
	RenewalTime time.Time
	Ready       string
}
