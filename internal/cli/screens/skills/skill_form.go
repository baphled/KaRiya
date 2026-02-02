package skills

import (
	"github.com/baphled/kariya/internal/cli/forms"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/domain/career"
)

// SkillFormState identifies the skill add/edit form view in the state
// matrix. On this screen the user sees an interactive form with fields for
// skill name, category, level, and years of experience. In add mode the
// fields are empty; in edit mode they are pre-populated with the existing
// skill data. Tab advances between fields, Enter submits the completed
// form, and Escape cancels without saving.
const SkillFormState = "skill_form"

// SkillFormScreen provides a form for adding or editing skills.
//
// This screen wraps FormScreen[*forms.SkillFormData] with skill-specific context:
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
// - FormScreen provides the form UI
// - forms.SkillFormData defines form structure
// - forms.NewSkillFormWithDataAndDimensions creates the huh form
// - docs/FORMS_GUIDE.md (Form patterns and best practices).
type SkillFormScreen struct {
	*base.FormScreen[*forms.SkillFormData]
}

// NewSkillFormScreen creates a new skill form screen.
//
// Expected:
//   - skill must be valid.
//
// Returns:
//   - A fully initialized SkillFormScreen ready for use.
//
// Side effects:
//   - None.
func NewSkillFormScreen(skill *career.Skill) *SkillFormScreen {
	var breadcrumbs []string
	if skill == nil {
		breadcrumbs = []string{"Main Menu", "Manage Skills", "Add Skill"}
	} else {
		breadcrumbs = []string{"Main Menu", "Manage Skills", "Edit Skill"}
	}

	var formData *forms.SkillFormData
	if skill == nil {
		formData = &forms.SkillFormData{}
	} else {
		formData = forms.GetSkillFormData(skill)
	}

	baseScreen := base.NewBaseFormScreen(breadcrumbs, forms.NewSkillFormWithDataAndDimensions, formData)

	return &SkillFormScreen{
		FormScreen: baseScreen,
	}
}
