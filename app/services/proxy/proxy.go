package proxy

import (
	"github.com/cloudness-io/cloudness/app/usererror"
	"github.com/cloudness-io/cloudness/types/enum"
)

type Service struct {
	proxies map[enum.ServerProxy]proxy
}

func New() *Service {
	cloudflareProxy := newCloudflareProxy()
	return &Service{
		proxies: map[enum.ServerProxy]proxy{
			enum.ServerProxyCloudflare: cloudflareProxy,
		},
	}
}

func (s *Service) Resolve(proxy enum.ServerProxy) (proxy, error) {
	if proxy, ok := s.proxies[proxy]; !ok {
		return nil, usererror.BadRequest("Invalid proxy")
	} else {
		return proxy, nil
	}
}
