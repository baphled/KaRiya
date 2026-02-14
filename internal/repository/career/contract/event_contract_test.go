package contract_test

import (
	"context"
	"fmt"
	"time"

	career_repo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/repository/career/memory"
	repo_sql "github.com/baphled/kariya/internal/repository/career/sql"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"
)

var _ = Describe("EventRepository Contract", func() {
	var (
		memRepo career_repo.EventRepository
		sqlRepo career_repo.EventRepository
		ctx     context.Context
		db      *gorm.DB
	)

	BeforeEach(func() {
		ctx = context.Background()
		db = setupContractTestDB()

		memEventRepo := memory.NewEventRepository()
		memSkillRepo := memory.NewSkillRepository()
		memEventRepo.SetSkillRepository(memSkillRepo)
		memSkillRepo.SetEventRepository(memEventRepo)
		memRepo = memEventRepo

		sqlRepo = repo_sql.NewEventRepository(db)
	})

	runForBoth := func(name string, fn func(repo career_repo.EventRepository)) {
		It(fmt.Sprintf("%s [memory]", name), func() { fn(memRepo) })
		It(fmt.Sprintf("%s [sql]", name), func() { fn(sqlRepo) })
	}

	Describe("Gap 2: Pagination limit=0 defaults to 100", func() {
		runForBoth("returns at most 100 items when limit is 0", func(repo career_repo.EventRepository) {
			for i := range 110 {
				event := fixtures.Event(fmt.Sprintf("event-%d", i))
				event.Text = fmt.Sprintf("Test career event number %d for pagination", i)
				event.Date = time.Now().Add(-time.Duration(i) * time.Hour)
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			results, err := repo.List(ctx, *fixtures.EventListFiltersWithLimit(0, 0))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(BeNumerically("<=", 100))
		})
	})

	Describe("Gap 3: Tag filtering uses exact match", func() {
		runForBoth("returns only exact tag matches", func(repo career_repo.EventRepository) {
			event1 := fixtures.Event("event-go")
			event1.Text = "Implemented microservice in Go language"
			event1.Tags = []string{"technical"}
			err := repo.Create(ctx, event1)
			Expect(err).NotTo(HaveOccurred())

			event2 := fixtures.Event("event-other")
			event2.Text = "Led project management initiative"
			event2.Tags = []string{"leadership"}
			err = repo.Create(ctx, event2)
			Expect(err).NotTo(HaveOccurred())

			results, err := repo.List(ctx, *fixtures.EventListFiltersWithTags([]string{"technical"}))
			Expect(err).NotTo(HaveOccurred())
			Expect(results).To(HaveLen(1))
			Expect(results[0].ID).To(Equal("event-go"))
		})
	})

	Describe("Gap 6: List ordering is deterministic", func() {
		runForBoth("returns events in consistent order across calls", func(repo career_repo.EventRepository) {
			now := time.Now()
			for i := range 5 {
				event := fixtures.Event(fmt.Sprintf("event-%d", i))
				event.Text = fmt.Sprintf("Deterministic ordering test event %d", i)
				event.Date = now.Add(-time.Duration(i) * 24 * time.Hour)
				event.CreatedAt = now.Add(-time.Duration(i) * time.Hour)
				err := repo.Create(ctx, event)
				Expect(err).NotTo(HaveOccurred())
			}

			results1, err := repo.List(ctx, *fixtures.EventListFilters())
			Expect(err).NotTo(HaveOccurred())
			results2, err := repo.List(ctx, *fixtures.EventListFilters())
			Expect(err).NotTo(HaveOccurred())

			Expect(results1).To(HaveLen(5))
			Expect(results2).To(HaveLen(5))
			for i := range results1 {
				Expect(results1[i].ID).To(Equal(results2[i].ID))
			}
		})
	})

	Describe("Gap 5: Delete cascade cleans up skill links", func() {
		It("memory: deleting event removes skill associations", func() {
			memEventRepo := memRepo.(*memory.EventRepository)
			memSkillRepo := memory.NewSkillRepository()
			memSkillRepo.SetEventRepository(memEventRepo)
			memEventRepo.SetSkillRepository(memSkillRepo)

			event := fixtures.Event("event-cascade")
			event.Text = "Cascade delete test event for skills"
			err := memEventRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			skill := fixtures.Skill("skill-cascade")
			err = memSkillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			err = memEventRepo.LinkSkill(ctx, "event-cascade", "skill-cascade")
			Expect(err).NotTo(HaveOccurred())

			err = memEventRepo.Delete(ctx, "event-cascade")
			Expect(err).NotTo(HaveOccurred())

			skills, err := memSkillRepo.GetSkillsForEvent(ctx, "event-cascade")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(BeEmpty())
		})
	})

	Describe("Gap 4: Skill sync roundtrip", func() {
		It("memory: LinkSkill + GetSkillsForEvent returns linked skill", func() {
			memEventRepo := memRepo.(*memory.EventRepository)
			memSkillRepo := memory.NewSkillRepository()
			memSkillRepo.SetEventRepository(memEventRepo)
			memEventRepo.SetSkillRepository(memSkillRepo)

			event := fixtures.Event("event-link")
			event.Text = "Skill link roundtrip test event"
			err := memEventRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			skill := fixtures.Skill("skill-link")
			err = memSkillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			err = memEventRepo.LinkSkill(ctx, "event-link", "skill-link")
			Expect(err).NotTo(HaveOccurred())

			skills, err := memSkillRepo.GetSkillsForEvent(ctx, "event-link")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].ID).To(Equal("skill-link"))
		})

		It("sql: LinkSkill + GetSkillsForEvent returns linked skill", func() {
			sqlEventRepo := sqlRepo.(*repo_sql.EventRepository)
			sqlSkillRepo := repo_sql.NewSkillRepository(db)

			event := fixtures.Event("event-link")
			event.Text = "Skill link roundtrip test event"
			err := sqlEventRepo.Create(ctx, event)
			Expect(err).NotTo(HaveOccurred())

			skill := fixtures.Skill("skill-link")
			err = sqlSkillRepo.Create(ctx, skill)
			Expect(err).NotTo(HaveOccurred())

			err = sqlEventRepo.LinkSkill(ctx, "event-link", "skill-link")
			Expect(err).NotTo(HaveOccurred())

			skills, err := sqlSkillRepo.GetSkillsForEvent(ctx, "event-link")
			Expect(err).NotTo(HaveOccurred())
			Expect(skills).To(HaveLen(1))
			Expect(skills[0].ID).To(Equal("skill-link"))
		})
	})

	Describe("Gap 1: Create validation rejects short text", func() {
		It("memory: rejects event with text shorter than 10 characters", func() {
			event := fixtures.EventWith("invalid-event", "Short", "", "")
			err := memRepo.Create(ctx, event)
			Expect(err).To(HaveOccurred())
		})
	})
})
