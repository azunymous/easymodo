// Package input provides the structs for transforming templates into Kubernetes/Kustomize resource files.
package input

import (
	"github.com/ghodss/yaml"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path"
)

// Application defines the struct for a single application and possible configuration that would be
// required for a deployment.
type Application struct {
	Name           string
	Namespace      string
	Stateful       bool
	Image          string
	ContainerName  string
	ContainerPort  int
	Protocol       string
	Host           string
	ConfigPath     string
	Replicas       int
	CpuRequests    string
	MemoryRequests string
	CpuLimits      string
	MemoryLimits   string
}

// GetBaseApp reads the base deployment file and returns the set application name and port
func GetBaseApp(fs afero.Fs, dir string) (string, string, int) {
	df, err := afero.ReadFile(fs, path.Join(dir, "base", "deployment.yaml"))
	if err != nil {
		log.Fatalf("Could not open base deployment.yaml file: %v. Make sure you have a base deployment or call easymodo create base", err)
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

	log.Debugf("Read base deployment file %s. Using %s as application name", deployment.Metadata.Name, deployment.Metadata.Name)
	containers := deployment.Spec.Template.Spec.Containers
	var port int
	var image string
	if len(containers) > 0 {
		port = containers[0].Ports[0].ContainerPort
		image = containers[0].Image
	} else {
		log.Warnf("Cannot determine container port from base deployment")
	}
	return deployment.Metadata.Name, image, port
}

func ValidateNamespaceOrSuffix(suffix string, appName string, args []string, c *cobra.Command) (string, string) {
	var (
		namespace string
		nsDir     string
	)

	if suffix != "" {
		namespace = appName + "-" + suffix
		nsDir = suffix
	} else if len(args) < 1 {
		println(c.UsageString())
		log.Fatalf("No namespace or namespace suffix provided!")
	} else {
		namespace = args[0]
		nsDir = args[0]
	}
	return namespace, nsDir
}
