package main

import (
	"fmt"

	"k8s.io/client-go/tools/cache"
)

type Pod struct {
	Name string
	Value int
}

func NewPod(name string, v int) Pod {
	return Pod{
		Name: name,
		Value: v,
	}
}

func PodKeyFunc(obj interface{}) (string, error) {
	return obj.(Pod).Name, nil
}

func main() {
	df := cache.NewDeltaFIFOWithOptions(cache.DeltaFIFOOptions{
		KeyFunction: PodKeyFunc,
	})

	pod1 := NewPod("pod-1", 1)
	pod2 := NewPod("pod-2", 2)
	pod3 := NewPod("pod-3", 3)
	df.Add(pod1)
	df.Add(pod2)
	df.Add(pod3)

	pod1.Value = 11
	df.Update(pod1)
	df.Update(pod1)

	fmt.Println(df.List())
	
	for {
		df.Pop(func(i interface{}) error {
			for _, delta := range i.(cache.Deltas) {
				switch delta.Type {
				case cache.Added:
					fmt.Printf("Add Event: %v \n", delta.Object)
					break
				case cache.Updated:
					fmt.Printf("Update Event: %v \n", delta.Object)
					break
				case cache.Deleted:
					fmt.Printf("Delete Event: %v \n", delta.Object)
					break
				case cache.Sync:
					fmt.Printf("Sync Event: %v \n", delta.Object)
					break
				case cache.Replaced:
					fmt.Printf("Replaced Event: %v \n", delta.Object)
					break
				}
			}
			return nil
		})
	}
}