package cmd

import (
	"fmt"

	"github.com/repsejnworb/keysej/internal/sshconf"
	"github.com/spf13/cobra"
)

var sshconfValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate keysej-managed SSH config fragments",
	RunE: func(cmd *cobra.Command, args []string) error {
		issues := sshconf.ValidateAll()
		if len(issues) == 0 {
			fmt.Println("✓ sshconf: all keysej configs look good")
			return nil
		}
		for _, i := range issues {
			fmt.Printf("✗ %s\n", i)
		}
		return fmt.Errorf("validation failed: %d issue(s)", len(issues))
	},
}

func init() { sshconfCmd.AddCommand(sshconfValidateCmd) }
