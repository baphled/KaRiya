package cv

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
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
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Led team standup",
					Date:      time.Now().AddDate(0, -1, 0),
					Company:   "Acme Corp",
					Tags:      []string{"leadership"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        "2",
					Text:      "Implemented feature",
					Date:      time.Now().AddDate(0, -2, 0),
					Company:   "Acme Corp",
					Tags:      []string{"technical"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(1))
			Expect(result["Acme Corp"]).NotTo(BeNil())
			Expect(result["Acme Corp"].Company).To(Equal("Acme Corp"))
			Expect(result["Acme Corp"].EventIDs).To(HaveLen(2))
		})

		It("should group multiple companies", func() {
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Worked at Acme",
					Date:      time.Now().AddDate(0, -1, 0),
					Company:   "Acme Corp",
					Tags:      []string{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        "2",
					Text:      "Worked at TechCo",
					Date:      time.Now().AddDate(0, -6, 0),
					Company:   "TechCo",
					Tags:      []string{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HaveLen(2))
			Expect(result["Acme Corp"]).NotTo(BeNil())
			Expect(result["TechCo"]).NotTo(BeNil())
		})

		It("should handle events with no company", func() {
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Personal project",
					Date:      time.Now(),
					Company:   "",
					Tags:      []string{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result["Other"]).NotTo(BeNil())
		})

		It("should set correct date ranges", func() {
			now := time.Now()
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Recent event",
					Date:      now,
					Company:   "Acme",
					Tags:      []string{},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Old event",
					Date:      now.AddDate(-1, 0, 0),
					Company:   "Acme",
					Tags:      []string{},
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

			result, err := dps.GroupEventsByCompany(svc, events)
			Expect(err).NotTo(HaveOccurred())

			group := result["Acme"]
			Expect(group.StartDate).To(Equal(now.AddDate(-1, 0, 0)))
			Expect(group.EndDate).To(Equal(now))
		})
	})

	Describe("ExtractAchievements", func() {
		It("should extract achievement from event", func() {
			event := &career.CareerEvent{
				ID:        "1",
				Text:      "Led team of 12 engineers",
				Date:      time.Now(),
				Company:   "Acme",
				Tags:      []string{"leadership"},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(achievements).To(HaveLen(1))
			Expect(achievements[0].EventID).To(Equal("1"))
			Expect(achievements[0].ActionVerb).NotTo(BeEmpty())
		})

		It("should extract metrics from achievement", func() {
			event := &career.CareerEvent{
				ID:        "1",
				Text:      "Increased performance by 25% with team of 12 people",
				Date:      time.Now(),
				Company:   "Acme",
				Tags:      []string{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(achievements[0].Metrics).NotTo(BeEmpty())
		})

		It("should extract facts as achievements", func() {
			event := &career.CareerEvent{
				ID:        "1",
				Text:      "Event text",
				Date:      time.Now(),
				Company:   "Acme",
				Tags:      []string{},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
			fact := &career.Fact{
				ID:            "f1",
				Text:          "Strong leadership capability",
				SourceEventID: "1",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}

			achievements, err := dps.ExtractAchievements(svc, event, []*career.Fact{fact})
			Expect(err).NotTo(HaveOccurred())
			Expect(achievements).To(HaveLen(2)) // Event + fact
			Expect(achievements[1].FactIDs).To(ContainElement("f1"))
		})
	})

	Describe("ExtractSkills", func() {
		It("should extract empty skills", func() {
			result, err := dps.ExtractSkills(svc, []*career.CareerEvent{}, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeEmpty())
		})

		It("should extract skills from event tags", func() {
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event",
					Date:      time.Now(),
					Company:   "Acme",
					Tags:      []string{"technical", "leadership"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			result, err := dps.ExtractSkills(svc, events, []*career.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
		})

		It("should aggregate skill categories", func() {
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event",
					Date:      time.Now(),
					Company:   "Acme",
					Tags:      []string{"technical"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        "2",
					Text:      "Event",
					Date:      time.Now(),
					Company:   "Acme",
					Tags:      []string{"technical"},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

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
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event",
					Date:      time.Now(),
					Company:   "Acme",
					Project:   "",
					Tags:      []string{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			result, err := dps.ExtractProjectsFromEvents(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeNil())
		})

		It("should extract single project", func() {
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event",
					Date:      time.Now(),
					Company:   "Acme",
					Project:   "ProjectX",
					Tags:      []string{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

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
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event",
					Date:      time.Now(),
					Company:   "Acme",
					Project:   "ProjectX",
					Tags:      []string{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
				{
					ID:        "2",
					Text:      "Event",
					Date:      time.Now(),
					Company:   "Acme",
					Project:   "ProjectY",
					Tags:      []string{},
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			result, err := dps.ExtractProjectsFromEvents(svc, events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(BeEmpty())
		})

		It("should set correct project date ranges", func() {
			now := time.Now()
			events := []*career.CareerEvent{
				{
					ID:        "1",
					Text:      "Event",
					Date:      now,
					Company:   "Acme",
					Project:   "ProjectX",
					Tags:      []string{},
					CreatedAt: now,
					UpdatedAt: now,
				},
				{
					ID:        "2",
					Text:      "Event",
					Date:      now.AddDate(0, -3, 0),
					Company:   "Acme",
					Project:   "ProjectX",
					Tags:      []string{},
					CreatedAt: now,
					UpdatedAt: now,
				},
			}

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
