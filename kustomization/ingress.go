package kustomization

import "text/template"

func Ingress() *template.Template {
	ingress :=
		`apiVersion: networking.k8s.io/v1beta1
kind: Ingress
metadata:
  name: {{.Name}}
spec:
  rules:
  - host: {{.Host}}
    http:
      paths:
      - backend:
          serviceName: {{.Name}}
          servicePort: {{.ContainerPort}}
`

	tmpl, err := template.New("service").Parse(ingress)
	if err != nil {
		panic("ingress spec template is misconfigured")
	}
	return tmpl
}
