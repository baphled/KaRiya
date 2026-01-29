package fixtures_test

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EventFactory", func() {
	BeforeEach(func() {
		// Set seed for reproducible tests
		fixtures.SetSeed(42)
	})

	Describe("MustCreate", func() {
		It("should create an event with all required fields", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			Expect(event).NotTo(BeNil())
			Expect(event.ID).NotTo(BeEmpty())
			Expect(event.Text).NotTo(BeEmpty())
			Expect(event.Date).NotTo(BeZero())
			Expect(event.CreatedAt).NotTo(BeZero())
			Expect(event.UpdatedAt).NotTo(BeZero())
		})

		It("should create events with unique sequential IDs", func() {
			event1 := fixtures.EventFactory.MustCreate().(*career.Event)
			event2 := fixtures.EventFactory.MustCreate().(*career.Event)
			Expect(event1.ID).NotTo(Equal(event2.ID))
		})

		It("should create events with company and project", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			Expect(event.Company).NotTo(BeEmpty())
			Expect(event.Project).NotTo(BeEmpty())
		})

		It("should create events with tags and categories", func() {
			event := fixtures.EventFactory.MustCreate().(*career.Event)
			Expect(event.Tags).NotTo(BeEmpty())
			Expect(event.Categories).NotTo(BeEmpty())
		})
	})

	Describe("MustCreateWithOption", func() {
		It("should allow overriding specific fields", func() {
			event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
				"Company": "CustomCorp",
				"Project": "CustomProject",
			}).(*career.Event)
			Expect(event.Company).To(Equal("CustomCorp"))
			Expect(event.Project).To(Equal("CustomProject"))
		})
	})
})

var _ = Describe("Event Quick Helper", func() {
	Describe("Event", func() {
		It("should create a minimal valid event with given ID", func() {
			event := fixtures.Event("test-id")
			Expect(event.ID).To(Equal("test-id"))
			Expect(event.Text).To(ContainSubstring("test-id"))
			Expect(event.Date).NotTo(BeZero())
			Expect(event.CreatedAt).NotTo(BeZero())
			Expect(event.UpdatedAt).NotTo(BeZero())
		})
	})

	Describe("EventWith", func() {
		It("should create an event with custom fields", func() {
			event := fixtures.EventWith("evt-1", "Led team", "TechCorp", "Platform")
			Expect(event.ID).To(Equal("evt-1"))
			Expect(event.Text).To(Equal("Led team"))
			Expect(event.Company).To(Equal("TechCorp"))
			Expect(event.Project).To(Equal("Platform"))
		})
	})
})

var _ = Describe("Events Batch Helper", func() {
	BeforeEach(func() {
		fixtures.SetSeed(42)
	})

	Describe("Events", func() {
		It("should create n events with unique IDs", func() {
			events := fixtures.Events(5)
			Expect(events).To(HaveLen(5))

			ids := make(map[string]bool)
			for _, e := range events {
				Expect(ids[e.ID]).To(BeFalse(), "Duplicate ID found: "+e.ID)
				ids[e.ID] = true
			}
		})
	})
})
