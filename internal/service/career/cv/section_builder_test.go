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
		sections, err := builder.BuildSections(ctx, []*career.CVBullet{}, []*career.Event{}, []*career.Fact{}, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(sections).To(BeEmpty())
	})

	It("should create experience section from bullets", func() {
		bullet1 := fixtures.CVBulletWithSources("bullet1", "", "Implemented authentication system", []string{"event1"}, []string{})
		bullet1.Rank = 0.9
		bullets := []*career.CVBullet{bullet1}

		event := fixtures.EventWith("event1", "Implemented authentication system", "TechCorp", "")
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(sections).ToNot(BeEmpty())

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
		Expect(expSection.Content).ToNot(BeEmpty())
		Expect(expSection.Content[0].Header).To(Equal("TechCorp"))
		Expect(expSection.Content[0].Bullets).To(HaveLen(1))
		Expect(expSection.Content[0].Bullets[0].Text).To(Equal("Implemented authentication system"))
	})

	It("should create skills section from event skills (Phase 11 - Task 40)", func() {
		bullets := []*career.CVBullet{
			fixtures.CVBulletWithSources("bullet1", "", "Built API", []string{"event1"}, nil),
		}

		event := fixtures.Event("event1")
		event.Skills = []string{"Go", "PostgreSQL"}
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(sections).ToNot(BeEmpty())

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
		Expect(skillsSection.Content).ToNot(BeEmpty())
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
			fixtures.CVBulletWithSources("bullet1", "", "Implemented feature", []string{"event1"}, []string{}),
		}

		event := fixtures.EventWith("event1", "Implemented feature", "TechCorp", "")
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
		Expect(err).NotTo(HaveOccurred())

		// Should only have experience section, no skills
		for _, section := range sections {
			Expect(section.SectionType).NotTo(Equal("skills"))
		}
	})

	It("should create summary section for principal role", func() {
		bullet1 := fixtures.CVBulletWithSources("bullet1", "", "Led architecture implementation", []string{"event1"}, []string{})
		bullet1.Rank = 0.9
		bullets := []*career.CVBullet{bullet1}

		event := fixtures.EventWith("event1", "Led architecture implementation", "TechCorp", "")
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
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
		// Summary should contain the actual bullet text, not a template
		Expect(summarySection.Summary).To(ContainSubstring("Led architecture implementation"))
	})

	It("should create summary section for all roles", func() {
		bullets := []*career.CVBullet{
			fixtures.CVBulletWithSources("bullet1", "", "Implemented feature", []string{"event1"}, []string{}),
		}

		event := fixtures.EventWith("event1", "Implemented feature", "TechCorp", "")
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil, nil)
		Expect(err).NotTo(HaveOccurred())

		// Should include summary section for senior_ic (all roles get summary now)
		var summarySection *career.CVSection
		for _, section := range sections {
			if section.SectionType == "summary" {
				summarySection = section
				break
			}
		}
		Expect(summarySection).NotTo(BeNil())
		Expect(summarySection.Title).To(Equal("Professional Summary"))
	})

	It("should order sections correctly", func() {
		bullets := []*career.CVBullet{
			fixtures.CVBulletWithSources("bullet1", "", "Implemented feature", []string{"event1"}, []string{"fact1"}),
		}

		event := fixtures.EventWith("event1", "Implemented feature", "TechCorp", "")
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
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
		bullet2 := fixtures.CVBulletWithSources("bullet2", "", "Feature B", []string{"event2"}, []string{})
		bullet2.Rank = 0.7
		bullets := []*career.CVBullet{
			fixtures.CVBulletWithSources("bullet1", "", "Feature A", []string{"event1"}, []string{}),
			bullet2,
		}

		event1 := fixtures.EventWith("event1", "Feature A", "CompanyA", "")
		event2 := fixtures.EventWith("event2", "Feature B", "CompanyB", "")
		events := []*career.Event{event1, event2}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
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
		Expect(expSection.Content).To(HaveLen(2))

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
			fixtures.CVBulletWithSources("bullet1", "", "Feature", []string{"event1"}, nil),
		}

		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := builder.BuildSections(cancelCtx, bullets, []*career.Event{}, []*career.Fact{}, "principal", &SkillsFormatConfig{Format: "flat", Limit: 0, SelectedTechnologies: nil}, nil)
		Expect(err).To(HaveOccurred())
	})

	It("should handle events with no company", func() {
		bullets := []*career.CVBullet{
			fixtures.CVBulletWithSources("bullet1", "", "Feature", []string{"event1"}, nil),
		}

		event := fixtures.EventWith("event1", "Feature", "", "MyProject") // Has project instead of company
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(sections).ToNot(BeEmpty())

		// Should have projects section instead of experience
		var projectsSection *career.CVSection
		for _, section := range sections {
			if section.SectionType == "projects" {
				projectsSection = section
				break
			}
		}
		Expect(projectsSection).NotTo(BeNil())
		Expect(projectsSection.Content).ToNot(BeEmpty())
		Expect(projectsSection.Content[0].Header).To(Equal("MyProject"))
		Expect(projectsSection.Content[0].Bullets[0].Text).To(Equal("Feature"))
	})

	It("should generate unique section IDs", func() {
		bullet2 := fixtures.CVBulletWithSources("bullet2", "", "Feature B", []string{}, []string{"fact1"})
		bullet2.Rank = 0.7
		bullets := []*career.CVBullet{
			fixtures.CVBulletWithSources("bullet1", "", "Feature A", []string{"event1"}, []string{}),
			bullet2,
		}

		event := fixtures.EventWith("event1", "Feature A", "TechCorp", "")
		events := []*career.Event{event}

		sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "principal", nil, nil)
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

			events := []*career.Event{eventA1, eventB, eventA2}

			// Bullets for each event
			bulletB := fixtures.CVBulletWithSources("bullet-b", "", "Achievement at B", []string{"b1"}, nil)
			bulletB.Rank = 0.7
			bulletA2 := fixtures.CVBulletWithSources("bullet-a2", "", "Achievement at A (second)", []string{"a2"}, nil)
			bulletA2.Rank = 0.9
			bullets := []*career.CVBullet{
				fixtures.CVBulletWithSources("bullet-a1", "", "Achievement at A (first)", []string{"a1"}, nil),
				bulletB,
				bulletA2,
			}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil, nil)
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
			Expect(expSection.Content).To(HaveLen(3))

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

			events := []*career.Event{eventA1, eventF, eventA2}

			bulletB2 := fixtures.CVBulletWithSources("b2", "", "Freelance", []string{"f1"}, nil)
			bulletB2.Rank = 0.7
			bulletB3 := fixtures.CVBulletWithSources("b3", "", "Work 2", []string{"a2"}, nil)
			bulletB3.Rank = 0.9
			bullets := []*career.CVBullet{
				fixtures.CVBulletWithSources("b1", "", "Work 1", []string{"a1"}, nil),
				bulletB2,
				bulletB3,
			}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil, nil)
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
			Expect(expSection.Content).To(HaveLen(3))
		})

		It("should order experience entries by most recent first", func() {
			// Events at different times
			eventOld := fixtures.EventWith("old", "Old work", "OldCorp", "")
			eventOld.Date = time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)

			eventNew := fixtures.EventWith("new", "New work", "NewCorp", "")
			eventNew.Date = time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)

			events := []*career.Event{eventOld, eventNew}

			bulletB2 := fixtures.CVBulletWithSources("b2", "", "New achievement", []string{"new"}, nil)
			bulletB2.Rank = 0.9
			bullets := []*career.CVBullet{
				fixtures.CVBulletWithSources("b1", "", "Old achievement", []string{"old"}, nil),
				bulletB2,
			}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil, nil)
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

			events := []*career.Event{beisEvent1, beisEvent2, wafEvent}

			// Bullet references all three events (cross-company SourceEventIDs from BUG-013).
			bullet := fixtures.CVBulletWithSources("bullet-cross", "", "Cross-company bullet", []string{"beis1", "beis2", "waf1"}, nil)
			bullets := []*career.CVBullet{bullet}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil, nil)
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

			events := []*career.Event{eventA1, eventA2, eventA3, eventB, eventC}

			bullet := fixtures.CVBulletWithSources("bullet-multi", "", "Multi-company bullet", []string{"a1", "a2", "a3", "b1", "c1"}, nil)
			bullets := []*career.CVBullet{bullet}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil, nil)
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

			events := []*career.Event{event1, event2}

			bullet := fixtures.CVBulletWithSources("bullet-single", "", "Single-company bullet", []string{"e1", "e2"}, nil)
			bullets := []*career.CVBullet{bullet}

			sections, err := builder.BuildSections(ctx, bullets, events, []*career.Fact{}, "senior_ic", nil, nil)
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

