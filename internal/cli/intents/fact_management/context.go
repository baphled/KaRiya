// Package fact_management implements the FactManagement intent for managing career facts.
package fact_management

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

// NewIntentContext creates a new context with default values.
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

// Validate ensures the context is complete.
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

// LoadFacts loads facts from the repository.
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

// GetPageFacts returns the facts for the current page.
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

// SelectFact selects a fact by index within the current page.
func (c *IntentContext) SelectFact(index int) {
	pageFacts := c.GetPageFacts()
	if index >= 0 && index < len(pageFacts) {
		c.SelectedFact = pageFacts[index]
		c.SelectedFactIndex = index
	}
}

// GetSelectedFact returns the currently selected fact.
func (c *IntentContext) GetSelectedFact() *domain.Fact {
	return c.SelectedFact
}

// ClearFormErrors clears all form validation errors.
func (c *IntentContext) ClearFormErrors() {
	c.FormErrors = make(map[string]string)
}

// SetFormError sets a validation error for a field.
func (c *IntentContext) SetFormError(field, message string) {
	c.FormErrors[field] = message
}

// HasFormErrors returns true if there are validation errors.
func (c *IntentContext) HasFormErrors() bool {
	return len(c.FormErrors) > 0
}

// CreateFact creates a new fact.
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

// UpdateFact updates an existing fact.
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

// DeleteFact deletes a fact by ID.
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

// StartNewFact prepares for creating a new fact.
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

// StartEditFact prepares for editing an existing fact.
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

// CancelEdit cancels the current edit operation.
func (c *IntentContext) CancelEdit() {
	c.EditingFact = nil
	c.IsNewFact = false
	c.ClearFormErrors()
}

// SaveEdit saves the edited fact.
func (c *IntentContext) SaveEdit() error {
	if c.EditingFact == nil {
		return ErrInvalidState
	}
	if c.IsNewFact {
		return c.CreateFact(c.EditingFact)
	}
	return c.UpdateFact(c.EditingFact)
}

// ToggleRowExpansion toggles the expansion state of a row.
func (c *IntentContext) ToggleRowExpansion(rowIndex int) {
	if c.ExpandedRows[rowIndex] {
		delete(c.ExpandedRows, rowIndex)
	} else {
		c.ExpandedRows[rowIndex] = true
	}
}

// IsRowExpanded returns whether a row is expanded.
func (c *IntentContext) IsRowExpanded(rowIndex int) bool {
	return c.ExpandedRows[rowIndex]
}
