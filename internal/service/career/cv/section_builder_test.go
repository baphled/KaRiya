package cv

import (
	"context"
	"strings"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
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
		builder = NewSectionBuilder(nil, log)
		ctx = context.Background()
	})

	It("should return empty sections for no bullets", func() {
		sections, err := builder.BuildSections(ctx, []*career.CVBullet{}, []*career.CareerEvent{}, []*career.Fact{}, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil})
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

		event := fixtures.EventWith("event1", "Implemented authentication system", "TechCorp", "")
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
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

	It("should create skills section from event skills (Phase 11 - Task 40)", func() {
		bullets := []*career.CVBullet{
			{
				ID:             "bullet1",
				Text:           "Built API",
				SourceEventIDs: []string{"event1"},
				Rank:           0.8,
			},
		}

		event := fixtures.Event("event1")
		event.Skills = []string{"Go", "PostgreSQL"}
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
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
		Expect(skillsSection.Title).To(Equal("Technical Skills"))
		Expect(len(skillsSection.Content)).To(BeNumerically(">", 0))
		// Skills section should have technical skills with counts
		hasSkill := false
		for _, group := range skillsSection.Content {
			for _, bullet := range group.Bullets {
				if strings.Contains(bullet.Text, "Go") || strings.Contains(bullet.Text, "PostgreSQL") {
					hasSkill = true
					break
				}
			}
		}
		Expect(hasSkill).To(BeTrue())
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

		event := fixtures.EventWith("event1", "Implemented feature", "TechCorp", "")
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
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

		event := fixtures.EventWith("event1", "Led architecture implementation", "TechCorp", "")
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
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

		event := fixtures.EventWith("event1", "Implemented feature", "TechCorp", "")
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil)
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

		event := fixtures.EventWith("event1", "Implemented feature", "TechCorp", "")
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
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

		event1 := fixtures.EventWith("event1", "Feature A", "CompanyA", "")
		event2 := fixtures.EventWith("event2", "Feature B", "CompanyB", "")
		events := []*career.CareerEvent{event1, event2}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
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

		_, err := builder.BuildSections(cancelCtx, bullets, []*career.CareerEvent{}, []*career.Fact{}, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil})
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

		event := fixtures.EventWith("event1", "Feature", "", "MyProject") // Has project instead of company
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
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

		event := fixtures.EventWith("event1", "Feature A", "TechCorp", "")
		events := []*career.CareerEvent{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil)
		Expect(err).NotTo(HaveOccurred())

		// All section IDs should be unique
		ids := make(map[string]bool)
		for _, section := range sections {
			Expect(ids[section.ID]).To(BeFalse())
			ids[section.ID] = true
		}
	})

	// BUG-009: Tenure-aware bullet grouping tests
	Context("tenure-aware bullet grouping", func() {
		It("should create separate experience entries for separate tenures", func() {
			// Timeline: Company A (Jan 2022) -> Company B (Jun 2022) -> Company A (Jan 2023)

			// First tenure at Company A
			eventA1 := fixtures.EventWith("a1", "First work at A", "Company A", "")
			eventA1.Date = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)

			// Company B in between
			eventB := fixtures.EventWith("b1", "Work at B", "Company B", "")
			eventB.Date = time.Date(2022, 6, 15, 0, 0, 0, 0, time.UTC)

			// Second tenure at Company A
			eventA2 := fixtures.EventWith("a2", "Back at A", "Company A", "")
			eventA2.Date = time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)

			events := []*career.CareerEvent{eventA1, eventB, eventA2}

			// Bullets for each event
			bullets := []*career.CVBullet{
				{ID: "bullet-a1", Text: "Achievement at A (first)", SourceEventIDs: []string{"a1"}, Rank: 0.8},
				{ID: "bullet-b", Text: "Achievement at B", SourceEventIDs: []string{"b1"}, Rank: 0.7},
				{ID: "bullet-a2", Text: "Achievement at A (second)", SourceEventIDs: []string{"a2"}, Rank: 0.9},
			}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil)
			Expect(err).NotTo(HaveOccurred())

			// Find experience section
			var expSection *career.CVSection
			for _, s := range sections {
				if s.SectionType == "experience" {
					expSection = s
					break
				}
			}
			Expect(expSection).NotTo(BeNil())

			// Should have 3 content groups: Company A (tenure 1), Company B, Company A (tenure 2)
			Expect(len(expSection.Content)).To(Equal(3))

			// Count Company A entries
			companyACount := 0
			for _, group := range expSection.Content {
				if group.Header == "Company A" {
					companyACount++
				}
			}
			Expect(companyACount).To(Equal(2))
		})

		It("should render separate entries when Freelance separates tenures", func() {
			// Company A -> Freelance -> Company A

			eventA1 := fixtures.EventWith("a1", "First at A", "Company A", "")
			eventA1.Date = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)

			eventF := fixtures.EventWith("f1", "Freelance work", "Freelance", "")
			eventF.Date = time.Date(2022, 6, 15, 0, 0, 0, 0, time.UTC)

			eventA2 := fixtures.EventWith("a2", "Back at A", "Company A", "")
			eventA2.Date = time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)

			events := []*career.CareerEvent{eventA1, eventF, eventA2}

			bullets := []*career.CVBullet{
				{ID: "b1", Text: "Work 1", SourceEventIDs: []string{"a1"}, Rank: 0.8},
				{ID: "b2", Text: "Freelance", SourceEventIDs: []string{"f1"}, Rank: 0.7},
				{ID: "b3", Text: "Work 2", SourceEventIDs: []string{"a2"}, Rank: 0.9},
			}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil)
			Expect(err).NotTo(HaveOccurred())

			// Find experience section
			var expSection *career.CVSection
			for _, s := range sections {
				if s.SectionType == "experience" {
					expSection = s
					break
				}
			}
			Expect(expSection).NotTo(BeNil())

			// Should have 3 groups
			Expect(len(expSection.Content)).To(Equal(3))
		})

		It("should order experience entries by most recent first", func() {
			// Events at different times
			eventOld := fixtures.EventWith("old", "Old work", "OldCorp", "")
			eventOld.Date = time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)

			eventNew := fixtures.EventWith("new", "New work", "NewCorp", "")
			eventNew.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			events := []*career.CareerEvent{eventOld, eventNew}

			bullets := []*career.CVBullet{
				{ID: "b1", Text: "Old achievement", SourceEventIDs: []string{"old"}, Rank: 0.8},
				{ID: "b2", Text: "New achievement", SourceEventIDs: []string{"new"}, Rank: 0.9},
			}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil)
			Expect(err).NotTo(HaveOccurred())

			// Find experience section
			var expSection *career.CVSection
			for _, s := range sections {
				if s.SectionType == "experience" {
					expSection = s
					break
				}
			}
			Expect(expSection).NotTo(BeNil())

			// First entry should be NewCorp (most recent)
			Expect(expSection.Content[0].Header).To(Equal("NewCorp"))
			// Second entry should be OldCorp
			Expect(expSection.Content[1].Header).To(Equal("OldCorp"))
		})
	})

	// BUG-014: Date calculation in groupBulletsByCompany ignores primary company.
	Context("BUG-014: date calculation filters by primary company", func() {
		It("should compute date range using only events from the primary company", func() {
			// Bullet has SourceEventIDs from two companies:
			// - BEIS event from Jun 2019
			// - BEIS event from Dec 2019
			// - We Are Friday event from Nov 2012
			// Primary company is BEIS (2 events vs 1), so dates should be Jun 2019 - Dec 2019.

			beisEvent1 := fixtures.EventWith("beis1", "Policy work", "BEIS", "")
			beisEvent1.Date = time.Date(2019, 6, 15, 0, 0, 0, 0, time.UTC)

			beisEvent2 := fixtures.EventWith("beis2", "Delivery work", "BEIS", "")
			beisEvent2.Date = time.Date(2019, 12, 15, 0, 0, 0, 0, time.UTC)

			wafEvent := fixtures.EventWith("waf1", "Old project", "We Are Friday", "")
			wafEvent.Date = time.Date(2012, 11, 15, 0, 0, 0, 0, time.UTC)

			events := []*career.CareerEvent{beisEvent1, beisEvent2, wafEvent}

			// Bullet references all three events (cross-company SourceEventIDs from BUG-013).
			bullet := &career.CVBullet{
				ID:             "bullet-cross",
				Text:           "Cross-company bullet",
				SourceEventIDs: []string{"beis1", "beis2", "waf1"},
				Rank:           0.8,
			}
			bullets := []*career.CVBullet{bullet}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil)
			Expect(err).NotTo(HaveOccurred())

			// Find experience section.
			var expSection *career.CVSection
			for _, s := range sections {
				if s.SectionType == "experience" {
					expSection = s
					break
				}
			}
			Expect(expSection).NotTo(BeNil())

			// The bullet should be assigned to BEIS (primary company).
			Expect(expSection.Content).To(HaveLen(1))
			Expect(expSection.Content[0].Header).To(Equal("BEIS"))

			// Dates should reflect BEIS events only: Jun 2019 - Dec 2019.
			Expect(expSection.Content[0].StartDate).To(Equal("Jun 2019"))
			Expect(expSection.Content[0].EndDate).To(Equal("Dec 2019"))
		})

		It("should produce correct dates for a cross-cutting bullet with many source companies", func() {
			// Bullet with events from 3 companies: CompanyA (x3, 2023), CompanyB (x1, 2020), CompanyC (x1, 2018).
			// Primary = CompanyA. Dates should only reflect CompanyA events.

			eventA1 := fixtures.EventWith("a1", "Work A1", "CompanyA", "")
			eventA1.Date = time.Date(2023, 3, 1, 0, 0, 0, 0, time.UTC)

			eventA2 := fixtures.EventWith("a2", "Work A2", "CompanyA", "")
			eventA2.Date = time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)

			eventA3 := fixtures.EventWith("a3", "Work A3", "CompanyA", "")
			eventA3.Date = time.Date(2023, 9, 1, 0, 0, 0, 0, time.UTC)

			eventB := fixtures.EventWith("b1", "Work B", "CompanyB", "")
			eventB.Date = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

			eventC := fixtures.EventWith("c1", "Work C", "CompanyC", "")
			eventC.Date = time.Date(2018, 5, 1, 0, 0, 0, 0, time.UTC)

			events := []*career.CareerEvent{eventA1, eventA2, eventA3, eventB, eventC}

			bullet := &career.CVBullet{
				ID:             "bullet-multi",
				Text:           "Multi-company bullet",
				SourceEventIDs: []string{"a1", "a2", "a3", "b1", "c1"},
				Rank:           0.8,
			}
			bullets := []*career.CVBullet{bullet}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil)
			Expect(err).NotTo(HaveOccurred())

			// Find experience section.
			var expSection *career.CVSection
			for _, s := range sections {
				if s.SectionType == "experience" {
					expSection = s
					break
				}
			}
			Expect(expSection).NotTo(BeNil())

			// Should be assigned to CompanyA (3 events, most frequent).
			Expect(expSection.Content).To(HaveLen(1))
			Expect(expSection.Content[0].Header).To(Equal("CompanyA"))

			// Dates should be Mar 2023 - Sep 2023 (CompanyA only), NOT May 2018 - Sep 2023.
			Expect(expSection.Content[0].StartDate).To(Equal("Mar 2023"))
			Expect(expSection.Content[0].EndDate).To(Equal("Sep 2023"))
		})

		It("should handle a single-company bullet without date corruption", func() {
			// Sanity check: a bullet with events from only one company should still work.

			event1 := fixtures.EventWith("e1", "Work 1", "SingleCo", "")
			event1.Date = time.Date(2021, 4, 1, 0, 0, 0, 0, time.UTC)

			event2 := fixtures.EventWith("e2", "Work 2", "SingleCo", "")
			event2.Date = time.Date(2021, 8, 1, 0, 0, 0, 0, time.UTC)

			events := []*career.CareerEvent{event1, event2}

			bullet := &career.CVBullet{
				ID:             "bullet-single",
				Text:           "Single-company bullet",
				SourceEventIDs: []string{"e1", "e2"},
				Rank:           0.8,
			}
			bullets := []*career.CVBullet{bullet}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil)
			Expect(err).NotTo(HaveOccurred())

			var expSection *career.CVSection
			for _, s := range sections {
				if s.SectionType == "experience" {
					expSection = s
					break
				}
			}
			Expect(expSection).NotTo(BeNil())
			Expect(expSection.Content).To(HaveLen(1))
			Expect(expSection.Content[0].Header).To(Equal("SingleCo"))
			Expect(expSection.Content[0].StartDate).To(Equal("Apr 2021"))
			Expect(expSection.Content[0].EndDate).To(Equal("Aug 2021"))
		})
	})
})
