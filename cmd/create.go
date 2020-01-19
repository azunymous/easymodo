package cmd

import (
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create kustomize resources",
	Long: `Create or bootstrap kustomize resources for multiple environments.
For example:
easymodo create base my-cool-app
easymodo create overlay my-cool-app-production`,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
