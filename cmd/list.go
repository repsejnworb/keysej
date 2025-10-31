package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/shell"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show currently loaded SSH keys and fingerprints",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		out, err := shell.AgentList(ctx)
		if err != nil {
			if strings.Contains(out, "The agent has no identities.") {
				fmt.Println("❕ No keys loaded in SSH agent.")
				return nil
			}
			return fmt.Errorf("ssh-add -l failed: %v", err)
		}

		out = strings.TrimSpace(out)
		if out == "" {
			fmt.Println("❕ No keys loaded in SSH agent.")
			return nil
		}

		fmt.Println(renderList(out))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func renderList(out string) string {
	rows := strings.Split(out, "\n")

	headerStyle := lipgloss.NewStyle().Bold(true).Underline(true)
	rowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("111"))
	commentStyle := lipgloss.NewStyle().Faint(true)

	var b strings.Builder
	fmt.Fprintf(&b, "%s\n", headerStyle.Render("Loaded SSH keys:"))
	for _, r := range rows {
		r = strings.TrimSpace(r)
		if r == "" {
			continue
		}
		parts := strings.Fields(r)
		if len(parts) < 3 {
			b.WriteString(rowStyle.Render(r) + "\n")
			continue
		}

		b.WriteString(
			fmt.Sprintf("%s %s %s\n",
				rowStyle.Render(parts[0]),                         // bits
				keyStyle.Render(parts[1]),                         // fingerprint
				commentStyle.Render(strings.Join(parts[2:], " ")), // comment
			),
		)
	}
	return b.String()
}
