package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/tools/cache"
)

const (
	ObjectNameField = "metadata.name"
)

func main() {

	updates := make(chan []*v1.Pod)

	cliset := mustClientset()

	lw := cache.NewListWatchFromClient(cliset.CoreV1().RESTClient(), "pods", metav1.NamespaceAll, fields.OneTermEqualSelector(ObjectNameField, "extended-resource-demo-2"))

	send := func(objs []interface{}) {
		var pods []*v1.Pod
		for _, o := range objs {
			pods = append(pods, o.(*v1.Pod))
		}
		updates <- pods
	}

	r := cache.NewReflector(lw, &v1.Pod{}, cache.NewUndeltaStore(send, cache.MetaNamespaceKeyFunc), 0 )

	go r.Run(wait.NeverStop)

	for pods := range updates {
		fmt.Println(len(pods), pods)
	}
}