# k8t - Kubernetes Template Generation

Simple template generation CLI. NOT intended to be a template management
tool. Rather: a quick way to generate commonly made templates.

Currently supports generating:

* Deployment
* StatefulSet
* Ingress + ClusterIP Service

## Things You Probably Should Change

* Ingress TLS host name
* Issuer in the Ingress. Default is `letsencrypt-staging`
* Number of replicas. Default is 1
