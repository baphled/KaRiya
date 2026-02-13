package features

import (
	"flag"
	"os"
	"testing"

	"github.com/baphled/kariya/features/steps"
	"github.com/baphled/kariya/features/support"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/format"
)

// opts holds the Godog options for test execution.
var opts = godog.Options{
	Output: colors.Colored(os.Stdout),
	Format: "progress",
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
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
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
	steps.RegisterConfigureSteps(sc)
	steps.RegisterNavigationSteps(sc)
	steps.RegisterCLISteps(sc)
}
