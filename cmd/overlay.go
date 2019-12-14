package cmd

import (
	"github.com/azunymous/easymodo/cmd/fs"
	"github.com/azunymous/easymodo/cmd/input"
	"github.com/azunymous/easymodo/cmd/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path/filepath"
)

// overlayCmd represents the overlay command
var overlayCmd = &cobra.Command{
	Use:   "overlay [namespace]",
	Short: "Defines a kustomize overlay",
	Long: `Defines kustomization files for an overlay. 
This is intended for different environments e.g dev, stage or production, setting common
kustomization options via flags`,
	Run:  overlayCommand,
	Args: cobra.MaximumNArgs(1),
}

func init() {
	rootCmd.AddCommand(overlayCmd)

	overlayCmd.PersistentFlags().StringToStringVarP(ConfigFilesFlag(), "configFile", "c", nil, "Configuration filename and file for generating config maps")

	// Shared with add command
	overlayCmd.PersistentFlags().StringVarP(SuffixFlag(), "suffix", "s", "", "Suffix to use for namespace for overlay")
	overlayCmd.Flags().BoolVarP(NamespaceResourceFlag(), "resource", "r", false, "Create namespace resource")
}

func overlayCommand(cmd *cobra.Command, args []string) {
	resourceFiles := fs.NewFileMap()
	appName := input.GetAppName(fs.Get(), Directory())

	var (
		namespace string
		nsDir     string
	)

	if cmd.Flags().Changed("suffix") {
		namespace = appName + "-" + Suffix()
		nsDir = Suffix()
	} else if len(args) < 1 {
		println(cmd.UsageString())
		log.Fatalf("No namespace or namespace suffix provided!")
	} else {
		namespace = args[0]
		nsDir = args[0]
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

	if NamespaceResource() {
		err := kustomization.Generate("namespace", kustomization.Namespace())(input.Application{Namespace: namespace}, resourceFiles)
		if err != nil {
			log.Warnf("Could not create namespace: %v", err)
		} else {
			k.AddResource("namespace.yaml")
		}
	}

	if len(ConfigFiles()) > 0 {
		err := kustomization.Generate("deployment-config-patch", kustomization.DeploymentConfigPatch())(application, resourceFiles)
		if err != nil {
			log.Fatalf("could not create deployment patch with given config file: %v", err)
		}

		k.AddPatch("deployment-config-patch.yaml")

		for fileName, content := range ConfigFiles() {
			resourceFiles.Add(fileName, content)
			k.AddConfig(appName+"-config", fileName)
		}

	}

	kustomization.Create(&k, resourceFiles)
	resourceFiles.WriteAll(Directory(), nsDir)

}
