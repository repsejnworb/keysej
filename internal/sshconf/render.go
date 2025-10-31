package sshconf

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

func isCIDR(s string) bool {
	// very lax: contains a slash
	return strings.Contains(s, "/")
}

// markerKey normalizes a unique marker line for a rule
func markerKey(kind, pattern string) string {
	return fmt.Sprintf("# BEGIN keysej:%s:%s", kind, pattern)
}

// RenderBlock returns fragment + its unique marker (BEGIN line)
func RenderBlock(key, pattern, user string, forward bool) (string, string) {
	kind := "host"
	if isCIDR(pattern) {
		kind = "cidr"
	}
	begin := markerKey(kind, pattern)
	end := "# END keysej"

	var body string
	if kind == "host" {
		body = fmt.Sprintf("Host %s\n", pattern)
	} else {
		body = fmt.Sprintf("Match address %s\n", pattern)
	}
	if user != "" {
		body += fmt.Sprintf("  User %s\n", user)
	}
	if forward {
		body += "  ForwardAgent yes\n"
	}
	// always specify explicit identity for the key
	identity := fmt.Sprintf("~/.ssh/%s", filepath.Base("id_ed25519_"+key))
	body += fmt.Sprintf("  IdentityFile %s\n", identity)

	frag := fmt.Sprintf("%s\n%s%s\n\n%s\n", begin, body, "", end)
	return frag, begin
}

var beginRE = regexp.MustCompile(`(?m)^# BEGIN keysej:(host|cidr):(.+)$`)
