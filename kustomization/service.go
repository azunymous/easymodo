package kustomization

import "text/template"

func Service() *template.Template {
	service :=
		`apiVersion: v1
kind: Service
metadata:
  name: {{.Name}}
spec:
  selector:
    app: {{.Name}}
  ports:
    - protocol: {{.Protocol}}
      port: {{.ContainerPort}}
      targetPort: {{.ContainerPort}}
`

	tmpl, err := template.New("service").Parse(service)
	if err != nil {
		panic("service spec template is misconfigured")
	}
	return tmpl
}
