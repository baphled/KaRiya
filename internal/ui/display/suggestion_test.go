package display_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/ui/display"
)

var _ = Describe("Suggestion display types", func() {
	Describe("BurstSuggestionFromDomain", func() {
		It("converts a fully populated burst suggestion", func() {
			suggestion := burstfact.BurstSuggestion{
				EventIDs:        []string{"event-1", "event-2"},
				ConfidenceScore: 0.91,
				Name:            "Platform Initiative",
				Description:     "Clustered platform improvement work",
			}

			Expect(display.BurstSuggestionFromDomain(suggestion)).To(Equal(display.BurstSuggestion{
				EventIDs:        []string{"event-1", "event-2"},
				ConfidenceScore: 0.91,
				Name:            "Platform Initiative",
				Description:     "Clustered platform improvement work",
			}))
		})

		It("handles zero-value input gracefully", func() {
			Expect(display.BurstSuggestionFromDomain(burstfact.BurstSuggestion{})).To(Equal(display.BurstSuggestion{}))
		})
	})

	Describe("BurstSuggestionsFromDomain", func() {
		It("converts a slice of burst suggestions", func() {
			input := []burstfact.BurstSuggestion{
				{Name: "First"},
				{Name: "Second"},
			}

			Expect(display.BurstSuggestionsFromDomain(input)).To(Equal([]display.BurstSuggestion{
				{Name: "First"},
				{Name: "Second"},
			}))
		})

		It("returns nil for nil input", func() {
			Expect(display.BurstSuggestionsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{})).To(Equal([]display.BurstSuggestion{}))
		})
	})

	Describe("SkillSuggestionFromDomain", func() {
		It("converts a fully populated skill suggestion", func() {
			suggestion := skillinference.SkillSuggestion{
				Name:       "Go",
				Category:   "backend",
				Confidence: 0.88,
				EventIDs:   []string{"event-1", "event-2"},
				Contexts:   []string{"Built services in Go", "Introduced Go tooling"},
			}

			Expect(display.SkillSuggestionFromDomain(suggestion)).To(Equal(display.SkillSuggestion{
				Name:       "Go",
				Category:   "backend",
				Confidence: 0.88,
				EventIDs:   []string{"event-1", "event-2"},
				Contexts:   []string{"Built services in Go", "Introduced Go tooling"},
			}))
		})

		It("handles zero-value input gracefully", func() {
			Expect(display.SkillSuggestionFromDomain(skillinference.SkillSuggestion{})).To(Equal(display.SkillSuggestion{}))
		})
	})

	Describe("SkillSuggestionsFromDomain", func() {
		It("converts a slice of skill suggestions", func() {
			input := []skillinference.SkillSuggestion{
				{Name: "Go"},
				{Name: "Ruby"},
			}

			Expect(display.SkillSuggestionsFromDomain(input)).To(Equal([]display.SkillSuggestion{
				{Name: "Go"},
				{Name: "Ruby"},
			}))
		})

		It("returns nil for nil input", func() {
			Expect(display.SkillSuggestionsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.SkillSuggestionsFromDomain([]skillinference.SkillSuggestion{})).To(Equal([]display.SkillSuggestion{}))
		})
	})
})
