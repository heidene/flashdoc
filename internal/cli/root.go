package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// exportValue implements pflag.Value for the --export flag
type exportValue struct {
	path *string
}

func (e *exportValue) String() string {
	if e.path == nil {
		return ""
	}
	return *e.path
}

func (e *exportValue) Set(val string) error {
	*e.path = val
	return nil
}

func (e *exportValue) Type() string {
	return "string"
}

var _ pflag.Value = (*exportValue)(nil)

var (
	title          string
	port           int
	noOpen         bool
	forceReinstall bool
	exportPath     string
)

// customArgsValidator validates arguments allowing for --export path
func customArgsValidator(cmd *cobra.Command, args []string) error {
	// We need exactly 1 positional arg (the directory)
	// But if --export is used with a space, there might be 2 args
	if len(args) < 1 {
		return cobra.MinimumNArgs(1)(cmd, args)
	}
	if len(args) > 2 {
		return cobra.MaximumNArgs(2)(cmd, args)
	}
	// If we have 2 args, the second one should be the export path
	// This will be handled in Parse()
	return nil
}

// NewRootCommand creates the root cobra command
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "stardoc <directory>",
		Short:   "Generate and serve a Starlight documentation site from markdown files",
		Version: FullVersion(),
		Args:    customArgsValidator,
		RunE:    runStardoc,
		SilenceUsage: true,
	}

	rootCmd.Flags().StringVar(&title, "title", "", "Title for the documentation site")
	rootCmd.Flags().IntVar(&port, "port", 4321, "Port for the dev server (1024-65535)")
	rootCmd.Flags().BoolVar(&noOpen, "no-open", false, "Don't automatically open the browser")
	rootCmd.Flags().BoolVar(&forceReinstall, "force-reinstall", false, "Force reinstall of dependencies even if cached")

	// Export flag with optional value
	exportFlag := rootCmd.Flags().VarPF(
		&exportValue{path: &exportPath},
		"export",
		"",
		"Export static build to directory (default: ./export-doc)",
	)
	exportFlag.NoOptDefVal = "./export-doc"

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
	nonFlagArgs := []string{}
	for _, arg := range args {
		if arg != "" && arg[0] != '-' {
			hasSourceDir = true
			nonFlagArgs = append(nonFlagArgs, arg)
		}
	}

	// If no source directory provided, it means --help or --version was used
	if !hasSourceDir {
		return nil, true, nil
	}

	// Extract the source directory from args if available
	sourceDir := ""
	if len(nonFlagArgs) > 0 {
		sourceDir = nonFlagArgs[0]
	}

	// Handle --export with space-separated path
	// If we have 2 non-flag args and exportPath is empty or is the default,
	// the second arg is the export path
	finalExportPath := exportPath
	if len(nonFlagArgs) == 2 {
		// Check if --export flag was present in args
		exportFlagPresent := false
		for _, arg := range args {
			if arg == "--export" {
				exportFlagPresent = true
				break
			}
		}
		if exportFlagPresent {
			// Second non-flag arg is the export path
			finalExportPath = nonFlagArgs[1]
		}
	}

	return &Config{
		SourceDir:      sourceDir,
		Title:          title,
		Port:           port,
		NoOpen:         noOpen,
		ForceReinstall: forceReinstall,
		ExportPath:     finalExportPath,
	}, false, nil
}
