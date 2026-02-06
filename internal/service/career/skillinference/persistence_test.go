package skillinference_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Skill Persistence", func() {
	var (
		service     skillinference.SkillInferenceService
		skillRepo   *mockSkillRepository
		eventRepo   *mockEventRepository
		ctx         context.Context
		testTime    time.Time
		suggestions []skillinference.SkillSuggestion
	)

	BeforeEach(func() {
		skillRepo = &mockSkillRepository{
			skills:        make(map[string]*career.Skill),
			skillsByID:    make(map[string]*career.Skill),
			notFoundError: career_repo.ErrSkillNotFound,
		}
		eventRepo = &mockEventRepository{
			events: make(map[string]*career.Event),
		}

		service = skillinference.NewSkillInferenceService(skillRepo, eventRepo)
		ctx = context.Background() //nolint:fatcontext // test setup
		testTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	})

	Describe("CreateSkillsFromSuggestions", func() {
		Context("when creating new skills", func() {
			BeforeEach(func() {
				suggestions = []skillinference.SkillSuggestion{
					{
						Name:       "Go",
						Category:   "Backend",
						Confidence: 0.95,
						EventIDs:   []string{"event-1", "event-2"},
					},
					{
						Name:       "PostgreSQL",
						Category:   "Database",
						Confidence: 0.75,
						EventIDs:   []string{"event-1"},
					},
				}

				// Add events to mock repo
				event1 := fixtures.Event("event-1")
				event1.Date = testTime
				eventRepo.events["event-1"] = event1
				event2 := fixtures.Event("event-2")
				event2.Date = testTime.Add(24 * time.Hour)
				eventRepo.events["event-2"] = event2
			})

			It("should create skills from suggestions", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(2))

				// Verify Go skill
				goSkill := findSkillByName(skills, "Go")
				Expect(goSkill).NotTo(BeNil())
				Expect(goSkill.Name).To(Equal("Go"))
				Expect(goSkill.Category).To(Equal("Backend"))

				// Verify PostgreSQL skill
				pgSkill := findSkillByName(skills, "PostgreSQL")
				Expect(pgSkill).NotTo(BeNil())
				Expect(pgSkill.Name).To(Equal("PostgreSQL"))
				Expect(pgSkill.Category).To(Equal("Database"))
			})

			It("should store skills in repository", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skillRepo.skills).To(HaveLen(2))

				// Verify skills are stored
				Expect(skillRepo.skills["go"]).NotTo(BeNil()) // Lowercase key
				Expect(skillRepo.skills["postgresql"]).NotTo(BeNil())

				// Verify returned skills have IDs
				Expect(skills[0].ID).NotTo(BeEmpty())
				Expect(skills[1].ID).NotTo(BeEmpty())
			})

			It("should link skills to events", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())

				goSkill := findSkillByName(skills, "Go")
				pgSkill := findSkillByName(skills, "PostgreSQL")

				// Verify event-skill links
				event1 := eventRepo.events["event-1"]
				Expect(event1.Skills).To(ConsistOf(goSkill.ID, pgSkill.ID))

				event2 := eventRepo.events["event-2"]
				Expect(event2.Skills).To(ConsistOf(goSkill.ID))
			})

			It("should set LastUsed to most recent event date", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())

				goSkill := findSkillByName(skills, "Go")
				Expect(goSkill.LastUsed).NotTo(BeNil())
				Expect(*goSkill.LastUsed).To(Equal(testTime.Add(24 * time.Hour))) // event-2 is more recent

				pgSkill := findSkillByName(skills, "PostgreSQL")
				Expect(pgSkill.LastUsed).NotTo(BeNil())
				Expect(*pgSkill.LastUsed).To(Equal(testTime)) // Only event-1
			})
		})

		Context("when skill already exists", func() {
			var existingTime time.Time

			BeforeEach(func() {
				existingTime = testTime.Add(-7 * 24 * time.Hour) // 1 week ago

				// Pre-create a skill
				existingSkill := fixtures.SkillWith("skill-1", "Go", "Backend", "intermediate")
				existingSkill.LastUsed = &existingTime
				skillRepo.skills["go"] = existingSkill
				skillRepo.skillsByID["skill-1"] = existingSkill

				suggestions = []skillinference.SkillSuggestion{
					{
						Name:       "Go", // Same skill
						Category:   "Backend",
						Confidence: 0.95,
						EventIDs:   []string{"event-1"},
					},
				}

				event1 := fixtures.Event("event-1")
				event1.Date = testTime
				eventRepo.events["event-1"] = event1
			})

			It("should reuse existing skill (case-insensitive)", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(1))

				// Should return the existing skill
				Expect(skills[0].ID).To(Equal("skill-1"))
				Expect(skillRepo.skills).To(HaveLen(1)) // No new skill created
			})

			It("should update LastUsed to more recent date", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())

				skill := skills[0]
				Expect(skill.LastUsed).NotTo(BeNil())
				Expect(*skill.LastUsed).To(Equal(testTime)) // Updated from 1 week ago
			})

			It("should link to new events", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())

				event := eventRepo.events["event-1"]
				Expect(event.Skills).To(ConsistOf(skills[0].ID))
			})
		})

		Context("when skill exists with different case", func() {
			BeforeEach(func() {
				existingSkill := fixtures.SkillWith("skill-1", "go", "Backend", "intermediate")
				skillRepo.skills["go"] = existingSkill
				skillRepo.skillsByID["skill-1"] = existingSkill

				suggestions = []skillinference.SkillSuggestion{
					{
						Name:     "Go", // Different case
						Category: "Backend",
						EventIDs: []string{"event-1"},
					},
				}

				event1 := fixtures.Event("event-1")
				event1.Date = testTime
				eventRepo.events["event-1"] = event1
			})

			It("should match existing skill case-insensitively", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(1))
				Expect(skills[0].ID).To(Equal("skill-1"))
			})
		})

		Context("when suggestions have duplicates", func() {
			BeforeEach(func() {
				suggestions = []skillinference.SkillSuggestion{
					{
						Name:       "Go",
						Category:   "Backend",
						Confidence: 0.95,
						EventIDs:   []string{"event-1"},
					},
					{
						Name:       "Go", // Duplicate
						Category:   "Backend",
						Confidence: 0.75, // Lower confidence (ignored)
						EventIDs:   []string{"event-2"},
					},
				}

				event1 := fixtures.Event("event-1")
				event1.Date = testTime
				eventRepo.events["event-1"] = event1
				event2 := fixtures.Event("event-2")
				event2.Date = testTime.Add(24 * time.Hour)
				eventRepo.events["event-2"] = event2
			})

			It("should deduplicate suggestions", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(1)) // Only one skill
			})

			It("should merge event IDs", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())

				skillID := skills[0].ID

				// Both events should be linked
				Expect(eventRepo.events["event-1"].Skills).To(ContainElement(skillID))
				Expect(eventRepo.events["event-2"].Skills).To(ContainElement(skillID))
			})

			It("should use most recent event for LastUsed", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills[0].LastUsed).NotTo(BeNil())
				Expect(*skills[0].LastUsed).To(Equal(testTime.Add(24 * time.Hour))) // Most recent
			})
		})

		Context("when repository returns ErrSkillNotFound for new skills", func() {
			BeforeEach(func() {
				skillRepo.notFoundError = career_repo.ErrSkillNotFound

				suggestions = []skillinference.SkillSuggestion{
					{
						Name:       "Go",
						Category:   "Backend",
						Confidence: 0.95,
						EventIDs:   []string{"event-1", "event-2"},
					},
					{
						Name:       "PostgreSQL",
						Category:   "Database",
						Confidence: 0.75,
						EventIDs:   []string{"event-1"},
					},
				}

				event1 := fixtures.Event("event-1")
				event1.Date = testTime
				eventRepo.events["event-1"] = event1
				event2 := fixtures.Event("event-2")
				event2.Date = testTime.Add(24 * time.Hour)
				eventRepo.events["event-2"] = event2
			})

			It("should create new skills successfully", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(2))

				goSkill := findSkillByName(skills, "Go")
				Expect(goSkill).NotTo(BeNil())
				Expect(goSkill.Category).To(Equal("Backend"))

				pgSkill := findSkillByName(skills, "PostgreSQL")
				Expect(pgSkill).NotTo(BeNil())
				Expect(pgSkill.Category).To(Equal("Database"))
			})

			It("should link skills to events", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())

				goSkill := findSkillByName(skills, "Go")
				pgSkill := findSkillByName(skills, "PostgreSQL")

				event1 := eventRepo.events["event-1"]
				Expect(event1.Skills).To(ConsistOf(goSkill.ID, pgSkill.ID))

				event2 := eventRepo.events["event-2"]
				Expect(event2.Skills).To(ConsistOf(goSkill.ID))
			})

			It("should handle mix of existing and new skills", func() {
				existingSkill := fixtures.SkillWith("existing-1", "Go", "Backend", "intermediate")
				skillRepo.skills["go"] = existingSkill
				skillRepo.skillsByID["existing-1"] = existingSkill

				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(2))

				goSkill := findSkillByName(skills, "Go")
				Expect(goSkill.ID).To(Equal("existing-1"))

				pgSkill := findSkillByName(skills, "PostgreSQL")
				Expect(pgSkill).NotTo(BeNil())
				Expect(pgSkill.ID).NotTo(Equal("existing-1"))
			})
		})

		Context("when no suggestions provided", func() {
			It("should return empty slice", func() {
				skills, err := service.CreateSkillsFromSuggestions(ctx, []skillinference.SkillSuggestion{})

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(BeEmpty())
			})
		})

		Context("when context is cancelled", func() {
			It("should return context error", func() {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately

				suggestions = []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", EventIDs: []string{"event-1"}},
				}

				_, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).To(Equal(context.Canceled))
			})
		})

		Context("when repository fails", func() {
			BeforeEach(func() {
				suggestions = []skillinference.SkillSuggestion{
					{Name: "Go", Category: "Backend", EventIDs: []string{"event-1"}},
				}
				event1 := fixtures.Event("event-1")
				event1.Date = testTime
				eventRepo.events["event-1"] = event1
			})

			It("should return error when skill creation fails", func() {
				skillRepo.createError = errors.New("database error")

				_, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("database error"))
			})

			It("should return error when skill update fails", func() {
				skillRepo.updateError = errors.New("update failed")

				// Pre-create skill to trigger update path
				existingSkill := fixtures.SkillWith("skill-1", "Go", "Backend", "intermediate")
				skillRepo.skills["go"] = existingSkill
				skillRepo.skillsByID["skill-1"] = existingSkill

				_, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("update failed"))
			})

			It("should return error when event linking fails", func() {
				eventRepo.updateError = errors.New("link failed")

				_, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("link failed"))
			})

			It("should return error when event not found", func() {
				// Remove event-1 from repo (added by BeforeEach)
				delete(eventRepo.events, "event-1")

				_, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("event not found"))
			})
		})
	})
})

