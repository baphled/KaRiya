package cv

import (
	"context"
	"sort"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/constants"
	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("DataProcessingService", func() {
	var (
		svc context.Context
		dps DataProcessingService
	)

	BeforeEach(func() {
		svc = context.Background()
		dps = NewDataProcessingService(logger.DefaultLogger())
	})

	Describe("GroupEventsByCompany", func() {
		It("should group empty events", func() {
			result, err := dps.GroupEventsByCompany(svc, []*career.Event{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should group single company events", func() {
			event1 := fixtures.EventWith("1", "Led team standup", "Acme Corp", "")
			event1.Date = time.Now().AddDate(0, -1, 0)
			event1.Tags = []string{"leadership"}

			event2 := fixtures.EventWith("2", "Implemented feature", "Acme Corp", "")
			event2.Date = time.Now().AddDate(0, -2, 0)
			event2.Tags = []string{"technical"}

			events := []*career.Event{event1, event2}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(1))
			Expect(result["Acme Corp"]).NotTo(BeNil())
			Expect(result["Acme Corp"].Company).To(Equal("Acme Corp"))
			Expect(result["Acme Corp"].EventIDs).To(HaveLen(2))
		})

		It("should group multiple companies", func() {
			event1 := fixtures.EventWith("1", "Worked at Acme", "Acme Corp", "")
			event1.Date = time.Now().AddDate(0, -1, 0)

			event2 := fixtures.EventWith("2", "Worked at TechCo", "TechCo", "")
			event2.Date = time.Now().AddDate(0, -6, 0)

			events := []*career.Event{event1, event2}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(2))
			Expect(result["Acme Corp"]).NotTo(BeNil())
			Expect(result["TechCo"]).NotTo(BeNil())
		})

		It("should handle events with no company", func() {
			event := fixtures.EventWith("1", "Personal project", "", "") // No company

			events := []*career.Event{event}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result[constants.DefaultCompanyName]).NotTo(BeNil())
		})

		It("should set correct date ranges", func() {
			now := time.Now()

			event1 := fixtures.EventWith("1", "Recent event", "Acme", "")
			event1.Date = now

			event2 := fixtures.EventWith("2", "Old event", "Acme", "")
			event2.Date = now.AddDate(-1, 0, 0)

			events := []*career.Event{event1, event2}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())

			group := result["Acme"]
			Expect(group.StartDate).To(Equal(now.AddDate(-1, 0, 0)))
			Expect(group.EndDate).To(Equal(now))
		})
	})

	Describe("ExtractAchievements", func() {
		It("should return error for cancelled context", func() {
			cancelCtx, cancel := context.WithCancel(svc)
			cancel()
			event := fixtures.EventWith("1", "Built API", "Acme", "")
			_, err := dps.ExtractAchievements(cancelCtx, event, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should extract achievement from event", func() {
			event := fixtures.EventWith("1", "Led team of 12 engineers", "Acme", "")
			event.Tags = []string{"leadership"}

			achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(achievements).To(HaveLen(1))
			Expect(achievements[0].EventID).To(Equal("1"))
			Expect(achievements[0].ActionVerb).NotTo(BeEmpty())
		})

		It("should extract metrics from achievement", func() {
			event := fixtures.EventWith("1", "Increased performance by 25% with team of 12 people", "Acme", "")

			achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(achievements[0].Metrics).NotTo(BeEmpty())
		})

		It("should extract facts as achievements", func() {
			event := fixtures.EventWith("1", "Event text", "Acme", "")
			fact := fixtures.Fact("f1", "1")
			fact.Text = "Strong leadership capability"

			achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact})
			Expect(err).NotTo(HaveOccurred())
			Expect(achievements).To(HaveLen(2)) // Event + fact
			Expect(achievements[1].FactIDs).To(ContainElement("f1"))
		})

		// BUG-008: Achievement.Category propagation tests
		Context("Category propagation", func() {
			Context("from event.Categories", func() {
				It("should use first category as primary", func() {
					event := fixtures.EventWithCategories("evt-1", "Led architecture review", []string{"leadership", "technical"})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements).To(HaveLen(1))
					Expect(achievements[0].Category).To(Equal(constants.CompetencyLeadership))
				})

				It("should handle single category", func() {
					event := fixtures.EventWithCategories("evt-1", "Built API", []string{"technical"})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyTechnical))
				})

				It("should handle empty categories slice", func() {
					event := fixtures.EventWith("evt-1", "Some work", "Acme", "")
					event.Categories = []string{}

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyCategory("")))
				})

				It("should handle nil categories slice", func() {
					event := fixtures.EventWith("evt-1", "Some work", "Acme", "")
					event.Categories = nil

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyCategory("")))
				})

				It("should normalize category to lowercase", func() {
					event := fixtures.EventWithCategories("evt-1", "Led team", []string{"LEADERSHIP"})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyLeadership))
				})

				It("should handle mixed case category", func() {
					event := fixtures.EventWithCategories("evt-1", "Built system", []string{"Technical"})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyTechnical))
				})

				It("should return empty for invalid category", func() {
					event := fixtures.EventWithCategories("evt-1", "Did stuff", []string{"invalid_category"})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyCategory("")))
				})

				It("should skip invalid first category and return empty", func() {
					// When first category is invalid, we still return empty (don't fall through to second)
					event := fixtures.EventWithCategories("evt-1", "Did stuff", []string{"invalid", "leadership"})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					// extractPrimaryCategory only looks at first category
					Expect(achievements[0].Category).To(Equal(constants.CompetencyCategory("")))
				})

				It("should handle whitespace-only category", func() {
					event := fixtures.EventWithCategories("evt-1", "Did stuff", []string{"   "})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyCategory("")))
				})

				It("should handle empty string category", func() {
					event := fixtures.EventWithCategories("evt-1", "Did stuff", []string{""})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[0].Category).To(Equal(constants.CompetencyCategory("")))
				})
			})

			Context("from fact.CompetencyCategories", func() {
				It("should use first category as primary", func() {
					event := fixtures.EventWith("evt-1", "Event text", "Acme", "")
					fact := fixtures.FactWithCategories("f1", "Mentored 5 junior engineers", "evt-1",
						[]string{"mentoring", "leadership"}, []string{"hiring_manager"})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements).To(HaveLen(2))
					Expect(achievements[1].Category).To(Equal(constants.CompetencyMentoring))
				})

				It("should handle empty categories", func() {
					event := fixtures.EventWith("evt-1", "Event text", "Acme", "")
					fact := fixtures.Fact("f1", "evt-1")
					fact.CompetencyCategories = []string{}

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[1].Category).To(Equal(constants.CompetencyCategory("")))
				})

				It("should handle nil categories", func() {
					event := fixtures.EventWith("evt-1", "Event text", "Acme", "")
					fact := fixtures.Fact("f1", "evt-1")
					fact.CompetencyCategories = nil

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[1].Category).To(Equal(constants.CompetencyCategory("")))
				})

				It("should normalize category to lowercase", func() {
					event := fixtures.EventWith("evt-1", "Event text", "Acme", "")
					fact := fixtures.FactWithCategories("f1", "Led project", "evt-1",
						[]string{"PRODUCT"}, []string{})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[1].Category).To(Equal(constants.CompetencyProduct))
				})

				It("should return empty for invalid category", func() {
					event := fixtures.EventWith("evt-1", "Event text", "Acme", "")
					fact := fixtures.FactWithCategories("f1", "Did something", "evt-1",
						[]string{"not_a_real_category"}, []string{})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements[1].Category).To(Equal(constants.CompetencyCategory("")))
				})
			})

			Context("with multiple facts", func() {
				It("should propagate different categories for different facts", func() {
					event := fixtures.EventWithCategories("evt-1", "Event", []string{"leadership"})
					fact1 := fixtures.FactWithCategories("f1", "Technical work", "evt-1", []string{"technical"}, []string{})
					fact2 := fixtures.FactWithCategories("f2", "Mentoring work", "evt-1", []string{"mentoring"}, []string{})

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact1, fact2})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements).To(HaveLen(3)) // 1 event + 2 facts

					Expect(achievements[0].Category).To(Equal(constants.CompetencyLeadership)) // from event
					Expect(achievements[1].Category).To(Equal(constants.CompetencyTechnical))  // from fact1
					Expect(achievements[2].Category).To(Equal(constants.CompetencyMentoring))  // from fact2
				})

				It("should handle mix of valid and empty categories", func() {
					event := fixtures.EventWithCategories("evt-1", "Event", []string{"leadership"})
					factWithCategory := fixtures.FactWithCategories("f1", "Technical", "evt-1", []string{"technical"}, []string{})
					factWithoutCategory := fixtures.Fact("f2", "evt-1")
					factWithoutCategory.CompetencyCategories = []string{}

					achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{factWithCategory, factWithoutCategory})
					Expect(err).NotTo(HaveOccurred())
					Expect(achievements).To(HaveLen(3))

					Expect(achievements[0].Category).To(Equal(constants.CompetencyLeadership))
					Expect(achievements[1].Category).To(Equal(constants.CompetencyTechnical))
					Expect(achievements[2].Category).To(Equal(constants.CompetencyCategory("")))
				})
			})
		})
	})

	Describe("ExtractSkills", func() {
		It("should extract empty skills", func() {
			result, err := dps.ExtractSkills(svc, []*career.Event{}, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should extract skills from event tags", func() {
			event := fixtures.EventWith("1", "Event", "Acme", "")
			event.Tags = []string{"technical", "leadership"}

			events := []*career.Event{event}

			result, err := dps.ExtractSkills(svc, events, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
		})

		It("should aggregate skill categories", func() {
			event1 := fixtures.EventWith("1", "Event", "Acme", "")
			event1.Tags = []string{"technical"}

			event2 := fixtures.EventWith("2", "Event", "Acme", "")
			event2.Tags = []string{"technical"}

			events := []*career.Event{event1, event2}

			result, err := dps.ExtractSkills(svc, events, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())

			// Should have aggregated technical skills
			Expect(result).NotTo(BeEmpty())
		})
	})

	Describe("CalculateMetrics", func() {
		It("should extract percentage metrics", func() {
			text := "Increased performance by 25%"
			metrics, err := dps.CalculateMetrics(svc, text)
			Expect(err).NotTo(HaveOccurred())
			Expect(metrics).NotTo(BeEmpty())

			found := false
			for _, m := range metrics {
				if m.Type == "percentage" && m.Value == "25" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should extract count metrics", func() {
			text := "Led team of 12 people"
			metrics, err := dps.CalculateMetrics(svc, text)
			Expect(err).NotTo(HaveOccurred())

			found := false
			for _, m := range metrics {
				if m.Type == "count" && m.Value == "12" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should extract currency metrics", func() {
			text := "Managed budget of $5M"
			metrics, err := dps.CalculateMetrics(svc, text)
			Expect(err).NotTo(HaveOccurred())

			found := false
			for _, m := range metrics {
				if m.Type == "currency" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should extract time metrics", func() {
			text := "Led project for 6 months"
			metrics, err := dps.CalculateMetrics(svc, text)
			Expect(err).NotTo(HaveOccurred())

			found := false
			for _, m := range metrics {
				if m.Type == "time" && m.Value == "6" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should extract ratio metrics", func() {
			text := "Improved performance 3x faster"
			metrics, err := dps.CalculateMetrics(svc, text)
			Expect(err).NotTo(HaveOccurred())

			found := false
			for _, m := range metrics {
				if m.Type == "ratio" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should extract multiple metrics", func() {
			text := "Led team of 12 people, improved performance by 25% in 6 months"
			metrics, err := dps.CalculateMetrics(svc, text)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(metrics)).To(BeNumerically(">", 1))
		})
	})

	Describe("ExtractProjectsFromEvents", func() {
		It("should extract no projects from events without projects", func() {
			event := fixtures.EventWith("1", "Event", "Acme", "")

			events := []*career.Event{event}

			result, err := dps.ExtractProjectsFromEvents(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})

		It("should extract single project", func() {
			event := fixtures.EventWith("1", "Event", "Acme", "ProjectX")

			events := []*career.Event{event}

			result, err := dps.ExtractProjectsFromEvents(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())

			found := false
			for _, p := range result {
				if p.Name == "ProjectX" {
					found = true
					break
				}
			}
			Expect(found).To(BeTrue())
		})

		It("should extract multiple projects", func() {
			event1 := fixtures.EventWith("1", "Event", "Acme", "ProjectX")
			event2 := fixtures.EventWith("2", "Event", "Acme", "ProjectY")

			events := []*career.Event{event1, event2}

			result, err := dps.ExtractProjectsFromEvents(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
		})

		It("should set correct project date ranges", func() {
			now := time.Now()

			event1 := fixtures.EventWith("1", "Event", "Acme", "ProjectX")
			event1.Date = now

			event2 := fixtures.EventWith("2", "Event", "Acme", "ProjectX")
			event2.Date = now.AddDate(0, -3, 0)

			events := []*career.Event{event1, event2}

			result, err := dps.ExtractProjectsFromEvents(svc, events)
			Expect(err).NotTo(HaveOccurred())

			var project *ProjectGroup
			for _, p := range result {
				if p.Name == "ProjectX" {
					project = p
					break
				}
			}

			Expect(project).NotTo(BeNil())
			Expect(project.StartDate).To(Equal(now.AddDate(0, -3, 0)))
			Expect(project.EndDate).To(Equal(now))
		})
	})

	Describe("Context handling", func() {
		It("should handle cancelled context", func() {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := dps.GroupEventsByCompany(ctx, []*career.Event{})
			Expect(err).To(HaveOccurred())
		})
	})

	// BUG-009: Tenure detection tests
	Describe("GroupEventsByCompany with tenure detection", func() {
		Context("when person worked at same company in separate periods", func() {
			It("should detect single tenure when no gaps exist", func() {
				// Company A events only - should be single tenure
				now := time.Now()

				event1 := fixtures.EventWith("1", "Work at Acme", "Acme Corp", "")
				event1.Date = now.AddDate(0, -2, 0) // 2 months ago

				event2 := fixtures.EventWith("2", "More work at Acme", "Acme Corp", "")
				event2.Date = now.AddDate(0, -1, 0) // 1 month ago

				events := []*career.Event{event1, event2}

				result, err := dps.GroupEventsByCompany(svc, events)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(HaveLen(1))
				Expect(result["Acme Corp"]).NotTo(BeNil())
				Expect(result["Acme Corp"].EventIDs).To(HaveLen(2))
			})

			It("should detect multiple tenures with intervening company", func() {
				// Timeline: Company A (Mar-Jul 2022) -> Company B (Aug 2022-Jun 2023) -> Company A (Jul 2023-Apr 2024)
				// Should create TWO separate groups for Company A

				// First tenure at Company A
				eventA1 := fixtures.EventWith("a1", "First stint at A", "Company A", "")
				eventA1.Date = time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC)

				eventA2 := fixtures.EventWith("a2", "More work at A", "Company A", "")
				eventA2.Date = time.Date(2022, 7, 15, 0, 0, 0, 0, time.UTC)

				// Company B in between
				eventB1 := fixtures.EventWith("b1", "Work at B", "Company B", "")
				eventB1.Date = time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)

				eventB2 := fixtures.EventWith("b2", "More work at B", "Company B", "")
				eventB2.Date = time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)

				// Second tenure at Company A
				eventA3 := fixtures.EventWith("a3", "Back at A", "Company A", "")
				eventA3.Date = time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)

				eventA4 := fixtures.EventWith("a4", "Still at A", "Company A", "")
				eventA4.Date = time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

				events := []*career.Event{eventA1, eventA2, eventB1, eventB2, eventA3, eventA4}

				result, err := dps.GroupEventsByCompany(svc, events)
				Expect(err).NotTo(HaveOccurred())

				// Should have 3 groups: Company A (tenure 1), Company B, Company A (tenure 2)
				Expect(result).To(HaveLen(3))

				// Check that Company A has two separate entries
				companyACount := 0
				for key := range result {
					if key == "Company A" || strings.HasPrefix(key, "Company A"+constants.TenureSeparator) {
						companyACount++
					}
				}
				Expect(companyACount).To(Equal(2))
			})

			It("should treat Freelance as tenure separator", func() {
				// Company A -> Freelance -> Company A should create two separate Company A tenures

				eventA1 := fixtures.EventWith("a1", "First at A", "Company A", "")
				eventA1.Date = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)

				eventF := fixtures.EventWith("f1", "Freelance work", "Freelance", "")
				eventF.Date = time.Date(2022, 6, 15, 0, 0, 0, 0, time.UTC)

				eventA2 := fixtures.EventWith("a2", "Back at A", "Company A", "")
				eventA2.Date = time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)

				events := []*career.Event{eventA1, eventF, eventA2}

				result, err := dps.GroupEventsByCompany(svc, events)
				Expect(err).NotTo(HaveOccurred())

				// Should have 3 groups: Company A (tenure 1), Freelance, Company A (tenure 2)
				Expect(result).To(HaveLen(3))

				// Check that Company A has two separate entries
				companyACount := 0
				for key := range result {
					if key == "Company A" || strings.HasPrefix(key, "Company A"+constants.TenureSeparator) {
						companyACount++
					}
				}
				Expect(companyACount).To(Equal(2))
			})

			It("should order tenures reverse-chronologically by end date", func() {
				// Create events with clear timeline
				eventA1 := fixtures.EventWith("a1", "First at A", "Company A", "")
				eventA1.Date = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)

				eventB := fixtures.EventWith("b1", "Work at B", "Company B", "")
				eventB.Date = time.Date(2022, 6, 15, 0, 0, 0, 0, time.UTC)

				eventA2 := fixtures.EventWith("a2", "Back at A", "Company A", "")
				eventA2.Date = time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)

				events := []*career.Event{eventA1, eventB, eventA2}

				result, err := dps.GroupEventsByCompany(svc, events)
				Expect(err).NotTo(HaveOccurred())

				// Collect all groups and sort by end date descending
				var groups []*CompanyGroup
				for _, g := range result {
					groups = append(groups, g)
				}

				// Sort by end date descending (newest first)
				sort.Slice(groups, func(i, j int) bool {
					return groups[i].EndDate.After(groups[j].EndDate)
				})

				// The most recent Company A tenure should have end date in 2023
				// The first entry should be the newest (2023 Company A tenure)
				Expect(groups[0].Company).To(Equal("Company A"))
				Expect(groups[0].EndDate.Year()).To(Equal(2023))
			})

			It("should handle adjacent same-company events as single tenure", func() {
				// Events at same company with no intervening work should be single tenure
				eventA1 := fixtures.EventWith("a1", "Work", "Company A", "")
				eventA1.Date = time.Date(2022, 1, 15, 0, 0, 0, 0, time.UTC)

				eventA2 := fixtures.EventWith("a2", "More work", "Company A", "")
				eventA2.Date = time.Date(2022, 2, 15, 0, 0, 0, 0, time.UTC)

				eventA3 := fixtures.EventWith("a3", "Even more", "Company A", "")
				eventA3.Date = time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC)

				events := []*career.Event{eventA1, eventA2, eventA3}

				result, err := dps.GroupEventsByCompany(svc, events)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(HaveLen(1))
				Expect(result["Company A"]).NotTo(BeNil())
				Expect(result["Company A"].EventIDs).To(HaveLen(3))
			})
		})
	})

	// BUG-012: Project-only events should not split tenures
	Describe("hasInterveningCompanyEvents", func() {
		Context("BUG-012: project-only events (empty Company)", func() {
			It("should return false when only project-only events exist between dates", func() {
				// Project-only event (Company="") between two dates should NOT count
				// as an intervening company event.
				projectEvent := fixtures.EventWith("p1", "n-vyro.io work", "", "n-vyro.io")
				projectEvent.Date = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

				allEvents := []*career.Event{projectEvent}
				sort.Slice(allEvents, func(i, j int) bool {
					return allEvents[i].Date.Before(allEvents[j].Date)
				})

				startDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

				result := hasInterveningCompanyEvents(
					startDate,
					endDate,
					"Mindful Chef",
					allEvents,
				)

				Expect(result).To(BeFalse())
			})

			It("should still return true when real company events intervene", func() {
				realCompanyEvent := fixtures.EventWith("c1", "Work at Other Corp", "Other Corp", "")
				realCompanyEvent.Date = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

				allEvents := []*career.Event{realCompanyEvent}
				sort.Slice(allEvents, func(i, j int) bool {
					return allEvents[i].Date.Before(allEvents[j].Date)
				})

				startDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

				result := hasInterveningCompanyEvents(
					startDate,
					endDate,
					"Mindful Chef",
					allEvents,
				)

				Expect(result).To(BeTrue())
			})

			It("should return false when mix of project-only and same-company events exist", func() {
				projectEvent := fixtures.EventWith("p1", "n-vyro.io work", "", "n-vyro.io")
				projectEvent.Date = time.Date(2024, 3, 10, 0, 0, 0, 0, time.UTC)

				sameCompanyEvent := fixtures.EventWith("mc1", "Mindful Chef work", "Mindful Chef", "")
				sameCompanyEvent.Date = time.Date(2024, 3, 20, 0, 0, 0, 0, time.UTC)

				allEvents := []*career.Event{projectEvent, sameCompanyEvent}
				sort.Slice(allEvents, func(i, j int) bool {
					return allEvents[i].Date.Before(allEvents[j].Date)
				})

				startDate := time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

				result := hasInterveningCompanyEvents(
					startDate,
					endDate,
					"Mindful Chef",
					allEvents,
				)

				Expect(result).To(BeFalse())
			})
		})
	})

	// BUG-012: Integration-level tenure detection with project-only events
	Describe("GroupEventsByCompany with project-only events", func() {
		It("BUG-012: should not split tenure when only project-only events fill the gap", func() {
			// Mindful Chef: Feb 2024 and Apr 2024 with n-vyro.io (project-only) at Mar 2024
			// Should produce ONE Mindful Chef tenure, not two.
			mcEvent1 := fixtures.EventWith("mc1", "Work at Mindful Chef", "Mindful Chef", "")
			mcEvent1.Date = time.Date(2024, 2, 15, 0, 0, 0, 0, time.UTC)

			mcEvent2 := fixtures.EventWith("mc2", "More at Mindful Chef", "Mindful Chef", "")
			mcEvent2.Date = time.Date(2024, 4, 15, 0, 0, 0, 0, time.UTC)

			// Project-only event (no company, just a project)
			projectEvent := fixtures.EventWith("p1", "n-vyro.io development", "", "n-vyro.io")
			projectEvent.Date = time.Date(2024, 3, 15, 0, 0, 0, 0, time.UTC)

			events := []*career.Event{mcEvent1, projectEvent, mcEvent2}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())

			mcCount := 0
			for _, group := range result {
				if group.Company == "Mindful Chef" {
					mcCount++
				}
			}

			// Mindful Chef should appear as ONE tenure, not split into two.
			Expect(mcCount).To(Equal(1), "Mindful Chef should be a single tenure, not split by project-only events")
		})

		It("BUG-012: should still split tenure when real company events intervene", func() {
			// Ensure the existing BUG-009 behavior is preserved:
			// Company A -> Company B -> Company A should still produce two tenures for A.
			eventA1 := fixtures.EventWith("a1", "First at A", "Company A", "")
			eventA1.Date = time.Date(2022, 3, 15, 0, 0, 0, 0, time.UTC)

			eventB := fixtures.EventWith("b1", "Work at B", "Company B", "")
			eventB.Date = time.Date(2022, 8, 15, 0, 0, 0, 0, time.UTC)

			eventA2 := fixtures.EventWith("a2", "Back at A", "Company A", "")
			eventA2.Date = time.Date(2023, 7, 15, 0, 0, 0, 0, time.UTC)

			events := []*career.Event{eventA1, eventB, eventA2}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())

			companyACount := 0
			for key := range result {
				if key == "Company A" || strings.HasPrefix(key, "Company A"+constants.TenureSeparator) {
					companyACount++
				}
			}
			Expect(companyACount).To(Equal(2), "Real company events should still split tenures")
		})
	})
})

