package cmd

import (
	"context"
	"path/filepath"
	"time"

	"github.com/repsejnworb/keysej/internal/config"
	"github.com/repsejnworb/keysej/internal/osutil"
	"github.com/repsejnworb/keysej/internal/shell"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install [name] [user@host]",
	Short: "Install the public key on a remote host (ssh-copy-id if available)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, target := args[0], args[1]
		if err := osutil.ValidateName(name); err != nil {
			return err
		}
		pub := filepath.Join(config.C.SSHDir, "id_ed25519_"+name+".pub")
		ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
		defer cancel()
		return shell.InstallPubkey(ctx, pub, target)
	},
}

func init() { rootCmd.AddCommand(installCmd) }
