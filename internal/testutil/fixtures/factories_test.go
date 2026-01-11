package fixtures_test

import (
	"testing"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestFixtures(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Fixtures Suite")
}

var _ = Describe("EventFactory", func() {
	BeforeEach(func() {
		// Set seed for reproducible tests
		fixtures.SetSeed(42)
	})

	Describe("MustCreate", func() {
		It("should create an event with all required fields", func() {
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			Expect(event).NotTo(BeNil())
			Expect(event.ID).NotTo(BeEmpty())
			Expect(event.Text).NotTo(BeEmpty())
			Expect(event.Date).NotTo(BeZero())
			Expect(event.CreatedAt).NotTo(BeZero())
			Expect(event.UpdatedAt).NotTo(BeZero())
		})

		It("should create events with unique sequential IDs", func() {
			event1 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			event2 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			Expect(event1.ID).NotTo(Equal(event2.ID))
		})

		It("should create events with company and project", func() {
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			Expect(event.Company).NotTo(BeEmpty())
			Expect(event.Project).NotTo(BeEmpty())
		})

		It("should create events with tags and categories", func() {
			event := fixtures.EventFactory.MustCreate().(*career.CareerEvent)
			Expect(event.Tags).NotTo(BeEmpty())
			Expect(event.Categories).NotTo(BeEmpty())
		})
	})

	Describe("MustCreateWithOption", func() {
		It("should allow overriding specific fields", func() {
			event := fixtures.EventFactory.MustCreateWithOption(map[string]interface{}{
				"Company": "CustomCorp",
				"Project": "CustomProject",
			}).(*career.CareerEvent)
			Expect(event.Company).To(Equal("CustomCorp"))
			Expect(event.Project).To(Equal("CustomProject"))
		})
	})
})

var _ = Describe("BurstFactory", func() {
	BeforeEach(func() {
		fixtures.SetSeed(42)
	})

	Describe("MustCreate", func() {
		It("should create a burst with all required fields", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst).NotTo(BeNil())
			Expect(burst.ID).NotTo(BeEmpty())
			Expect(burst.Name).NotTo(BeEmpty())
			Expect(burst.EventIDs).To(HaveLen(2)) // Minimum required
			Expect(burst.CreatedAt).NotTo(BeZero())
			Expect(burst.UpdatedAt).NotTo(BeZero())
		})

		It("should create bursts with unique sequential IDs", func() {
			burst1 := fixtures.BurstFactory.MustCreate().(*career.Burst)
			burst2 := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst1.ID).NotTo(Equal(burst2.ID))
		})

		It("should create bursts with description", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst.Description).NotTo(BeEmpty())
		})

		It("should create unconfirmed bursts by default", func() {
			burst := fixtures.BurstFactory.MustCreate().(*career.Burst)
			Expect(burst.Confirmed).To(BeFalse())
			Expect(burst.ConfirmedAt).To(BeNil())
		})
	})

	Describe("MustCreateWithOption", func() {
		It("should allow overriding confirmed status", func() {
			burst := fixtures.BurstFactory.MustCreateWithOption(map[string]interface{}{
				"Confirmed": true,
			}).(*career.Burst)
			Expect(burst.Confirmed).To(BeTrue())
		})
	})
})

var _ = Describe("FactFactory", func() {
	BeforeEach(func() {
		fixtures.SetSeed(42)
	})

	Describe("MustCreate", func() {
		It("should create a fact with all required fields", func() {
			fact := fixtures.FactFactory.MustCreate().(*career.Fact)
			Expect(fact).NotTo(BeNil())
			Expect(fact.ID).NotTo(BeEmpty())
			Expect(fact.Text).NotTo(BeEmpty())
			Expect(fact.CompetencyCategories).NotTo(BeEmpty())
			Expect(fact.RoleFit).NotTo(BeEmpty())
			Expect(fact.AudienceRelevance).NotTo(BeEmpty())
			Expect(fact.StrengthSignal).NotTo(BeEmpty())
			Expect(fact.SourceEventID).NotTo(BeEmpty())
			Expect(fact.CreatedAt).NotTo(BeZero())
			Expect(fact.UpdatedAt).NotTo(BeZero())
		})

		It("should create facts with unique sequential IDs", func() {
			fact1 := fixtures.FactFactory.MustCreate().(*career.Fact)
			fact2 := fixtures.FactFactory.MustCreate().(*career.Fact)
			Expect(fact1.ID).NotTo(Equal(fact2.ID))
		})

		It("should create facts with valid role fit values", func() {
			fact := fixtures.FactFactory.MustCreate().(*career.Fact)
			validRoleFits := []career.RoleFit{
				career.RoleFitStaff,
				career.RoleFitSeniorIC,
				career.RoleFitPrincipal,
				career.RoleFitEM,
			}
			Expect(validRoleFits).To(ContainElement(fact.RoleFit))
		})
	})
})

