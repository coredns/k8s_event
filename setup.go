package k8s_event

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

const (
	pluginName = "k8s_event"

	defaultRateCacheSize = 4096
	minRateCacheSize     = 1
	maxRateCacheSize     = 65535
	defaultRateBurst     = 25
	minRateBurst         = 1
	maxRateBurst         = 256
	defaultRateQPS       = 1. / 300.
	minRateQPS           = 1. / 3600.
	maxRateQPS           = 1

	defaultLevel = 1<<Error | 1<<Warning
)

var log = clog.NewWithPlugin(pluginName)

func init() { plugin.Register(pluginName, setup) }

func setup(c *caddy.Controller) error {
	k, err := parse(c)
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	err = k.Init()
	if err != nil {
		return plugin.Error(pluginName, err)
	}

	c.OnStartup(k.Startup(dnsserver.GetConfig(c)))
	c.OnShutdown(k.Shutdown())

	return nil
}

func parse(c *caddy.Controller) (*k8sEvent, error) {
	var (
		ke  *k8sEvent
		err error
	)
	i := 0
	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i++
		ke, err = parseStanza(c)
		if err != nil {
			return ke, err
		}
	}
	return ke, nil
}

// parseStanza parses a k8sEvent stanza
func parseStanza(c *caddy.Controller) (*k8sEvent, error) {
	ke := &k8sEvent{
		levels:    defaultLevel,
		qps:       defaultRateQPS,
		burst:     defaultRateBurst,
		cacheSize: defaultRateCacheSize,
	}
	for c.NextBlock() {
		switch c.Val() {
		case "level":
			levelsArgs := c.RemainingArgs()
			if len(levelsArgs) == 0 {
				return nil, c.ArgErr()
			}
			levels := 0
			for _, l := range levelsArgs {
				level, err := levelFromString(strings.ToLower(l))
				if err != nil {
					return nil, err
				}
				levels |= 1 << level
			}
			ke.levels = levels
		case "rate":
			args := c.RemainingArgs()
			argsLen := len(args)
			if argsLen == 0 || argsLen > 3 {
				return nil, c.ArgErr()
			}

			qps, err := strconv.ParseFloat(args[0], 32)
			if err != nil {
				return nil, err
			}
			if qps < minRateQPS || qps > maxRateQPS {
				return nil, c.Errf("qps must be in range [%f, %f]: %d", minRateQPS, maxRateQPS, qps)
			}
			ke.qps = float32(qps)

			if argsLen--; argsLen == 0 {
				break
			}

			burst, err := strconv.Atoi(args[1])
			if err != nil {
				return nil, err
			}
			if burst < minRateBurst || burst > maxRateBurst {
				return nil, c.Errf("burst must be in range [%d, %d]: %d", minRateBurst, maxRateBurst, burst)
			}
			ke.burst = burst

			if argsLen--; argsLen == 0 {
				break
			}

			cacheSize, err := strconv.Atoi(args[2])
			if err != nil {
				return nil, err
			}
			if cacheSize < minRateCacheSize || cacheSize > maxRateCacheSize {
				return nil, c.Errf("cacheSize must be in range [%d, %d]: %d", minRateCacheSize, maxRateCacheSize, cacheSize)
			}
			ke.cacheSize = cacheSize
		default:
			return nil, c.Errf("unknown property '%s'", c.Val())
		}
	}
	return ke, nil
}

type Level int

const (
	Debug Level = iota
	Error
	Fatal
	Info
	Warning
)

func levelFromString(s string) (Level, error) {
	switch s {
	case "debug":
		return Debug, nil
	case "error":
		return Error, nil
	case "fatal":
		return Fatal, nil
	case "info":
		return Info, nil
	case "warning":
		return Warning, nil
	}
	return Debug, fmt.Errorf("invalid Level: %s", s)
}
