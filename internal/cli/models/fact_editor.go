package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Focus indices for fact editor navigation
const (
	FactTextFieldIdx = iota
	FactCompetenciesFieldIdx
	FactRoleFitFieldIdx
	FactAudienceFieldIdx
	FactSaveButtonIdx
	FactCancelButtonIdx
	FactEditorFieldCount // Total number of focus positions
)

// FactEditorModel represents the fact editor form state
type FactEditorModel struct {
	*BaseStandardModel
	fact                 *career.Fact
	originalFact         *career.Fact // For reverting changes
	service              *careerservice.Service
	ctx                  context.Context
	textInput            textinput.Model
	focusIndex           int
	competencyIndex      int // Index for competency navigation within competencies
	roleFitIndex         int // Index for role fit selection
	audienceIndex        int // Index for audience navigation within audience
	err                  error
	submitted            bool
	cancelled            bool
	competencySelector   *components.CategorySelector // Reuse category selector for competencies
	roleFitOptions       []career.RoleFit
	audienceSelector     *components.AudienceRelevanceSelector
	fieldErrors          map[int]string
	width                int
	height               int
	characterCount       int
	originalCompetencies []string // Track original for change detection
	originalRoleFit      career.RoleFit
	originalAudience     []string
	helpFooter           components.HelpFooterModel // Help footer for keyboard shortcuts
}

// NewFactEditorModel creates a new fact editor model
func NewFactEditorModel(fact *career.Fact, service *careerservice.Service, ctx context.Context) *FactEditorModel {
	// Create a copy of the fact for reverting
	factCopy := *fact

	// Create text input
	textInput := textinput.New()
	textInput.Placeholder = "Fact text (1-2000 characters, no aspirational language)"
	textInput.SetValue(fact.Text)
	textInput.Width = 70
	textInput.Focus()

	// Create competency selector from category selector
	competencySelector := components.NewCategorySelector()
	_ = competencySelector.SetSelected(fact.CompetencyCategories)

	// Create audience selector
	audienceSelector := components.NewAudienceRelevanceSelector()
	audienceSelector.SetSelected(fact.AudienceRelevance)

	// Create role fit options
	roleFitOptions := []career.RoleFit{
		career.RoleFitPrincipal,
		career.RoleFitEM,
		career.RoleFitStaff,
		career.RoleFitSeniorIC,
	}

	// Find current role fit index
	roleFitIndex := 0
	for i, rf := range roleFitOptions {
		if rf == fact.RoleFit {
			roleFitIndex = i
			break
		}
	}

	return &FactEditorModel{
		BaseStandardModel:    NewBaseStandardModel(),
		fact:                 fact,
		originalFact:         &factCopy,
		service:              service,
		ctx:                  ctx,
		textInput:            textInput,
		focusIndex:           0,
		competencyIndex:      0,
		roleFitIndex:         roleFitIndex,
		audienceIndex:        0,
		err:                  nil,
		submitted:            false,
		cancelled:            false,
		competencySelector:   competencySelector,
		roleFitOptions:       roleFitOptions,
		audienceSelector:     audienceSelector,
		fieldErrors:          make(map[int]string),
		characterCount:       len(fact.Text),
		originalCompetencies: append([]string{}, fact.CompetencyCategories...),
		originalRoleFit:      fact.RoleFit,
		originalAudience:     append([]string{}, fact.AudienceRelevance...),
		helpFooter:           components.NewHelpFooter("fact_editor", 80),
	}
}

// Init initializes the model
func (m *FactEditorModel) Init() tea.Cmd {
	return nil
}

// Update handles messages
func (m *FactEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKeyMsg(msg)
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}

	// Update text input
	if m.focusIndex == FactTextFieldIdx {
		var cmd tea.Cmd
		m.textInput, cmd = m.textInput.Update(msg)
		m.fact.Text = m.textInput.Value()
		m.characterCount = len(m.fact.Text)
		return m, cmd
	}

	return m, nil
}

