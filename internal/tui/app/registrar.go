package app

import (
	"context"

	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/tui/intents"
	"github.com/baphled/kariya/internal/tui/intents/browsetimeline"
	burstmanagement "github.com/baphled/kariya/internal/tui/intents/burst_management"
	"github.com/baphled/kariya/internal/tui/intents/captureevent"
	"github.com/baphled/kariya/internal/tui/intents/configure"
	"github.com/baphled/kariya/internal/tui/intents/factmanagement"
	"github.com/baphled/kariya/internal/tui/intents/generatecv"
	"github.com/baphled/kariya/internal/tui/intents/skillsmanagement"
)

// IntentRegisterer defines the interface for registering intents with the router.
// This abstraction enables dependency injection for testing.
type IntentRegisterer interface {
	RegisterAll(ctx context.Context, router *intents.DefaultIntentRouter) error
}

// RegistrarConfig contains dependencies needed for intent registration.
type RegistrarConfig struct {
	CLIService            *service.CLIEventService
	CareerService         *careerservice.Service
	SkillInferenceService skillinference.SkillInferenceService
	Log                   *logger.Logger
	CVGenService          cv.CVGenerationService
	CVExportService       *cv.ExportService
}

// DefaultIntentRegisterer implements IntentRegisterer with production logic.
type DefaultIntentRegisterer struct {
	config *RegistrarConfig
}

// NewDefaultIntentRegisterer creates a new default intent registrar.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A fully initialized DefaultIntentRegisterer ready for use.
//
// Side effects:
//   - None.
func NewDefaultIntentRegisterer(cfg *RegistrarConfig) *DefaultIntentRegisterer {
	return &DefaultIntentRegisterer{config: cfg}
}

// RegisterAll registers all intents with the router.
//
// Expected:
//   - defaultintentrouter must be valid.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (r *DefaultIntentRegisterer) RegisterAll(ctx context.Context, router *intents.DefaultIntentRouter) error {
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

func (r *DefaultIntentRegisterer) registerCaptureEvent(router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("capture_event", func() intents.Intent {
		if r.config.CLIService == nil || r.config.CareerService == nil {
			r.config.Log.Error("Failed to create CaptureEvent intent: missing required services")
			return nil
		}
		captureCtx := &captureevent.IntentValidator{
			CaptureStrategy:       "manual",
			Metadata:              make(map[string]string),
			CLIEventService:       r.config.CLIService,
			CareerService:         r.config.CareerService,
			SkillInferenceService: r.config.SkillInferenceService,
		}
		intent, err := captureevent.NewIntent(captureCtx)
		if err != nil {
			r.config.Log.Error("Failed to create CaptureEvent intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegisterer) registerBrowseTimeline(ctx context.Context, router *intents.DefaultIntentRouter) error {
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
		browserCtx := &browsetimeline.IntentValidator{
			Events:                events,
			CLIEventService:       r.config.CLIService,
			CLISkillCreator:       service.NewCLISkillCreator(r.config.CareerService.GetSkillRepository()),
			SkillInferenceService: r.config.SkillInferenceService,
		}
		intent, err := browsetimeline.NewIntent(browserCtx)
		if err != nil {
			r.config.Log.Error("Failed to create BrowseTimeline intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegisterer) registerManageSkills(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("manage_skills", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create ManageSkills intent: missing CareerService")
			return nil
		}
		skillsCtx := skillsmanagement.NewIntentValidator(
			ctx,
			r.config.CareerService.GetSkillRepository(),
		)
		skillsCtx.EventRepository = r.config.CareerService.GetEventRepository()
		if r.config.SkillInferenceService != nil {
			skillsCtx.SkillInferenceService = r.config.SkillInferenceService
		}
		intent, err := skillsmanagement.NewIntent(skillsCtx)
		if err != nil {
			r.config.Log.Error("Failed to create ManageSkills intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegisterer) registerGenerateCV(ctx context.Context, router *intents.DefaultIntentRouter) error {
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
		appCfg, err := config.LoadConfig()
		if err != nil {
			r.config.Log.Error("Failed to load config, using defaults: %v", err)
			appCfg = config.DefaultConfig()
		}
		config.MigrateProfileConfig(appCfg)
		profileCfg = &appCfg.Profile
		scoringCfg = &appCfg.Scoring

		cvCtx := &generatecv.IntentValidator{
			Events:                events,
			Facts:                 facts,
			AvailableProfiles:     createDefaultCVProfiles(),
			DefaultProfile:        createDefaultCVProfiles()[0],
			CVGenerationService:   r.config.CVGenService,
			DataProcessingService: cv.NewDataProcessingService(r.config.Log),
			BulletGenerator:       cv.NewBulletGenerator(r.config.Log, scoringCfg),
			ExportService:         r.config.CVExportService,
			SkillRepository:       r.config.CareerService.GetSkillRepository(),
			EventRepository:       r.config.CareerService.GetEventRepository(),
			ProfileConfig:         profileCfg,
		}
		intent, err := generatecv.NewIntent(cvCtx)
		if err != nil {
			r.config.Log.Error("Failed to create GenerateCV intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegisterer) registerConfigureSystem(_ context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("configure_system", func() intents.Intent {
		cfg, err := config.LoadConfig()
		if err != nil {
			cfg = config.DefaultConfig()
		}
		intentCtx := &configure.IntentValidator{
			Cfg:      cfg,
			Settings: configure.SettingsFromConfig(cfg),
		}
		intent, err := configure.NewIntent(intentCtx)
		if err != nil {
			r.config.Log.Error("Failed to create ConfigureSystem intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegisterer) registerBurstManagement(_ context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("burst_management", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create BurstManagement intent: missing CareerService")
			return nil
		}
		burstRepo := r.config.CareerService.GetBurstRepository()
		burstCtx := &burstmanagement.IntentValidator{
			Service:               r.config.CareerService,
			SkillInferenceService: r.config.SkillInferenceService,
			BurstRepository:       burstRepo,
			SkillRepository:       r.config.CareerService.GetSkillRepository(),
		}
		intent, err := burstmanagement.NewIntent(burstCtx)
		if err != nil || intent == nil {
			r.config.Log.Error("Failed to create BurstManagement intent: %v", err)
			return nil
		}
		return intent
	})
}

func (r *DefaultIntentRegisterer) registerFactManagement(ctx context.Context, router *intents.DefaultIntentRouter) error {
	return router.RegisterIntent("fact_management", func() intents.Intent {
		if r.config.CareerService == nil {
			r.config.Log.Error("Failed to create FactManagement intent: missing CareerService")
			return nil
		}
		factRepo := r.config.CareerService.GetFactRepository()
		factCtx := factmanagement.NewIntentValidator(ctx, factRepo)
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
