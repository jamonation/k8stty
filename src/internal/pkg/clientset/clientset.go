// Borrowed from and inspired by https://github.com/a4abhishek/Client-Go-Examples/blob/master/exec_to_pod/exec_to_pod.go

package clientset

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	core "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type ClientManager interface {
	Configure() error
	BuildClientSet() error
}

// K8sClient implements ClientManager and makes
// config & clientsets available for use in a struct
type K8sClient struct {
	V1Client  *core.CoreV1Client
	Clientset *kubernetes.Clientset
	RestCfg   *rest.Config
	ClientManager
}

// NewClient returns an empty K8sClient
// which gets properly configured in Configure()
func NewK8sClient() ClientManager {
	return &K8sClient{
		V1Client:  &core.CoreV1Client{},
		Clientset: &kubernetes.Clientset{},
		RestCfg:   &rest.Config{},
	}
}

// GetConfig from service account token, or ~/.kube/config
func (c *K8sClient) Configure() error {
	var err, errNotInCluster, errKubeConfig error

	c.RestCfg, errNotInCluster = rest.InClusterConfig()
	if errNotInCluster != nil {
		kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
		c.RestCfg, errKubeConfig = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if errKubeConfig != nil {
			err = fmt.Errorf("could not build k8s client config: %v", errKubeConfig)
			return err
		}
	}
	return nil
}

// BuildClientSet takes a config and populates the Core cilent & Clientset
func (c *K8sClient) BuildClientSet() error {
	var err error

	c.V1Client, err = core.NewForConfig(c.RestCfg)
	if err != nil {
		return fmt.Errorf("error creating k8s client: %v", err)
	}

	c.Clientset, err = kubernetes.NewForConfig(c.RestCfg)
	if err != nil {
		return fmt.Errorf("error creating k8s clientset: %v", err)
	}

	return nil
}