func (m *FactEditorModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		// Allow quit from fact editor
		return m, func() tea.Msg { return QuitMsg{} }

	case "esc":
		// Cancel editor and return to parent
		m.cancelled = true
		return m, nil

	case "tab":
		m.clearFieldErrors()
		m.focusIndex = (m.focusIndex + 1) % FactEditorFieldCount
		m.updateInputFocus()
		return m, nil

	case "shift+tab":
		m.clearFieldErrors()
		m.focusIndex = (m.focusIndex - 1 + FactEditorFieldCount) % FactEditorFieldCount
		m.updateInputFocus()
		return m, nil

	case "up":
		if m.focusIndex == FactCompetenciesFieldIdx {
			if m.competencyIndex > 0 {
				m.competencyIndex--
			}
		} else if m.focusIndex == FactRoleFitFieldIdx {
			if m.roleFitIndex > 0 {
				m.roleFitIndex--
			}
		} else if m.focusIndex == FactAudienceFieldIdx {
			if m.audienceIndex > 0 {
				m.audienceIndex--
			}
		}
		return m, nil

	case "down":
		if m.focusIndex == FactCompetenciesFieldIdx {
			competencies := m.competencySelector.AvailableCategories()
			if m.competencyIndex < len(competencies)-1 {
				m.competencyIndex++
			}
		} else if m.focusIndex == FactRoleFitFieldIdx {
			if m.roleFitIndex < len(m.roleFitOptions)-1 {
				m.roleFitIndex++
			}
		} else if m.focusIndex == FactAudienceFieldIdx {
			audiences := []string{"hiring_manager", "recruiter", "peer"}
			if m.audienceIndex < len(audiences)-1 {
				m.audienceIndex++
			}
		}
		return m, nil

	case " ":
		if m.focusIndex == FactCompetenciesFieldIdx {
			// Toggle selected competency
			competencies := m.competencySelector.AvailableCategories()
			if m.competencyIndex < len(competencies) {
				selectedComp := competencies[m.competencyIndex]
				_ = m.competencySelector.ToggleCategory(selectedComp)
				m.fact.CompetencyCategories = m.competencySelector.SelectedCategories()
			}
			return m, nil
		} else if m.focusIndex == FactAudienceFieldIdx {
			// Toggle selected audience
			audiences := []string{"hiring_manager", "recruiter", "peer"}
			if m.audienceIndex < len(audiences) {
				selected := m.audienceSelector.GetSelected()
				if contains(selected, audiences[m.audienceIndex]) {
					// Remove
					newSelected := []string{}
					for _, s := range selected {
						if s != audiences[m.audienceIndex] {
							newSelected = append(newSelected, s)
						}
					}
					m.audienceSelector.SetSelected(newSelected)
				} else {
					// Add
					newSelected := append(selected, audiences[m.audienceIndex])
					m.audienceSelector.SetSelected(newSelected)
				}
				m.fact.AudienceRelevance = m.audienceSelector.GetSelected()
			}
			return m, nil
		}
		return m, nil

	case "enter":
		if m.focusIndex == FactSaveButtonIdx {
			return m.submitFact()
		} else if m.focusIndex == FactCancelButtonIdx {
			m.cancelled = true
			return m, nil
		}
		return m, nil
	}

	return m, nil
}

func (m *FactEditorModel) updateInputFocus() {
	isFocused := m.focusIndex == FactTextFieldIdx
	if isFocused && !m.textInput.Focused() {
		m.textInput.Focus()
	} else if !isFocused && m.textInput.Focused() {
		m.textInput.Blur()
	}
}

func (m *FactEditorModel) clearFieldErrors() {
	m.fieldErrors = make(map[int]string)
}

func (m *FactEditorModel) submitFact() (tea.Model, tea.Cmd) {
	m.clearFieldErrors()

	// Validate text
	if strings.TrimSpace(m.fact.Text) == "" {
		m.fieldErrors[FactTextFieldIdx] = "Fact text cannot be empty"
		m.err = fmt.Errorf("validation error: fact text is required")
		return m, nil
	}

	if len(m.fact.Text) > 2000 {
		m.fieldErrors[FactTextFieldIdx] = "Fact text cannot exceed 2000 characters"
		m.err = fmt.Errorf("validation error: fact text exceeds 2000 characters")
		return m, nil
	}

	// Validate competencies
	if len(m.fact.CompetencyCategories) == 0 {
		m.fieldErrors[FactCompetenciesFieldIdx] = "At least one competency category is required"
		m.err = fmt.Errorf("validation error: no competencies selected")
		return m, nil
	}

	// Validate role fit
	if !m.isValidRoleFit(m.fact.RoleFit) {
		m.fieldErrors[FactRoleFitFieldIdx] = "Invalid role fit selection"
		m.err = fmt.Errorf("validation error: invalid role fit")
		return m, nil
	}

	// Validate audience relevance
	if len(m.fact.AudienceRelevance) == 0 {
		m.fieldErrors[FactAudienceFieldIdx] = "At least one audience type is required"
		m.err = fmt.Errorf("validation error: no audience selected")
		return m, nil
	}

	// Check for aspirational language
	if m.hasAspirationLanguage(m.fact.Text) {
		m.fieldErrors[FactTextFieldIdx] = "Text contains aspirational language (will, should, could, etc.)"
		m.err = fmt.Errorf("validation error: aspirational language detected")
		return m, nil
	}

	// Update timestamp
	m.fact.UpdatedAt = time.Now()

	// Validate the fact
	if err := m.fact.Validate(); err != nil {
		m.err = err
		return m, nil
	}

	m.submitted = true
	return m, nil
}

