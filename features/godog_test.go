package features

import (
	"flag"
	"os"
	"testing"

	"github.com/baphled/kariya/features/steps"
	"github.com/baphled/kariya/features/support"
	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// opts holds the Godog options for test execution.
var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "pretty",
	Tags:   "~@wip",
	Strict: true,
}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opts)
}

// TestFeatures runs all BDD feature tests.
//
// Expected: Feature files exist in the features/ directory.
// Returns: None.
// Side effects: Runs all scenarios and reports results.
func TestFeatures(t *testing.T) {
	RegisterTestingT(t)
	support.SetTestingT(t)

	// Increase Gomega's format.MaxLength to prevent truncation of long view outputs
	format.MaxLength = 0 // 0 = unlimited

	opts.Paths = []string{"./"}
	opts.TestingT = t

	suite := godog.TestSuite{
		TestSuiteInitializer: InitializeSuite,
		ScenarioInitializer:  InitializeScenario,
		Options:              &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

// InitializeSuite registers suite-level hooks for shared database setup and teardown.
//
// Expected: sc is a valid TestSuiteContext.
// Returns: None.
// Side effects: Registers suite-level hooks via support.RegisterSuiteHooks.
func InitializeSuite(sc *godog.TestSuiteContext) {
	support.RegisterSuiteHooks(sc)
}

// InitializeScenario sets up the scenario context with step definitions and hooks.
//
// Expected: sc is a valid ScenarioContext.
// Returns: None.
// Side effects: Registers steps and hooks with the scenario context.
func InitializeScenario(sc *godog.ScenarioContext) {
	support.RegisterHooks(sc)
	steps.RegisterCommonSteps(sc)
	steps.RegisterOnboardingSteps(sc)
	steps.RegisterCaptureSteps(sc)
	steps.RegisterBrowseSteps(sc)
	steps.RegisterSkillsSteps(sc)
	steps.RegisterBurstsSteps(sc)
	steps.RegisterFactsSteps(sc)
	steps.RegisterCVSteps(sc)
}

// InitializeSuite sets up suite-level hooks for shared database management.
//
// Expected: ctx is a valid TestSuiteContext.
// Returns: None.
// Side effects: Registers BeforeSuite and AfterSuite hooks.
func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		e2e.SetupShared()
	})
	ctx.AfterSuite(func() {
		e2e.CleanupShared()
	})
}
