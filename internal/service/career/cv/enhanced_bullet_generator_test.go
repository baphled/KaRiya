package cv

import (
	"context"
	"io"
	"strings"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EnhancedBulletGenerator", func() {
	var (
		generator EnhancedBulletGenerator
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		generator = NewEnhancedBulletGenerator(log)
		ctx = context.Background()
	})

	It("should return empty list for no input", func() {
		bullets, err := generator.GenerateBullets(ctx, nil, nil, nil, "principal", "hiring_manager")
		Expect(err).NotTo(HaveOccurred())
		Expect(bullets).To(BeEmpty())
	})

	It("should generate bullets from events", func() {
		events := []*career.CareerEvent{
			{
				ID:   "e1",
				Text: "Led team to deliver microservices",
				Date: time.Now(),
			},
		}

		bullets, err := generator.GenerateBullets(ctx, events, nil, nil, "principal", "hiring_manager")
		Expect(err).NotTo(HaveOccurred())
		Expect(len(bullets)).To(BeNumerically(">", 0))
	})

	It("should filter by role confidence", func() {
		bullets := []*EnhancedBullet{
			{ID: "b1", Confidence: 0.85},
			{ID: "b2", Confidence: 0.70},
		}

		filtered := generator.FilterByRole(bullets, "principal")
		Expect(len(filtered)).To(Equal(1))
	})

	It("should rank bullets by score", func() {
		bullets := []*EnhancedBullet{
			{ID: "b1", RoleScore: 0.8, AudienceScore: 0.7, MetricScore: 0.6, ImpactScore: 0.7, Confidence: 0.8},
			{ID: "b2", RoleScore: 0.5, AudienceScore: 0.5, MetricScore: 0.5, ImpactScore: 0.5, Confidence: 0.5},
		}

		ranked := generator.RankByRelevance(bullets, "principal", "hiring_manager")
		Expect(ranked[0].ID).To(Equal("b1"))
		Expect(ranked[1].ID).To(Equal("b2"))
	})

	It("should enhance bullet wording", func() {
		bullet := &EnhancedBullet{
			ID:   "b1",
			Text: "Worked on authentication system",
		}

		enhanced, err := generator.EnhanceWording(bullet, "principal")
		Expect(err).NotTo(HaveOccurred())
		Expect(strings.Contains(enhanced.EnhancedText, "Led")).To(BeTrue())
	})
})
