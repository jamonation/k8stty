package objectmanager

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8stty/internal/pkg/clientset"
)

type namespaceManager ManagerImpl

// NewNamespaceManager returns an objectmanager.Manager interface for namespaces
func NewNamespaceManager(c clientset.K8sClient) Manager {
	return &namespaceManager{Client: c}
}

// Create a namespace
func (k *namespaceManager) Create(ctx context.Context, reqOpts map[string]string) error {
	id := reqOpts["id"]

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
	}
	opts := metav1.CreateOptions{}
	if _, err := k.Client.Clientset.CoreV1().Namespaces().Create(ctx, ns, opts); err != nil {
		return fmt.Errorf("error creating namespace: %v", err)
	}
	return nil
}

// Delete a namespace
func (k *namespaceManager) Delete(ctx context.Context, id string) error {
	var zero int64 = 0
	var deletionPolicy = metav1.DeletePropagationBackground
	opts := metav1.DeleteOptions{
		GracePeriodSeconds: &zero,
		PropagationPolicy:  &deletionPolicy,
	}
	if err := k.Client.Clientset.CoreV1().Namespaces().Delete(ctx, id, opts); err != nil {
		return fmt.Errorf("error deleting namespace: %v", err)
	}
	return nil
}
