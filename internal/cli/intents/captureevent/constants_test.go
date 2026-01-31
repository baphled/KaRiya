package captureevent_test

import (
	ce "github.com/baphled/kariya/internal/cli/intents/captureevent"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Constants", func() {
	Describe("State", func() {
		It("should define all workflow states", func() {
			Expect(string(ce.StateChooseStrategy)).To(Equal("choose_strategy"))
			Expect(string(ce.StateForm)).To(Equal("form"))
			Expect(string(ce.StateReview)).To(Equal("review"))
			Expect(string(ce.StateSubmit)).To(Equal("submit"))
		})

		It("should have distinct values for each state", func() {
			states := []ce.State{
				ce.StateChooseStrategy,
				ce.StateForm,
				ce.StateReview,
				ce.StateSubmit,
			}
			seen := make(map[ce.State]bool)
			for _, s := range states {
				Expect(seen[s]).To(BeFalse(), "duplicate state: %s", s)
				seen[s] = true
			}
		})
	})

	Describe("EditingMode", func() {
		It("should define all editing modes", func() {
			Expect(string(ce.EditingModeNone)).To(Equal(""))
			Expect(string(ce.EditingModeMetadata)).To(Equal("metadata"))
			Expect(string(ce.EditingModeBursts)).To(Equal("bursts"))
			Expect(string(ce.EditingModeFacts)).To(Equal("facts"))
		})

		It("should have distinct values for each mode", func() {
			modes := []ce.EditingMode{
				ce.EditingModeMetadata,
				ce.EditingModeBursts,
				ce.EditingModeFacts,
			}
			seen := make(map[ce.EditingMode]bool)
			for _, m := range modes {
				Expect(seen[m]).To(BeFalse(), "duplicate editing mode: %s", m)
				seen[m] = true
			}
		})

		It("should use empty string for EditingModeNone", func() {
			Expect(ce.EditingModeNone).To(Equal(ce.EditingMode("")))
		})
	})

	Describe("CaptureStrategy re-exports", func() {
		It("should re-export StrategyQuick and StrategyManual", func() {
			Expect(string(ce.StrategyQuick)).To(Equal("quick"))
			Expect(string(ce.StrategyManual)).To(Equal("manual"))
		})
	})
})
