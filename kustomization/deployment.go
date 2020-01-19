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
            - mountPath: {{.ConfigPath}}
              name: {{.Name}}-config
      volumes:
        - name: {{.Name}}-config
          configMap:
            name: {{.Name}}-config`

	tmpl, err := template.New("deployment-config").Parse(deploymentConfigPatch)
	if err != nil {
		panic("deploymentConfigPatch spec template is misconfigured")
	}
	return tmpl
}

func DeploymentSecretPatch() *template.Template {
	deploymentSecretPatch :=
		`apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
spec:
  template:
    spec:
      containers:
        - name: {{.ContainerName}}
          envFrom:
            - secretRef:
                name: {{.Name}}-secret`

	tmpl, err := template.New("deployment-secret").Parse(deploymentSecretPatch)
	if err != nil {
		panic("deploymentSecretPatch spec template is misconfigured")
	}
	return tmpl
}

func DeploymentImagePatch() *template.Template {
	deploymentVersionPatch :=
		`apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Name}}
spec:
  template:
    spec:
      containers:
        - name: {{.ContainerName}}
          image: {{.Image}}`
	tmpl, err := template.New("deployment-secret").Parse(deploymentVersionPatch)
	if err != nil {
		panic("deploymentVersionPatch spec template is misconfigured")
	}
	return tmpl
}
