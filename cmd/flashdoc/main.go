package main

import (
	"fmt"
	"os"
	"time"

	"github.com/heidene/flashdoc/internal/browser"
	"github.com/heidene/flashdoc/internal/builder"
	"github.com/heidene/flashdoc/internal/cleanup"
	"github.com/heidene/flashdoc/internal/cli"
	"github.com/heidene/flashdoc/internal/exporter"
	"github.com/heidene/flashdoc/internal/installer"
	"github.com/heidene/flashdoc/internal/pkgmanager"
	"github.com/heidene/flashdoc/internal/processor"
	"github.com/heidene/flashdoc/internal/shared"
	"github.com/heidene/flashdoc/internal/signal"
	"github.com/heidene/flashdoc/internal/staticserver"
	"github.com/heidene/flashdoc/internal/template"
	"github.com/heidene/flashdoc/internal/workspace"
)

func main() {
	// Parse CLI arguments
	cfg, helpOrVersion, err := cli.Parse(os.Args[1:])
	if err != nil {
		// Error occurred during parsing
		os.Exit(1)
	}
	if helpOrVersion {
		// --help or --version was used, exit successfully
		os.Exit(0)
	}

	// Create shared project manager
	sharedMgr, err := shared.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create shared manager: %v\n", err)
		os.Exit(1)
	}

	// Ensure directories exist
	if err := sharedMgr.EnsureDirectories(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create directories: %v\n", err)
		os.Exit(1)
	}

	// Cleanup old runs (older than 24 hours)
	if err := sharedMgr.CleanupOldRuns(24 * time.Hour); err != nil {
		// Log warning but don't fail
		fmt.Fprintf(os.Stderr, "Warning: failed to cleanup old runs: %v\n", err)
	}

	// Get package hash for cache invalidation
	packageHash, err := template.GetEmbeddedPackageHash()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to get package hash: %v\n", err)
		os.Exit(1)
	}

	// Check if shared project is current
	isCurrent, err := sharedMgr.IsSharedProjectCurrent(packageHash)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to check shared project: %v\n", err)
		os.Exit(1)
	}

	// Install to shared directory if needed
	if !isCurrent || cfg.ForceReinstall {
		// Acquire lock to prevent concurrent installs
		if err := sharedMgr.AcquireLock(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer func() { _ = sharedMgr.ReleaseLock() }()

		// Extract template to shared directory
		if err := template.ExtractToShared(sharedMgr.GetSharedDir()); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to extract template: %v\n", err)
			os.Exit(1)
		}

		// Detect package manager
		pm, err := pkgmanager.Detect()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Install dependencies to shared directory
		if err := installer.InstallShared(sharedMgr.GetSharedDir(), pm); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Save version hash
		if err := sharedMgr.SaveVersion(packageHash); err != nil {
			fmt.Fprintf(os.Stderr, "Error: failed to save version: %v\n", err)
			os.Exit(1)
		}
	}

	// Generate unique run ID
	runID := sharedMgr.GenerateRunID()
	runDir := sharedMgr.GetRunDir(runID)

	// Create workspace with symlinks to shared project
	ws, err := workspace.New(runDir, sharedMgr.GetSharedDir())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to create workspace: %v\n", err)
		os.Exit(1)
	}

	// Setup workspace structure (creates symlinks)
	if err := ws.Setup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to setup workspace: %v\n", err)
		_ = ws.Cleanup()
		os.Exit(1)
	}

	// Setup cleanup manager
	cleanupMgr := cleanup.New(ws)
	defer func() { _ = cleanupMgr.Cleanup() }()

	// Setup signal handling
	sigHandler := signal.New(cleanupMgr.Cleanup)
	sigHandler.Setup()

	// Log workspace path
	fmt.Printf("📦 Workspace: %s\n", ws.Path)

	// Extract config files only (not package.json, which is symlinked)
	if err := template.ExtractConfigOnly(ws.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to extract config: %v\n", err)
		os.Exit(1)
	}

	// Generate config with title
	siteTitle := cfg.Title
	if siteTitle == "" {
		siteTitle = template.GenerateTitle(cfg.SourceDir)
	}

	if err := template.GenerateConfig(ws.Path, siteTitle); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to generate config: %v\n", err)
		os.Exit(1)
	}

	// Process markdown files
	targetDir := ws.GetDocsDir()
	proc := processor.New(cfg.SourceDir, targetDir)

	if err := proc.Process(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Detect package manager for build command
	pm, err := pkgmanager.Detect()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Build static site
	bldr := builder.NewBuilder(ws.Path, pm.String(), os.Stdout)
	if err := bldr.Build(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if export mode is enabled
	if cfg.ExportPath != "" {
		// Export mode - copy built files and exit
		distPath := ws.GetDistDir()
		exp := exporter.New(distPath, cfg.ExportPath, os.Stdout)

		if err := exp.Export(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Exit successfully after export
		os.Exit(0)
	}

	// Normal mode - start dev server
	distPath := ws.GetDistDir()
	srv := staticserver.NewServer(distPath, cfg.Port, os.Stdout)

	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Register server for cleanup
	cleanupMgr.RegisterServer(srv)

	// Open browser unless --no-open flag is set
	serverURL := srv.GetURL()
	if !cfg.NoOpen {
		if err := browser.Open(serverURL); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to open browser: %v\n", err)
			fmt.Printf("Please open %s manually\n", serverURL)
		}
	} else {
		fmt.Println("(browser not opened due to --no-open flag)")
	}

	// Wait for signals
	fmt.Println("\nPress Ctrl+C to exit")
	sigHandler.Wait()
}
