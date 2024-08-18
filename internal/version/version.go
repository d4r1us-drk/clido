package version

import (
	"fmt"
	"runtime"
)

type Info struct {
	Version   string
	BuildDate string
	GitCommit string
}

func Get() Info {
	return Info{
		Version:   "dev",
		BuildDate: "unknown",
		GitCommit: "unknown",
	}
}

func FullVersion() string {
	info := Get()
	return fmt.Sprintf(
		"Clido version %s\nBuild date: %s\nGit commit: %s\nGo version: %s\nOS/Arch: %s/%s",
		info.Version,
		info.BuildDate,
		info.GitCommit,
		runtime.Version(),
		runtime.GOOS,
		runtime.GOARCH,
	)
}
