package event_test

import (
	"time"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReviewEnrichment", func() {
	var (
		eventData display.Event
		config    event.ReviewEnrichmentConfig
		view      *event.ReviewEnrichment
	)

	BeforeEach(func() {
		domainEvent := fixtures.EventWith("test-id", "Test Event for enrichment", "", "")
		domainEvent.Date = time.Now().Add(-time.Hour)
		domainEvent.Tags = []string{"project", "achievement"}
		domainEvent.Categories = []string{"technical", "leadership"}
		domainEvent.Skills = []string{"skill1"}
		eventData = display.EventFromDomain(domainEvent)
		config = event.ReviewEnrichmentConfig{
			AvailableTags:       []string{"tag1", "tag2", "tag3"},
			AvailableCategories: []string{"cat1", "cat2"},
			AvailableSkills:     []string{"skill1", "skill2"},
			TerminalWidth:       80,
			TerminalHeight:      24,
		}
		view = event.NewReviewEnrichment(eventData, config)
	})

	It("should create model with event data", func() {
		Expect(view.GetEvent()).To(Equal(eventData))
	})

	It("should return non-nil cmd from Init", func() {
		cmd := view.Init()
		Expect(cmd).NotTo(BeNil())
	})

	It("should render non-empty content", func() {
		content := view.RenderContent()
		Expect(content).NotTo(BeEmpty())
	})

	It("should return non-empty help text", func() {
		help := view.HelpText()
		Expect(help).NotTo(BeEmpty())
	})

	It("should emit CancelViewResult on Esc key", func() {
		msg := tea.KeyMsg{Type: tea.KeyEsc}
		_, result := view.Update(msg)
		_, ok := result.(*widgets.CancelViewResult)
		Expect(ok).To(BeTrue())
	})

	It("should emit CancelViewResult on Ctrl+C", func() {
		msg := tea.KeyMsg{Type: tea.KeyCtrlC}
		cmd, result := view.Update(msg)
		_, ok := result.(*widgets.CancelViewResult)
		Expect(ok).To(BeTrue())
		Expect(cmd).To(BeNil())
		Expect(view.IsCancelled()).To(BeTrue())
		Expect(view.IsSubmitted()).To(BeFalse())
	})

	It("should emit SubmitViewResult on form submit", func() {
		msg := tea.KeyMsg{Type: tea.KeyCtrlS}
		_, result := view.Update(msg)
		_, ok := result.(*widgets.SubmitViewResult)
		Expect(ok).To(BeTrue())
	})

	It("should submit from a zero-value view on Ctrl+S", func() {
		zeroView := &event.ReviewEnrichment{}

		cmd, result := zeroView.Update(tea.KeyMsg{Type: tea.KeyCtrlS})

		submitResult, ok := result.(*widgets.SubmitViewResult)
		Expect(ok).To(BeTrue())
		Expect(cmd).To(BeNil())
		Expect(zeroView.IsSubmitted()).To(BeTrue())
		Expect(zeroView.IsCancelled()).To(BeFalse())
		formData, ok := submitResult.FormData.(event.ReviewEnrichmentResult)
		Expect(ok).To(BeTrue())
		Expect(formData.Event).NotTo(BeZero())
		Expect(formData.FormData).To(BeNil())
		Expect(zeroView.GetError()).NotTo(HaveOccurred())
	})

	It("should emit CancelViewResult on form cancel", func() {
		msg := tea.KeyMsg{Type: tea.KeyEsc}
		_, result := view.Update(msg)
		_, ok := result.(*widgets.CancelViewResult)
		Expect(ok).To(BeTrue())
	})

	It("should expose edited event data", func() {
		input := view.GetEvent()
		Expect(input).NotTo(BeZero())
	})

	It("should track IsSubmitted and IsCancelled state", func() {
		Expect(view.IsSubmitted()).To(BeFalse())
		Expect(view.IsCancelled()).To(BeFalse())
		view.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
		Expect(view.IsSubmitted()).To(BeTrue())
		view.Update(tea.KeyMsg{Type: tea.KeyEsc})
		Expect(view.IsCancelled()).To(BeTrue())
	})

	It("should revert event to original", func() {
		original := view.GetEvent()
		view.Update(tea.KeyMsg{Type: tea.KeyEnter})
		view.Revert()
		Expect(view.GetEvent()).To(Equal(original))
	})

	It("should return nil error initially", func() {
		Expect(view.GetError()).To(Succeed())
	})

	It("should return strings for GetTitle, GetContent, GetFooter", func() {
		Expect(view.GetTitle()).NotTo(BeEmpty())
		Expect(view.GetContent()).NotTo(BeEmpty())
		Expect(view.GetFooter()).NotTo(BeEmpty())
	})

	It("should update dimensions on window size message", func() {
		msg := tea.WindowSizeMsg{Width: 100, Height: 40}
		view.Update(msg)
		Expect(view.GetEvent()).To(Equal(eventData))
	})

	It("should handle default width and height when config has zeros", func() {
		zeroConfig := event.ReviewEnrichmentConfig{
			AvailableTags:       []string{"tag1"},
			AvailableCategories: []string{"cat1"},
			AvailableSkills:     []string{"skill1"},
			TerminalWidth:       0,
			TerminalHeight:      0,
		}
		zeroView := event.NewReviewEnrichment(eventData, zeroConfig)
		Expect(zeroView).NotTo(BeNil())
		content := zeroView.RenderContent()
		Expect(content).NotTo(BeEmpty())
	})

	It("should handle regular key messages via form update", func() {
		msg := tea.KeyMsg{Type: tea.KeyTab}
		cmd, result := view.Update(msg)
		Expect(result).To(BeNil())
		_ = cmd
	})

	It("should handle enter key via form update", func() {
		msg := tea.KeyMsg{Type: tea.KeyEnter}
		cmd, result := view.Update(msg)
		_ = cmd
		_ = result
	})

	It("should handle down key via form update", func() {
		msg := tea.KeyMsg{Type: tea.KeyDown}
		cmd, result := view.Update(msg)
		_ = cmd
		_ = result
	})

	It("should handle shift+tab key via form update", func() {
		msg := tea.KeyMsg{Type: tea.KeyShiftTab}
		cmd, result := view.Update(msg)
		_ = cmd
		_ = result
	})

	It("should handle rune key via form update", func() {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
		cmd, result := view.Update(msg)
		_ = cmd
		_ = result
	})

	It("should submit after advancing through the form with enter", func() {
		runCmd := func(cmd tea.Cmd) tea.Msg {
			if cmd == nil {
				return nil
			}

			messages := make(chan tea.Msg, 1)
			go func() {
				messages <- cmd()
			}()

			select {
			case msg := <-messages:
				return msg
			case <-time.After(5 * time.Millisecond):
				return nil
			}
		}

		var result widgets.ViewResult
		for attempts := 0; attempts < 8 && result == nil; attempts++ {
			cmd, updateResult := view.Update(tea.KeyMsg{Type: tea.KeyEnter})
			if updateResult != nil {
				result = updateResult
				break
			}

			for steps := 0; steps < 4 && cmd != nil && result == nil; steps++ {
				msg := runCmd(cmd)
				switch msg.(type) {
				case nil, cursor.BlinkMsg:
					cmd = nil
				default:
					cmd, result = view.Update(msg)
				}
			}
		}

		Expect(result).NotTo(BeNil())
		Expect(result.Type()).To(Equal(widgets.ResultSubmit))
		Expect(view.IsSubmitted()).To(BeTrue())
		Expect(view.IsCancelled()).To(BeFalse())
	})
})
