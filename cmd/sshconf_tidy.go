package cmd

import (
	"fmt"

	"github.com/repsejnworb/keysej/internal/sshconf"
	"github.com/spf13/cobra"
)

var sshconfTidyCmd = &cobra.Command{
	Use:   "tidy",
	Short: "Normalize keysej config files (sort blocks, trim whitespace, ensure newline)",
	RunE: func(cmd *cobra.Command, args []string) error {
		n, err := sshconf.TidyAll()
		if err != nil {
			return err
		}
		fmt.Printf("✓ tidied %d file(s)\n", n)
		return nil
	},
}

func init() { sshconfCmd.AddCommand(sshconfTidyCmd) }
