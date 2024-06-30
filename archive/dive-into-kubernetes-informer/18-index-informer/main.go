package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
)

type CmdHandler struct {
}

func (this *CmdHandler) OnAdd(obj interface{}) {
	fmt.Println("Add: ", obj.(*v1.Pod).Name)
}

func (this *CmdHandler) OnUpdate(obj interface{}, newObj interface{}) {
	fmt.Println("Update: ", newObj.(*v1.Pod).Name)
}

func (this *CmdHandler) OnDelete(obj interface{}) {
	fmt.Println("Delete: ", obj.(*v1.Pod).Name)
}

func LabelsIndexFunc(obj interface{}) ([]string, error) {
	metaD, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	return []string{metaD.GetLabels()["run"]}, nil
}

func main() {
	cliset := mustClientset()

	listWatcher := cache.NewListWatchFromClient(
		cliset.CoreV1().RESTClient(),
		"pods",
		"default",
		fields.Everything())

	myIndexer := cache.Indexers{
		"run": LabelsIndexFunc,
	}

	i, c := cache.NewIndexerInformer(listWatcher, &v1.Pod{}, 0, &CmdHandler{}, myIndexer)

	go c.Run(wait.NeverStop)

	time.Sleep(time.Second*3)

	objList, err := i.ByIndex("run", "nginx")
	if err != nil {
		panic(err)
	}

	fmt.Println(objList[0].(*v1.Pod).Name)
}
