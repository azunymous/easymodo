package cmd

import (
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	"github.com/azunymous/easymodo/cmd/resources"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [application name]",
	Short: "Generate kustomize YAML",
	Long: `Generates kustomization files in the ./platform
directory. Defaults to creating a base and an overlay for
deploying locally.`,
	Run:     initCommand,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"create", "generate"},
}

var directory string
var imageName string
var port int
var protocol string
var ingress string

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&directory, "directory", "d", "platform", "directory for kustomization files and folders")
	initCmd.Flags().StringVarP(&imageName, "imageName", "i", "", "Set image name")
	initCmd.Flags().IntVarP(&port, "port", "p", 8080, "Set container port")
	initCmd.Flags().StringVar(&protocol, "protocol", "TCP", "Set protocol")
	initCmd.Flags().StringVar(&ingress, "ingress", "", "Enable ingress resource generation with given host")
}

func initCommand(cmd *cobra.Command, args []string) {
	resourceFiles := resources.NewFileMap()

	app := input.Application{
		Name:          args[0],
		Stateful:      false,
		ImageName:     useDefault(args[0], imageName) + ":latest",
		ContainerName: args[0],
		ContainerPort: port,
		Protocol:      "TCP",
		Host:          ingress,
	}

	log.Infof("initializing current directory for application %s", app.Name)

	createDirectory()

	createBase(app, resourceFiles, kustomization.Generators(cmd.Flags().Changed("ingress")))
	createKustomization(input.NewKustomization(resourceFiles.GetResources(), ""), resourceFiles)
	writeFiles(resourceFiles, "base")
	relativeBasePath := filepath.Join("../", "base", "kustomization.yaml")

	resourceFiles = resources.NewFileMap()
	res := []string{relativeBasePath}
	createKustomization(input.NewKustomization(res, app.Name+"-dev"), resourceFiles)
	writeFiles(resourceFiles, "dev")

}

func useDefault(def string, flag string) string {
	if len(flag) == 0 {
		return def
	}
	return flag
}

func createDirectory() {
	stat, err := os.Stat(directory)

	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Could not read local directory %v", err)
	}

	if err == nil && stat.IsDir() {
		log.Warnf("Platform directory %s already exists", directory)
	}

	err = os.MkdirAll(filepath.Join("./", directory, "base"), 0755)

	log.Infof("Creating directory %s", directory)

	if err != nil {
		log.Fatalf("Cannot create platform directory %s %v", err)
	}
}

func createBase(app input.Application, files resources.Files, generators []kustomization.Generator) {
	for _, generate := range generators {
		err := generate(app, files)
		if err != nil {
			log.Fatalf("Could not create resource %v", err)
		}
	}
}

func createKustomization(resources *input.Kustomization, files *resources.FileMap) {
	err := kustomization.GenerateKustomization(resources, files)
	if err != nil {
		log.Fatalf("Could not create kustomization.yaml: %v", err)
	}
}

func writeFiles(files resources.Files, subDir string) {
	wd, _ := os.Getwd()
	_ = os.Chdir(directory)
	_ = os.Mkdir(subDir, 0755)
	_ = os.Chdir(subDir)
	_ = files.Write()

	_ = os.Chdir(wd)
}