var _ = Describe("determineSkillLevelFromFact", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should return expert for principal role fit", func() {
		fact := fixtures.Fact("f1", "e1")
		fact.RoleFit = career.RoleFitPrincipal
		Expect(dps.determineSkillLevelFromFact(fact)).To(Equal("expert"))
	})

	It("should return advanced for staff role fit", func() {
		fact := fixtures.Fact("f2", "e1")
		fact.RoleFit = career.RoleFitStaff
		Expect(dps.determineSkillLevelFromFact(fact)).To(Equal("advanced"))
	})

	It("should return advanced for senior IC role fit", func() {
		fact := fixtures.Fact("f3", "e1")
		fact.RoleFit = career.RoleFitSeniorIC
		Expect(dps.determineSkillLevelFromFact(fact)).To(Equal("advanced"))
	})

	It("should return advanced for EM role fit", func() {
		fact := fixtures.Fact("f4", "e1")
		fact.RoleFit = career.RoleFitEM
		Expect(dps.determineSkillLevelFromFact(fact)).To(Equal("advanced"))
	})

	It("should return intermediate for unknown role fit", func() {
		fact := fixtures.Fact("f5", "e1")
		fact.RoleFit = "unknown"
		Expect(dps.determineSkillLevelFromFact(fact)).To(Equal("intermediate"))
	})
})

