package models

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("StringSlice", func() {
	Describe("Scan", func() {
		var s StringSlice

		BeforeEach(func() {
			s = nil
		})

		Context("when value is nil", func() {
			It("sets the slice to nil", func() {
				err := s.Scan(nil)

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})

		Context("when value is a string", func() {
			It("splits a comma-separated string", func() {
				err := s.Scan("a,b,c")

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"a", "b", "c"}))
			})

			It("handles a single element", func() {
				err := s.Scan("solo")

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"solo"}))
			})

			It("sets to nil for an empty string", func() {
				err := s.Scan("")

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})

		Context("when value is []byte", func() {
			It("splits a comma-separated byte slice", func() {
				err := s.Scan([]byte("x,y,z"))

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"x", "y", "z"}))
			})

			It("handles a single element", func() {
				err := s.Scan([]byte("only"))

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(Equal(StringSlice{"only"}))
			})

			It("sets to nil for empty bytes", func() {
				err := s.Scan([]byte(""))

				Expect(err).NotTo(HaveOccurred())
				Expect(s).To(BeNil())
			})
		})

		Context("when value is an unsupported type", func() {
			It("returns an error for int", func() {
				err := s.Scan(42)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported type"))
			})

			It("returns an error for bool", func() {
				err := s.Scan(true)

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unsupported type"))
			})
		})
	})

	Describe("Value", func() {
		Context("when slice is nil", func() {
			It("returns an empty string", func() {
				var s StringSlice
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(""))
			})
		})

		Context("when slice is empty", func() {
			It("returns an empty string", func() {
				s := StringSlice{}
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal(""))
			})
		})

		Context("when slice has one element", func() {
			It("returns the element without commas", func() {
				s := StringSlice{"single"}
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("single"))
			})
		})

		Context("when slice has multiple elements", func() {
			It("returns a comma-separated string", func() {
				s := StringSlice{"a", "b", "c"}
				val, err := s.Value()

				Expect(err).NotTo(HaveOccurred())
				Expect(val).To(Equal("a,b,c"))
			})
		})
	})

	Describe("Scan-Value round trip", func() {
		It("preserves data through serialisation and deserialisation", func() {
			original := StringSlice{"technical", "leadership", "mentoring"}
			val, err := original.Value()
			Expect(err).NotTo(HaveOccurred())

			var restored StringSlice
			err = restored.Scan(val)
			Expect(err).NotTo(HaveOccurred())
			Expect(restored).To(Equal(original))
		})
	})
})

var _ = Describe("Skill Model", func() {
	Describe("TableName", func() {
		It("returns skills", func() {
			Expect(Skill{}.TableName()).To(Equal("skills"))
		})
	})

	Describe("ToDomain", func() {
		It("converts all fields to domain Skill", func() {
			now := time.Now()
			years := 5
			lastUsed := now.AddDate(0, -1, 0)
			model := &Skill{
				ID:        "skill-1",
				Name:      "Go",
				Category:  "backend",
				Level:     "expert",
				YearsUsed: &years,
				LastUsed:  &lastUsed,
				CreatedAt: now,
				UpdatedAt: now,
			}

			domain := model.ToDomain()

			Expect(domain.ID).To(Equal("skill-1"))
			Expect(domain.Name).To(Equal("Go"))
			Expect(domain.Category).To(Equal("backend"))
			Expect(domain.Level).To(Equal("expert"))
			Expect(*domain.YearsUsed).To(Equal(5))
			Expect(*domain.LastUsed).To(Equal(lastUsed))
			Expect(domain.CreatedAt).To(Equal(now))
			Expect(domain.UpdatedAt).To(Equal(now))
		})

		It("handles nil optional fields", func() {
			model := &Skill{
				ID:        "skill-2",
				Name:      "Python",
				Category:  "backend",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			domain := model.ToDomain()

			Expect(domain.YearsUsed).To(BeNil())
			Expect(domain.LastUsed).To(BeNil())
		})
	})

	Describe("SkillFromDomain", func() {
		It("converts all fields from domain Skill", func() {
			years := 3
			lastUsed := time.Now().AddDate(0, -2, 0)
			domainSkill := fixtures.SkillWithYears("skill-d1", "Ruby", "backend", years)
			domainSkill.Level = "advanced"
			domainSkill.LastUsed = &lastUsed

			model := SkillFromDomain(domainSkill)

			Expect(model.ID).To(Equal("skill-d1"))
			Expect(model.Name).To(Equal("Ruby"))
			Expect(model.Category).To(Equal("backend"))
			Expect(model.Level).To(Equal("advanced"))
			Expect(*model.YearsUsed).To(Equal(3))
			Expect(*model.LastUsed).To(Equal(lastUsed))
			Expect(model.CreatedAt).To(Equal(domainSkill.CreatedAt))
			Expect(model.UpdatedAt).To(Equal(domainSkill.UpdatedAt))
		})

		It("handles nil optional fields", func() {
			domainSkill := fixtures.SkillWith("skill-d2", "Docker", "devops", "")

			model := SkillFromDomain(domainSkill)

			Expect(model.YearsUsed).To(BeNil())
			Expect(model.LastUsed).To(BeNil())
		})
	})

	Describe("ToDomain-FromDomain round trip", func() {
		It("preserves all Skill data", func() {
			years := 7
			lastUsed := time.Now().Truncate(time.Second).AddDate(0, -3, 0)
			original := fixtures.SkillWithYears("skill-rt", "Kubernetes", "devops", years)
			original.Level = "expert"
			original.LastUsed = &lastUsed
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := SkillFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.ID).To(Equal(original.ID))
			Expect(restored.Name).To(Equal(original.Name))
			Expect(restored.Category).To(Equal(original.Category))
			Expect(restored.Level).To(Equal(original.Level))
			Expect(*restored.YearsUsed).To(Equal(*original.YearsUsed))
			Expect(*restored.LastUsed).To(Equal(*original.LastUsed))
			Expect(restored.CreatedAt).To(Equal(original.CreatedAt))
			Expect(restored.UpdatedAt).To(Equal(original.UpdatedAt))
		})
	})
})

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

