package input

import (
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
)

type Application struct {
	Name          string
	Stateful      bool
	Image         string
	ContainerName string
	ContainerPort int
	Protocol      string
	Host          string
	Namespace     string
}

func GetAppName(fs afero.Fs, dir string) string {
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
	}

	deployment := Deployment{}
	err = yaml.Unmarshal(df, &deployment)

	if err != nil {
		log.Fatalf("Error unmarshalling deployment.yaml: %v", err)
	}

	if deployment.Kind != "Deployment" {
		log.Fatalf("Kubernetes resource file is not a Deployment")
	}

	log.Infof("Read base deployment file %s. Using %s as namespace prefix", deployment.Metadata.Name, deployment.Metadata.Name)
	return deployment.Metadata.Name
}
