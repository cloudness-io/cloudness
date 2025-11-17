package kube

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	appsapplyv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	coreapplyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	metaapplyv1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/util/retry"
)

func (m *K8sManager) AddOrUpdateRegistry(ctx context.Context, server *types.Server, name string, size int64, enableNodePort bool) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}

	var nodePort int32 = 0
	if enableNodePort {
		nodePort = 30050
	}
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := m.createOrUpdateStatefulset(ctx, client, m.getRegistryStatefulSet(name, size)); err != nil {
			return err
		}

		if err := m.createOrUpdateService(ctx, client, m.getRegistryServiceConfig(name, nodePort)); err != nil {
			return err
		}
		return nil
	})
}

func (m *K8sManager) RemoveRegistry(ctx context.Context, server *types.Server, name string) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		if err := m.deleteStatefulset(ctx, client, defaultK8sCloudnessNamespace, name); err != nil {
			return err
		}
		if err := m.deleteService(ctx, client, defaultK8sCloudnessNamespace, name); err != nil {
			return err
		}
		if err := m.deletePVC(ctx, client, defaultK8sCloudnessNamespace, fmt.Sprintf("data-%s-0", name)); err != nil {
			return err
		}
		return nil
	})
}

// helpers
func (m *K8sManager) getRegistryStatefulSet(name string, size int64) *appsapplyv1.StatefulSetApplyConfiguration {
	label := map[string]string{
		"app": name,
	}
	probe := coreapplyv1.Probe().
		WithHTTPGet(coreapplyv1.HTTPGetAction().WithPath("/").WithPort(intstr.FromInt(5000))).
		WithFailureThreshold(5).
		WithInitialDelaySeconds(5).
		WithPeriodSeconds(5).
		WithTimeoutSeconds(6)

	resource := map[v1.ResourceName]resource.Quantity{
		v1.ResourceStorage: resource.MustParse(fmt.Sprintf("%dGi", size)),
	}

	sst := appsapplyv1.StatefulSet(name, defaultK8sCloudnessNamespace).
		WithLabels(label).
		WithSpec(
			appsapplyv1.StatefulSetSpec().
				WithServiceName(name).
				WithReplicas(1).
				WithSelector(metaapplyv1.LabelSelector().WithMatchLabels(label)).
				WithTemplate(
					coreapplyv1.PodTemplateSpec().
						WithLabels(label).
						WithSpec(
							coreapplyv1.PodSpec().WithContainers(
								coreapplyv1.Container().
									WithName(name).
									WithImage(defaultRegistryImage).
									WithImagePullPolicy(v1.PullIfNotPresent).
									WithPorts(
										coreapplyv1.ContainerPort().
											WithContainerPort(5000).
											WithName(name),
									).
									WithLivenessProbe(probe).
									WithReadinessProbe(probe).
									WithStartupProbe(probe).
									WithVolumeMounts(coreapplyv1.VolumeMount().WithName("data").WithMountPath("/var/lib/registry")),
							),
						),
				).
				WithVolumeClaimTemplates(
					coreapplyv1.PersistentVolumeClaim(name, defaultK8sCloudnessNamespace).
						WithName("data").
						WithSpec(
							coreapplyv1.PersistentVolumeClaimSpec().
								WithAccessModes(v1.ReadWriteOnce).
								WithResources(coreapplyv1.VolumeResourceRequirements().WithRequests(resource))),
				),
		)

	return sst
}

func (m *K8sManager) getRegistryServiceConfig(name string, nodePort int32) *coreapplyv1.ServiceApplyConfiguration {
	label := map[string]string{
		"app": name,
	}

	svc := coreapplyv1.Service(name, defaultK8sCloudnessNamespace).
		WithLabels(label).
		WithSpec(
			coreapplyv1.ServiceSpec().
				WithType(v1.ServiceTypeNodePort).
				WithSelector(label).
				WithPorts(
					coreapplyv1.ServicePort().
						WithName(name).
						WithAppProtocol(string(v1.ProtocolTCP)).
						WithPort(5000).
						WithTargetPort(intstr.FromInt(5000)),
				),
		)

	if nodePort > 0 {
		svc.Spec.Ports[0].WithNodePort(int32(nodePort))
	}

	return svc
}
