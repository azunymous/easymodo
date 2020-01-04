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
{{if .Config -}}
configMapGenerator:
{{- range $key, $value := .Config}}
  - name: {{$key}}
    files:{{range $index, $filename := $value }}
      - {{$filename}}
{{- end}}
{{- end}}
{{- end}}
{{if .Secrets -}}
secretGenerator:
{{- range $key, $value := .Secrets}}
  - name: {{$key}}
    envs:{{range $index, $filename := $value }}
      - {{$filename}}
{{- end}}
{{- end}}
{{- end}}

{{- if .Patches}}
patchesStrategicMerge:
{{range $key, $value := .Patches }}  - {{$value}}
{{end}}
{{- end}}
`

	tmpl, err := template.New("kustomization").Parse(kustomization)
	if err != nil {
		panic("kustomization spec template is misconfigured")
	}
	return tmpl
}
