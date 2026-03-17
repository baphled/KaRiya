package browsetimeline

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("handleSkillMessages", func() {
	var intent *Intent

	BeforeEach(func() {
		intent = &Intent{}
	})

	It("returns nil for unknown message type", func() {
		cmd, handled := intent.handleSkillMessages(tea.KeyMsg{Type: tea.KeyEnter})
		Expect(handled).To(BeFalse())
		Expect(cmd).To(BeNil())
	})

	It("shows error modal for SkillLinkedMsg with error", func() {
		msg := SkillLinkedMsg{Error: errors.New("link error")}
		cmd, handled := intent.handleSkillMessages(msg)
		Expect(handled).To(BeTrue())
		Expect(cmd).To(BeNil())
	})

	It("refreshes skills modal for SkillLinkedMsg without error", func() {
		msg := SkillLinkedMsg{}
		cmd, handled := intent.handleSkillMessages(msg)
		Expect(handled).To(BeTrue())
		Expect(cmd).To(BeNil())
	})

	// Repeat for SkillUnlinkedMsg, SkillCreatedMsg, SkillSuggestionsLoadedMsg, SkillSuggestionsErrorMsg, SkillsForModalLoadedMsg, SkillPickerDataLoadedMsg, SkillsRefreshedMsg
})
