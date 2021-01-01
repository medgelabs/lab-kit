# Lab Kit

A collection of infra / ops tooling to jump start new projects
with Developer-oriented tooling.

Central to all this is a Container Orchestrator. Here, we use
Kubernetes.

## K8T

Simple K8S template generation

## CI

## Deployments


## Logging

## Monitoring / Metrics

Prometheus scrapes metrics for both the K8S infrastructure and app Services.
These services are deployed to the `monitoring` namespace.

For an application to be scraped, simply add the following annotations to your
Deployment YAML:

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: APP_NAME
spec:
  ...
  template:
    metadata:
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/path: /test
      labels:
        name: APP_NAME
```

`prometheus.io/scrape` tells Prometheus to scrape metrics. `prometheus.io/path` allows you
to override the default metrics URI. It is set to `/metrics` by default.

## Chaos Testing
