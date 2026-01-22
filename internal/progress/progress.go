package progress

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/charmbracelet/lipgloss"
)

var (
	// Styles
	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ff00")).
			Bold(true)

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00d7ff")).
			Bold(true)

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))
)

// Spinner provides a styled progress indicator
type Spinner struct {
	spinner   *spinner.Spinner
	startTime time.Time
	message   string
}

// New creates a new progress spinner
func New(message string) *Spinner {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Suffix = " " + infoStyle.Render(message)
	s.Writer = os.Stderr

	return &Spinner{
		spinner:   s,
		startTime: time.Now(),
		message:   message,
	}
}

// Start begins the spinner animation
func (s *Spinner) Start() {
	s.startTime = time.Now()
	s.spinner.Start()
}

// Stop stops the spinner and shows a success message
func (s *Spinner) Stop(successMessage string) {
	s.spinner.Stop()
	duration := time.Since(s.startTime)

	fmt.Fprintf(os.Stderr, "\r%s %s %s\n",
		successStyle.Render("✓"),
		successMessage,
		dimStyle.Render(fmt.Sprintf("(%s)", formatDuration(duration))),
	)
}

// StopWithError stops the spinner and shows an error message
func (s *Spinner) StopWithError(errorMessage string) {
	s.spinner.Stop()

	errorStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#ff0000")).
		Bold(true)

	fmt.Fprintf(os.Stderr, "\r%s %s\n",
		errorStyle.Render("✗"),
		errorMessage,
	)
}

// Update changes the spinner message
func (s *Spinner) Update(message string) {
	s.message = message
	s.spinner.Suffix = " " + infoStyle.Render(message)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

// DiscardWriter returns a writer that discards all output
func DiscardWriter() io.Writer {
	return io.Discard
}

// Witty messages for different operations
var (
	InstallMessages = []string{
		"Summoning the npm spirits...",
		"Bribing the package manager...",
		"Convincing dependencies to cooperate...",
		"Downloading the internet, one package at a time...",
		"Teaching packages to play nicely together...",
	}

	BuildMessages = []string{
		"Weaving markdown into HTML magic...",
		"Teaching Astro some new tricks...",
		"Transforming your docs into pixel perfection...",
		"Compiling dreams into reality...",
		"Building something beautiful...",
		"Converting markdown to awesome...",
	}
)
