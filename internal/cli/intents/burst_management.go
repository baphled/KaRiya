package intents

import (
	"context"
	"errors"
	"time"

	domain "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
)

var (
	// ErrServiceNotAvailable is returned when the service is not available
	ErrServiceNotAvailable = errors.New("service not available")

	// ErrInvalidState is returned when an operation is attempted in an invalid state
	ErrInvalidState = errors.New("invalid state")

	// ErrNoBurstSelected is returned when no burst is selected
	ErrNoBurstSelected = errors.New("no burst selected")

	// ErrInvalidBurstData is returned when burst data is invalid
	ErrInvalidBurstData = errors.New("invalid burst data")
)

// BurstManagementState represents the current state of the BurstManagement intent
type BurstManagementState string

const (
	BurstListState          BurstManagementState = "list"
	BurstViewState          BurstManagementState = "view"
	BurstEditorState        BurstManagementState = "editor"
	BurstDeleteConfirmState BurstManagementState = "delete_confirm"
	BurstSuggestState       BurstManagementState = "suggest"
	BurstCompletedState     BurstManagementState = "completed"
)

// BurstManagementContext holds the state and data for BurstManagement intent
type BurstManagementContext struct {
	// Current state
	CurrentState BurstManagementState

	// Burst data
	Bursts             []*domain.Burst
	SelectedBurst      *domain.Burst
	SelectedBurstIndex int

	// Filter/Search
	FilterCompetency string
	SearchText       string
	SortBy           string // "name", "createdAt", "updatedAt", "eventCount"
	SortOrder        string // "asc", "desc"

	// Pagination
	CurrentPage int
	PageSize    int // 20
	TotalBursts int

	// Editor state
	EditingBurst *domain.Burst
	FormErrors   map[string]string
	IsNewBurst   bool

	// Suggestions
	Suggestions             []*BurstSuggestion
	SelectedSuggestionIndex int

	// Delete confirmation
	BurstToDelete *domain.Burst

	// UI state
	ScrollPosition int
	ExpandedRows   map[int]bool

	// Services
	Service         *careerservice.Service
	BurstRepository career.BurstRepository

	// Context
	Context context.Context

	// Metadata
	PreviousState       BurstManagementState
	ScrollRestoreNeeded bool
}

// BurstSuggestion represents an AI-suggested burst grouping
type BurstSuggestion struct {
	Title                string
	Description          string
	Events               []*domain.CareerEvent
	ConfidenceScore      float64 // 0.0 to 1.0
	RecommendedStartDate time.Time
	RecommendedEndDate   time.Time
	Skills               []string
	CompetencyFocus      string
}

// BurstManagementResult is the result returned when BurstManagement intent completes
type BurstManagementResult struct {
	// Action performed: "created", "updated", "deleted", "none"
	Action string

	// The burst that was affected (if applicable)
	Burst *domain.Burst

	// All bursts (for list view)
	Bursts []*domain.Burst

	// Error information (if any)
	Error   error
	Message string

	// Metadata for context restoration
	ScrollPosition  int
	SelectedBurstID string
}

// NewBurstManagementContext creates a new context for BurstManagement intent
func NewBurstManagementContext(service *careerservice.Service, burstRepo career.BurstRepository, ctx context.Context) *BurstManagementContext {
	return &BurstManagementContext{
		CurrentState:            BurstListState,
		Bursts:                  make([]*domain.Burst, 0),
		SelectedBurstIndex:      -1,
		FilterCompetency:        "",
		SearchText:              "",
		SortBy:                  "name",
		SortOrder:               "asc",
		CurrentPage:             0,
		PageSize:                20,
		FormErrors:              make(map[string]string),
		Suggestions:             make([]*BurstSuggestion, 0),
		SelectedSuggestionIndex: -1,
		ExpandedRows:            make(map[int]bool),
		Service:                 service,
		BurstRepository:         burstRepo,
		Context:                 ctx,
		PreviousState:           BurstListState,
		ScrollRestoreNeeded:     false,
	}
}

// LoadBursts loads all bursts from the repository
func (c *BurstManagementContext) LoadBursts() error {
	if c.BurstRepository == nil {
		return ErrServiceNotAvailable
	}

	// Use empty filters to get all bursts
	bursts, err := c.BurstRepository.List(c.Context, career.BurstListFilters{})
	if err != nil {
		return err
	}

	c.Bursts = bursts
	c.TotalBursts = len(bursts)
	c.CurrentPage = 0
	c.SelectedBurstIndex = -1
	c.SelectedBurst = nil

	return nil
}

// GetPageBursts returns the bursts for the current page
func (c *BurstManagementContext) GetPageBursts() []*domain.Burst {
	if len(c.Bursts) == 0 {
		return make([]*domain.Burst, 0)
	}

	start := c.CurrentPage * c.PageSize
	end := start + c.PageSize

	if start >= len(c.Bursts) {
		return make([]*domain.Burst, 0)
	}

	if end > len(c.Bursts) {
		end = len(c.Bursts)
	}

	return c.Bursts[start:end]
}

