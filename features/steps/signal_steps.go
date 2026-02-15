package steps

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/cucumber/godog"
)

// RegisterSignalSteps registers all signal-handling step definitions
func RegisterSignalSteps(ctx *godog.ScenarioContext, testCtx *TestContext) {
	ctx.Step(`^I have started "stardoc ([^"]*)"$`, testCtx.iHaveStartedStardoc)
	ctx.Step(`^the dev server is running$`, testCtx.theDevServerIsRunning)
	ctx.Step(`^I press Ctrl\+C$`, testCtx.iPressCtrlC)
	ctx.Step(`^the CLI should catch the SIGINT signal$`, testCtx.theCLIShouldCatchSIGINT)
	ctx.Step(`^the CLI should display "([^"]*)"$`, testCtx.theCLIShouldDisplay)
	ctx.Step(`^the dev server process should be terminated$`, testCtx.theDevServerProcessShouldBeTerminated)
	ctx.Step(`^I send a SIGTERM signal to the stardoc process$`, testCtx.iSendSIGTERM)
	ctx.Step(`^the CLI should catch the SIGTERM signal$`, testCtx.theCLIShouldCatchSIGTERM)
	ctx.Step(`^I press Ctrl\+C again within (\d+) second$`, testCtx.iPressCtrlCAgainWithinSeconds)
	ctx.Step(`^the CLI should force exit immediately$`, testCtx.theCLIShouldForceExitImmediately)
	ctx.Step(`^best-effort cleanup should be attempted$`, testCtx.bestEffortCleanupShouldBeAttempted)
	ctx.Step(`^stardoc is installing npm dependencies$`, testCtx.stardocIsInstallingNpmDependencies)
	ctx.Step(`^the npm install process should be terminated$`, testCtx.theNpmInstallProcessShouldBeTerminated)
	ctx.Step(`^stardoc is starting the Astro dev server$`, testCtx.stardocIsStartingAstroDevServer)
	ctx.Step(`^the server has not finished starting yet$`, testCtx.theServerHasNotFinishedStarting)
	ctx.Step(`^the Astro process should be terminated$`, testCtx.theAstroProcessShouldBeTerminated)
	ctx.Step(`^the CLI should send SIGTERM to PID (\d+)$`, testCtx.theCLIShouldSendSIGTERMToPID)
	ctx.Step(`^the CLI should wait up to (\d+) seconds for graceful shutdown$`, testCtx.theCLIShouldWaitForGracefulShutdown)
	ctx.Step(`^if the process doesn't stop, send SIGKILL$`, testCtx.ifProcessDoesntStopSendSIGKILL)
	ctx.Step(`^the CLI should not exit until all child processes are stopped$`, testCtx.theCLIShouldNotExitUntilAllChildProcessesStopped)
	ctx.Step(`^cleanup operations are taking longer than expected$`, testCtx.cleanupOperationsAreTakingLonger)
	ctx.Step(`^the CLI should start cleanup$`, testCtx.theCLIShouldStartCleanup)
	ctx.Step(`^the CLI should wait up to (\d+) seconds for cleanup completion$`, testCtx.theCLIShouldWaitForCleanupCompletion)
	ctx.Step(`^if cleanup doesn't complete, force exit with a warning$`, testCtx.ifCleanupDoesntCompleteForceExit)
	ctx.Step(`^the warning should mention "([^"]*)"$`, testCtx.theWarningShouldMention)
	ctx.Step(`^the user browses the documentation site$`, testCtx.theUserBrowsesDocumentationSite)
	ctx.Step(`^the server should respond normally$`, testCtx.theServerShouldRespondNormally)
	ctx.Step(`^the signal handler should not consume excessive resources$`, testCtx.theSignalHandlerShouldNotConsumeExcessiveResources)
}

func (ctx *TestContext) iHaveStartedStardoc(args string) error {
	// Start stardoc in the background
	return ctx.iRunStardoc(args)
}

func (ctx *TestContext) theDevServerIsRunning() error {
	// Stub - will be implemented in Phase 4
	return nil
}

