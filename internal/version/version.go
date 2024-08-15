package version

import (
	"fmt"
	"runtime"
)

var (
	version   = "1.1.2"
	buildDate = "unknown"
	gitCommit = "unknown"
)

func Version() string   { return version }
func BuildDate() string { return buildDate }
func GitCommit() string { return gitCommit }

func FullVersion() string {
	return fmt.Sprintf(
		"Clido version %s\nBuild date: %s\nGit commit: %s\nGo version: %s\nOS/Arch: %s/%s",
		Version(),
		BuildDate(),
		GitCommit(),
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
