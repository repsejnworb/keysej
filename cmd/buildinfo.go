// cmd/buildinfo.go
package cmd

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

func SetBuildInfo(v, c, d string) {
	buildVersion, buildCommit, buildDate = v, c, d
	rootCmd.Version = buildVersion // enables `--version`
}
