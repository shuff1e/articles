package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
)

type CmdHandler struct {
}

func (this *CmdHandler) OnAdd(obj interface{}) {
	fmt.Println("Add: ", obj.(*v1.ConfigMap).Name)
}

func (this *CmdHandler) OnUpdate(obj interface{}, newObj interface{}) {
	fmt.Println("Update: ", newObj.(*v1.ConfigMap).Name)
}

func (this *CmdHandler) OnDelete(obj interface{}) {
	fmt.Println("Delete: ", obj.(*v1.ConfigMap).Name)
}

func main() {
	cliset := mustClientset()

	informerFactory := informers.NewSharedInformerFactoryWithOptions(
		cliset,
		0,
		informers.WithNamespace("default"),
	)

	cmGVR := schema.GroupVersionResource{
		Group: "",
		Version: "v1",
		Resource: "configmaps",
	}
	cmInformer, err := informerFactory.ForResource(cmGVR)
	if err != nil {
		panic(err)
	}
	cmInformer.Informer().AddEventHandler(&CmdHandler{})

	podGVR := schema.GroupVersionResource{
		Group: "",
		Version: "v1",
		Resource: "pods",
	}
	_, _ = informerFactory.ForResource(podGVR)

	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)

	fmt.Println("Configmap:")
	listConfigMap, _ := informerFactory.Core().V1().ConfigMaps().Lister().List(labels.Everything())
	for _, obj := range listConfigMap {
		fmt.Printf("%s/%s \n", obj.Namespace, obj.Name)
	}

	fmt.Println("Pod:")
	listPod, _ := informerFactory.Core().V1().Pods().Lister().List(labels.Everything())
	for _, obj := range listPod {
		fmt.Printf("%s/%s \n", obj.Namespace, obj.Name)
	}

	select {}
}