package features

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/nicovandenhove/stardoc/features/steps"
)

func InitializeScenario(sc *godog.ScenarioContext) {
	testCtx := steps.NewTestContext()

	// Reset context before each scenario
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		testCtx.Reset()
		return ctx, nil
	})

	// Phase 1: Foundation
	steps.RegisterCLISteps(sc, testCtx)
	steps.RegisterWorkspaceSteps(sc, testCtx)
	steps.RegisterSignalSteps(sc, testCtx)
	steps.RegisterCleanupSteps(sc, testCtx)

	// Phase 2: Markdown Processing
	steps.RegisterScannerSteps(sc, testCtx)
	steps.RegisterFrontmatterSteps(sc, testCtx)
	steps.RegisterProcessorSteps(sc, testCtx)

	// Phase 3: Starlight Setup
	steps.RegisterTemplateSteps(sc, testCtx)
	steps.RegisterConfigSteps(sc, testCtx)
	steps.RegisterPkgManagerSteps(sc, testCtx)
	steps.RegisterInstallerSteps(sc, testCtx)

	// Phase 4: Server & Browser
	steps.RegisterOutputSteps(sc, testCtx)
	steps.RegisterBrowserSteps(sc, testCtx)
	steps.RegisterServerSteps(sc, testCtx)
	steps.RegisterExportSteps(sc, testCtx)
}

// runPhase is a helper function to run tests for a specific phase
func runPhase(t *testing.T, phasePath string) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{phasePath},
			TestingT: t,
			Strict:   false, // Set to false to allow pending steps
		},
	}

	if suite.Run() != 0 {
		t.Errorf("Phase %s: non-zero status returned", phasePath)
	}
}

// TestPhase1 runs Phase 1: Foundation tests
func TestPhase1(t *testing.T) {
	runPhase(t, "phase1-foundation")
}

// TestPhase2 runs Phase 2: Markdown Processing tests
func TestPhase2(t *testing.T) {
	runPhase(t, "phase2-markdown")
}

// TestPhase3 runs Phase 3: Starlight Setup tests
func TestPhase3(t *testing.T) {
	runPhase(t, "phase3-starlight")
}

// TestPhase4 runs Phase 4: Server & Browser tests
func TestPhase4(t *testing.T) {
	runPhase(t, "phase4-server")
}

// TestFeatures runs all phases together
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format: "pretty",
			Paths: []string{
				"phase1-foundation",
				"phase2-markdown",
				"phase3-starlight",
				"phase4-server",
			},
			TestingT: t,
			Strict:   false, // Set to false to allow pending steps
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned")
	}
}
