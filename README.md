# k8s_event

## Name

*k8s_event* - reports CoreDNS status to Kubernetes events.

## Description

*k8s_event* listens for log printings, and reports them as Events to Kubernetes APIServer.

This plugin requires ...
* the [_kubeapi_ plugin](http://github.com/coredns/kubeapi) to create a connection to the Kubernetes API.
* create/patch/update permission to the Events API.

Enabling this plugin is process-wide: enabling k8s_event in at least one server block enables it globally.

## Syntax

```
k8s_event {
    level LEVELS...
    rate [QPS] [Burst] [CacheSize]
}
```

* `levels` selects what level of logs should be reported as Kubernetes events.
  LEVELS is a space-separated list of log levels, supported levels are `debug`, `error`, `fatal`, `info`, and `warning`.
  The level of the log will be reflected on the `Reason` field of event, e.g. it will use CoreDNSWarning as `Reason` field for a warning log.
  If no level is specified, it defaults to `error` and `warning`.
* `rate` is used to control the throttling of events.
  * QPS is the fill rate of the token bucket in queries per second, which is 1/300 by default.
  * Burst is the burst size used by the token bucket rate filtering, which is 25 by default.
  * CacheSize is the lru cache size used for event caching locally, which is 4096 by default.

## External Plugin

To use this plugin, compile CoreDNS with this plugin added to the `plugin.cfg`.
This plugin also requires the _kubeapi_ plugin, which should be added to the end of `plugin.cfg`.


## Usage

By default, this plugin reports events on behalf of its own CoreDNS Pod,
PodName and Namespace are collected through the [Downward API](https://kubernetes.io/docs/tasks/inject-data-application/environment-variable-expose-pod-information/#the-downward-api).

When deploying CoreDNS, you should include the following environment variables.

```yaml
env:
  - name: COREDNS_POD_NAME
    valueFrom:
      fieldRef:
        fieldPath: metadata.name
  - name: COREDNS_NAMESPACE
    valueFrom:
      fieldRef:
        fieldPath: metadata.namespace
```

When these environment variables are missing, this plugin reports events on behalf of the `default` namespace.

## Examples

It listens for log printings of `info`, `error`, and `warning` levels, reporting them via in-cluster Kubernetes API.
The event sending rate is controlled by `QPS 0.15 token/sec`, `Burst 10 tokens`, and `LRUCacheSize 1024 tokens`.
```
    kubeapi
    k8s_event {
      level info error warning
      rate 0.15 10 1024
    }
```
