package main

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

func LabelsIndexFunc(obj interface{}) ([]string, error) {
	metaD, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	return []string{metaD.GetLabels()["app"]}, nil
}

func main() {
	idxs := cache.Indexers{
		"app": LabelsIndexFunc,
	}

	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod1",
			Namespace: "ns1",
			Labels: map[string]string{
				"app": "l1",
			},
		},
	}

	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod2",
			Namespace: "ns2",
			Labels: map[string]string{
				"app": "l2",
			},
		},
	}

	myIdx := cache.NewIndexer(cache.MetaNamespaceKeyFunc, idxs)
	myIdx.Add(pod1)
	myIdx.Add(pod2)
	fmt.Println(myIdx.IndexKeys("app", "l1"))
}
