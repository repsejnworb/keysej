package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/sshconf"
)

var listHost string

var (
	warnStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	okStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	dimStyle  = lipgloss.NewStyle().Faint(true)
	pathStyle = lipgloss.NewStyle().Italic(true)
	headStyle = lipgloss.NewStyle().Bold(true).Underline(true)
)

var sshconfListCmd = &cobra.Command{
	Use:   "list [key]",
	Short: "List rules for a key (or all), optionally filtered by --host",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// small safety timeout (kept from earlier pattern)
		_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if len(args) == 1 {
			key := args[0]
			path := sshconf.FileForKey(key)

			// If the file doesn’t exist, print a friendly hint and exit 0.
			if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
				fmt.Println(warnStyle.Render("No sshconf file for key:"), key)
				fmt.Printf("%s %s\n", dimStyle.Render("Expected:"), pathStyle.Render(path))
				fmt.Println(dimStyle.Render("Create one with:"), okStyle.Render(fmt.Sprintf("keysej sshconf new %s <host-or-cidr>", key)))
				return nil
			} else if err != nil {
				// Real filesystem error
				return fmt.Errorf("cannot stat %s: %w", path, err)
			}

			b, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("read %s: %w", path, err)
			}
			out := sshconf.FilterByHostOrAll(string(b), listHost)
			out = strings.TrimSpace(out)
			if out == "" {
				if listHost != "" {
					fmt.Println(dimStyle.Render("No matching rules for host:"), pathStyle.Render(listHost))
				} else {
					fmt.Println(dimStyle.Render("File exists but contains no keysej-managed blocks yet."))
					fmt.Println(dimStyle.Render("Add one with:"), okStyle.Render(fmt.Sprintf("keysej sshconf new %s <host-or-cidr>", key)))
				}
				return nil
			}
			fmt.Println(headStyle.Render(path))
			fmt.Println(out)
			return nil
		}

		// No key given: list all keysej files (optionally filter by --host)
		files, err := sshconf.AllKeysejFiles()
		if err != nil {
			return err
		}
		if len(files) == 0 {
			fmt.Println(dimStyle.Render("No keysej sshconf files found under ~/.ssh/config.d"))
			fmt.Println(dimStyle.Render("Create one with:"), okStyle.Render("keysej sshconf new <key> <host-or-cidr>"))
			return nil
		}

		shown := 0
		for _, f := range files {
			b, err := os.ReadFile(f)
			if err != nil {
				fmt.Println(warnStyle.Render("skip:"), f, dimStyle.Render(err.Error()))
				continue
			}
			out := strings.TrimSpace(sshconf.FilterByHostOrAll(string(b), listHost))
			if out == "" {
				continue
			}
			if shown > 0 {
				fmt.Println() // spacing between files
			}
			fmt.Println(headStyle.Render(f))
			fmt.Println(out)
			shown++
		}
		if shown == 0 {
			if listHost != "" {
				fmt.Println(dimStyle.Render("No keysej rules match host:"), pathStyle.Render(listHost))
			} else {
				fmt.Println(dimStyle.Render("No keysej rules found."))
			}
		}
		return nil
	},
}

func init() {
	sshconfCmd.AddCommand(sshconfListCmd)
	sshconfListCmd.Flags().StringVar(&listHost, "host", "", "Filter rules matching this host/ip pattern")
}
