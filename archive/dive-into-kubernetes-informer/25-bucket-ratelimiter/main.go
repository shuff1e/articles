package main

import (
	"fmt"
	"strconv"

	"golang.org/x/time/rate"
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
	queue := workqueue.NewRateLimitingQueue(&workqueue.BucketRateLimiter{
		Limiter: rate.NewLimiter(1, 1),
	})

	go func() {
		for {
			item, _ := queue.Get()
			fmt.Println(item.(reconcile.Request).NamespacedName)
			queue.Done(item)
		}
	}()

	for i := 0; i < 100; i ++ {
		queue.AddRateLimited(newItem("abc" + strconv.Itoa(i), "default"))
	}

	select {}
}
