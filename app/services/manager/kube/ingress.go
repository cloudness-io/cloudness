package kube

import (
	"context"

	"github.com/cloudness-io/cloudness/types"

	"github.com/rs/zerolog/log"
	// Or the appropriate API version
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	v1 "sigs.k8s.io/gateway-api/applyconfiguration/apis/v1"
)

func (m *K8sManager) AddHttpRoute(ctx context.Context, server *types.Server, namespace, name, service string, port int32, host string, httpScheme string) error {
	gatewayClientset, err := m.getGatewayClient(ctx, server)
	if err != nil {
		return err
	}

	gatewaySection := "websecure"
	if httpScheme == "http" {
		gatewaySection = "web"
	}

	httpRouteApply := v1.HTTPRoute(name, namespace).WithSpec(
		v1.HTTPRouteSpec().
			WithParentRefs(v1.ParentReference().WithNamespace(gatewayv1.Namespace(defaultK8sGatewayNamespace)).WithName(defaultK8sGatewayName).WithSectionName(gatewayv1.SectionName(gatewaySection))).
			WithHostnames(gatewayv1.Hostname(host)).
			WithRules(
				v1.HTTPRouteRule().
					WithMatches(v1.HTTPRouteMatch().WithPath(v1.HTTPPathMatch().WithType(gatewayv1.PathMatchPathPrefix).WithValue("/"))).
					WithBackendRefs(
						v1.HTTPBackendRef().WithName(gatewayv1.ObjectName(service)).WithNamespace(gatewayv1.Namespace(namespace)).WithPort(gatewayv1.PortNumber(port)),
					),
			),
	)

	log.Ctx(ctx).Debug().Any("httproute", httpRouteApply).Msg("HTTPRoute")
	if _, err := gatewayClientset.GatewayV1().HTTPRoutes(namespace).Apply(ctx, httpRouteApply, applyOptions); err != nil {
		return err
	}
	return nil
}

func (m *K8sManager) RemoveHttpRoute(ctx context.Context, server *types.Server, namespace, name string) error {
	gatewayClientSet, err := m.getGatewayClient(ctx, server)
	if err != nil {
		return err
	}

	if err := gatewayClientSet.GatewayV1().HTTPRoutes(namespace).Delete(ctx, name, metav1.DeleteOptions{}); err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}
