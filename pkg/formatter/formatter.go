package formatter

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
)

func PrintTable(res interface{}) {
	pods, ok := res.([]corev1.Pod)
	if !ok {
		fmt.Println("No pods found or wrong type")
		return
	}
	fmt.Printf("[DEBUG] PrintTable received %d pods\n", len(pods))
	fmt.Printf("NAME\tNAMESPACE\tSTATUS\n")
	for _, pod := range pods {
		fmt.Printf("%s\t%s\t%s\n", pod.Name, pod.Namespace, string(pod.Status.Phase))
	}
}

func PrintJSON(res interface{}) {
	fmt.Println("[json output]")
}

func PrintYAML(res interface{}) {
	fmt.Println("[yaml output]")
}

func PrintCSV(res interface{}) {
	fmt.Println("[csv output]")
}
