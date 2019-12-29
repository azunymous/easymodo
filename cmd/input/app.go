// Package input provides the structs for transforming templates into Kubernetes/Kustomize resource files.
package input

import (
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
)

type Application struct {
	Name          string
	Namespace     string
	Stateful      bool
	Image         string
	ContainerName string
	ContainerPort int
	Protocol      string
	Host          string
	ConfigPath    string
}

func GetAppName(fs afero.Fs, dir string) (string, int) {
	df, err := afero.ReadFile(fs, path.Join(dir, "base", "deployment.yaml"))
	if err != nil {
		log.Fatalf("Could not open base deployment.yaml file: %v. Make sure you have a base deployment or call easymodo init", err)
	}

	if err != nil {
		log.Fatalf("Could not read base deployment.yaml file: %v", err)
	}

	type Deployment struct {
		Kind     string `yaml:"kind"`
		Metadata struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
		Spec struct {
			Template struct {
				Spec struct {
					Containers []struct {
						Name  string `yaml:"name"`
						Image string `yaml:"image"`
						Ports []struct {
							ContainerPort int `yaml:"containerPort"`
						} `yaml:"ports"`
					} `yaml:"containers"`
				} `yaml:"spec"`
			} `yaml:"template"`
		} `yaml:"spec"`
	}

	deployment := Deployment{}
	err = yaml.Unmarshal(df, &deployment)

	if err != nil {
		log.Fatalf("Error unmarshalling deployment.yaml: %v", err)
	}

	if deployment.Kind != "Deployment" {
		log.Fatalf("Kubernetes resource file is not a Deployment")
	}

	log.Infof("Read base deployment file %s. Using %s as application name", deployment.Metadata.Name, deployment.Metadata.Name)
	containers := deployment.Spec.Template.Spec.Containers
	var port int
	if len(containers) > 0 {
		port = containers[0].Ports[0].ContainerPort
	} else {
		log.Warnf("Cannot determine container port from base deployment")
	}
	return deployment.Metadata.Name, port
}
