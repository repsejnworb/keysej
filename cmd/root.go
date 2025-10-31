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
	// Friendlier flag-parse errors (unknown/missing flags)
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		style.Hint("Usage: " + style.UsageLine(cmd))
		style.Hint("Help : " + style.Code(cmd.CommandPath()+" --help"))
		fmt.Fprintln(os.Stderr)
		style.Error(err.Error())
		return err
	})
}

func Execute() {
	applySilence(rootCmd)

	// Resolve leaf command AND the args meant for it
	cand, candArgs, _ := rootCmd.Find(os.Args[1:])
	if cand == nil {
		cand = rootCmd
	}

	if err := rootCmd.Execute(); err != nil {
		msg := err.Error()

		// If this looks like a Cobra positional-arg error, rewrite it nicely
		if isCobraArgError(msg) {
			msg = prettyArgError(cand, candArgs)
		}

		// Hints first
		style.Hint("Usage: " + style.UsageLine(cand))
		style.Hint("Help : " + style.Code(cand.CommandPath()+" --help"))
		fmt.Fprintln(os.Stderr) // spacer

		// Error last
		style.Error(msg)
		os.Exit(1)
	}
}

// -------- helpers --------

func isCobraArgError(s string) bool {
	ls := strings.ToLower(s)
	return strings.Contains(ls, "accepts ") ||
		strings.Contains(ls, "requires at least") ||
		strings.Contains(ls, "accepts at most")
}

var argTokenRE = regexp.MustCompile(`<[^>]+>`) // matches <key>, <host-or-cidr>, etc.

// prettyArgError builds:
//   - "Missing required arguments: <key> <host-or-cidr>"
//   - "Missing required arguments: <host-or-cidr>"
//   - "Too many arguments: got [extra], expected <key> <host-or-cidr>"
func prettyArgError(cmd *cobra.Command, candArgs []string) string {
	required := argTokenRE.FindAllString(cmd.UseLine(), -1) // e.g. ["<key>", "<host-or-cidr>"]

	// take positional args only (stop at first flag)
	var pos []string
	for _, a := range candArgs {
		if strings.HasPrefix(a, "-") {
			break
		}
		pos = append(pos, a)
	}

	switch {
	case len(pos) < len(required):
		missing := required[len(pos):]
		return "Missing required arguments: " + strings.Join(missing, " ")
	case len(pos) > len(required):
		extra := pos[len(required):]
		return "Too many arguments: got [" + strings.Join(extra, " ") + "], expected " + strings.Join(required, " ")
	default:
		// Same count but still errored? Fallback to concise "need …"
		need := strings.TrimSpace(strings.TrimPrefix(cmd.UseLine(), cmd.CommandPath()+" "))
		if need == "" {
			need = "[flags]"
		}
		return "need " + need
	}
}

// Recursively silence Cobra’s default usage/error spam for all commands.
func applySilence(cmd *cobra.Command) {
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	for _, c := range cmd.Commands() {
		applySilence(c)
	}
}