var _ = Describe("extractSkillNameFromFact", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should return full text when 5 words or fewer", func() {
		Expect(dps.extractSkillNameFromFact("Go programming")).To(Equal("Go programming"))
	})

	It("should truncate to first 5 words when text is longer", func() {
		result := dps.extractSkillNameFromFact("one two three four five six seven")
		Expect(result).To(Equal("one two three four five"))
	})

	It("should return exact text with exactly 5 words", func() {
		Expect(dps.extractSkillNameFromFact("one two three four five")).To(Equal("one two three four five"))
	})

	It("should handle empty string", func() {
		Expect(dps.extractSkillNameFromFact("")).To(Equal(""))
	})
})

var _ = Describe("determineSkillCategory", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should use first category when categories exist", func() {
		skill := &Skill{Name: "test", Categories: []string{"Technical", "Leadership"}}
		Expect(dps.determineSkillCategory(skill)).To(Equal("technical"))
	})

	It("should detect technical from name", func() {
		skill := &Skill{Name: "software engineering"}
		Expect(dps.determineSkillCategory(skill)).To(Equal("technical"))
	})

	It("should detect leadership from name", func() {
		skill := &Skill{Name: "team management"}
		Expect(dps.determineSkillCategory(skill)).To(Equal("leadership"))
	})

	It("should detect product from name", func() {
		skill := &Skill{Name: "product strategy"}
		Expect(dps.determineSkillCategory(skill)).To(Equal("product"))
	})

	It("should return other for unrecognised name", func() {
		skill := &Skill{Name: "communication"}
		Expect(dps.determineSkillCategory(skill)).To(Equal("other"))
	})

	It("should detect engineer keyword as technical", func() {
		skill := &Skill{Name: "data engineer"}
		Expect(dps.determineSkillCategory(skill)).To(Equal("technical"))
	})

	It("should detect manage keyword as leadership", func() {
		skill := &Skill{Name: "project manager"}
		Expect(dps.determineSkillCategory(skill)).To(Equal("leadership"))
	})
})

