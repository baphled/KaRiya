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
		sections, err := builder.BuildSections(ctx, []*career.CVBullet{}, []*career.CareerEvent{}, "principal")
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
				Rank:           0.8,
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

		sections, err := builder.BuildSections(ctx, bullets, events, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sections)).To(BeNumerically(">", 0))

		// First section should be experience
		Expect(sections[0].SectionType).To(Equal("experience"))
		Expect(sections[0].Title).To(Equal("Experience"))
		Expect(sections[0].Content).To(ContainSubstring("Implemented authentication system"))
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

		sections, err := builder.BuildSections(ctx, bullets, []*career.CareerEvent{}, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sections)).To(BeNumerically(">", 0))

		// Should have skills section
		hasSkillsSection := false
		for _, section := range sections {
			if section.SectionType == "skills" {
				hasSkillsSection = true
				Expect(section.Title).To(Equal("Core Competencies"))
				Expect(section.Content).To(ContainSubstring("Improved system performance"))
			}
		}
		Expect(hasSkillsSection).To(BeTrue())
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

		sections, err := builder.BuildSections(ctx, bullets, events, "principal")
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

		sections, err := builder.BuildSections(ctx, bullets, events, "principal")
		Expect(err).NotTo(HaveOccurred())

		// Should include summary section for principal role
		hasSummarySection := false
		for _, section := range sections {
			if section.SectionType == "summary" {
				hasSummarySection = true
				Expect(section.Title).To(Equal("Professional Summary"))
				Expect(section.Content).To(ContainSubstring("Experienced"))
			}
		}
		Expect(hasSummarySection).To(BeTrue())
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

		sections, err := builder.BuildSections(ctx, bullets, events, "senior_ic")
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

		sections, err := builder.BuildSections(ctx, bullets, events, "principal")
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

		sections, err := builder.BuildSections(ctx, bullets, events, "principal")
		Expect(err).NotTo(HaveOccurred())

		// Experience section should contain both companies
		Expect(sections[0].Content).To(ContainSubstring("CompanyA"))
		Expect(sections[0].Content).To(ContainSubstring("CompanyB"))
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

		_, err := builder.BuildSections(cancelCtx, bullets, []*career.CareerEvent{}, "principal")
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
				ID:   "event1",
				Text: "Feature",
				Date: time.Now(),
				// No Company field
			},
		}

		sections, err := builder.BuildSections(ctx, bullets, events, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(sections)).To(BeNumerically(">", 0))
		Expect(sections[0].Content).To(ContainSubstring("Feature"))
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

		sections, err := builder.BuildSections(ctx, bullets, events, "principal")
		Expect(err).NotTo(HaveOccurred())

		// All section IDs should be unique
		ids := make(map[string]bool)
		for _, section := range sections {
			Expect(ids[section.ID]).To(BeFalse())
			ids[section.ID] = true
		}
	})
})

