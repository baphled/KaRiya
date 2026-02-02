// Package factmanagement implements the FactManagement intent for managing career facts.
package factmanagement

import (
	"context"
	"errors"
	"time"

	domain "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

var (
	// ErrFactNotFound indicates the requested fact was not found.
	ErrFactNotFound = errors.New("fact not found")

	// ErrServiceNotAvailable indicates the fact repository is not available.
	ErrServiceNotAvailable = errors.New("fact service not available")

	// ErrInvalidState indicates an invalid state for the current operation.
	ErrInvalidState = errors.New("invalid state")
)

// IntentContext holds input parameters and business logic for the FactManagement intent.
type IntentContext struct {
	// Facts is the list of facts loaded from the repository.
	Facts []*domain.Fact

	// SelectedFact is the currently selected fact.
	SelectedFact *domain.Fact

	// SelectedFactIndex is the index of the selected fact.
	SelectedFactIndex int

	// SearchText is the current search filter.
	SearchText string

	// MinQuality is the minimum quality filter.
	MinQuality float64

	// MaxQuality is the maximum quality filter.
	MaxQuality float64

	// FilterSource filters by source type.
	FilterSource string

	// SortBy is the current sort field.
	SortBy string

	// SortOrder is the current sort direction.
	SortOrder string

	// CurrentPage is the current page number.
	CurrentPage int

	// PageSize is the number of items per page.
	PageSize int

	// TotalFacts is the total count of facts.
	TotalFacts int

	// EditingFact is the fact being edited (copy for editing).
	EditingFact *domain.Fact

	// FormErrors contains validation errors by field.
	FormErrors map[string]string

	// IsNewFact indicates whether we are creating a new fact.
	IsNewFact bool

	// FactToDelete is the fact pending deletion.
	FactToDelete *domain.Fact

	// ScrollPosition tracks the scroll offset.
	ScrollPosition int

	// ExpandedRows tracks which rows are expanded.
	ExpandedRows map[int]bool

	// FactRepository is the repository for CRUD operations.
	FactRepository careerrepo.FactRepository

	// Context is the context for repository operations.
	Context context.Context
}

// NewIntentContext initialises a FactManagement context with sensible defaults for pagination, sorting, and filtering.
//
// Expected: ctx must be non-nil. factRepo may be nil for offline/test usage, but repository operations will fail.
//
// Returns: a fully initialised IntentContext ready for use by the intent.
//
// Side effects: None.
func NewIntentContext(ctx context.Context, factRepo careerrepo.FactRepository) *IntentContext {
	return &IntentContext{
		Facts:             make([]*domain.Fact, 0),
		SelectedFactIndex: -1,
		SearchText:        "",
		MinQuality:        0.0,
		MaxQuality:        1.0,
		FilterSource:      "",
		SortBy:            "date",
		SortOrder:         "desc",
		CurrentPage:       0,
		PageSize:          20,
		FormErrors:        make(map[string]string),
		ExpandedRows:      make(map[int]bool),
		FactRepository:    factRepo,
		Context:           ctx,
	}
}

// Validate ensures all required collection fields are initialised, replacing nil maps and slices with empty defaults.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (c *IntentContext) Validate() error {
	if c.Facts == nil {
		c.Facts = make([]*domain.Fact, 0)
	}
	if c.FormErrors == nil {
		c.FormErrors = make(map[string]string)
	}
	if c.ExpandedRows == nil {
		c.ExpandedRows = make(map[int]bool)
	}
	return nil
}

// LoadFacts fetches all facts from the repository and resets pagination and selection state.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (c *IntentContext) LoadFacts() error {
	if c.FactRepository == nil {
		return ErrServiceNotAvailable
	}
	facts, err := c.FactRepository.List(c.Context, careerrepo.FactListFilters{})
	if err != nil {
		return err
	}
	c.Facts = facts
	c.TotalFacts = len(facts)
	c.CurrentPage = 0
	if len(facts) > 0 {
		c.SelectedFactIndex = 0
		c.SelectedFact = facts[0]
	} else {
		c.SelectedFactIndex = -1
		c.SelectedFact = nil
	}
	return nil
}

// GetPageFacts provides the slice of facts visible on the current page, bounded by PageSize.
//
// Returns:
//   - A []*domain.Fact value.
//
// Side effects:
//   - None.
func (c *IntentContext) GetPageFacts() []*domain.Fact {
	if len(c.Facts) == 0 {
		return make([]*domain.Fact, 0)
	}
	start := c.CurrentPage * c.PageSize
	end := start + c.PageSize
	if start >= len(c.Facts) {
		return make([]*domain.Fact, 0)
	}
	if end > len(c.Facts) {
		end = len(c.Facts)
	}
	return c.Facts[start:end]
}

// SelectFact updates the current selection to the fact at the given page-relative index.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (c *IntentContext) SelectFact(index int) {
	pageFacts := c.GetPageFacts()
	if index >= 0 && index < len(pageFacts) {
		c.SelectedFact = pageFacts[index]
		c.SelectedFactIndex = index
	}
}

// GetSelectedFact provides access to the currently highlighted fact for detail views and actions.
//
// Returns:
//   - A fully initialized domain.Fact ready for use.
//
// Side effects:
//   - None.
func (c *IntentContext) GetSelectedFact() *domain.Fact {
	return c.SelectedFact
}

// ClearFormErrors resets the validation error state, typically called when entering or re-entering a form.
//
// Side effects:
//   - None.
func (c *IntentContext) ClearFormErrors() {
	c.FormErrors = make(map[string]string)
}

