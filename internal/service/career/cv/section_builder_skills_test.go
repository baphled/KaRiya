package cv

import (
	"context"
	"io"

	"github.com/baphled/kariya/internal/constants"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
			events  []*career.Event
			facts   []*career.Fact
		)

		BeforeEach(func() {
			bullet1 := fixtures.CVBulletWith("bullet1", "", "Built API")
			bullet1.SourceEventIDs = []string{"event1"}
			bullet1.InclusionReason = string(constants.InclusionReasonEventDirect)
			bullets = []*career.CVBullet{bullet1}

			event1 := fixtures.EventWith("event1", "Built API", "TechCorp", "")
			event1.Skills = []string{"Go", "PostgreSQL", "Docker"}
			events = []*career.Event{event1}

			facts = []*career.Fact{}
		})

		Context("Language Agnostic (no selected technologies)", func() {
			It("should create skills section with all event skills", func() {
				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil}, nil)

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
				Expect(skillsSection.Content).ToNot(BeEmpty())

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

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs}, nil)

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

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs}, nil)

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

				Expect(skillsSection.Content[0].Bullets).ToNot(BeEmpty())
			})
		})

		Context("Specialist with single technology", func() {
			It("should prominently feature the selected technology", func() {
				selectedTechs := []string{"Go"}

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs}, nil)

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
				event2 := fixtures.EventWith("event2", "Developed service", "StartupCo", "")
				event2.Skills = []string{"Go", "Redis"}
				events = append(events, event2)

				bullet2 := fixtures.CVBulletWith("bullet2", "", "Developed service")
				bullet2.SourceEventIDs = []string{"event2"}
				bullet2.InclusionReason = string(constants.InclusionReasonEventDirect)
				bullets = append(bullets, bullet2)
			})

			It("should show unique skills without duplicates", func() {
				selectedTechs := []string{"Go"}

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs}, nil)

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
				Expect(skillsSection.Content[0].Bullets).To(HaveLen(4))
			})
		})

		Context("No events have skills", func() {
			BeforeEach(func() {
				events[0].Skills = []string{}
			})

			It("should return nil skills section", func() {
				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil}, nil)

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
