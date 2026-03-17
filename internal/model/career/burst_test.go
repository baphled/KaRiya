package career

import (
	"time"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Burst Model", func() {
	Describe("TableName", func() {
		It("returns bursts", func() {
			Expect(Burst{}.TableName()).To(Equal("bursts"))
		})
	})

	Describe("ToDomain", func() {
		It("converts all fields to domain Burst", func() {
			now := time.Now()
			confirmedAt := now.Add(-time.Hour)
			model := &Burst{
				ID:          "burst-1",
				Name:        "Platform Initiative",
				Description: "Series of platform improvements",
				EventIDs:    StringSlice{"evt-1", "evt-2", "evt-3"},
				Confirmed:   true,
				ConfirmedAt: &confirmedAt,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			domain := model.ToDomain()

			Expect(domain.ID).To(Equal("burst-1"))
			Expect(domain.Name).To(Equal("Platform Initiative"))
			Expect(domain.Description).To(Equal("Series of platform improvements"))
			Expect(domain.EventIDs).To(Equal([]string{"evt-1", "evt-2", "evt-3"}))
			Expect(domain.Confirmed).To(BeTrue())
			Expect(domain.ConfirmedAt).NotTo(BeNil())
			Expect(*domain.ConfirmedAt).To(Equal(confirmedAt))
			Expect(domain.CreatedAt).To(Equal(now))
			Expect(domain.UpdatedAt).To(Equal(now))
		})

		It("handles unconfirmed burst with nil ConfirmedAt", func() {
			model := &Burst{
				ID:        "burst-2",
				Name:      "Unconfirmed Burst",
				EventIDs:  StringSlice{"evt-1", "evt-2"},
				Confirmed: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			domain := model.ToDomain()

			Expect(domain.Confirmed).To(BeFalse())
			Expect(domain.ConfirmedAt).To(BeNil())
		})
	})

	Describe("BurstFromDomain", func() {
		It("converts all fields from domain Burst", func() {
			domainBurst := fixtures.BurstConfirmed("burst-d1", "evt-a", "evt-b")
			domainBurst.Description = "Database migration activities"

			model := BurstFromDomain(domainBurst)

			Expect(model.ID).To(Equal("burst-d1"))
			Expect(model.Name).To(Equal(domainBurst.Name))
			Expect(model.Description).To(Equal("Database migration activities"))
			Expect([]string(model.EventIDs)).To(Equal([]string{"evt-a", "evt-b"}))
			Expect(model.Confirmed).To(BeTrue())
			Expect(model.ConfirmedAt).NotTo(BeNil())
			Expect(model.CreatedAt).To(Equal(domainBurst.CreatedAt))
			Expect(model.UpdatedAt).To(Equal(domainBurst.UpdatedAt))
		})

		It("handles unconfirmed burst", func() {
			domainBurst := fixtures.Burst("burst-d2", "evt-x", "evt-y")

			model := BurstFromDomain(domainBurst)

			Expect(model.Confirmed).To(BeFalse())
			Expect(model.ConfirmedAt).To(BeNil())
		})
	})

	Describe("ToDomain-FromDomain round trip", func() {
		It("preserves all Burst data including confirmation", func() {
			original := fixtures.BurstConfirmed("burst-rt", "evt-r1", "evt-r2")
			original.Description = "Testing round trip"
			original.EventIDs = append(original.EventIDs, "evt-r3")
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := BurstFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.ID).To(Equal(original.ID))
			Expect(restored.Name).To(Equal(original.Name))
			Expect(restored.Description).To(Equal(original.Description))
			Expect(restored.EventIDs).To(Equal(original.EventIDs))
			Expect(restored.Confirmed).To(Equal(original.Confirmed))
			Expect(restored.ConfirmedAt).NotTo(BeNil())
			Expect(restored.CreatedAt).To(Equal(original.CreatedAt))
			Expect(restored.UpdatedAt).To(Equal(original.UpdatedAt))
		})

		It("preserves unconfirmed state through round trip", func() {
			original := fixtures.Burst("burst-rt-unconfirmed", "evt-u1", "evt-u2")
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := BurstFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.Confirmed).To(BeFalse())
			Expect(restored.ConfirmedAt).To(BeNil())
		})
	})
})
