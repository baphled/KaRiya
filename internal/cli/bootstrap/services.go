package bootstrap

import (
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
)

// Services holds all initialized services needed by the application.
type Services struct {
	ConfigManager   cv.ConfigManager
	CVGenService    cv.CVGenerationService
	CVExportService *cv.ExportService
}

// InitServices initializes all CV-related services.
//
// Expected:
//   - service must be valid.
//   - config must be a valid configuration object.
//   - logger must be valid.
//
// Returns:
//   - A fully initialized Services ready for use.
//
// Side effects:
//   - None.
func InitServices(careerService *careerservice.Service, cfg *config.Config, log *logger.Logger) *Services {
	configMgr := initConfigManager(log)
	cvGenService := initCVGenerationService(careerService, configMgr, &cfg.Scoring, log)
	cvExportService := cv.NewExportServiceWithDeps(log, &cfg.Profile, careerService.GetSkillRepository())

	return &Services{
		ConfigManager:   configMgr,
		CVGenService:    cvGenService,
		CVExportService: cvExportService,
	}
}

// initConfigManager initializes the CV config manager.
// Falls back to in-memory storage if YAML config fails to load.
func initConfigManager(log *logger.Logger) cv.ConfigManager {
	yamlMgr, err := cv.NewYAMLConfigManager(log)
	if err != nil {
		log.Error("Failed to initialize CV config manager: %v", err)
		return cv.NewMemoryConfigManager()
	}
	return yamlMgr
}

// initCVGenerationService initializes the CV generation service with all dependencies.
func initCVGenerationService(
	careerService *careerservice.Service,
	configMgr cv.ConfigManager,
	scoringCfg *config.ScoringConfig,
	log *logger.Logger,
) cv.CVGenerationService {
	bulletGenerator := cv.NewBulletGenerator(log, scoringCfg)
	dataProcessor := cv.NewDataProcessingService(log)
	sectionBuilder := cv.NewSectionBuilder(
		careerService.GetSkillRepository(),
		log,
	)
	return cv.NewCVGenerationService(
		careerService.GetEventRepository(),
		careerService.GetFactRepository(),
		configMgr,
		bulletGenerator,
		dataProcessor,
		sectionBuilder,
		log,
	)
}
