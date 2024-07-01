package main

import (
	"fmt"
	"strconv"
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
	limiter := workqueue.NewItemExponentialFailureRateLimiter(time.Second*5, time.Second*10)
	queue := workqueue.NewRateLimitingQueue(limiter)
	go func() {
		for {
			item, _ := queue.Get()
			fmt.Println(item.(reconcile.Request).NamespacedName)
			queue.Done(item)
		}
	}()

	for i := 0;i < 100;i++ {
		queue.AddRateLimited(newItem("abc"+strconv.Itoa(i),"default"))
	}

	select {}
}