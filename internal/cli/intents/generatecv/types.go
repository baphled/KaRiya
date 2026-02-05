package generatecv

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	cvmodals "github.com/baphled/kariya/internal/cli/screens/cv/modals"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/service/career/cv"
)

// Intent orchestrates the multi-step CV generation workflow, managing
// state transitions between configuration, technology extraction, generation,
// review, preview, and export phases.
//
// Side effects:
//   - None at construction; state mutations occur through Update and Init.
type Intent struct {
	*intents.BaseIntent

	context *IntentContext
	state   *model
	active  bool
	result  *intents.IntentResult[*Result]
	logger  *logger.Logger

	wizardModal   *cvmodals.ConfigWizardModal
	progressModal *cvmodals.ProgressModal
	exportModal   *cvmodals.ExportModal

	wizardReviewScreen  screens.Screen
	wizardPreviewScreen screens.Screen
}

// model tracks the internal mutable state of the GenerateCV intent throughout
// the wizard workflow: selected profile, audience, technology/focus preferences,
// skills configuration, generated CV, and export progress.
type model struct {
	context      *IntentContext
	currentState State

	selectedProfile  *CVProfile
	selectedAudience string
	generatedCV      *career.CVView

	extractedTechnologies []*ExtractedTechnology
	focusAreaSuggestion   *FocusAreaSuggestion

	selectedTechnologyFocus cv.TechnologyFocus
	selectedTechnologies    []string
	selectedFocusArea       cv.FocusArea

	selectedSkillsFormat string
	selectedSkillsLimit  int
	selectedCVLength     string

	selectedExportFormat ExportFormat
	selectedExportOption ExportOption
	exportedPath         string
	exportError          error
	isExporting          bool
}
