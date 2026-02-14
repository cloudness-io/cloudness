package kube

import (
	"context"
	"fmt"

	"github.com/cloudness-io/cloudness/types"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/util/retry"
)

func (m *K8sManager) DeleteResources(ctx context.Context, server *types.Server, namespace string, identifier string) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}
	dynamicClient, err := m.getDynamicClient(ctx, server)
	if err != nil {
		return err
	}

	httpGvr := schema.GroupVersionResource{
		Group:    "gateway.networking.k8s.io",
		Version:  "v1",
		Resource: "httproutes",
	}
	label := fmt.Sprintf("app.kubernetes.io/instance=%s", identifier)

	deleteOption := metav1.DeleteOptions{}
	listOption := metav1.ListOptions{
		LabelSelector: label,
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err = dynamicClient.Resource(httpGvr).Namespace(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.AppsV1().StatefulSets(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.AppsV1().Deployments(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		return nil
	})
}

func (m *K8sManager) DeleteNamespace(ctx context.Context, server *types.Server, namespace string) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}

	err = client.CoreV1().Namespaces().Delete(ctx, namespace, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	return nil
}

func (m *K8sManager) DeleteVolume(ctx context.Context, server *types.Server, volume *types.Volume) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}

	err = client.CoreV1().PersistentVolumeClaims(volume.ParentSlug).Delete(ctx, volume.GetIdentifierStr(), metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}

	return nil
}

func (m *K8sManager) DeleteApplication(ctx context.Context, server *types.Server, app *types.Application) error {
	client, err := m.getInterface(ctx, server)
	if err != nil {
		return err
	}
	dynamicClient, err := m.getDynamicClient(ctx, server)
	if err != nil {
		return err
	}

	namespace := app.ParentSlug
	httpGvr := schema.GroupVersionResource{
		Group:    "gateway.networking.k8s.io",
		Version:  "v1",
		Resource: "httproutes",
	}
	label := fmt.Sprintf("app.kubernetes.io/instance=%s", app.GetIdentifierStr())

	deleteOption := metav1.DeleteOptions{}
	listOption := metav1.ListOptions{
		LabelSelector: label,
	}

	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		err = dynamicClient.Resource(httpGvr).Namespace(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		//services
		services, err := client.CoreV1().Services(namespace).List(ctx, listOption)
		if err != nil {
			return err
		}
		for _, service := range services.Items {
			err = client.CoreV1().Services(namespace).Delete(ctx, service.Name, deleteOption)
			if err != nil && !errors.IsNotFound(err) {
				return err
			}
		}

		err = client.AppsV1().StatefulSets(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.AppsV1().Deployments(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.CoreV1().ConfigMaps(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.CoreV1().Secrets(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.RbacV1().RoleBindings(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.RbacV1().Roles(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}

		err = client.CoreV1().ServiceAccounts(namespace).DeleteCollection(ctx, deleteOption, listOption)
		if err != nil {
			return err
		}
		return nil
	})
}
