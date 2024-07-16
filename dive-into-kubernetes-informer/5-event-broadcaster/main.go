package main

import (
	"time"

	apicorev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

func main() {
	cliset := mustClientset()

	broadcaster := record.NewBroadcaster()

	broadcaster.StartLogging(klog.Infof)
	broadcaster.StartRecordingToSink(&corev1.EventSinkImpl{
		Interface: cliset.CoreV1().Events(""),
	})

	eventRecorder := broadcaster.NewRecorder(scheme.Scheme, apicorev1.EventSource{
		Component: "test-leader-election",
	})

	pod := &apicorev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "abc",
		},
	}
	for {
		eventRecorder.AnnotatedEventf(pod, map[string]string{"a": "b"}, apicorev1.EventTypeWarning, "Evicted", "test %s %s", "a", "b")
		eventRecorder.Eventf(pod, apicorev1.EventTypeWarning, "Evicted", "test %s %s", "a", "b")
		eventRecorder.Event(pod, apicorev1.EventTypeWarning, "Evicted", "test")
		time.Sleep(time.Second*10)
	}
}