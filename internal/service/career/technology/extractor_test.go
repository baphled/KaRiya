package technology_test

import (
	"context"
	"errors"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	careerRepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/service/career/technology"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

type failingSkillRepository struct {
	listErr error
}

var errNotImplemented = errors.New("not implemented")

func (r *failingSkillRepository) Create(_ context.Context, _ *career.Skill) error { return nil }
func (r *failingSkillRepository) GetByID(_ context.Context, _ string) (*career.Skill, error) {
	return nil, errNotImplemented
}
func (r *failingSkillRepository) GetByName(_ context.Context, _ string) (*career.Skill, error) {
	return nil, errNotImplemented
}
func (r *failingSkillRepository) List(_ context.Context, _ *careerRepo.SkillListFilters) ([]*career.Skill, error) {
	return nil, r.listErr
}
func (r *failingSkillRepository) Update(_ context.Context, _ *career.Skill) error { return nil }
func (r *failingSkillRepository) Delete(_ context.Context, _ string) error        { return nil }
func (r *failingSkillRepository) GetByCategory(_ context.Context, _ string) ([]*career.Skill, error) {
	return nil, errNotImplemented
}
func (r *failingSkillRepository) GetSkillsForEvent(_ context.Context, _ string) ([]*career.Skill, error) {
	return nil, errNotImplemented
}
func (r *failingSkillRepository) GetEventCountsForSkills(_ context.Context) (map[string]int, error) {
	return nil, errNotImplemented
}
func (r *failingSkillRepository) GetLastUsedForSkills(_ context.Context) (map[string]time.Time, error) {
	return nil, errNotImplemented
}
func (r *failingSkillRepository) GetEventsUsingSkill(_ context.Context, _ string) ([]*career.Event, error) {
	return nil, errNotImplemented
}

type failingEventRepository struct {
	listErr error
}

func (r *failingEventRepository) Create(_ context.Context, _ *career.Event) error { return nil }
func (r *failingEventRepository) GetByID(_ context.Context, _ string) (*career.Event, error) {
	return nil, errNotImplemented
}
func (r *failingEventRepository) Update(_ context.Context, _ *career.Event) error { return nil }
func (r *failingEventRepository) Delete(_ context.Context, _ string) error        { return nil }
func (r *failingEventRepository) List(_ context.Context, _ careerRepo.EventListFilters) ([]*career.Event, error) {
	return nil, r.listErr
}
func (r *failingEventRepository) Count(_ context.Context, _ careerRepo.EventListFilters) (int, error) {
	return 0, nil
}
func (r *failingEventRepository) LinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}
func (r *failingEventRepository) UnlinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

var _ = Describe("Extractor", func() {
	var (
		extractor *technology.Extractor
		skillRepo careerRepo.SkillRepository
		eventRepo careerRepo.EventRepository
		ctx       context.Context
	)

	BeforeEach(func() {
		ctx = context.Background()

		// Create in-memory repositories
		skillRepo = careermemory.NewSkillRepository()
		eventRepo = careermemory.NewEventRepository()

		skill1 := fixtures.SkillWith("skill-1", "Ruby", "backend", "expert")
		skill2 := fixtures.SkillWith("skill-2", "PostgreSQL", "database", "advanced")
		skill3 := fixtures.SkillWith("skill-3", "React", "frontend", "intermediate")

		event1 := fixtures.EventWith("event-1", "Backend Developer at Company A - Worked with Ruby and PostgreSQL", "Company A", "")
		event1.Date = time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
		event1.Skills = []string{"skill-1", "skill-2"}
		event1.Tags = []string{"project"}

		event2 := fixtures.EventWith("event-2", "Senior Backend Developer at Company A - Continued Ruby work", "Company A", "")
		event2.Date = time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
		event2.Skills = []string{"skill-1"}
		event2.Tags = []string{"project"}

		event3 := fixtures.EventWith("event-3", "Fullstack Developer at Company B - Worked with React and Ruby", "Company B", "")
		event3.Date = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		event3.Skills = []string{"skill-1", "skill-3"}
		event3.Tags = []string{"project"}

		event4 := fixtures.EventWith("event-4", "Project without skills - Legacy event", "Company C", "")
		event4.Date = time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		event4.Skills = []string{}
		event4.Tags = []string{"project"}

		for _, skill := range []*career.Skill{skill1, skill2, skill3} {
			err := skillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())
		}

		for _, event := range []*career.Event{event1, event2, event3, event4} {
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
				emptySkillRepo := careermemory.NewSkillRepository()
				emptyEventRepo := careermemory.NewEventRepository()
				emptyExtractor := technology.NewExtractor(emptySkillRepo, emptyEventRepo)

				techs, err := emptyExtractor.ExtractFromUser(ctx)

				Expect(err).NotTo(HaveOccurred())
				Expect(techs).To(BeEmpty())
			})
		})

		Context("when skill repository fails", func() {
			It("should propagate the error", func() {
				skillListErr := errors.New("skill repository unavailable")
				failingSkillRepo := &failingSkillRepository{listErr: skillListErr}
				failingExtractor := technology.NewExtractor(failingSkillRepo, eventRepo)

				_, err := failingExtractor.ExtractFromUser(ctx)

				Expect(err).To(MatchError(skillListErr))
			})
		})

		Context("when event repository fails", func() {
			It("should propagate the error", func() {
				eventListErr := errors.New("event repository unavailable")
				failingEventRepo := &failingEventRepository{listErr: eventListErr}
				failingExtractor := technology.NewExtractor(skillRepo, failingEventRepo)

				_, err := failingExtractor.ExtractFromUser(ctx)

				Expect(err).To(MatchError(eventListErr))
			})
		})

		Context("when context is cancelled", func() {
			It("should return the context error", func() {
				cancelledCtx, cancel := context.WithCancel(ctx)
				cancel()

				cancellingSkillRepo := &failingSkillRepository{listErr: cancelledCtx.Err()}
				cancellingExtractor := technology.NewExtractor(cancellingSkillRepo, eventRepo)

				_, err := cancellingExtractor.ExtractFromUser(cancelledCtx)

				Expect(err).To(MatchError(context.Canceled))
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
