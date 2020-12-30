# Provisioning unRaid VMs

Notes:

* Kubernetes Manager node assumed to be `192.168.1.131`

## Setup

* Create an Ubuntu Server VM with desired specs. Use Bridge Network.
* When installing (over VNC) - import SSH identity from Github

The machine can now be provisioned with the Makefile.

## Make

Usage:

```
NHOST="IP_ADDR" make provision
```

Provisions a new host with a Firewall, SSHD config, and other tools.

## Kubernetes



### Manager Node

```
NHOST="IP_ADDR" make kube-manager
```

Installs a Manager node. Make also downloads `/etc/rancher/k3s/k3s.yaml` to `./kube.yaml`.
You can either `mv ./kube.yaml ~/.kube/config` or merge an existing config manually.

### Init Manager Services

We use MetalLB as a Load Balancer and Ingress-NGINX as a reverse proxy/TLS termination point.
Cert-Manager and LetsEncrypt provide the TLS certificates.

Run `make init-k8s` to install these services.

Note: you _may_ see an error about a secret already being created. This is fine, as the secret
is created on first setup. If you re-run `make init-k8s` again, it tries to recreate the secret.
We ignore this failure in make and proceed anyway.

> NOTE: ingress-nginx baremetal deploys as a `NodePort` by default. Change this to `LoadBalancer` so it gets an IP from MetalLB

### Worker Nodes

```
NHOST="IP_ADDR" make kube-worker
```

Installs a Worker node, including grabbing the NODE_TOKEN from the manager node.

### Test the Cluster

`make test-k8s` will run a simple test to ensure the cluster is configured properly.
It deploys an NGINX container, registers a Service, grabs the ingress IP, and tries to curl
the service through `IP/test-cluster`. If all is well, the NGINX default HTML page should
show on the console. Make will remove these pods once complete.
