package version

var (
	Version = "0.0.0"
	Commit  = ""
)

func String() string {
	if Commit == "" {
		return Version
	}
	return Version + "+" + Commit
}
