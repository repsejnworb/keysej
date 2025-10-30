package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var names string

var configSnippetCmd = &cobra.Command{
	Use:   "config-snippet",
	Short: "Print a sane ~/.ssh/config snippet for given names",
	RunE: func(cmd *cobra.Command, args []string) error {
		ids := []string{}
		for _, n := range strings.Split(names, ",") {
			n = strings.TrimSpace(n)
			if n == "" {
				continue
			}
			ids = append(ids, fmt.Sprintf("IdentityFile %s/.ssh/id_ed25519_%s", os.Getenv("HOME"), n))
		}
		fmt.Println("Host *")
		fmt.Println("  AddKeysToAgent yes")
		fmt.Println("  IdentitiesOnly yes")
		if runtime.GOOS == "darwin" {
			fmt.Println("  UseKeychain yes")
		}
		for _, l := range ids {
			fmt.Println("  ", l)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configSnippetCmd)
	configSnippetCmd.Flags().StringVar(&names, "names", "", "comma-separated key names (e.g. work,personal)")
}
