package util

import (
	"os"
	"path/filepath"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
)

func BuildKubeClients() (*rest.Config, *restmapper.DeferredDiscoveryRESTMapper, *discovery.DiscoveryClient, *kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		home := os.Getenv("HOME")
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			config, err = rest.InClusterConfig()
		}
	}
	if err != nil {
		return nil, nil, nil, nil, err
	}
	disco, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	restMapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(disco))
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	return config, restMapper, disco, clientset, nil
}

func DetectNamespace() (string, error) {
	if ns := os.Getenv("NAMESPACE"); ns != "" {
		return ns, nil
	}
	data, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/namespace")
	if err == nil {
		return string(data), nil
	}
	return "default", nil
}
