package proxy

import "context"

type proxy interface {
	ValidateAPIKeyForDNS01(ctx context.Context, apitoken string, zone string) error
}
