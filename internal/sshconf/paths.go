package sshconf

import (
	"os"
	"path/filepath"
)

func dir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssh", "config.d")
}

func ensureDir() error { return os.MkdirAll(dir(), 0o700) }

func FileForKey(key string) string {
	return filepath.Join(dir(), "keysej."+key+".conf")
}

func AllKeysejFiles() ([]string, error) {
	_ = ensureDir()
	entries, err := os.ReadDir(dir())
	if err != nil {
		return nil, err
	}
	var res []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if len(name) >= 8 && name[:7] == "keysej." && filepath.Ext(name) == ".conf" {
			res = append(res, filepath.Join(dir(), name))
		}
	}
	return res, nil
}
