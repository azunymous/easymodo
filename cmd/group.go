package cmd

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	"github.com/azunymous/easymodo/kustomization"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
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
	Run:  newGroupCommand,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(groupCmd)

	groupCmd.Flags().StringArrayVarP(KustomizationsFlag(), "kustomization", "k", []string{}, "Kustomization folder to add to a new generated kustomization")
	_ = groupCmd.MarkFlagDirname("kustomization")
	_ = groupCmd.MarkFlagRequired("kustomization")

	groupCmd.Flags().BoolVarP(VerifyFlag(), "verify", "v", false, "Verify kustomizations exist")
	groupCmd.Flags().StringVarP(OutputFlag(), "output", "o", ".", "Output folder for kustomization file")
}

func newGroupCommand(_ *cobra.Command, _ []string) {
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
		if Verify() {
			exists, err := afero.DirExists(fs.Get(), kFolder)
			if err != nil {
				log.Fatalf("Error looking for directory %s: %v", kFolder, err)
			}
			if !exists {
				log.Errorf("%s does not exist or is not a directory", kFolder)
				panic(1)
			}
		}
		if path.IsAbs(kFolder) {
			kFolder = getRelativePathFor(outputDir, kFolder)
		} else if !path.IsAbs(outputDir) {
			tempOutput, _ := filepath.Abs(outputDir)
			tempKFolder, _ := filepath.Abs(kFolder)
			kFolder = getRelativePathFor(tempOutput, tempKFolder)
		}
		k.AddResource(kFolder)
	}

	kustomization.Create(&k, resourceFiles)
	resourceFiles.WriteAll("", outputDir)
	log.Info("Created kustomization yaml in ", outputDir)
}

func getRelativePathFor(baseDir string, kustomizeDir string) string {
	base := baseDir
	if !path.IsAbs(base) {
		base = getFullDirPath(base)
	}
	rel, err := filepath.Rel(base, kustomizeDir)
	if err != nil {
		log.Fatalf("Cannot calculate relative path from %s to %s due to %v", baseDir, kustomizeDir, err)
	}
	return rel
}

func getFullDirPath(d string) string {
	wd, err := filepath.Abs(d)
	if err != nil {
		log.Fatalf("Cannot get absolute path to output directory %v", err)
	}
	log.Infof("Using output directory %s", wd)
	return wd
}
