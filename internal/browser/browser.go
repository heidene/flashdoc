package browser

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

// Open opens the given URL in the default browser
func Open(url string) error {
	fmt.Printf("Opening browser at %s...\n", url)

	// Slight delay to ensure server is fully ready
	time.Sleep(1 * time.Second)

	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	// Run with timeout to avoid hanging
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("failed to open browser: %w", err)
		}
		return nil
	case <-time.After(5 * time.Second):
		// Timeout - but don't fail, just warn
		fmt.Println("Warning: browser open command timed out")
		return nil
	}
}
