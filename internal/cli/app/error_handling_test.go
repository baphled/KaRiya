package app

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/cli/service"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("Error Recovery & Edge Cases", func() {
	var (
		repo       *careerrepo.MemoryRepository
		svc        *careerservice.Service
		cliSvc     *service.CLIEventService
		model      *Model
		ctx        context.Context
	)

	BeforeEach(func() {
		repo = careerrepo.NewMemoryRepository()
		svc = careerservice.NewService(repo)
		cliSvc = service.NewCLIEventService(svc)
		model = NewModel(cliSvc, svc)
		ctx = context.Background()
	})

	Context("Service Layer Error Handling", func() {
		It("should handle service validation errors gracefully", func() {
			// Try to capture event with empty text
			err := cliSvc.CaptureEvent(ctx, "", time.Now(), careerservice.TimelineJournaling)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("text cannot be empty"))
		})

		It("should handle timeline mode 30-day window validation", func() {
			// Try to capture event from more than 30 days ago in timeline mode
			oldDate := time.Now().AddDate(0, 0, -31)
			err := cliSvc.CaptureEvent(ctx, "Old event", oldDate, careerservice.TimelineJournaling)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("timeline events must be recent"))
		})

		It("should handle future date rejection", func() {
			// Try to capture event with future date
			futureDate := time.Now().AddDate(0, 0, 1)
			err := cliSvc.CaptureEvent(ctx, "Future event", futureDate, careerservice.CVBackfill)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("cannot be in the future"))
		})

		It("should reject whitespace-only text", func() {
			err := cliSvc.CaptureEvent(ctx, "   ", time.Now(), careerservice.ManualEntry)
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Repository Error Recovery", func() {
		It("should handle list operation with empty repository", func() {
			events, err := cliSvc.ListEvents(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(0))
		})

		It("should handle retrieval of non-existent event", func() {
			event, err := cliSvc.GetEventByID(ctx, "non-existent-id")
			Expect(err).To(HaveOccurred())
			Expect(event).To(BeNil())
		})

		It("should accept valid default filter", func() {
			// Capture an event first
			_ = cliSvc.CaptureEvent(ctx, "Test event", time.Now(), careerservice.ManualEntry)

			// List with default/nil filters should work
			_, err := cliSvc.ListEvents(ctx, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Application State Recovery", func() {
		It("should maintain valid state after validation error", func() {
			// Try to capture invalid event
			err := cliSvc.CaptureEvent(ctx, "", time.Now(), careerservice.TimelineJournaling)
			Expect(err).To(HaveOccurred())

			// App should still be functional for next event
			err2 := cliSvc.CaptureEvent(ctx, "Valid event", time.Now(), careerservice.ManualEntry)
			Expect(err2).NotTo(HaveOccurred())
		})
	})

	Context("Long Event Text Edge Cases", func() {
		It("should handle very long text at boundary (1999 chars)", func() {
			longText := ""
			for i := 0; i < 1999; i++ {
				longText += "a"
			}
			err := cliSvc.CaptureEvent(ctx, longText, time.Now(), careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject text exceeding 2000 character limit", func() {
			longText := ""
			for i := 0; i < 2001; i++ {
				longText += "a"
			}
			err := cliSvc.CaptureEvent(ctx, longText, time.Now(), careerservice.ManualEntry)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("2000"))
		})
	})

	Context("Navigation State Recovery", func() {
		It("should handle screen navigation without errors", func() {
			model.currentScreen = CaptureScreen
			Expect(model.currentScreen).To(Equal(CaptureScreen))

			model.currentScreen = ListScreen
			Expect(model.currentScreen).To(Equal(ListScreen))

			model.currentScreen = HomeScreen
			Expect(model.currentScreen).To(Equal(HomeScreen))
		})

		It("should maintain previous screen history", func() {
			initialPrevious := model.previousScreen
			model.currentScreen = CaptureScreen
			model.previousScreen = HomeScreen
			Expect(model.previousScreen).To(Equal(HomeScreen))
			model.previousScreen = initialPrevious
		})
	})

	Context("Message Handling Robustness", func() {
		It("should handle unknown message types gracefully", func() {
			customMsg := struct{}{}
			newModel, cmd := model.Update(customMsg)
			// Should not panic, just return nil cmd
			Expect(newModel).NotTo(BeNil())
			Expect(cmd).To(BeNil())
		})

		It("should handle window size changes", func() {
			winMsg := tea.WindowSizeMsg{
				Width:  120,
				Height: 40,
			}
			newModel, _ := model.Update(winMsg)
			Expect(newModel).NotTo(BeNil())
		})

		It("should handle rapid successive updates without panic", func() {
			// Should not crash under rapid input
			for i := 0; i < 10; i++ {
				_, _ = model.Update(tea.KeyMsg{Type: tea.KeyEscape})
			}
			Expect(model).NotTo(BeNil())
		})
	})

	Context("Special Character Handling", func() {
		It("should accept text with newlines and tabs", func() {
			text := "Valid event\nwith newlines\tand tabs"
			err := cliSvc.CaptureEvent(ctx, text, time.Now(), careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept unicode characters", func() {
			text := "Event with unicode: 你好 🎉 العربية"
			err := cliSvc.CaptureEvent(ctx, text, time.Now(), careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept very long word without spaces", func() {
			// Long word (no spaces), should still be valid
			longWord := ""
			for i := 0; i < 100; i++ {
				longWord += "word"
			}
			Expect(len(longWord)).To(BeNumerically("<", 2000))
			err := cliSvc.CaptureEvent(ctx, longWord, time.Now(), careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("Date Edge Cases", func() {
		It("should accept event from 29 days ago in timeline mode", func() {
			// Recent events within 30 days should be acceptable
			twentyNineDaysAgo := time.Now().AddDate(0, 0, -29)
			err := cliSvc.CaptureEvent(ctx, "Event", twentyNineDaysAgo, careerservice.TimelineJournaling)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should reject event from 31 days ago in timeline mode", func() {
			thirtyOneDaysAgo := time.Now().AddDate(0, 0, -31)
			err := cliSvc.CaptureEvent(ctx, "Event", thirtyOneDaysAgo, careerservice.TimelineJournaling)
			Expect(err).To(HaveOccurred())
		})

		It("should accept very old dates in cv backfill mode", func() {
			// Year 2000
			y2k := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
			err := cliSvc.CaptureEvent(ctx, "Old event", y2k, careerservice.CVBackfill)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should accept very old dates in manual entry mode", func() {
			// Year 1990
			oldDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
			err := cliSvc.CaptureEvent(ctx, "Very old event", oldDate, careerservice.ManualEntry)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
