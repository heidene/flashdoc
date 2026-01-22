package staticserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// Server wraps Go's HTTP file server for serving static sites
type Server struct {
	distPath string
	port     int
	server   *http.Server
	output   io.Writer
}

// NewServer creates a new static file server
func NewServer(distPath string, port int, output io.Writer) *Server {
	if output == nil {
		output = os.Stdout
	}

	return &Server{
		distPath: distPath,
		port:     port,
		output:   output,
	}
}

// Start starts the HTTP server in a goroutine
func (s *Server) Start() error {
	// Check if port is available
	if !s.isPortAvailable() {
		return fmt.Errorf("port %d is already in use", s.port)
	}

	// Create file server handler
	fs := http.FileServer(http.Dir(s.distPath))

	// Create mux and handle all routes
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	// Create HTTP server
	s.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		fmt.Fprintf(s.output, "🚀 Server started at http://localhost:%d\n", s.port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(s.output, "Server error: %v\n", err)
		}
	}()

	// Wait for server to be ready
	if err := s.WaitReady(5 * time.Second); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}

// Stop gracefully shuts down the server
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	fmt.Fprintf(s.output, "🛑 Server stopped\n")
	return nil
}

// WaitReady waits for the server to be ready to accept connections
func (s *Server) WaitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", s.port), 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("server did not become ready within %v", timeout)
}

// GetURL returns the server URL
func (s *Server) GetURL() string {
	return fmt.Sprintf("http://localhost:%d", s.port)
}

// GetPort returns the server port
func (s *Server) GetPort() int {
	return s.port
}

// isPortAvailable checks if the port is available
func (s *Server) isPortAvailable() bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return false
	}
	ln.Close()
	return true
}
