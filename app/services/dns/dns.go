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

func (s *Service) ValidateHost(ctx context.Context, hostname string, ipV4 string, dnsServers string, proxy enum.ServerProxy) error {
	for dnsServer := range strings.SplitSeq(dnsServers, ",") {
		ips, err := s.resolveARecord(dnsServer, hostname)
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Any("host", hostname).Any("dns", dnsServer).Msg("Error looking up host")
			continue
		}

		switch proxy {
		case enum.ServerProxyCloudflare:
			if slices.ContainsFunc(ips, s.isCloudflareProxied) {
				return nil
			}
		case enum.ServerProxyNone:
			if slices.Contains(ips, ipV4) {
				return nil
			}
		}
	}

	return usererror.BadRequest("DNS validation failed, no DNS entry found for host pointing to server IP")
}
