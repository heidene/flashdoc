package server

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/nicovandenhove/flashdoc/internal/pkgmanager"
)

// Server manages the Astro dev server
type Server struct {
	workspacePath  string
	packageManager pkgmanager.PackageManager
	port           int
	cmd            *exec.Cmd
	serverURL      string
	ready          chan bool
}

// New creates a new server manager
func New(workspacePath string, pm pkgmanager.PackageManager, port int) *Server {
	return &Server{
		workspacePath:  workspacePath,
		packageManager: pm,
		port:           port,
		ready:          make(chan bool, 1),
	}
}

// Start starts the Astro dev server
func (s *Server) Start() error {
	fmt.Println("Dev server starting...")

	// Build command based on package manager
	var cmdArgs []string
	switch s.packageManager {
	case pkgmanager.Pnpm:
		cmdArgs = []string{"pnpm", "run", "dev", "--", "--port", fmt.Sprintf("%d", s.port)}
	case pkgmanager.Bun:
		cmdArgs = []string{"bun", "run", "dev", "--port", fmt.Sprintf("%d", s.port)}
	case pkgmanager.Npm:
		cmdArgs = []string{"npm", "run", "dev", "--", "--port", fmt.Sprintf("%d", s.port)}
	default:
		cmdArgs = []string{"npm", "run", "dev", "--", "--port", fmt.Sprintf("%d", s.port)}
	}

	s.cmd = exec.Command(cmdArgs[0], cmdArgs[1:]...)
	s.cmd.Dir = s.workspacePath

	// Set up pipes for output streaming
	stdout, err := s.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := s.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the server process
	if err := s.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start dev server: %w", err)
	}

	// Stream output and detect when server is ready
	go s.streamOutput(stdout, false)
	go s.streamOutput(stderr, true)

	return nil
}

// streamOutput streams server output and detects readiness
func (s *Server) streamOutput(reader io.Reader, isStderr bool) {
	scanner := bufio.NewScanner(reader)
	urlPattern := regexp.MustCompile(`Local\s+http://localhost:(\d+)`)

	for scanner.Scan() {
		line := scanner.Text()

		// Print to appropriate stream
		if isStderr {
			fmt.Fprintln(os.Stderr, line)
		} else {
			fmt.Println(line)
		}

		// Detect server ready
		if matches := urlPattern.FindStringSubmatch(line); matches != nil {
			port := matches[1]
			s.serverURL = fmt.Sprintf("http://localhost:%s", port)

			// Signal that server is ready
			select {
			case s.ready <- true:
			default:
			}

			fmt.Printf("\nServer ready at %s\n", s.serverURL)
		}
	}
}

// WaitReady waits for the server to be ready or times out
func (s *Server) WaitReady(timeout time.Duration) (string, error) {
	select {
	case <-s.ready:
		return s.serverURL, nil
	case <-time.After(timeout):
		// Fallback: assume server is ready
		if s.serverURL == "" {
			s.serverURL = fmt.Sprintf("http://localhost:%d", s.port)
		}
		fmt.Println("Server detection timed out, assuming ready")
		return s.serverURL, nil
	}
}

// GetPID returns the server process ID
func (s *Server) GetPID() int {
	if s.cmd != nil && s.cmd.Process != nil {
		return s.cmd.Process.Pid
	}
	return 0
}

// Wait waits for the server process to exit
func (s *Server) Wait() error {
	if s.cmd == nil {
		return nil
	}
	return s.cmd.Wait()
}

// IsRunning checks if the server is still running
func (s *Server) IsRunning() bool {
	if s.cmd == nil || s.cmd.Process == nil {
		return false
	}

	// Try to signal with signal 0 to check if process exists
	err := s.cmd.Process.Signal(os.Signal(nil))
	return err == nil
}

// GetURL returns the server URL
func (s *Server) GetURL() string {
	if s.serverURL == "" {
		return fmt.Sprintf("http://localhost:%d", s.port)
	}
	return s.serverURL
}
