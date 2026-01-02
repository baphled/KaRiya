package cv

import (
	"context"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultBulletGenerator", func() {
	var (
		generator *DefaultBulletGenerator
		log       *logger.Logger
		ctx       context.Context
	)

	BeforeEach(func() {
		log = logger.DefaultLogger()
		generator = &DefaultBulletGenerator{
			eventRepo: nil,
			factRepo:  nil,
			logger:    log,
		}
		ctx = context.Background()
	})

	It("should return empty list for no events", func() {
		bullets, err := generator.GenerateBullets(ctx, []*career.CareerEvent{}, []*career.Fact{}, "principal", []string{"hiring_manager"})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(bullets)).To(Equal(0))
	})

	It("should generate bullets from valid events", func() {
		events := []*career.CareerEvent{
			{
				ID:         "event1",
				Text:       "Implemented authentication system",
				Date:       time.Now(),
				Categories: []string{"technical"},
			},
		}

		bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(bullets)).To(BeNumerically(">", 0))
	})

	It("should filter out aspirational language", func() {
		events := []*career.CareerEvent{
			{
				ID:         "event1",
				Text:       "Will implement new system",
				Date:       time.Now(),
				Categories: []string{"technical"},
			},
		}

		bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(bullets)).To(Equal(0))
	})

	It("should cap principal role at 4 bullets", func() {
		events := make([]*career.CareerEvent, 10)
		for i := 0; i < 10; i++ {
			events[i] = &career.CareerEvent{
				ID:         "event" + string(rune('0'+i)),
				Text:       "Implemented feature " + string(rune('A'+i)),
				Date:       time.Now(),
				Categories: []string{"technical"},
			}
		}

		bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
		Expect(err).NotTo(HaveOccurred())
		Expect(len(bullets)).To(BeNumerically("<=", 4))
	})

	It("should preserve event source IDs", func() {
		events := []*career.CareerEvent{
			{
				ID:         "event1",
				Text:       "Implemented feature",
				Date:       time.Now(),
				Categories: []string{"technical"},
			},
		}

		bullets, err := generator.GenerateBullets(ctx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
		Expect(err).NotTo(HaveOccurred())

		if len(bullets) > 0 {
			Expect(bullets[0].SourceEventIDs).To(ContainElement("event1"))
		}
	})

	It("should handle context cancellation", func() {
		events := []*career.CareerEvent{
			{
				ID:   "event1",
				Text: "Some event",
				Date: time.Now(),
			},
		}

		cancelCtx, cancel := context.WithCancel(context.Background())
		cancel()

		_, err := generator.GenerateBullets(cancelCtx, events, []*career.Fact{}, "principal", []string{"hiring_manager"})
		Expect(err).To(HaveOccurred())
	})
})
