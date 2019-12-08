package cmd

import (
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	"github.com/azunymous/easymodo/cmd/resources"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
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

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:

	initCmd.Flags().StringVarP(&directory, "directory", "d", "platform", "directory for kustomization files and folders")
	initCmd.Flags().StringVarP(&imageName, "imageName", "i", "", "Set image name")
	initCmd.Flags().IntVarP(&port, "port", "p", 8080, "Set container port")
	initCmd.Flags().StringVar(&protocol, "protocol", "TCP", "Set protocol")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initCommand(_ *cobra.Command, args []string) {
	resourceFiles := resources.NewFileMap()

	appName := args[0]

	app := input.Application{
		Name:          appName,
		Stateful:      false,
		ImageName:     useDefault(appName, imageName) + ":latest",
		ContainerName: appName,
		ContainerPort: port,
		Protocol:      "TCP",
	}

	log.Infof("initializing current directory for application %s", app.Name)

	createDirectory()
	createBase(app, resourceFiles, kustomization.Generators())

	log.Fatal(kustomization.GenerateKustomization(resourceFiles))

	writeBaseFiles(resourceFiles)

}

func writeBaseFiles(files resources.Files) {
	_ = os.Chdir(directory)
	_ = os.Chdir("base")
	_ = files.Write()
}

func useDefault(def string, flag string) string {
	if len(flag) == 0 {
		return def
	}
	return flag
}

func createBase(app input.Application, files resources.Files, generators []kustomization.Generator) {
	for _, generate := range generators {
		err := generate(app, files)
		if err != nil {
			log.Fatalf("Could not create resource %v", err)
		}
	}
}

func createDirectory() {
	stat, err := os.Stat(directory)

	if err != nil && !os.IsNotExist(err) {
		log.Fatalf("Could not read local directory %v", err)
	}

	if os.IsExist(err) && stat.IsDir() {
		log.Warnf("Platform directory %s already exists", directory)
	}

	err = os.MkdirAll("./platform/base", 0755)

	log.Infof("Creating directory %s", directory)

	if err != nil {
		log.Fatalf("Cannot create platform directory %s %v", err)
	}
}
