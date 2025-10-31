// internal/cli/style/style.go
package style

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	RedBold = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	// Use a mid-gray instead of Faint(); better contrast
	HintStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	Code      = lipgloss.NewStyle().Foreground(lipgloss.Color("111")).Render
)

func Error(msg string) { fmt.Fprintln(os.Stderr, RedBold.Render("✖ ")+msg) }
func Hint(msg string)  { fmt.Fprintln(os.Stderr, HintStyle.Render(msg)) }

func UsageLine(cmd *cobra.Command) string { return Code(cmd.UseLine()) }
