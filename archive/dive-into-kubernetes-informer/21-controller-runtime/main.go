package main

import (
	"context"
	"fmt"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
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
		fmt.Printf("%+v\n", mgr.GetClient().Scheme())

		for gvk, _ := range mgr.GetScheme().AllKnownTypes() {
			fmt.Println(gvk)
		}

		pod := &corev1.Pod{}
		err = mgr.GetClient().Get(context.TODO(), types.NamespacedName{
			Namespace: "default",
			Name: "nginx",
		}, pod)
		if err != nil {
			panic(err)
		}
		fmt.Printf("api version: %v, kind: %v\n", pod.APIVersion, pod.Kind)

		gvk, err := apiutil.GVKForObject(pod, scheme.Scheme)
		if err != nil {
			panic(err)
		}
		fmt.Println(gvk)

		pod = &corev1.Pod{}
		ptr, err := conversion.EnforcePtr(pod)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Kind: %+v \n", ptr.Type())

	}()

	mgr.Start(context.TODO())
}
