package cmd

import (
	"github.com/azunymous/easymodo/fs"
	"github.com/azunymous/easymodo/input"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
)

// verify represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Verify kustomizations and directory structure",
	Long: `Verify kustomization files correctly build.
This command will build all kustomizations in the provided directory (default: platform).

Kustomize must be installed. kubectl kustomize is not currently supported.
`,
	Run:  newVerifyCommand,
	Args: cobra.NoArgs,
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}

func newVerifyCommand(_ *cobra.Command, _ []string) {
	_, err := exec.LookPath("kustomize")
	if err != nil {
		log.Fatalf("kustomize is not installed")
	}

	appName, _, _ := input.GetBaseApp(fs.Get(), Directory())

	log.Infof("Verifying %s directory for application %s", Directory(), appName)

	err = afero.Walk(fs.Get(), Directory(), returnWalkFunc())

	if err != nil {
		log.Fatalf("Could not traverse directory: %v", err)
	}
	log.Infof("SUCCESS")
}

func returnWalkFunc() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() || path == Directory() {
			return nil
		}

		if kustExists, err := afero.Exists(fs.Get(), filepath.Join(path, "kustomization.yaml")); !kustExists && err == nil {
			log.Infof("Treating %s as context directory", path)
			return nil
		}

		log.Infof("Building kustomization %s", path)
		kustomize := exec.Command("kustomize", "build", path)
		kustomize.Stderr = os.Stderr
		err = kustomize.Run()
		if err != nil {
			log.Fatalf("Failed to build kustomization %s: %v", path, err)
		}
		return nil
	}
}
