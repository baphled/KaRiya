package career

import (
	"time"

	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Event Model", func() {
	Describe("TableName", func() {
		It("returns career_events", func() {
			Expect(Event{}.TableName()).To(Equal("career_events"))
		})
	})

	Describe("ToDomain", func() {
		It("converts all fields to domain Event", func() {
			now := time.Now()
			model := &Event{
				ID:         "evt-1",
				Text:       "Led platform migration",
				Date:       now,
				Company:    "TechCorp",
				Project:    "Platform",
				Tags:       StringSlice{"technical", "leadership"},
				Categories: StringSlice{"technical"},
				CreatedAt:  now,
				UpdatedAt:  now,
				Skills: []Skill{
					{ID: "sk-1"},
					{ID: "sk-2"},
				},
			}

			domain := model.ToDomain()

			Expect(domain.ID).To(Equal("evt-1"))
			Expect(domain.Text).To(Equal("Led platform migration"))
			Expect(domain.Date).To(Equal(now))
			Expect(domain.Company).To(Equal("TechCorp"))
			Expect(domain.Project).To(Equal("Platform"))
			Expect(domain.Tags).To(Equal([]string{"technical", "leadership"}))
			Expect(domain.Categories).To(Equal([]string{"technical"}))
			Expect(domain.Skills).To(Equal([]string{"sk-1", "sk-2"}))
			Expect(domain.CreatedAt).To(Equal(now))
			Expect(domain.UpdatedAt).To(Equal(now))
		})

		It("handles nil tags and categories", func() {
			model := &Event{
				ID:        "evt-2",
				Text:      "Simple event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			domain := model.ToDomain()

			Expect(domain.Tags).To(BeNil())
			Expect(domain.Categories).To(BeNil())
		})

		It("handles empty skills slice", func() {
			model := &Event{
				ID:        "evt-3",
				Text:      "No skills event",
				Date:      time.Now(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Skills:    []Skill{},
			}

			domain := model.ToDomain()

			Expect(domain.Skills).To(BeEmpty())
		})
	})

	Describe("EventFromDomain", func() {
		It("converts all fields from domain Event", func() {
			domainEvent := fixtures.EventWith("evt-d1", "Architected microservices", "StartupXYZ", "Arch Refresh")
			domainEvent.Tags = []string{"achievement", "technical"}
			domainEvent.Categories = []string{"technical", "leadership"}

			model := EventFromDomain(domainEvent)

			Expect(model.ID).To(Equal("evt-d1"))
			Expect(model.Text).To(Equal("Architected microservices"))
			Expect(model.Company).To(Equal("StartupXYZ"))
			Expect(model.Project).To(Equal("Arch Refresh"))
			Expect([]string(model.Tags)).To(Equal([]string{"achievement", "technical"}))
			Expect([]string(model.Categories)).To(Equal([]string{"technical", "leadership"}))
			Expect(model.CreatedAt).To(Equal(domainEvent.CreatedAt))
			Expect(model.UpdatedAt).To(Equal(domainEvent.UpdatedAt))
		})

		It("handles nil tags and categories", func() {
			domainEvent := fixtures.EventWith("evt-d2", "Minimal event", "", "")
			domainEvent.Tags = nil
			domainEvent.Categories = nil

			model := EventFromDomain(domainEvent)

			Expect(model.Tags).To(BeNil())
			Expect(model.Categories).To(BeNil())
		})
	})

	Describe("ToDomain-FromDomain round trip", func() {
		It("preserves all Event data except Skills", func() {
			original := fixtures.EventWith("evt-rt", "Round trip event", "RoundTripCo", "RT Project")
			original.Tags = []string{"technical", "project"}
			original.Categories = []string{"technical"}
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)
			original.Date = original.Date.Truncate(time.Second)

			model := EventFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.ID).To(Equal(original.ID))
			Expect(restored.Text).To(Equal(original.Text))
			Expect(restored.Date).To(Equal(original.Date))
			Expect(restored.Company).To(Equal(original.Company))
			Expect(restored.Project).To(Equal(original.Project))
			Expect(restored.Tags).To(Equal(original.Tags))
			Expect(restored.Categories).To(Equal(original.Categories))
			Expect(restored.CreatedAt).To(Equal(original.CreatedAt))
			Expect(restored.UpdatedAt).To(Equal(original.UpdatedAt))
		})
	})
})
