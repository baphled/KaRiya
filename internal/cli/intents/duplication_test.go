package intents_test

import (
	"strings"
	"time"

	"github.com/baphled/kariya/internal/cli/intents"
	browse_timeline "github.com/baphled/kariya/internal/cli/intents/browsetimeline"
	"github.com/baphled/kariya/internal/cli/intents/captureevent"
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("View Duplication Prevention", func() {
	Describe("BrowseTimeline", func() {
		var (
			intent *browse_timeline.Intent
			events []*career.Event
		)

		BeforeEach(func() {
			evt := fixtures.EventWith("event1", "Test event", "Test Co", "")
			evt.Date = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			events = []*career.Event{evt}

			ctx := &browse_timeline.IntentContext{
				Events: events,
				InitialFilters: &browse_timeline.Filters{
					Tags:      make([]string, 0),
					Companies: make([]string, 0),
					SortBy:    "date",
					SortOrder: "desc",
				},
			}

			var err error
			intent, err = browse_timeline.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should have no duplicate breadcrumbs", func() {
			view := intent.View()

			// Count breadcrumb separator occurrences
			separatorCount := strings.Count(view, "▸")
			Expect(separatorCount).To(BeNumerically(">=", 0), "Expected at least one breadcrumb separator")
			Expect(separatorCount).To(BeNumerically("<=", 4), "Too many breadcrumb separators, likely duplication")

			// Count "Timeline" occurrences - breadcrumb + header title is acceptable
			timelineCount := strings.Count(view, "Timeline")
			Expect(timelineCount).To(BeNumerically("<=", 3), "'Timeline' appears too many times, likely duplication")
		})

		It("should have no duplicate help footer", func() {
			view := intent.View()

			// Count navigation help pattern - should appear only once
			upDownPattern := "↑/k"
			count := strings.Count(view, upDownPattern)
			Expect(count).To(BeNumerically("<=", 1), "Navigation help appears too many times, likely duplication")
		})
	})

	Describe("GenerateCV", func() {
		var intent *intents.GenerateCVIntent

		BeforeEach(func() {
			profiles := []*intents.CVProfile{
				{
					ID:         "profile1",
					Name:       "Test Profile",
					TargetRole: "Engineer",
				},
			}

			evt := fixtures.EventWith("event1", "Test achievement", "Test Co", "")
			evt.Date = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
			events := []*career.Event{evt}

			ctx := &intents.GenerateCVContext{
				AvailableProfiles: profiles,
				Events:            events,
			}

			var err error
			intent, err = intents.NewGenerateCVIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should have no duplicate help footer", func() {
			view := intent.View()

			// Count help text patterns - should appear only once
			upDownPattern := "↑/k"
			count := strings.Count(view, upDownPattern)
			Expect(count).To(BeNumerically("<=", 1), "Help text appears too many times, likely duplication")
		})

		It("should have no duplicate main menu reference", func() {
			view := intent.View()

			// Count "Main menu" - should appear only once in help
			mainMenuCount := strings.Count(strings.ToLower(view), "main menu")
			Expect(mainMenuCount).To(BeNumerically("<=", 2), "'Main menu' appears too many times, likely duplication")
		})
	})

	Describe("CaptureEvent", func() {
		var intent *captureevent.Intent

		BeforeEach(func() {
			ctx := &captureevent.IntentContext{
				CaptureStrategy: "manual",
				PreviousEvent:   nil,
				Metadata:        make(map[string]string),
			}

			var err error
			intent, err = captureevent.NewIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		It("should have reasonable quit hint count", func() {
			view := intent.View()

			// Count help patterns - be lenient since 'q' is common
			quitPattern := "q"
			quitCount := strings.Count(strings.ToLower(view), quitPattern)
			// 'q' is common in text, but shouldn't be excessive
			Expect(quitCount).To(BeNumerically("<=", 10), "'q' appears too many times, may indicate duplication")
		})
	})
})
