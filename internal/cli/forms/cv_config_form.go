package forms

import (
	"github.com/charmbracelet/huh"
)

// CVConfigFormData holds the wizard configuration data.
type CVConfigFormData struct {
	ProfileID       string
	Audience        string
	TechFocus       string
	Technologies    []string
	Technology      string
	FocusArea       string
	SkillsFormat    string
	SkillsLimit     int
	CVLength        string
	SubmitConfirmed bool
}

// SkillsLimitOption represents a preset option for skills limit.
type SkillsLimitOption struct {
	Value int
	Label string
}

// SkillsLimitOptions returns the preset options for skills limit selection.
//
// Returns:
//   - A []SkillsLimitOption value.
//
// Side effects:
//   - None.
func SkillsLimitOptions() []SkillsLimitOption {
	return []SkillsLimitOption{
		{Value: 3, Label: "3"},
		{Value: 5, Label: "5"},
		{Value: 10, Label: "10"},
		{Value: 15, Label: "15"},
		{Value: 0, Label: "All"},
	}
}

// ProfileOption represents a CV profile option.
type ProfileOption struct {
	ID   string
	Name string
}

// ExtractedTechnology represents a technology extracted from events.
type ExtractedTechnology struct {
	Name string
}

// NewCVConfigForm creates the CV configuration wizard form.
// singleTechSelect: when true, uses single-select for technologies (specialist mode),
// when false, uses multi-select (generalist mode).
//
// Expected:
//   - data must be a valid CVConfigFormData pointer.
//   - profileoptions must be a valid slice of ProfileOption.
//   - extractedtechs must be a valid slice of ExtractedTechnology.
//   - width must be a positive integer.
//   - height must be a positive integer.
//   - singletechselect must be a valid boolean.
//
// Returns:
//   - A fully initialized huh.Form ready for use.
//
// Side effects:
//   - None.
func NewCVConfigForm(
	data *CVConfigFormData, profileOptions []ProfileOption, extractedTechs []ExtractedTechnology,
	width, height int, singleTechSelect bool,
) *huh.Form {
	// Step 1: WHO - Profile and Audience
	profileOpts := make([]huh.Option[string], 0, len(profileOptions))
	for _, profile := range profileOptions {
		profileOpts = append(profileOpts, huh.NewOption(profile.Name, profile.ID))
	}

	step1Fields := []huh.Field{
		huh.NewSelect[string]().
			Key("profile").
			Title("CV Profile").
			Description("Choose which profile to use for this CV").
			Options(profileOpts...).
			Value(&data.ProfileID),

		huh.NewSelect[string]().
			Key("audience").
			Title("Target Audience").
			Description("Who will be reading this CV?").
			Options(
				huh.NewOption("Hiring Manager", "hiring_manager"),
				huh.NewOption("Recruiter", "recruiter"),
				huh.NewOption("Technical Peer", "peer"),
			).
			Value(&data.Audience),
	}

	// Step 2: TECH - Technology focus
	step2Fields := []huh.Field{
		huh.NewSelect[string]().
			Key("tech_focus").
			Title("Technology Focus").
			Description("How should we present your technical expertise?").
			Options(
				huh.NewOption("Language Agnostic", "language_agnostic"),
				huh.NewOption("Generalist", "generalist"),
				huh.NewOption("Specialist", "specialist"),
			).
			Value(&data.TechFocus),
	}

	if len(extractedTechs) > 0 {
		techOptions := make([]huh.Option[string], 0, len(extractedTechs))
		for _, tech := range extractedTechs {
			techOptions = append(techOptions, huh.NewOption(tech.Name, tech.Name))
		}

		// Use single-select for specialist mode, multi-select for generalist
		var techField huh.Field
		if singleTechSelect {
			techField = huh.NewSelect[string]().
				Key("technology").
				Title("Select Technology to Highlight").
				Description("Choose ONE technology to specialize in").
				Options(techOptions...).
				Value(&data.Technology)
		} else {
			techField = huh.NewMultiSelect[string]().
				Key("technologies").
				Title("Select Technologies to Highlight").
				Description("Choose technologies to emphasize (leave empty for all)").
				Options(techOptions...).
				Value(&data.Technologies)
		}

		step2Fields = append(step2Fields,
			techField,

			huh.NewSelect[string]().
				Key("focus_area").
				Title("Focus Area").
				Description("Primary area of expertise").
				Options(
					huh.NewOption("Backend", "backend"),
					huh.NewOption("Frontend", "frontend"),
					huh.NewOption("Full Stack", "fullstack"),
					huh.NewOption("DevOps", "devops"),
				).
				Value(&data.FocusArea),
		)
	}

	// Set default skills limit if not already set
	if data.SkillsLimit == 0 {
		data.SkillsLimit = 5
	}

	// Build skills limit options for the select
	skillsLimitOpts := make([]huh.Option[int], 0, 5)
	for _, opt := range SkillsLimitOptions() {
		skillsLimitOpts = append(skillsLimitOpts, huh.NewOption(opt.Label, opt.Value))
	}

	// Step 3: FORMAT - Skills formatting and CV length
	step3Fields := []huh.Field{
		huh.NewSelect[string]().
			Key("skills_format").
			Title("Skills Presentation").
			Description("How should skills be organized?").
			Options(
				huh.NewOption("Grouped", "grouped"),
				huh.NewOption("Flat", "flat"),
			).
			Value(&data.SkillsFormat),

		huh.NewSelect[int]().
			Key("skills_limit").
			Title("Skills Limit").
			Description("Maximum skills to show").
			Options(skillsLimitOpts...).
			Value(&data.SkillsLimit),

		huh.NewSelect[string]().
			Key("cv_length").
			Title("CV Length").
			Description("Target length for the CV").
			Options(
				huh.NewOption("1 Page", "1_page"),
				huh.NewOption("2 Pages", "2_page"),
				huh.NewOption("Standard", "standard"),
				huh.NewOption("Detailed", "detailed"),
			).
			Value(&data.CVLength),
	}

	// Create groups for each step
	groups := []*huh.Group{
		huh.NewGroup(step1Fields...).Title("Step 1: WHO - Profile & Audience"),
		huh.NewGroup(step2Fields...).Title("Step 2: TECH - Technology Focus"),
		huh.NewGroup(step3Fields...).Title("Step 3: FORMAT - Presentation"),
	}

	// Use standard form creation with theme
	form := huh.NewForm(groups...).
		WithTheme(Theme()).
		WithWidth(width)

	if height > 0 {
		form = form.WithHeight(height)
	}

	return form
}
