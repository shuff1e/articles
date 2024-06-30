package main

import (
	"os"

	"github.com/spongeprojects/magicconch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func mustClientset() kubernetes.Interface {
	kubeconfig := os.Getenv("KUBECONFIG")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	magicconch.Must(err)

	clientset, err := kubernetes.NewForConfig(config)
	magicconch.Must(err)

	return clientset
}
