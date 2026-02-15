package cli

import (
	"fmt"
	"runtime/debug"
)

// Config holds the CLI configuration parsed from flags and arguments
type Config struct {
	SourceDir      string
	Title          string
	Port           int
	NoOpen         bool
	ForceReinstall bool
	ExportPath     string // Path to export static build, empty means no export
}

// Version variables - injected at build time via ldflags
var (
	Version = "dev"       // Semantic version (e.g., "0.2.0")
	Commit  = "unknown"   // Git commit hash
	Date    = "unknown"   // Build date
)

// FullVersion returns a detailed version string including commit and build date
func FullVersion() string {
	// Try to get version from build info if available (for go install)
	if Version == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			Version = info.Main.Version
		}
	}

	if Commit != "unknown" && Date != "unknown" {
		return fmt.Sprintf("%s (commit: %s, built: %s)", Version, Commit, Date)
	}
	return Version
}
