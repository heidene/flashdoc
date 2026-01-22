package cli

// Config holds the CLI configuration parsed from flags and arguments
type Config struct {
	SourceDir      string
	Title          string
	Port           int
	NoOpen         bool
	ForceReinstall bool
}

// Version is the current version of stardoc
const Version = "0.1.0"