var _ = Describe("determineSkillLevel", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should return expert for technical tag", func() {
		Expect(dps.determineSkillLevel("technical")).To(Equal("expert"))
	})

	It("should return advanced for leadership tag", func() {
		Expect(dps.determineSkillLevel("leadership")).To(Equal("advanced"))
	})

	It("should return advanced for product tag", func() {
		Expect(dps.determineSkillLevel("product")).To(Equal("advanced"))
	})

	It("should return intermediate for unknown tag", func() {
		Expect(dps.determineSkillLevel("other")).To(Equal("intermediate"))
	})

	It("should be case-insensitive", func() {
		Expect(dps.determineSkillLevel("TECHNICAL")).To(Equal("expert"))
		Expect(dps.determineSkillLevel("Leadership")).To(Equal("advanced"))
	})
})

var _ = Describe("extractPosition", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should return Professional when no position keywords found", func() {
		events := []*career.Event{fixtures.EventWith("e1", "did some work", "Acme", "")}
		Expect(dps.extractPosition(events)).To(Equal("Professional"))
	})

	It("should detect engineer from event text", func() {
		events := []*career.Event{fixtures.EventWith("e1", "Senior engineer built API", "Acme", "")}
		Expect(dps.extractPosition(events)).To(Equal("Engineer"))
	})

	It("should detect manager from event text", func() {
		events := []*career.Event{fixtures.EventWith("e1", "Project manager led delivery", "Acme", "")}
		Expect(dps.extractPosition(events)).To(Equal("Manager"))
	})

	It("should detect architect from event text", func() {
		events := []*career.Event{fixtures.EventWith("e1", "Solutions architect designed system", "Acme", "")}
		Expect(dps.extractPosition(events)).To(Equal("Architect"))
	})

	It("should detect director from event text", func() {
		events := []*career.Event{fixtures.EventWith("e1", "Director of product operations", "Acme", "")}
		Expect(dps.extractPosition(events)).To(Equal("Director"))
	})

	It("should detect lead from event text", func() {
		events := []*career.Event{fixtures.EventWith("e1", "Tech lead guided team", "Acme", "")}
		Expect(dps.extractPosition(events)).To(Equal("Lead"))
	})

	It("should return most common position across multiple events", func() {
		events := []*career.Event{
			fixtures.EventWith("e1", "Senior engineer built API", "Acme", ""),
			fixtures.EventWith("e2", "Lead engineer designed system", "Acme", ""),
			fixtures.EventWith("e3", "Staff engineer optimised performance", "Acme", ""),
		}
		Expect(dps.extractPosition(events)).To(Equal("Engineer"))
	})
})

