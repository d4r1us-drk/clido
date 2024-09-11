package version

import (
    "fmt"
    "runtime"
)

// Info represents the build information for the application.
// It includes fields for the version, build date, and Git commit.
type Info struct {
    Version   string  // The version of the application
    BuildDate string  // The date when the application was built
    GitCommit string  // The Git commit hash corresponding to the build
}

// Get returns an Info struct containing the current build information.
// By default, it returns "dev" for the version, and "unknown" for the build date and Git commit.
// These values can be replaced at build time using linker flags.
func Get() Info {
    return Info{
        Version:   "dev",
        BuildDate: "unknown",
        GitCommit: "unknown",
    }
}

// FullVersion returns a formatted string with full version details, including the application version, 
// build date, Git commit hash, the Go runtime version, and the operating system and architecture details.
//
// Example output:
// Clido version dev
// Build date: unknown
// Git commit: unknown
// Go version: go1.16.4
// OS/Arch: linux/amd64
func FullVersion() string {
    info := Get()
    return fmt.Sprintf(
        "Clido version %s\nBuild date: %s\nGit commit: %s\nGo version: %s\nOS/Arch: %s/%s",
        info.Version,
        info.BuildDate,
        info.GitCommit,
        runtime.Version(),  // The Go runtime version
        runtime.GOOS,       // The operating system
        runtime.GOARCH,     // The system architecture (e.g., amd64, arm)
    )
}
