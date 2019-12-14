package cmd

import (
	"fmt"
	"github.com/azunymous/easymodo/cmd/fs"
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path/filepath"
)

// overlayCmd represents the overlay command
var overlayCmd = &cobra.Command{
	Use:   "overlay",
	Short: "Defines a kustomize overlay",
	Long: `Defines kustomization files for an overlay. 
This is intended for different environments e.g dev, stage or production, setting common
kustomization options via flags`,
	Run: overlayCommand,
}

func init() {
	rootCmd.AddCommand(overlayCmd)

	overlayCmd.PersistentFlags().StringToStringVarP(input.ConfigFilesFlag(), "configFile", "c", nil, "Configuration filename and file for generating config maps")

	// Shared with add command
	overlayCmd.PersistentFlags().StringVarP(input.SuffixFlag(), "suffix", "s", "dev", "Suffix to use for namespace for overlay")
	overlayCmd.PersistentFlags().StringVarP(input.NamespaceFlag(), "namespace", "n", "", "Specify a full namespace instead of using a suffix")
	overlayCmd.Flags().BoolVarP(input.NamespaceResourceFlag(), "resource", "r", false, "Create namespace resource")
}

func overlayCommand(cmd *cobra.Command, args []string) {
	fmt.Println("overlay called")
	resourceFiles := fs.NewFileMap()

	// TODO remove below duplicated code between add and overlay
	appName := input.GetAppName(fs.Get(), input.Directory())
	namespace := appName + "-" + input.Suffix()
	nsDir := input.Suffix()

	if cmd.Flags().Changed("namespace") {
		namespace = input.Namespace()
		nsDir = input.Namespace()
	}

	k := input.Kustomization{
		Res:       []string{},
		Patches:   []string{},
		Config:    map[string][]string{},
		Namespace: namespace,
	}

	application := input.Application{
		Name:          appName,
		ContainerName: appName,
		Namespace:     namespace,
	}

	relativeBasePath := filepath.Join("../", "base", "kustomization.yaml")
	k.AddResource(relativeBasePath)

	if input.NamespaceResource() {
		err := kustomization.Generate("namespace", kustomization.Namespace())(input.Application{Namespace: namespace}, resourceFiles)
		if err != nil {
			log.Warnf("Could not create namespace: %v", err)
		} else {
			k.AddResource("namespace.yaml")
		}
	}

	if len(input.ConfigFiles()) > 0 {
		err := kustomization.Generate("deployment-config-patch", kustomization.DeploymentConfigPatch())(application, resourceFiles)
		if err != nil {
			log.Fatalf("could not create deployment patch with given config file: %v", err)
		}

		k.AddPatch("deployment-config-patch.yaml")

		for fileName, content := range input.ConfigFiles() {
			resourceFiles.Add(fileName, content)
			k.AddConfig(appName+"-config", fileName)
		}

	}

	kustomization.Create(&k, resourceFiles)
	fs.WriteAll(resourceFiles, input.Directory(), nsDir)

}