var _ = Describe("Quick Helper Functions", func() {
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

	Describe("Burst", func() {
		It("should create a minimal valid burst with given ID and event IDs", func() {
			burst := fixtures.Burst("burst-1", "evt-1", "evt-2")
			Expect(burst.ID).To(Equal("burst-1"))
			Expect(burst.Name).To(ContainSubstring("burst-1"))
			Expect(burst.EventIDs).To(Equal([]string{"evt-1", "evt-2"}))
			Expect(burst.Confirmed).To(BeFalse())
		})

		It("should use default event IDs if less than 2 provided", func() {
			burst := fixtures.Burst("burst-1", "evt-1")
			Expect(burst.EventIDs).To(HaveLen(2))
		})
	})

	Describe("BurstConfirmed", func() {
		It("should create a confirmed burst", func() {
			burst := fixtures.BurstConfirmed("burst-1", "evt-1", "evt-2")
			Expect(burst.Confirmed).To(BeTrue())
			Expect(burst.ConfirmedAt).NotTo(BeNil())
		})
	})

	Describe("Fact", func() {
		It("should create a minimal valid fact with given ID and source", func() {
			fact := fixtures.Fact("fact-1", "evt-1")
			Expect(fact.ID).To(Equal("fact-1"))
			Expect(fact.SourceEventID).To(Equal("evt-1"))
			Expect(fact.Text).To(ContainSubstring("fact-1"))
			Expect(fact.CompetencyCategories).To(Equal([]string{"technical"}))
			Expect(fact.RoleFit).To(Equal(career.RoleFitStaff))
			Expect(fact.AudienceRelevance).To(Equal([]string{"hiring_manager"}))
			Expect(fact.StrengthSignal).To(Equal("high"))
		})
	})

	Describe("FactFromBurst", func() {
		It("should create a fact linked to a burst", func() {
			fact := fixtures.FactFromBurst("fact-1", "burst-1")
			Expect(fact.SourceBurstID).To(Equal("burst-1"))
			Expect(fact.SourceEventID).To(BeEmpty())
		})
	})
})

var _ = Describe("Batch Creation Helpers", func() {
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

	Describe("Bursts", func() {
		It("should create n bursts", func() {
			events := fixtures.Events(6)
			bursts := fixtures.Bursts(3, events)
			Expect(bursts).To(HaveLen(3))
		})

		It("should link bursts to provided events", func() {
			events := fixtures.Events(4)
			bursts := fixtures.Bursts(2, events)
			Expect(bursts[0].EventIDs).To(HaveLen(2))
			Expect(events[0].ID).To(Equal(bursts[0].EventIDs[0]))
		})
	})

	Describe("Facts", func() {
		It("should create n facts", func() {
			events := fixtures.Events(3)
			facts := fixtures.Facts(5, events)
			Expect(facts).To(HaveLen(5))
		})

		It("should link facts to provided events", func() {
			events := fixtures.Events(2)
			facts := fixtures.Facts(4, events)
			// Facts should cycle through available events
			Expect(facts[0].SourceEventID).To(Equal(events[0].ID))
			Expect(facts[1].SourceEventID).To(Equal(events[1].ID))
			Expect(facts[2].SourceEventID).To(Equal(events[0].ID))
		})
	})
})

var _ = Describe("Reproducibility", func() {
	It("should produce same data with same seed", func() {
		fixtures.SetSeed(12345)
		event1 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)

		fixtures.SetSeed(12345)
		event2 := fixtures.EventFactory.MustCreate().(*career.CareerEvent)

		// With same seed, random parts should match
		// Note: IDs are sequential so they reset with new factory
		Expect(event1.Company).To(Equal(event2.Company))
		Expect(event1.Project).To(Equal(event2.Project))
	})
})
