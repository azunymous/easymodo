package cmd

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/azunymous/easymodo/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path/filepath"
)

// initCmd represents the init command for initializing a kustomize base.
var baseCmd = &cobra.Command{
	Use:   "base [application name]",
	Short: "Define base kustomize YAML",
	Long: `Generates base kustomization files for the given application name. 
Defaults to creating a base in the platform directory.`,
	Run:     newBaseCommand,
	Args:    cobra.MinimumNArgs(1),
	Aliases: []string{"init"},
}

var imageUri string
var port int
var protocol string

func init() {
	createCmd.AddCommand(baseCmd)

	baseCmd.Flags().StringVarP(&imageUri, "image", "i", "", "Set image e.g nginx:1.7.9")
	baseCmd.Flags().IntVarP(&port, "port", "p", 8080, "Set container port")
	baseCmd.Flags().StringVar(&protocol, "protocol", "TCP", "Set protocol")
	baseCmd.Flags().StringVar(IngressFlag(), "ingress", "", "Enable ingress resource generation with given host")
}

func newBaseCommand(_ *cobra.Command, args []string) {
	resourceFiles := fs.NewFileMap()
	app := input.Application{
		Name:          args[0],
		Stateful:      false,
		Image:         useDefault(args[0]+":latest", imageUri),
		ContainerName: args[0],
		ContainerPort: port,
		Protocol:      "TCP",
		Host:          Ingress(),
	}

	log.Infof("initializing current directory for application %s", app.Name)

	createDirectory()

	createBase(app, resourceFiles, kustomization.BaseGenerators(Ingress() != ""))
	kustomization.Create(input.NewKustomization(resourceFiles.GetFilenames(), ""), resourceFiles)

	resourceFiles.WriteAll(Directory(), "base")
}

func useDefault(def string, flag string) string {
	if len(flag) == 0 {
		return def
	}
	return flag
}

func createDirectory() {
	_, err := fs.Get().Stat(Directory())
	dirExists, err := afero.DirExists(fs.Get(), Directory())

	if dirExists {
		log.Warnf("Platform directory %s already exists", Directory())

	} else if err != nil {
		log.Fatalf("Could not read local directory %v", err)
	}

	err = fs.Get().MkdirAll(filepath.Join("./", Directory(), "base"), 0755)

	log.Infof("Creating directory %s", Directory())

	if err != nil {
		log.Fatalf("Cannot create platform directory %s %v", Directory(), err)
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
