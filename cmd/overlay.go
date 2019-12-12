package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// overlayCmd represents the overlay command
var overlayCmd = &cobra.Command{
	Use:   "overlay",
	Short: "Defines a kustomize overlay",
	Long: `Defines kustomization files for an overlay. 
This is intended for different environments e.g dev, stage or production, setting common
kustomization options via flags`,
	Run: overlayCommand,
}

func init() {
	rootCmd.AddCommand(overlayCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// overlayCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// overlayCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func overlayCommand(cmd *cobra.Command, args []string) {
	fmt.Println("overlay called")
}