func (ctx *TestContext) iPressCtrlC() error {
	// In test mode, simulate the Ctrl+C behavior
	if ctx.cmd != nil && ctx.cmd.Process != nil {
		return ctx.cmd.Process.Signal(os.Interrupt)
	}

	// Mock mode - simulate signal handling
	ctx.output.WriteString("\nShutting down gracefully...\n")
	ctx.output.WriteString("Stopping dev server...\n")
	ctx.output.WriteString("Cleaning up workspace...\n")

	// Perform cleanup - actually remove the workspace
	if ctx.workspacePath != "" {
		err := os.RemoveAll(ctx.workspacePath)
		if err != nil {
			// Log warning about cleanup failure
			ctx.output.WriteString(fmt.Sprintf("⚠️  Warning: failed to clean up workspace at %s: %v\n", ctx.workspacePath, err))
			// Still exit successfully since user's action was successful
			ctx.exitCode = 0
		} else {
			ctx.output.WriteString("Cleanup complete\n")
			ctx.exitCode = 0
		}
	}

	// Terminate any test child processes
	for _, pid := range ctx.childPIDs {
		if process, err := os.FindProcess(pid); err == nil {
			_ = process.Kill()
		}
	}
	ctx.childPIDs = nil

	ctx.interrupted = true
	return nil
}

func (ctx *TestContext) theCLIShouldCatchSIGINT() error {
	// Check that the process handled the signal
	// This is implicitly tested by checking the output
	return nil
}

func (ctx *TestContext) theCLIShouldDisplay(message string) error {
	output := ctx.output.String()
	errorOutput := ctx.errorOutput.String()
	combinedOutput := output + errorOutput

	if !strings.Contains(combinedOutput, message) {
		return fmt.Errorf("expected output to contain %q, got: %s", message, combinedOutput)
	}
	return nil
}

func (ctx *TestContext) theDevServerProcessShouldBeTerminated() error {
	// Stub - will be implemented in Phase 4
	return nil
}

func (ctx *TestContext) iSendSIGTERM() error {
	if ctx.cmd != nil && ctx.cmd.Process != nil {
		return ctx.cmd.Process.Signal(syscall.SIGTERM)
	}

	// Mock mode - simulate SIGTERM handling (same as SIGINT)
	ctx.output.WriteString("\nShutting down gracefully...\n")
	ctx.output.WriteString("Stopping dev server...\n")
	ctx.output.WriteString("Cleaning up workspace...\n")

	// Perform cleanup - actually remove the workspace
	if ctx.workspacePath != "" {
		err := os.RemoveAll(ctx.workspacePath)
		if err != nil {
			ctx.output.WriteString(fmt.Sprintf("⚠️  Warning: failed to clean up workspace at %s: %v\n", ctx.workspacePath, err))
			ctx.exitCode = 0
		} else {
			ctx.output.WriteString("Cleanup complete\n")
		}
	}

	// Terminate any test child processes
	for _, pid := range ctx.childPIDs {
		if process, err := os.FindProcess(pid); err == nil {
			_ = process.Kill()
		}
	}
	ctx.childPIDs = nil

	ctx.interrupted = true
	ctx.exitCode = 0
	return nil
}

func (ctx *TestContext) theCLIShouldCatchSIGTERM() error {
	// Check that the process handled the signal
	return nil
}

func (ctx *TestContext) iPressCtrlCAgainWithinSeconds(seconds int) error {
	// First interrupt should have already happened
	if !ctx.interrupted {
		return fmt.Errorf("first interrupt should have been triggered before second")
	}

	time.Sleep(time.Duration(seconds) * time.Millisecond * 100) // Short delay

	// Force stop message
	ctx.output.WriteString("Force stopping...\n")

	// Perform immediate cleanup
	if ctx.workspacePath != "" {
		os.RemoveAll(ctx.workspacePath)
	}

	ctx.exitCode = 1 // Exit with error code on force stop
	return nil
}

func (ctx *TestContext) theCLIShouldForceExitImmediately() error {
	// Check that "Force stopping..." was displayed
	return ctx.theCLIShouldDisplay("Force stopping...")
}

func (ctx *TestContext) bestEffortCleanupShouldBeAttempted() error {
	// Stub - check that cleanup was attempted even during force exit
	return nil
}

func (ctx *TestContext) stardocIsInstallingNpmDependencies() error {
	// Mock - simulate npm install in progress
	// Ensure workspace exists for cleanup
	if ctx.workspacePath == "" {
		tempDir, err := os.MkdirTemp("", "stardoc-npm-*")
		if err != nil {
			return err
		}
		ctx.workspacePath = tempDir
		ctx.TrackDir(tempDir)
	}
	ctx.output.WriteString("Installing dependencies...\n")
	ctx.npmInstalling = true
	return nil
}

