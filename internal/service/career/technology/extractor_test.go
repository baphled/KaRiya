package technology_test

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/technology"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Extractor", func() {
	var (
		extractor  *technology.Extractor
		skillRepo  careerRepo.SkillRepository
		eventRepo  careerRepo.EventRepository
		ctx        context.Context
		testSkills []*career.Skill
		testEvents []*career.CareerEvent
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Create in-memory repositories
		skillRepo = careermemory.NewSkillRepository()
		eventRepo = careermemory.NewEventRepository()

		// Create test skills
		now := time.Now()
		testSkills = []*career.Skill{
			{
				ID:        "skill-1",
				Name:      "Ruby",
				Category:  "backend",
				Level:     "expert",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "skill-2",
				Name:      "PostgreSQL",
				Category:  "database",
				Level:     "advanced",
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "skill-3",
				Name:      "React",
				Category:  "frontend",
				Level:     "intermediate",
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		// Create test events (using correct CareerEvent fields)
		testEvents = []*career.CareerEvent{
			{
				ID:        "event-1",
				Text:      "Backend Developer at Company A - Worked with Ruby and PostgreSQL",
				Company:   "Company A",
				Date:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Skills:    []string{"skill-1", "skill-2"}, // Skill IDs
				Tags:      []string{"project"},
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "event-2",
				Text:      "Senior Backend Developer at Company A - Continued Ruby work",
				Company:   "Company A",
				Date:      time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC),
				Skills:    []string{"skill-1"},
				Tags:      []string{"project"},
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "event-3",
				Text:      "Fullstack Developer at Company B - Worked with React and Ruby",
				Company:   "Company B",
				Date:      time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				Skills:    []string{"skill-1", "skill-3"},
				Tags:      []string{"project"},
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        "event-4",
				Text:      "Project without skills - Legacy event",
				Company:   "Company C",
				Date:      time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
				Skills:    []string{}, // No skills
				Tags:      []string{"project"},
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		// Populate repositories
		for _, skill := range testSkills {
			err := skillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
		}

		for _, event := range testEvents {
			err := eventRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())
		}

		// Create extractor
		extractor = technology.NewExtractor(skillRepo, eventRepo)
	})

	Describe("ExtractFromUser", func() {
		Context("when user has skills with event associations", func() {
			It("should extract technologies with event counts", func() {
				techs, err := extractor.ExtractFromUser(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(techs).To(HaveLen(3))

				// Find Ruby skill (should have 3 events)
				var rubyTech *technology.ExtractedTechnology
				for _, tech := range techs {
					if tech.Name == "Ruby" {
						rubyTech = tech
						break
					}
				}

				Expect(rubyTech).NotTo(BeNil())
				Expect(rubyTech.ID).To(Equal("skill-1"))
				Expect(rubyTech.Category).To(Equal("backend"))
				Expect(rubyTech.EventCount).To(Equal(3))
				Expect(rubyTech.EventIDs).To(ConsistOf("event-1", "event-2", "event-3"))
			})

			It("should include technologies with different event counts", func() {
				techs, err := extractor.ExtractFromUser(ctx)

				Expect(err).NotTo(HaveOccurred())

				// PostgreSQL: 1 event
				var postgresqlTech *technology.ExtractedTechnology
				for _, tech := range techs {
					if tech.Name == "PostgreSQL" {
						postgresqlTech = tech
						break
					}
				}

				Expect(postgresqlTech).NotTo(BeNil())
				Expect(postgresqlTech.EventCount).To(Equal(1))
				Expect(postgresqlTech.EventIDs).To(ConsistOf("event-1"))

				// React: 1 event
				var reactTech *technology.ExtractedTechnology
				for _, tech := range techs {
					if tech.Name == "React" {
						reactTech = tech
						break
					}
				}

				Expect(reactTech).NotTo(BeNil())
				Expect(reactTech.EventCount).To(Equal(1))
				Expect(reactTech.EventIDs).To(ConsistOf("event-3"))
			})

			It("should sort technologies by event count (descending)", func() {
				techs, err := extractor.ExtractFromUser(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(techs).To(HaveLen(3))

				// Ruby should be first (3 events)
				Expect(techs[0].Name).To(Equal("Ruby"))
				Expect(techs[0].EventCount).To(Equal(3))

				// PostgreSQL and React should follow (1 event each)
				// Order between them may vary, but both should have 1 event
				Expect(techs[1].EventCount).To(Equal(1))
				Expect(techs[2].EventCount).To(Equal(1))
			})
		})

		Context("when user has no skills", func() {
			It("should return empty list", func() {
				// Create empty repositories
				emptySkillRepo := careermemory.NewSkillRepository()
				emptyEventRepo := careermemory.NewEventRepository()
				emptyExtractor := technology.NewExtractor(emptySkillRepo, emptyEventRepo)

				techs, err := emptyExtractor.ExtractFromUser(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(techs).To(BeEmpty())
			})
		})
	})

	Describe("FilterByThreshold", func() {
		var extractedTechs []*technology.ExtractedTechnology

		BeforeEach(func() {
			var err error
			extractedTechs, err = extractor.ExtractFromUser(ctx)
			Expect(err).NotTo(HaveOccurred())
		})

		Context("when threshold is 3", func() {
			It("should only include technologies with 3+ events", func() {
				filtered := extractor.FilterByThreshold(extractedTechs, 3)

				Expect(filtered).To(HaveLen(1))
				Expect(filtered[0].Name).To(Equal("Ruby"))
				Expect(filtered[0].EventCount).To(Equal(3))
			})
		})

		Context("when threshold is 1", func() {
			It("should include all technologies with 1+ events", func() {
				filtered := extractor.FilterByThreshold(extractedTechs, 1)

				Expect(filtered).To(HaveLen(3))
			})
		})

		Context("when threshold is 0", func() {
			It("should include all technologies", func() {
				filtered := extractor.FilterByThreshold(extractedTechs, 0)

				Expect(filtered).To(HaveLen(3))
			})
		})

		Context("when input is empty", func() {
			It("should return empty list", func() {
				filtered := extractor.FilterByThreshold([]*technology.ExtractedTechnology{}, 3)

				Expect(filtered).To(BeEmpty())
			})
		})
	})
})
