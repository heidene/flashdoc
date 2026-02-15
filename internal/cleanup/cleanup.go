package cleanup

import (
	"fmt"
	"os"
	"sync"

	"github.com/nicovandenhove/flashdoc/internal/staticserver"
	"github.com/nicovandenhove/flashdoc/internal/workspace"
)

// Manager handles cleanup of resources when stardoc exits
type Manager struct {
	workspace    *workspace.Workspace
	server       *staticserver.Server
	shutdownOnce sync.Once
	mu           sync.Mutex
}

// New creates a new cleanup manager
func New(ws *workspace.Workspace) *Manager {
	return &Manager{
		workspace: ws,
	}
}

// RegisterServer adds the HTTP server to be stopped on cleanup
func (m *Manager) RegisterServer(server *staticserver.Server) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.server = server
}

// Cleanup stops the server and removes the workspace
func (m *Manager) Cleanup() error {
	var cleanupErr error

	m.shutdownOnce.Do(func() {
		// Stop server first
		if err := m.StopServer(); err != nil {
			fmt.Fprintf(os.Stderr, "Stopping server...\n")
			cleanupErr = fmt.Errorf("failed to stop server: %w", err)
		}

		// Remove workspace
		fmt.Fprintf(os.Stderr, "Cleaning up workspace...\n")
		if m.workspace != nil {
			if err := m.workspace.Cleanup(); err != nil {
				if cleanupErr != nil {
					cleanupErr = fmt.Errorf("%v; failed to cleanup workspace: %w", cleanupErr, err)
				} else {
					// Log warning but don't fail if workspace cleanup fails
					fmt.Fprintf(os.Stderr, "Warning: failed to remove workspace %s: %v\n", m.workspace.Path, err)
					fmt.Fprintf(os.Stderr, "You may need to manually remove this directory\n")
				}
			}
		}

		if cleanupErr == nil {
			fmt.Fprintf(os.Stderr, "Cleanup complete\n")
		}
	})

	return cleanupErr
}

// StopServer gracefully shuts down the HTTP server
func (m *Manager) StopServer() error {
	m.mu.Lock()
	server := m.server
	m.mu.Unlock()

	if server == nil {
		return nil
	}

	return server.Stop()
}
