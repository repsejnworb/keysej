package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/config"
)

var (
	flagSSHDir string
)

var rootCmd = &cobra.Command{
	Use:   "keysej",
	Short: "A tiny, secure SSH key helper",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize config once for all subcommands.
		return config.Init(flagSSHDir)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagSSHDir, "ssh-dir", "", "SSH directory (default $HOME/.ssh or $KEYSEJ_SSH_DIR)")
}
