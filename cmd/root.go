package cmd

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/cli/style"
)

var rootCmd = &cobra.Command{
	Use:   "keysej",
	Short: "A tiny, secure SSH key helper",
}

func init() {
	// Optional: friendlier flag-parse errors (unknown/missing flags)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		style.Error(err.Error())
		style.Hint("Usage: " + style.UsageLine(cmd))
		style.Hint("Help : " + style.Code(cmd.CommandPath()+" --help"))
		return err
	})
}

func Execute() {
	applySilence(rootCmd)

	cand, _, _ := rootCmd.Find(os.Args[1:])
	if cand == nil {
		cand = rootCmd
	}

	if err := rootCmd.Execute(); err != nil {
		msg := err.Error()

		// Rewrite Cobra's default "accepts X arg(s)" errors into friendlier text
		if acceptsArgsRE.MatchString(msg) {
			usage := cand.UseLine()
			need := strings.TrimSpace(strings.TrimPrefix(usage, cand.CommandPath()+" "))
			if need == "" {
				need = "[flags]"
			}
			msg = "need " + need
		}

		// --- print hints first ---
		style.Hint("Usage: " + style.UsageLine(cand))
		style.Hint("Help : " + style.Code(cand.CommandPath()+" --help"))
		fmt.Fprintln(os.Stderr) // visual spacer

		// --- then the error last ---
		style.Error(msg)

		os.Exit(1)
	}
}

var acceptsArgsRE = regexp.MustCompile(`^accepts \d+ arg.*received \d+`)

// Recursively silence Cobra’s default usage/error spam for all commands.
func applySilence(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	for _, c := range cmd.Commands() {
		applySilence(c)
	}
}
