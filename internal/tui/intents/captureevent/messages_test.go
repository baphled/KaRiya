package captureevent_test

import (
	"errors"

	ce "github.com/baphled/kariya/internal/tui/intents/captureevent"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Messages", func() {
	Describe("SubmitCompleteMsg", func() {
		It("should be constructible as a zero-value struct", func() {
			msg := ce.SubmitCompleteMsg{}
			Expect(msg).To(Equal(ce.SubmitCompleteMsg{}))
		})
	})

	Describe("SubmitErrorMsg", func() {
		It("should carry code, message and cause", func() {
			cause := errors.New("db connection failed")
			msg := ce.SubmitErrorMsg{
				Code:    "PERSISTENCE_ERROR",
				Message: "Failed to save event",
				Cause:   cause,
			}
			Expect(msg.Code).To(Equal("PERSISTENCE_ERROR"))
			Expect(msg.Message).To(Equal("Failed to save event"))
			Expect(msg.Cause).To(MatchError("db connection failed"))
		})

		It("should allow nil cause", func() {
			msg := ce.SubmitErrorMsg{
				Code:    "VALIDATION_ERROR",
				Message: "Invalid event",
			}
			Expect(msg.Cause).ToNot(HaveOccurred())
		})
	})

	Describe("InferenceCompleteMsg", func() {
		It("should be constructible as a zero-value struct", func() {
			msg := ce.InferenceCompleteMsg{}
			Expect(msg).To(Equal(ce.InferenceCompleteMsg{}))
		})
	})

	Describe("DismissModalMsg", func() {
		It("should be constructible as a zero-value struct", func() {
			msg := ce.DismissModalMsg{}
			Expect(msg).To(Equal(ce.DismissModalMsg{}))
		})
	})

	Describe("Type distinctiveness", func() {
		It("should be distinct types for tea.Msg dispatch", func() {
			var msgs []interface{}
			msgs = append(msgs, ce.SubmitCompleteMsg{})
			msgs = append(msgs, ce.SubmitErrorMsg{})
			msgs = append(msgs, ce.DismissModalMsg{})
			msgs = append(msgs, ce.InferenceCompleteMsg{})

			_, isComplete := msgs[0].(ce.SubmitCompleteMsg)
			_, isError := msgs[0].(ce.SubmitErrorMsg)
			Expect(isComplete).To(BeTrue())
			Expect(isError).To(BeFalse())

			_, isComplete2 := msgs[1].(ce.SubmitCompleteMsg)
			_, isError2 := msgs[1].(ce.SubmitErrorMsg)
			Expect(isComplete2).To(BeFalse())
			Expect(isError2).To(BeTrue())

			_, isInference := msgs[3].(ce.InferenceCompleteMsg)
			_, isSubmit := msgs[3].(ce.SubmitCompleteMsg)
			Expect(isInference).To(BeTrue())
			Expect(isSubmit).To(BeFalse())
		})
	})
})
