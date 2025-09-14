package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print version information including build commit and build time.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("tailscale_exporter version %s (commit: %s, built: %s)\n", version, commit, buildTime)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

