package proxy

import (
	"context"

	"github.com/cloudness-io/cloudness/types/enum"
)

func (s *Service) ValidateToken(ctx context.Context, proxy enum.ServerProxy, token string, hostname string) error {
	proxySvc, err := s.Resolve(proxy)
	if err != nil {
		return nil
	}

	return proxySvc.ValidateAPIKeyForDNS01(ctx, token, hostname)
}
