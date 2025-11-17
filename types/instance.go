package types

import "fmt"

type Instance struct {
	ID                    int64  `db:"instance_id"                        json:"-"`
	SuperAdmin            *int64 `db:"instance_super_admin"               json:"-"`
	PublicIPV4            string `db:"instance_public_ipv4"               json:"public_ipv4"`
	PublicIPV6            string `db:"instance_public_ipv6"               json:"public_ipv6"`
	FQDN                  string `db:"instance_fqdn"                      json:"fqdn"`
	UpdateEnabled         bool   `db:"instance_update_enabled"            json:"update_enabled,string"`
	UpdateCheckFrequency  string `db:"instance_update_check_frequency"    json:"update_check_frequency"`
	DNSValidationEnabled  bool   `db:"instance_dns_validation_enabled"    json:"dns_validation_enabled,string"`
	DNSServers            string `db:"instance_dns_servers"               json:"dns_servers"`
	UserSignupEnabled     bool   `db:"instance_user_signup_enabled"       json:"user_signup_enabled,string"`
	DemoUserEnabled       bool   `db:"instance_demo_user_enabled"         json:"demo_user_enabled,string"`
	RegistryEnabled       bool   `db:"instance_registry_enabled"          json:"registry_enabled,string"`
	RegistrySize          int64  `db:"instance_registry_size"             json:"registry_size,string"`
	RegistryMirrorEnabled bool   `db:"instance_registry_mirror_enabled"   json:"registry_mirror_enabled,string"`
	RegistryMirrorSize    int64  `db:"instance_registry_mirror_size"      json:"registry_mirror_size"`
	ExternalScripts       string `db:"instance_external_scripts"          json:"external_scripts"`
	Created               int64  `db:"instance_created"                   json:"-"`
	Updated               int64  `db:"instance_updated"                   json:"-"`
}

type DemoUser struct {
	Email    string `json:"demo_user_email"`
	Password string `json:"-"`
}

// helpers

// GetHttpDomain returns the http domain for this instance
func (i *Instance) GetHttpDomain() string {
	if i.FQDN == "" {
		return fmt.Sprintf("https://%s", i.PublicIPV4)
	}
	return i.FQDN
}
