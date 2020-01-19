package cmd

import (
	"fmt"
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/azunymous/easymodo/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

// imageCmd represents the version command
var imageCmd = &cobra.Command{
	Use:   "image [namespace]",
	Short: "Change the version of an application via kustomize patch",
	Long: `Create a kustomize overlay with a deployment patch changing the image of a container
in a deployment.

e.g easymodo modify image my-cool-app-production -i gcr.io/cool/my-app:v2.0.0

Outputs the directory the version kustomization is stored.
`,
	Run:     newVersionCommand,
	Args:    cobra.MaximumNArgs(1),
	Aliases: []string{"set", "change"},
}

var w = os.Stdout

var image string

func init() {
	modifyCmd.AddCommand(imageCmd)
	imageCmd.PersistentFlags().StringVarP(SuffixFlag(), "suffix", "s", "", "Suffix to use for the existing namespace kustomization directory")

	imageCmd.Flags().StringVarP(&image, "image", "i", "", "Image (required)")
	_ = imageCmd.MarkFlagRequired("image")
}

func newVersionCommand(c *cobra.Command, args []string) {
	resourceFiles := fs.NewFileMap()
	appName, appImage, appPort := input.GetBaseApp(fs.Get(), Directory())

	namespace, nsDir := input.ValidateNamespaceOrSuffix(Suffix(), appName, args, c)
	_, err := afero.ReadFile(fs.Get(), path.Join(Directory(), nsDir, "kustomization.yaml"))
	if err != nil {
		log.Fatalf("Could not open %s kustomization.yaml file: %v.", nsDir, err)
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
		Image:         image,
		ConfigPath:    configPath,
		Host:          Ingress(),
	}

	if appImage == image {
		log.Fatalf("Base image is the same as the input image: %s", appImage)
	}

	relativeBasePath := filepath.Join("../", nsDir)
	k.AddResource(relativeBasePath)

	_ = kustomization.Generate("deployment-image-patch", kustomization.DeploymentImagePatch())(application, resourceFiles)

	k.AddPatch("deployment-image-patch.yaml")

	kustomization.Create(&k, resourceFiles)
	resourceFiles.WriteAll(Directory(), nsDir+"-temp")
	abs, _ := filepath.Abs(path.Join(Directory(), nsDir+"-temp"))
	_, _ = fmt.Fprintln(w, abs)
}
