package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const (
	ObjectNameField = "metadata.name"
)

func main() {
	cliset := mustClientset()

	kubeInformers := informers.NewSharedInformerFactoryWithOptions(cliset, 0, informers.WithTweakListOptions(func(options *metav1.ListOptions) {
		options.FieldSelector = fields.Set{ObjectNameField: "extended-resource-demo-2"}.String()
	}))

	podInformer := kubeInformers.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			fmt.Println("add", obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			fmt.Println("update", newObj)
		},
		DeleteFunc: func(obj interface{}) {
			fmt.Println("delete", obj)
		},
	})

	kubeInformers.Start(wait.NeverStop)

	select {}
}