package cmd

import (
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	"github.com/azunymous/easymodo/cmd/resources"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

var appFs = afero.NewOsFs()

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [application name]",
	Short: "Generate kustomize YAML",
	Long: `Generates kustomization files in the ./platform
directory. Defaults to creating a base and an overlay for
deploying locally.`,
	Run:     newInitCommand,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"base", "generate"},
}

var imageUri string
var port int
var protocol string
var ingress string

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&imageUri, "image", "i", "", "Set image e.g nginx:1.7.9")
	initCmd.Flags().IntVarP(&port, "port", "p", 8080, "Set container port")
	initCmd.Flags().StringVar(&protocol, "protocol", "TCP", "Set protocol")
	initCmd.Flags().StringVar(&ingress, "ingress", "", "Enable ingress resource generation with given host")
}

func newInitCommand(cmd *cobra.Command, args []string) {
	resourceFiles := resources.NewFileMap()
	app := input.Application{
		Name:          args[0],
		Stateful:      false,
		Image:         useDefault(args[0]+":latest", imageUri),
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
}

func useDefault(def string, flag string) string {
	if len(flag) == 0 {
		return def
	}
	return flag
}

func createDirectory() {
	_, err := appFs.Stat(directory)
	dirExists, err := afero.DirExists(appFs, directory)

	if dirExists {
		log.Warnf("Platform directory %s already exists", directory)

	} else if err != nil {
		log.Fatalf("Could not read local directory %v", err)
	}

	err = appFs.MkdirAll(filepath.Join("./", directory, "base"), 0755)

	log.Infof("Creating directory %s", directory)

	if err != nil {
		log.Fatalf("Cannot create platform directory %s %v", directory, err)
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
	_ = appFs.Mkdir(path.Join(directory, subDir), 0755)

	_ = files.Write(appFs, directory, subDir)
}
