package kube

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/cloudness-io/cloudness/types"
	"github.com/cloudness-io/cloudness/types/enum"

	cmacme "github.com/cert-manager/cert-manager/pkg/apis/acme/v1"
	cmv1 "github.com/cert-manager/cert-manager/pkg/apis/certmanager/v1"
	cmmeta "github.com/cert-manager/cert-manager/pkg/apis/meta/v1"
	cmclientset "github.com/cert-manager/cert-manager/pkg/client/clientset/versioned"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/util/retry"
	gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"
	gwapiv1b1 "sigs.k8s.io/gateway-api/apis/v1beta1"
	xlsv1alpha1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
	gatewayclientset "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"
)

func (m *K8sManager) RemoveWildcardSSL(ctx context.Context, server *types.Server) error {
	return m.RemoveSSLCertificate(ctx, server, DefaultK8sGatewayNamespace, DefaultWidlcardCertificateKey)
}

func (m *K8sManager) RemoveSSLCertificate(ctx context.Context, server *types.Server, namespace string, certKey string) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}
	cmclient, err := m.getACMEClient(ctx, server)
	if err != nil {
		return err
	}

	gwClient, err := m.getGatewayClient(ctx, server)
	if err != nil {
		return err
	}

	//remove cert resources
	if err := client.CoreV1().Secrets(namespace).Delete(ctx, certLetsEncryptKey(certKey), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := client.CoreV1().Secrets(namespace).Delete(ctx, certDNSProxyAPISecretName(certKey), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := client.CoreV1().Secrets(namespace).Delete(ctx, certTLSSecretName(certKey), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := cmclient.CertmanagerV1().Issuers(namespace).Delete(ctx, certIssuerName(certKey), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}
	if err := cmclient.CertmanagerV1().Certificates(namespace).Delete(ctx, certificateName(certKey), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}

	if err := gwClient.ExperimentalV1alpha1().XListenerSets(namespace).Delete(ctx, listenerSetName(certKey), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (m *K8sManager) AddSSLCertificate(ctx context.Context, server *types.Server, namespace, dns, certKey string, dnsProvider enum.DNSProvider, dnsAuthKey string) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}
	cmclient, err := m.getACMEClient(ctx, server)
	if err != nil {
		return err
	}
	gwClient, err := m.getGatewayClient(ctx, server)
	if err != nil {
		return err
	}
	isWildcardDomain := strings.HasPrefix(dns, "*.")

	isProxied := dnsProvider != enum.DNSProviderNone && dnsAuthKey != ""
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if isProxied {
			// 1: Create Cloudflare proxy api key
			if err := m.createOrUpdateProxyAuthSecret(ctx, client, namespace, certKey, dnsAuthKey); err != nil {
				return err
			}
		}

		// 2. Create ClusterIssuer
		if err := m.createOrUpdateIssuer(ctx, cmclient, namespace, dns, certKey, dnsProvider, dnsAuthKey); err != nil {
			return err
		}

		// 3. Create certificate
		if err := m.createOrUpdateTLSCertificate(ctx, cmclient, namespace, dns, certKey); err != nil {
			return err
		}

		// 4. Add listener set to gateway
		if err := m.addListenerSet(ctx, gwClient, dns, namespace, certKey, isWildcardDomain); err != nil {
			return err
		}

		return nil
	})
}

func (m *K8sManager) AddWildcardDomainWithSSL(ctx context.Context, server *types.Server) error {
	wUrl, err := url.Parse(server.WildCardDomain)
	if err != nil {
		return nil
	}

	if wUrl.Scheme == "https" {
		return m.AddSSLCertificate(ctx, server, DefaultK8sGatewayNamespace, fmt.Sprintf("*.%s", wUrl.Hostname()), DefaultWidlcardCertificateKey, server.DNSProvider, server.DNSProviderAuth)
	}

	return nil
}

