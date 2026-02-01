package skillsmanagement_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/skillsmanagement"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

// MockSkillRepository implements SkillRepository for testing.
type MockSkillRepository struct {
	skills        []*career.Skill
	eventCounts   map[string]int
	events        []*career.Event
	createErr     error
	updateErr     error
	deleteErr     error
	listErr       error
	eventsErr     error
	eventsBySkill map[string][]*career.Event
}

func NewMockSkillRepository() *MockSkillRepository {
	return &MockSkillRepository{
		skills:        make([]*career.Skill, 0),
		eventCounts:   make(map[string]int),
		eventsBySkill: make(map[string][]*career.Event),
	}
}

func (m *MockSkillRepository) Create(_ context.Context, skill *career.Skill) error {
	if m.createErr != nil {
		return m.createErr
	}
	if skill.ID == "" {
		skill.ID = "skill-" + time.Now().Format("20060102150405")
	}
	m.skills = append(m.skills, skill)
	return nil
}

func (m *MockSkillRepository) GetByID(_ context.Context, id string) (*career.Skill, error) {
	for _, s := range m.skills {
		if s.ID == id {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockSkillRepository) Update(_ context.Context, skill *career.Skill) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i, s := range m.skills {
		if s.ID == skill.ID {
			m.skills[i] = skill
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockSkillRepository) Delete(_ context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	for i, s := range m.skills {
		if s.ID == id {
			m.skills = append(m.skills[:i], m.skills[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockSkillRepository) List(_ context.Context, _ *careerrepo.SkillListFilters) ([]*career.Skill, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.skills, nil
}

func (m *MockSkillRepository) GetByName(_ context.Context, name string) (*career.Skill, error) {
	for _, s := range m.skills {
		if s.Name == name {
			return s, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockSkillRepository) GetByCategory(_ context.Context, category string) ([]*career.Skill, error) {
	var result []*career.Skill
	for _, s := range m.skills {
		if s.Category == category {
			result = append(result, s)
		}
	}
	return result, nil
}

func (m *MockSkillRepository) GetSkillsForEvent(_ context.Context, _ string) ([]*career.Skill, error) {
	return m.skills, nil
}

func (m *MockSkillRepository) GetEventCountsForSkills(_ context.Context) (map[string]int, error) {
	return m.eventCounts, nil
}

func (m *MockSkillRepository) GetLastUsedForSkills(_ context.Context) (map[string]time.Time, error) {
	return make(map[string]time.Time), nil
}

func (m *MockSkillRepository) GetEventsUsingSkill(_ context.Context, skillID string) ([]*career.Event, error) {
	if m.eventsErr != nil {
		return nil, m.eventsErr
	}
	if events, ok := m.eventsBySkill[skillID]; ok {
		return events, nil
	}
	return m.events, nil
}

var _ = Describe("Context", func() {
	Describe("Errors", func() {
		It("should define ErrNoSkillSelected", func() {
			Expect(skillsmanagement.ErrNoSkillSelected).To(MatchError("no skill selected"))
		})

		It("should define ErrInvalidSkillData", func() {
			Expect(skillsmanagement.ErrInvalidSkillData).To(MatchError("invalid skill data"))
		})
	})

	Describe("Filters", func() {
		It("should have Category field for filtering by category", func() {
			filters := &skillsmanagement.Filters{
				Category: "Programming",
			}
			Expect(filters.Category).To(Equal("Programming"))
		})

		It("should have Level field for filtering by proficiency level", func() {
			filters := &skillsmanagement.Filters{
				Level: "advanced",
			}
			Expect(filters.Level).To(Equal("advanced"))
		})

		It("should have MinEvents field for filtering by minimum event count", func() {
			filters := &skillsmanagement.Filters{
				MinEvents: 5,
			}
			Expect(filters.MinEvents).To(Equal(5))
		})

		It("should have SearchText field for text search", func() {
			filters := &skillsmanagement.Filters{
				SearchText: "golang",
			}
			Expect(filters.SearchText).To(Equal("golang"))
		})

		It("should have SortBy field for sort column", func() {
			filters := &skillsmanagement.Filters{
				SortBy: "name",
			}
			Expect(filters.SortBy).To(Equal("name"))
		})

		It("should have SortOrder field for sort direction", func() {
			filters := &skillsmanagement.Filters{
				SortOrder: "asc",
			}
			Expect(filters.SortOrder).To(Equal("asc"))
		})

		Describe("HasActiveFilters", func() {
			It("should return false when no filters are set", func() {
				filters := &skillsmanagement.Filters{}
				Expect(filters.HasActiveFilters()).To(BeFalse())
			})

			It("should return true when Category is set", func() {
				filters := &skillsmanagement.Filters{Category: "Programming"}
				Expect(filters.HasActiveFilters()).To(BeTrue())
			})

			It("should return true when Level is set", func() {
				filters := &skillsmanagement.Filters{Level: "advanced"}
				Expect(filters.HasActiveFilters()).To(BeTrue())
			})

			It("should return true when MinEvents is set", func() {
				filters := &skillsmanagement.Filters{MinEvents: 1}
				Expect(filters.HasActiveFilters()).To(BeTrue())
			})

			It("should return true when SearchText is set", func() {
				filters := &skillsmanagement.Filters{SearchText: "go"}
				Expect(filters.HasActiveFilters()).To(BeTrue())
			})

			It("should return true when SortBy is set", func() {
				filters := &skillsmanagement.Filters{SortBy: "name"}
				Expect(filters.HasActiveFilters()).To(BeTrue())
			})
		})

		Describe("Clear", func() {
			It("should clear all filter fields in FIFO order - search first", func() {
				filters := &skillsmanagement.Filters{
					SearchText: "go",
					Category:   "Programming",
					SortBy:     "name",
				}
				filters.Clear()
				Expect(filters.SearchText).To(BeEmpty())
				Expect(filters.Category).To(Equal("Programming"))
			})

			It("should clear category/level filters after search", func() {
				filters := &skillsmanagement.Filters{
					Category: "Programming",
					Level:    "advanced",
					SortBy:   "name",
				}
				filters.Clear()
				Expect(filters.Category).To(BeEmpty())
				Expect(filters.Level).To(BeEmpty())
				Expect(filters.SortBy).To(Equal("name"))
			})

			It("should clear sort last", func() {
				filters := &skillsmanagement.Filters{
					SortBy:    "name",
					SortOrder: "asc",
				}
				filters.Clear()
				Expect(filters.SortBy).To(BeEmpty())
				Expect(filters.SortOrder).To(BeEmpty())
			})
		})
	})

	Describe("IntentContext", func() {
		var (
			ctx      context.Context
			mockRepo *MockSkillRepository
		)

		BeforeEach(func() {
			ctx = context.Background()
			mockRepo = NewMockSkillRepository()
		})

		Describe("Construction", func() {
			It("should create a context with repository", func() {
				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx).NotTo(BeNil())
			})

			It("should store the context", func() {
				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx.Ctx).To(Equal(ctx))
			})

			It("should store the repository", func() {
				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx.SkillRepository).To(Equal(mockRepo))
			})

			It("should initialize empty filters", func() {
				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx.Filters).NotTo(BeNil())
			})
		})

		Describe("Validate", func() {
			Context("when repository is nil", func() {
				It("should return error", func() {
					intentCtx := &skillsmanagement.IntentContext{
						Ctx:             ctx,
						SkillRepository: nil,
					}
					err := intentCtx.Validate()
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when context is nil", func() {
				It("should return error", func() {
					intentCtx := &skillsmanagement.IntentContext{
						Ctx:             nil,
						SkillRepository: mockRepo,
					}
					err := intentCtx.Validate()
					Expect(err).To(HaveOccurred())
				})
			})

			Context("when all fields are valid", func() {
				It("should return no error", func() {
					intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when Filters is nil", func() {
				It("should initialize to empty filters", func() {
					intentCtx := &skillsmanagement.IntentContext{
						Ctx:             ctx,
						SkillRepository: mockRepo,
						Filters:         nil,
					}
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Filters).NotTo(BeNil())
				})
			})
		})

		Describe("LoadSkills", func() {
			var testSkill *career.Skill

			BeforeEach(func() {
				testSkill = fixtures.SkillWith("skill-1", "Go Programming", "Programming", "advanced")
			})

			Context("with valid repository", func() {
				It("should load skills successfully", func() {
					mockRepo.skills = []*career.Skill{testSkill}
					intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
					skills, err := intentCtx.LoadSkills()
					Expect(err).NotTo(HaveOccurred())
					Expect(skills).To(HaveLen(1))
				})

				It("should return empty slice when no skills", func() {
					intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
					skills, err := intentCtx.LoadSkills()
					Expect(err).NotTo(HaveOccurred())
					Expect(skills).To(BeEmpty())
				})
			})

			Context("when repository returns error", func() {
				It("should propagate the error", func() {
					mockRepo.listErr = errors.New("database error")
					intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
					_, err := intentCtx.LoadSkills()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("database error"))
				})
			})
		})

		Describe("GetEventCounts", func() {
			It("should return event counts for skills", func() {
				mockRepo.eventCounts = map[string]int{
					"skill-1": 5,
					"skill-2": 3,
				}
				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				counts, err := intentCtx.GetEventCounts()
				Expect(err).NotTo(HaveOccurred())
				Expect(counts["skill-1"]).To(Equal(5))
			})
		})

		Describe("GetEventsForSkill", func() {
			It("should return events for a specific skill", func() {
				events := fixtures.Events(1)
				mockRepo.eventsBySkill["skill-1"] = events
				intentCtx := skillsmanagement.NewIntentContext(ctx, mockRepo)
				result, err := intentCtx.GetEventsForSkill("skill-1")
				Expect(err).NotTo(HaveOccurred())
				Expect(result).To(HaveLen(1))
			})
		})
	})
})
