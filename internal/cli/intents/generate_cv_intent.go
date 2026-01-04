package intents

import (
	"fmt"
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// GenerateCVIntent implements the Intent interface for generating CVs.
// It owns the complete lifecycle of CV generation, including:
// - Selecting a profile (target role and audience)
// - Generating the CV from career events and facts
// - Previewing and reviewing the CV
// - Returning the generated CV or cancelling
type GenerateCVIntent struct {
	// context is the input context passed to the intent.
	context *GenerateCVContext

	// state represents the current state of the intent.
	state *GenerateCVModel

	// active indicates whether this intent is currently active.
	active bool

	// result is the final result of the intent (set when complete).
	result *IntentResult[*GenerateCVResult]
}

// NewGenerateCVIntent creates a new GenerateCV intent.
func NewGenerateCVIntent(context *GenerateCVContext) (*GenerateCVIntent, error) {
	// Validate the context.
	if err := context.Validate(); err != nil {
		return nil, err
	}

	selectedProfile := context.DefaultProfile
	if selectedProfile == nil && len(context.AvailableProfiles) > 0 {
		selectedProfile = context.AvailableProfiles[0]
	}

	return &GenerateCVIntent{
		context: context,
		state: &GenerateCVModel{
			context:         context,
			currentState:    GenerateCVStateSelectProfile,
			selectedProfile: selectedProfile,
			selectedIndex:   0,
		},
		active: true,
	}, nil
}

// Init is called when the intent is activated.
func (i *GenerateCVIntent) Init() tea.Cmd {
	// Initialize with default profile if available.
	if i.context.DefaultProfile != nil {
		i.state.selectedProfile = i.context.DefaultProfile
	}
	return nil
}

// Update processes a message in the intent.
func (i *GenerateCVIntent) Update(msg tea.Msg) tea.Cmd {
	if !i.active {
		return nil
	}

	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return i.updateSelectProfile(msg)

	case GenerateCVStateSelectAudience:
		return i.updateSelectAudience(msg)

	case GenerateCVStatePreview:
		return i.updatePreview(msg)

	case GenerateCVStateReview:
		return i.updateReview(msg)

	case GenerateCVStateConfirm:
		return i.updateConfirm(msg)
	}

	return nil
}

// updateSelectProfile handles messages while selecting a profile.
func (i *GenerateCVIntent) updateSelectProfile(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			// Move selection up.
			if i.state.selectedIndex > 0 {
				i.state.selectedIndex--
				i.state.selectedProfile = i.context.AvailableProfiles[i.state.selectedIndex]
			}
			return nil

		case "down", "j":
			// Move selection down.
			if i.state.selectedIndex < len(i.context.AvailableProfiles)-1 {
				i.state.selectedIndex++
				i.state.selectedProfile = i.context.AvailableProfiles[i.state.selectedIndex]
			}
			return nil

		case "enter":
			// Confirm profile selection and move to audience selection.
			if i.state.selectedProfile != nil {
				i.state.currentState = GenerateCVStateSelectAudience
				i.state.selectedAudiences = i.state.selectedProfile.TargetAudience
			}
			return nil

		case "q", "ctrl+c":
			// Cancel.
			i.setCancelled()
			return nil

		case "esc":
			// Go back (no-op at profile selection).
			i.setCancelled()
			return nil
		}

	case ProfileSelectedMsg:
		// Profile was selected (possibly by router or other component).
		i.state.selectedProfile = msg.Profile
		i.state.selectedIndex = msg.Index
		i.state.currentState = GenerateCVStateSelectAudience
		return nil
	}

	return nil
}

