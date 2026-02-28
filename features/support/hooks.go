package support

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/baphled/kariya/internal/config"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	"github.com/cucumber/godog"
	"gorm.io/gorm"
)

// testingT holds the current testing.T for scenario setup.
var testingT *testing.T

// sharedGormDB holds the suite-level GORM connection reused across all scenarios.
var sharedGormDB *gorm.DB

// sharedSQLDB holds the underlying sql.DB for closing during suite teardown.
var sharedSQLDB *sql.DB

// sharedTmpDir holds the temporary directory containing the shared test database.
var sharedTmpDir string

// txKey is the context key for storing a per-scenario GORM transaction.
type txKey struct{}

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

// RegisterSuiteHooks registers suite-level hooks that open a shared database
// and run migrations once before all scenarios, then tear down after all scenarios.
//
// Expected:
//   - sc is a valid *godog.TestSuiteContext.
//
// Side effects:
//   - Creates a temporary database directory and opens a shared SQLite DB.
//   - Runs all migrations once.
//   - Cleans up the database and directory after all scenarios complete.
func RegisterSuiteHooks(sc *godog.TestSuiteContext) {
	sc.BeforeSuite(func() {
		var err error
		sharedTmpDir, err = os.MkdirTemp("", "bdd_suite_*")
		if err != nil {
			panic("failed to create temp dir: " + err.Error())
		}

		configPath := filepath.Join(sharedTmpDir, "config.yaml")
		config.SetConfigPathForTesting(configPath)

		dbPath := filepath.Join(sharedTmpDir, "bdd_suite.db")

		sharedSQLDB, err = careersql.OpenDB(dbPath)
		if err != nil {
			panic("failed to open shared test db: " + err.Error())
		}

		if err := careerrepo.RunMigrationsForTests(sharedSQLDB); err != nil {
			_ = sharedSQLDB.Close()
			panic("failed to run migrations: " + err.Error())
		}

		sharedGormDB, err = careersql.NewGormDB(sharedSQLDB)
		if err != nil {
			_ = sharedSQLDB.Close()
			panic("failed to create GORM connection: " + err.Error())
		}
	})

	sc.AfterSuite(func() {
		if sharedSQLDB != nil {
			_ = sharedSQLDB.Close()
		}
		config.ResetConfigPath()
		if sharedTmpDir != "" {
			_ = os.RemoveAll(sharedTmpDir)
		}
		sharedGormDB = nil
		sharedSQLDB = nil
	})
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
// When a shared database is available, it begins a transaction for isolation.
func setupScenarioEnv(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
	if needsAppEnv(sc) && testingT != nil {
		if sharedGormDB != nil {
			tx := sharedGormDB.Begin()
			ctx = context.WithValue(ctx, txKey{}, tx)
			env := NewAppEnvFromGormDB(testingT, tx)
			ctx = WithAppEnv(ctx, env)
		} else {
			env := NewAppEnv(testingT)
			ctx = WithAppEnv(ctx, env)
		}
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
// It rolls back the per-scenario transaction to reset data without re-migrating.
func afterScenario(ctx context.Context, _ *godog.Scenario, _ error) (context.Context, error) {
	tx, ok := ctx.Value(txKey{}).(*gorm.DB)
	if ok && tx != nil {
		tx.Rollback()
	}
	env := GetAppEnv(ctx)
	if env != nil {
		env.Cleanup()
	}
	return ctx, nil
}