var _ = Describe("computeExperienceYears", func() {
	It("should return 0 for empty events", func() {
		years := computeExperienceYears([]*career.Event{})
		Expect(years).To(Equal(0))
	})

	It("should return 0 for nil events", func() {
		years := computeExperienceYears(nil)
		Expect(years).To(Equal(0))
	})

	It("should compute years from single event to now", func() {
		event := fixtures.Event("e1")
		event.Date = time.Now().AddDate(-5, 0, 0) // 5 years ago
		years := computeExperienceYears([]*career.Event{event})
		Expect(years).To(BeNumerically(">=", 4))
		Expect(years).To(BeNumerically("<=", 6))
	})

	It("should compute years from earliest to latest event when latest is in the past", func() {
		event1 := fixtures.Event("e1")
		event1.Date = time.Date(2015, 1, 1, 0, 0, 0, 0, time.UTC)

		event2 := fixtures.Event("e2")
		event2.Date = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

		// Latest is in the past, so it uses now as end date
		years := computeExperienceYears([]*career.Event{event1, event2})
		// From 2015 to now (~11 years as of 2026)
		Expect(years).To(BeNumerically(">=", 10))
	})

	It("should use min/max dates, not sum durations", func() {
		// Timeline: 2018, 2019, 2020 - overlapping events
		event1 := fixtures.Event("e1")
		event1.Date = time.Date(2018, 6, 1, 0, 0, 0, 0, time.UTC)

		event2 := fixtures.Event("e2")
		event2.Date = time.Date(2019, 3, 1, 0, 0, 0, 0, time.UTC)

		event3 := fixtures.Event("e3")
		event3.Date = time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC)

		years := computeExperienceYears([]*career.Event{event1, event2, event3})
		// From mid-2018 to now (~7-8 years)
		Expect(years).To(BeNumerically(">=", 5))
		Expect(years).To(BeNumerically("<=", 10))
	})
})

