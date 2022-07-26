package k8s_event

import (
	"os"
	"testing"

	corev1 "k8s.io/api/core/v1"
)

func TestK8sEventInit(t *testing.T) {
	tests := []struct {
		namespace string
		pod       string
		ref       corev1.ObjectReference
	}{
		{
			"",
			"",
			corev1.ObjectReference{
				Kind: "Namespace",
				Name: "default",
			},
		},
		{
			"ns1",
			"",
			corev1.ObjectReference{
				Kind: "Namespace",
				Name: "default",
			},
		},
		{
			"",
			"pod1",
			corev1.ObjectReference{
				Kind: "Namespace",
				Name: "default",
			},
		},
		{
			"ns1",
			"pod1",
			corev1.ObjectReference{
				Kind:      "Pod",
				Name:      "pod1",
				Namespace: "ns1",
			},
		},
	}
	for i, test := range tests {
		err := os.Setenv("COREDNS_NAMESPACE", test.namespace)
		if err != nil {
			t.Errorf("Test %d: k8sEvent set env failed, err: %s", i, err)
		}
		err = os.Setenv("COREDNS_POD_NAME", test.pod)
		if err != nil {
			t.Errorf("Test %d: k8sEvent set env failed, err: %s", i, err)
		}
		ke := &k8sEvent{}
		err = ke.Init()
		if err != nil {
			t.Errorf("Test %d: k8sEvent init failed, err: %s", i, err)
		}
		if ke.ref.Kind != test.ref.Kind || ke.ref.Namespace != test.ref.Namespace || ke.ref.Name != test.ref.Name {
			t.Errorf("Test %d: Expected object reference is %s, instead found %s", i, test.ref.String(), ke.ref.String())
		}
		err = os.Unsetenv("COREDNS_NAMESPACE")
		if err != nil {
			t.Errorf("Test %d: k8sEvent unset env failed, err: %s", i, err)
		}
		err = os.Unsetenv("COREDNS_POD_NAME")
		if err != nil {
			t.Errorf("Test %d: k8sEvent unset env failed, err: %s", i, err)
		}
	}
}
