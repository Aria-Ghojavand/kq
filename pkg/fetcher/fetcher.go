package fetcher

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func Fetch(ctx context.Context, cfg interface{}, restMapper interface{}, disco interface{}, clientset interface{}, q interface{}, namespace string, allNS bool) (interface{}, error) {
	cs, ok := clientset.(*kubernetes.Clientset)
	if !ok {
		fmt.Println("[DEBUG] clientset type assertion failed")
		return nil, nil
	}
	fmt.Printf("[DEBUG] Fetching pods in namespace: %s\n", namespace)
	pods, err := cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		fmt.Printf("[DEBUG] Error fetching pods: %v\n", err)
		return nil, err
	}
	fmt.Printf("[DEBUG] pods.Items length: %d\n", len(pods.Items))
	return pods.Items, nil
}
