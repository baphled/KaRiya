// Package fact_management implements the FactManagement intent for managing career facts.
package fact_management

import (
	"github.com/baphled/kariya/internal/cli/behaviors"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/fact"
	factmodals "github.com/baphled/kariya/internal/cli/screens/fact/modals"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	domain "github.com/baphled/kariya/internal/domain/career"
)

// Ensure Intent implements ScreenResultHandler interface.
var _ intents.ScreenResultHandler = (*Intent)(nil)

// Intent orchestrates the fact management workflow.
type Intent struct {
	*intents.BaseIntent

	// Context holds input parameters and business logic.
	context *IntentContext

	// State machine fields.
	state  State
	active bool
	result *intents.IntentResult[*Result]

	// Table behavior for list navigation.
	tableBehavior *behaviors.TableBehavior[*domain.Fact]

	// Screens (one per major state).
	listScreen *fact.ListScreen

	// Active screen pointer.
	activeScreen screens.Screen

	// Modals (overlays).
	detailModal *factmodals.DetailModal
	editModal   *factmodals.EditModal
	deleteModal *feedback.ConfirmModal
	errorModal  *feedback.Modal

	// Modal registry for unified handling.
	modalRegistry *intents.ModalRegistry
}