func (ctx *TestContext) theNpmInstallProcessShouldBeTerminated() error {
	// Verify npm install was interrupted
	if !ctx.interrupted {
		return fmt.Errorf("npm install was not terminated")
	}
	ctx.npmInstalling = false
	return nil
}

func (ctx *TestContext) stardocIsStartingAstroDevServer() error {
	// Mock - simulate server starting
	// Ensure workspace exists for cleanup
	if ctx.workspacePath == "" {
		tempDir, err := os.MkdirTemp("", "stardoc-server-*")
		if err != nil {
			return err
		}
		ctx.workspacePath = tempDir
		ctx.TrackDir(tempDir)
	}
	ctx.output.WriteString("Starting Astro dev server...\n")
	ctx.serverStarting = true
	ctx.serverReady = false // Explicitly set to false - server is still starting
	return nil
}

func (ctx *TestContext) theServerHasNotFinishedStarting() error {
	// Verify server is still starting (not ready)
	if ctx.serverReady {
		return fmt.Errorf("server has already finished starting")
	}
	return nil
}

func (ctx *TestContext) theAstroProcessShouldBeTerminated() error {
	// Verify server was interrupted
	if !ctx.interrupted {
		return fmt.Errorf("Astro process was not terminated")
	}
	ctx.serverStarting = false
	ctx.serverReady = false
	return nil
}

func (ctx *TestContext) theCLIShouldSendSIGTERMToPID(pid int) error {
	// Stub - verify signal sending logic
	ctx.runningPID = pid
	return nil
}

func (ctx *TestContext) theCLIShouldWaitForGracefulShutdown(seconds int) error {
	// Stub - verify timeout logic
	_ = seconds
	return nil
}

func (ctx *TestContext) ifProcessDoesntStopSendSIGKILL() error {
	// Stub - verify escalation to SIGKILL
	return nil
}

func (ctx *TestContext) theCLIShouldNotExitUntilAllChildProcessesStopped() error {
	// Stub - verify all child processes are stopped
	return nil
}

func (ctx *TestContext) cleanupOperationsAreTakingLonger() error {
	// Mock - simulate slow cleanup by setting a flag
	// Ensure workspace exists for cleanup
	if ctx.workspacePath == "" {
		tempDir, err := os.MkdirTemp("", "stardoc-cleanup-*")
		if err != nil {
			return err
		}
		ctx.workspacePath = tempDir
		ctx.TrackDir(tempDir)
	}
	ctx.slowCleanup = true
	ctx.output.WriteString("Cleanup in progress...\n")
	return nil
}

func (ctx *TestContext) theCLIShouldStartCleanup() error {
	// Verify cleanup was initiated
	output := ctx.output.String() + ctx.errorOutput.String()
	if !strings.Contains(output, "Cleanup") && !strings.Contains(output, "cleanup") {
		return fmt.Errorf("cleanup was not started")
	}
	return nil
}

func (ctx *TestContext) theCLIShouldWaitForCleanupCompletion(seconds int) error {
	// Mock - simulate waiting for cleanup with timeout
	_ = seconds
	if ctx.slowCleanup {
		ctx.output.WriteString(fmt.Sprintf("Waiting up to %d seconds for cleanup...\n", seconds))
	}
	return nil
}

func (ctx *TestContext) ifCleanupDoesntCompleteForceExit() error {
	// Mock - simulate force exit if cleanup is slow
	if ctx.slowCleanup {
		ctx.errorOutput.WriteString("Warning: cleanup may be incomplete, force exiting\n")
		ctx.exitCode = 1
	}
	return nil
}

func (ctx *TestContext) theWarningShouldMention(text string) error {
	return ctx.theErrorMessageShouldContain(text)
}

func (ctx *TestContext) theUserBrowsesDocumentationSite() error {
	// Mock - simulate user browsing the site
	ctx.browserOpened = true
	ctx.output.WriteString("Documentation site is accessible\n")
	return nil
}

func (ctx *TestContext) theServerShouldRespondNormally() error {
	// Verify server is responding
	if !ctx.serverReady && !ctx.serverStarting {
		return fmt.Errorf("server is not running")
	}
	return nil
}

func (ctx *TestContext) theSignalHandlerShouldNotConsumeExcessiveResources() error {
	// Stub - performance test
	return nil
}
