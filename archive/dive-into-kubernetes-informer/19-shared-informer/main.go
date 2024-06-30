package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
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

	listWarcher := cache.NewListWatchFromClient(
		cliset.CoreV1().RESTClient(),
		"configmaps",
		"default",
		fields.Everything(),
	)

	sharedInformer := cache.NewSharedInformer(listWarcher, &v1.ConfigMap{}, 0)
	sharedInformer.AddEventHandler(&CmdHandler{})

	sharedInformer.Run(wait.NeverStop)
}
