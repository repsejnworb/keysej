package sshconf

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type pendingDelete struct {
	path string
	data string
}

var pd *pendingDelete

func BackupFile(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	tmp := filepath.Join(os.TempDir(), filepath.Base(path)+".bak")
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return "", err
	}
	return tmp, nil
}

func DeleteBlock(path, marker string) (backupPath string, changed bool, err error) {
	old, err := os.ReadFile(path)
	if err != nil {
		return "", false, err
	}
	s := string(old)

	start := strings.Index(s, marker+"\n")
	if start == -1 {
		return "", false, nil
	}
	end := strings.Index(s[start:], "\n# END keysej")
	if end == -1 {
		return "", false, errors.New("malformed block (no END)")
	}
	endIdx := start + end
	endLineIdx := strings.Index(s[endIdx:], "\n")
	if endLineIdx == -1 {
		endLineIdx = len(s) - endIdx
	}

	backupPath, berr := BackupFile(path)
	if berr != nil {
		return "", false, berr
	}

	newS := strings.TrimSpace(s[:start] + s[endIdx+endLineIdx:])
	if newS != "" {
		newS += "\n"
	} // ensure trailing newline

	pd = &pendingDelete{path: path, data: newS}
	return backupPath, true, nil
}

func CommitDelete(path string) error {
	if pd == nil || pd.path != path {
		return errors.New("no pending delete")
	}
	defer func() { pd = nil }()
	return os.WriteFile(path, []byte(pd.data), 0o600)
}
