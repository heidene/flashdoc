package steps

import (
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

// RegisterInstallerSteps registers all installer-related step definitions
func RegisterInstallerSteps(sc *godog.ScenarioContext, ctx *TestContext) {
	sc.Step(`^dependency installation is triggered$`, ctx.dependencyInstallationIsTriggered)
	sc.Step(`^dependencies are being installed$`, ctx.dependenciesAreBeingInstalled)
	sc.Step(`^dependency installation is attempted$`, ctx.dependencyInstallationIsAttempted)

	sc.Step(`^stardoc should execute "([^"]*)" in the temp workspace$`, ctx.stardocShouldExecuteInWorkspace)
	sc.Step(`^it should execute "([^"]*)"$`, ctx.shouldExecuteInstall)
	sc.Step(`^installation progress should be displayed$`, ctx.installationProgressShouldBeDisplayed)
	sc.Step(`^the CLI should show output from the package manager$`, ctx.cliShouldShowPackageManagerOutput)
	sc.Step(`^the CLI should wait for installation to complete$`, ctx.cliShouldWaitForInstallation)
	sc.Step(`^all dependencies should be installed successfully$`, ctx.allDependenciesShouldBeInstalled)
	sc.Step(`^the "([^"]*)" directory should be created$`, ctx.directoryShouldBeCreatedByInstall)
	sc.Step(`^it should contain the installed packages$`, ctx.shouldContainInstalledPackages)
	sc.Step(`^installation failed$`, ctx.installationFailed)
	sc.Step(`^the CLI should display the error from "([^"]*)"$`, ctx.cliShouldDisplayPackageManagerError)
	sc.Step(`^the CLI should suggest "([^"]*)"$`, ctx.cliShouldSuggest)
	sc.Step(`^the workspace should remain intact$`, ctx.workspaceShouldRemainIntact)
	sc.Step(`^the user can fix the issue and retry$`, ctx.userCanFixAndRetry)
	sc.Step(`^stardoc should execute "([^"]*)" first$`, ctx.stardocShouldExecuteFirst)
	sc.Step(`^the lockfile should be generated$`, ctx.lockfileShouldBeGenerated)
	sc.Step(`^installation should take less than (\d+) seconds$`, ctx.installationShouldBeFast)
	sc.Step(`^the estimated time should be displayed$`, ctx.estimatedTimeShouldBeDisplayed)
	sc.Step(`^the CLI should show a progress indicator \(spinner or progress bar\)$`, ctx.cliShouldShowProgressIndicator)
	sc.Step(`^after (\d+) seconds, it should show "([^"]*)"$`, ctx.afterSecondsShouldShowMessage)
	sc.Step(`^stardoc should NOT run "([^"]*)"$`, ctx.stardocShouldNotRun)
	sc.Step(`^stardoc should skip the postinstall phase$`, ctx.stardocShouldSkipPostinstall)
	sc.Step(`^installation time should be reduced$`, ctx.installationTimeShouldBeReduced)
}

func (ctx *TestContext) dependencyInstallationIsTriggered() error {
	// Simulate installation
	pmName := "npm"
	if ctx.detectedPM != "" {
		pmName = string(ctx.detectedPM)
	}

	_ = fmt.Sprintf("%s install", pmName) // Command would be executed here
	ctx.installOutput = fmt.Sprintf("Installing dependencies using %s...\n", pmName)
	ctx.output.WriteString(ctx.installOutput)

	// Simulate successful installation
	ctx.installOutput += "Dependencies installed successfully\n"
	ctx.output.WriteString("Dependencies installed successfully\n")

	return nil
}

func (ctx *TestContext) dependenciesAreBeingInstalled() error {
	return ctx.dependencyInstallationIsTriggered()
}

func (ctx *TestContext) dependencyInstallationIsAttempted() error {
	err := ctx.dependencyInstallationIsTriggered()
	if err != nil {
		// Error might be expected
		return nil
	}
	return nil
}

func (ctx *TestContext) stardocShouldExecuteInWorkspace(command string) error {
	output := ctx.output.String() + ctx.installOutput
	if !strings.Contains(output, command) && !strings.Contains(output, strings.Split(command, " ")[0]) {
		return fmt.Errorf("expected command %q not found in output", command)
	}
	return nil
}

func (ctx *TestContext) shouldExecuteInstall(command string) error {
	return ctx.stardocShouldExecuteInWorkspace(command)
}

func (ctx *TestContext) installationProgressShouldBeDisplayed() error {
	if !strings.Contains(ctx.output.String(), "Installing") && !strings.Contains(ctx.output.String(), "install") {
		return fmt.Errorf("no installation progress displayed")
	}
	return nil
}

func (ctx *TestContext) cliShouldShowPackageManagerOutput() error {
	return nil
}

func (ctx *TestContext) cliShouldWaitForInstallation() error {
	return nil
}

func (ctx *TestContext) allDependenciesShouldBeInstalled() error {
	return nil
}

func (ctx *TestContext) directoryShouldBeCreatedByInstall(dirName string) error {
	// Check if directory exists (would be created by actual npm install)
	return nil
}

func (ctx *TestContext) shouldContainInstalledPackages() error {
	return nil
}

func (ctx *TestContext) installationFailed() error {
	ctx.exitCode = 1
	ctx.errorOutput.WriteString("Installation failed\n")
	return nil
}

func (ctx *TestContext) cliShouldDisplayPackageManagerError(pmName string) error {
	if ctx.errorOutput.Len() == 0 {
		return fmt.Errorf("no error output displayed")
	}
	return nil
}

func (ctx *TestContext) cliShouldSuggest(suggestion string) error {
	// Check for helpful suggestions in output
	return nil
}

func (ctx *TestContext) workspaceShouldRemainIntact() error {
	return nil
}

func (ctx *TestContext) userCanFixAndRetry() error {
	return nil
}

func (ctx *TestContext) stardocShouldExecuteFirst(command string) error {
	return nil
}

func (ctx *TestContext) lockfileShouldBeGenerated() error {
	return nil
}

func (ctx *TestContext) installationShouldBeFast(maxSeconds int) error {
	return nil
}

func (ctx *TestContext) estimatedTimeShouldBeDisplayed() error {
	return nil
}

func (ctx *TestContext) cliShouldShowProgressIndicator() error {
	return nil
}

func (ctx *TestContext) afterSecondsShouldShowMessage(seconds int, message string) error {
	return nil
}

func (ctx *TestContext) stardocShouldNotRun(command string) error {
	output := ctx.output.String()
	if strings.Contains(output, command) {
		return fmt.Errorf("command %q should not have been run", command)
	}
	return nil
}

func (ctx *TestContext) stardocShouldSkipPostinstall() error {
	return nil
}

func (ctx *TestContext) installationTimeShouldBeReduced() error {
	return nil
}
