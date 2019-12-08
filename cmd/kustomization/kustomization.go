package kustomization

import (
	"text/template"
)

func Kustomization() *template.Template {
	kustomization :=
		`apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
{{if $namespace := .Namespace}}namespace: {{$namespace}}{{end}}
resources:
{{range $key, $value := .Res }}- {{$value}}
{{end}}
`

	tmpl, err := template.New("kustomization").Parse(kustomization)
	if err != nil {
		panic("kustomization spec template is misconfigured")
	}
	return tmpl
}
