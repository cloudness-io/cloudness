package types

import (
	"fmt"
	"net/url"

	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/types/enum"
)

type Server struct {
	ID                            int64            `db:"server_id"                                   json:"-"`
	UID                           int64            `db:"server_uid"                                  json:"uid"`
	Type                          enum.ServerType  `db:"server_type"                                 json:"type"`
	Name                          string           `db:"server_name"                                 json:"name"`
	Description                   string           `db:"server_description"                          json:"description"`
	IPV4                          string           `db:"server_ipv4"                                 json:"ipv4"`
	IPV6                          string           `db:"server_ipv6"                                 json:"ipv6"`
	WildCardDomain                string           `db:"server_wildcard_domain"                      json:"wildcard_domain"`
	DNSProxy                      enum.ServerProxy `db:"server_dns_proxy"                            json:"dns_proxy"`
	ProxyAuthKey                  string           `db:"server_proxy_auth_key"                       json:"-"`
	User                          string           `db:"server_user"                                 json:"user"`
	Port                          int64            `db:"server_port"                                 json:"port"`
	VolumeSupportsOnlineExpansion bool             `db:"server_volume_supports_online_expansion"     json:"volume_supports_online_expansion"`
	VolumeMinSize                 int64            `db:"server_volume_min_size"                      json:"volume_min_size"`
	BuildEnabled                  bool             `db:"server_builder_is_enabled"                   json:"build_enabled"`
	IsBuildServer                 bool             `db:"server_builder_is_build_server"              json:"is_build_server"`
	MaxConcurrentBuilds           int64            `db:"server_builder_max_concurrent_builds"        json:"max_concurrent_builds"`
	MaxCPUPerBuild                float64          `db:"server_builder_max_cpu"                      json:"max_cpu_per_build"`
	MaxMemoryPerBuild             float64          `db:"server_builder_max_memory"                   json:"max_memory_per_build"`

	Created int64 `db:"server_created"          json:"created"`
	Updated int64 `db:"server_updated"          json:"updated"`
}

type ServerDomain struct {
	Hostname string
	Scheme   string
}

func (s *Server) GetDomain() (*ServerDomain, error) {
	if s.WildCardDomain != "" {
		domain, err := url.Parse(s.WildCardDomain)
		if err != nil {
			return nil, err
		}
		return &ServerDomain{
			Hostname: domain.Hostname(),
			Scheme:   domain.Scheme,
		}, nil
	}

	if s.IPV4 != "" {
		return &ServerDomain{
			Hostname: fmt.Sprintf("%s.sslip.io", s.IPV4),
			Scheme:   "http",
		}, nil
	}

	return nil, usererror.BadRequest("No wildcard domain or ipv4 found")
}
