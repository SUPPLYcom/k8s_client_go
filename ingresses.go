package k8s_client_go

import (
	"context"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (client *Client) ListIngresses(namespace string) ([]v1beta1.Ingress, error) {
	ingressList, err := client.
		K8sClientSet.
		ExtensionsV1beta1().
		Ingresses(namespace).
		List(context.TODO(), metav1.ListOptions{})

	if(err != nil) {
		return nil, err
	}

	ingresses := make([]v1beta1.Ingress, 0)
	for _, ingress := range ingressList.Items {
		ingresses = append(ingresses, ingress)
	}

	return ingresses, nil

	//ingresses := make([]v1beta1.Ingress, 0)
	//for _, ingress := range ingressList.Items {
	//	if(ingress.Spec.Rules[0].IngressRuleValue.HTTP.Paths[0].Backend.ServiceName == "service-nest-dev") {
	//		ingresses = append(ingresses, ingress)
	//	}
	//}
	//
	//return ingresses, nil
}
