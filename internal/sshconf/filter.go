package sshconf

import "strings"

// FilterByHostOrAll returns only blocks matching a host pattern (substring match)
// If host == "", returns the whole content.
func FilterByHostOrAll(s, host string) string {
	if host == "" {
		return s
	}
	var out []string
	blocks := strings.Split(s, "# BEGIN keysej:")
	for _, b := range blocks {
		b = strings.TrimSpace(b)
		if b == "" {
			continue
		}
		blk := "# BEGIN keysej:" + b
		// crude substring match on the header line
		firstLineEnd := strings.Index(blk, "\n")
		header := blk
		if firstLineEnd > 0 {
			header = blk[:firstLineEnd]
		}
		if strings.Contains(header, host) {
			out = append(out, blk)
		}
	}
	return strings.Join(out, "\n\n")
}
