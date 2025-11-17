package enum

type ServerProxy string

const (
	ServerProxyNone       ServerProxy = "none"
	ServerProxyCloudflare ServerProxy = "cloudflare"
)

var ServerProxys = sortEnum([]ServerProxy{
	ServerProxyNone,
	ServerProxyCloudflare,
})

var ServerProxysStr = []string{
	string(ServerProxyNone),
	string(ServerProxyCloudflare),
}

func ServerProxyFromString(s string) ServerProxy {
	switch s {
	case string(ServerProxyNone):
		return ServerProxyNone
	case string(ServerProxyCloudflare):
		return ServerProxyCloudflare
	default:
		return ""
	}
}
