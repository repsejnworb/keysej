package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/repsejnworb/keysej/internal/config"
	"github.com/repsejnworb/keysej/internal/shell"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show loaded SSH keys and keysej-generated keys on disk",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		sshDir := config.C.SSHDir // from config

		// 1) Loaded keys from agent (ssh-add -l)
		loadedRaw, _ := shell.AgentList(ctx)
		loaded := parseAgentList(strings.TrimSpace(loadedRaw)) // map[fpr]agentEntry

		// 2) On-disk keysej keys: scan SSHDir for *.pub with ::keysej:: in comment
		onDisk, _ := findKeysejKeys(ctx, sshDir) // map[fpr]diskEntry

		// 3) Render
		fmt.Println(renderLoaded(loaded, onDisk))
		fmt.Println()
		fmt.Println(renderOnDisk(onDisk, loaded, sshDir))
		return nil
	},
}

func init() { rootCmd.AddCommand(listCmd) }

// ----- Types

type agentEntry struct {
	Bits        string
	Fingerprint string // "SHA256:....."
	Comment     string
	Alg         string // e.g. (ED25519) w/o parens
	Raw         string
}

type diskEntry struct {
	Name        string // base filename w/o .pub (e.g. id_ed25519_work)
	PubPath     string
	Fingerprint string // "SHA256:....."
	Alg         string
	Comment     string
}

// ----- Parse helpers

// ssh-add -l line format:
// "256 SHA256:abc... comment (ED25519)"
var agentLineRE = regexp.MustCompile(`^(\d+)\s+(SHA256:[A-Za-z0-9+/=]+)\s+(.*)\s+$begin:math:text$([^)]+)$end:math:text$$`)

// ssh-keygen -lf output format:
// "256 SHA256:abc... comment (ED25519)"
var fpLineRE = regexp.MustCompile(`^\s*\d+\s+(SHA256:[A-Za-z0-9+/=]+)\s+(.*)\s+$begin:math:text$([^)]+)$end:math:text$\s*$`)

func parseAgentList(out string) map[string]agentEntry {
	m := map[string]agentEntry{}
	if out == "" || strings.Contains(out, "The agent has no identities.") {
		return m
	}
	for _, line := range strings.Split(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if ms := agentLineRE.FindStringSubmatch(line); ms != nil {
			bits, fpr, rest, alg := ms[1], ms[2], ms[3], ms[4]
			// rest may contain spaces; it's the comment we want as-is
			m[fpr] = agentEntry{Bits: bits, Fingerprint: fpr, Comment: rest, Alg: alg, Raw: line}
		} else {
			// Fallback: try to at least grab the fpr
			fields := strings.Fields(line)
			if len(fields) >= 2 && strings.HasPrefix(fields[1], "SHA256:") {
				m[fields[1]] = agentEntry{Bits: fields[0], Fingerprint: fields[1], Comment: strings.Join(fields[2:], " "), Alg: ""}
			}
		}
	}
	return m
}

func findKeysejKeys(ctx context.Context, sshDir string) (map[string]diskEntry, error) {
	res := map[string]diskEntry{}
	entries, err := os.ReadDir(sshDir)
	if err != nil {
		return res, nil // quietly ignore if dir unreadable
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if !strings.HasSuffix(name, ".pub") {
			continue
		}
		pubPath := filepath.Join(sshDir, name)
		// read .pub to ensure it's a keysej key (comment contains ::keysej::)
		b, err := os.ReadFile(pubPath)
		if err != nil {
			continue
		}
		line := strings.TrimSpace(string(b))
		if !strings.Contains(line, "::keysej::") {
			continue
		}
		// fingerprint via ssh-keygen -lf
		fpOut, err := shell.Fingerprint(ctx, pubPath)
		if err != nil {
			continue
		}
		fpOut = strings.TrimSpace(fpOut)
		fpr, comment, alg := parseFingerprintLine(fpOut)
		if fpr == "" {
			continue
		}
		res[fpr] = diskEntry{
			Name:        strings.TrimSuffix(name, ".pub"),
			PubPath:     pubPath,
			Fingerprint: fpr,
			Alg:         alg,
			Comment:     comment,
		}
	}
	return res, nil
}

func parseFingerprintLine(line string) (fpr, comment, alg string) {
	if ms := fpLineRE.FindStringSubmatch(line); ms != nil {
		return ms[1], ms[2], ms[3]
	}
	// fallback best-effort
	fields := strings.Fields(line)
	if len(fields) >= 2 && strings.HasPrefix(fields[1], "SHA256:") {
		return fields[1], strings.Join(fields[2:], " "), ""
	}
	return "", "", ""
}

// ----- Rendering

func renderLoaded(loaded map[string]agentEntry, onDisk map[string]diskEntry) string {
	title := lipgloss.NewStyle().Bold(true).Underline(true).Render("Loaded SSH keys:")
	greenBold := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	redBold := lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true)
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	dim := lipgloss.NewStyle().Faint(true)

	if len(loaded) == 0 {
		return title + "\n" + dim.Render("No keys loaded in SSH agent.")
	}

	var b strings.Builder
	b.WriteString(title + "\n")
	for _, v := range loaded {
		_, hasDisk := onDisk[v.Fingerprint]
		// Try to detect if it's keysej-generated (comment contains tag)
		isKeysej := strings.Contains(v.Comment, "::keysej::")

		line := fmt.Sprintf("%s %s %s",
			v.Bits,
			v.Fingerprint,
			v.CommentWithAlg(),
		)

		switch {
		case hasDisk:
			// Loaded and we have the matching keysej file on disk
			b.WriteString(greenBold.Render(line) + "\n")
		case !isKeysej:
			// Loaded, but not keysej
			b.WriteString(redBold.Render(line) + "\n")
		default:
			// Loaded key with no file on disk (forwarded or removed)
			b.WriteString(red.Render(line) + " " + dim.Render("(no local file)") + "\n")
		}
	}
	return b.String()
}

func (a agentEntry) CommentWithAlg() string {
	if a.Alg == "" {
		return a.Comment
	}
	return fmt.Sprintf("%s (%s)", a.Comment, a.Alg)
}

func renderOnDisk(onDisk map[string]diskEntry, loaded map[string]agentEntry, sshDir string) string {
	title := lipgloss.NewStyle().Bold(true).Underline(true).
		Render(fmt.Sprintf("Available keysej keys (under %s):", sshDir))
	greenBold := lipgloss.NewStyle().Foreground(lipgloss.Color("42")).Bold(true)
	red := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	dim := lipgloss.NewStyle().Faint(true)

	if len(onDisk) == 0 {
		return title + "\n" + dim.Render("No keysej-generated keys found on disk.")
	}

	// We'll show "name" with color by loaded state.
	var b strings.Builder
	b.WriteString(title + "\n")
	for fpr, v := range onDisk {
		_, isLoaded := loaded[fpr]
		if isLoaded {
			b.WriteString(greenBold.Render(v.Name) + "\n")
		} else {
			b.WriteString(red.Render(v.Name) + "\n")
		}
	}
	return b.String()
}
