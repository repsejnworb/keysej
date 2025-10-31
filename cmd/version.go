package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() { rootCmd.AddCommand(versionCmd) }

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Printf("keysej %s\ncommit: %s\ndate:   %s\n", buildVersion, buildCommit, buildDate)
	},
}
