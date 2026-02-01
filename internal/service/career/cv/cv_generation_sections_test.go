package cv

import (
	"context"
	"io"
	"time"

	career "github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/logger"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CVGenerationService - Sections Should Be Populated", func() {
	var (
		log *logger.Logger
		ctx context.Context
	)

	BeforeEach(func() {
		log = logger.New(io.Discard, logger.InfoLevel)
		ctx = context.Background()
	})

	It("should populate sections when generating CV with real bullet generator and section builder", func() {
		// Create test events
		event := fixtures.EventWith("event-1", "Architected microservices platform", "TechCorp", "")
		event.Date = time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		event.Tags = []string{"architecture"}
		event.Categories = []string{"technical"}
		events := []*career.Event{event}

		fact := fixtures.Fact("fact-1", event.ID)
		fact.Text = "Expert in distributed systems"
		fact.CompetencyCategories = []string{"technical"}
		fact.RoleFit = "staff"
		fact.AudienceRelevance = []string{"hiring_manager"}
		facts := []*career.Fact{fact}

		config := fixtures.CVConfigWithFilters("test-cv", make(map[string]interface{}))

		// Create repositories that return test data
		eventRepo := &TestEventRepository{events: events}
		factRepo := &TestFactRepository{facts: facts}

		// Create services with REAL implementations (not empty mocks)
		// BUG-008: use BulletGenerator for role-based scoring
		configManager := NewMemoryConfigManager()
		bulletGenerator := NewBulletGenerator(log, nil) // nil uses default scoring config
		sectionBuilder := NewSectionBuilder(nil, log)   // Use REAL section builder, not empty
		dataProcessor := NewDataProcessingService(log)

		service := NewCVGenerationService(
			eventRepo,
			factRepo,
			configManager,
			bulletGenerator,
			dataProcessor,
			sectionBuilder,
			log,
		)

		// Generate CV
		cv, err := service.GenerateCVFromConfig(ctx, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(cv).NotTo(BeNil())

		// THIS IS THE KEY TEST - sections should be populated
		Expect(cv.Sections).NotTo(BeEmpty(), "CV SECTIONS MUST NOT BE EMPTY when generating from events")
		Expect(len(cv.Sections)).To(BeNumerically(">", 0), "CV should have at least one section")

		// Verify sections have content
		for _, section := range cv.Sections {
			Expect(section.ID).NotTo(BeEmpty())
			Expect(section.Title).NotTo(BeEmpty())
			// Summary sections use Summary field, others use Content array
			if section.SectionType == "summary" {
				Expect(section.Summary).NotTo(BeEmpty(), "Summary section should have Summary text")
			} else {
				Expect(len(section.Content)).To(BeNumerically(">", 0), "Section should have content groups")
			}
		}
	})
})

// TestEventRepository provides test events
type TestEventRepository struct {
	events []*career.Event
}

func (r *TestEventRepository) List(_ context.Context, _ careerrepo.EventListFilters) ([]*career.Event, error) {
	return r.events, nil
}

func (r *TestEventRepository) GetByID(ctx context.Context, id string) (*career.Event, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *TestEventRepository) Create(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *TestEventRepository) Update(ctx context.Context, event *career.Event) error {
	return nil
}

func (r *TestEventRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *TestEventRepository) Count(ctx context.Context, filters careerrepo.EventListFilters) (int, error) {
	return len(r.events), nil
}

func (r *TestEventRepository) LinkSkill(_ context.Context, _ string, _ string) error {
	return nil
}

// TestFactRepository provides test facts
type TestFactRepository struct {
	facts []*career.Fact
}

func (r *TestFactRepository) List(ctx context.Context, filters careerrepo.FactListFilters) ([]*career.Fact, error) {
	return r.facts, nil
}

func (r *TestFactRepository) GetByID(ctx context.Context, id string) (*career.Fact, error) {
	return nil, nil //nolint:nilnil // test stub - method not used in these tests
}

func (r *TestFactRepository) Create(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *TestFactRepository) Update(ctx context.Context, fact *career.Fact) error {
	return nil
}

func (r *TestFactRepository) Delete(ctx context.Context, id string) error {
	return nil
}

func (r *TestFactRepository) Count(ctx context.Context, filters careerrepo.FactListFilters) (int, error) {
	return len(r.facts), nil
}

func (r *TestFactRepository) GetBySourceEventID(ctx context.Context, eventID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}

func (r *TestFactRepository) GetBySourceBurstID(ctx context.Context, burstID string) ([]*career.Fact, error) {
	return []*career.Fact{}, nil
}
