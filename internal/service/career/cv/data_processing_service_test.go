package cv

import (
	"context"
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
			result, err := dps.GroupEventsByCompany(svc, []*career.CareerEvent{})
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

			events := []*career.CareerEvent{event1, event2}

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

			events := []*career.CareerEvent{event1, event2}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(2))
			Expect(result["Acme Corp"]).NotTo(BeNil())
			Expect(result["TechCo"]).NotTo(BeNil())
		})

		It("should handle events with no company", func() {
			event := fixtures.EventWith("1", "Personal project", "", "") // No company

			events := []*career.CareerEvent{event}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result["Other"]).NotTo(BeNil())
		})

		It("should set correct date ranges", func() {
			now := time.Now()

			event1 := fixtures.EventWith("1", "Recent event", "Acme", "")
			event1.Date = now

			event2 := fixtures.EventWith("2", "Old event", "Acme", "")
			event2.Date = now.AddDate(-1, 0, 0)

			events := []*career.CareerEvent{event1, event2}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())

			group := result["Acme"]
			Expect(group.StartDate).To(Equal(now.AddDate(-1, 0, 0)))
			Expect(group.EndDate).To(Equal(now))
		})
	})

	Describe("ExtractAchievements", func() {
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
			result, err := dps.ExtractSkills(svc, []*career.CareerEvent{}, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should extract skills from event tags", func() {
			event := fixtures.EventWith("1", "Event", "Acme", "")
			event.Tags = []string{"technical", "leadership"}

			events := []*career.CareerEvent{event}

			result, err := dps.ExtractSkills(svc, events, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
		})

		It("should aggregate skill categories", func() {
			event1 := fixtures.EventWith("1", "Event", "Acme", "")
			event1.Tags = []string{"technical"}

			event2 := fixtures.EventWith("2", "Event", "Acme", "")
			event2.Tags = []string{"technical"}

			events := []*career.CareerEvent{event1, event2}

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

			events := []*career.CareerEvent{event}

			result, err := dps.ExtractProjectsFromEvents(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})

		It("should extract single project", func() {
			event := fixtures.EventWith("1", "Event", "Acme", "ProjectX")

			events := []*career.CareerEvent{event}

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

			events := []*career.CareerEvent{event1, event2}

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

			events := []*career.CareerEvent{event1, event2}

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

			_, err := dps.GroupEventsByCompany(ctx, []*career.CareerEvent{})
			Expect(err).To(HaveOccurred())
		})
	})
})
