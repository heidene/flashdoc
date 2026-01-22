package builder

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/nicovandenhove/stardoc/internal/progress"
)

// Builder handles static site building with Astro
type Builder struct {
	workspacePath string
	packageMgr    string
	output        io.Writer
}

// NewBuilder creates a new builder instance
func NewBuilder(workspacePath, packageMgr string, output io.Writer) *Builder {
	if output == nil {
		output = os.Stdout
	}

	return &Builder{
		workspacePath: workspacePath,
		packageMgr:    packageMgr,
		output:        output,
	}
}

// Build executes the astro build command
func (b *Builder) Build() error {
	// Pick a random witty message
	message := progress.BuildMessages[rand.Intn(len(progress.BuildMessages))]

	sp := progress.New(message)
	sp.Start()

	// Determine build command based on package manager
	var cmd *exec.Cmd
	switch b.packageMgr {
	case "npm":
		cmd = exec.Command("npm", "run", "build")
	case "pnpm":
		cmd = exec.Command("pnpm", "run", "build")
	case "bun":
		cmd = exec.Command("bun", "run", "build")
	default:
		sp.StopWithError("Unsupported package manager")
		return fmt.Errorf("unsupported package manager: %s", b.packageMgr)
	}

	cmd.Dir = b.workspacePath
	// Discard verbose build output
	cmd.Stdout = progress.DiscardWriter()
	cmd.Stderr = progress.DiscardWriter()

	if err := cmd.Run(); err != nil {
		sp.StopWithError("Build failed")
		return fmt.Errorf("build failed: %w", err)
	}

	// Verify dist directory was created
	distPath := filepath.Join(b.workspacePath, "dist")
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		sp.StopWithError("Build completed but dist directory not found")
		return fmt.Errorf("build completed but dist directory not found")
	}

	sp.Stop("Documentation site ready")

	return nil
}

// GetDistPath returns the path to the build output directory
func (b *Builder) GetDistPath() string {
	return filepath.Join(b.workspacePath, "dist")
}
