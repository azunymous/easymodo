package cmd

import (
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	"github.com/azunymous/easymodo/cmd/resources"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add namespaced overlay on top of a kustomize base",
	Long: `Creates an overlay kustomization from the base using
a namespace formed of the [app name]-[suffix].`,
	Run:  add,
	Args: cobra.MaximumNArgs(0),
}

var suffix string
var flagNamespace string
var nsResource bool

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(&suffix, "suffix", "s", "dev", "Suffix to use for namespace for overlay")
	addCmd.PersistentFlags().StringVarP(&flagNamespace, "namespace", "n", "", "Specify a full namespace instead of using a suffix")

	addCmd.Flags().BoolVarP(&nsResource, "resource", "r", false, "Create namespace resource")
}
func add(cmd *cobra.Command, args []string) {
	exists, err := afero.DirExists(appFs, path.Join(directory, "base"))

	if err != nil || !exists {
		log.Fatalf("Base directory %s doesn't exist", path.Join(directory, "base"))
	}

	var namespace string
	var dir string

	if flagNamespace == "" {
		namespace = input.GetAppName(appFs, directory) + "-" + suffix
		dir = suffix
	} else {
		namespace = flagNamespace
		dir = flagNamespace
	}

	relativeBasePath := filepath.Join("../", "base", "kustomization.yaml")
	resourceFiles := resources.NewFileMap()
	res := []string{relativeBasePath}

	if nsResource {
		err := kustomization.Generate("namespace", kustomization.Namespace())(input.Application{Namespace: namespace}, resourceFiles)
		if err != nil {
			log.Warnf("Could not create namespace: %v", err)
		} else {
			res = append(res, "namespace.yaml")
		}
	}

	createKustomization(input.NewKustomization(res, namespace), resourceFiles)
	writeFiles(resourceFiles, dir)
}
