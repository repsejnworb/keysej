package osutil

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

var SSHDir string

func InitSSHDir(dryRun bool) error {
	if dryRun {
		tmp, err := os.MkdirTemp("", "keysej_ssh_*")
		if err != nil {
			return fmt.Errorf("create temp ssh dir: %w", err)
		}
		SSHDir = tmp
		fmt.Println("🔬 Dry-run: using temporary SSH dir:", SSHDir)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("home dir: %w", err)
		}
		SSHDir = filepath.Join(home, ".ssh")
	}
	// mkdir -p + perms
	if err := os.MkdirAll(SSHDir, 0o700); err != nil {
		return fmt.Errorf("mkdir %s: %w", SSHDir, err)
	}
	// ensure perms in case dir existed
	_ = os.Chmod(SSHDir, 0o700)
	return nil
}

var reName = regexp.MustCompile(`^[a-z0-9._-]+$`)

func ValidateName(n string) error {
	if !reName.MatchString(n) {
		return errors.New("name must match [a-z0-9._-]+, e.g. 'work' or 'homelab'")
	}
	return nil
}

func IsDarwin() bool { return runtime.GOOS == "darwin" }

func Hostname() string {
	h, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return h
}
