package skillinference_test

import (
	"context"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careersql "github.com/baphled/kariya/internal/repository/career/sql"
	"github.com/baphled/kariya/internal/service/career/skillinference"
	"github.com/baphled/kariya/internal/testutil/fixtures"
)

var _ = Describe("Skill Persistence (SQL-backed)", func() {
	var (
		service skillinference.SkillInferenceService
		repos   *careerrepo.Repositories
		ctx     context.Context
	)

	BeforeEach(func() {
		tmpDir := GinkgoT().TempDir()
		dbPath := filepath.Join(tmpDir, "test_skills.db")

		var err error
		repos, err = careersql.NewRepositoriesFromPath(dbPath)
		Expect(err).NotTo(HaveOccurred())

		service = skillinference.NewSkillInferenceService(repos.Skill, repos.Event)
		ctx = context.Background() //nolint:fatcontext // test setup
	})

	AfterEach(func() {
		if repos != nil {
			repos.Close()
		}
	})

	Describe("CreateSkillsFromSuggestions against SQL", func() {
		Context("when creating new skills", func() {
			It("should persist skills to the SQL database", func() {
				testTime := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
				event1 := fixtures.EventWith("event-1", "Built REST API using Go", "", "")
				event1.Date = testTime
				event2 := fixtures.EventWith("event-2", "Designed PostgreSQL schema", "", "")
				event2.Date = testTime.Add(24 * time.Hour)

				Expect(repos.Event.Create(ctx, event1)).To(Succeed())
				Expect(repos.Event.Create(ctx, event2)).To(Succeed())

				suggestions := []skillinference.SkillSuggestion{
					{
						Name:       "Go",
						Category:   "Backend",
						Confidence: 0.95,
						EventIDs:   []string{"event-1", "event-2"},
					},
					{
						Name:       "PostgreSQL",
						Category:   "Database",
						Confidence: 0.85,
						EventIDs:   []string{"event-1"},
					},
				}

				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(2))

				goSkill := findSkillByName(skills, "Go")
				Expect(goSkill).NotTo(BeNil())
				Expect(goSkill.ID).NotTo(BeEmpty())
				Expect(goSkill.Category).To(Equal("Backend"))

				pgSkill := findSkillByName(skills, "PostgreSQL")
				Expect(pgSkill).NotTo(BeNil())
				Expect(pgSkill.ID).NotTo(BeEmpty())
				Expect(pgSkill.Category).To(Equal("Database"))

				storedGo, err := repos.Skill.GetByName(ctx, "Go")
				Expect(err).NotTo(HaveOccurred())
				Expect(storedGo.Name).To(Equal("Go"))

				storedPg, err := repos.Skill.GetByName(ctx, "PostgreSQL")
				Expect(err).NotTo(HaveOccurred())
				Expect(storedPg.Name).To(Equal("PostgreSQL"))
			})
		})

		Context("when skill already exists with different case", func() {
			It("should reuse the existing skill instead of creating a duplicate", func() {
				testTime := time.Date(2024, 6, 15, 0, 0, 0, 0, time.UTC)
				event1 := fixtures.EventWith("event-1", "Built API with Go", "", "")
				event1.Date = testTime
				Expect(repos.Event.Create(ctx, event1)).To(Succeed())

				existingSkill := fixtures.SkillWith("", "go", "Backend", "intermediate")
				Expect(repos.Skill.Create(ctx, existingSkill)).To(Succeed())

				suggestions := []skillinference.SkillSuggestion{
					{
						Name:       "Go",
						Category:   "Backend",
						Confidence: 0.95,
						EventIDs:   []string{"event-1"},
					},
				}

				skills, err := service.CreateSkillsFromSuggestions(ctx, suggestions)

				Expect(err).NotTo(HaveOccurred())
				Expect(skills).To(HaveLen(1))
				Expect(skills[0].ID).To(Equal(existingSkill.ID),
					"should reuse existing 'go' skill when suggestion is 'Go'")

				allSkills, err := repos.Skill.List(ctx, nil)
				Expect(err).NotTo(HaveOccurred())
				Expect(allSkills).To(HaveLen(1),
					"should NOT create a duplicate skill with different casing")
			})
		})

		Context("when looking up existing skills by name", func() {
			It("should find skill regardless of case", func() {
				existingSkill := fixtures.SkillWith("", "Go", "Backend", "intermediate")
				Expect(repos.Skill.Create(ctx, existingSkill)).To(Succeed())

				found, err := repos.Skill.GetByName(ctx, "go")
				Expect(err).NotTo(HaveOccurred())
				Expect(found).NotTo(BeNil())
				Expect(found.ID).To(Equal(existingSkill.ID))
			})
		})
	})
})
