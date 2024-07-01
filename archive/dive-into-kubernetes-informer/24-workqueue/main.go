package main

import (
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func newItem(namespace, name string) reconcile.Request {
	return reconcile.Request{
		types.NamespacedName{
			Namespace: namespace,
			Name: name,
		},
	}
}

func main() {
	q := workqueue.New()
	go func() {
		for {
			item, shutdown := q.Get()
			if shutdown {
				return
			}
			fmt.Println(item.(reconcile.Request).NamespacedName)
			time.Sleep(time.Millisecond*10)
			q.Done(item)
		}
	}()

	for {
		q.Add(newItem("abc", "default"))
		q.Add(newItem("abc", "default"))
		q.Add(newItem("abc2", "default"))
		time.Sleep(time.Second*1)
		fmt.Println("Insert new item")
	}
}
