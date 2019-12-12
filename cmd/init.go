package cmd

import (
	"github.com/azunymous/easymodo/cmd/fs"
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path/filepath"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init [application name]",
	Short: "Define base kustomize YAML",
	Long: `Generates base kustomization files for the given application name. 
Defaults to creating a base in the platform directory.`,
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
	resourceFiles := fs.NewFileMap()
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
	kustomization.Create(input.NewKustomization(resourceFiles.GetResources(), ""), resourceFiles)

	fs.WriteAll(resourceFiles, input.Directory(), "base")
}

func useDefault(def string, flag string) string {
	if len(flag) == 0 {
		return def
	}
	return flag
}

func createDirectory() {
	_, err := fs.Get().Stat(input.Directory())
	dirExists, err := afero.DirExists(fs.Get(), input.Directory())

	if dirExists {
		log.Warnf("Platform directory %s already exists", input.Directory())

	} else if err != nil {
		log.Fatalf("Could not read local directory %v", err)
	}

	err = fs.Get().MkdirAll(filepath.Join("./", input.Directory(), "base"), 0755)

	log.Infof("Creating directory %s", input.Directory())

	if err != nil {
		log.Fatalf("Cannot create platform directory %s %v", input.Directory(), err)
	}
}

func createBase(app input.Application, files fs.Files, generators []kustomization.Generator) {
	for _, generate := range generators {
		err := generate(app, files)
		if err != nil {
			log.Fatalf("Could not create resource %v", err)
		}
	}
}
