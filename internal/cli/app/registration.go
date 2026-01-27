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
// All RegisterIntent calls below use nolint:errcheck because RegisterIntent only returns
// an error on duplicate registration, which cannot happen in this initialization code.
func registerAllIntents(
	router *intents.DefaultIntentRouter,
	cliService *service.CLIEventService,
	careerService *careerservice.Service,
	log *logger.Logger,
	ctx context.Context,
	cvGenService cv.CVGenerationService,
	cvExportService *cv.ExportService,
) {
	registerCaptureEventIntent(router, cliService, careerService, log)
	registerBrowseTimelineIntent(router, cliService, careerService, log, ctx)
	registerManageSkillsIntent(router, careerService, log, ctx)
	registerGenerateCVIntent(router, careerService, log, ctx, cvGenService, cvExportService)
	registerConfigureSystemIntent(router, log, ctx)
	registerBurstManagementIntent(router, careerService, log, ctx)
	registerFactManagementIntent(router, careerService, log, ctx)
}

func registerCaptureEventIntent(
	router *intents.DefaultIntentRouter,
	cliService *service.CLIEventService,
	careerService *careerservice.Service,
	log *logger.Logger,
) {
	err := router.RegisterIntent("capture_event", func() intents.Intent {
		captureCtx := &intents.CaptureEventContext{
			CaptureStrategy: "manual",
			Metadata:        make(map[string]string),
			CLIEventService: cliService,
			CareerService:   careerService,
		}
		intent, err := intents.NewCaptureEventIntent(captureCtx)
		if err != nil {
			log.Error("Failed to create CaptureEvent intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		log.Error("Failed to register capture_event intent: %v", err)
	}
}

func registerBrowseTimelineIntent(
	router *intents.DefaultIntentRouter,
	cliService *service.CLIEventService,
	careerService *careerservice.Service,
	log *logger.Logger,
	ctx context.Context,
) {
	err := router.RegisterIntent("browse_timeline", func() intents.Intent {
		events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{
			Limit:     1000,
			SortBy:    "date",
			SortOrder: "desc",
		})
		if err != nil {
			log.Error("Failed to load events: %v", err)
			events = make([]*career.CareerEvent, 0)
		}
		browserCtx := &browse_timeline.IntentContext{
			Events:          events,
			CLIEventService: cliService,
		}
		intent, err := browse_timeline.NewIntent(browserCtx)
		if err != nil {
			log.Error("Failed to create BrowseTimeline intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		log.Error("Failed to register browse_timeline intent: %v", err)
	}
}

func registerManageSkillsIntent(
	router *intents.DefaultIntentRouter,
	careerService *careerservice.Service,
	log *logger.Logger,
	ctx context.Context,
) {
	err := router.RegisterIntent("manage_skills", func() intents.Intent {
		skillsCtx := &intents.ManageSkillsContext{
			Ctx:             ctx,
			SkillRepository: careerService.GetSkillRepository(),
			Service:         careerService,
		}
		return intents.NewManageSkillsIntent(skillsCtx)
	})
	if err != nil {
		log.Error("Failed to register manage_skills intent: %v", err)
	}
}

func registerGenerateCVIntent(
	router *intents.DefaultIntentRouter,
	careerService *careerservice.Service,
	log *logger.Logger,
	ctx context.Context,
	cvGenService cv.CVGenerationService,
	cvExportService *cv.ExportService,
) {
	// BUG-004: Removed stub data fallback - empty state is now handled by showing
	// an info modal in handleMenuInput before this intent is activated.
	err := router.RegisterIntent("generate_cv", func() intents.Intent {
		events, err := careerService.GetEventRepository().List(ctx, careerrepo.ListFilters{Limit: 100})
		if err != nil {
			log.Error("Failed to load events for CV generation: %v", err)
			events = []*career.CareerEvent{}
		}
		facts, err := careerService.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 100})
		if err != nil {
			log.Error("Failed to load facts for CV generation: %v", err)
			facts = []*career.Fact{}
		}

		// Load user's profile and scoring config for CV generation.
		var profileCfg *config.ProfileConfig
		var scoringCfg *config.ScoringConfig
		if cfg, err := config.LoadConfig(); err == nil {
			profileCfg = &cfg.Profile
			scoringCfg = &cfg.Scoring
		}

		cvCtx := &intents.GenerateCVContext{
			Events:                events,
			Facts:                 facts,
			AvailableProfiles:     createDefaultCVProfiles(),
			DefaultProfile:        createDefaultCVProfiles()[0],
			CVGenerationService:   cvGenService,
			DataProcessingService: cv.NewDataProcessingService(log),
			BulletGenerator:       cv.NewBulletGenerator(log, scoringCfg),
			ExportService:         cvExportService,
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
			log.Error("Failed to create GenerateCV intent: %v", err)
			return nil
		}
		// Enable wizard flow by default (Phase 7 - Full Integration).
		intent.EnableWizardFlow()
		return intent
	})
	if err != nil {
		log.Error("Failed to register generate_cv intent: %v", err)
	}
}

func registerConfigureSystemIntent(
	router *intents.DefaultIntentRouter,
	log *logger.Logger,
	ctx context.Context,
) {
	err := router.RegisterIntent("configure_system", func() intents.Intent {
		intent, err := intents.NewConfigureSystemIntent(ctx)
		if err != nil {
			log.Error("Failed to create ConfigureSystem intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		log.Error("Failed to register configure_system intent: %v", err)
	}
}

func registerBurstManagementIntent(
	router *intents.DefaultIntentRouter,
	careerService *careerservice.Service,
	log *logger.Logger,
	ctx context.Context,
) {
	err := router.RegisterIntent("burst_management", func() intents.Intent {
		burstRepo := careerService.GetBurstRepository()
		burstCtx := intents.NewBurstManagementContext(careerService, burstRepo, ctx)
		if burstCtx == nil {
			log.Error("Failed to create BurstManagement context")
			return nil
		}
		intent, err := intents.NewBurstManagementIntent(burstCtx)
		if err != nil || intent == nil {
			log.Error("Failed to create BurstManagement intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		log.Error("Failed to register burst_management intent: %v", err)
	}
}

func registerFactManagementIntent(
	router *intents.DefaultIntentRouter,
	careerService *careerservice.Service,
	log *logger.Logger,
	ctx context.Context,
) {
	err := router.RegisterIntent("fact_management", func() intents.Intent {
		factRepo := careerService.GetFactRepository()
		factCtx := factmanagement.NewIntentContext(ctx, factRepo)
		if factCtx == nil {
			log.Error("Failed to create FactManagement context")
			return nil
		}
		intent, err := factmanagement.NewIntent(factCtx)
		if err != nil {
			log.Error("Failed to create FactManagement intent: %v", err)
			return nil
		}
		return intent
	})
	if err != nil {
		log.Error("Failed to register fact_management intent: %v", err)
	}
}
