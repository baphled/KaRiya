package modals_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/screens/timeline/modals"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/domain/career"
)

var _ = Describe("Helpers", func() {
	var theme themes.Theme

	BeforeEach(func() {
		theme = themes.NewDefaultTheme()
	})

	Describe("RenderOverlayModal", func() {
		It("renders modal over background", func() {
			modalView := "Modal Content"
			backgroundView := "Background Content"

			result := modals.RenderOverlayModal(modalView, backgroundView)

			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Modal Content"))
		})

		It("handles empty modal view", func() {
			result := modals.RenderOverlayModal("", "Background")

			Expect(result).NotTo(BeEmpty())
		})

		It("handles empty background view", func() {
			result := modals.RenderOverlayModal("Modal", "")

			Expect(result).NotTo(BeEmpty())
			Expect(result).To(ContainSubstring("Modal"))
		})
	})

	Describe("RenderEventDetailContent", func() {
		It("renders nil event as no event selected", func() {
			result := modals.RenderEventDetailContent(nil, theme)

			Expect(result).To(Equal("No event selected."))
		})

		It("renders event with all fields", func() {
			event := &career.CareerEvent{
				ID:         "test-id",
				Text:       "Test event description",
				Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
				Company:    "Test Company",
				Project:    "Test Project",
				Tags:       []string{"tag1", "tag2"},
				Categories: []string{"category1"},
				Skills:     []string{"skill1", "skill2", "skill3"},
			}

			result := modals.RenderEventDetailContent(event, theme)

			Expect(result).To(ContainSubstring("Event Details"))
			Expect(result).To(ContainSubstring("2024-01-15"))
			Expect(result).To(ContainSubstring("Test Company"))
			Expect(result).To(ContainSubstring("Test Project"))
			Expect(result).To(ContainSubstring("Test event description"))
			Expect(result).To(ContainSubstring("tag1, tag2"))
			Expect(result).To(ContainSubstring("category1"))
			Expect(result).To(ContainSubstring("3 associated"))
		})

		It("renders event with minimal fields", func() {
			event := &career.CareerEvent{
				ID:   "test-id",
				Text: "Minimal event",
				Date: time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
			}

			result := modals.RenderEventDetailContent(event, theme)

			Expect(result).To(ContainSubstring("Event Details"))
			Expect(result).To(ContainSubstring("2024-01-15"))
			Expect(result).To(ContainSubstring("Minimal event"))
			Expect(result).NotTo(ContainSubstring("Company:"))
			Expect(result).NotTo(ContainSubstring("Project:"))
		})

		It("omits empty company and project", func() {
			event := &career.CareerEvent{
				ID:      "test-id",
				Text:    "Event text",
				Date:    time.Now(),
				Company: "",
				Project: "",
			}

			result := modals.RenderEventDetailContent(event, theme)

			Expect(result).NotTo(ContainSubstring("Company:"))
			Expect(result).NotTo(ContainSubstring("Project:"))
		})
	})

	Describe("RenderSkillsContent", func() {
		It("renders empty skills message", func() {
			result := modals.RenderSkillsContent([]*career.Skill{}, theme)

			Expect(result).To(ContainSubstring("No skills associated"))
		})

		It("renders nil skills slice as empty", func() {
			result := modals.RenderSkillsContent(nil, theme)

			Expect(result).To(ContainSubstring("No skills associated"))
		})

		It("renders single skill", func() {
			years := 3
			skills := []*career.Skill{
				{
					ID:        "skill-1",
					Name:      "Go",
					Category:  "Programming",
					Level:     "Expert",
					YearsUsed: &years,
				},
			}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("Programming"))
			Expect(result).To(ContainSubstring("Expert"))
			Expect(result).To(ContainSubstring("3 years"))
		})

		It("renders single year correctly", func() {
			years := 1
			skills := []*career.Skill{
				{
					ID:        "skill-1",
					Name:      "Rust",
					YearsUsed: &years,
				},
			}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("1 year"))
			Expect(result).NotTo(ContainSubstring("1 years"))
		})

		It("renders multiple skills", func() {
			skills := []*career.Skill{
				{ID: "1", Name: "Go", Category: "Backend"},
				{ID: "2", Name: "React", Category: "Frontend"},
				{ID: "3", Name: "PostgreSQL", Category: "Database"},
			}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("React"))
			Expect(result).To(ContainSubstring("PostgreSQL"))
		})

		It("skips nil skills in slice", func() {
			skills := []*career.Skill{
				{ID: "1", Name: "Go"},
				nil,
				{ID: "2", Name: "Python"},
			}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Go"))
			Expect(result).To(ContainSubstring("Python"))
		})

		It("handles skill with nil years", func() {
			skills := []*career.Skill{
				{
					ID:        "skill-1",
					Name:      "Docker",
					Category:  "DevOps",
					Level:     "Intermediate",
					YearsUsed: nil,
				},
			}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Docker"))
			Expect(result).NotTo(ContainSubstring("year"))
		})

		It("handles skill with zero years", func() {
			years := 0
			skills := []*career.Skill{
				{
					ID:        "skill-1",
					Name:      "Kubernetes",
					YearsUsed: &years,
				},
			}

			result := modals.RenderSkillsContent(skills, theme)

			Expect(result).To(ContainSubstring("Kubernetes"))
			Expect(result).NotTo(ContainSubstring("0 year"))
		})
	})
})
