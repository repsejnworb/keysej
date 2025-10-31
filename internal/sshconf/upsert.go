package sshconf

import (
	"errors"
	"os"
	"strings"
)

// UpsertBlock writes/updates a block identified by its BEGIN marker.
// Returns (changed, err)
func UpsertBlock(path, marker, frag string) (bool, error) {
	_ = ensureDir()
	// if file doesn't exist, write fresh
	old, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return false, err
	}
	if errors.Is(err, os.ErrNotExist) {
		return writeFile(path, frag)
	}
	// find existing block by marker
	s := string(old)
	start := strings.Index(s, marker+"\n")
	if start == -1 {
		// append
		return writeFile(path, strings.TrimRight(s, "\n")+"\n\n"+frag)
	}
	// replace from marker to END
	end := strings.Index(s[start:], "\n# END keysej")
	if end == -1 {
		// malformed; append new safely
		return writeFile(path, strings.TrimRight(s, "\n")+"\n\n"+frag)
	}
	endIdx := start + end
	// include END line (+len)
	endLineIdx := strings.Index(s[endIdx:], "\n")
	if endLineIdx == -1 {
		endLineIdx = len(s) - endIdx
	}
	replaced := s[:start] + frag + s[endIdx+endLineIdx:]
	if replaced == s {
		return false, nil
	}
	return writeFile(path, replaced)
}

func writeFile(path, content string) (bool, error) {
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		return false, err
	}
	return true, nil
}
