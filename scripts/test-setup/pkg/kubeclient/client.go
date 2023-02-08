package kubeclient

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type client struct {
	kubeConfigPath *string
	clientSet      *kubernetes.Clientset
}

func (c *client) GetCurrentContext() (string, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: *c.kubeConfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.RawConfig()
	if err != nil {
		return "", fmt.Errorf("error getting raw config: %w", err)
	}
	return config.CurrentContext, nil
}

func NewClient(kubeConfigPath *string) (*client, error) {
	client := &client{
		kubeConfigPath: kubeConfigPath,
	}

	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfigPath)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	client.clientSet = clientset

	return client, nil
}
