package cmd

import (
	"github.com/azunymous/easymodo/cmd/fs"
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

// addCmd represents the newAddCommand command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add namespaced overlay on top of a kustomize base",
	Long: `Creates an overlay kustomization from the base using
a namespace formed of the [app name]-[suffix].`,
	Run:  newAddCommand,
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
func newAddCommand(cmd *cobra.Command, args []string) {
	exists, err := afero.DirExists(fs.Get(), path.Join(input.Directory(), "base"))

	if err != nil || !exists {
		log.Fatalf("Base directory %s doesn't exist", path.Join(input.Directory(), "base"))
	}

	var namespace string
	var nsDir string

	if flagNamespace == "" {
		namespace = input.GetAppName(fs.Get(), input.Directory()) + "-" + suffix
		nsDir = suffix
	} else {
		namespace = flagNamespace
		nsDir = flagNamespace
	}

	exists, err = afero.DirExists(fs.Get(), path.Join(input.Directory(), nsDir))

	if err == nil && exists {
		log.Fatalf("%s directory (%s) already exists", nsDir, path.Join(input.Directory(), nsDir))
	}

	relativeBasePath := filepath.Join("../", "base", "kustomization.yaml")
	resourceFiles := fs.NewFileMap()
	res := []string{relativeBasePath}

	if nsResource {
		err := kustomization.Generate("namespace", kustomization.Namespace())(input.Application{Namespace: namespace}, resourceFiles)
		if err != nil {
			log.Warnf("Could not create namespace: %v", err)
		} else {
			res = append(res, "namespace.yaml")
		}
	}

	kustomization.Create(input.NewKustomization(res, namespace), resourceFiles)
	fs.WriteAll(resourceFiles, input.Directory(), nsDir)
}
