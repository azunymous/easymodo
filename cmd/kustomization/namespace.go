package kustomization

import "text/template"

func Namespace() *template.Template {
	namespace :=
		`apiVersion: v1
kind: Namespace
metadata:
  name: {{.Namespace}}
`

	tmpl, err := template.New("namespace").Parse(namespace)
	if err != nil {
		panic("namespace spec template is misconfigured")
	}
	return tmpl
}
