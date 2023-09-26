package kube

import (
	"fmt"

	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NewKubeClient Generating clientset based on kubeconfig path
func NewKubeClient(kubeConfigPath string) (*kubernetes.Clientset, error) {
	kubeConf, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kubeConf, %v", err)
	}
	kubeCli, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeCli, %v", err)
	}
	return kubeCli, nil
}

// NewKubeClientByConfContent Generating clientset based on kubeconfig file content
func NewKubeClientByConfContent(kubeConfig string) (*kubernetes.Clientset, error) {
	kubeConf, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeConfig))
	if err != nil {
		return nil, fmt.Errorf("failed to generate kubeConf, %v", err)
	}
	kubeCli, err := kubernetes.NewForConfig(kubeConf)
	if err != nil {
		return nil, fmt.Errorf("failed to get kubeCli, %v", err)
	}
	return kubeCli, nil
}

// NewDynamicClient based on kube config content generate dynamic client
func NewDynamicClient(kubeConfigPath string) (dynamic.Interface, error) {
	kubeConf, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to generate kubeConf, %v", err)
	}

	dynClient, err := dynamic.NewForConfig(kubeConf)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client, err: %v", err)
	}
	return dynClient, nil
}
