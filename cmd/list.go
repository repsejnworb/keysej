package cmd

import (
	"fmt"

	"github.com/repsejnworb/keysej/internal/shell"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List keys loaded in ssh-agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		out, err := shell.AgentList(cmd.Context())
		if err != nil {
			return err
		}
		fmt.Print(out)
		return nil
	},
}

func init() { rootCmd.AddCommand(listCmd) }