// Helper function to find skill by name.
func findSkillByName(skills []*career.Skill, name string) *career.Skill {
	for _, skill := range skills {
		if skill.Name == name {
			return skill
		}
	}
	return nil
}

// Mock repositories.
type mockSkillRepository struct {
	skills        map[string]*career.Skill // Key: lowercase name
	skillsByID    map[string]*career.Skill
	createError   error
	updateError   error
	notFoundError error
	nextID        int
}

func (m *mockSkillRepository) Create(ctx context.Context, skill *career.Skill) error {
	if m.createError != nil {
		return m.createError
	}

	m.nextID++
	skill.ID = generateID("skill", m.nextID)

	key := normalizeSkillName(skill.Name)
	m.skills[key] = skill
	m.skillsByID[skill.ID] = skill

	return nil
}

func (m *mockSkillRepository) Update(ctx context.Context, skill *career.Skill) error {
	if m.updateError != nil {
		return m.updateError
	}

	key := normalizeSkillName(skill.Name)
	m.skills[key] = skill
	m.skillsByID[skill.ID] = skill

	return nil
}

func (m *mockSkillRepository) GetByName(ctx context.Context, name string) (*career.Skill, error) {
	key := normalizeSkillName(name)
	skill, exists := m.skills[key]
	if !exists {
		return nil, m.notFoundError
	}
	return skill, nil
}

