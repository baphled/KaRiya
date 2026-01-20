package fixtures_test

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstFactory", func() {
	BeforeEach(func() {
		fixtures.SetSeed(42)
	})

	Describe("MustCreate", func() {
		It("should create a burst with all required fields", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst).NotTo(BeNil())
			Expect(burst.ID).NotTo(BeEmpty())
			Expect(burst.Name).NotTo(BeEmpty())
			Expect(burst.EventIDs).To(HaveLen(2)) // Minimum required
			Expect(burst.CreatedAt).NotTo(BeZero())
			Expect(burst.UpdatedAt).NotTo(BeZero())
		})

		It("should create bursts with unique sequential IDs", func() {
			burst1 := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst2 := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst1.ID).NotTo(Equal(burst2.ID))
		})

		It("should create bursts with description", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst.Description).NotTo(BeEmpty())
		})

		It("should create unconfirmed bursts by default", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst.Confirmed).To(BeFalse())
			Expect(burst.ConfirmedAt).To(BeNil())
		})
	})

	Describe("MustCreateWithOption", func() {
		It("should allow overriding confirmed status", func() {
			burst := fixtures.BurstFactory.MustCreateWithOption(map[string]interface{}{
				"Confirmed": true,
			}).(*career.Burst)
			Expect(burst.Confirmed).To(BeTrue())
		})
	})
})

var _ = Describe("Burst Quick Helper", func() {
	Describe("Burst", func() {
		It("should create a minimal valid burst with given ID and event IDs", func() {
			burst := fixtures.Burst("burst-1", "evt-1", "evt-2")
			Expect(burst.ID).To(Equal("burst-1"))
			Expect(burst.Name).To(ContainSubstring("burst-1"))
			Expect(burst.EventIDs).To(Equal([]string{"evt-1", "evt-2"}))
			Expect(burst.Confirmed).To(BeFalse())
		})

		It("should use default event IDs if less than 2 provided", func() {
			burst := fixtures.Burst("burst-1", "evt-1")
			Expect(burst.EventIDs).To(HaveLen(2))
		})
	})

	Describe("BurstConfirmed", func() {
		It("should create a confirmed burst", func() {
			burst := fixtures.BurstConfirmed("burst-1", "evt-1", "evt-2")
			Expect(burst.Confirmed).To(BeTrue())
			Expect(burst.ConfirmedAt).NotTo(BeNil())
		})
	})
})

var _ = Describe("Bursts Batch Helper", func() {
	BeforeEach(func() {
		fixtures.SetSeed(42)
	})

	Describe("Bursts", func() {
		It("should create n bursts", func() {
			events := fixtures.Events(6)
			bursts := fixtures.Bursts(3, events)
			Expect(bursts).To(HaveLen(3))
		})

		It("should link bursts to provided events", func() {
			events := fixtures.Events(4)
			bursts := fixtures.Bursts(2, events)
			Expect(bursts[0].EventIDs).To(HaveLen(2))
			Expect(events[0].ID).To(Equal(bursts[0].EventIDs[0]))
		})
	})
})
