// Package event provides view components for timeline event management.
package event

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/themes"
	"github.com/baphled/kariya/internal/ui/uikit/primitives"
	overlay "github.com/rmhubbert/bubbletea-overlay"
)

// staticViewModel implements overlay.Viewable for rendering static content.
// This is used with bubbletea-overlay for modal compositing.
type staticViewModel struct {
	content string
}

// View returns the pre-rendered string content stored in this model. It satisfies the
// overlay.Viewable interface so that static markup (such as a modal or background view)
// can be composed via bubbletea-overlay without requiring a full Bubble Tea model.
//
// Returns:
//   - string: the pre-rendered content.
//
// Side effects:
//   - None.
func (m staticViewModel) View() string { return m.content }

// RenderOverlayModal renders a modal view over a background using bubbletea-overlay.
//
// Expected:
//   - Must be a valid string.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderOverlayModal(modalView, backgroundView string) string {
	modalContent := staticViewModel{content: modalView}
	bgModel := staticViewModel{content: backgroundView}

	overlayModel := overlay.New(
		modalContent,
		bgModel,
		overlay.Center,
		overlay.Center,
		0,
		-2,
	)

	return overlayModel.View()
}

// RenderEventDetailContent renders career event details as formatted content.
//
// Expected:
//   - event must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderEventDetailContent(event display.Event, _ themes.Theme) string {
	if isEmptyDisplayEvent(event) {
		return "No event selected."
	}

	var content strings.Builder
	content.WriteString("\nEvent Details\n\n")

	_, _ = fmt.Fprintf(&content, "Date: %s\n", event.Date.Format("2006-01-02"))

	if event.Company != "" {
		_, _ = fmt.Fprintf(&content, "Company: %s\n", event.Company)
	}
	if event.Project != "" {
		_, _ = fmt.Fprintf(&content, "Project: %s\n", event.Project)
	}

	_, _ = fmt.Fprintf(&content, "\nText:\n%s\n", event.Text)

	if len(event.Tags) > 0 {
		_, _ = fmt.Fprintf(&content, "\nTags: %s\n", strings.Join(event.Tags, ", "))
	}
	if len(event.Categories) > 0 {
		_, _ = fmt.Fprintf(&content, "Categories: %s\n", strings.Join(event.Categories, ", "))
	}

	if len(event.Skills) > 0 {
		_, _ = fmt.Fprintf(&content, "Skills: %d associated\n", len(event.Skills))
	}

	return content.String()
}

// RenderSkillsContent renders a list of skills as formatted content.
//
// Expected:
//   - skill must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func RenderSkillsContent(skills []display.Skill, theme themes.Theme) string {
	if len(skills) == 0 {
		return primitives.Muted("No skills associated with this event.", theme).Italic().Render()
	}
	visibleSkills := filterVisibleSkills(skills)
	if len(visibleSkills) == 0 {
		return primitives.Muted("No skills associated with this event.", theme).Italic().Render()
	}
	return renderVisibleSkills(visibleSkills, theme)
}

// filterVisibleSkills returns the skills with meaningful display content.
//
// Expected:
//   - skills must be valid.
//
// Returns:
//   - A slice of display.Skill values.
//
// Side effects:
//   - None.
func filterVisibleSkills(skills []display.Skill) []display.Skill {
	visibleSkills := make([]display.Skill, 0, len(skills))
	for i := range skills {
		skill := skills[i]
		if isEmptyDisplaySkill(skill) {
			continue
		}
		visibleSkills = append(visibleSkills, skill)
	}
	return visibleSkills
}

// renderVisibleSkills renders the visible skills list.
//
// Expected:
//   - skills must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func renderVisibleSkills(skills []display.Skill, theme themes.Theme) string {
	var content strings.Builder
	separator := primitives.Muted(" | ", theme).Render()
	for i := range skills {
		skill := skills[i]
		content.WriteString(primitives.NewText("- "+skill.Name, theme).Bold().Render())
		content.WriteString("\n")
		details := buildSkillDetails(skill, theme)
		if len(details) > 0 {
			content.WriteString("  ")
			content.WriteString(strings.Join(details, separator))
			content.WriteString("\n")
		}
		if i < len(skills)-1 {
			content.WriteString("\n")
		}
	}
	return content.String()
}

// buildSkillDetails builds the detail strings for a skill.
//
// Expected:
//   - skill must be valid.
//   - th must be a valid theme instance (can be nil).
//
// Returns:
//   - A string slice value.
//
// Side effects:
//   - None.
func buildSkillDetails(skill display.Skill, theme themes.Theme) []string {
	var details []string
	if skill.Category != "" {
		details = append(details, primitives.Muted(skill.Category, theme).Render())
	}
	if skill.Level != "" {
		details = append(details, primitives.SuccessText(skill.Level, theme).Render())
	}
	if skill.YearsUsed != nil && *skill.YearsUsed > 0 {
		details = append(details, primitives.InfoText(formatYearsUsed(*skill.YearsUsed), theme).Render())
	}
	return details
}

// formatYearsUsed formats the years used label.
//
// Expected:
//   - yearsUsed must be valid.
//
// Returns:
//   - A string value.
//
// Side effects:
//   - None.
func formatYearsUsed(yearsUsed int) string {
	label := "year"
	if yearsUsed > 1 {
		label = "years"
	}
	return fmt.Sprintf("%d %s", yearsUsed, label)
}

func isEmptyDisplaySkill(skill display.Skill) bool {
	return skill.ID == "" &&
		skill.Name == "" &&
		skill.Category == "" &&
		skill.Level == "" &&
		skill.YearsUsed == nil &&
		skill.LastUsed == nil &&
		skill.CreatedAt.IsZero() &&
		skill.UpdatedAt.IsZero()
}
