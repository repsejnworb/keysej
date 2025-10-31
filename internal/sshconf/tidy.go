package sshconf

import (
	"os"
	"regexp"
	"sort"
	"strings"
)

var blockRE = regexp.MustCompile(`(?ms)^# BEGIN keysej:.*?# END keysej\n*`)

func TidyAll() (int, error) {
	files, err := AllKeysejFiles()
	if err != nil {
		return 0, err
	}
	changed := 0
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		s := string(b)

		blocks := blockRE.FindAllString(s, -1)
		if len(blocks) == 0 {
			continue
		}

		// normalize whitespace inside each block
		for i := range blocks {
			blk := strings.TrimRight(blocks[i], "\n")
			blk = strings.ReplaceAll(blk, "\r\n", "\n")
			blk = blk + "\n\n" // two newlines between blocks
			blocks[i] = blk
		}

		// sort by BEGIN line alphabetically for stability
		sort.SliceStable(blocks, func(i, j int) bool {
			// compare first line of each
			li := blocks[i][:strings.Index(blocks[i], "\n")]
			lj := blocks[j][:strings.Index(blocks[j], "\n")]
			return li < lj
		})

		out := strings.Join(blocks, "")
		if !strings.HasSuffix(out, "\n") {
			out += "\n"
		}

		if out != s {
			if err := os.WriteFile(f, []byte(out), 0o600); err == nil {
				changed++
			}
		}
	}
	return changed, nil
}