// SelectBurst selects a burst at the given index
func (c *BurstManagementContext) SelectBurst(index int) {
	pageBursts := c.GetPageBursts()
	if index >= 0 && index < len(pageBursts) {
		c.SelectedBurst = pageBursts[index]
		c.SelectedBurstIndex = index
	}
}

// GetSelectedBurst returns the currently selected burst
func (c *BurstManagementContext) GetSelectedBurst() *domain.Burst {
	return c.SelectedBurst
}

// ClearFormErrors clears all form errors
func (c *BurstManagementContext) ClearFormErrors() {
	c.FormErrors = make(map[string]string)
}

// SetFormError sets an error for a specific field
func (c *BurstManagementContext) SetFormError(field, message string) {
	c.FormErrors[field] = message
}

// HasFormErrors returns true if there are any form errors
func (c *BurstManagementContext) HasFormErrors() bool {
	return len(c.FormErrors) > 0
}

// CreateBurst creates a new burst in the repository
func (c *BurstManagementContext) CreateBurst(burst *domain.Burst) error {
	if c.BurstRepository == nil {
		return ErrServiceNotAvailable
	}

	if err := burst.Validate(); err != nil {
		return err
	}

	if err := c.BurstRepository.Create(c.Context, burst); err != nil {
		return err
	}

	c.Bursts = append(c.Bursts, burst)
	c.TotalBursts = len(c.Bursts)

	return nil
}

// UpdateBurst updates an existing burst in the repository
func (c *BurstManagementContext) UpdateBurst(burst *domain.Burst) error {
	if c.BurstRepository == nil {
		return ErrServiceNotAvailable
	}

	if err := burst.Validate(); err != nil {
		return err
	}

	if err := c.BurstRepository.Update(c.Context, burst); err != nil {
		return err
	}

	// Update the burst in the list
	for i, b := range c.Bursts {
		if b.ID == burst.ID {
			c.Bursts[i] = burst
			break
		}
	}

	return nil
}

// DeleteBurst deletes a burst from the repository
func (c *BurstManagementContext) DeleteBurst(burstID string) error {
	if c.BurstRepository == nil {
		return ErrServiceNotAvailable
	}

	if err := c.BurstRepository.Delete(c.Context, burstID); err != nil {
		return err
	}

	// Remove from the list
	for i, b := range c.Bursts {
		if b.ID == burstID {
			c.Bursts = append(c.Bursts[:i], c.Bursts[i+1:]...)
			break
		}
	}

	c.TotalBursts = len(c.Bursts)
	c.SelectedBurst = nil
	c.SelectedBurstIndex = -1

	return nil
}

// StartNewBurst initializes a new burst for editing
func (c *BurstManagementContext) StartNewBurst() {
	c.EditingBurst = &domain.Burst{
		ID:              "",
		Name:            "",
		Description:     "",
		EventIDs:        make([]string, 0),
		CompetencyFocus: "",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	c.IsNewBurst = true
	c.ClearFormErrors()
}

// StartEditBurst initializes editing of an existing burst
func (c *BurstManagementContext) StartEditBurst(burst *domain.Burst) {
	// Create a copy to avoid modifying the original
	c.EditingBurst = &domain.Burst{
		ID:              burst.ID,
		Name:            burst.Name,
		Description:     burst.Description,
		EventIDs:        append([]string{}, burst.EventIDs...),
		CompetencyFocus: burst.CompetencyFocus,
		CreatedAt:       burst.CreatedAt,
		UpdatedAt:       burst.UpdatedAt,
	}
	c.IsNewBurst = false
	c.ClearFormErrors()
}

// CancelEdit cancels the current edit without saving
func (c *BurstManagementContext) CancelEdit() {
	c.EditingBurst = nil
	c.IsNewBurst = false
	c.ClearFormErrors()
}

// SaveEdit saves the current edit
func (c *BurstManagementContext) SaveEdit() error {
	if c.EditingBurst == nil {
		return ErrInvalidState
	}

	if c.IsNewBurst {
		return c.CreateBurst(c.EditingBurst)
	} else {
		return c.UpdateBurst(c.EditingBurst)
	}
}

// ToggleRowExpansion toggles the expansion state of a row
func (c *BurstManagementContext) ToggleRowExpansion(rowIndex int) {
	if c.ExpandedRows[rowIndex] {
		delete(c.ExpandedRows, rowIndex)
	} else {
		c.ExpandedRows[rowIndex] = true
	}
}

// IsRowExpanded returns true if a row is expanded
func (c *BurstManagementContext) IsRowExpanded(rowIndex int) bool {
	return c.ExpandedRows[rowIndex]
}

// Validate ensures the context is complete.
func (c *BurstManagementContext) Validate() error {
	if c.Bursts == nil {
		c.Bursts = make([]*domain.Burst, 0)
	}
	return nil
}
