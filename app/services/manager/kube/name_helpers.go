package kube

//certificate

func certIssuerName(certKey string) string {
	return certKey + "-issuer"
}

func certLetsEncryptKey(certKey string) string {
	return certKey + "-letsencrypt-key"
}

func certificateName(certKey string) string {
	return certKey + "-certificate"
}

func certTLSSecretName(certKey string) string {
	return certKey + "-certificate-tls"
}

func certDNSProxyAPISecretName(certKey string) string {
	return certKey + "-api-token-secret"
}

func gatewayReferenctGrant(key string) string {
	return key + "-reference-grant"
}

func listenerSetName(key string) string {
	return key + "-listener"
}

func listenerSectionName(key string) string {
	return key + "-listener-section"
}

func httpRouteName(key string) string {
	return key + "-http-route"
}
