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
          image: {{.ImageName}}
          ports:
            - containerPort: {{.ContainerPort}}
`

	tmpl, err := template.New("deployment").Parse(deployment)
	if err != nil {
		panic("deployment spec template is misconfigured")
	}
	return tmpl
}