var _ = Describe("calculateDateRange", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should return now for empty events", func() {
		start, end := dps.calculateDateRange([]*career.Event{})
		Expect(start).To(BeTemporally("~", time.Now(), time.Second))
		Expect(end).To(BeTemporally("~", time.Now(), time.Second))
	})

	It("should return same date for single event", func() {
		date := time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC)
		ev := fixtures.EventWith("e1", "Built API", "Acme", "")
		ev.Date = date
		start, end := dps.calculateDateRange([]*career.Event{ev})
		Expect(start).To(Equal(date))
		Expect(end).To(Equal(date))
	})

	It("should find earliest and latest dates", func() {
		early := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
		mid := time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC)
		late := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
		ev1 := fixtures.EventWith("e1", "First", "Acme", "")
		ev1.Date = mid
		ev2 := fixtures.EventWith("e2", "Second", "Acme", "")
		ev2.Date = early
		ev3 := fixtures.EventWith("e3", "Third", "Acme", "")
		ev3.Date = late
		start, end := dps.calculateDateRange([]*career.Event{ev1, ev2, ev3})
		Expect(start).To(Equal(early))
		Expect(end).To(Equal(late))
	})
})

var _ = Describe("extractContext", func() {
	It("should return empty for no match", func() {
		Expect(extractContext("hello world", "missing")).To(BeEmpty())
	})

	It("should extract context around match", func() {
		text := "The quick brown fox jumps over the lazy dog"
		result := extractContext(text, "fox")
		Expect(result).To(ContainSubstring("fox"))
	})

	It("should handle match at start of text", func() {
		result := extractContext("fox jumps over the lazy dog", "fox")
		Expect(result).To(HavePrefix("fox"))
	})

	It("should handle match at end of text", func() {
		result := extractContext("The quick brown fox", "fox")
		Expect(result).To(ContainSubstring("fox"))
	})
})

