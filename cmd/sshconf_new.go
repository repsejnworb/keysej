package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/sshconf"
)

var (
	newUser    string
	newForward bool
)

var sshconfNewCmd = &cobra.Command{
	Use:   "new <key> <host-or-cidr>",
	Short: "Create or update a rule for a key (Host pattern or CIDR Match)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		pattern := args[1]

		frag, marker := sshconf.RenderBlock(key, pattern, newUser, newForward)
		path := sshconf.FileForKey(key)

		changed, err := sshconf.UpsertBlock(path, marker, frag)
		if err != nil {
			return err
		}

		if changed {
			fmt.Printf("✓ wrote %s\n", path)
		} else {
			fmt.Printf("= no change %s (identical)\n", path)
		}
		return nil
	},
}

func init() {
	sshconfCmd.AddCommand(sshconfNewCmd)
	sshconfNewCmd.Flags().StringVar(&newUser, "user", "", "Set User <name> in the block")
	sshconfNewCmd.Flags().BoolVar(&newForward, "forward", false, "Set ForwardAgent yes (default no)")
}
