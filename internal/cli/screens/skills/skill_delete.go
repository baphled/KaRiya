package skills

import (
	"fmt"

	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
)

// SkillDeleteConfirmState identifies the skill deletion confirmation view
// in the state matrix. On this screen the user sees the skill name and a
// yes/no prompt to confirm permanent deletion. The skill name appears
// inside a confirmation dialog with Delete and Cancel buttons. A warning
// states the action cannot be undone. The Cancel button is focused by
// default for safety. Left/right arrows or h/l toggle selection, y/n
// submit directly, and Escape cancels.
const SkillDeleteConfirmState = "skill_delete_confirm"

// SkillDeleteConfirmScreen provides a confirmation dialog for deleting a skill.
//
// This screen wraps ConfirmScreen with skill-specific context:
// - Shows skill name in confirmation message
// - Customizes button text ("Delete" / "Cancel")
// - Preserves skill data for the caller
//
// Usage:
//
//	screen := skills.NewSkillDeleteConfirmScreen(skill)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultNavigate {
//	    confirmed := result.Data().(bool)
//	    if confirmed {
//	        skill := screen.GetSkill()
//	        // Proceed with deletion
//	    } else {
//	        // User cancelled - go back
//	    }
//	}
//
// Related:
// - ConfirmScreen provides the confirmation UI
// - docs/TUI_DEVELOPER_GUIDE.md (Screen patterns)
// - docs/TUI_STANDARDS.md (Keyboard shortcuts).
type SkillDeleteConfirmScreen struct {
	*base.ConfirmScreen
	skill *career.Skill
}

// NewSkillDeleteConfirmScreen creates a new skill deletion confirmation screen.
//
// The screen:
// - Displays skill name in the confirmation message
// - Uses "Delete" / "Cancel" button text
// - Defaults to "No" (Cancel) for safety
// - Supports standard navigation: ←→/hl to toggle, y/n for direct, Enter to confirm
//
// Parameters:
//   - skill: The skill to delete
func NewSkillDeleteConfirmScreen(skill *career.Skill) *SkillDeleteConfirmScreen {
	breadcrumbs := []string{"Main Menu", "Manage Skills", "Delete Confirmation"}
	title := "Delete Skill"
	message := fmt.Sprintf(
		"Are you sure you want to delete '%s'?\nThis action cannot be undone.",
		skill.Name,
	)

	confirmScreen := base.NewBaseConfirmScreen(breadcrumbs, title, message)
	confirmScreen.SetYesText("Delete")
	confirmScreen.SetNoText("Cancel")

	return &SkillDeleteConfirmScreen{
		ConfirmScreen: confirmScreen,
		skill:         skill,
	}
}

// GetSkill returns the skill being considered for deletion.
func (s *SkillDeleteConfirmScreen) GetSkill() *career.Skill {
	return s.skill
}
