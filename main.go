package main

import "github.com/repsejnworb/keysej/cmd"

var (
	// These get overridden by -ldflags at build/release time
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	cmd.SetBuildInfo(version, commit, date)
	cmd.Execute()
}
