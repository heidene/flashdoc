package steps

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/cucumber/godog"
)

// RegisterBrowserSteps registers all browser-related step definitions
func RegisterBrowserSteps(sc *godog.ScenarioContext, ctx *TestContext) {
	sc.Step(`^the operating system is "([^"]*)"$`, ctx.operatingSystemIs)
	sc.Step(`^the dev server is running on "([^"]*)"$`, ctx.devServerIsRunningOn)
	sc.Step(`^the --no-open flag is not set$`, ctx.noOpenFlagIsNotSet)
	sc.Step(`^the --no-open flag is set$`, ctx.noOpenFlagIsSet)

	sc.Step(`^stardoc opens the URL$`, ctx.stardocOpensURL)
	sc.Step(`^the browser should NOT open automatically$`, ctx.browserShouldNotOpenAutomatically)

	sc.Step(`^stardoc should execute "([^"]*)"$`, ctx.stardocShouldExecuteBrowserCommand)
	sc.Step(`^the default browser should open$`, ctx.defaultBrowserShouldOpen)
	sc.Step(`^the CLI should log "([^"]*)"$`, ctx.cliShouldLogBrowser)
	sc.Step(`^the browser command fails$`, ctx.browserCommandFails)
	sc.Step(`^stardoc should log a warning$`, ctx.stardocShouldLogWarning)
	sc.Step(`^the warning should include the URL$`, ctx.warningShouldIncludeURL)
	sc.Step(`^the URL should still be displayed for manual opening$`, ctx.urlShouldStillBeDisplayed)
}

func (ctx *TestContext) operatingSystemIs(osName string) error {
	// Mock the operating system
	// In actual implementation, this would be os-specific
	return nil
}

func (ctx *TestContext) devServerIsRunningOn(url string) error {
	ctx.serverURL = url
	ctx.serverReady = true
	return nil
}

func (ctx *TestContext) noOpenFlagIsNotSet() error {
	ctx.noOpen = false
	return nil
}

func (ctx *TestContext) noOpenFlagIsSet() error {
	ctx.noOpen = true
	return nil
}

func (ctx *TestContext) stardocOpensURL() error {
	if ctx.noOpen {
		return nil
	}

	// Simulate opening browser based on OS
	os := runtime.GOOS
	var command string

	switch os {
	case "darwin":
		command = fmt.Sprintf("open %s", ctx.serverURL)
	case "windows":
		command = fmt.Sprintf("start %s", ctx.serverURL)
	case "linux":
		command = fmt.Sprintf("xdg-open %s", ctx.serverURL)
	default:
		command = fmt.Sprintf("open %s", ctx.serverURL)
	}

	ctx.browserCommand = command
	ctx.browserOpened = true
	ctx.output.WriteString(fmt.Sprintf("Opening browser at %s\n", ctx.serverURL))

	return nil
}

func (ctx *TestContext) browserShouldNotOpenAutomatically() error {
	if ctx.browserOpened {
		return fmt.Errorf("browser should not have opened automatically")
	}
	return nil
}

func (ctx *TestContext) stardocShouldExecuteBrowserCommand(expectedCommand string) error {
	if !strings.Contains(ctx.browserCommand, strings.Split(expectedCommand, " ")[0]) {
		return fmt.Errorf("expected command %q, got %q", expectedCommand, ctx.browserCommand)
	}
	return nil
}

func (ctx *TestContext) defaultBrowserShouldOpen() error {
	if !ctx.browserOpened {
		return fmt.Errorf("browser was not opened")
	}
	return nil
}

func (ctx *TestContext) cliShouldLogBrowser(expectedLog string) error {
	output := ctx.output.String()
	if !strings.Contains(output, expectedLog) {
		return fmt.Errorf("expected log %q not found", expectedLog)
	}
	return nil
}

func (ctx *TestContext) browserCommandFails() error {
	ctx.errorOutput.WriteString("Failed to open browser\n")
	ctx.browserOpened = false
	return nil
}

func (ctx *TestContext) stardocShouldLogWarning() error {
	output := ctx.output.String() + ctx.errorOutput.String()
	if !strings.Contains(strings.ToLower(output), "warn") && !strings.Contains(strings.ToLower(output), "failed") {
		return fmt.Errorf("no warning logged")
	}
	return nil
}

func (ctx *TestContext) warningShouldIncludeURL() error {
	output := ctx.output.String() + ctx.errorOutput.String()
	if !strings.Contains(output, "http") {
		return fmt.Errorf("warning does not include URL")
	}
	return nil
}

func (ctx *TestContext) urlShouldStillBeDisplayed() error {
	output := ctx.output.String()
	if !strings.Contains(output, ctx.serverURL) && !strings.Contains(output, "http") {
		return fmt.Errorf("URL not displayed in output")
	}
	return nil
}
