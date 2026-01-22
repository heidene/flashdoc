package cli

import (
	"github.com/spf13/cobra"
)

var (
	title          string
	port           int
	noOpen         bool
	forceReinstall bool
)

// NewRootCommand creates the root cobra command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "stardoc <directory>",
		Short:   "Generate and serve a Starlight documentation site from markdown files",
		Version: Version,
		Args:    cobra.ExactArgs(1),
		RunE:    runStardoc,
		SilenceUsage: true,
	}

	rootCmd.Flags().StringVar(&title, "title", "", "Title for the documentation site")
	rootCmd.Flags().IntVar(&port, "port", 4321, "Port for the dev server (1024-65535)")
	rootCmd.Flags().BoolVar(&noOpen, "no-open", false, "Don't automatically open the browser")
	rootCmd.Flags().BoolVar(&forceReinstall, "force-reinstall", false, "Force reinstall of dependencies even if cached")

	return rootCmd
}

func runStardoc(cmd *cobra.Command, args []string) error {
	sourceDir := args[0]

	// Validate the source directory
	if err := ValidatePath(sourceDir); err != nil {
		return err
	}

	// Validate the port
	if err := ValidatePort(port); err != nil {
		return err
	}

	// Store the configuration
	config := &Config{
		SourceDir:      sourceDir,
		Title:          title,
		Port:           port,
		NoOpen:         noOpen,
		ForceReinstall: forceReinstall,
	}

	// For now, just store it - actual execution will be wired up in main.go
	_ = config

	return nil
}

// Parse parses the command line arguments and returns the configuration
func Parse(args []string) (*Config, bool, error) {
	rootCmd := NewRootCommand()
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		return nil, false, err
	}

	// Check if --help or --version was used (indicated by no args or just flags)
	hasSourceDir := false
	for _, arg := range args {
		if arg != "" && arg[0] != '-' {
			hasSourceDir = true
			break
		}
	}

	// If no source directory provided, it means --help or --version was used
	if !hasSourceDir {
		return nil, true, nil
	}

	// Extract the source directory from args if available
	sourceDir := ""
	if len(args) > 0 && args[0] != "" && args[0][0] != '-' {
		sourceDir = args[0]
	}

	return &Config{
		SourceDir:      sourceDir,
		Title:          title,
		Port:           port,
		NoOpen:         noOpen,
		ForceReinstall: forceReinstall,
	}, false, nil
}
