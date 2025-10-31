package sshconf

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	beginLineRE = regexp.MustCompile(`(?m)^# BEGIN keysej:(host|cidr):(.+)$`)
	hostLineRE  = regexp.MustCompile(`(?m)^Host\s+(.+)$`)
	matchRE     = regexp.MustCompile(`(?m)^Match\s+address\s+(.+)$`)
	identRE     = regexp.MustCompile(`(?m)^\s*IdentityFile\s+(.+)$`)
)

func ValidateAll() []string {
	var issues []string
	files, err := AllKeysejFiles()
	if err != nil {
		return []string{fmt.Sprintf("cannot list files: %v", err)}
	}

	seenMarker := map[string]string{} // marker -> file

	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: read error: %v", f, err))
			continue
		}
		s := string(b)

		// check duplicate markers
		for _, m := range beginLineRE.FindAllStringSubmatch(s, -1) {
			marker := m[0]
			if other, ok := seenMarker[marker]; ok {
				issues = append(issues, fmt.Sprintf("%s: duplicate rule also in %s: %s", f, other, marker))
			} else {
				seenMarker[marker] = f
			}
		}

		// find each block
		blocks := strings.Split(s, "# BEGIN keysej:")
		for _, blk := range blocks {
			if strings.TrimSpace(blk) == "" {
				continue
			}
			blk = "# BEGIN keysej:" + blk

			// host vs cidr correctness
			if m := hostLineRE.FindStringSubmatch(blk); m != nil {
				// Host supports globs; no deep validation needed
			} else if m := matchRE.FindStringSubmatch(blk); m != nil {
				addr := strings.TrimSpace(m[1])
				if _, _, err := net.ParseCIDR(addr); err != nil {
					issues = append(issues, fmt.Sprintf("%s: invalid CIDR in Match address: %q", f, addr))
				}
			} else {
				issues = append(issues, fmt.Sprintf("%s: block missing Host/Match header", f))
			}

			// identity exists + keysej tag present
			if m := identRE.FindStringSubmatch(blk); m != nil {
				path := expandTilde(m[1])
				if _, err := os.Stat(path); err != nil {
					issues = append(issues, fmt.Sprintf("%s: IdentityFile not found: %s", f, path))
				}
				pub := path + ".pub"
				if b, err := os.ReadFile(pub); err == nil {
					if !strings.Contains(string(b), "::keysej::") {
						issues = append(issues, fmt.Sprintf("%s: public key missing ::keysej:: tag: %s", f, pub))
					}
				} else {
					issues = append(issues, fmt.Sprintf("%s: public key not found: %s", f, pub))
				}
			} else {
				issues = append(issues, fmt.Sprintf("%s: missing IdentityFile in block", f))
			}
		}
	}
	return issues
}

func expandTilde(p string) string {
	if strings.HasPrefix(p, "~") {
		home, _ := os.UserHomeDir()
		if p == "~" {
			return home
		}
		if strings.HasPrefix(p, "~/") {
			return filepath.Join(home, p[2:])
		}
	}
	return p
}
