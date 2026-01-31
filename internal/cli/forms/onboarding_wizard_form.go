package forms

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
)

// OnboardingFormData holds the data collected from the onboarding wizard.
type OnboardingFormData struct {
	Name      string
	Email     string
	Location  string
	Title     string
	GitHub    string
	Portfolio string
}

// NewOnboardingWizardForm creates a 3-step onboarding wizard form.
// The form collects: Welcome+Name, Contact info, Professional details.
// Data fields are bound via pointers so huh updates them directly.
func NewOnboardingWizardForm(data *OnboardingFormData, width, height int) *huh.Form {
	step1 := huh.NewGroup(
		huh.NewNote().
			Title("Welcome to KaRiya!").
			Description("Let's set up your profile for CV generation.\nThis information will appear on your CVs."),
		huh.NewInput().
			Key("name").
			Title("Your Name").
			Description("Required - appears at the top of your CV").
			Placeholder("e.g., Jane Doe").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("name is required")
				}
				return nil
			}).
			Value(&data.Name),
	).Title("Step 1 of 3: Welcome")

	step2 := huh.NewGroup(
		huh.NewInput().
			Key("email").
			Title("Email Address").
			Description("Required - contact information for your CV").
			Placeholder("e.g., jane@example.com").
			Validate(func(s string) error {
				if strings.TrimSpace(s) == "" {
					return fmt.Errorf("email is required")
				}
				if !strings.Contains(s, "@") {
					return fmt.Errorf("please enter a valid email address")
				}
				return nil
			}).
			Value(&data.Email),
		huh.NewInput().
			Key("location").
			Title("Location").
			Description("Optional - e.g., city, country, or 'Remote'").
			Placeholder("e.g., London, UK").
			Value(&data.Location),
	).Title("Step 2 of 3: Contact")

	step3 := huh.NewGroup(
		huh.NewInput().
			Key("title").
			Title("Professional Title").
			Description("Optional - your current role or target role").
			Placeholder("e.g., Senior Software Engineer").
			Value(&data.Title),
		huh.NewInput().
			Key("github").
			Title("GitHub Username").
			Description("Optional - your GitHub username (not full URL)").
			Placeholder("e.g., baphled").
			Validate(GitHubUsername).
			Value(&data.GitHub),
		huh.NewInput().
			Key("portfolio").
			Title("Portfolio/Website").
			Description("Optional - personal website or portfolio").
			Placeholder("e.g., https://janedoe.dev").
			Value(&data.Portfolio),
	).Title("Step 3 of 3: Professional Details")

	form := huh.NewForm(step1, step2, step3).
		WithTheme(Theme()).
		WithWidth(width).
		WithHeight(DefaultFormHeight(height)).
		WithShowHelp(true).
		WithShowErrors(true)

	return form
}