func (m *mockSkillRepository) GetByID(ctx context.Context, id string) (*career.Skill, error) {
	skill, exists := m.skillsByID[id]
	if !exists {
		return nil, errors.New("skill not found")
	}
	return skill, nil
}

type mockEventRepository struct {
	events      map[string]*career.Event
	updateError error
}

func (m *mockEventRepository) GetByID(ctx context.Context, id string) (*career.Event, error) {
	event, exists := m.events[id]
	if !exists {
		return nil, errors.New("event not found")
	}
	return event, nil
}

func (m *mockEventRepository) Update(ctx context.Context, event *career.Event) error {
	if m.updateError != nil {
		return m.updateError
	}

	m.events[event.ID] = event
	return nil
}

func (m *mockEventRepository) LinkSkill(ctx context.Context, eventID string, skillID string) error {
	if m.updateError != nil {
		return m.updateError
	}

	event, exists := m.events[eventID]
	if !exists {
		return errors.New("event not found")
	}

	// Add skill ID to event.Skills if not already present
	for _, sid := range event.Skills {
		if sid == skillID {
			return nil // Already linked
		}
	}

	event.Skills = append(event.Skills, skillID)
	return nil
}

func (m *mockEventRepository) UnlinkSkill(ctx context.Context, eventID string, skillID string) error {
	if m.updateError != nil {
		return m.updateError
	}

	event, exists := m.events[eventID]
	if !exists {
		return errors.New("event not found")
	}

	newSkills := make([]string, 0, len(event.Skills))
	for _, sid := range event.Skills {
		if sid != skillID {
			newSkills = append(newSkills, sid)
		}
	}

	event.Skills = newSkills
	return nil
}

// Helper functions.
func normalizeSkillName(name string) string {
	return strings.ToLower(name)
}

func generateID(prefix string, num int) string {
	return fmt.Sprintf("%s-%d", prefix, num)
}
