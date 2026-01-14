package cv

import (
	"context"
	"io"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultSectionBuilder - Skills Section", func() {
	var (
		builder *DefaultSectionBuilder
		ctx     context.Context
		log     *logger.Logger
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		builder = NewSectionBuilder(nil, log)
		ctx = context.Background()
	})

	Describe("BuildSections with selected technologies", func() {
		var (
			bullets []*career.CVBullet
			events  []*career.CareerEvent
			facts   []*career.Fact
		)

		BeforeEach(func() {
			// Create sample bullets and events
			bullets = []*career.CVBullet{
				{
					ID:              "bullet1",
					Text:            "Built API",
					SourceEventIDs:  []string{"event1"},
					InclusionReason: "event_direct",
				},
			}

			events = []*career.CareerEvent{
				{
					ID:      "event1",
					Text:    "Built API",
					Company: "TechCorp",
					Skills:  []string{"Go", "PostgreSQL", "Docker"},
				},
			}

			facts = []*career.Fact{}
		})

		Context("Language Agnostic (no selected technologies)", func() {
			It("should create skills section with all event skills", func() {
				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", nil)

				Expect(err).NotTo(HaveOccurred())

				// Find skills section
				var skillsSection *career.CVSection
				for _, section := range sections {
					if section.SectionType == "skills" {
						skillsSection = section
						break
					}
				}

				Expect(skillsSection).NotTo(BeNil())
				Expect(skillsSection.Title).To(Equal("Technical Skills"))

				// Should have bullets for all skills
				Expect(len(skillsSection.Content)).To(BeNumerically(">", 0))

				// Extract skill names from bullets
				skillNames := make([]string, 0)
				for _, group := range skillsSection.Content {
					for _, bullet := range group.Bullets {
						skillNames = append(skillNames, bullet.Text)
					}
				}

				Expect(skillNames).To(ContainElement(ContainSubstring("Go")))
				Expect(skillNames).To(ContainElement(ContainSubstring("PostgreSQL")))
				Expect(skillNames).To(ContainElement(ContainSubstring("Docker")))
			})
		})

		Context("Generalist with selected technologies", func() {
			It("should prioritize selected technologies", func() {
				selectedTechs := []string{"Go", "Docker"}

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", selectedTechs)

				Expect(err).NotTo(HaveOccurred())

				// Find skills section
				var skillsSection *career.CVSection
				for _, section := range sections {
					if section.SectionType == "skills" {
						skillsSection = section
						break
					}
				}

				Expect(skillsSection).NotTo(BeNil())

				// Get first two skills (should be selected ones)
				firstGroup := skillsSection.Content[0]
				Expect(len(firstGroup.Bullets)).To(BeNumerically(">=", 2))

				firstTwoSkills := []string{
					firstGroup.Bullets[0].Text,
					firstGroup.Bullets[1].Text,
				}

				// Selected technologies should appear first
				Expect(firstTwoSkills).To(ContainElement(ContainSubstring("Go")))
				Expect(firstTwoSkills).To(ContainElement(ContainSubstring("Docker")))
			})

			It("should show skill names without counts", func() {
				selectedTechs := []string{"Go"}

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", selectedTechs)

				Expect(err).NotTo(HaveOccurred())

				// Find skills section
				var skillsSection *career.CVSection
				for _, section := range sections {
					if section.SectionType == "skills" {
						skillsSection = section
						break
					}
				}

				Expect(skillsSection).NotTo(BeNil())

				// Skills should just show names (no counts)
				// Note: With nil skillRepo, it will use the skill ID as the name
				// In production, it would look up the actual skill name
				Expect(len(skillsSection.Content[0].Bullets)).To(BeNumerically(">", 0))
			})
		})

		Context("Specialist with single technology", func() {
			It("should prominently feature the selected technology", func() {
				selectedTechs := []string{"Go"}

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", selectedTechs)

				Expect(err).NotTo(HaveOccurred())

				// Find skills section
				var skillsSection *career.CVSection
				for _, section := range sections {
					if section.SectionType == "skills" {
						skillsSection = section
						break
					}
				}

				Expect(skillsSection).NotTo(BeNil())

				// First skill should be the selected one
				firstBullet := skillsSection.Content[0].Bullets[0]
				Expect(firstBullet.Text).To(ContainSubstring("Go"))
			})
		})

		Context("Multiple events with same skill", func() {
			BeforeEach(func() {
				events = append(events, &career.CareerEvent{
					ID:      "event2",
					Text:    "Developed service",
					Company: "StartupCo",
					Skills:  []string{"Go", "Redis"},
				})

				bullets = append(bullets, &career.CVBullet{
					ID:              "bullet2",
					Text:            "Developed service",
					SourceEventIDs:  []string{"event2"},
					InclusionReason: "event_direct",
				})
			})

			It("should show unique skills without duplicates", func() {
				selectedTechs := []string{"Go"}

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", selectedTechs)

				Expect(err).NotTo(HaveOccurred())

				// Find skills section
				var skillsSection *career.CVSection
				for _, section := range sections {
					if section.SectionType == "skills" {
						skillsSection = section
						break
					}
				}

				Expect(skillsSection).NotTo(BeNil())

				// Should have unique skills (no duplicates)
				// Count unique skill IDs across both events: Go (appears twice but counted once), PostgreSQL, Docker, Redis = 4 unique
				Expect(len(skillsSection.Content[0].Bullets)).To(Equal(4))
			})
		})

		Context("No events have skills", func() {
			BeforeEach(func() {
				events[0].Skills = []string{}
			})

			It("should return nil skills section", func() {
				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", nil)

				Expect(err).NotTo(HaveOccurred())

				// Find skills section
				var skillsSection *career.CVSection
				for _, section := range sections {
					if section.SectionType == "skills" {
						skillsSection = section
						break
					}
				}

				// Should not create skills section if no skills exist
				Expect(skillsSection).To(BeNil())
			})
		})
	})
})
