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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (m *K8sManager) AddWildcardDomainWithSSL(ctx context.Context, server *types.Server) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}
	cmclient, err := m.getACMEClient(ctx, server)
	if err != nil {
		return err
	}

	//remove cert resources
	if server.WildCardDomain == "" || !strings.HasPrefix(server.WildCardDomain, "https") {
		if err := client.CoreV1().Secrets(defaultK8sGatewayNamespace).Delete(ctx, defaultProxyAPIKeySecretName, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return err
		}
		if err := cmclient.CertmanagerV1().Issuers(defaultK8sGatewayNamespace).Delete(ctx, defaultClusterIssuerName, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return err
		}
		if err := cmclient.CertmanagerV1().Certificates(defaultK8sGatewayNamespace).Delete(ctx, defaultWidlcardCertificateName, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
			return err
		}

		return nil
	}

	isProxied := server.DNSProxy != enum.ServerProxyNone && server.ProxyAuthKey != ""
	wUrl, err := url.Parse(server.WildCardDomain)
	if err != nil {
		return nil
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if wUrl.Scheme == "https" && isProxied {
			var err error
			// 1. Create secret for proxy
			secret := &v1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Name:      defaultProxyAPIKeySecretName,
					Namespace: defaultK8sGatewayNamespace,
				},
				StringData: map[string]string{
					defaultProxyAPIKeySecretKey: server.ProxyAuthKey,
				},
				Type: v1.SecretTypeOpaque,
			}

			if _, err = client.CoreV1().Secrets(defaultK8sGatewayNamespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
				if errors.IsAlreadyExists(err) {
					if _, err = client.CoreV1().Secrets(defaultK8sGatewayNamespace).Update(ctx, secret, metav1.UpdateOptions{}); err != nil {
						return err
					}
				} else {
					return err
				}
			}
		} else {
			client.CoreV1().Secrets(defaultK8sGatewayNamespace).Delete(ctx, defaultProxyAPIKeySecretName, metav1.DeleteOptions{})
		}

		// 2. Create ClusterIssuer
		if err := m.checkIssuer(ctx, server); err != nil {
			return err
		}

		// 3. Create certificate
		exisitingCert, err := cmclient.CertmanagerV1().Certificates(defaultK8sGatewayNamespace).Get(ctx, defaultWidlcardCertificateName, metav1.GetOptions{})
		if err != nil && !errors.IsNotFound(err) {
			return err
		}
		cert := &cmv1.Certificate{
			ObjectMeta: metav1.ObjectMeta{
				Name:      defaultWidlcardCertificateName,
				Namespace: defaultK8sGatewayNamespace,
			},
			Spec: cmv1.CertificateSpec{
				DNSNames: []string{
					fmt.Sprintf("*.%s", wUrl.Hostname()),
				},
				SecretName: defaultWidlcardCertificateName,
				IssuerRef: cmmeta.ObjectReference{
					Name: defaultClusterIssuerName,
					Kind: "Issuer",
				},
			},
		}
		if errors.IsNotFound(err) {
			if _, err = cmclient.CertmanagerV1().Certificates(defaultK8sGatewayNamespace).Create(ctx, cert, metav1.CreateOptions{}); err != nil {
				return err
			}
		} else {
			cert.ObjectMeta.ResourceVersion = exisitingCert.ObjectMeta.ResourceVersion
			if cert.Annotations == nil {
				cert.Annotations = make(map[string]string)
			}
			cert.Annotations["cert-manager.io/force-renewal"] = time.Now().Format(time.RFC3339)
			if _, err = cmclient.CertmanagerV1().Certificates(defaultK8sGatewayNamespace).Update(ctx, cert, metav1.UpdateOptions{}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (m *K8sManager) checkIssuer(ctx context.Context, server *types.Server) error {
	issuerExisits := true
	var issuer *cmv1.Issuer
	var err error

	cmclient, err := m.getACMEClient(ctx, server)
	if err != nil {
		return err
	}

	issuer, err = cmclient.CertmanagerV1().Issuers(defaultK8sGatewayNamespace).Get(ctx, defaultClusterIssuerName, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			issuerExisits = false
		} else {
			return err
		}
	}

	if !issuerExisits {
		issuer = &cmv1.Issuer{
			ObjectMeta: metav1.ObjectMeta{
				Name:      defaultClusterIssuerName,
				Namespace: defaultK8sGatewayNamespace,
			},
			Spec: cmv1.IssuerSpec{
				IssuerConfig: cmv1.IssuerConfig{
					ACME: &cmacme.ACMEIssuer{
						Email:  "selfhost@cloudness.io",
						Server: defaultLetsEncryptServerURL,
						PrivateKey: cmmeta.SecretKeySelector{
							LocalObjectReference: cmmeta.LocalObjectReference{
								Name: "cert-wildcard-certificate-key",
							},
						},
						Solvers: []cmacme.ACMEChallengeSolver{},
					},
				},
			},
		}
	}

	switch server.DNSProxy {
	case enum.ServerProxyCloudflare:
		issuer.Spec.IssuerConfig.ACME.Solvers = []cmacme.ACMEChallengeSolver{
			{
				DNS01: &cmacme.ACMEChallengeSolverDNS01{
					Cloudflare: &cmacme.ACMEIssuerDNS01ProviderCloudflare{
						Email: "selfhost@cloudness.io",
						APIToken: &cmmeta.SecretKeySelector{
							LocalObjectReference: cmmeta.LocalObjectReference{
								Name: defaultProxyAPIKeySecretName,
							},
							Key: defaultProxyAPIKeySecretKey,
						},
					},
				},
			},
		}
	default:
		issuer.Spec.IssuerConfig.ACME.Solvers = []cmacme.ACMEChallengeSolver{}
	}

	if !issuerExisits {
		if _, err := cmclient.CertmanagerV1().Issuers(defaultK8sGatewayNamespace).Create(ctx, issuer, metav1.CreateOptions{}); err != nil {
			return err
		}
	} else {
		if _, err = cmclient.CertmanagerV1().Issuers(defaultK8sGatewayNamespace).Update(ctx, issuer, metav1.UpdateOptions{}); err != nil {
			return err
		}
	}
	return nil
}

func (m *K8sManager) ListCertificates(ctx context.Context, server *types.Server) ([]*types.Certificate, error) {
	cmclient, err := m.getACMEClient(ctx, server)
	if err != nil {
		return nil, err
	}

	certs, err := cmclient.CertmanagerV1().Certificates(defaultK8sGatewayNamespace).List(ctx, metav1.ListOptions{})
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
