package support

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
)

// testingT holds the current testing.T for scenario setup.
var testingT *testing.T

// SetTestingT sets the testing.T to use for scenario setup.
//
// Expected:
//   - t is a valid *testing.T.
//
// Side effects:
//   - Stores t in package-level variable testingT.
//
//nolint:thelper // Not a test helper, configuration function for BDD framework.
func SetTestingT(t *testing.T) {
	testingT = t
}

// RegisterHooks registers BeforeScenario and AfterScenario hooks with Godog.
//
// Expected:
//   - sc is a valid *godog.ScenarioContext.
//
// Side effects:
//   - Registers Before and After hooks with Godog.
func RegisterHooks(sc *godog.ScenarioContext) {
	sc.Before(beforeScenario)
	sc.After(afterScenario)
}

// beforeScenario runs before each scenario.
func beforeScenario(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	return setupScenarioEnv(ctx, sc)
}

// setupScenarioEnv creates environment for app-level scenarios.
func setupScenarioEnv(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	if needsAppEnv(sc) && testingT != nil {
		env := NewAppEnv(testingT)
		ctx = WithAppEnv(ctx, env)
		ctx = WithEventData(ctx, &EventData{})
	}
	return ctx, nil
}

// needsAppEnv checks if scenario requires full app environment.
// It checks for feature-level tags that indicate the main app is needed.
func needsAppEnv(sc *godog.Scenario) bool {
	appTags := map[string]bool{
		"@capture":    true,
		"@browse":     true,
		"@skills":     true,
		"@bursts":     true,
		"@facts":      true,
		"@cv":         true,
		"@configure":  true,
		"@navigation": true,
	}
	for _, tag := range sc.Tags {
		if appTags[tag.Name] {
			return true
		}
	}
	return false
}

// afterScenario runs after each scenario.
func afterScenario(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	env := GetAppEnv(ctx)
	if env != nil {
		env.Cleanup()
	}
	return ctx, nil
}
