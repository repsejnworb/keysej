package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/repsejnworb/keysej/internal/osutil"
	"github.com/repsejnworb/keysej/internal/shell"
	"github.com/repsejnworb/keysej/internal/tui"
	"github.com/spf13/cobra"
)

var (
	flagTTL string
	dryRun  bool
)

var newCmd = &cobra.Command{
	Use:   "new [name]",
	Short: "Create a new purpose-scoped SSH key (one-prompt flow)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := osutil.ValidateName(name); err != nil {
			return err
		}

		// Init SSHDir (temp when dry-run) and set shell dry-run
		if err := osutil.InitSSHDir(dryRun); err != nil {
			return err
		}
		shell.DryRun = dryRun

		u, _ := user.Current()
		host := osutil.Hostname()
		date := time.Now().Format("2006-01-02")
		comment := fmt.Sprintf("%s@%s::%s::keysej::%s", u.Username, host, date, name)

		keyPath := filepath.Join(osutil.SSHDir, "id_ed25519_"+name)
		pubPath := keyPath + ".pub"

		m := tui.NewModel(name, comment, keyPath, pubPath, flagTTL)

		// No alt-screen, as you prefer
		p := tea.NewProgram(m)
		if _, err := p.Run(); err != nil {
			return err
		}

		// IMPORTANT: check the mutated pointer instance
		if m.Aborted() || !m.Confirmed() {
			fmt.Println("✖ Cancelled. No changes made.")
			return nil
		}

		// After TUI finishes, perform the side effects using captured passphrase (already zeroed in model post-add)
		if m.Aborted() {
			return nil
		}

		// 1) mkdir -p ~/.ssh ; chmod 700
		if err := os.MkdirAll(filepath.Join(osutil.SSHDir), 0o700); err != nil {
			return err
		}

		// 2) ssh-keygen -t ed25519 -o -a 100 -C comment -f keyPath -N pass
		ctx, cancel := context.WithTimeout(cmd.Context(), 30*time.Second)
		defer cancel()
		if err := shell.Keygen(ctx, keyPath, comment, m.Passphrase()); err != nil {
			return err
		}

		// 3) ssh-add (macOS: --apple-use-keychain) with TTL
		if err := shell.AgentAdd(ctx, keyPath, flagTTL, m.Passphrase(), dryRun); err != nil {
			return err
		}
		m.ZeroPass() // wipe

		// 4) Offer config snippet (append or copy) and copy pubkey already handled inside TUI via subsequent screens, or print here
		fmt.Printf("\n✅ Created and added: %s\n", keyPath)
		fmt.Printf("🔖 Comment: %s\n", comment)
		fp, _ := shell.Fingerprint(ctx, pubPath)
		if strings.TrimSpace(fp) != "" {
			fmt.Printf("🔑 Fingerprint: %s\n", strings.TrimSpace(fp))
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVar(&flagTTL, "ttl", "0", "Agent TTL (e.g. 1h); 0 = no TTL")
	newCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Run safely in temp dir and echo all commands except ssh-keygen")

}
