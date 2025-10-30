package cmd

import (
	"fmt"
	"os"

	"github.com/repsejnworb/keysej/internal/version"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "keysej",
	Short: "A tiny, secure SSH key helper",
	Long:  "keysej wraps OpenSSH for key generation, agent add, install, and config with a nice TUI.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true
	rootCmd.Version = version.String()
}
