package installer

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"

	"github.com/heidene/flashdoc/internal/pkgmanager"
	"github.com/heidene/flashdoc/internal/progress"
)

// Installer handles dependency installation
type Installer struct {
	workspacePath  string
	packageManager pkgmanager.PackageManager
}

// New creates a new installer
func New(workspacePath string, pm pkgmanager.PackageManager) *Installer {
	return &Installer{
		workspacePath:  workspacePath,
		packageManager: pm,
	}
}

// Install runs the package manager install command
func (i *Installer) Install() error {
	// Pick a random witty message
	message := progress.InstallMessages[rand.Intn(len(progress.InstallMessages))]

	sp := progress.New(message)
	sp.Start()

	// Get install command
	cmdArgs := i.packageManager.InstallCommand()
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = i.workspacePath

	// Discard verbose output
	cmd.Stdout = progress.DiscardWriter()
	cmd.Stderr = progress.DiscardWriter()

	// Run the install command
	if err := cmd.Run(); err != nil {
		sp.StopWithError("Dependencies installation failed")
		return fmt.Errorf("dependency installation failed: %w", err)
	}

	sp.Stop("Dependencies installed")

	return nil
}

// InstallShared installs dependencies to the shared directory
func InstallShared(sharedDir string, pm pkgmanager.PackageManager) error {
	// Pick a random witty message
	message := progress.InstallMessages[rand.Intn(len(progress.InstallMessages))]

	sp := progress.New(message)
	sp.Start()

	// Get install command
	cmdArgs := pm.InstallCommand()
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = sharedDir

	// Discard verbose output
	cmd.Stdout = progress.DiscardWriter()
	cmd.Stderr = progress.DiscardWriter()

	// Run the install command
	if err := cmd.Run(); err != nil {
		sp.StopWithError("Dependencies installation failed")
		return fmt.Errorf("shared dependency installation failed: %w", err)
	}

	sp.Stop("Dependencies installed")

	return nil
}

// IsInstalled checks if node_modules exists in the given directory
func IsInstalled(dir string) bool {
	nodeModules := fmt.Sprintf("%s/node_modules", dir)
	info, err := os.Stat(nodeModules)
	if err != nil {
		return false
	}
	return info.IsDir()
}
