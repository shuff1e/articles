package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/cache"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func main() {
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(time.Second * 2)

		pod := &corev1.Pod{}
		err = mgr.GetClient().Get(context.TODO(), types.NamespacedName{
			Namespace: "default",
			Name: "nginx",
		}, pod)
		if err != nil {
			panic(err)
		}
		fmt.Printf("api version: %v, kind: %v\n", pod.APIVersion, pod.Kind)

		informer, _ := mgr.GetCache().GetInformer(context.TODO(), pod)
		fmt.Println("pods", informer.(cache.SharedIndexInformer).GetIndexer().ListKeys())
	}()

	mgr.Start(context.TODO())

}