var _ = Describe("buildSummarySection with SummaryConfig", func() {
	var (
		builder *DefaultSectionBuilder
		log     *logger.Logger
	)

	BeforeEach(func() {
		log = logger.DefaultLogger()
		builder = NewSectionBuilder(nil, log)
	})

	It("should render prose from top 4 bullets only", func() {
		bullets := []*career.CVBullet{
			fixtures.CVBulletWith("b1", "", "First achievement"),
			fixtures.CVBulletWith("b2", "", "Second achievement"),
			fixtures.CVBulletWith("b3", "", "Third achievement"),
			fixtures.CVBulletWith("b4", "", "Fourth achievement"),
			fixtures.CVBulletWith("b5", "", "Fifth achievement"),
		}
		for i, b := range bullets {
			b.SourceEventIDs = []string{"e1"}
			b.Rank = float64(5-i) / 5.0
		}

		section := builder.buildSummarySection(bullets, 0, nil, 10)

		Expect(section).NotTo(BeNil())
		// Should contain only top 4 bullets joined with space
		Expect(section.Summary).To(ContainSubstring("First achievement"))
		Expect(section.Summary).To(ContainSubstring("Second achievement"))
		Expect(section.Summary).To(ContainSubstring("Third achievement"))
		Expect(section.Summary).To(ContainSubstring("Fourth achievement"))
		// Should NOT contain 5th bullet
		// Should use period separation and end with period
		Expect(section.Summary).To(Equal("First achievement. Second achievement. Third achievement. Fourth achievement."))
	})

	It("should render heading template with Title and Years", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Some achievement")
		bullet.SourceEventIDs = []string{"e1"}

		summaryCfg := &SummaryConfig{
			SummaryHeading: "**{{.Title}} | {{.Years}}+ Years Experience**",
			ProfileTitle:   "Senior Ruby Developer",
		}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, summaryCfg, 15)

		Expect(section).NotTo(BeNil())
		Expect(section.Summary).To(HavePrefix("**Senior Ruby Developer | 15+ Years Experience**"))
		Expect(section.Summary).To(ContainSubstring("\n"))
		Expect(section.Summary).To(ContainSubstring("Some achievement"))
	})

	It("should return prose only when SummaryHeading is empty", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Some achievement")
		bullet.SourceEventIDs = []string{"e1"}

		summaryCfg := &SummaryConfig{
			SummaryHeading: "",
			ProfileTitle:   "",
		}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, summaryCfg, 10)

		Expect(section).NotTo(BeNil())
		// Should NOT have a newline at the start (no heading)
		Expect(section.Summary).NotTo(HavePrefix("\n"))
		Expect(section.Summary).To(Equal("Some achievement."))
	})

	It("should auto-generate heading from ProfileTitle when SummaryHeading is empty", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Built production systems serving millions of users")
		bullet.SourceEventIDs = []string{"e1"}

		summaryCfg := &SummaryConfig{
			SummaryHeading: "", // empty — no template
			ProfileTitle:   "Staff Software Engineer",
		}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, summaryCfg, 20)

		Expect(section).NotTo(BeNil())
		Expect(section.Summary).To(HavePrefix("**Staff Software Engineer | 20+ Years Experience**"))
		Expect(section.Summary).To(ContainSubstring("\n"))
		Expect(section.Summary).To(ContainSubstring("Built production systems"))
	})

	It("should gracefully handle invalid template syntax", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Achievement text")
		bullet.SourceEventIDs = []string{"e1"}

		summaryCfg := &SummaryConfig{
			SummaryHeading: "**{{.InvalidSyntax", // Missing closing braces
			ProfileTitle:   "Title",
		}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, summaryCfg, 10)

		Expect(section).NotTo(BeNil())
		// Should fall back to prose only
		Expect(section.Summary).To(Equal("Achievement text."))
	})

	It("should gracefully handle template execution errors", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Achievement text")
		bullet.SourceEventIDs = []string{"e1"}

		summaryCfg := &SummaryConfig{
			// Valid syntax but references non-existent field
			SummaryHeading: "{{.NonExistentField}}",
			ProfileTitle:   "Title",
		}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, summaryCfg, 10)

		Expect(section).NotTo(BeNil())
		// Template execution with missing field produces empty string, not error
		// So we get empty heading which is treated as no heading
		Expect(section.Summary).To(ContainSubstring("Achievement text."))
	})

	It("should return nil for empty bullets", func() {
		section := builder.buildSummarySection([]*career.CVBullet{}, 0, nil, 10)
		Expect(section).To(BeNil())
	})

	It("should strip bullet markers from text", func() {
		bullets := []*career.CVBullet{
			fixtures.CVBulletWith("b1", "", "- Bullet with dash"),
			fixtures.CVBulletWith("b2", "", "* Bullet with asterisk"),
		}
		for _, b := range bullets {
			b.SourceEventIDs = []string{"e1"}
		}

		section := builder.buildSummarySection(bullets, 0, nil, 5)

		Expect(section).NotTo(BeNil())
		Expect(section.Summary).NotTo(ContainSubstring("- Bullet"))
		Expect(section.Summary).NotTo(ContainSubstring("* Bullet"))
		Expect(section.Summary).To(ContainSubstring("Bullet with dash"))
		Expect(section.Summary).To(ContainSubstring("Bullet with asterisk"))
	})

	It("should ensure prose ends with period", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Achievement without period")
		bullet.SourceEventIDs = []string{"e1"}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, nil, 5)

		Expect(section).NotTo(BeNil())
		Expect(section.Summary).To(HaveSuffix("."))
	})

	It("should not double period if bullet already ends with period", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Achievement with period.")
		bullet.SourceEventIDs = []string{"e1"}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, nil, 5)

		Expect(section).NotTo(BeNil())
		Expect(section.Summary).NotTo(HaveSuffix(".."))
		Expect(section.Summary).To(Equal("Achievement with period."))
	})

})

