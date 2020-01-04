package cmd

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/azunymous/easymodo/kustomization"
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

var configPath string

func init() {
	rootCmd.AddCommand(overlayCmd)

	overlayCmd.PersistentFlags().StringToStringVarP(ConfigFilesFlag(), "configFile", "c", nil, "Configuration filename and file for generating config maps")
	overlayCmd.PersistentFlags().StringVarP(&configPath, "configPath", "p", "/config/", "Configuration folder for mounting config map contents")

	overlayCmd.PersistentFlags().StringToStringVarP(SecretEnvsFlag(), "secretEnv", "e", nil, "Secret .env filename and env file for generating secrets")

	overlayCmd.Flags().StringVar(IngressFlag(), "ingress", "", "Enable ingress resource generation with given host")

	overlayCmd.PersistentFlags().StringVarP(SuffixFlag(), "suffix", "s", "", "Suffix to use for namespace for overlay")
	overlayCmd.Flags().BoolVarP(NamespaceResourceFlag(), "resource", "r", false, "Create namespace resource")
}

func overlayCommand(cmd *cobra.Command, args []string) {
	resourceFiles := fs.NewFileMap()
	appName, appPort := input.GetAppName(fs.Get(), Directory())

	var (
		namespace string
		nsDir     string
	)

	if Suffix() != "" {
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
		Secrets:   map[string][]string{},
		Namespace: namespace,
	}

	application := input.Application{
		Name:          appName,
		ContainerName: appName,
		ContainerPort: appPort,
		Namespace:     namespace,
		ConfigPath:    configPath,
		Host:          Ingress(),
	}

	relativeBasePath := filepath.Join("../", "base")
	k.AddResource(relativeBasePath)

	if NamespaceResource() {
		err := kustomization.Generate("namespace", kustomization.Namespace())(input.Application{Namespace: namespace}, resourceFiles)
		if err != nil {
			log.Warnf("Could not create namespace: %v", err)
		} else {
			k.AddResource("namespace.yaml")
		}
	}

	addConfigGenerator(application, resourceFiles, &k, appName)
	addSecretGenerator(application, resourceFiles, &k, appName)

	if Ingress() != "" {
		err := kustomization.Generate("ingress", kustomization.Ingress())(application, resourceFiles)
		if err != nil {
			log.Warnf("Could not create ingress: %v", err)
		} else {
			k.AddResource("ingress.yaml")
		}
	}

	kustomization.Create(&k, resourceFiles)
	resourceFiles.WriteAll(Directory(), nsDir)

}

func addConfigGenerator(application input.Application, resourceFiles fs.Files, k *input.Kustomization, appName string) {
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
}

func addSecretGenerator(application input.Application, resourceFiles fs.Files, k *input.Kustomization, appName string) {
	if len(SecretEnvs()) > 0 {
		err := kustomization.Generate("deployment-secret-patch", kustomization.DeploymentSecretPatch())(application, resourceFiles)
		if err != nil {
			log.Fatalf("could not create deployment patch with given secret file: %v", err)
		}

		k.AddPatch("deployment-secret-patch.yaml")

		for fileName, content := range SecretEnvs() {
			resourceFiles.Add(fileName, content)
			k.AddSecret(appName+"-secret", fileName)
		}
	}
}