// SetFormError records a validation error message against a specific form field for display to the user.
//
// Expected:
//   - Must be a valid string.
//
// Side effects:
//   - None.
func (c *IntentContext) SetFormError(field, message string) {
	c.FormErrors[field] = message
}

// HasFormErrors indicates whether any validation errors exist, used to gate form submission.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func (c *IntentContext) HasFormErrors() bool {
	return len(c.FormErrors) > 0
}

// CreateFact validates and persists a new fact, then appends it to the in-memory collection.
//
// Expected: fact must be non-nil and pass domain validation. FactRepository must be set.
//
// Returns: ErrServiceNotAvailable if FactRepository is nil, or any validation/repository error.
//
// Side effects: On success, appends to Facts and updates TotalFacts. Persists the fact via the repository.
func (c *IntentContext) CreateFact(fact *domain.Fact) error {
	if c.FactRepository == nil {
		return ErrServiceNotAvailable
	}
	if err := fact.Validate(); err != nil {
		return err
	}
	if err := c.FactRepository.Create(c.Context, fact); err != nil {
		return err
	}
	c.Facts = append(c.Facts, fact)
	c.TotalFacts = len(c.Facts)
	return nil
}

// UpdateFact validates and persists changes to an existing fact, then updates the in-memory collection.
//
// Expected: fact must be non-nil with a valid ID matching an existing fact. FactRepository must be set.
//
// Returns: ErrServiceNotAvailable if FactRepository is nil, or any validation/repository error.
//
// Side effects: On success, replaces the matching entry in Facts and persists via the repository.
func (c *IntentContext) UpdateFact(fact *domain.Fact) error {
	if c.FactRepository == nil {
		return ErrServiceNotAvailable
	}
	if err := fact.Validate(); err != nil {
		return err
	}
	if err := c.FactRepository.Update(c.Context, fact); err != nil {
		return err
	}
	for i, f := range c.Facts {
		if f.ID == fact.ID {
			c.Facts[i] = fact
			break
		}
	}
	return nil
}

// DeleteFact removes a fact from both the repository and the in-memory collection, then clears the selection.
//
// Expected: factID must be a non-empty string identifying an existing fact. FactRepository must be set.
//
// Returns: ErrServiceNotAvailable if FactRepository is nil, or any repository error.
//
// Side effects: Removes the fact from Facts, updates TotalFacts, resets SelectedFact and SelectedFactIndex. Deletes via the repository.
func (c *IntentContext) DeleteFact(factID string) error {
	if c.FactRepository == nil {
		return ErrServiceNotAvailable
	}
	if err := c.FactRepository.Delete(c.Context, factID); err != nil {
		return err
	}
	for i, f := range c.Facts {
		if f.ID == factID {
			c.Facts = append(c.Facts[:i], c.Facts[i+1:]...)
			break
		}
	}
	c.TotalFacts = len(c.Facts)
	c.SelectedFact = nil
	c.SelectedFactIndex = -1
	return nil
}

// StartNewFact initialises a blank editing fact with default timestamps for the creation workflow.
//
// Side effects:
//   - None.
func (c *IntentContext) StartNewFact() {
	c.EditingFact = &domain.Fact{
		ID:                   "",
		Text:                 "",
		CompetencyCategories: make([]string, 0),
		AudienceRelevance:    make([]string, 0),
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	c.IsNewFact = true
	c.ClearFormErrors()
}

// StartEditFact creates a deep copy of the given fact for safe editing without mutating the original.
//
// Expected:
//   - fact must be valid.
//
// Side effects:
//   - None.
func (c *IntentContext) StartEditFact(fact *domain.Fact) {
	c.EditingFact = &domain.Fact{
		ID:                   fact.ID,
		Text:                 fact.Text,
		CompetencyCategories: append([]string{}, fact.CompetencyCategories...),
		RoleFit:              fact.RoleFit,
		AudienceRelevance:    append([]string{}, fact.AudienceRelevance...),
		StrengthSignal:       fact.StrengthSignal,
		SourceEventID:        fact.SourceEventID,
		SourceBurstID:        fact.SourceBurstID,
		CreatedAt:            fact.CreatedAt,
		UpdatedAt:            fact.UpdatedAt,
	}
	c.IsNewFact = false
	c.ClearFormErrors()
}

// CancelEdit discards any in-progress fact editing and resets the form state.
//
// Side effects:
//   - None.
func (c *IntentContext) CancelEdit() {
	c.EditingFact = nil
	c.IsNewFact = false
	c.ClearFormErrors()
}

// SaveEdit persists the current EditingFact, dispatching to CreateFact or UpdateFact based on IsNewFact.
//
// Returns:
//   - A error value.
//
// Side effects:
//   - None.
func (c *IntentContext) SaveEdit() error {
	if c.EditingFact == nil {
		return ErrInvalidState
	}
	if c.IsNewFact {
		return c.CreateFact(c.EditingFact)
	}
	return c.UpdateFact(c.EditingFact)
}

// ToggleRowExpansion flips the expanded/collapsed state of a table row for detail visibility.
//
// Expected:
//   - int must be valid.
//
// Side effects:
//   - None.
func (c *IntentContext) ToggleRowExpansion(rowIndex int) {
	if c.ExpandedRows[rowIndex] {
		delete(c.ExpandedRows, rowIndex)
	} else {
		c.ExpandedRows[rowIndex] = true
	}
}

// IsRowExpanded checks whether a specific table row is in the expanded state.
//
// Expected: rowIndex must be a valid zero-based row offset.
//
// Returns: true if the row is expanded, false otherwise.
//
// Side effects: None.
func (c *IntentContext) IsRowExpanded(rowIndex int) bool {
	return c.ExpandedRows[rowIndex]
}
