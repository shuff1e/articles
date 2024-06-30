package main

import (
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

func main() {
	cliset := mustClientset()
	store := cache.NewStore(cache.MetaNamespaceKeyFunc)

	lwc := cache.NewListWatchFromClient(cliset.CoreV1().RESTClient(),
		"pods",
		"default",
		fields.Everything())

	df := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction: cache.MetaNamespaceKeyFunc,
		KnownObjects: store,
	})

	rf := cache.NewReflector(lwc, &v1.Pod{}, df, time.Second * 0)
	rsCH := make(chan struct{})

	go func() {
		rf.Run(rsCH)
	}()

	for {
		df.Pop(func(i interface{}) error {
			for _, d := range i.(cache.Deltas) {
				fmt.Println(d.Type, ":", d.Object.(*v1.Pod).Name,
					"-", d.Object.(*v1.Pod).Status.Phase)
				switch d.Type {
				case cache.Sync, cache.Added:
					store.Add(d.Object)
				case cache.Updated:
					store.Update(d.Object)
				case cache.Deleted:
					store.Delete(d.Object)
				}
			}
			return nil
		})
	}
}