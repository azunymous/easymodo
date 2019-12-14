package kustomization

import (
	"text/template"
)

func Deployment() *template.Template {
	deployment :=
		`apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
  labels:
    app: {{.Name}}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{.Name}}
  template:
    metadata:
      labels:
        app: {{.Name}}
    spec:
      containers:
        - name: {{.ContainerName}}
          image: {{.Image}}
          ports:
            - containerPort: {{.ContainerPort}}
`

	tmpl, err := template.New("deployment").Parse(deployment)
	if err != nil {
		panic("deployment spec template is misconfigured")
	}
	return tmpl
}

func DeploymentConfigPatch() *template.Template {
	deploymentConfigPatch :=
		`apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
spec:
  template:
    spec:
      containers:
        - name: {{.ContainerName}}
          volumeMounts:
            - mountPath: /config/
              name: {{.Name}}-config
      volumes:
        - name: {{.Name}}-config
          configMap:
            name: {{.Name}}-config`

	tmpl, err := template.New("deployment").Parse(deploymentConfigPatch)
	if err != nil {
		panic("deploymentConfigPatch spec template is misconfigured")
	}
	return tmpl
}
