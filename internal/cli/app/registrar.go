package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/intents/browsetimeline"
	burstmanagement "github.com/baphled/kariya/internal/cli/intents/burst_management"
	"github.com/baphled/kariya/internal/cli/intents/factmanagement"
	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
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

// IntentRegistrar defines the interface for registering intents with the router.
// This abstraction enables dependency injection for testing.
type IntentRegistrar interface {
	RegisterAll(ctx context.Context, router *intents.DefaultIntentRouter) error
}

// RegistrarConfig contains dependencies needed for intent registration.
type RegistrarConfig struct {
	CLIService      *service.CLIEventService
	CareerService   *careerservice.Service
	Log             *logger.Logger
	CVGenService    cv.CVGenerationService
	CVExportService *cv.ExportService
}

// DefaultIntentRegistrar implements IntentRegistrar with production logic.
type DefaultIntentRegistrar struct {
	config *RegistrarConfig
}

// NewDefaultIntentRegistrar creates a new default intent registrar.
func NewDefaultIntentRegistrar(cfg *RegistrarConfig) *DefaultIntentRegistrar {
	return &DefaultIntentRegistrar{config: cfg}
}

// RegisterAll registers all intents with the router.
func (r *DefaultIntentRegistrar) RegisterAll(ctx context.Context, router *intents.DefaultIntentRouter) error {
	var firstErr error
	recordErr := func(err error) {
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}

	recordErr(r.registerCaptureEvent(router))
	recordErr(r.registerBrowseTimeline(ctx, router))
	recordErr(r.registerManageSkills(ctx, router))
	recordErr(r.registerGenerateCV(ctx, router))
	recordErr(r.registerConfigureSystem(ctx, router))
	recordErr(r.registerBurstManagement(ctx, router))
	recordErr(r.registerFactManagement(ctx, router))

	return firstErr
}

func (r *DefaultIntentRegistrar) registerCaptureEvent(router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("capture_event", func() intents.Intent {
		if r.config.CLIService == nil || r.config.CareerService == nil {
			r.config.Log.Error("Failed to create CaptureEvent intent: missing required services")
			return nil
		}
		captureCtx := &intents.CaptureEventContext{
			CaptureStrategy: "manual",
			Metadata:        make(map[string]string),
			CLIEventService: r.config.CLIService,
			CareerService:   r.config.CareerService,
		}
		intent, err := intents.NewCaptureEventIntent(captureCtx)
		if err != nil {
			r.config.Log.Error("Failed to create CaptureEvent intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegistrar) registerBrowseTimeline(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("browse_timeline", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create BrowseTimeline intent: missing CareerService")
			return nil
		}
		events, err := r.config.CareerService.GetEventRepository().List(ctx, careerrepo.EventListFilters{
			Limit:     1000,
			SortBy:    "date",
			SortOrder: "desc",
		})
		if err != nil {
			r.config.Log.Error("Failed to load events: %v", err)
			events = make([]*career.Event, 0)
		}
		browserCtx := &browsetimeline.IntentContext{
			Events:          events,
			CLIEventService: r.config.CLIService,
		}
		intent, err := browsetimeline.NewIntent(browserCtx)
		if err != nil {
			r.config.Log.Error("Failed to create BrowseTimeline intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegistrar) registerManageSkills(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("manage_skills", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create ManageSkills intent: missing CareerService")
			return nil
		}
		skillsCtx := skillsmanagement.NewIntentContext(
			ctx,
			r.config.CareerService.GetSkillRepository(),
		)
		intent, err := skillsmanagement.NewIntent(skillsCtx)
		if err != nil {
			r.config.Log.Error("Failed to create ManageSkills intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegistrar) registerGenerateCV(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("generate_cv", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create GenerateCV intent: missing CareerService")
			return nil
		}
		events, err := r.config.CareerService.GetEventRepository().List(ctx, careerrepo.EventListFilters{Limit: 100})
		if err != nil {
			r.config.Log.Error("Failed to load events for CV generation: %v", err)
			events = []*career.Event{}
		}
		facts, err := r.config.CareerService.GetFactRepository().List(ctx, careerrepo.FactListFilters{Limit: 100})
		if err != nil {
			r.config.Log.Error("Failed to load facts for CV generation: %v", err)
			facts = []*career.Fact{}
		}

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
			CVGenerationService:   r.config.CVGenService,
			DataProcessingService: cv.NewDataProcessingService(r.config.Log),
			BulletGenerator:       cv.NewBulletGenerator(r.config.Log, scoringCfg),
			ExportService:         r.config.CVExportService,
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
			r.config.Log.Error("Failed to create GenerateCV intent: %v", err)
			return nil
		}
		intent.EnableWizardFlow()
		return intent
	})
}

func (r *DefaultIntentRegistrar) registerConfigureSystem(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("configure_system", func() intents.Intent {
		intent, err := intents.NewConfigureSystemIntent(ctx)
		if err != nil {
			r.config.Log.Error("Failed to create ConfigureSystem intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegistrar) registerBurstManagement(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("burst_management", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create BurstManagement intent: missing CareerService")
			return nil
		}
		burstRepo := r.config.CareerService.GetBurstRepository()
		burstCtx := &burstmanagement.IntentContext{
			Service:         r.config.CareerService,
			BurstRepository: burstRepo,
			Context:         ctx,
		}
		intent, err := burstmanagement.NewIntent(burstCtx)
		if err != nil || intent == nil {
			r.config.Log.Error("Failed to create BurstManagement intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegistrar) registerFactManagement(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("fact_management", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create FactManagement intent: missing CareerService")
			return nil
		}
		factRepo := r.config.CareerService.GetFactRepository()
		factCtx := factmanagement.NewIntentContext(ctx, factRepo)
		if factCtx == nil {
			r.config.Log.Error("Failed to create FactManagement context")
			return nil
		}
		intent, err := factmanagement.NewIntent(factCtx)
		if err != nil {
			r.config.Log.Error("Failed to create FactManagement intent: %v", err)
			return nil
		}
		return intent
	})
}