func (m *K8sManager) ListCertificates(ctx context.Context, server *types.Server) ([]*types.Certificate, error) {
	cmclient, err := m.getACMEClient(ctx, server)
	if err != nil {
		return nil, err
	}

	certs, err := cmclient.CertmanagerV1().Certificates(DefaultK8sGatewayNamespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	result := make([]*types.Certificate, len(certs.Items))
	for i, cert := range certs.Items {
		result[i] = &types.Certificate{
			Name:       cert.Name,
			Namespace:  cert.Namespace,
			DNSNames:   cert.Spec.DNSNames,
			IssuerRef:  fmt.Sprintf("%s/%s", cert.Spec.IssuerRef.Kind, cert.Spec.IssuerRef.Name),
			SecretName: cert.Spec.SecretName,
			Ready:      string(cert.Status.Conditions[len(cert.Status.Conditions)-1].Status),
		}
		if cert.Status.NotBefore != nil {
			result[i].NotBefore = cert.Status.NotBefore.Time
		}
		if cert.Status.NotAfter != nil {
			result[i].NotAfter = cert.Status.NotAfter.Time
		}
		if cert.Status.RenewalTime != nil {
			result[i].RenewalTime = cert.Status.RenewalTime.Time
		}
	}

	return result, nil
}

func (m *K8sManager) createOrUpdateIssuer(ctx context.Context, cmclient *cmclientset.Clientset, namespace, dns, certKey string, dnsProvider enum.DNSProvider, dnsAuth string) error {
	issuerExisits := true
	var issuer *cmv1.Issuer
	var err error
	var issuerName = certIssuerName(certKey)

	issuer, err = cmclient.CertmanagerV1().Issuers(namespace).Get(ctx, issuerName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			issuerExisits = false
		} else {
			return err
		}
	}

	if !issuerExisits {
		acmeUrl, acmeEmail := m.configSvc.GetAcmeUrl()
		issuer = &cmv1.Issuer{
			ObjectMeta: metav1.ObjectMeta{
				Name:      issuerName,
				Namespace: namespace,
			},
			Spec: cmv1.IssuerSpec{
				IssuerConfig: cmv1.IssuerConfig{
					ACME: &cmacme.ACMEIssuer{
						Email:  acmeEmail,
						Server: acmeUrl,
						PrivateKey: cmmeta.SecretKeySelector{
							LocalObjectReference: cmmeta.LocalObjectReference{
								Name: certLetsEncryptKey(certKey),
							},
						},
						Solvers: []cmacme.ACMEChallengeSolver{},
					},
				},
			},
		}
	}

	switch dnsProvider {
	case enum.DNSProviderCloudflare:
		issuer.Spec.IssuerConfig.ACME.Solvers = []cmacme.ACMEChallengeSolver{
			{
				DNS01: &cmacme.ACMEChallengeSolverDNS01{
					Cloudflare: &cmacme.ACMEIssuerDNS01ProviderCloudflare{
						Email: "selfhost@cloudness.io",
						APIToken: &cmmeta.SecretKeySelector{
							LocalObjectReference: cmmeta.LocalObjectReference{
								Name: certDNSProxyAPISecretName(certKey),
							},
							Key: defaultProxyAPIKeySecretKey,
						},
					},
				},
			},
		}
	default:
		issuer.Spec.IssuerConfig.ACME.Solvers = []cmacme.ACMEChallengeSolver{
			{
				HTTP01: &cmacme.ACMEChallengeSolverHTTP01{
					GatewayHTTPRoute: &cmacme.ACMEChallengeSolverHTTP01GatewayHTTPRoute{
						ParentRefs: []gwapiv1.ParentReference{
							{
								Name:      DefaultK8sGatewayName,
								Namespace: ptr(gwapiv1.Namespace(DefaultK8sGatewayNamespace)),
							},
						},
					},
				},
			},
		}
	}

	if !issuerExisits {
		if _, err := cmclient.CertmanagerV1().Issuers(namespace).Create(ctx, issuer, metav1.CreateOptions{}); err != nil {
			return err
		}
	} else {
		if _, err = cmclient.CertmanagerV1().Issuers(namespace).Update(ctx, issuer, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (m *K8sManager) createOrUpdateProxyAuthSecret(ctx context.Context, client kubernetes.Interface, namespace, certKey, dnsAuthValue string) error {
	// 1. Create secret for dns provider auth
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      certDNSProxyAPISecretName(certKey),
			Namespace: namespace,
		},
		StringData: map[string]string{
			defaultProxyAPIKeySecretKey: dnsAuthValue,
		},
		Type: v1.SecretTypeOpaque,
	}

	if _, err := client.CoreV1().Secrets(namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
		if errors.IsAlreadyExists(err) {
			if _, err := client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (m *K8sManager) createOrUpdateTLSCertificate(ctx context.Context, cmclient *cmclientset.Clientset, namespace, dns, certKey string) error {
	exisitingCert, err := cmclient.CertmanagerV1().Certificates(namespace).Get(ctx, certificateName(certKey), metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	cert := &cmv1.Certificate{
		ObjectMeta: metav1.ObjectMeta{
			Name:      certificateName(certKey),
			Namespace: namespace,
		},
		Spec: cmv1.CertificateSpec{
			DNSNames:   []string{dns},
			SecretName: certTLSSecretName(certKey),
			IssuerRef: cmmeta.ObjectReference{
				Name: certIssuerName(certKey),
				Kind: "Issuer",
			},
		},
	}
	if errors.IsNotFound(err) {
		if _, err = cmclient.CertmanagerV1().Certificates(namespace).Create(ctx, cert, metav1.CreateOptions{}); err != nil {
			return err
		}
	} else {
		cert.ObjectMeta.ResourceVersion = exisitingCert.ObjectMeta.ResourceVersion
		if cert.Annotations == nil {
			cert.Annotations = make(map[string]string)
		}
		cert.Annotations["cert-manager.io/force-renewal"] = time.Now().Format(time.RFC3339)
		if _, err = cmclient.CertmanagerV1().Certificates(namespace).Update(ctx, cert, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}

	return nil
}

func (m *K8sManager) createOrUpdateReferenceGrant(ctx context.Context, gwClient *gatewayclientset.Clientset, namespace, certKey string) error {
	refGrant := &gwapiv1b1.ReferenceGrant{
		ObjectMeta: metav1.ObjectMeta{
			Name:      gatewayReferenctGrant(certKey),
			Namespace: namespace,
		},
		Spec: gwapiv1b1.ReferenceGrantSpec{
			From: []gwapiv1b1.ReferenceGrantFrom{
				{
					Group:     gwapiv1.GroupName,
					Kind:      "Gateway",
					Namespace: gwapiv1.Namespace(DefaultK8sGatewayNamespace),
				},
			},
			To: []gwapiv1b1.ReferenceGrantTo{
				{
					Group: "",
					Kind:  "Secret",
				},
			},
		},
	}

	_, err := gwClient.GatewayV1beta1().ReferenceGrants(namespace).Get(ctx, refGrant.Name, metav1.GetOptions{})
	log.Ctx(ctx).Debug().Err(err).Msg("error getting reference grant")
	if err != nil {
		if errors.IsNotFound(err) {
			if _, err = gwClient.GatewayV1beta1().ReferenceGrants(namespace).Create(ctx, refGrant, metav1.CreateOptions{}); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("error creating reference grant")
				return err
			}
		}
		return err
	}

	_, err = gwClient.GatewayV1beta1().ReferenceGrants(namespace).Update(ctx, refGrant, metav1.UpdateOptions{})
	log.Ctx(ctx).Debug().Err(err).Msg("error updating reference grant")
	return err
}

func (m *K8sManager) addCertificateToGateway(ctx context.Context, gwClient *gatewayclientset.Clientset, secretNamespace, secretName string) error {
	gateway, err := gwClient.GatewayV1().Gateways(DefaultK8sGatewayNamespace).Get(ctx, DefaultK8sGatewayName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get gateway: %w", err)
	}

	// Find the HTTP Listener
	var httpsListener *gwapiv1.Listener
	listenerIndex := -1
	for i := range gateway.Spec.Listeners {
		if string(gateway.Spec.Listeners[i].Name) == "websecure" {
			httpsListener = &gateway.Spec.Listeners[i]
			listenerIndex = i
			break
		}
	}

	if httpsListener == nil {
		return fmt.Errorf("HTTPS listener websecure not found in gateway")
	}

	// Check if certificate reference already exisits
	// Check if certificate reference already exists
	certRef := gwapiv1.SecretObjectReference{
		Group:     (*gwapiv1.Group)(strPtr("")),
		Kind:      (*gwapiv1.Kind)(strPtr("Secret")),
		Name:      gwapiv1.ObjectName(secretName),
		Namespace: (*gwapiv1.Namespace)(strPtr(secretNamespace)),
	}

	for _, ref := range httpsListener.TLS.CertificateRefs {
		if string(ref.Name) == secretName && string(*ref.Namespace) == secretNamespace {
			log.Ctx(ctx).Info().Msgf("Certificate reference already exists: %s/%s", secretNamespace, secretName)
			return nil
		}
	}

	// Add the new certificate reference
	httpsListener.TLS.CertificateRefs = append(httpsListener.TLS.CertificateRefs, certRef)
	gateway.Spec.Listeners[listenerIndex] = *httpsListener

	// Update the Gateway
	_, err = gwClient.GatewayV1().Gateways(DefaultK8sGatewayNamespace).Update(ctx, gateway, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update gateway: %w", err)
	}

	return nil
}

func (m *K8sManager) addListenerSet(ctx context.Context, gwclient *gatewayclientset.Clientset, hostname, namespace, key string, isWildcardDomain bool) error {
	allowedRoutesFrom := ptr(gwapiv1.NamespacesFromSame)
	if isWildcardDomain {
		allowedRoutesFrom = ptr(gwapiv1.NamespacesFromAll)
	}
	listener := &xlsv1alpha1.XListenerSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      listenerSetName(key),
			Namespace: namespace,
		},
		Spec: xlsv1alpha1.ListenerSetSpec{
			ParentRef: xlsv1alpha1.ParentGatewayReference{
				Group:     ptr(xlsv1alpha1.Group(gwapiv1.GroupName)),
				Kind:      ptr(xlsv1alpha1.Kind("Gateway")),
				Name:      DefaultK8sGatewayName,
				Namespace: ptr(xlsv1alpha1.Namespace(DefaultK8sGatewayNamespace)),
			},
			Listeners: []xlsv1alpha1.ListenerEntry{
				{
					Name:     xlsv1alpha1.SectionName(listenerSectionName(key)),
					Port:     443,
					Protocol: gwapiv1.HTTPSProtocolType,
					Hostname: ptr(xlsv1alpha1.Hostname(hostname)),
					AllowedRoutes: &gwapiv1.AllowedRoutes{
						Namespaces: &gwapiv1.RouteNamespaces{
							From: allowedRoutesFrom,
						},
					},
					TLS: &gwapiv1.ListenerTLSConfig{
						Mode: ptr(gwapiv1.TLSModeTerminate),
						CertificateRefs: []gwapiv1.SecretObjectReference{
							{
								Kind:      (*gwapiv1.Kind)(strPtr("Secret")),
								Name:      gwapiv1.ObjectName(certTLSSecretName(key)),
								Namespace: (*gwapiv1.Namespace)(strPtr(namespace)),
							},
						},
					},
				},
			},
		},
	}

	log.Ctx(ctx).Debug().Any("listener", listener).Msg("Trying to add listener set")

	_, err := gwclient.ExperimentalV1alpha1().XListenerSets(namespace).Update(ctx, listener, metav1.UpdateOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err := gwclient.ExperimentalV1alpha1().XListenerSets(namespace).Create(ctx, listener, metav1.CreateOptions{})
			if err != nil {
				return fmt.Errorf("failed to create listener set: %w", err)
			}
			return nil
		}
	}

	return nil
}

func (m *K8sManager) removeCertificateFromGateway(ctx context.Context, gwclient *gatewayclientset.Clientset, secretNamespace, secretName string) error {
	gateway, err := gwclient.GatewayV1().Gateways(DefaultK8sGatewayNamespace).Get(ctx, DefaultK8sGatewayName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get gateway: %w", err)
	}

	// Find the HTTP Listener
	var httpsListener *gwapiv1.Listener
	listenerIndex := -1
	for i := range gateway.Spec.Listeners {
		if string(gateway.Spec.Listeners[i].Name) == "websecure" {
			httpsListener = &gateway.Spec.Listeners[i]
			listenerIndex = i
			break
		}
	}

	if httpsListener == nil {
		return fmt.Errorf("HTTPS listener websecure not found in gateway")
	}

	// Remove the certificate reference
	for i, ref := range httpsListener.TLS.CertificateRefs {
		if string(ref.Name) == secretName && string(*ref.Namespace) == secretNamespace {
			httpsListener.TLS.CertificateRefs = append(httpsListener.TLS.CertificateRefs[:i], httpsListener.TLS.CertificateRefs[i+1:]...)
			break
		}
	}

	gateway.Spec.Listeners[listenerIndex] = *httpsListener

	// Update the Gateway
	_, err = gwclient.GatewayV1().Gateways(DefaultK8sGatewayNamespace).Update(ctx, gateway, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("failed to update gateway: %w", err)
	}

	return nil
}

// Helper function to create string pointer
func strPtr(s string) *string {
	return &s
}

func ptr[T any](v T) *T {
	return &v
}
