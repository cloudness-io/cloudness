package dns

import (
	"net"

	"github.com/rs/zerolog/log"
)

var cloudflareIPs = []string{
	"173.245.48.0/20",
	"103.21.244.0/22",
	"103.22.200.0/22",
	"103.31.4.0/22",
	"141.101.64.0/18",
	"108.162.192.0/18",
	"190.93.240.0/20",
	"188.114.96.0/20",
	"197.234.240.0/22",
	"198.41.128.0/17",
	"162.158.0.0/15",
	"104.16.0.0/13",
	"104.24.0.0/14",
	"172.64.0.0/13",
	"131.0.72.0/22",
}

func (s *Service) isCloudflareProxied(ipv4 string) bool {
	hostIP := net.ParseIP(ipv4)
	if hostIP == nil {
		log.Error().Msg("Error parsing host IP")
	}
	for _, ip := range cloudflareIPs {
		_, network, err := net.ParseCIDR(ip)
		if err != nil {
			log.Error().Err(err).Msgf("Error parsing CIDR %s: %v", ip, err)
			continue
		}
		if network.Contains(hostIP) {
			return true
		}
	}

	return false
}
