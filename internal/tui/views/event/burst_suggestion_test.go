package event_test

import (
	"github.com/baphled/kariya/internal/domain/career"
	burstfact "github.com/baphled/kariya/internal/service/career/burstfact"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/display"
	widgets "github.com/baphled/kariya/internal/ui/uikit/widgets"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("BurstSuggestion", func() {
	var (
		suggestions []burstfact.BurstSuggestion
		view        *event.BurstSuggestion
	)

	BeforeEach(func() {
		suggestions = []burstfact.BurstSuggestion{
			{
				EventIDs:        []string{"ev1"},
				ConfidenceScore: 0.8,
				Name:            "Burst One",
				Description:     "Desc One",
			},
			{
				EventIDs:        []string{"ev2"},
				ConfidenceScore: 0.6,
				Name:            "Burst Two",
				Description:     "Desc Two",
			},
		}
		view = event.NewBurstSuggestion(display.BurstSuggestionsFromDomain(suggestions))
	})

	It("renders no-suggestions message when empty", func() {
		v := event.NewBurstSuggestion(display.BurstSuggestionsFromDomain([]burstfact.BurstSuggestion{}))
		Expect(v.RenderContent()).To(ContainSubstring("No burst suggestions available"))
	})

	It("renders content and help text", func() {
		Expect(view.RenderContent()).NotTo(BeEmpty())
		Expect(view.HelpText()).NotTo(BeEmpty())
	})

	It("navigates up/down with j/k and arrow keys", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyDown})
		Expect(view.CurrentIdx()).To(Equal(1))
		view.Update(tea.KeyMsg{Type: tea.KeyUp})
		Expect(view.CurrentIdx()).To(Equal(0))
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
		Expect(view.CurrentIdx()).To(Equal(1))
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
		Expect(view.CurrentIdx()).To(Equal(0))
	})

	It("confirms suggestion with 'y' and emits SubmitViewResult", func() {
		_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		Expect(result).To(BeAssignableToTypeOf(&widgets.SubmitViewResult{}))
		confirmed := view.GetConfirmed()
		Expect(confirmed).To(HaveLen(1))
		Expect(confirmed[0].Name).To(Equal("Burst One"))
	})

	It("rejects suggestion with 'n' and emits NavigateViewResult", func() {
		_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
		Expect(result).To(BeAssignableToTypeOf(&widgets.NavigateViewResult{}))
		rejected := view.GetRejected()
		Expect(rejected).To(HaveLen(1))
		Expect(rejected[0].Name).To(Equal("Burst One"))
	})

	It("cancels with Esc and emits CancelViewResult", func() {
		_, result := view.Update(tea.KeyMsg{Type: tea.KeyEsc})
		Expect(result).To(BeAssignableToTypeOf(&widgets.CancelViewResult{}))
	})

	It("toggles edit mode with 'e'", func() {
		Expect(view.IsEditing()).To(BeFalse())
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		Expect(view.IsEditing()).To(BeTrue())
	})

	It("SetRelatedEvents provides event data for rendering", func() {
		events := []*career.Event{fixtures.EventWith("ev1", "Event One", "", "")}
		view.SetRelatedEvents(0, display.EventsFromDomain(events))
		content := view.RenderContent()
		Expect(content).To(ContainSubstring("Event One"))
	})

	It("IsDone returns true when all suggestions processed", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		Expect(view.IsDone()).To(BeTrue())
	})

	It("Window size message updates dimensions", func() {
		msg := tea.WindowSizeMsg{Width: 80, Height: 24}
		view.Update(msg)
		Expect(view.Width()).To(Equal(80))
		Expect(view.Height()).To(Equal(24))
	})

	It("Init returns nil", func() {
		Expect(view.Init()).To(BeNil())
	})

	It("SetEditedName stores edited name for index", func() {
		view.SetEditedName(0, "New Name")
		content := view.RenderContent()
		Expect(content).NotTo(BeEmpty())
	})

	It("SetEditedDescription stores edited description for index", func() {
		view.SetEditedDescription(0, "New Desc")
		content := view.RenderContent()
		Expect(content).NotTo(BeEmpty())
	})

	It("ExitEditMode clears editing state", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		Expect(view.IsEditing()).To(BeTrue())
		view.ExitEditMode()
		Expect(view.IsEditing()).To(BeFalse())
	})

	It("renders edit view when in edit mode", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		Expect(view.IsEditing()).To(BeTrue())
		content := view.RenderContent()
		Expect(content).To(ContainSubstring("Edit"))
	})

	It("exits edit mode with 'e' key while editing", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		Expect(view.IsEditing()).To(BeTrue())
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		Expect(view.IsEditing()).To(BeFalse())
	})

	It("ignores non-edit keys while in edit mode", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
		_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		Expect(result).To(BeNil())
		Expect(view.GetConfirmed()).To(BeEmpty())
	})

	It("returns complete result when confirming last suggestion", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})
		Expect(result).To(BeAssignableToTypeOf(&widgets.SubmitViewResult{}))
		submitResult := result.(*widgets.SubmitViewResult)
		formData, ok := submitResult.FormData.(event.BurstSuggestionCompleteResult)
		Expect(ok).To(BeTrue())
		Expect(formData.Action).To(Equal("complete"))
	})

	It("returns complete result when rejecting last suggestion", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
		_, result := view.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
		Expect(result).To(BeAssignableToTypeOf(&widgets.SubmitViewResult{}))
		submitResult := result.(*widgets.SubmitViewResult)
		formData, ok := submitResult.FormData.(event.BurstSuggestionCompleteResult)
		Expect(ok).To(BeTrue())
		Expect(formData.Action).To(Equal("complete"))
	})

	It("does not navigate below last suggestion", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyDown})
		Expect(view.CurrentIdx()).To(Equal(1))
		view.Update(tea.KeyMsg{Type: tea.KeyDown})
		Expect(view.CurrentIdx()).To(Equal(1))
	})

	It("does not navigate above first suggestion", func() {
		view.Update(tea.KeyMsg{Type: tea.KeyUp})
		Expect(view.CurrentIdx()).To(Equal(0))
	})

	It("renders related events when empty slice is set", func() {
		view.SetRelatedEvents(0, display.EventsFromDomain([]*career.Event{}))
		content := view.RenderContent()
		Expect(content).NotTo(ContainSubstring("Related Events"))
	})
})
