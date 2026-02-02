// Package modals provides timeline-specific modal components.
// These modals handle filtering, searching, sorting, and event CRUD operations
// for the BrowseTimeline intent.
package modals

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/domain/career"
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
func RenderEventDetailContent(event *career.Event, theme themes.Theme) string {
	if event == nil {
		return "No event selected."
	}

	var content strings.Builder
	content.WriteString("\nEvent Details\n\n")

	content.WriteString(fmt.Sprintf("Date: %s\n", event.Date.Format("2006-01-02")))

	if event.Company != "" {
		content.WriteString(fmt.Sprintf("Company: %s\n", event.Company))
	}
	if event.Project != "" {
		content.WriteString(fmt.Sprintf("Project: %s\n", event.Project))
	}

	content.WriteString(fmt.Sprintf("\nText:\n%s\n", event.Text))

	if len(event.Tags) > 0 {
		content.WriteString(fmt.Sprintf("\nTags: %s\n", strings.Join(event.Tags, ", ")))
	}
	if len(event.Categories) > 0 {
		content.WriteString(fmt.Sprintf("Categories: %s\n", strings.Join(event.Categories, ", ")))
	}

	if len(event.Skills) > 0 {
		content.WriteString(fmt.Sprintf("Skills: %d associated\n", len(event.Skills)))
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
func RenderSkillsContent(skills []*career.Skill, theme themes.Theme) string {
	if len(skills) == 0 {
		return primitives.Muted("No skills associated with this event.", theme).Italic().Render()
	}

	var content strings.Builder

	separator := primitives.Muted(" | ", theme).Render()

	for i, skill := range skills {
		if skill == nil {
			continue
		}

		content.WriteString(primitives.NewText(fmt.Sprintf("- %s", skill.Name), theme).Bold().Render())
		content.WriteString("\n")

		var details []string

		if skill.Category != "" {
			details = append(details, primitives.Muted(skill.Category, theme).Render())
		}

		if skill.Level != "" {
			details = append(details, primitives.SuccessText(skill.Level, theme).Render())
		}

		if skill.YearsUsed != nil && *skill.YearsUsed > 0 {
			yearText := "year"
			if *skill.YearsUsed > 1 {
				yearText = "years"
			}
			details = append(details, primitives.InfoText(fmt.Sprintf("%d %s", *skill.YearsUsed, yearText), theme).Render())
		}

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
