package app

import (
	"github.com/baphled/kariya/internal/cli/intents"
)

// createDefaultCVProfiles creates a set of default CV profiles for the GenerateCV intent.
func createDefaultCVProfiles() []*intents.CVProfile {
	return []*intents.CVProfile{
		{
			ID:             "profile-staff-engineer",
			Name:           "Staff Engineer",
			TargetRole:     "staff",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for staff engineering roles",
		},
		{
			ID:             "profile-principal-engineer",
			Name:           "Principal Engineer",
			TargetRole:     "principal",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for principal/architect roles",
		},
		{
			ID:             "profile-engineering-manager",
			Name:           "Engineering Manager",
			TargetRole:     "em",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for engineering management roles",
		},
		{
			ID:             "profile-senior-engineer",
			Name:           "Senior Engineer",
			TargetRole:     "senior_ic",
			TargetAudience: "hiring_manager",
			Description:    "CV tailored for senior individual contributor roles",
		},
	}
}
