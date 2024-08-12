package version

import (
	"fmt"
	"runtime"
)

var (
	Version = "1.1.0"

	BuildDate = "unknown"

	GitCommit = "unknown"
)

func FullVersion() string {
	return fmt.Sprintf(
		"Clido version %s\nBuild date: %s\nGit commit: %s\nGo version: %s\nOS/Arch: %s/%s",
		Version,
		BuildDate,
		GitCommit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
