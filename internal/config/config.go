package config

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	SSHDir string
}

var C Config

// Init sets SSHDir from flag > env > default $HOME/.ssh
func Init(sshDirFlag string) error {
	if sshDirFlag != "" {
		C.SSHDir = sshDirFlag
		return nil
	}
	if v := os.Getenv("KEYSEJ_SSH_DIR"); v != "" {
		C.SSHDir = v
		return nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return errors.New("cannot determine home directory")
	}
	C.SSHDir = filepath.Join(home, ".ssh")
	return nil
}

// Pretty returns "~" based path for display.
func Pretty(p string) string {
	home, err := os.UserHomeDir()
	if err != nil || p == "" {
		return p
	}
	if len(p) >= len(home) && p[:len(home)] == home {
		rest := p[len(home):]
		if len(rest) == 0 {
			return "~"
		}
		if rest[0] == os.PathSeparator {
			return filepath.Join("~", rest[1:])
		}
	}
	return p
}
