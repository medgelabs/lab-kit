# Lab Kit

A collection of infra / ops tooling to jump start new projects
with Developer-oriented tooling.

Central to all this is a Docker Orchestrator. Here, we use
Kubernetes.

## CI

## Deployments

Deployments to Kubernetes can be direct or through the Deployment CLI.
The CLI simplifies managing raw Kubernetes templates which are, often,
hard to keep straight.

## Monitoring

Monitoring is partially done by Kubernetes, with its built in tooling,
and partially with Prometheus. The AlertManager is used for this purpose.
They can also be exported to tools like Sentry with little additional
work.

## Logging

Logging is handled by the ELK Stack (Elastisearch, Logstash, Kibana). Logs
are piped to Logstash, indexed by Elastisearch, and visualized by Kibana.

## Metrics

Metrics are collected through Prometheus scrapes, read by Grafana Dashboards.
Each application exposes these metrics in Prometheus format via a REST API

## Chaos Testing

The Lab Kit includes a CLI for chaos testing when running locally in Kubernetes.
It simulates network / container failures, high network load, and etcetera.
