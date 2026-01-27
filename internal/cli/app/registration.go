package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	browse_timeline "github.com/baphled/kariya/internal/cli/intents/browse_timeline"
	factmanagement "github.com/baphled/kariya/internal/cli/intents/factmanagement"
	"github.com/baphled/kariya/internal/cli/screens"
	cvscreens "github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
)

// registrationConfig bundles the services needed for intent registration.
type registrationConfig struct {
	router          *intents.DefaultIntentRouter
	cliService      *service.CLIEventService
	careerService   *careerservice.Service
	log             *logger.Logger
	cvGenService    cv.CVGenerationService
	cvExportService *cv.ExportService
}

// createDefaultCVProfiles creates a set of default CV profiles for the GenerateCV intent.
func createDefaultCVProfiles() []*intents.CVProfile {
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

// registerAllIntents registers all intents with the router.
func registerAllIntents(
	ctx context.Context,
	cfg *registrationConfig,
) {
	registerCaptureEventIntent(cfg)
	registerBrowseTimelineIntent(ctx, cfg)
	registerManageSkillsIntent(ctx, cfg)
	registerGenerateCVIntent(ctx, cfg)
	registerConfigureSystemIntent(ctx, cfg)
	registerBurstManagementIntent(ctx, cfg)
	registerFactManagementIntent(ctx, cfg)
}

func registerCaptureEventIntent(cfg *registrationConfig) {
	err := cfg.router.RegisterIntent("capture_event", func() intents.Intent {
		captureCtx := &intents.CaptureEventContext{
			CaptureStrategy: "manual",
			Metadata:        make(map[string]string),
			CLIEventService: cfg.cliService,
			CareerService:   cfg.careerService,
		}
		intent, err := intents.NewCaptureEventIntent(captureCtx)
		if err != nil {
			cfg.log.Error("Failed to create CaptureEvent intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		cfg.log.Error("Failed to register capture_event intent: %v", err)
	}
}

func registerBrowseTimelineIntent(ctx context.Context, cfg *registrationConfig) {
	err := cfg.router.RegisterIntent("browse_timeline", func() intents.Intent {
		events, err := cfg.careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{
			Limit:     1000,
			SortBy:    "date",
			SortOrder: "desc",
		})
		if err != nil {
			cfg.log.Error("Failed to load events: %v", err)
			events = make([]*career.CareerEvent, 0)
		}
		browserCtx := &browse_timeline.IntentContext{
			Events:          events,
			CLIEventService: cfg.cliService,
		}
		intent, err := browse_timeline.NewIntent(browserCtx)
		if err != nil {
			cfg.log.Error("Failed to create BrowseTimeline intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		cfg.log.Error("Failed to register browse_timeline intent: %v", err)
	}
}

func registerManageSkillsIntent(ctx context.Context, cfg *registrationConfig) {
	err := cfg.router.RegisterIntent("manage_skills", func() intents.Intent {
		skillsCtx := &intents.ManageSkillsContext{
			Ctx:             ctx,
			SkillRepository: cfg.careerService.GetSkillRepository(),
			Service:         cfg.careerService,
		}
		return intents.NewManageSkillsIntent(skillsCtx)
	})
	if err != nil {
		cfg.log.Error("Failed to register manage_skills intent: %v", err)
	}
}

func registerGenerateCVIntent(ctx context.Context, cfg *registrationConfig) {
	err := cfg.router.RegisterIntent("generate_cv", func() intents.Intent {
		events, err := cfg.careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{Limit: 100})
		if err != nil {
			cfg.log.Error("Failed to load events for CV generation: %v", err)
			events = []*career.CareerEvent{}
		}
		facts, err := cfg.careerService.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 100})
		if err != nil {
			cfg.log.Error("Failed to load facts for CV generation: %v", err)
			facts = []*career.Fact{}
		}

		// Load user's profile and scoring config for CV generation.
		var profileCfg *config.ProfileConfig
		var scoringCfg *config.ScoringConfig
		if appCfg, err := config.LoadConfig(); err == nil {
			profileCfg = &appCfg.Profile
			scoringCfg = &appCfg.Scoring
		}

		cvCtx := &intents.GenerateCVContext{
			Events:                events,
			Facts:                 facts,
			AvailableProfiles:     createDefaultCVProfiles(),
			DefaultProfile:        createDefaultCVProfiles()[0],
			CVGenerationService:   cfg.cvGenService,
			DataProcessingService: cv.NewDataProcessingService(cfg.log),
			BulletGenerator:       cv.NewBulletGenerator(cfg.log, scoringCfg),
			ExportService:         cfg.cvExportService,
			ProfileConfig:         profileCfg,
			AppContext:            ctx,
			ReviewScreenFactory: func(cvView *career.CVView) screens.Screen {
				return cvscreens.NewCVReviewScreen(cvView)
			},
			PreviewScreenFactory: func(cvView *career.CVView) screens.Screen {
				return cvscreens.NewCVPreviewScreenWithProfile(cvView, profileCfg)
			},
		}
		intent, err := intents.NewGenerateCVIntent(cvCtx)
		if err != nil {
			cfg.log.Error("Failed to create GenerateCV intent: %v", err)
			return nil
		}
		// Enable wizard flow by default (Phase 7 - Full Integration).
		intent.EnableWizardFlow()
		return intent
	})
	if err != nil {
		cfg.log.Error("Failed to register generate_cv intent: %v", err)
	}
}

func registerConfigureSystemIntent(ctx context.Context, cfg *registrationConfig) {
	err := cfg.router.RegisterIntent("configure_system", func() intents.Intent {
		intent, err := intents.NewConfigureSystemIntent(ctx)
		if err != nil {
			cfg.log.Error("Failed to create ConfigureSystem intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		cfg.log.Error("Failed to register configure_system intent: %v", err)
	}
}

func registerBurstManagementIntent(ctx context.Context, cfg *registrationConfig) {
	err := cfg.router.RegisterIntent("burst_management", func() intents.Intent {
		burstRepo := cfg.careerService.GetBurstRepository()
		burstCtx := intents.NewBurstManagementContext(cfg.careerService, burstRepo, ctx)
		if burstCtx == nil {
			cfg.log.Error("Failed to create BurstManagement context")
			return nil
		}
		intent, err := intents.NewBurstManagementIntent(burstCtx)
		if err != nil || intent == nil {
			cfg.log.Error("Failed to create BurstManagement intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		cfg.log.Error("Failed to register burst_management intent: %v", err)
	}
}

func registerFactManagementIntent(ctx context.Context, cfg *registrationConfig) {
	err := cfg.router.RegisterIntent("fact_management", func() intents.Intent {
		factRepo := cfg.careerService.GetFactRepository()
		factCtx := factmanagement.NewIntentContext(ctx, factRepo)
		if factCtx == nil {
			cfg.log.Error("Failed to create FactManagement context")
			return nil
		}
		intent, err := factmanagement.NewIntent(factCtx)
		if err != nil {
			cfg.log.Error("Failed to create FactManagement intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		cfg.log.Error("Failed to register fact_management intent: %v", err)
	}
}
