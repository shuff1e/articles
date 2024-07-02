package main

import (
	"context"
	"fmt"
	"os"
	"time"

	apicorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/uuid"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/client-go/tools/record"
)

func main() {
	cliset := mustClientset()

	broadcaster := record.NewBroadcaster()
	broadcaster.StartRecordingToSink(&corev1.EventSinkImpl{
		Interface: cliset.CoreV1().Events(""),
	})
	eventRecorder := broadcaster.NewRecorder(scheme.Scheme, apicorev1.EventSource{
		Component: "test-leader-election",
	})

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	// add a uniquifier so that two processes on the same host don't accidentally both become active
	id := hostname + "_" + string(uuid.NewUUID())


	rl, err := resourcelock.New(resourcelock.ConfigMapsResourceLock,
		"default",
		"configmap-leader-election-4",
		cliset.CoreV1(),
		cliset.CoordinationV1(),
		resourcelock.ResourceLockConfig{
			Identity: id,
			EventRecorder: eventRecorder,
		})
	if err != nil {
		panic(err)
	}

	cfg := &leaderelection.LeaderElectionConfig{
		Lock: rl,
		LeaseDuration: time.Second * 30,
		RenewDeadline: time.Second * 10,
		RetryPeriod: time.Second * 6,
		WatchDog: leaderelection.NewLeaderHealthzAdaptor(time.Second*20),
		Name: "test-leader-election",
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) {
				fmt.Println("start")
				eventRecorder.Event(&apicorev1.Pod{},apicorev1.EventTypeNormal, "test", "test")
				select {}
			},
			OnStoppedLeading: func() {
				fmt.Println("stopped")
			},
		},
	}

	leaderElector, err := leaderelection.NewLeaderElector(*cfg)

	leaderElector.Run(context.TODO())
}

