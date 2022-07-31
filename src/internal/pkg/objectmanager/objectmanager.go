package objectmanager

import (
	"context"

	"k8stty/internal/pkg/clientset"
)

// Manager is an abstract interface for all k8s objects managed by k8stty
type Manager interface {
	Create(context.Context, map[string]string) error // map is used since difference create functions take different sets of parameters
	Delete(context.Context, string) error            // single string is `id` of the resource to delete
}

// ManagerImpl is a generic struct for Manager interface
type ManagerImpl struct {
	Client clientset.K8sClient
	Manager
}
