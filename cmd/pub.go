package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/repsejnworb/keysej/internal/osutil"
	"github.com/spf13/cobra"
)

var copyFlag bool

var pubCmd = &cobra.Command{
	Use:   "pub [name]",
	Short: "Print (and optionally copy) the public key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := osutil.ValidateName(name); err != nil {
			return err
		}
		pub := filepath.Join(os.Getenv("HOME"), ".ssh", "id_ed25519_"+name+".pub")
		b, err := os.ReadFile(pub)
		if err != nil {
			return err
		}
		fmt.Print(string(b))
		if copyFlag {
			return clipboard.WriteAll(string(b))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pubCmd)
	pubCmd.Flags().BoolVar(&copyFlag, "copy", false, "Copy to clipboard")
}
