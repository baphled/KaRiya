package cv

import (
	"context"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultSectionBuilder", func() {
	var (
		builder *DefaultSectionBuilder
		log     *logger.Logger
		ctx     context.Context
	)

	BeforeEach(func() {
		log = logger.DefaultLogger()
		builder = NewSectionBuilder(log)
		ctx = context.Background()
	})

	It("should return empty sections for no bullets", func() {
		sections, err := builder.BuildSections(ctx, []*career.CVBullet{}, []*career.CareerEvent{}, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sections)).To(Equal(0))
	})

	It("should create experience section from bullets", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Implemented authentication system",
				SourceEventIDs: []string{"event1"},
				SourceFactIDs:  []string{},
				Rank:           0.9,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Implemented authentication system",
				Company: "TechCorp",
				Date:    time.Now(),
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sections)).To(BeNumerically(">", 0))

		// Find experience section (order changed: summary is first for principal)
		var expSection *career.CVSection
		for _, section := range sections {
			if section.SectionType == "experience" {
				expSection = section
				break
			}
		}
		Expect(expSection).NotTo(BeNil())
		Expect(expSection.Title).To(Equal("Experience"))
		Expect(len(expSection.Content)).To(BeNumerically(">", 0))
		Expect(expSection.Content[0].Header).To(Equal("TechCorp"))
		Expect(len(expSection.Content[0].Bullets)).To(Equal(1))
		Expect(expSection.Content[0].Bullets[0].Text).To(Equal("Implemented authentication system"))
	})

	It("should create skills section from fact-based bullets", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Improved system performance",
				SourceEventIDs: []string{},
				SourceFactIDs:  []string{"fact1"},
				Rank:           0.8,
			},
		}

		facts := []*career.Fact{
			{
				ID:                   "fact1",
				CompetencyCategories: []string{"Performance Optimization", "System Design"},
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, []*career.CareerEvent{}, facts, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sections)).To(BeNumerically(">", 0))

		// Should have skills section
		var skillsSection *career.CVSection
		for _, section := range sections {
			if section.SectionType == "skills" {
				skillsSection = section
				break
			}
		}
		Expect(skillsSection).NotTo(BeNil())
		Expect(skillsSection.Title).To(Equal("Core Competencies"))
		Expect(len(skillsSection.Content)).To(BeNumerically(">", 0))
		// Skills section should have competency categories, not bullet text
		hasCategory := false
		for _, group := range skillsSection.Content {
			for _, bullet := range group.Bullets {
				if bullet.Text == "Performance Optimization" || bullet.Text == "System Design" {
					hasCategory = true
					break
				}
			}
		}
		Expect(hasCategory).To(BeTrue())
	})

	It("should not create skills section without fact-based bullets", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Implemented feature",
				SourceEventIDs: []string{"event1"},
				SourceFactIDs:  []string{},
				Rank:           0.8,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Implemented feature",
				Company: "TechCorp",
				Date:    time.Now(),
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())

		// Should only have experience section, no skills
		for _, section := range sections {
			Expect(section.SectionType).NotTo(Equal("skills"))
		}
	})

	It("should create summary section for principal role", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Led architecture implementation",
				SourceEventIDs: []string{"event1"},
				SourceFactIDs:  []string{},
				Rank:           0.9,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Led architecture implementation",
				Company: "TechCorp",
				Date:    time.Now(),
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())

		// Should include summary section for principal role
		var summarySection *career.CVSection
		for _, section := range sections {
			if section.SectionType == "summary" {
				summarySection = section
				break
			}
		}
		Expect(summarySection).NotTo(BeNil())
		Expect(summarySection.Title).To(Equal("Professional Summary"))
		// Summary section uses Summary field, not Content
		Expect(summarySection.Summary).To(ContainSubstring("Experienced"))
	})

	It("should not create summary section for junior roles", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Implemented feature",
				SourceEventIDs: []string{"event1"},
				SourceFactIDs:  []string{},
				Rank:           0.8,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Implemented feature",
				Company: "TechCorp",
				Date:    time.Now(),
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic")
		Expect(err).NotTo(HaveOccurred())

		// Should not include summary section for senior_ic
		for _, section := range sections {
			Expect(section.SectionType).NotTo(Equal("summary"))
		}
	})

	It("should order sections correctly", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Implemented feature",
				SourceEventIDs: []string{"event1"},
				SourceFactIDs:  []string{"fact1"},
				Rank:           0.8,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Implemented feature",
				Company: "TechCorp",
				Date:    time.Now(),
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())

		// Should have experience first, then skills, then summary
		if len(sections) > 1 {
			Expect(sections[0].Order).To(Equal(0))
			Expect(sections[1].Order).To(Equal(1))
			if len(sections) > 2 {
				Expect(sections[2].Order).To(Equal(2))
			}
		}
	})

	It("should group bullets by company", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Feature A",
				SourceEventIDs: []string{"event1"},
				SourceFactIDs:  []string{},
				Rank:           0.8,
			},
			{
				ID:             "bullet2",
				Text:           "Feature B",
				SourceEventIDs: []string{"event2"},
				SourceFactIDs:  []string{},
				Rank:           0.7,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Feature A",
				Company: "CompanyA",
				Date:    time.Now(),
			},
			{
				ID:      "event2",
				Text:    "Feature B",
				Company: "CompanyB",
				Date:    time.Now(),
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())

		// Find experience section
		var expSection *career.CVSection
		for _, section := range sections {
			if section.SectionType == "experience" {
				expSection = section
				break
			}
		}
		Expect(expSection).NotTo(BeNil())

		// Experience section should have two content groups (one per company)
		Expect(len(expSection.Content)).To(Equal(2))

		// Collect company headers
		companies := []string{}
		for _, group := range expSection.Content {
			companies = append(companies, group.Header)
		}

		// Should contain both companies
		Expect(companies).To(ContainElement("CompanyA"))
		Expect(companies).To(ContainElement("CompanyB"))
	})

	It("should handle context cancellation", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Feature",
				SourceEventIDs: []string{"event1"},
				Rank:           0.8,
			},
		}

		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := builder.BuildSections(cancelCtx, bullets, []*career.CareerEvent{}, []*career.Fact{}, "principal")
		Expect(err).To(HaveOccurred())
	})

	It("should handle events with no company", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Feature",
				SourceEventIDs: []string{"event1"},
				Rank:           0.8,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Feature",
				Date:    time.Now(),
				Project: "MyProject", // Has project instead of company
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sections)).To(BeNumerically(">", 0))

		// Should have projects section instead of experience
		var projectsSection *career.CVSection
		for _, section := range sections {
			if section.SectionType == "projects" {
				projectsSection = section
				break
			}
		}
		Expect(projectsSection).NotTo(BeNil())
		Expect(len(projectsSection.Content)).To(BeNumerically(">", 0))
		Expect(projectsSection.Content[0].Header).To(Equal("MyProject"))
		Expect(projectsSection.Content[0].Bullets[0].Text).To(Equal("Feature"))
	})

	It("should generate unique section IDs", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Feature A",
				SourceEventIDs: []string{"event1"},
				SourceFactIDs:  []string{},
				Rank:           0.8,
			},
			{
				ID:             "bullet2",
				Text:           "Feature B",
				SourceEventIDs: []string{},
				SourceFactIDs:  []string{"fact1"},
				Rank:           0.7,
			},
		}

		events := []*career.CareerEvent{
			{
				ID:      "event1",
				Text:    "Feature A",
				Company: "TechCorp",
				Date:    time.Now(),
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal")
		Expect(err).NotTo(HaveOccurred())

		// All section IDs should be unique
		ids := make(map[string]bool)
		for _, section := range sections {
			Expect(ids[section.ID]).To(BeFalse())
			ids[section.ID] = true
		}
	})
})