var _ = Describe("Fact Model", func() {
	Describe("TableName", func() {
		It("returns facts", func() {
			Expect(Fact{}.TableName()).To(Equal("facts"))
		})
	})

	Describe("ToDomain", func() {
		It("converts all fields to domain Fact", func() {
			now := time.Now()
			model := &Fact{
				ID:                   "fact-1",
				Text:                 "Led critical migration",
				CompetencyCategories: StringSlice{"technical", "leadership"},
				RoleFit:              "staff",
				AudienceRelevance:    StringSlice{"hiring_manager", "peer"},
				StrengthSignal:       "high",
				SourceEventID:        "evt-src",
				SourceBurstID:        "burst-src",
				CreatedAt:            now,
				UpdatedAt:            now,
			}

			domain := model.ToDomain()

			Expect(domain.ID).To(Equal("fact-1"))
			Expect(domain.Text).To(Equal("Led critical migration"))
			Expect(domain.CompetencyCategories).To(Equal([]string{"technical", "leadership"}))
			Expect(domain.RoleFit).To(Equal(career.RoleFit("staff")))
			Expect(domain.AudienceRelevance).To(Equal([]string{"hiring_manager", "peer"}))
			Expect(domain.StrengthSignal).To(Equal("high"))
			Expect(domain.SourceEventID).To(Equal("evt-src"))
			Expect(domain.SourceBurstID).To(Equal("burst-src"))
			Expect(domain.CreatedAt).To(Equal(now))
			Expect(domain.UpdatedAt).To(Equal(now))
		})

		It("handles nil slice fields", func() {
			model := &Fact{
				ID:      "fact-2",
				Text:    "Minimal fact",
				RoleFit: "senior_ic",
			}

			domain := model.ToDomain()

			Expect(domain.CompetencyCategories).To(BeNil())
			Expect(domain.AudienceRelevance).To(BeNil())
		})
	})

	Describe("FactFromDomain", func() {
		It("converts all fields from domain Fact", func() {
			domainFact := fixtures.FactWithCategories("fact-d1", "Architected platform", "evt-1", []string{"technical"}, []string{"peer"})
			domainFact.RoleFit = career.RoleFitPrincipal
			domainFact.StrengthSignal = "leadership capability"

			model := FactFromDomain(domainFact)

			Expect(model.ID).To(Equal("fact-d1"))
			Expect(model.Text).To(Equal("Architected platform"))
			Expect([]string(model.CompetencyCategories)).To(Equal([]string{"technical"}))
			Expect(model.RoleFit).To(Equal(string(career.RoleFitPrincipal)))
			Expect([]string(model.AudienceRelevance)).To(Equal([]string{"peer"}))
			Expect(model.StrengthSignal).To(Equal("leadership capability"))
			Expect(model.SourceEventID).To(Equal("evt-1"))
			Expect(model.CreatedAt).To(Equal(domainFact.CreatedAt))
			Expect(model.UpdatedAt).To(Equal(domainFact.UpdatedAt))
		})

		It("handles nil slice fields", func() {
			domainFact := fixtures.FactWith("fact-d2", "Minimal fact")
			domainFact.RoleFit = career.RoleFitStaff

			model := FactFromDomain(domainFact)

			Expect(model.CompetencyCategories).To(BeNil())
			Expect(model.AudienceRelevance).To(BeNil())
		})
	})

	Describe("ToDomain-FromDomain round trip", func() {
		It("preserves all Fact data", func() {
			original := fixtures.FactWithCategories("fact-rt", "Round trip fact", "evt-rt", []string{"technical", "mentoring"}, []string{"hiring_manager", "recruiter"})
			original.RoleFit = career.RoleFitEM
			original.StrengthSignal = "team growth"
			original.SourceBurstID = "burst-rt"
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := FactFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.ID).To(Equal(original.ID))
			Expect(restored.Text).To(Equal(original.Text))
			Expect(restored.CompetencyCategories).To(Equal(original.CompetencyCategories))
			Expect(restored.RoleFit).To(Equal(original.RoleFit))
			Expect(restored.AudienceRelevance).To(Equal(original.AudienceRelevance))
			Expect(restored.StrengthSignal).To(Equal(original.StrengthSignal))
			Expect(restored.SourceEventID).To(Equal(original.SourceEventID))
			Expect(restored.SourceBurstID).To(Equal(original.SourceBurstID))
			Expect(restored.CreatedAt).To(Equal(original.CreatedAt))
			Expect(restored.UpdatedAt).To(Equal(original.UpdatedAt))
		})
	})
})

