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

var _ = Describe("FactRepository Contract", func() {
	var (
		memRepo career_repo.FactRepository
		sqlRepo career_repo.FactRepository
		ctx     context.Context
		db      *gorm.DB
	)

	BeforeEach(func() {
		ctx = context.Background()
		db = setupContractTestDB()

		memRepo = memory.NewFactRepository()
		sqlRepo = repo_sql.NewFactRepository(db)
	})

	runForBoth := func(name string, fn func(repo career_repo.FactRepository)) {
		It(fmt.Sprintf("%s [memory]", name), func() { fn(memRepo) })
		It(fmt.Sprintf("%s [sql]", name), func() { fn(sqlRepo) })
	}

	Describe("Gap 2: Pagination limit=0 defaults to 100", func() {
		runForBoth("returns at most 100 facts when limit is 0", func(repo career_repo.FactRepository) {
			for i := range 110 {
				fact := fixtures.Fact(fmt.Sprintf("fact-%d", i), fmt.Sprintf("event-%d", i))
				fact.Text = fmt.Sprintf("Reduced system latency by %d percent through redesign", i+10)
				err := repo.Create(ctx, fact)
				Expect(err).NotTo(HaveOccurred())
			}

			results, err := repo.List(ctx, *fixtures.FactListFiltersWithLimit(0, 0))
			Expect(err).NotTo(HaveOccurred())
			Expect(len(results)).To(BeNumerically("<=", 100))
		})
	})

	Describe("Gap 6: List ordering is deterministic", func() {
		runForBoth("returns facts in consistent order across calls", func(repo career_repo.FactRepository) {
			now := time.Now()
			for i := range 5 {
				fact := fixtures.Fact(fmt.Sprintf("fact-%d", i), "event-source")
				fact.Text = fmt.Sprintf("Deterministic ordering test fact number %d", i)
				fact.CreatedAt = now.Add(-time.Duration(i) * time.Hour)
				err := repo.Create(ctx, fact)
				Expect(err).NotTo(HaveOccurred())
			}

			results1, err := repo.List(ctx, *fixtures.FactListFilters())
			Expect(err).NotTo(HaveOccurred())
			results2, err := repo.List(ctx, *fixtures.FactListFilters())
			Expect(err).NotTo(HaveOccurred())

			Expect(results1).To(HaveLen(5))
			Expect(results2).To(HaveLen(5))
			for i := range results1 {
				Expect(results1[i].ID).To(Equal(results2[i].ID))
			}
		})
	})
})
