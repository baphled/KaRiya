package career_test

import (
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repository Setters", func() {
	var (
		repo      *careermemory.EventRepository
		svc       *careerservice.Service
		factRepo  *careermemory.FactRepository
		burstRepo *careermemory.BurstRepository
	)

	BeforeEach(func() {
		repo = careermemory.NewEventRepository()
		svc = careerservice.NewService(repo)
		factRepo = careermemory.NewFactRepository()
		burstRepo = careermemory.NewBurstRepository()
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
