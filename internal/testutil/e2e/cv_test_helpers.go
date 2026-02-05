// Package e2e provides E2E test utilities for integration testing of the KaRiya TUI.
package e2e

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"github.com/baphled/kariya/internal/cli/app"
	"github.com/baphled/kariya/internal/cli/bootstrap"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	cvscreens "github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
)

// CVMockConfig holds mock service configuration for CV E2E tests.
// Uses the standardized GoMock types from testutil/mocks/service.
type CVMockConfig struct {
	CVGenService    cv.CVGenerationService
	ClipboardWriter cv.ClipboardWriter
}

// SetupWithCVMocks creates an E2E test environment with mock CV services.
// The returned TestEnv uses a custom intent registrar that injects the mock
// services into the GenerateCV intent, allowing tests to control CV generation
// and export behavior.
//
// Expected:
//   - t must be a valid TestingT.
//   - mockCfg must be a valid CVMockConfig with configured mock services.
//
// Returns:
//   - A fully initialized TestEnv with mock CV services.
//
// Side effects:
//   - Creates temporary directory and SQLite database.
func SetupWithCVMocks(t TestingT, mockCfg *CVMockConfig) *TestEnv {
	t.Helper()

	ctx := context.Background()
	tmpDir, err := os.MkdirTemp("", "e2e_cv_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	dbPath := filepath.Join(tmpDir, "e2e_cv_test.db")

	configPath := filepath.Join(tmpDir, "config.yaml")
	prevConfigPath := config.SwapConfigPathForTesting(configPath)

	db, err := careersql.OpenDB(dbPath)
	if err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := careerrepo.RunMigrationsForTests(db); err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		_ = db.Close()
		t.Fatalf("failed to run migrations: %v", err)
	}

	repos, err := careersql.NewRepositories(db)
	if err != nil {
		config.SetConfigPathForTesting(prevConfigPath)
		_ = db.Close()
		t.Fatalf("failed to create repositories: %v", err)
	}

	svc := careerservice.NewService(repos.Event)
	svc.SetBurstRepository(repos.Burst)
	svc.SetFactRepository(repos.Fact)
	svc.SetSkillRepository(repos.Skill)

	cliService := service.NewCLIEventService(svc)

	log := logger.DefaultLogger()
	bootstrapResult := bootstrap.SkipOnboarding(config.DefaultConfig(), svc, log)

	cvExportService := cv.NewExportServiceWithClipboard(log, mockCfg.ClipboardWriter)

	registrar := newCVMockRegistrar(
		cliService,
		svc,
		bootstrapResult,
		mockCfg.CVGenService,
		cvExportService,
		log,
	)

	model := app.NewModel(
		cliService,
		svc,
		bootstrapResult,
		app.WithIntentRegistrar(registrar),
	)

	cleanup := func() {
		config.SetConfigPathForTesting(prevConfigPath)
		if err := db.Close(); err != nil {
			t.Errorf("warning: failed to close db: %v", err)
		}
		time.Sleep(100 * time.Millisecond)
		_ = os.RemoveAll(tmpDir)
	}

	return &TestEnv{
		T:          t,
		Model:      model,
		DB:         db,
		DBPath:     dbPath,
		EventRepo:  repos.Event,
		BurstRepo:  repos.Burst,
		FactRepo:   repos.Fact,
		SkillRepo:  repos.Skill,
		Service:    svc,
		CLIService: cliService,
		Ctx:        ctx,
		cleanup:    cleanup,
	}
}

// cvMockRegistrar implements app.IntentRegistrar, delegating to
// DefaultIntentRegistrar but overriding the generate_cv intent
// with mock services.
type cvMockRegistrar struct {
	delegate    *app.DefaultIntentRegistrar
	cliService  *service.CLIEventService
	svc         *careerservice.Service
	cvGenSvc    cv.CVGenerationService
	cvExportSvc *cv.ExportService
	log         *logger.Logger
}

