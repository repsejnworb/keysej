package osutil

import (
	"errors"
	"os"
	"regexp"
	"runtime"
)

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