func (m *FactEditorModel) isValidRoleFit(rf career.RoleFit) bool {
	return career.AllowedRoleFits[string(rf)]
}

func (m *FactEditorModel) hasAspirationLanguage(text string) bool {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)

	for _, word := range words {
		// Remove punctuation
		cleanWord := strings.TrimFunc(word, func(r rune) bool {
			return !isLetterChar(r)
		})

		if career.AspirationKeywords[cleanWord] {
			return true
		}
	}
	return false
}

func isLetterChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

// GetFact returns the edited fact
func (m *FactEditorModel) GetFact() *career.Fact {
	return m.fact
}

// IsSubmitted returns true if changes were saved
func (m *FactEditorModel) IsSubmitted() bool {
	return m.submitted
}

// IsCancelled returns true if operation was cancelled
func (m *FactEditorModel) IsCancelled() bool {
	return m.cancelled
}

// Revert reverts changes to the original fact
func (m *FactEditorModel) Revert() {
	m.fact.Text = m.originalFact.Text
	m.fact.CompetencyCategories = append([]string{}, m.originalCompetencies...)
	m.fact.RoleFit = m.originalRoleFit
	m.fact.AudienceRelevance = append([]string{}, m.originalAudience...)
	m.textInput.SetValue(m.originalFact.Text)
	m.characterCount = len(m.originalFact.Text)
}

// GetError returns the current error
func (m *FactEditorModel) GetError() error {
	return m.err
}

// View renders the editor UI using FormFieldContainers
func (m *FactEditorModel) View() string {
	// Render form content using FormFieldContainers
	formContent := m.renderFormContentWithContainers()

	// Use header and footer components
	headerView := components.NewHeader("Fact Editor", m.width).View()
	footerView := components.NewFooter(m.width).View()

	// Render help footer
	m.helpFooter.SetWidth(m.width)
	helpFooterContent := m.helpFooter.View()

	// Combine all sections
	fullContent := strings.Join([]string{
		headerView,
		"",
		formContent,
		"",
		footerView,
		"",
		helpFooterContent,
	}, "\n")

	return fullContent
}

