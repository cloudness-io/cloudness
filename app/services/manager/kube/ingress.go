package kube

import (
	"context"
	"strings"

	"github.com/cloudness-io/cloudness/types"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	xlsv1alpha1 "sigs.k8s.io/gateway-api/apisx/v1alpha1"
	v1 "sigs.k8s.io/gateway-api/applyconfiguration/apis/v1"
)

func (m *K8sManager) AddHttpRoute(ctx context.Context, server *types.Server, namespace, key, service string, port int32, host string, httpScheme string) error {
	gatewayClientset, err := m.getGatewayClient(ctx, server)
	if err != nil {
		return err
	}

	isWebListener := false
	domain, err := server.GetDomain()
	if err != nil {
		return err
	}

	if strings.HasSuffix(host, domain.Hostname) && httpScheme == "http" {
		isWebListener = true
	}

	httpRouteApply := v1.HTTPRoute(httpRouteName(key), namespace).WithSpec(
		v1.HTTPRouteSpec().
			WithHostnames(gatewayv1.Hostname(host)).
			WithRules(
				v1.HTTPRouteRule().
					WithMatches(v1.HTTPRouteMatch().WithPath(v1.HTTPPathMatch().WithType(gatewayv1.PathMatchPathPrefix).WithValue("/"))).
					WithBackendRefs(
						v1.HTTPBackendRef().WithName(gatewayv1.ObjectName(service)).WithNamespace(gatewayv1.Namespace(namespace)).WithPort(gatewayv1.PortNumber(port)),
					),
			),
	)

	if isWebListener {
		httpRouteApply.Spec = httpRouteApply.Spec.
			WithParentRefs(
				v1.ParentReference().
					WithNamespace(gatewayv1.Namespace(DefaultK8sGatewayNamespace)).
					WithName(DefaultK8sGatewayName).
					WithSectionName(gatewayv1.SectionName("web")),
			)
	} else {
		httpRouteApply.Spec = httpRouteApply.Spec.WithParentRefs(
			v1.ParentReference().
				WithGroup(xlsv1alpha1.GroupName).
				WithKind(gatewayv1.Kind("XListenerSet")).
				WithName(gatewayv1.ObjectName(listenerSetName(key))).
				WithSectionName(gatewayv1.SectionName(listenerSectionName(key))),
		)
	}

	_, err = gatewayClientset.GatewayV1().HTTPRoutes(namespace).Apply(ctx, httpRouteApply, applyOptions)
	return err
}

func (m *K8sManager) RemoveHttpRoute(ctx context.Context, server *types.Server, namespace, key string) error {
	gatewayClientSet, err := m.getGatewayClient(ctx, server)
	if err != nil {
		return err
	}

	if err := gatewayClientSet.GatewayV1().HTTPRoutes(namespace).Delete(ctx, httpRouteName(key), metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
