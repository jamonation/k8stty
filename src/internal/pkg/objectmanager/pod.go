package objectmanager

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"

	"k8stty/internal/pkg/clientset"
)

type podManager ManagerImpl

// NewPodManager returns an objectmanager.Manager interface for pods
func NewPodManager(c clientset.K8sClient) Manager {
	return &podManager{Client: c}
}

// Create a pod
func (k *podManager) Create(ctx context.Context, reqInfo map[string]string) error {
	var false bool = false

	// presence of these is checked in PodServer.CreatePod() before this function is called
	id := reqInfo["id"]
	image := reqInfo["image"]
        runtimeClass := "sysbox-runc"

	// pods use sysbox as their runtime, so need the io.kubernetes.cri-o.userns-mode annotation
	// along with the RuntimeClassName set to sysbox-runc
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: id,
			Annotations: map[string]string{
				"io.kubernetes.cri-o.userns-mode": "auto:size=65536",
			},
		},
		Spec: corev1.PodSpec{
			// HostAliases are used to override the kube-dns resolver response of 169.254.169.254
			// For some reason 169.254.169.254 doesn't seem directly reachable from pods
			HostAliases: []corev1.HostAlias{
				{
					IP:        "10.254.4.10",
					Hostnames: []string{"metadata.google.internal", "metadata"},
				},
			},
			Containers: []corev1.Container{
				{
					Name:    id,
					Image:   image,
					Stdin:   true,
					ImagePullPolicy: "Always",
					Resources: corev1.ResourceRequirements{
						Limits: corev1.ResourceList{
							"memory": resource.MustParse("600M"),
						},
					},
				},
			},
			AutomountServiceAccountToken: &false, // Don't expose account token in /var/run
			EnableServiceLinks:           &false, // Don't expose env variables
			ServiceAccountName:           "default",
			Hostname:                     id[:8], //use first 8 characters of ID as name
			RuntimeClassName:             &runtimeClass,
			DNSPolicy:                    "None",
			DNSConfig: &corev1.PodDNSConfig{
				Nameservers: []string{"8.8.8.8", "8.8.4.4"},
			},
		},
	}

	opts := metav1.CreateOptions{}
	pod, err := k.Client.Clientset.CoreV1().Pods(id).Create(ctx, pod, opts)
	if err != nil {
		return fmt.Errorf("error creating pod: %v", err)
	}

	// copied from https://miminar.fedorapeople.org/_preview/openshift-enterprise/registry-redeploy/go_client/executing_remote_processes.html
	// But there has to be a better way?
	watcher, err := k.Client.Clientset.CoreV1().Pods(id).Watch(
		context.TODO(),
		metav1.SingleObject(pod.ObjectMeta),
	)
	if err != nil {
		return fmt.Errorf("error watching pod status: %v", err)
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		switch event.Type {
		case watch.Error:
			// this is usually noise, but leaving it here in case actual errors show up
			status := event.Object.(*metav1.Status)
			for _, detail := range status.Details.Causes {
				if detail.Message != "unable to decode an event from the watch stream: http2: response body closed" {
					log.Printf("pod status error: %v", detail)
					return fmt.Errorf("error watching pod status: %s", detail.Type)
				}
			}
		case watch.Modified:
			pod = event.Object.(*corev1.Pod)

			// If pod contains a status condition Ready == True, stop watching.
			// Status.Conditions is a slice, so as inelegant as it is, it needs ranging over
			// every time an event occurs in order to detect a ready condition
			for _, cond := range pod.Status.Conditions {
				if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
					log.Printf("pod %s ready\n", id)
					watcher.Stop()
				} else {
					// the rest of these statuses are mostly for logging and erroring out with
					// meaningful error messages
					for _, status := range pod.Status.ContainerStatuses {
						if status.State.Terminated != nil {
							log.Printf("Termination state: %#v", status.State)
							watcher.Stop()
							return fmt.Errorf("pod was terminated: %s", status.State.Terminated.Reason)
						}
						if status.LastTerminationState.Terminated != nil {
							log.Printf("error starting pod: %v - %v\n", status.LastTerminationState.Terminated.Reason, status.LastTerminationState.Terminated.Message)
							watcher.Stop()
							return fmt.Errorf("error starting pod: %s", status.LastTerminationState.Terminated.Reason)
						}
						if status.State.Waiting != nil {
							if status.State.Waiting.Reason == "ErrImagePull" || status.State.Waiting.Reason == "ImagePullbackOff" {
								log.Printf("error pulling image %s: %s", image, status.State.Waiting.Message)
								watcher.Stop()
								return fmt.Errorf("error pulling image %s", image)
							}
							//log.Printf("waiting for pod: %s\n", status.State.Waiting.Reason)
						}
						// if status.State.Running != nil {
						// 	log.Printf("pod is running: %v", status.State.Running.StartedAt)
						// 	watcher.Stop()
						// }
						//log.Printf("what is this status? %#v\n", status.State)
					}
				}
			}
		default:
			log.Printf("unexpected event type but trying to continue: %v", event)
		}
	}

	return nil
}

// Delete is not imlpemented since deleting the containing namespace will remove everything
func (k *podManager) Delete(ctx context.Context, id string) error {
	return nil
}