func newCVMockRegistrar(
	cliService *service.CLIEventService,
	svc *careerservice.Service,
	bootstrapResult *bootstrap.Result,
	cvGenSvc cv.CVGenerationService,
	cvExportSvc *cv.ExportService,
	log *logger.Logger,
) *cvMockRegistrar {
	delegate := app.NewDefaultIntentRegistrar(&app.RegistrarConfig{
		CLIService:      cliService,
		CareerService:   svc,
		Log:             log,
		CVGenService:    bootstrapResult.Services.CVGenService,
		CVExportService: bootstrapResult.Services.CVExportService,
	})

	return &cvMockRegistrar{
		delegate:    delegate,
		cliService:  cliService,
		svc:         svc,
		cvGenSvc:    cvGenSvc,
		cvExportSvc: cvExportSvc,
		log:         log,
	}
}

// RegisterAll registers all intents, overriding generate_cv with mock services.
//
// Expected:
//   - router must be a valid DefaultIntentRouter.
//
// Returns:
//   - An error if registration fails.
//
// Side effects:
//   - Registers intents with the router.
func (r *cvMockRegistrar) RegisterAll(ctx context.Context, router *intents.DefaultIntentRouter) error {
	if err := r.delegate.RegisterAll(ctx, router); err != nil {
		return err
	}

	return r.overrideGenerateCV(ctx, router)
}

func (r *cvMockRegistrar) overrideGenerateCV(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("generate_cv", func() intents.Intent {
		events, err := r.svc.GetEventRepository().List(ctx, careerrepo.EventListFilters{Limit: 100})
		if err != nil {
			r.log.Error("Failed to load events for CV generation: %v", err)
			events = []*career.Event{}
		}
		facts, err := r.svc.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 100})
		if err != nil {
			r.log.Error("Failed to load facts for CV generation: %v", err)
			facts = []*career.Fact{}
		}

		var profileCfg *config.ProfileConfig
		var scoringCfg *config.ScoringConfig
		if appCfg, loadErr := config.LoadConfig(); loadErr == nil {
			profileCfg = &appCfg.Profile
			scoringCfg = &appCfg.Scoring
		}

		cvCtx := &intents.GenerateCVContext{
			Events:                events,
			Facts:                 facts,
			AvailableProfiles:     createDefaultCVProfilesForTest(),
			DefaultProfile:        createDefaultCVProfilesForTest()[0],
			CVGenerationService:   r.cvGenSvc,
			DataProcessingService: cv.NewDataProcessingService(r.log),
			BulletGenerator:       cv.NewBulletGenerator(r.log, scoringCfg),
			ExportService:         r.cvExportSvc,
			SkillRepository:       r.svc.GetSkillRepository(),
			EventRepository:       r.svc.GetEventRepository(),
			ProfileConfig:         profileCfg,
			AppContext:            ctx,
			ReviewScreenFactory: func(cvView *career.CVView) screens.Screen {
				return cvscreens.NewCVReviewScreen(cvView)
			},
			PreviewScreenFactory: func(cvView *career.CVView) screens.Screen {
				return cvscreens.NewCVPreviewScreenWithProfile(cvView, profileCfg)
			},
		}

		intent, intentErr := intents.NewGenerateCVIntent(cvCtx)
		if intentErr != nil {
			r.log.Error("Failed to create GenerateCV intent: %v", intentErr)
			return nil
		}
		return intent
	})
}

// createDefaultCVProfilesForTest mirrors the production createDefaultCVProfiles
// from registration.go, providing the same profile options for test consistency.
func createDefaultCVProfilesForTest() []*intents.CVProfile {
	return []*intents.CVProfile{
		{
			ID:             "profile-staff-engineer",
			Name:           "Staff Engineer",
			TargetRole:     "staff",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for staff engineering roles",
		},
		{
			ID:             "profile-principal-engineer",
			Name:           "Principal Engineer",
			TargetRole:     "principal",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for principal/architect roles",
		},
		{
			ID:             "profile-engineering-manager",
			Name:           "Engineering Manager",
			TargetRole:     "em",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for engineering management roles",
		},
		{
			ID:             "profile-senior-engineer",
			Name:           "Senior Engineer",
			TargetRole:     "senior_ic",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for senior individual contributor roles",
		},
	}
}
