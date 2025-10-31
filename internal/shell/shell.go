package shell

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

var DryRun bool

func Keygen(ctx context.Context, keyPath, comment, pass string) error {
	args := []string{"-t", "ed25519", "-o", "-a", "100", "-C", comment, "-f", keyPath, "-N", pass}
	cmd := exec.CommandContext(ctx, "ssh-keygen", args...)
	cmd.Stdin = nil
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	return cmd.Run()
}

func AgentAdd(ctx context.Context, keyPath, ttl, pass string, dryRun bool) error {
	args := []string{}
	if ttl != "0" {
		args = append(args, "-t", ttl)
	}
	args = append(args, keyPath)

	if dryRun {
		fmt.Println("🔬 Dry-run: would run:", "ssh-add", args)
		return nil
	}

	if runtime.GOOS == "darwin" {
		// seed Keychain with the passphrase so --apple-use-keychain can consume it non-interactively
		label := fmt.Sprintf("SSH: %s", keyPath)
		sec := exec.CommandContext(ctx, "security", "add-generic-password", "-a", os.Getenv("USER"), "-l", label, "-s", label, "-w", pass, "-T", "/usr/bin/ssh", "-T", "/usr/bin/ssh-add")
		_ = sec.Run() // ignore if exists; ssh-add will still succeed

		cmd := exec.CommandContext(ctx, "ssh-add", append([]string{"--apple-use-keychain"}, args...)...)
		cmd.Stdin = nil
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Linux/other: use askpass to avoid a second prompt when we already have the passphrase
	// We'll call the current binary in askpass mode, passing the secret via env **in-memory** only.
	exe, _ := os.Executable()
	cmd := exec.CommandContext(ctx, "ssh-add", args...)
	cmd.Env = append(os.Environ(),
		"SSH_ASKPASS="+exe,
		"SSH_ASKPASS_REQUIRE=force",
		"KEYSEJ_MODE=askpass",
		"KEYSEJ_SECRET="+pass,
	)
	cmd.Stdin = bytes.NewReader(nil) // close stdin to force askpass
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AgentRemove(ctx context.Context, keyPath string) error {
	cmd := exec.CommandContext(ctx, "ssh-add", "-d", keyPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func AgentList(ctx context.Context) (string, error) {
	out, err := exec.CommandContext(ctx, "ssh-add", "-l").CombinedOutput()
	return string(out), err
}

func Fingerprint(ctx context.Context, pubPath string) (string, error) {
	out, err := exec.CommandContext(ctx, "ssh-keygen", "-lf", pubPath).CombinedOutput()
	return string(out), err
}

func InstallPubkey(ctx context.Context, pubPath, target string) error {
	// If ssh-copy-id exists, use it
	if _, err := exec.LookPath("ssh-copy-id"); err == nil {
		cmd := exec.CommandContext(ctx, "ssh-copy-id", "-i", pubPath, target)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	// Fallback: print a helpful message with a safe one-liner
	fmt.Println("ssh-copy-id not found. Run this:")
	fmt.Printf("cat %q | ssh %s 'mkdir -p ~/.ssh && chmod 700 ~/.ssh && \\\n  touch ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys && \\\n  grep -qxF \"$(cat)\" ~/.ssh/authorized_keys || echo \"$(cat)\" >> ~/.ssh/authorized_keys'\n", pubPath, target)
	return nil
}

// macOS-only: best-effort removal of keychain entry (non-fatal on error)
func MacDeleteKeychainEntry(keyPath string) error {
	if runtime.GOOS != "darwin" {
		return nil
	}
	label := fmt.Sprintf("SSH: %s", keyPath)
	cmd := exec.Command("security", "delete-generic-password", "-l", label)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
