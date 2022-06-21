package k8s_event

import (
	"os"

	plog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/kubeapi"
	corev1 "k8s.io/api/core/v1"

	"github.com/coredns/coredns/core/dnsserver"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"

	"k8s.io/client-go/tools/record"
)

const (
	componentName = "CoreDNS"
)

type k8sEvent struct {
	client kubernetes.Interface
	ref    *corev1.ObjectReference

	rateSet   bool
	qps       float32
	burst     int
	cacheSize int

	levels int

	broadcaster record.EventBroadcaster
}

func (k *k8sEvent) Init() error {
	ns := os.Getenv("COREDNS_NAMESPACE")
	pod := os.Getenv("COREDNS_POD_NAME")
	if len(ns) > 0 && len(pod) > 0 {
		k.ref = &corev1.ObjectReference{
			Kind:      "Pod",
			Name:      pod,
			Namespace: ns,
		}
	} else {
		k.ref = &corev1.ObjectReference{
			Kind: "Namespace",
			Name: "default",
		}
		log.Warning("COREDNS_NAMESPACE or COREDNS_POD_NAME is not set in environment variables, reporting events to default namespace")
	}
	return nil
}

func (k *k8sEvent) Startup(config *dnsserver.Config) func() error {
	return func() error {
		var err error
		k.client, err = kubeapi.Client(config)
		if err != nil {
			return err
		}

		if k.rateSet {
			k.broadcaster = record.NewBroadcasterWithCorrelatorOptions(record.CorrelatorOptions{
				LRUCacheSize: k.cacheSize,
				QPS:          k.qps,
				BurstSize:    k.burst,
			})
		} else {
			k.broadcaster = record.NewBroadcaster()
		}

		source := corev1.EventSource{Component: componentName}
		recorder := k.broadcaster.NewRecorder(scheme.Scheme, source)

		k.broadcaster.StartRecordingToSink(&typedv1.EventSinkImpl{
			Interface: typedv1.New(k.client.CoreV1().RESTClient()).Events(""),
		})

		i := NewInterceptor(k.ref, recorder, k.levels)
		err = plog.LoadInterceptor(i)
		if err != nil {
			return err
		}
		return nil
	}
}

func (k *k8sEvent) Shutdown() func() error {
	return func() error {
		k.broadcaster.Shutdown()
		return nil
	}
}
