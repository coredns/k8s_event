package k8s_event

import (
	"fmt"
	"time"

	"github.com/coredns/coredns/plugin/pkg/rand"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
)

var rn = rand.New(time.Now().UnixNano())

type listener struct {
	recorder record.EventRecorder
	ref      *corev1.ObjectReference
	levels   int
	id       int
}

func newListener(ref *corev1.ObjectReference, recorder record.EventRecorder, levels int) *listener {
	return &listener{
		recorder: recorder,
		ref:      ref,
		levels:   levels,
		id:       rn.Int(),
	}
}

func (l *listener) Name() string {
	return fmt.Sprintf("%s-%d", pluginName, l.id)
}

func (l *listener) Debug(plugin string, v ...interface{}) {
	if l.levels&(1<<Debug) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeNormal, "CoreDNSDebug", plugin+fmt.Sprint(v...))
	}
}

func (l *listener) Debugf(plugin string, format string, v ...interface{}) {
	if l.levels&(1<<Debug) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeNormal, "CoreDNSDebug", plugin+fmt.Sprintf(format, v...))
	}
}

func (l *listener) Info(plugin string, v ...interface{}) {
	if l.levels&(1<<Info) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeNormal, "CoreDNSInfo", plugin+fmt.Sprint(v...))
	}
}

func (l *listener) Infof(plugin string, format string, v ...interface{}) {
	if l.levels&(1<<Info) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeNormal, "CoreDNSInfo", plugin+fmt.Sprintf(format, v...))
	}
}

func (l *listener) Warning(plugin string, v ...interface{}) {
	if l.levels&(1<<Warning) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeWarning, "CoreDNSWarning", plugin+fmt.Sprint(v...))
	}
}

func (l *listener) Warningf(plugin string, format string, v ...interface{}) {
	if l.levels&(1<<Warning) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeWarning, "CoreDNSWarning", plugin+fmt.Sprintf(format, v...))
	}
}

func (l *listener) Error(plugin string, v ...interface{}) {
	if l.levels&(1<<Error) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeWarning, "CoreDNSError", plugin+fmt.Sprint(v...))
	}
}

func (l *listener) Errorf(plugin string, format string, v ...interface{}) {
	if l.levels&(1<<Error) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeWarning, "CoreDNSError", plugin+fmt.Sprintf(format, v...))
	}
}

func (l *listener) Fatal(plugin string, v ...interface{}) {
	if l.levels&(1<<Fatal) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeWarning, "CoreDNSFatal", plugin+fmt.Sprint(v...))
	}
}

func (l *listener) Fatalf(plugin string, format string, v ...interface{}) {
	if l.levels&(1<<Fatal) > 0 {
		l.recorder.Eventf(l.ref, corev1.EventTypeWarning, "CoreDNSFatal", plugin+fmt.Sprintf(format, v...))
	}
}