var _ = Describe("firstSentence", func() {
	It("extracts first sentence from multi-sentence text", func() {
		text := "Migrated backend to Node.js. This reduced costs by 70% and improved response times."
		Expect(firstSentence(text, 150)).To(Equal("Migrated backend to Node.js."))
	})

	It("returns full text when no sentence boundary exists and under maxChars", func() {
		text := "Short achievement without period boundary"
		Expect(firstSentence(text, 150)).To(Equal("Short achievement without period boundary."))
	})

	It("truncates at word boundary when first sentence exceeds maxChars", func() {
		text := "Migrated QuikCV backend from Ruby on Rails to Node.js reducing server costs by 70% and improving response times by 40% through async processing and connection pooling optimisations across the entire platform infrastructure"
		result := firstSentence(text, 150)
		Expect(len(result)).To(BeNumerically("<=", 151))
		Expect(result).To(HaveSuffix("."))
		Expect(result).To(HavePrefix("Migrated QuikCV backend"))
	})

	It("returns empty string for empty input", func() {
		Expect(firstSentence("", 150)).To(Equal(""))
	})

	It("handles text that is exactly at maxChars", func() {
		text := "Exactly at limit"
		Expect(firstSentence(text, 200)).To(Equal("Exactly at limit."))
	})

	It("handles text with trailing punctuation", func() {
		text := "Achievement with period."
		Expect(firstSentence(text, 150)).To(Equal("Achievement with period."))
	})

	It("extracts first sentence even when it has trailing punctuation", func() {
		text := "First sentence here. Second sentence follows."
		Expect(firstSentence(text, 150)).To(Equal("First sentence here."))
	})
})

