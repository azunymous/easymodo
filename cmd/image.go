package cmd

import (
	"fmt"
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/azunymous/easymodo/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"path"
	"path/filepath"
	"regexp"

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

func init() {
	modifyCmd.AddCommand(imageCmd)
	imageCmd.PersistentFlags().StringVarP(SuffixFlag(), "suffix", "s", "", "Suffix to use for the existing namespace kustomization directory")

	imageCmd.Flags().StringVarP(ImageFlag(), "image", "i", "", "Image (required)")
	_ = imageCmd.MarkFlagRequired("image")

	imageCmd.Flags().StringVarP(OutputFlag(), "output", "o", "", "Output folder for kustomization files. Defaults to '<namespace folder name>-<version>'")
}

func newVersionCommand(c *cobra.Command, args []string) {
	resourceFiles := fs.NewFileMap()
	appName, appImage, appPort := input.GetBaseApp(fs.Get(), Directory())
	outputDir := Output()

	namespace, nsDir := input.ValidateNamespaceOrSuffix(Suffix(), appName, args, c)
	_, err := afero.ReadFile(fs.Get(), path.Join(Directory(), Context(), nsDir, "kustomization.yaml"))
	if err != nil {
		log.Fatalf("Could not open %s kustomization.yaml file: %v.", nsDir, err)
	}

	if outputDir == "" {
		version := parseVersion(Image())
		outputDir = path.Join(Context(), nsDir+"-"+version)
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
		Image:         Image(),
		ConfigPath:    configPath,
		Host:          Ingress(),
	}

	if appImage == Image() {
		log.Fatalf("Base image is the same as the input image: %s", appImage)
	}

	relativeBasePath := filepath.Join("../", nsDir)
	k.AddResource(relativeBasePath)

	_ = kustomization.Generate("deployment-image-patch", kustomization.DeploymentImagePatch())(application, resourceFiles)

	k.AddPatch("deployment-image-patch.yaml")

	kustomization.Create(&k, resourceFiles)
	resourceFiles.WriteAll(Directory(), outputDir)
	abs, _ := filepath.Abs(path.Join(Directory(), outputDir))
	_, _ = fmt.Fprintln(w, abs)
}

var tag = regexp.MustCompile(`.+:(.+)`)

func parseVersion(image string) string {
	m := tag.FindStringSubmatch(image)
	if len(m) < 2 {
		return "UNKNOWN"
	}
	return m[1]
}
