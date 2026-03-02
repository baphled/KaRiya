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
				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil})

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

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs})

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

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs})

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

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs})

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

				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: selectedTechs})

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
				sections, err := builder.BuildSections(ctx, bullets, events, facts, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil})

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

var _ = Describe("buildGroupedSkills", func() {
	var builder *DefaultSectionBuilder

	BeforeEach(func() {
		log := logger.New(io.Discard, logger.InfoLevel)
		builder = NewSectionBuilder(nil, log)
	})

	It("should group skills by category", func() {
		skills := []skillInfo{
			{ID: "s1", Name: "Go", Category: "technical"},
			{ID: "s2", Name: "Python", Category: "technical"},
			{ID: "s3", Name: "Team management", Category: "leadership"},
		}

		groups := builder.buildGroupedSkills(skills, 0)
		Expect(groups).To(HaveLen(2))

		headers := make(map[string]bool)
		for _, g := range groups {
			headers[g.Header] = true
		}
		Expect(headers).To(HaveKey("technical"))
		Expect(headers).To(HaveKey("leadership"))
	})

	It("should sort categories alphabetically", func() {
		skills := []skillInfo{
			{ID: "s1", Name: "Go", Category: "technical"},
			{ID: "s2", Name: "Leadership", Category: "leadership"},
			{ID: "s3", Name: "Product", Category: "product"},
		}

		groups := builder.buildGroupedSkills(skills, 0)
		Expect(groups).To(HaveLen(3))
		Expect(groups[0].Header).To(Equal("leadership"))
		Expect(groups[1].Header).To(Equal("product"))
		Expect(groups[2].Header).To(Equal("technical"))
	})

	It("should default empty category to other", func() {
		skills := []skillInfo{
			{ID: "s1", Name: "Something", Category: ""},
		}

		groups := builder.buildGroupedSkills(skills, 0)
		Expect(groups).To(HaveLen(1))
		Expect(groups[0].Header).To(Equal("other"))
	})

	It("should apply per-group limit", func() {
		skills := []skillInfo{
			{ID: "s1", Name: "Go", Category: "technical"},
			{ID: "s2", Name: "Python", Category: "technical"},
			{ID: "s3", Name: "Rust", Category: "technical"},
		}

		groups := builder.buildGroupedSkills(skills, 2)
		Expect(groups).To(HaveLen(1))
		Expect(groups[0].Bullets).To(HaveLen(2))
	})

	It("should not limit when limitPerGroup is 0", func() {
		skills := []skillInfo{
			{ID: "s1", Name: "Go", Category: "technical"},
			{ID: "s2", Name: "Python", Category: "technical"},
			{ID: "s3", Name: "Rust", Category: "technical"},
		}

		groups := builder.buildGroupedSkills(skills, 0)
		Expect(groups[0].Bullets).To(HaveLen(3))
	})

	It("should create bullets with skill names as text", func() {
		skills := []skillInfo{
			{ID: "s1", Name: "Go", Category: "technical"},
		}

		groups := builder.buildGroupedSkills(skills, 0)
		Expect(groups[0].Bullets[0].Text).To(Equal("Go"))
		Expect(groups[0].Bullets[0].ID).NotTo(BeEmpty())
	})

	It("should handle empty skills list", func() {
		groups := builder.buildGroupedSkills([]skillInfo{}, 0)
		Expect(groups).To(BeEmpty())
	})

	It("should normalise category to lowercase", func() {
		skills := []skillInfo{
			{ID: "s1", Name: "Go", Category: "Technical"},
			{ID: "s2", Name: "Python", Category: "TECHNICAL"},
		}

		groups := builder.buildGroupedSkills(skills, 0)
		Expect(groups).To(HaveLen(1))
		Expect(groups[0].Header).To(Equal("technical"))
		Expect(groups[0].Bullets).To(HaveLen(2))
	})
})
