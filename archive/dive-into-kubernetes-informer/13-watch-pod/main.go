package main

import (
	"fmt"
	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func main() {
	cliset := mustClientset()
	lwc := cache.NewListWatchFromClient(cliset.CoreV1().RESTClient(), "pods", "default", fields.Everything())
	watcher, err := lwc.Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	for {
		select {
		case v, ok := <-watcher.ResultChan():
			if ok {
				fmt.Println(v.Type, ":", v.Object.(*v1.Pod).Name, "-", v.Object.(*v1.Pod).Status.Phase)
			}
		}
	}
}
