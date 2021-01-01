package main

import (
	"flag"
	"fmt"
	"strings"
	"text/template"
)

type Template struct {
	Name               string
	Namespace          string
	ContainerImage     string
	ContainerPort      int
	ContainerMountPath string
	IngressURI         string
}

func main() {
	isDeployment := flag.Bool("deployment", true, "kind: Deployment")
	isStatefulSet := flag.Bool("stateful", false, "kind: StatefulSet")
	genIngress := flag.Bool("ingress", false, "Generate an Ingress/Service")
	name := flag.String("name", "app", "Name of the app (no spaces) to use as the name label in metadata")
	containerImage := flag.String("image", "", "Container image")
	containerPort := flag.Int("port", 80, "Container port to bind to in Deployment")
	containerMountPath := flag.String("mount", "/data", "Mount path in the container for StatefulSet PersistentVolumes")
	ingressURI := flag.String("uri", "/", "URI for the Ingress")
	namespace := flag.String("ns", "default", "namespace for generated resources")

	flag.Parse()

	data := Template{
		Name:               *name,
		Namespace:          *namespace,
		ContainerImage:     *containerImage,
		ContainerPort:      *containerPort,
		ContainerMountPath: *containerMountPath,
		IngressURI:         *ingressURI,
	}

	var buf strings.Builder

	if *genIngress {
		ingress().Execute(&buf, data)
	}

	if *isStatefulSet {
		stateful().Execute(&buf, data)
	} else if *isDeployment {
		deployment().Execute(&buf, data)
	}

	fmt.Println(
		strings.TrimSpace(buf.String()))
}

// Ingress produces an Ingress YAML
func ingress() *template.Template {
	return template.Must(
		template.New("ingress").Parse(`
---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: {{.Name}}-ingress
  namespace: {{.Namespace}}
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/issuer: "letsencrypt-staging"
spec:
  tls:
  - hosts:
    - k8s.local
    secretName: {{.Name}}-cert
  rules:
  - http:
      paths:
      - path: {{.IngressURI}}
        pathType: Prefix
        backend:
          service:
            name: {{.Name}}
            port:
              number: 80
---

apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  selector:
    name: {{.Name}}
  type: ClusterIP
  ports:
  - protocol: TCP
    port: 80
    targetPort: {{.ContainerPort}}`),
	)
}

// deployment produces a Deployment YAML
func deployment() *template.Template {
	return template.Must(
		template.New("deployment").Parse(`
---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  namespace: {{.Namespace}}
spec:
  selector:
    matchLabels:
      name: {{.Name}}
  replicas: 1
  template:
    metadata:
      labels:
        name: {{.Name}}
    spec:
      containers:
      - name: app
        image: {{.ContainerImage}}
        ports:
        - containerPort: {{.ContainerPort}}`),
	)
}

func stateful() *template.Template {
	return template.Must(
		template.New("statefulSet").Parse(`
---

apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: {{.Name}}
spec:
  selector:
    matchLabels:
	  name: {{.Name}}
  serviceName: {{.Name}}
  replicas: 1
  template:
    metadata:
      labels:
		name: {{.Name}}
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: {{.Name}}
        image: {{.ContainerImage}}
        ports:
        - containerPort: {{.ContainerPort}}
        volumeMounts:
        - name: {{.Name}}-volume
          mountPath: {{.ContainerMountPath}}
  volumeClaimTemplates:
  - metadata:
      name: {{.Name}}-volume
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: "my-storage-class"
      resources:
        requests:
          storage: 1Gi
`),
	)
}
