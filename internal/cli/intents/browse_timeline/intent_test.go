package browse_timeline

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Intent", func() {
	Describe("Messages", func() {
		It("should have EventSelectedMsg defined", func() {
			msg := EventSelectedMsg{}
			Expect(msg).ToNot(BeNil())
		})

		It("should have FilterChangedMsg defined", func() {
			msg := FilterChangedMsg{}
			Expect(msg).ToNot(BeNil())
		})

		It("should have EventDeletedMsg defined", func() {
			msg := EventDeletedMsg{}
			Expect(msg).ToNot(BeNil())
		})
	})

	Describe("Intent Initialization", func() {
		It("should initialize without error", func() {
			ctx := &IntentContext{
				Events:          make([]*career.CareerEvent, 0),
				InitialFilters:  nil,
				SelectedEventID: "",
			}
			intent, err := NewIntent(ctx)
			Expect(err).To(BeNil())
			Expect(intent).ToNot(BeNil())
		})
	})
})
