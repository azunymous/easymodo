package kustomization

import (
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/resources"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"strings"
	"text/template"
)

type Generator func(input.Application, resources.Files) error

func Generate(resourceName string, template *template.Template) Generator {
	return func(app input.Application, files resources.Files) error {
		content := strings.Builder{}
		err := template.Execute(&content, app)

		if err != nil {
			return errors.Wrapf(err, "Could not create %s.yaml"+resourceName)
		}
		log.Info(content.String())
		files.Add(resourceName+".yaml", content.String())
		return nil
	}
}

func GenerateKustomization(files resources.Files) error {
	content := strings.Builder{}
	err := Kustomization().Execute(&content, files.Get())

	if err != nil {
		return errors.Wrap(err, "Could not create kustomization.yaml")
	}
	log.Info(content.String())
	files.Add("kustomization.yaml", content.String())
	return nil

}

func Generators() []Generator {
	return []Generator{
		Generate("deployment", Deployment()),
		Generate("service", Service()),
	}
}