var _ = Describe("Burst Model", func() {
	Describe("TableName", func() {
		It("returns bursts", func() {
			Expect(Burst{}.TableName()).To(Equal("bursts"))
		})
	})

	Describe("ToDomain", func() {
		It("converts all fields to domain Burst", func() {
			now := time.Now()
			confirmedAt := now.Add(-time.Hour)
			model := &Burst{
				ID:          "burst-1",
				Name:        "Platform Initiative",
				Description: "Series of platform improvements",
				EventIDs:    StringSlice{"evt-1", "evt-2", "evt-3"},
				Confirmed:   true,
				ConfirmedAt: &confirmedAt,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			domain := model.ToDomain()

			Expect(domain.ID).To(Equal("burst-1"))
			Expect(domain.Name).To(Equal("Platform Initiative"))
			Expect(domain.Description).To(Equal("Series of platform improvements"))
			Expect(domain.EventIDs).To(Equal([]string{"evt-1", "evt-2", "evt-3"}))
			Expect(domain.Confirmed).To(BeTrue())
			Expect(domain.ConfirmedAt).NotTo(BeNil())
			Expect(*domain.ConfirmedAt).To(Equal(confirmedAt))
			Expect(domain.CreatedAt).To(Equal(now))
			Expect(domain.UpdatedAt).To(Equal(now))
		})

		It("handles unconfirmed burst with nil ConfirmedAt", func() {
			model := &Burst{
				ID:        "burst-2",
				Name:      "Unconfirmed Burst",
				EventIDs:  StringSlice{"evt-1", "evt-2"},
				Confirmed: false,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			domain := model.ToDomain()

			Expect(domain.Confirmed).To(BeFalse())
			Expect(domain.ConfirmedAt).To(BeNil())
		})
	})

	Describe("BurstFromDomain", func() {
		It("converts all fields from domain Burst", func() {
			domainBurst := fixtures.BurstConfirmed("burst-d1", "evt-a", "evt-b")
			domainBurst.Description = "Database migration activities"

			model := BurstFromDomain(domainBurst)

			Expect(model.ID).To(Equal("burst-d1"))
			Expect(model.Name).To(Equal(domainBurst.Name))
			Expect(model.Description).To(Equal("Database migration activities"))
			Expect([]string(model.EventIDs)).To(Equal([]string{"evt-a", "evt-b"}))
			Expect(model.Confirmed).To(BeTrue())
			Expect(model.ConfirmedAt).NotTo(BeNil())
			Expect(model.CreatedAt).To(Equal(domainBurst.CreatedAt))
			Expect(model.UpdatedAt).To(Equal(domainBurst.UpdatedAt))
		})

		It("handles unconfirmed burst", func() {
			domainBurst := fixtures.Burst("burst-d2", "evt-x", "evt-y")

			model := BurstFromDomain(domainBurst)

			Expect(model.Confirmed).To(BeFalse())
			Expect(model.ConfirmedAt).To(BeNil())
		})
	})

	Describe("ToDomain-FromDomain round trip", func() {
		It("preserves all Burst data including confirmation", func() {
			original := fixtures.BurstConfirmed("burst-rt", "evt-r1", "evt-r2")
			original.Description = "Testing round trip"
			original.EventIDs = append(original.EventIDs, "evt-r3")
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := BurstFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.ID).To(Equal(original.ID))
			Expect(restored.Name).To(Equal(original.Name))
			Expect(restored.Description).To(Equal(original.Description))
			Expect(restored.EventIDs).To(Equal(original.EventIDs))
			Expect(restored.Confirmed).To(Equal(original.Confirmed))
			Expect(restored.ConfirmedAt).NotTo(BeNil())
			Expect(restored.CreatedAt).To(Equal(original.CreatedAt))
			Expect(restored.UpdatedAt).To(Equal(original.UpdatedAt))
		})

		It("preserves unconfirmed state through round trip", func() {
			original := fixtures.Burst("burst-rt-unconfirmed", "evt-u1", "evt-u2")
			original.CreatedAt = original.CreatedAt.Truncate(time.Second)
			original.UpdatedAt = original.UpdatedAt.Truncate(time.Second)

			model := BurstFromDomain(original)
			restored := model.ToDomain()

			Expect(restored.Confirmed).To(BeFalse())
			Expect(restored.ConfirmedAt).To(BeNil())
		})
	})
})
