package intents

import (
	"context"
	"errors"
	"time"

	domain "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

var (
	ErrFactNotFound = errors.New("fact not found")
)

type FactManagementState string

const (
	FactListState          FactManagementState = "list"
	FactViewState          FactManagementState = "view"
	FactEditorState        FactManagementState = "editor"
	FactDeleteConfirmState FactManagementState = "delete_confirm"
	FactResultsState       FactManagementState = "results"
	FactCompletedState     FactManagementState = "completed"
)

type FactManagementContext struct {
	CurrentState      FactManagementState
	Facts             []*domain.Fact
	SelectedFact      *domain.Fact
	SelectedFactIndex int
	SearchText        string
	MinQuality        float64
	MaxQuality        float64
	FilterSource      string
	SortBy            string
	SortOrder         string
	CurrentPage       int
	PageSize          int
	TotalFacts        int
	EditingFact       *domain.Fact
	FormErrors        map[string]string
	IsNewFact         bool
	FactToDelete      *domain.Fact
	ScrollPosition    int
	ExpandedRows      map[int]bool
	FactRepository    careerrepo.FactRepository
	Context           context.Context
	PreviousState     FactManagementState
}

type FactManagementResult struct {
	Action         string
	Fact           *domain.Fact
	Facts          []*domain.Fact
	Error          error
	Message        string
	ScrollPosition int
}

func NewFactManagementContext(factRepo careerrepo.FactRepository, ctx context.Context) *FactManagementContext {
	return &FactManagementContext{
		CurrentState:      FactListState,
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
		PreviousState:     FactListState,
	}
}

func (c *FactManagementContext) LoadFacts() error {
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
	// Initialize to 0 instead of -1 to fix pagination off-by-one error
	if len(facts) > 0 {
		c.SelectedFactIndex = 0
		c.SelectedFact = facts[0]
	} else {
		c.SelectedFactIndex = -1
		c.SelectedFact = nil
	}
	return nil
}

func (c *FactManagementContext) GetPageFacts() []*domain.Fact {
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

func (c *FactManagementContext) SelectFact(index int) {
	pageFacts := c.GetPageFacts()
	if index >= 0 && index < len(pageFacts) {
		c.SelectedFact = pageFacts[index]
		c.SelectedFactIndex = index
	}
}

func (c *FactManagementContext) GetSelectedFact() *domain.Fact {
	return c.SelectedFact
}

func (c *FactManagementContext) ClearFormErrors() {
	c.FormErrors = make(map[string]string)
}

func (c *FactManagementContext) SetFormError(field, message string) {
	c.FormErrors[field] = message
}

func (c *FactManagementContext) HasFormErrors() bool {
	return len(c.FormErrors) > 0
}

func (c *FactManagementContext) CreateFact(fact *domain.Fact) error {
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

func (c *FactManagementContext) UpdateFact(fact *domain.Fact) error {
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

func (c *FactManagementContext) DeleteFact(factID string) error {
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

func (c *FactManagementContext) StartNewFact() {
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

func (c *FactManagementContext) StartEditFact(fact *domain.Fact) {
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

func (c *FactManagementContext) CancelEdit() {
	c.EditingFact = nil
	c.IsNewFact = false
	c.ClearFormErrors()
}

func (c *FactManagementContext) SaveEdit() error {
	if c.EditingFact == nil {
		return ErrInvalidState
	}
	if c.IsNewFact {
		return c.CreateFact(c.EditingFact)
	} else {
		return c.UpdateFact(c.EditingFact)
	}
}

func (c *FactManagementContext) ToggleRowExpansion(rowIndex int) {
	if c.ExpandedRows[rowIndex] {
		delete(c.ExpandedRows, rowIndex)
	} else {
		c.ExpandedRows[rowIndex] = true
	}
}

func (c *FactManagementContext) IsRowExpanded(rowIndex int) bool {
	return c.ExpandedRows[rowIndex]
}
