package cmd

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/azunymous/easymodo/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
)

// group represents the group command
var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Create kustomization from multiple kustomizations",
	Long: `Create a kustomize file specifying other kustomize files for the purpose of deploying
a set of applications together, such as applications with their mocks, databases and other 
dependencies.

Specified files are expected to be available relative to the output directory.
'
`,
	Run:     newGroupCommand,
	Args:    cobra.NoArgs,
	Aliases: []string{"set", "change"},
}

func init() {
	rootCmd.AddCommand(groupCmd)

	groupCmd.Flags().StringArrayVarP(KustomizationsFlag(), "kustomization", "k", []string{}, "Kustomization folder to add to a new generated kustomization")
	_ = groupCmd.MarkFlagDirname("kustomization")
	_ = groupCmd.MarkFlagRequired("kustomization")

	groupCmd.Flags().StringVarP(OutputFlag(), "output", "o", ".", "Output folder for kustomization file.")
}

func newGroupCommand(c *cobra.Command, args []string) {
	resourceFiles := fs.NewFileMap()
	outputDir := path.Clean(Output())

	k := input.Kustomization{
		Res:     []string{},
		Patches: []string{},
		Config:  map[string][]string{},
		Secrets: map[string][]string{},
	}

	for _, kFolder := range Kustomizations() {
		kFolder = path.Clean(kFolder)
		if path.IsAbs(kFolder) {
			kFolder = getRelativePathFor(outputDir, kFolder)
		}
		k.AddResource(kFolder)
	}

	kustomization.Create(&k, resourceFiles)
	resourceFiles.WriteAll("", outputDir)
	log.Info("Created kustomization yaml in ", outputDir)
}

func getRelativePathFor(baseDir string, kustomizeDir string) string {
	base := baseDir
	if base == "." {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Cannot get working directory %v", err)
		}
		log.Infof("Using working directory %s", wd)
		base = wd
	}
	rel, err := filepath.Rel(base, kustomizeDir)
	if err != nil {
		log.Fatalf("Cannot calculate relative path from %s to %s due to %v", baseDir, kustomizeDir, err)
	}
	return rel
}
