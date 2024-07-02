package main

import (
	"context"
	"fmt"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func main() {
	cliset := mustClientset()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	id := hostname + "_" + string(uuid.NewUUID())

	rl := resourcelock.EndpointsLock{
		EndpointsMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name: "leader-election-1",
		},
		Client: cliset.CoreV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}

	ctx := context.TODO()
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock: &rl,
		LeaseDuration: time.Second * 30,
		RenewDeadline: time.Second * 10,
		RetryPeriod: time.Second * 6,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				fmt.Println("start")
				select {}
			},
			OnStoppedLeading: func() {
				fmt.Println("stopped")
			},
		},
	})
}
