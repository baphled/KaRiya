package skills

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
)

// State constant for state matrix tracking (REQUIRED)
const SkillFormState = "skill_form"

// SkillFormScreen provides a form for adding or editing skills.
//
// This screen wraps BaseFormScreen[*forms.SkillFormData] with skill-specific context:
// - Pre-populates form with existing skill data when editing
// - Uses SkillFormData with validation from forms package
// - Handles terminal resize by rebuilding form
// - Returns SubmitResult with form data on submission
// - Returns CancelResult on escape
//
// Usage (New Skill):
//
//	screen := skills.NewSkillFormScreen(nil)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultSubmit {
//	    submitResult := result.(*screens.SubmitResult)
//	    data := submitResult.FormData.(*forms.SkillFormData)
//	    if data.SubmitConfirmed {
//	        newSkill := &career.Skill{}
//	        forms.ApplySkillFormData(newSkill, data)
//	        // Save new skill
//	    }
//	}
//
// Usage (Edit Skill):
//
//	screen := skills.NewSkillFormScreen(existingSkill)
//	cmd, result := screen.Update(msg)
//	if result != nil && result.Type() == screens.ResultSubmit {
//	    submitResult := result.(*screens.SubmitResult)
//	    data := submitResult.FormData.(*forms.SkillFormData)
//	    if data.SubmitConfirmed {
//	        forms.ApplySkillFormData(existingSkill, data)
//	        // Save updated skill
//	    }
//	}
//
// Related:
// - BaseFormScreen provides the form UI
// - forms.SkillFormData defines form structure
// - forms.NewSkillFormWithDataAndDimensions creates the huh form
// - docs/FORMS_GUIDE.md (Form patterns and best practices)
type SkillFormScreen struct {
	*base.BaseFormScreen[*forms.SkillFormData]
}

// NewSkillFormScreen creates a new skill form screen.
//
// The screen:
// - Pre-populates form data if skill is not nil (edit mode)
// - Creates empty form if skill is nil (add mode)
// - Uses forms.NewSkillFormWithDataAndDimensions for form building
// - Supports standard navigation: Esc to cancel, Tab for next field
//
// Parameters:
//   - skill: The skill to edit (nil for new skill)
func NewSkillFormScreen(skill *career.Skill) *SkillFormScreen {
	// Determine breadcrumbs based on mode
	var breadcrumbs []string
	if skill == nil {
		breadcrumbs = []string{"Main Menu", "Manage Skills", "Add Skill"}
	} else {
		breadcrumbs = []string{"Main Menu", "Manage Skills", "Edit Skill"}
	}

	// Get form data (empty if skill is nil, populated if not)
	var formData *forms.SkillFormData
	if skill == nil {
		formData = &forms.SkillFormData{}
	} else {
		formData = forms.GetSkillFormData(skill)
	}

	// Use forms.NewSkillFormWithDataAndDimensions directly as the builder.
	baseScreen := base.NewBaseFormScreen(breadcrumbs, forms.NewSkillFormWithDataAndDimensions, formData)

	return &SkillFormScreen{
		BaseFormScreen: baseScreen,
	}
}
