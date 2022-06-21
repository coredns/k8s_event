package k8s_event

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

type Interceptor struct {
	recorder record.EventRecorder
	ref      *corev1.ObjectReference
	levels   int
}

func NewInterceptor(ref *corev1.ObjectReference, recorder record.EventRecorder, levels int) *Interceptor {
	return &Interceptor{
		recorder: recorder,
		ref:      ref,
		levels:   levels,
	}
}

func (i *Interceptor) Debug(plugin string, v ...interface{}) {
	if i.levels&(1<<Debug) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeNormal, "CoreDNSDebug", plugin+fmt.Sprint(v...))
	}
}

func (i *Interceptor) Debugf(plugin string, format string, v ...interface{}) {
	if i.levels&(1<<Debug) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeNormal, "CoreDNSDebug", plugin+fmt.Sprintf(format, v...))
	}
}

func (i *Interceptor) Info(plugin string, v ...interface{}) {
	if i.levels&(1<<Info) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeNormal, "CoreDNSInfo", plugin+fmt.Sprint(v...))
	}
}

func (i *Interceptor) Infof(plugin string, format string, v ...interface{}) {
	if i.levels&(1<<Info) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeNormal, "CoreDNSInfo", plugin+fmt.Sprintf(format, v...))
	}
}

func (i *Interceptor) Warning(plugin string, v ...interface{}) {
	if i.levels&(1<<Warning) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeWarning, "CoreDNSWarning", plugin+fmt.Sprint(v...))
	}
}

func (i *Interceptor) Warningf(plugin string, format string, v ...interface{}) {
	if i.levels&(1<<Warning) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeWarning, "CoreDNSWarning", plugin+fmt.Sprintf(format, v...))
	}
}

func (i *Interceptor) Error(plugin string, v ...interface{}) {
	if i.levels&(1<<Error) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeWarning, "CoreDNSError", plugin+fmt.Sprint(v...))
	}
}

func (i *Interceptor) Errorf(plugin string, format string, v ...interface{}) {
	if i.levels&(1<<Error) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeWarning, "CoreDNSError", plugin+fmt.Sprintf(format, v...))
	}
}

func (i *Interceptor) Fatal(plugin string, v ...interface{}) {
	if i.levels&(1<<Fatal) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeWarning, "CoreDNSFatal", plugin+fmt.Sprint(v...))
	}
}

func (i *Interceptor) Fatalf(plugin string, format string, v ...interface{}) {
	if i.levels&(1<<Fatal) > 0 {
		i.recorder.Eventf(i.ref, corev1.EventTypeWarning, "CoreDNSFatal", plugin+fmt.Sprintf(format, v...))
	}
}