var _ = Describe("buildSummarySection wall-of-text regression", func() {
	var (
		builder *DefaultSectionBuilder
		log     *logger.Logger
	)

	BeforeEach(func() {
		log = logger.DefaultLogger()
		builder = NewSectionBuilder(nil, log)
	})

	It("produces short prose from long multi-sentence bullets (top 4)", func() {
		bullets := []*career.CVBullet{
			fixtures.CVBulletWith("b1", "", "Migrated QuikCV backend from Ruby on Rails to Node.js, reducing server costs by 70% and improving response times by 40% through async processing and connection pooling optimisations across the entire platform infrastructure. This was a major undertaking that took six months."),
			fixtures.CVBulletWith("b2", "", "Founded n-vyro.io IoT platform, delivering production-ready firmware in C/C++ and backend services in Go and Node.js for real-time device control. The platform served thousands of connected devices across multiple regions."),
			fixtures.CVBulletWith("b3", "", "Adopted Jest early for QuikCV testing, establishing a test-first culture that reduced regression bugs by 60% across the engineering team. This approach was later adopted company-wide as the standard testing methodology."),
			fixtures.CVBulletWith("b4", "", "Built AI-powered customer service dashboards at Digital Genius integrating ML APIs for intelligent response routing and ticket classification. The system processed millions of customer interactions daily."),
		}
		for i, b := range bullets {
			b.SourceEventIDs = []string{"e1"}
			b.Rank = float64(4-i) / 4.0
		}

		section := builder.buildSummarySection(bullets, 0, nil, 10)

		Expect(section).NotTo(BeNil())
		Expect(len(section.Summary)).To(BeNumerically("<", 700), "summary must be under 700 chars with 4 bullets, got: "+section.Summary)
		Expect(section.Summary).To(ContainSubstring("Migrated QuikCV backend"))
		Expect(section.Summary).To(ContainSubstring("Founded n-vyro.io"))
		Expect(section.Summary).To(ContainSubstring("Adopted Jest early"))
		Expect(section.Summary).To(ContainSubstring("Built AI-powered customer service dashboards"))
		Expect(section.Summary).NotTo(ContainSubstring("This was a major undertaking"))
		Expect(section.Summary).NotTo(ContainSubstring("The platform served"))
		Expect(section.Summary).NotTo(ContainSubstring("This approach was later"))
		Expect(section.Summary).NotTo(ContainSubstring("The system processed"))
	})

	It("produces short prose from single long bullet without sentence boundary", func() {
		bullet := fixtures.CVBulletWith("b1", "", "Migrated QuikCV backend from Ruby on Rails to Node.js reducing server costs by 70% and improving response times by 40% through async processing and connection pooling optimisations across the entire platform infrastructure which was a significant engineering effort")
		bullet.SourceEventIDs = []string{"e1"}

		section := builder.buildSummarySection([]*career.CVBullet{bullet}, 0, nil, 5)

		Expect(section).NotTo(BeNil())
		Expect(len(section.Summary)).To(BeNumerically("<", 350), "single bullet summary must be under 350 chars with maxChars=300")
		Expect(section.Summary).To(HaveSuffix("."))
	})
})
