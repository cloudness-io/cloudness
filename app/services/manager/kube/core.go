package kube

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsapplyv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	coreapplyv1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
)

var applyOptions metav1.ApplyOptions = metav1.ApplyOptions{
	FieldManager: "cloudness",
	Force:        true,
}

func (m *K8sManager) createOrUpdateStatefulset(ctx context.Context, client kubernetes.Interface, statefulset *appsapplyv1.StatefulSetApplyConfiguration) error {
	apps := client.AppsV1().StatefulSets(*statefulset.Namespace)

	_, err := apps.Apply(ctx, statefulset, applyOptions)
	return err
}

func (m *K8sManager) deleteStatefulset(ctx context.Context, client kubernetes.Interface, namespace string, name string) error {
	apps := client.AppsV1().StatefulSets(namespace)

	err := apps.Delete(ctx, name, metav1.DeleteOptions{})
	if err == nil {
		return nil
	}

	if !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

func (m *K8sManager) createOrUpdateService(ctx context.Context, client kubernetes.Interface, service *coreapplyv1.ServiceApplyConfiguration) error {
	svcs := client.CoreV1().Services(*service.Namespace)

	_, err := svcs.Apply(ctx, service, applyOptions)

	return err
}

func (m *K8sManager) deleteService(ctx context.Context, client kubernetes.Interface, namespace string, name string) error {
	svcs := client.CoreV1().Services(namespace)

	err := svcs.Delete(ctx, name, metav1.DeleteOptions{})
	if err == nil {
		return nil
	}

	if !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}

func (m *K8sManager) deletePVC(ctx context.Context, client kubernetes.Interface, namespace string, name string) error {
	pvc := client.CoreV1().PersistentVolumeClaims(namespace)

	err := pvc.Delete(ctx, name, metav1.DeleteOptions{})
	if err == nil {
		return nil
	}

	if !apierrors.IsNotFound(err) {
		return err
	}

	return nil
}
