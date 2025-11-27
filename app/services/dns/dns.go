package dns

import (
	"context"
	"slices"
	"strings"

	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/types/enum"

	"github.com/rs/zerolog/log"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) ValidateHost(ctx context.Context, hostname string, ipV4 string, dnsServers string, proxy enum.DNSProvider) error {
	for dnsServer := range strings.SplitSeq(dnsServers, ",") {
		ips, err := s.resolveARecord(dnsServer, hostname)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any("host", hostname).Any("dns", dnsServer).Msg("Error looking up host")
			continue
		}

		switch proxy {
		case enum.DNSProviderCloudflare:
			if slices.ContainsFunc(ips, s.isCloudflareProxied) {
				return nil
			}
		case enum.DNSProviderNone:
			if slices.Contains(ips, ipV4) {
				return nil
			}
		}
		log.Ctx(ctx).Error().Str("host", hostname).Any("dns", dnsServers).Str("ipv4", ipV4).Any("proxy", proxy).Any("ips", ips).Msg("Error validating DNS")
	}

	return usererror.BadRequest("DNS validation failed, no DNS entry found for host pointing to server IP")
}
