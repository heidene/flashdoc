package signal

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Handler manages signal handling for graceful shutdown
type Handler struct {
	cleanupFunc func() error
	forceCount  int
	shutdownCh  chan struct{}
	signalCh    chan os.Signal
	mu          sync.Mutex
}

// New creates a new signal handler with the given cleanup function
func New(cleanup func() error) *Handler {
	return &Handler{
		cleanupFunc: cleanup,
		forceCount:  0,
		shutdownCh:  make(chan struct{}),
		signalCh:    make(chan os.Signal, 1),
	}
}

// Setup registers signal handlers for SIGINT and SIGTERM
func (h *Handler) Setup() {
	signal.Notify(h.signalCh, os.Interrupt, syscall.SIGTERM)

	go h.handleSignals()
}

func (h *Handler) handleSignals() {
	var firstSignalTime time.Time

	for {
		select {
		case sig := <-h.signalCh:
			h.mu.Lock()
			h.forceCount++
			currentCount := h.forceCount
			h.mu.Unlock()

			if currentCount == 1 {
				// First signal - graceful shutdown
				firstSignalTime = time.Now()
				fmt.Fprintln(os.Stderr, "\nShutting down gracefully...")
				fmt.Fprintln(os.Stderr, "Press Ctrl+C again to force exit")

				go h.performCleanup()

			} else if currentCount == 2 {
				// Second signal within timeout - force exit
				if time.Since(firstSignalTime) <= 1*time.Second {
					fmt.Fprintln(os.Stderr, "Force stopping...")
					os.Exit(1)
				} else {
					// Reset if too much time has passed
					h.mu.Lock()
					h.forceCount = 1
					h.mu.Unlock()
					firstSignalTime = time.Now()
					fmt.Fprintln(os.Stderr, "\nShutting down gracefully...")
					fmt.Fprintln(os.Stderr, "Press Ctrl+C again to force exit")
					go h.performCleanup()
				}
			} else {
				// Multiple signals - force exit immediately
				fmt.Fprintln(os.Stderr, "Force stopping...")
				os.Exit(1)
			}

			_ = sig // Acknowledge signal usage
		case <-h.shutdownCh:
			return
		}
	}
}

func (h *Handler) performCleanup() {
	// Create a timeout for cleanup
	cleanupDone := make(chan error, 1)

	go func() {
		cleanupDone <- h.cleanupFunc()
	}()

	select {
	case err := <-cleanupDone:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cleanup error: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)

	case <-time.After(10 * time.Second):
		fmt.Fprintln(os.Stderr, "Cleanup timeout - some resources may not have been cleaned up properly")
		os.Exit(1)
	}
}

// Wait blocks until a shutdown signal is received
func (h *Handler) Wait() {
	<-h.shutdownCh
}

// TriggerShutdown manually triggers a shutdown
func (h *Handler) TriggerShutdown() {
	close(h.shutdownCh)
}