// updateSelectAudience handles messages while selecting audience(s).
func (i *GenerateCVIntent) updateSelectAudience(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Confirm audience selection and generate CV.
			i.state.currentState = GenerateCVStatePreview
			// In a real implementation, this would call a service to generate the CV.
			// For now, we'll create a minimal CV.
			i.state.generatedCV = &career.CVView{
				ID:               fmt.Sprintf("cv_%d", time.Now().Unix()),
				Name:             i.state.selectedProfile.Name,
				TargetRole:       i.state.selectedProfile.TargetRole,
				TargetAudience:   i.state.selectedAudiences,
				GeneratedAt:      time.Now(),
				SourceEventCount: len(i.context.Events),
				SourceFactCount:  len(i.context.Facts),
			}
			return nil

		case "q", "ctrl+c":
			// Cancel.
			i.setCancelled()
			return nil

		case "esc":
			// Go back to profile selection.
			i.state.currentState = GenerateCVStateSelectProfile
			return nil
		}

	case AudienceSelectedMsg:
		// Audience was selected (possibly by router or other component).
		i.state.selectedAudiences = msg.Audiences
		i.state.currentState = GenerateCVStatePreview
		return nil
	}

	return nil
}

// updatePreview handles messages while previewing the CV.
func (i *GenerateCVIntent) updatePreview(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", "e":
			// Edit the CV (move to review state).
			i.state.currentState = GenerateCVStateReview
			return nil

		case "c":
			// Confirm the CV (move to confirm state).
			i.state.currentState = GenerateCVStateConfirm
			return nil

		case "q", "ctrl+c":
			// Cancel.
			i.setCancelled()
			return nil

		case "esc":
			// Go back to audience selection.
			i.state.currentState = GenerateCVStateSelectAudience
			return nil
		}
	}

	return nil
}

// updateReview handles messages while reviewing/editing the CV.
func (i *GenerateCVIntent) updateReview(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// Save changes and move to confirm state.
			i.state.currentState = GenerateCVStateConfirm
			return nil

		case "q", "ctrl+c":
			// Cancel.
			i.setCancelled()
			return nil

		case "esc":
			// Go back to preview.
			i.state.currentState = GenerateCVStatePreview
			return nil
		}
	}

	return nil
}

// updateConfirm handles messages while confirming the CV.
func (i *GenerateCVIntent) updateConfirm(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "enter":
			// Confirm and complete.
			i.setCompleted()
			return nil

		case "n", "esc":
			// Go back to review.
			i.state.currentState = GenerateCVStateReview
			return nil

		case "q", "ctrl+c":
			// Cancel.
			i.setCancelled()
			return nil
		}
	}

	return nil
}

// View renders the intent's current state.
func (i *GenerateCVIntent) View() string {
	switch i.state.currentState {
	case GenerateCVStateSelectProfile:
		return i.viewSelectProfile()

	case GenerateCVStateSelectAudience:
		return i.viewSelectAudience()

	case GenerateCVStatePreview:
		return i.viewPreview()

	case GenerateCVStateReview:
		return i.viewReview()

	case GenerateCVStateConfirm:
		return i.viewConfirm()
	}

	return ""
}

