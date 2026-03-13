// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

import (
	"github.com/baphled/kariya/internal/tui/intents"
	factviews "github.com/baphled/kariya/internal/tui/views/fact"
	"github.com/baphled/kariya/internal/ui/behaviors"
	"github.com/baphled/kariya/internal/ui/display"
)

// Intent orchestrates the fact management workflow.
type Intent struct {
	*intents.BaseIntent

	// Context holds input parameters and business logic.
	context *IntentValidator

	// State machine fields.
	state  State
	active bool
	result *intents.IntentResult[*Result]

	// Table behavior for list navigation.
	tableBehavior *behaviors.TableBehavior[display.Fact]

	// Edit modal for creating/editing facts.
	editModal *factviews.EditFact
}
