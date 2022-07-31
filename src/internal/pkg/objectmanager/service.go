package objectmanager

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8stty/internal/pkg/clientset"
)

type serviceManager ManagerImpl

// NewServiceManager returns an objectmanager.Manager interface for pods
func NewServiceManager(c clientset.K8sClient) Manager {
	return &serviceManager{Client: c}
}

func (k *serviceManager) Create(ctx context.Context, reqOpts map[string]string) error {
	id := reqOpts["id"]
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
		},
		// TODO actually fill out the spec
		Spec: corev1.ServiceSpec{},
	}
	opts := metav1.CreateOptions{}
	if _, err := k.Client.Clientset.CoreV1().Services(id).Create(ctx, service, opts); err != nil {
		return fmt.Errorf("error creating service: %v", err)
	}
	return nil
}

// Delete is not used since the containing namespace will remove everything
func (k *serviceManager) Delete(ctx context.Context, id string) error {
	return nil
}
