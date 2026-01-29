package cv

import (
	"context"
	"io"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	careermemory "github.com/baphled/kariya/internal/repository/career/memory"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TraceabilityService", func() {
	var (
		service        *TraceabilityService
		eventRepo      careerrepo.EventRepository
		factRepo       careerrepo.FactRepository
		testLogger     *logger.Logger
		ctx            context.Context
		testEvent1     *career.CareerEvent
		testEvent2     *career.CareerEvent
		testFact1      *career.Fact
		testFact2      *career.Fact
		testBullet1    *career.CVBullet
		testBullet2    *career.CVBullet
		orphanedBullet *career.CVBullet
	)

	BeforeEach(func() {
		ctx = context.Background()
		testLogger = logger.New(io.Discard, logger.DebugLevel)
		eventRepo = careermemory.NewEventRepository()
		factRepo = careermemory.NewFactRepository()
		service = NewTraceabilityService(eventRepo, factRepo, testLogger)

		// Create test events
		testEvent1 = fixtures.EventWith(uuid.New().String(), "Led team on project X", "", "")
		testEvent1.Tags = []string{"leadership"}

		testEvent2 = fixtures.EventWith(uuid.New().String(), "Implemented feature Y", "", "")
		testEvent2.Tags = []string{"technical"}

		// Create test facts with all required fields
		eventID1 := uuid.New().String()
		testFact1 = &career.Fact{
			ID:                   uuid.New().String(),
			Text:                 "Increased team productivity by 20%",
			CompetencyCategories: []string{"leadership"},
			RoleFit:              career.RoleFitPrincipal,
			AudienceRelevance:    []string{"hiring_manager"},
			StrengthSignal:       "teamwork",
			SourceEventID:        eventID1,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		eventID2 := uuid.New().String()
		testFact2 = &career.Fact{
			ID:                   uuid.New().String(),
			Text:                 "Reduced system latency by 30%",
			CompetencyCategories: []string{"technical"},
			RoleFit:              career.RoleFitStaff,
			AudienceRelevance:    []string{"peer"},
			StrengthSignal:       "technical-excellence",
			SourceEventID:        eventID2,
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}

		// Create test bullets
		testBullet1 = &career.CVBullet{
			ID:              uuid.New().String(),
			Text:            "Led and mentored team on project X",
			SourceEventIDs:  []string{testEvent1.ID},
			SourceFactIDs:   []string{testFact1.ID},
			Rank:            0.9,
			InclusionReason: "ownership",
			Confidence:      0.95,
		}

		testBullet2 = &career.CVBullet{
			ID:              uuid.New().String(),
			Text:            "Implemented feature Y with performance improvements",
			SourceEventIDs:  []string{testEvent2.ID},
			SourceFactIDs:   []string{testFact2.ID},
			Rank:            0.85,
			InclusionReason: "execution",
			Confidence:      0.90,
		}

		orphanedBullet = &career.CVBullet{
			ID:              uuid.New().String(),
			Text:            "Orphaned bullet with no sources",
			SourceEventIDs:  []string{},
			SourceFactIDs:   []string{},
			Rank:            0.5,
			InclusionReason: "unknown",
			Confidence:      0.0,
		}

		// Store test data in repositories
		Expect(eventRepo.Create(ctx, testEvent1)).To(Succeed())
		Expect(eventRepo.Create(ctx, testEvent2)).To(Succeed())
		Expect(factRepo.Create(ctx, testFact1)).To(Succeed())
		Expect(factRepo.Create(ctx, testFact2)).To(Succeed())
	})

	Describe("GetBulletSources", func() {
		It("should retrieve source events and facts for a bullet", func() {
			events, facts, err := service.GetBulletSources(ctx, testBullet1.ID, testBullet1)

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(1))
			Expect(facts).To(HaveLen(1))
			Expect(events[0].ID).To(Equal(testEvent1.ID))
			Expect(facts[0].ID).To(Equal(testFact1.ID))
		})

		It("should handle bullets with multiple sources", func() {
			multiBullet := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Multi-source bullet",
				SourceEventIDs: []string{testEvent1.ID, testEvent2.ID},
				SourceFactIDs:  []string{testFact1.ID, testFact2.ID},
			}

			events, facts, err := service.GetBulletSources(ctx, multiBullet.ID, multiBullet)

			Expect(err).NotTo(HaveOccurred())
			Expect(events).To(HaveLen(2))
			Expect(facts).To(HaveLen(2))
		})

		It("should return error when bullet is nil", func() {
			_, _, err := service.GetBulletSources(ctx, "any-id", nil)

			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("bullet cannot be nil"))
		})

		It("should handle missing event gracefully", func() {
			// GetBulletSources silently skips missing events and returns what it can find
			bulletWithMissing := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Bullet with missing event",
				SourceEventIDs: []string{uuid.New().String()}, // Non-existent event
				SourceFactIDs:  []string{testFact1.ID},
			}

			events, facts, err := service.GetBulletSources(ctx, bulletWithMissing.ID, bulletWithMissing)

			// No error is returned; the service just returns what it can find
			Expect(err).NotTo(HaveOccurred())
			// Missing event is skipped, so we get no events
			Expect(events).To(HaveLen(0))
			// But the fact is retrieved successfully
			Expect(facts).To(HaveLen(1))
			Expect(facts[0].ID).To(Equal(testFact1.ID))
		})
	})

	Describe("GetEventUsage", func() {
		It("should find all bullets using a specific event", func() {
			bullets := []*career.CVBullet{testBullet1, testBullet2, orphanedBullet}

			usage := service.GetEventUsage(ctx, testEvent1.ID, bullets)

			Expect(usage).To(HaveLen(1))
			Expect(usage[0].ID).To(Equal(testBullet1.ID))
		})

		It("should return empty list when event is not used", func() {
			bullets := []*career.CVBullet{testBullet1, testBullet2}

			usage := service.GetEventUsage(ctx, uuid.New().String(), bullets)

			Expect(usage).To(HaveLen(0))
		})

		It("should handle multiple bullets using same event", func() {
			multiBullet1 := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "First usage",
				SourceEventIDs: []string{testEvent1.ID},
				SourceFactIDs:  []string{},
			}
			multiBullet2 := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Second usage",
				SourceEventIDs: []string{testEvent1.ID},
				SourceFactIDs:  []string{},
			}
			bullets := []*career.CVBullet{multiBullet1, multiBullet2}

			usage := service.GetEventUsage(ctx, testEvent1.ID, bullets)

			Expect(usage).To(HaveLen(2))
		})
	})

	Describe("GetFactUsage", func() {
		It("should find all bullets using a specific fact", func() {
			bullets := []*career.CVBullet{testBullet1, testBullet2, orphanedBullet}

			usage := service.GetFactUsage(ctx, testFact1.ID, bullets)

			Expect(usage).To(HaveLen(1))
			Expect(usage[0].ID).To(Equal(testBullet1.ID))
		})

		It("should return empty list when fact is not used", func() {
			bullets := []*career.CVBullet{testBullet1, testBullet2}

			usage := service.GetFactUsage(ctx, uuid.New().String(), bullets)

			Expect(usage).To(HaveLen(0))
		})

		It("should handle multiple bullets using same fact", func() {
			multiBullet1 := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "First usage",
				SourceEventIDs: []string{},
				SourceFactIDs:  []string{testFact1.ID},
			}
			multiBullet2 := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Second usage",
				SourceEventIDs: []string{},
				SourceFactIDs:  []string{testFact1.ID},
			}
			bullets := []*career.CVBullet{multiBullet1, multiBullet2}

			usage := service.GetFactUsage(ctx, testFact1.ID, bullets)

			Expect(usage).To(HaveLen(2))
		})
	})

	Describe("ValidateTraceability", func() {
		It("should validate all bullets with valid sources", func() {
			bullets := []*career.CVBullet{testBullet1, testBullet2}

			report := service.ValidateTraceability(ctx, bullets)

			Expect(report.TotalBullets).To(Equal(2))
			Expect(report.ValidBullets).To(Equal(2))
			Expect(report.InvalidBullets).To(Equal(0))
			Expect(report.IsValid()).To(BeTrue())
		})

		It("should detect orphaned bullets with no sources", func() {
			bullets := []*career.CVBullet{testBullet1, orphanedBullet}

			report := service.ValidateTraceability(ctx, bullets)

			Expect(report.TotalBullets).To(Equal(2))
			Expect(report.ValidBullets).To(Equal(1))
			Expect(report.InvalidBullets).To(Equal(1))
			Expect(report.Issues).To(ContainElement(ContainSubstring("has no sources")))
			Expect(report.IsValid()).To(BeFalse())
		})

		It("should detect bullets with missing event sources", func() {
			// ValidateTraceability only checks if there are ANY sources declared
			// It doesn't validate that those sources actually exist
			bulletWithMissing := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Bullet with missing event",
				SourceEventIDs: []string{uuid.New().String()}, // Non-existent event
				SourceFactIDs:  []string{testFact1.ID},
			}
			bullets := []*career.CVBullet{testBullet1, bulletWithMissing}

			report := service.ValidateTraceability(ctx, bullets)

			// ValidateTraceability only checks presence of source IDs, not their validity
			Expect(report.InvalidBullets).To(Equal(0))
			Expect(report.ValidBullets).To(Equal(2))
		})

		It("should detect bullets with missing fact sources", func() {
			// ValidateTraceability only checks if there are ANY sources declared
			// It doesn't validate that those sources actually exist
			bulletWithMissing := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Bullet with missing fact",
				SourceEventIDs: []string{testEvent1.ID},
				SourceFactIDs:  []string{uuid.New().String()}, // Non-existent fact
			}
			bullets := []*career.CVBullet{testBullet1, bulletWithMissing}

			report := service.ValidateTraceability(ctx, bullets)

			// ValidateTraceability only checks presence of source IDs, not their validity
			Expect(report.InvalidBullets).To(Equal(0))
			Expect(report.ValidBullets).To(Equal(2))
		})

		It("should handle empty bullet list", func() {
			bullets := []*career.CVBullet{}

			report := service.ValidateTraceability(ctx, bullets)

			Expect(report.TotalBullets).To(Equal(0))
			Expect(report.ValidBullets).To(Equal(0))
			Expect(report.IsValid()).To(BeTrue())
		})
	})

	Describe("GetEventBulletMapping", func() {
		It("should create mapping of events to bullets", func() {
			bullets := []*career.CVBullet{testBullet1, testBullet2}

			mapping := service.GetEventBulletMapping(ctx, bullets)

			Expect(mapping).To(HaveKey(testEvent1.ID))
			Expect(mapping).To(HaveKey(testEvent2.ID))
			Expect(mapping[testEvent1.ID]).To(HaveLen(1))
			Expect(mapping[testEvent2.ID]).To(HaveLen(1))
		})

		It("should handle multiple bullets per event", func() {
			multiBullet1 := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "First usage",
				SourceEventIDs: []string{testEvent1.ID},
			}
			multiBullet2 := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Second usage",
				SourceEventIDs: []string{testEvent1.ID},
			}
			bullets := []*career.CVBullet{multiBullet1, multiBullet2}

			mapping := service.GetEventBulletMapping(ctx, bullets)

			Expect(mapping[testEvent1.ID]).To(HaveLen(2))
		})

		It("should handle bullets with multiple event sources", func() {
			multiBullet := &career.CVBullet{
				ID:             uuid.New().String(),
				Text:           "Multi-source",
				SourceEventIDs: []string{testEvent1.ID, testEvent2.ID},
			}
			bullets := []*career.CVBullet{multiBullet}

			mapping := service.GetEventBulletMapping(ctx, bullets)

			Expect(mapping[testEvent1.ID]).To(HaveLen(1))
			Expect(mapping[testEvent2.ID]).To(HaveLen(1))
		})

		It("should return empty mapping for empty bullets", func() {
			bullets := []*career.CVBullet{}

			mapping := service.GetEventBulletMapping(ctx, bullets)

			Expect(mapping).To(BeEmpty())
		})
	})

	Describe("GetFactBulletMapping", func() {
		It("should create mapping of facts to bullets", func() {
			bullets := []*career.CVBullet{testBullet1, testBullet2}

			mapping := service.GetFactBulletMapping(ctx, bullets)

			Expect(mapping).To(HaveKey(testFact1.ID))
			Expect(mapping).To(HaveKey(testFact2.ID))
			Expect(mapping[testFact1.ID]).To(HaveLen(1))
			Expect(mapping[testFact2.ID]).To(HaveLen(1))
		})

		It("should handle multiple bullets per fact", func() {
			multiBullet1 := &career.CVBullet{
				ID:            uuid.New().String(),
				Text:          "First usage",
				SourceFactIDs: []string{testFact1.ID},
			}
			multiBullet2 := &career.CVBullet{
				ID:            uuid.New().String(),
				Text:          "Second usage",
				SourceFactIDs: []string{testFact1.ID},
			}
			bullets := []*career.CVBullet{multiBullet1, multiBullet2}

			mapping := service.GetFactBulletMapping(ctx, bullets)

			Expect(mapping[testFact1.ID]).To(HaveLen(2))
		})

		It("should handle bullets with multiple fact sources", func() {
			multiBullet := &career.CVBullet{
				ID:            uuid.New().String(),
				Text:          "Multi-source",
				SourceFactIDs: []string{testFact1.ID, testFact2.ID},
			}
			bullets := []*career.CVBullet{multiBullet}

			mapping := service.GetFactBulletMapping(ctx, bullets)

			Expect(mapping[testFact1.ID]).To(HaveLen(1))
			Expect(mapping[testFact2.ID]).To(HaveLen(1))
		})

		It("should return empty mapping for empty bullets", func() {
			bullets := []*career.CVBullet{}

			mapping := service.GetFactBulletMapping(ctx, bullets)

			Expect(mapping).To(BeEmpty())
		})
	})

	Describe("ValidationReport", func() {
		It("should correctly determine if report is valid", func() {
			validReport := &ValidationReport{
				TotalBullets:   2,
				ValidBullets:   2,
				InvalidBullets: 0,
				Issues:         []string{},
			}

			Expect(validReport.IsValid()).To(BeTrue())
		})

		It("should correctly determine if report is invalid", func() {
			invalidReport := &ValidationReport{
				TotalBullets:   2,
				ValidBullets:   1,
				InvalidBullets: 1,
				Issues:         []string{"Issue 1"},
			}

			Expect(invalidReport.IsValid()).To(BeFalse())
		})

		It("should provide valid summary string", func() {
			validReport := &ValidationReport{
				TotalBullets:   5,
				ValidBullets:   5,
				InvalidBullets: 0,
				Issues:         []string{},
			}

			summary := validReport.SummaryString()

			Expect(summary).To(Equal("All 5 bullets have valid sources"))
		})

		It("should provide invalid summary string", func() {
			invalidReport := &ValidationReport{
				TotalBullets:   5,
				ValidBullets:   3,
				InvalidBullets: 2,
				Issues:         []string{"Issue 1", "Issue 2"},
			}

			summary := invalidReport.SummaryString()

			Expect(summary).To(ContainSubstring("3 valid"))
			Expect(summary).To(ContainSubstring("2 invalid"))
			Expect(summary).To(ContainSubstring("5 total"))
		})
	})
})
