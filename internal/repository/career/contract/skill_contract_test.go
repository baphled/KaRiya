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

var _ = Describe("SkillRepository Contract", func() {
	var (
		memRepo career_repo.SkillRepository
		sqlRepo career_repo.SkillRepository
		ctx     context.Context
		db      *gorm.DB
	)

	BeforeEach(func() {
		ctx = context.Background()
		db = setupContractTestDB()

		memRepo = memory.NewSkillRepository()
		sqlRepo = repo_sql.NewSkillRepository(db)
	})

	runForBoth := func(name string, fn func(repo career_repo.SkillRepository)) {
		It(fmt.Sprintf("%s [memory]", name), func() { fn(memRepo) })
		It(fmt.Sprintf("%s [sql]", name), func() { fn(sqlRepo) })
	}

	Describe("Gap 2: Pagination limit=0 defaults to 100", func() {
		It("memory: returns at most 100 skills when limit is 0", func() {
			for i := range 110 {
				skill := fixtures.SkillWith(
					fmt.Sprintf("skill-%d", i),
					fmt.Sprintf("Skill Name %d", i),
					"backend",
					"intermediate",
				)
				err := memRepo.Create(ctx, skill)
				Expect(err).NotTo(HaveOccurred())
			}

			results, err := memRepo.List(ctx, fixtures.SkillListFiltersWithLimit(0, 0))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(BeNumerically("<=", 100))
		})
	})

	Describe("Gap 6: List ordering is deterministic", func() {
		runForBoth("returns skills in consistent order across calls", func(repo career_repo.SkillRepository) {
			now := time.Now()
			for i := range 5 {
				skill := fixtures.SkillWith(
					fmt.Sprintf("skill-%d", i),
					fmt.Sprintf("Skill %d", i),
					"backend",
					"intermediate",
				)
				skill.CreatedAt = now.Add(-time.Duration(i) * time.Hour)
				err := repo.Create(ctx, skill)
				Expect(err).NotTo(HaveOccurred())
			}

			results1, err := repo.List(ctx, fixtures.SkillListFilters())
			Expect(err).NotTo(HaveOccurred())
			results2, err := repo.List(ctx, fixtures.SkillListFilters())
			Expect(err).NotTo(HaveOccurred())

			Expect(results1).To(HaveLen(5))
			Expect(results2).To(HaveLen(5))
			for i := range results1 {
				Expect(results1[i].ID).To(Equal(results2[i].ID))
			}
		})
	})
})
