package burst_fact

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TemporalGrouper", func() {
	var grouper *TemporalGrouper

	BeforeEach(func() {
		grouper = NewTemporalGrouper()
	})

	Context("when grouping events within 6-month window", func() {
		It("should group events within 6 months as related", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			event1Date := baseDate
			event2Date := baseDate.AddDate(0, 3, 0) // 3 months later

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeTrue())
		})

		It("should not group events more than 6 months apart", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			event1Date := baseDate
			event2Date := baseDate.AddDate(0, 7, 0) // 7 months later

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeFalse())
		})

		It("should handle exactly 6 months as related", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			event1Date := baseDate
			event2Date := baseDate.AddDate(0, 6, 0) // exactly 6 months

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeTrue())
		})

		It("should handle events in reverse chronological order", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			event1Date := baseDate.AddDate(0, 3, 0) // 3 months later
			event2Date := baseDate                  // earlier

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeTrue())
		})

		It("should handle same-day events as related", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			event1Date := baseDate
			event2Date := baseDate // same day

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeTrue())
		})

		It("should handle 1-day events as related", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			event1Date := baseDate
			event2Date := baseDate.AddDate(0, 0, 1) // 1 day later

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeTrue())
		})

		It("should handle edge case at 6 month + 1 day as not related", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			event1Date := baseDate
			event2Date := baseDate.AddDate(0, 6, 1) // 6 months + 1 day

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeFalse())
		})

		It("should handle year boundaries correctly", func() {
			event1Date := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
			event2Date := time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC) // 5 months later

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeTrue())
		})

		It("should handle year boundaries crossing 6+ months", func() {
			event1Date := time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC)
			event2Date := time.Date(2025, 7, 1, 0, 0, 0, 0, time.UTC) // 7 months later

			isRelated := grouper.IsTemporallyRelated(event1Date, event2Date)

			Expect(isRelated).To(BeFalse())
		})
	})

	Context("when computing time differences", func() {
		It("should return 0 for same date", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

			diff := grouper.MonthsDifference(baseDate, baseDate)

			Expect(diff).To(Equal(0))
		})

		It("should return 3 for 3 months difference", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			laterDate := baseDate.AddDate(0, 3, 0)

			diff := grouper.MonthsDifference(baseDate, laterDate)

			Expect(diff).To(Equal(3))
		})

		It("should return 12 for 1 year difference", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			laterDate := baseDate.AddDate(1, 0, 0)

			diff := grouper.MonthsDifference(baseDate, laterDate)

			Expect(diff).To(Equal(12))
		})

		It("should handle reverse order (earlier date first)", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			earlierDate := baseDate.AddDate(0, -3, 0)

			diff := grouper.MonthsDifference(baseDate, earlierDate)

			Expect(diff).To(Equal(3))
		})

		It("should handle leap year correctly", func() {
			// 2024 is a leap year
			baseDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
			laterDate := baseDate.AddDate(0, 1, 0) // 1 month later (March)

			diff := grouper.MonthsDifference(baseDate, laterDate)

			Expect(diff).To(Equal(1))
		})

		It("should handle month boundaries correctly", func() {
			date1 := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)
			date2 := time.Date(2025, 2, 28, 0, 0, 0, 0, time.UTC) // February has 28 days

			diff := grouper.MonthsDifference(date1, date2)

			Expect(diff).To(Equal(1))
		})
	})

	Context("when grouping multiple events", func() {
		It("should identify temporal clusters of events", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			events := []time.Time{
				baseDate,                  // Jan 1
				baseDate.AddDate(0, 1, 0), // Feb 1 (1 month later)
				baseDate.AddDate(0, 2, 0), // Mar 1 (2 months later)
				baseDate.AddDate(0, 8, 1), // Sep 2 (8 months + 1 day - new cluster, >6 months from Mar 1)
				baseDate.AddDate(0, 9, 1), // Oct 2 (9 months + 1 day)
			}

			clusters := grouper.GroupEventsByTemporal(events)

			Expect(len(clusters)).To(Equal(2))
			Expect(len(clusters[0])).To(Equal(3)) // Jan, Feb, Mar
			Expect(len(clusters[1])).To(Equal(2)) // Sep 2, Oct 2
		})

		It("should handle empty event list", func() {
			events := []time.Time{}

			clusters := grouper.GroupEventsByTemporal(events)

			Expect(clusters).To(BeEmpty())
		})

		It("should handle single event", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			events := []time.Time{baseDate}

			clusters := grouper.GroupEventsByTemporal(events)

			Expect(len(clusters)).To(Equal(1))
			Expect(len(clusters[0])).To(Equal(1))
		})

		It("should handle all events in single cluster", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			events := []time.Time{
				baseDate,
				baseDate.AddDate(0, 1, 0),
				baseDate.AddDate(0, 2, 0),
				baseDate.AddDate(0, 3, 0),
				baseDate.AddDate(0, 4, 0),
			}

			clusters := grouper.GroupEventsByTemporal(events)

			Expect(len(clusters)).To(Equal(1))
			Expect(len(clusters[0])).To(Equal(5))
		})

		It("should handle all events as separate clusters (>6 months apart)", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			events := []time.Time{
				baseDate,
				baseDate.AddDate(0, 7, 0),
				baseDate.AddDate(0, 14, 0),
				baseDate.AddDate(0, 21, 0),
			}

			clusters := grouper.GroupEventsByTemporal(events)

			Expect(len(clusters)).To(Equal(4))
			for _, cluster := range clusters {
				Expect(len(cluster)).To(Equal(1))
			}
		})

		It("should preserve chronological order within clusters", func() {
			baseDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			events := []time.Time{
				baseDate.AddDate(0, 2, 0),
				baseDate,
				baseDate.AddDate(0, 1, 0),
			}

			clusters := grouper.GroupEventsByTemporal(events)

			Expect(len(clusters)).To(Equal(1))
			Expect(len(clusters[0])).To(Equal(3))
			// Check chronological order
			Expect(clusters[0][0].Before(clusters[0][1])).To(BeTrue())
			Expect(clusters[0][1].Before(clusters[0][2])).To(BeTrue())
		})
	})

	Context("when handling edge cases", func() {
		It("should handle zero time values gracefully", func() {
			zeroTime := time.Time{}
			normalTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

			// Should not panic
			isRelated := grouper.IsTemporallyRelated(zeroTime, normalTime)

			Expect(isRelated).To(BeFalse())
		})

		It("should handle very old dates", func() {
			oldDate := time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
			recentDate := oldDate.AddDate(0, 3, 0)

			isRelated := grouper.IsTemporallyRelated(oldDate, recentDate)

			Expect(isRelated).To(BeTrue())
		})

		It("should handle future dates", func() {
			futureDate := time.Date(2100, 1, 1, 0, 0, 0, 0, time.UTC)
			anotherFutureDate := futureDate.AddDate(0, 3, 0)

			isRelated := grouper.IsTemporallyRelated(futureDate, anotherFutureDate)

			Expect(isRelated).To(BeTrue())
		})
	})
})
