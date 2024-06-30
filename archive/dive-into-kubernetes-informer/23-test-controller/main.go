package main

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/ratelimiter"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type TestController struct {
}

func (this TestController) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	fmt.Println(req.NamespacedName)
	return reconcile.Result{}, nil
}

type MyWeb struct {
	h handler.EventHandler
}

func NewMyWeb(h handler.EventHandler) *MyWeb {
	return &MyWeb{
		h: h,
	}
}

func (m *MyWeb) Start(ctx context.Context) error {
	r := gin.New()
	r.GET("/add", func(c *gin.Context) {
		p := &corev1.Pod{}
		p.Name = "hello-world"
		p.Namespace = "testNamespace"
		m.h.Create(c, event.CreateEvent{Object: p}, getQueue())
	})
	return r.Run(":8081")
}

func main() {
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		panic(err)
	}

	c, err := controller.New("App", mgr, controller.Options{
		// 真正执行的reconsiler pkg.internal.controller.Controller.processNextWorkItem 从这往下看
		Reconciler: &TestController{},
		// 开启多个 Reconcile, 如果是多个同样的资源仅会触发一次，直到该资源done
		MaxConcurrentReconciles: 1,
		NewQueue: func(controllerName string, rateLimiter ratelimiter.RateLimiter) workqueue.RateLimitingInterface {
			ret :=  workqueue.NewRateLimitingQueueWithConfig(rateLimiter, workqueue.RateLimitingQueueConfig{
				Name: controllerName,
			})
			setQueue(ret)
			return ret
		},
	})
	if err != nil {
		panic(err)
	}

	h := &handler.TypedEnqueueRequestForObject[client.Object]{}
	err = c.Watch(source.Kind[client.Object](mgr.GetCache(), &corev1.Pod{}, h))
	if err != nil {
		panic(err)
	}

	mgr.Add(NewMyWeb(h))

	mgr.Start(context.TODO())
}

var queue workqueue.RateLimitingInterface

func setQueue(q workqueue.RateLimitingInterface) {
	queue = q
}

func getQueue() workqueue.RateLimitingInterface {
	return queue
}