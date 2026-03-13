package display_test

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/ui/display"
)

var _ = Describe("Burst display type", func() {
	Describe("BurstFromDomain", func() {
		It("converts a fully populated burst", func() {
			confirmedAt := time.Date(2024, time.March, 20, 10, 0, 0, 0, time.UTC)
			createdAt := time.Date(2024, time.March, 1, 10, 0, 0, 0, time.UTC)
			updatedAt := time.Date(2024, time.March, 2, 10, 0, 0, 0, time.UTC)

			burst := fixtures.BurstConfirmed("burst-1", "event-1", "event-2")
			burst.Name = "Platform Modernisation"
			burst.Description = "Related platform delivery work"
			burst.ConfirmedAt = &confirmedAt
			burst.CreatedAt = createdAt
			burst.UpdatedAt = updatedAt

			result := display.BurstFromDomain(burst)

			Expect(result).To(Equal(display.Burst{
				ID:          "burst-1",
				Name:        "Platform Modernisation",
				Description: "Related platform delivery work",
				EventIDs:    []string{"event-1", "event-2"},
				Confirmed:   true,
				ConfirmedAt: &confirmedAt,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
			}))
		})

		It("handles nil input gracefully", func() {
			Expect(display.BurstFromDomain(nil)).To(Equal(display.Burst{}))
		})

		It("preserves nil optional fields", func() {
			burst := fixtures.Burst("burst-1")
			result := display.BurstFromDomain(burst)

			Expect(result.ConfirmedAt).To(BeNil())
		})
	})

	Describe("BurstsFromDomain", func() {
		It("converts a slice of bursts", func() {
			b1 := fixtures.Burst("burst-1")
			b1.Name = "First"
			b2 := fixtures.Burst("burst-2")
			b2.Name = "Second"
			bursts := []*career.Burst{b1, b2}

			result := display.BurstsFromDomain(bursts)

			Expect(result).To(HaveLen(2))
			Expect(result[0].ID).To(Equal("burst-1"))
			Expect(result[0].Name).To(Equal("First"))
			Expect(result[1].ID).To(Equal("burst-2"))
			Expect(result[1].Name).To(Equal("Second"))
		})

		It("returns nil for nil input", func() {
			Expect(display.BurstsFromDomain(nil)).To(BeNil())
		})

		It("returns empty slice for empty input", func() {
			Expect(display.BurstsFromDomain([]*career.Burst{})).To(Equal([]display.Burst{}))
		})
	})
})
