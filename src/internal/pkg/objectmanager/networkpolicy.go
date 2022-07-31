package objectmanager

import (
	"context"
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"k8stty/internal/pkg/clientset"
)

type networkpolicyManager ManagerImpl

// NewNetworkpolicyManager returns an objectmanager.Manager interface for pods
func NewNetworkpolicyManager(c clientset.K8sClient) Manager {
	return &networkpolicyManager{Client: c}
}

func (k *networkpolicyManager) Create(ctx context.Context, reqInfo map[string]string) error {
	id := reqInfo["id"]
	tcp := v1.ProtocolTCP
	udp := v1.ProtocolUDP
	dns := intstr.IntOrString{IntVal: 53}
	http := intstr.IntOrString{IntVal: 80}
	https := intstr.IntOrString{IntVal: 443}
	httpAlt := intstr.IntOrString{IntVal: 8080}

	allowedPorts := []netv1.NetworkPolicyPort{
		{
			Protocol: &tcp,
			Port:     &http,
		},
		{
			Protocol: &tcp,
			Port:     &httpAlt,
		},
		{
			Protocol: &tcp,
			Port:     &https,
		},
		{
			Protocol: &tcp,
			Port:     &dns,
		},
		{
			Protocol: &udp,
			Port:     &dns,
		},
	}

	defaultDenyPolicy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deny-egress",
			Namespace: id,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			PolicyTypes: []netv1.PolicyType{"Ingress", "Egress"},
		},
	}
	allowPolicy := &netv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      id,
			Namespace: id,
		},
		Spec: netv1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{},
			Egress: []netv1.NetworkPolicyEgressRule{
				{
					Ports: allowedPorts,
				},
			},
			PolicyTypes: []netv1.PolicyType{"Egress"},
		},
	}
	opts := metav1.CreateOptions{}
	log.Printf("creating defaultDenyPolicy\n")
	if _, err := k.Client.Clientset.NetworkingV1().NetworkPolicies(id).Create(ctx, defaultDenyPolicy, opts); err != nil {
		return fmt.Errorf("error creating network policy: %v", err)
	}
	log.Printf("creating allowPolicy\n")
	if _, err := k.Client.Clientset.NetworkingV1().NetworkPolicies(id).Create(ctx, allowPolicy, opts); err != nil {
		return fmt.Errorf("error creating network policy: %v", err)
	}
	return nil
}

// Delete is not used since the containing namespace will remove everything
func (k *networkpolicyManager) Delete(ctx context.Context, id string) error {
	return nil
}
