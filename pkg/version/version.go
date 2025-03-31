package version

import "fmt"

// The following variables will be overwritten at build time using -ldflags.
var (
	Version   = "unknown"
	Commit    = "unknown"
	BuildDate = "unknown"
)

func GetVersion() string {
	return fmt.Sprintf("Version: %s\nCommit: %s\nBuildDate: %s", Version, Commit, BuildDate)
}
