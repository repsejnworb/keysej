package cmd

import "github.com/spf13/cobra"

var sshconfCmd = &cobra.Command{
	Use:   "sshconf",
	Short: "Manage ~/.ssh/config.d/keysej.<key>.conf rules",
}

func init() {
	rootCmd.AddCommand(sshconfCmd)
}
