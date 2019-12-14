package cmd

import (
	"github.com/azunymous/easymodo/cmd/fs"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"path"
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

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.PersistentFlags().StringVarP(SuffixFlag(), "suffix", "s", "dev", "Suffix to use for namespace for overlay")
	addCmd.PersistentFlags().StringVarP(NamespaceFlag(), "namespace", "n", "", "Specify a full namespace instead of using a suffix")

	addCmd.Flags().BoolVarP(NamespaceResourceFlag(), "resource", "r", false, "Create namespace resource")
}
func newAddCommand(cmd *cobra.Command, args []string) {
	_, _, nsDir := appNameAndNamespaceFromFlags(cmd)

	exists, err := afero.DirExists(fs.Get(), path.Join(Directory(), nsDir))

	if err == nil && exists && !Force() {
		log.Fatalf("%s directory (%s) already exists", nsDir, path.Join(Directory(), nsDir))
	}

	overlayCommand(cmd, args)
}
