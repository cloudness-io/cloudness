package kube

const (
	DefaultK8sCloudnessNamespace       = "cloudness"
	DefaultK8sCloudnessName            = "cloudness"
	DefaultK8sCloudnessService         = "cloudness"
	DefaultK8sCloudnessPort      int32 = 8000

	defaultRegistryImage = "registry:2.8.3"

	defaultK8sGatewayNamespace = "traefik"
	defaultK8sGatewayName      = "traefik"
	defaultK8sGatewayService   = "traefik"

	//certificates
	defaultProxyAPIKeySecretKey    = "token"
	defaultLetsEncryptServerURL    = "https://acme-staging-v02.api.letsencrypt.org/directory"
	defaultWidlcardCertificateName = "cert-wildcard-certificate"
)

var defaultCertificateLabel = map[string]string{
	"traefik.ingress.kubernetes.io/tls.cert": "true",
}
