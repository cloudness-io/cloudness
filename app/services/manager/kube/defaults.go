package kube

const (
	defaultK8sCloudnessNamespace = "cloudness"

	defaultRegistryImage = "registry:2.8.3"

	defaultK8sGatewayNamespace = "traefik"
	defaultK8sGatewayService   = "traefik"

	//certificates
	defaultProxyAPIKeySecretName   = "cert-proxy-api-key"
	defaultProxyAPIKeySecretKey    = "token"
	defaultClusterIssuerName       = "wildcard-issuer"
	defaultLetsEncryptServerURL    = "https://acme-staging-v02.api.letsencrypt.org/directory"
	defaultWidlcardCertificateName = "cert-wildcard-certificate"
)

var defaultCertificateLabel = map[string]string{
	"traefik.ingress.kubernetes.io/tls.cert": "true",
}