var _ = Describe("mergeSkills", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should return empty skill for empty input", func() {
		result := dps.mergeSkills([]*Skill{})
		Expect(result.Name).To(BeEmpty())
	})

	It("should merge projects and endorsements", func() {
		skills := []*Skill{
			{Name: "Go", Level: "expert", Projects: 3, Endorsements: 1, Categories: []string{"technical"}},
			{Name: "Go", Level: "expert", Projects: 2, Endorsements: 2, Categories: []string{"backend"}},
		}
		result := dps.mergeSkills(skills)
		Expect(result.Name).To(Equal("Go"))
		Expect(result.Projects).To(Equal(5))
		Expect(result.Endorsements).To(Equal(3))
		Expect(result.Categories).To(ContainElements("technical", "backend"))
	})

	It("should deduplicate categories", func() {
		skills := []*Skill{
			{Name: "Go", Level: "expert", Projects: 1, Categories: []string{"technical"}},
			{Name: "Go", Level: "expert", Projects: 1, Categories: []string{"technical"}},
		}
		result := dps.mergeSkills(skills)
		Expect(result.Categories).To(HaveLen(1))
	})
})

var _ = Describe("ExtractSkills", func() {
	var (
		dps DataProcessingService
		ctx context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()
		dps = NewDataProcessingService(logger.DefaultLogger())
	})

	It("should return nil for cancelled context", func() {
		cancelCtx, cancel := context.WithCancel(ctx)
		cancel()
		result, err := dps.ExtractSkills(cancelCtx, nil, nil)
		Expect(err).To(HaveOccurred())
		Expect(result).To(BeNil())
	})

	It("should extract skills from event tags", func() {
		ev := fixtures.EventWith("e1", "Built microservices", "Acme", "")
		ev.Tags = []string{"technical"}
		result, err := dps.ExtractSkills(ctx, []*career.Event{ev}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).NotTo(BeEmpty())
	})

	It("should extract skills from facts", func() {
		fact := fixtures.FactWithCategories("f1", "Implemented Go microservices", "e1", []string{"technical"}, []string{"peer"})
		result, err := dps.ExtractSkills(ctx, nil, []*career.Fact{fact})
		Expect(err).NotTo(HaveOccurred())
		Expect(result).NotTo(BeEmpty())
	})

	It("should categorise skills by category", func() {
		ev := fixtures.EventWith("e1", "Led team", "Acme", "")
		ev.Tags = []string{"leadership"}
		result, err := dps.ExtractSkills(ctx, []*career.Event{ev}, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(result).To(HaveKey("leadership"))
	})
})

var _ = Describe("extractPrimaryCategory (data processing)", func() {
	var dps *DefaultDataProcessingService

	BeforeEach(func() {
		dps = NewDataProcessingService(logger.DefaultLogger()).(*DefaultDataProcessingService)
	})

	It("should return empty for no categories", func() {
		Expect(dps.extractPrimaryCategory(nil)).To(Equal(constants.CompetencyCategory("")))
	})

	It("should return valid category", func() {
		result := dps.extractPrimaryCategory([]string{"technical"})
		Expect(result).To(Equal(constants.CompetencyTechnical))
	})

	It("should return empty for invalid category", func() {
		result := dps.extractPrimaryCategory([]string{"nonexistent"})
		Expect(result).To(Equal(constants.CompetencyCategory("")))
	})
})
