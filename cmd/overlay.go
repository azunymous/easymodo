package cmd

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/azunymous/easymodo/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"path"
	"path/filepath"
)

// overlayCmd represents the overlay command for creating a kustomize overlay
var overlayCmd = &cobra.Command{
	Use:   "overlay [namespace]",
	Short: "Defines a kustomize overlay",
	Long: `Defines kustomization files for an overlay. 
This is intended for different environments e.g dev, stage or production, setting common
kustomization options via flags`,
	Run:  newOverlayCommand,
	Args: cobra.MaximumNArgs(1),
}

var configPath string

func init() {
	createCmd.AddCommand(overlayCmd)

	overlayCmd.PersistentFlags().StringToStringVarP(ConfigFilesFlag(), "configFile", "c", nil, "Configuration filename and file for generating config maps")
	overlayCmd.PersistentFlags().StringVarP(&configPath, "configPath", "p", "/config/", "Configuration folder for mounting config map contents")

	overlayCmd.PersistentFlags().StringToStringVarP(SecretEnvsFlag(), "secretEnv", "e", nil, "Secret .env filename and env file for generating secrets")

	overlayCmd.Flags().StringVar(IngressFlag(), "ingress", "", "Enable ingress resource generation with given host")
	overlayCmd.Flags().IntVarP(ReplicasFlag(), "replicas", "r", 1, "Enable ingress resource generation with given host")
	overlayCmd.Flags().StringToStringVar(LimitsFlag(), "limits", map[string]string{}, "The resource requirement limits for this container.  For example, 'cpu=200m,memory=512Mi'")
	overlayCmd.Flags().StringToStringVar(RequestsFlag(), "requests", map[string]string{}, "The resource requirement requests for this container.  For example, 'cpu=200m,memory=512Mi'")

	overlayCmd.PersistentFlags().StringVarP(SuffixFlag(), "suffix", "s", "", "Suffix to use for namespace for overlay")
	overlayCmd.Flags().BoolVarP(NamespaceResourceFlag(), "namespace-resource", "n", false, "Create namespace resource")
}

func newOverlayCommand(c *cobra.Command, args []string) {
	resourceFiles := fs.NewFileMap()
	appName, _, appPort := input.GetBaseApp(fs.Get(), Directory())

	namespace, nsDir := input.ValidateNamespaceOrSuffix(Suffix(), appName, args, c)
	validateContainerResources(Requests(), "Requests")
	validateContainerResources(Limits(), "Limits")

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
		Replicas:      Replicas(),
	}

	k.AddResource(relativeBasePath())

	if NamespaceResource() {
		err := kustomization.Generate("namespace", kustomization.Namespace())(input.Application{Namespace: namespace}, resourceFiles)
		if err != nil {
			log.Warnf("Could not create namespace: %v", err)
		} else {
			k.AddResource("namespace.yaml")
		}
	}

	addContainerResourceGenerator(application, resourceFiles, &k)
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

	if Replicas() != 1 {
		err := kustomization.Generate("deployment-replica-patch", kustomization.DeploymentReplicaPatch())(application, resourceFiles)
		if err != nil {
			log.Fatalf("Could not create replica patch %v", err)
		} else {
			k.AddPatch("deployment-replica-patch.yaml")
		}
	}

	kustomization.Create(&k, resourceFiles)
	resourceFiles.WriteAll(Directory(), addContext(nsDir))

}

func addContainerResourceGenerator(application input.Application, resourceFiles fs.Files, k *input.Kustomization) {
	if len(Requests()) > 0 {
		setContainerResource(Requests(), "cpu", &application.CpuRequests)
		setContainerResource(Requests(), "memory", &application.MemoryRequests)
		err := kustomization.Generate("deployment-requests-patch", kustomization.DeploymentRequestsPatch())(application, resourceFiles)
		if err != nil {
			log.Warnf("Could not create request patch: %v", err)
		} else {
			k.AddPatch("deployment-requests-patch.yaml")
		}
	}
	if len(Limits()) > 0 {
		setContainerResource(Limits(), "cpu", &application.CpuLimits)
		setContainerResource(Limits(), "memory", &application.MemoryLimits)
		err := kustomization.Generate("deployment-limits-patch", kustomization.DeploymentLimitsPatch())(application, resourceFiles)
		if err != nil {
			log.Warnf("Could not create limits patch: %v", err)
		} else {
			k.AddPatch("deployment-limits-patch.yaml")
		}
	}
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

func validateContainerResources(m map[string]string, name string) {
	if len(m) > 0 && len(m) > 2 {
		log.Fatalf("%s flag is not correctly defined. Too many elements set, expected only memory/cpu.", name)
	}
}

func setContainerResource(m map[string]string, cpuOrMemory string, valuePtr *string) {
	if val, ok := m[cpuOrMemory]; ok {
		*valuePtr = val
	}
}

//
// Context specific, should be moved
//

func relativeBasePath() string {
	if Context() == "" {
		return filepath.Join("../", "base")
	}
	return filepath.Join("../../", "base")
}

func addContext(dir string) string {
	return path.Join(Context(), dir)
}
