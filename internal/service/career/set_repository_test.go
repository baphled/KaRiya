package career_test

import (
	"github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository Setters", func() {
	var (
		repo      *career.MemoryRepository
		svc       *careerservice.Service
		factRepo  *career.MemoryFactRepository
		burstRepo *career.MemoryBurstRepository
	)

	BeforeEach(func() {
		repo = career.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		factRepo = career.NewMemoryFactRepository()
		burstRepo = career.NewMemoryBurstRepository()
	})

	Describe("SetFactRepository", func() {
		It("sets the fact repository", func() {
			Expect(func() {
				svc.SetFactRepository(factRepo)
			}).NotTo(Panic())
		})
	})

	Describe("SetBurstRepository", func() {
		It("sets the burst repository", func() {
			Expect(func() {
				svc.SetBurstRepository(burstRepo)
			}).NotTo(Panic())
		})
	})
})
