package kube

const (
	defaultK8sCloudnessNamespace = "cloudness"

	defaultRegistryImage = "registry:2.8.3"

	defaultK8sGatewayNamespace = "traefik"
	defaultK8sGatewayService   = "traefik"

	//certificates
	defaultProxyAPIKeySecretName   = "cert-proxy-api-key"
	defaultProxyAPIKeySecretKey    = "proxy-api-key"
	defaultClusterIssuerName       = "wildcard-issuer"
	defaultLetsEncryptServerURL    = "https://acme-v02.api.letsencrypt.org/directory"
	defaultWidlcardCertificateName = "cert-wildcard-certificate"
)
