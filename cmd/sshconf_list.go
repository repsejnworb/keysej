package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/sshconf"
)

var listHost string

var sshconfListCmd = &cobra.Command{
	Use:   "list [key]",
	Short: "List rules for a key (or all), optionally filtered by --host",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			key := args[0]
			path := sshconf.FileForKey(key)
			b, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read %s: %w", path, err)
			}
			out := sshconf.FilterByHostOrAll(string(b), listHost)
			fmt.Print(out)
			return nil
		}
		// list all keysej.* files
		files, err := sshconf.AllKeysejFiles()
		if err != nil {
			return err
		}
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				continue
			}
			out := sshconf.FilterByHostOrAll(string(b), listHost)
			if out != "" {
				fmt.Printf("# %s\n%s\n", f, out)
			}
		}
		return nil
	},
}

func init() {
	sshconfCmd.AddCommand(sshconfListCmd)
	sshconfListCmd.Flags().StringVar(&listHost, "host", "", "Filter rules matching this host/ip pattern")
}
