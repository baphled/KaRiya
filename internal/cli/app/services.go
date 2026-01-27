package app

import (
	"github.com/baphled/kariya/internal/config"
	"github.com/baphled/kariya/internal/logger"
	careerservice "github.com/baphled/kariya/internal/service/career"
	cv "github.com/baphled/kariya/internal/service/career/cv"
)

// initConfigManager initializes the CV config manager.
// Falls back to in-memory storage if YAML config fails to load.
func initConfigManager(log *logger.Logger) cv.ConfigManager {
	var configMgr cv.ConfigManager
	yamlMgr, err := cv.NewYAMLConfigManager(log)
	if err != nil {
		log.Error("Failed to initialize CV config manager: %v", err)
		configMgr = cv.NewMemoryConfigManager()
	} else {
		configMgr = yamlMgr
	}
	return configMgr
}

// initCVGenerationService initializes the CV generation service with all dependencies.
func initCVGenerationService(
	careerService *careerservice.Service,
	configMgr cv.ConfigManager,
	scoringCfg *config.ScoringConfig,
	log *logger.Logger,
) cv.CVGenerationService {
	// BUG-008: Use BulletGenerator for role-based scoring with config.
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
