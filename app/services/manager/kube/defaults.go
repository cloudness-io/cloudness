package kube

const (
	DefaultK8sCloudnessNamespace       = "cloudness"
	DefaultK8sCloudnessName            = "cloudness"
	DefaultK8sCloudnessService         = "cloudness"
	DefaultK8sCloudnessPort      int32 = 8000

	defaultRegistryImage = "registry:2.8.3"

	DefaultK8sGatewayNamespace = "kgateway-system"
	DefaultK8sGatewayName      = "cloudness-gateway"
	DefaultK8sGatewayService   = "cloudness-gateway"

	//certificates
	defaultProxyAPIKeySecretKey   = "token"
	DefaultWidlcardCertificateKey = "cloudness-wildcard"
)
