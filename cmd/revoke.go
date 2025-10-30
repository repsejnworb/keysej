package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/repsejnworb/keysej/internal/osutil"
	"github.com/repsejnworb/keysej/internal/shell"
	"github.com/spf13/cobra"
)

var deleteFiles bool

var revokeCmd = &cobra.Command{
	Use:   "revoke [name]",
	Short: "Remove key from agent/keychain and optionally delete files",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := osutil.ValidateName(name); err != nil {
			return err
		}
		key := filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519_"+name)
		ctx, cancel := context.WithTimeout(cmd.Context(), 20*time.Second)
		defer cancel()
		if err := shell.AgentRemove(ctx, key); err != nil {
			return err
		}
		if osutil.IsDarwin() {
			_ = shell.MacDeleteKeychainEntry(key)
		}
		if deleteFiles {
			_ = os.Remove(key)
			_ = os.Remove(key + ".pub")
			fmt.Println("🗑️  Deleted key files.")
		}
		fmt.Println("✅ Revoked. Remember to remove the public key from remote authorized_keys.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(revokeCmd)
	revokeCmd.Flags().BoolVar(&deleteFiles, "delete-files", false, "Delete private and public key files")
}
