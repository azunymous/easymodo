/*
Package kustomization provides the go templates for generating Kubernetes and Kustomize resource
files. It provides utility functions for generating files and adding them to a given file map.

For example, after creating a file map use a generator to execute a template and add it to the
files map.
*/
package kustomization

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"text/template"
)

type Generator func(input.Application, fs.Files) error

// Generate returns a function which, when called, will execute the template and add it to the files map.
func Generate(resourceName string, template *template.Template) Generator {
	return func(app input.Application, files fs.Files) error {
		content := strings.Builder{}
		err := template.Execute(&content, app)

		if err != nil {
			return errors.Wrapf(err, "Could not create %s.yaml"+resourceName)
		}
		log.Infof("# Generated %s\n%s", resourceName, content.String())
		files.Add(resourceName+".yaml", content.String())
		return nil
	}
}

func BaseGenerators(ingressEnabled bool) []Generator {
	generators := []Generator{
		Generate("deployment", Deployment()),
		Generate("service", Service()),
	}

	if ingressEnabled {
		generators = append(generators, Generate("ingress", Ingress()))
	}

	return generators
}

func Create(kustomization *input.Kustomization, files fs.Files) {
	content := strings.Builder{}
	err := Kustomization().Execute(&content, kustomization)

	if err != nil {
		log.Fatalf("Could not create kustomization.yaml: %v", err)
	}
	log.Infof("# Generated kustomization\n%s", content.String())
	files.Add("kustomization.yaml", content.String())

}
