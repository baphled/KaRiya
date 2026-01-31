package fixtures

import (
	"fmt"
	"time"

	"github.com/bluele/factory-go/factory"
	"github.com/brianvoe/gofakeit/v7"

	"github.com/baphled/kariya/internal/domain/career"
)

// CVConfigFactory creates CVConfig fixtures with realistic fake data.
var CVConfigFactory = factory.NewFactory(
	&career.CVConfig{},
).Attr("Name", func(args factory.Args) (interface{}, error) {
	prefixes := []string{"Backend", "Frontend", "Platform", "Infrastructure", "Data"}
	suffixes := []string{"Engineer", "Architect", "Lead", "Manager"}
	return fmt.Sprintf("%s %s CV", prefixes[gofakeit.Number(0, len(prefixes)-1)], suffixes[gofakeit.Number(0, len(suffixes)-1)]), nil
}).Attr("TargetRole", func(args factory.Args) (interface{}, error) {
	roles := []string{"staff", "principal", "em", "senior_ic"}
	return roles[gofakeit.Number(0, len(roles)-1)], nil
}).Attr("TargetAudience", func(args factory.Args) (interface{}, error) {
	audiences := []string{"hiring_manager", "recruiter", "peer"}
	return audiences[gofakeit.Number(0, len(audiences)-1)], nil
}).Attr("CreatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
}).Attr("UpdatedAt", func(args factory.Args) (interface{}, error) {
	return time.Now(), nil
})

// CVConfig creates a minimal valid CVConfig with the given name.
func CVConfig(name string) *career.CVConfig {
	now := time.Now()
	return &career.CVConfig{
		Name:           name,
		TargetRole:     "staff",
		TargetAudience: "hiring_manager",
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// CVConfigWith creates a CVConfig with custom key fields.
func CVConfigWith(name, targetRole, targetAudience string) *career.CVConfig {
	now := time.Now()
	return &career.CVConfig{
		Name:           name,
		TargetRole:     targetRole,
		TargetAudience: targetAudience,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// CVConfigWithTech creates a CVConfig with technology focus settings.
func CVConfigWithTech(name, techFocus, focusArea string, technologies []string) *career.CVConfig {
	cfg := CVConfig(name)
	cfg.TechnologyFocus = techFocus
	cfg.FocusArea = focusArea
	cfg.SelectedTechnologies = technologies
	return cfg
}

// CVConfigWithFilters creates a CVConfig with event filters.
func CVConfigWithFilters(name string, filters map[string]interface{}) *career.CVConfig {
	cfg := CVConfig(name)
	cfg.EventFilters = filters
	return cfg
}

// CVConfigWithFormat creates a CVConfig with length and skills format settings.
func CVConfigWithFormat(name, lengthFormat, skillsFormat string, skillsLimit int) *career.CVConfig {
	cfg := CVConfig(name)
	cfg.LengthFormat = lengthFormat
	cfg.SkillsFormat = skillsFormat
	cfg.SkillsLimit = skillsLimit
	return cfg
}
