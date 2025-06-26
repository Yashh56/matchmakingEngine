package gameorchestrator

import (
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubeClient() (*kubernetes.Clientset, error) {
	// kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", "C:\\Users\\yashs\\.kube\\config")
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	return os.Getenv("USERPROFILE")
}
