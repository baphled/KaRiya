package intents

import (
	"strings"
	"time"

	careerdomain "github.com/baphled/kariya/internal/domain/career"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Export Format Functions", func() {
	Describe("marshalToJSON", func() {
		It("should marshal simple data to JSON", func() {
			data := map[string]string{"key": "value"}
			result, err := marshalToJSON(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(`"key": "value"`))
		})

		It("should marshal struct to JSON with indentation", func() {
			type TestStruct struct {
				Name string `json:"name"`
				Age  int    `json:"age"`
			}
			data := TestStruct{Name: "John", Age: 30}
			result, err := marshalToJSON(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(`"name": "John"`))
			Expect(result).To(ContainSubstring(`"age": 30`))
		})

		It("should handle empty data", func() {
			data := map[string]string{}
			result, err := marshalToJSON(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("{}"))
		})
	})

	Describe("marshalToYAML", func() {
		It("should marshal simple data to YAML", func() {
			data := map[string]string{"key": "value"}
			result, err := marshalToYAML(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("key: value"))
		})

		It("should marshal struct to YAML", func() {
			type TestStruct struct {
				Name string `yaml:"name"`
				Age  int    `yaml:"age"`
			}
			data := TestStruct{Name: "John", Age: 30}
			result, err := marshalToYAML(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("name: John"))
			Expect(result).To(ContainSubstring("age: 30"))
		})

		It("should handle empty data", func() {
			data := map[string]string{}
			result, err := marshalToYAML(data)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("{}\n"))
		})
	})

	Describe("marshalEventsToCSV", func() {
		var events []*careerdomain.CareerEvent

		BeforeEach(func() {
			events = []*careerdomain.CareerEvent{
				{
					ID:         "event-1",
					Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					Text:       "Led team of 5 engineers",
					Company:    "TechCorp",
					Project:    "Platform",
					Tags:       []string{"leadership", "engineering"},
					Categories: []string{"technical", "management"},
				},
			}
		})

		It("should include CSV header row", func() {
			result, err := marshalEventsToCSV(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HavePrefix("id,date,text,company,project,tags,categories\n"))
		})

		It("should marshal event fields correctly", func() {
			result, err := marshalEventsToCSV(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("event-1"))
			Expect(result).To(ContainSubstring("2024-01-15"))
			Expect(result).To(ContainSubstring("TechCorp"))
			Expect(result).To(ContainSubstring("Platform"))
		})

		It("should join tags with semicolon", func() {
			result, err := marshalEventsToCSV(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("leadership;engineering"))
		})

		It("should join categories with semicolon", func() {
			result, err := marshalEventsToCSV(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("technical;management"))
		})

		It("should handle multiple events", func() {
			events = append(events, &careerdomain.CareerEvent{
				ID:   "event-2",
				Date: time.Date(2024, 2, 20, 0, 0, 0, 0, time.UTC),
				Text: "Second event",
			})
			result, err := marshalEventsToCSV(events)
			Expect(err).NotTo(HaveOccurred())
			lines := strings.Split(result, "\n")
			Expect(lines).To(HaveLen(4)) // header + 2 events + empty trailing
		})

		It("should handle empty events slice", func() {
			result, err := marshalEventsToCSV([]*careerdomain.CareerEvent{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("id,date,text,company,project,tags,categories\n"))
		})
	})

	Describe("marshalEventsToText", func() {
		var events []*careerdomain.CareerEvent

		BeforeEach(func() {
			events = []*careerdomain.CareerEvent{
				{
					ID:         "event-1",
					Date:       time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					Text:       "Led team of 5 engineers",
					Company:    "TechCorp",
					Project:    "Platform",
					Tags:       []string{"leadership", "engineering"},
					Categories: []string{"technical", "management"},
				},
			}
		})

		It("should include header", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Career Events"))
		})

		It("should include separator line", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(strings.Repeat("=", 80)))
		})

		It("should include event number", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Event 1"))
		})

		It("should include date in correct format", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Date: 2024-01-15"))
		})

		It("should include company when present", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Company: TechCorp"))
		})

		It("should include project when present", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Project: Platform"))
		})

		It("should include text", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Text: Led team of 5 engineers"))
		})

		It("should include tags when present", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Tags: leadership, engineering"))
		})

		It("should include categories when present", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Categories: technical, management"))
		})

		It("should include total count", func() {
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Total Events: 1"))
		})

		It("should omit company when empty", func() {
			events[0].Company = ""
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(ContainSubstring("Company:"))
		})

		It("should omit project when empty", func() {
			events[0].Project = ""
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(ContainSubstring("Project:"))
		})

		It("should omit tags when empty", func() {
			events[0].Tags = nil
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(ContainSubstring("Tags:"))
		})

		It("should omit categories when empty", func() {
			events[0].Categories = nil
			result, err := marshalEventsToText(events)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(ContainSubstring("Categories:"))
		})

		It("should handle empty events slice", func() {
			result, err := marshalEventsToText([]*careerdomain.CareerEvent{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Career Events"))
			Expect(result).To(ContainSubstring("Total Events: 0"))
		})
	})

	Describe("marshalFactsToCSV", func() {
		var facts []*careerdomain.Fact

		BeforeEach(func() {
			facts = []*careerdomain.Fact{
				{
					ID:                   "fact-1",
					Text:                 "Led migration to microservices",
					CompetencyCategories: []string{"architecture", "leadership"},
					RoleFit:              careerdomain.RoleFitStaff,
					AudienceRelevance:    []string{"hiring_manager", "peer"},
					StrengthSignal:       "strong",
					SourceEventID:        "event-1",
					SourceBurstID:        "",
				},
			}
		})

		It("should include CSV header row", func() {
			result, err := marshalFactsToCSV(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HavePrefix("id,text,competency_categories,role_fit,audience_relevance,strength_signal,source_event_id,source_burst_id\n"))
		})

		It("should marshal fact fields correctly", func() {
			result, err := marshalFactsToCSV(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("fact-1"))
			Expect(result).To(ContainSubstring("Led migration to microservices"))
			Expect(result).To(ContainSubstring("staff"))
			Expect(result).To(ContainSubstring("strong"))
		})

		It("should join competency categories with semicolon", func() {
			result, err := marshalFactsToCSV(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("architecture;leadership"))
		})

		It("should join audience relevance with semicolon", func() {
			result, err := marshalFactsToCSV(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("hiring_manager;peer"))
		})

		It("should handle empty facts slice", func() {
			result, err := marshalFactsToCSV([]*careerdomain.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("id,text,competency_categories,role_fit,audience_relevance,strength_signal,source_event_id,source_burst_id\n"))
		})
	})

	Describe("marshalFactsToText", func() {
		var facts []*careerdomain.Fact

		BeforeEach(func() {
			facts = []*careerdomain.Fact{
				{
					ID:                   "fact-1",
					Text:                 "Led migration to microservices",
					CompetencyCategories: []string{"architecture", "leadership"},
					RoleFit:              careerdomain.RoleFitStaff,
					AudienceRelevance:    []string{"hiring_manager", "peer"},
					StrengthSignal:       "strong",
					SourceEventID:        "event-1",
					SourceBurstID:        "",
				},
			}
		})

		It("should include header", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Facts"))
		})

		It("should include separator line", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(strings.Repeat("=", 80)))
		})

		It("should include fact number and text", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Fact 1: Led migration to microservices"))
		})

		It("should include competency categories", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Competency Categories: architecture, leadership"))
		})

		It("should include role fit", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Role Fit: staff"))
		})

		It("should include audience relevance", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Audience Relevance: hiring_manager, peer"))
		})

		It("should include strength signal", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Strength Signal: strong"))
		})

		It("should include source event when present", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Source Event: event-1"))
		})

		It("should include source burst when present", func() {
			facts[0].SourceBurstID = "burst-1"
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Source Burst: burst-1"))
		})

		It("should omit source event when empty", func() {
			facts[0].SourceEventID = ""
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(ContainSubstring("Source Event:"))
		})

		It("should include total count", func() {
			result, err := marshalFactsToText(facts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Total Facts: 1"))
		})

		It("should handle empty facts slice", func() {
			result, err := marshalFactsToText([]*careerdomain.Fact{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Facts"))
			Expect(result).To(ContainSubstring("Total Facts: 0"))
		})
	})

	Describe("marshalBurstsToCSV", func() {
		var bursts []*careerdomain.Burst

		BeforeEach(func() {
			bursts = []*careerdomain.Burst{
				{
					ID:          "burst-1",
					Name:        "Platform Migration",
					Description: "Migrated legacy platform to cloud",
					EventIDs:    []string{"event-1", "event-2", "event-3"},
					Confirmed:   true,
				},
			}
		})

		It("should include CSV header row", func() {
			result, err := marshalBurstsToCSV(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(HavePrefix("id,name,description,event_ids,confirmed\n"))
		})

		It("should marshal burst fields correctly", func() {
			result, err := marshalBurstsToCSV(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("burst-1"))
			Expect(result).To(ContainSubstring("Platform Migration"))
			Expect(result).To(ContainSubstring("true"))
		})

		It("should join event IDs with semicolon", func() {
			result, err := marshalBurstsToCSV(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("event-1;event-2;event-3"))
		})

		It("should handle unconfirmed bursts", func() {
			bursts[0].Confirmed = false
			result, err := marshalBurstsToCSV(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("false"))
		})

		It("should handle empty bursts slice", func() {
			result, err := marshalBurstsToCSV([]*careerdomain.Burst{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal("id,name,description,event_ids,confirmed\n"))
		})
	})

	Describe("marshalBurstsToText", func() {
		var bursts []*careerdomain.Burst

		BeforeEach(func() {
			confirmedAt := time.Date(2024, 1, 20, 14, 30, 0, 0, time.UTC)
			bursts = []*careerdomain.Burst{
				{
					ID:          "burst-1",
					Name:        "Platform Migration",
					Description: "Migrated legacy platform to cloud",
					EventIDs:    []string{"event-1", "event-2", "event-3"},
					Confirmed:   true,
					ConfirmedAt: &confirmedAt,
				},
			}
		})

		It("should include header", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Bursts"))
		})

		It("should include separator line", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring(strings.Repeat("=", 80)))
		})

		It("should include burst number and name", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Burst 1: Platform Migration"))
		})

		It("should include description when present", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Description: Migrated legacy platform to cloud"))
		})

		It("should include event IDs", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Event IDs: event-1, event-2, event-3"))
		})

		It("should include confirmed status", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Confirmed: true"))
		})

		It("should include confirmed at timestamp when present", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Confirmed At: 2024-01-20 14:30:00"))
		})

		It("should omit confirmed at when nil", func() {
			bursts[0].ConfirmedAt = nil
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(ContainSubstring("Confirmed At:"))
		})

		It("should omit description when empty", func() {
			bursts[0].Description = ""
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).NotTo(ContainSubstring("Description:"))
		})

		It("should include total count", func() {
			result, err := marshalBurstsToText(bursts)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Total Bursts: 1"))
		})

		It("should handle empty bursts slice", func() {
			result, err := marshalBurstsToText([]*careerdomain.Burst{})
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(ContainSubstring("Bursts"))
			Expect(result).To(ContainSubstring("Total Bursts: 0"))
		})
	})

	Describe("escapeCSV", func() {
		It("should return plain string unchanged", func() {
			result := escapeCSV("simple text")
			Expect(result).To(Equal("simple text"))
		})

		It("should escape strings containing commas", func() {
			result := escapeCSV("hello, world")
			Expect(result).To(Equal(`"hello, world"`))
		})

		It("should escape strings containing quotes", func() {
			result := escapeCSV(`say "hello"`)
			Expect(result).To(Equal(`"say ""hello"""`))
		})

		It("should escape strings containing newlines", func() {
			result := escapeCSV("line1\nline2")
			Expect(result).To(Equal("\"line1\nline2\""))
		})

		It("should handle strings with multiple special characters", func() {
			result := escapeCSV(`hello, "world"\ntest`)
			Expect(result).To(Equal(`"hello, ""world""\ntest"`))
		})

		It("should handle empty string", func() {
			result := escapeCSV("")
			Expect(result).To(Equal(""))
		})
	})
})
