package enum

type DNSProvider string

const (
	DNSProviderNone       DNSProvider = "none"
	DNSProviderCloudflare DNSProvider = "cloudflare"
)

var DNSProviders = sortEnum([]DNSProvider{
	DNSProviderNone,
	DNSProviderCloudflare,
})

var DNSProvidersStr = []string{
	string(DNSProviderNone),
	string(DNSProviderCloudflare),
}

func DNSProviderFromString(s string) DNSProvider {
	switch s {
	case string(DNSProviderNone):
		return DNSProviderNone
	case string(DNSProviderCloudflare):
		return DNSProviderCloudflare
	default:
		return ""
	}
}
