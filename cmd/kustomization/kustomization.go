package kustomization

import (
	"text/template"
)

func Kustomization() *template.Template {
	kustomization :=
		`apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

resources:
{{range $key, $value := . }}- {{$key}}
{{end}}
`

	tmpl, err := template.New("kustomization").Parse(kustomization)
	if err != nil {
		panic("kustomization spec template is misconfigured")
	}
	return tmpl
}