// viewSelectProfile renders the profile selection view.
func (i *GenerateCVIntent) viewSelectProfile() string {
	var content strings.Builder
	content.WriteString("\nSelect CV Profile\n\n")

	if len(i.context.AvailableProfiles) == 0 {
		content.WriteString("No profiles available.\n")
	} else {
		for idx, profile := range i.context.AvailableProfiles {
			prefix := "  "
			if idx == i.state.selectedIndex {
				prefix = "> "
			}

			content.WriteString(fmt.Sprintf("%s%s (%s)\n", prefix, profile.Name, profile.TargetRole))
			if profile.Description != "" {
				content.WriteString(fmt.Sprintf("    %s\n", profile.Description))
			}
		}
	}

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("↑/↓ or k/j to navigate, Enter to select, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewSelectAudience renders the audience selection view.
func (i *GenerateCVIntent) viewSelectAudience() string {
	var content strings.Builder
	content.WriteString("\nSelect Target Audience\n\n")

	if i.state.selectedProfile != nil {
		content.WriteString(fmt.Sprintf("Profile: %s\n\n", i.state.selectedProfile.Name))
	}

	audiences := []string{"hiring_manager", "recruiter", "peer"}
	for _, audience := range audiences {
		selected := false
		for _, selected_aud := range i.state.selectedAudiences {
			if selected_aud == audience {
				selected = true
				break
			}
		}

		prefix := "  "
		if selected {
			prefix = "[x] "
		} else {
			prefix = "[ ] "
		}

		content.WriteString(fmt.Sprintf("%s%s\n", prefix, audience))
	}

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter to generate CV, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewPreview renders the CV preview view.
func (i *GenerateCVIntent) viewPreview() string {
	var content strings.Builder
	content.WriteString("\nCV Preview\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Name: %s\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Target Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Target Audiences: %s\n", strings.Join(i.state.generatedCV.TargetAudience, ", ")))
		content.WriteString(fmt.Sprintf("Source Events: %d\n", i.state.generatedCV.SourceEventCount))
		content.WriteString(fmt.Sprintf("Source Facts: %d\n", i.state.generatedCV.SourceFactCount))
	}

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter/e to edit, c to confirm, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewReview renders the review/edit view.
func (i *GenerateCVIntent) viewReview() string {
	var content strings.Builder
	content.WriteString("\nReview & Edit CV\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Name: %s\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Target Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Target Audiences: %s\n", strings.Join(i.state.generatedCV.TargetAudience, ", ")))
		content.WriteString("\n[Edit functionality would go here]\n")
	}

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("Enter to confirm, Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// viewConfirm renders the confirmation view.
func (i *GenerateCVIntent) viewConfirm() string {
	var content strings.Builder
	content.WriteString("\nConfirm CV Generation\n\n")

	if i.state.generatedCV != nil {
		content.WriteString(fmt.Sprintf("Generate CV: %s?\n", i.state.generatedCV.Name))
		content.WriteString(fmt.Sprintf("Role: %s\n", i.state.generatedCV.TargetRole))
		content.WriteString(fmt.Sprintf("Audiences: %s\n", strings.Join(i.state.generatedCV.TargetAudience, ", ")))
	}

	// Apply card styling.
	cardStyle := lipgloss.NewStyle().
		Padding(1, 2).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(styles.ColorBorder).
		Background(styles.ColorBackgroundCard).
		Foreground(styles.ColorTextPrimary)

	card := cardStyle.Render(content.String())

	// Add footer with instructions.
	footerStyle := lipgloss.NewStyle().
		Foreground(styles.ColorTextSecondary).
		MarginTop(1)

	footer := footerStyle.Render("y/Enter to confirm, n/Esc to go back, q to cancel")

	return lipgloss.JoinVertical(lipgloss.Left, card, footer)
}

// Result returns the final result of the intent.
func (i *GenerateCVIntent) Result() *IntentResult[interface{}] {
	if i.result == nil {
		return nil
	}

	return &IntentResult[interface{}]{
		Status:   i.result.Status,
		Data:     i.result.Data,
		Error:    i.result.Error,
		Metadata: i.result.Metadata,
	}
}

// setCompleted marks the intent as completed with success.
func (i *GenerateCVIntent) setCompleted() {
	i.result = &IntentResult[*GenerateCVResult]{
		Status: Completed,
		Data: &GenerateCVResult{
			GeneratedCV:     i.state.generatedCV,
			SelectedProfile: i.state.selectedProfile,
			AcceptedFields:  make(map[string]bool),
		},
		Metadata: map[string]interface{}{
			"profile":     i.state.selectedProfile.ID,
			"audiences":   strings.Join(i.state.selectedAudiences, ","),
			"timestamp":   time.Now(),
			"event_count": len(i.context.Events),
			"fact_count":  len(i.context.Facts),
		},
	}
	i.active = false
}

// setCancelled marks the intent as cancelled by the user.
func (i *GenerateCVIntent) setCancelled() {
	i.result = &IntentResult[*GenerateCVResult]{
		Status: Cancelled,
	}
	i.active = false
}
