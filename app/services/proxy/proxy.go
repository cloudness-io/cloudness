package proxy

import (
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/types/enum"
)

type Service struct {
	proxies map[enum.DNSProvider]proxy
}

func New() *Service {
	cloudflareProxy := newCloudflareProxy()
	return &Service{
		proxies: map[enum.DNSProvider]proxy{
			enum.DNSProviderCloudflare: cloudflareProxy,
		},
	}
}

func (s *Service) Resolve(proxy enum.DNSProvider) (proxy, error) {
	if proxy, ok := s.proxies[proxy]; !ok {
		return nil, usererror.BadRequest("Invalid proxy")
	} else {
		return proxy, nil
	}
}
