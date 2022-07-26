package k8s_event

import (
	"fmt"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/reference"
)

type mockEvent struct {
	object    runtime.Object
	eventType string
	reason    string
	message   string
}

type mockRecorder struct {
	events []mockEvent
}

func (r *mockRecorder) Event(object runtime.Object, eventtype, reason, message string) {
	return
}

func (r *mockRecorder) Eventf(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) {
	r.events = append(r.events, mockEvent{
		object:    object,
		eventType: eventtype,
		reason:    reason,
		message:   fmt.Sprintf(messageFmt, args...),
	})
}

func (r *mockRecorder) AnnotatedEventf(object runtime.Object, annotations map[string]string, eventtype, reason, messageFmt string, args ...interface{}) {
}

func (r *mockRecorder) CountEvent() int {
	return len(r.events)
}

func (r *mockRecorder) ContainsEvent(object runtime.Object, eventtype, reason, messageFmt string, args ...interface{}) bool {
	for _, e := range r.events {
		ref1, err := reference.GetReference(runtime.NewScheme(), object)
		if err != nil {
			continue
		}
		ref2, err := reference.GetReference(runtime.NewScheme(), e.object)
		if err != nil {
			continue
		}
		if ref1.String() != ref2.String() {
			continue
		}
		if eventtype != e.eventType {
			continue
		}
		if reason != e.reason {
			continue
		}
		if fmt.Sprintf(messageFmt, args...) != e.message {
			continue
		}
		return true
	}
	return false
}

func TestListenerEventCount(t *testing.T) {
	tests := []struct {
		ref                 corev1.ObjectReference
		levels              int
		pluginName          string
		fmt                 string
		args                string
		expectedEventsCount int
	}{
		{
			corev1.ObjectReference{Kind: "Pod", Name: "pod1", Namespace: "ns1"},
			1<<Debug | 1<<Info | 1<<Warning | 1<<Error | 1<<Fatal,
			"plugin1",
			"fmt %s",
			"args",
			10,
		},
		{
			corev1.ObjectReference{Kind: "Pod", Name: "pod1", Namespace: "ns1"},
			1<<Debug | 1<<Error,
			"plugin1",
			"fmt %s",
			"args",
			4,
		},
		{
			corev1.ObjectReference{Kind: "Pod", Name: "pod1", Namespace: "ns1"},
			1 << Debug,
			"plugin1",
			"fmt %s",
			"args",
			2,
		},
	}

	for i, test := range tests {
		r := &mockRecorder{}
		l := newListener(&test.ref, r, test.levels)
		l.Debug(test.pluginName, test.args)
		l.Debugf(test.pluginName, test.fmt, test.args)
		l.Info(test.pluginName, test.args)
		l.Infof(test.pluginName, test.fmt, test.args)
		l.Warning(test.pluginName, test.args)
		l.Warningf(test.pluginName, test.fmt, test.args)
		l.Error(test.pluginName, test.args)
		l.Errorf(test.pluginName, test.fmt, test.args)
		l.Fatal(test.pluginName, test.args)
		l.Fatalf(test.pluginName, test.fmt, test.args)

		if r.CountEvent() != test.expectedEventsCount {
			t.Errorf("Test %d: Expected %d events, instead found %d events", i, test.expectedEventsCount, r.CountEvent())
		}
	}
}

func TestListenerEventExists(t *testing.T) {
	tests := []struct {
		ref        corev1.ObjectReference
		levels     int
		pluginName string
		fmt        string
		args       string
	}{
		{
			corev1.ObjectReference{Kind: "Pod", Name: "pod1", Namespace: "ns1"},
			1<<Debug | 1<<Info | 1<<Warning | 1<<Error | 1<<Fatal,
			"plugin1",
			"fmt %s",
			"args",
		},
		{
			corev1.ObjectReference{Kind: "Namespace", Namespace: "ns1"},
			1<<Debug | 1<<Info | 1<<Warning | 1<<Error | 1<<Fatal,
			"plugin1",
			"fmt %s",
			"args",
		},
	}

	for i, test := range tests {
		r := &mockRecorder{}
		l := newListener(&test.ref, r, test.levels)
		l.Debug(test.pluginName, test.args)
		l.Debugf(test.pluginName, test.fmt, test.args)
		l.Info(test.pluginName, test.args)
		l.Infof(test.pluginName, test.fmt, test.args)
		l.Warning(test.pluginName, test.args)
		l.Warningf(test.pluginName, test.fmt, test.args)
		l.Error(test.pluginName, test.args)
		l.Errorf(test.pluginName, test.fmt, test.args)
		l.Fatal(test.pluginName, test.args)
		l.Fatalf(test.pluginName, test.fmt, test.args)

		if !r.ContainsEvent(&test.ref, corev1.EventTypeNormal, "CoreDNSDebug", test.pluginName+test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeNormal, "CoreDNSDebug", test.pluginName+test.fmt, test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeNormal, "CoreDNSInfo", test.pluginName+test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeNormal, "CoreDNSInfo", test.pluginName+test.fmt, test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeWarning, "CoreDNSWarning", test.pluginName+test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeWarning, "CoreDNSWarning", test.pluginName+test.fmt, test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeWarning, "CoreDNSError", test.pluginName+test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeWarning, "CoreDNSError", test.pluginName+test.fmt, test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeWarning, "CoreDNSFatal", test.pluginName+test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
		if !r.ContainsEvent(&test.ref, corev1.EventTypeWarning, "CoreDNSFatal", test.pluginName+test.fmt, test.args) {
			t.Errorf("Test %d: expected event not found", i)
		}
	}
}
