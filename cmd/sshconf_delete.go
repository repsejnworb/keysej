package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/sshconf"
)

var sshconfDeleteCmd = &cobra.Command{
	Use:   "delete <key> [host-or-cidr]",
	Short: "Delete a key's entire file or just one rule by host/cidr",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		path := sshconf.FileForKey(key)

		if len(args) == 1 {
			// delete the whole file (backup first)
			if _, err := os.Stat(path); err != nil {
				return fmt.Errorf("no such conf: %s", path)
			}
			tmp, err := sshconf.BackupFile(path)
			if err != nil {
				return err
			}
			if !askYesNo(fmt.Sprintf("Delete %s? (backed up to %s) [y/N]: ", path, tmp)) {
				fmt.Println("aborted")
				return nil
			}
			return os.Remove(path)
		}

		// delete a specific block by pattern (host or cidr)
		pattern := args[1]
		_, marker := sshconf.RenderBlock(key, pattern, "", false) // just to get normalized marker
		tmp, changed, err := sshconf.DeleteBlock(path, marker)
		if err != nil {
			return err
		}
		if !changed {
			fmt.Println("no such rule; use `sshconf list", key, "` to inspect")
			return nil
		}
		if !askYesNo(fmt.Sprintf("Delete rule %q in %s? (backup: %s) [y/N]: ", pattern, path, tmp)) {
			fmt.Println("aborted")
			return nil
		}
		return sshconf.CommitDelete(path)
	},
}

func askYesNo(prompt string) bool {
	fmt.Print(prompt)
	r := bufio.NewReader(os.Stdin)
	s, _ := r.ReadString('\n')
	s = strings.TrimSpace(strings.ToLower(s))
	return s == "y" || s == "yes"
}

func init() {
	sshconfCmd.AddCommand(sshconfDeleteCmd)
}
