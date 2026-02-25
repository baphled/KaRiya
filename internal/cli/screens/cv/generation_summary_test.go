package cv_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/cv"
	"github.com/baphled/kariya/internal/cli/types"
)

var _ = Describe("GenerationSummary", func() {
	It("can be created with the expected fields", func() {
		prof := &types.CVProfile{ID: "p1", Name: "Profile 1", TargetRole: "senior_ic", TargetAudience: "hiring_manager"}

		gs := cv.GenerationSummary{
			SelectedProfile:  prof,
			SelectedAudience: "Hiring Manager",
			TechnologyFocus:  "Backend",
			Technologies:     []string{"Go", "Docker"},
			FocusArea:        "System Design",
			SkillsFormat:     "list",
			SkillsLimit:      10,
			CVLength:         "long",
			SourceEventCount: 5,
			SourceFactCount:  7,
			SectionCount:     3,
			TotalBullets:     12,
		}

		Expect(gs.SelectedProfile).To(Equal(prof))
		Expect(gs.SelectedAudience).To(Equal("Hiring Manager"))
		Expect(gs.TechnologyFocus).To(Equal("Backend"))
		Expect(gs.Technologies).To(ConsistOf("Go", "Docker"))
		Expect(gs.FocusArea).To(Equal("System Design"))
		Expect(gs.SkillsFormat).To(Equal("list"))
		Expect(gs.SkillsLimit).To(Equal(10))
		Expect(gs.CVLength).To(Equal("long"))
		Expect(gs.SourceEventCount).To(Equal(5))
		Expect(gs.SourceFactCount).To(Equal(7))
		Expect(gs.SectionCount).To(Equal(3))
		Expect(gs.TotalBullets).To(Equal(12))
	})
})