// renderFormContentWithContainers renders all form fields using FormFieldContainers
// renderFormContentWithContainers renders all form fields using the smart FormContainer
func (m *FactEditorModel) renderFormContentWithContainers() string {
	// Create form container with responsive layout
	formContainer := components.NewFormContainer().
		SetWidth(m.width).
		SetHeight(m.height).
		SetLayout(components.Responsive).
		SetPadding(2).
		SetVerticalSpacing(2).
		SetColumnGap(4).
		SetMinFieldWidth(30)

	// Add text field
	textFieldErr := ""
	if err, ok := m.fieldErrors[FactTextFieldIdx]; ok {
		textFieldErr = err
	}
	charInfo := fmt.Sprintf("Characters: %d/2000 %s",
		m.characterCount, m.getCharCountIndicator())

	formContainer.AddField(components.FormField{
		Label:      "Fact Text (required):",
		Input:      m.textInput.View(),
		Error:      textFieldErr,
		Hint:       charInfo,
		IsFocused:  m.focusIndex == FactTextFieldIdx,
		IsRequired: true,
		FullWidth:  true,
	})

	// Add competencies field
	competencies := m.competencySelector.AvailableCategories()
	var badges []string
	for i, comp := range competencies {
		isSelected := m.competencySelector.IsSelected(comp)
		isFocused := m.focusIndex == FactCompetenciesFieldIdx && m.competencyIndex == i
		badge := m.renderBadge(comp, isSelected, isFocused)
		badges = append(badges, badge)
	}
	badgesStr := strings.Join(badges, "  ")

	competenciesErr := ""
	if err, ok := m.fieldErrors[FactCompetenciesFieldIdx]; ok {
		competenciesErr = err
	}
	competenciesHint := fmt.Sprintf("Select with Space ↑↓ | %d selected", len(m.fact.CompetencyCategories))

	formContainer.AddField(components.FormField{
		Label:      "Competency Categories (required):",
		Input:      badgesStr,
		Error:      competenciesErr,
		Hint:       competenciesHint,
		IsFocused:  m.focusIndex == FactCompetenciesFieldIdx,
		IsRequired: true,
		FullWidth:  false,
	})

	// Add role fit field
	var options []string
	for i, rf := range m.roleFitOptions {
		isSelected := m.fact.RoleFit == rf
		isFocused := m.focusIndex == FactRoleFitFieldIdx && m.roleFitIndex == i

		indicator := "○"
		if isSelected {
			indicator = "●"
		}

		text := fmt.Sprintf("%s %s", indicator, capitalize(string(rf)))
		if isFocused {
			text = styles.ButtonPrimaryFocused.Render(text)
		} else if isSelected {
			text = styles.ButtonPrimary.Render(text)
		}

		options = append(options, text)
	}
	optionsStr := strings.Join(options, "  ")

	roleFitErr := ""
	if err, ok := m.fieldErrors[FactRoleFitFieldIdx]; ok {
		roleFitErr = err
	}

	formContainer.AddField(components.FormField{
		Label:      "Role Fit (required):",
		Input:      optionsStr,
		Error:      roleFitErr,
		Hint:       "Navigate with ↑↓ | Select with Space",
		IsFocused:  m.focusIndex == FactRoleFitFieldIdx,
		IsRequired: true,
		FullWidth:  false,
	})

	// Add audience field
	audiences := []string{"hiring_manager", "recruiter", "peer"}
	var audienceBadges []string
	for i, audience := range audiences {
		isSelected := contains(m.fact.AudienceRelevance, audience)
		isFocused := m.focusIndex == FactAudienceFieldIdx && m.audienceIndex == i
		badge := m.renderBadge(capitalize(audience), isSelected, isFocused)
		audienceBadges = append(audienceBadges, badge)
	}
	audienceBadgesStr := strings.Join(audienceBadges, "  ")

	audienceErr := ""
	if err, ok := m.fieldErrors[FactAudienceFieldIdx]; ok {
		audienceErr = err
	}
	audienceHint := fmt.Sprintf("Select with Space ↑↓ | %d selected", len(m.fact.AudienceRelevance))

	formContainer.AddField(components.FormField{
		Label:      "Audience Relevance (required):",
		Input:      audienceBadgesStr,
		Error:      audienceErr,
		Hint:       audienceHint,
		IsFocused:  m.focusIndex == FactAudienceFieldIdx,
		IsRequired: true,
		FullWidth:  false,
	})

	// Add buttons field
	saveBtn := "[ Save ]"
	cancelBtn := "[ Cancel ]"

	if m.focusIndex == FactSaveButtonIdx {
		saveBtn = styles.ButtonPrimaryFocused.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	} else if m.focusIndex == FactCancelButtonIdx {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondaryFocused.Render(cancelBtn)
	} else {
		saveBtn = styles.ButtonPrimary.Render(saveBtn)
		cancelBtn = styles.ButtonSecondary.Render(cancelBtn)
	}

	buttonsStr := strings.Join([]string{saveBtn, cancelBtn}, "  ")

	formContainer.AddField(components.FormField{
		Input:     buttonsStr,
		IsFocused: m.focusIndex >= FactSaveButtonIdx,
		FullWidth: true,
	})

	// Get rendered form
	formContent := formContainer.Render()

	// Add model-level error if present
	if m.err != nil {
		formContent += "\n\n" + styles.ErrorBox.Render(m.err.Error())
	}

	return formContent
}

// getCharCountIndicator returns a visual indicator for character count
func (m *FactEditorModel) getCharCountIndicator() string {
	if m.characterCount > 1800 {
		return styles.Warning.Render("⚠")
	}
	return ""
}

func (m *FactEditorModel) renderBadge(text string, selected, focused bool) string {
	badge := fmt.Sprintf("[%s]", text)

	if focused {
		return styles.BadgeFocused.Render(badge)
	} else if selected {
		return styles.BadgeSelected.Render(badge)
	}
	return styles.Badge.Render(badge)
}

// Helper functions
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(string(word[0])) + word[1:]
		}
	}
	return strings.Join(words, " ")
}
