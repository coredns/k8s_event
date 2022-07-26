package k8s_event

import (
	"strings"
	"testing"

	"github.com/coredns/caddy"
)

func TestParseStanza(t *testing.T) {
	tests := []struct {
		input              string
		shouldErr          bool
		expectedErrContent string
		expectedQPS        float32
		expectedBurst      int
		expectedCacheSize  int
		expectedLevels     int
	}{
		{
			`k8s_event`,
			false, "",
			defaultRateQPS, defaultRateBurst, defaultRateCacheSize, defaultLevel,
		},
		{
			`k8s_event {
    level debug error
    rate 0.01 16 2048
}`,
			false, "",
			0.01, 16, 2048, 1<<Debug | 1<<Error,
		},
		{
			`k8s_event {
    level err
}`,
			true, "invalid Level",
			0, 0, 0, 0,
		},
		{
			`k8s_event {
    rate
}`,
			true, "Wrong argument count or unexpected line ending",
			0, 0, 0, 0,
		},
		{
			`k8s_event {
    rate 10
}`,
			true, "qps must be in range",
			0, 0, 0, 0,
		},
		{
			`k8s_event {
    rate 0.01 1024
}`,
			true, "burst must be in range",
			0, 0, 0, 0,
		},
		{
			`k8s_event {
    rate 0.01 256 65536
}`,
			true, "cacheSize must be in range",
			0, 0, 0, 0,
		},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		k8sEvent, err := parse(c)

		if test.shouldErr && err == nil {
			t.Errorf("Test %d: Expected error, but did not find error for input '%s'. Error was: '%v'", i, test.input, err)
		}

		if err != nil {
			if !test.shouldErr {
				t.Errorf("Test %d: Expected no error but found one for input %s. Error was: %v", i, test.input, err)
				continue
			}

			if test.shouldErr && (len(test.expectedErrContent) < 1) {
				t.Fatalf("Test %d: Test marked as expecting an error, but no expectedErrContent provided for input '%s'. Error was: '%v'", i, test.input, err)
			}
			if !strings.Contains(err.Error(), test.expectedErrContent) {
				t.Errorf("Test %d: Expected error to contain: %v, found error: %v, input: %s", i, test.expectedErrContent, err, test.input)
			}
			continue
		}

		if k8sEvent.levels != test.expectedLevels {
			t.Errorf("Test %d: Expected k8sEvent to be initialized with %d level, instead found level: '%v' for input '%s'",
				i, test.expectedLevels, k8sEvent.levels, test.input)
		}

		if k8sEvent.qps != test.expectedQPS {
			t.Errorf("Test %d: Expected k8sEvent to be initialized with %f qps, instead found qps: '%v' for input '%s'",
				i, test.expectedQPS, k8sEvent.qps, test.input)
		}

		if k8sEvent.burst != test.expectedBurst {
			t.Errorf("Test %d: Expected k8sEvent to be initialized with %d burst, instead found burst: '%v' for input '%s'",
				i, test.expectedBurst, k8sEvent.burst, test.input)
		}

		if k8sEvent.cacheSize != test.expectedCacheSize {
			t.Errorf("Test %d: Expected k8sEvent to be initialized with %d cacheSize, instead found cacheSize: '%v' for input '%s'",
				i, test.expectedCacheSize, k8sEvent.cacheSize, test.input)
		}
	}
}
