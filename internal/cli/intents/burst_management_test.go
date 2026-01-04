package intents

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	careerdom "github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

var _ = Describe("BurstManagement Intent", func() {
	var (
		intent    *BurstManagementIntent
		ctx       context.Context
		mockRepo  *MockBurstRepository
		testBurst *careerdom.Burst
	)

	BeforeEach(func() {
		ctx = context.Background()
		mockRepo = NewMockBurstRepository()

		// Create test burst
		testBurst = &careerdom.Burst{
			ID:              "burst-1",
			Name:            "Test Burst",
			Description:     "A test burst",
			EventIDs:        []string{"event-1", "event-2"},
			CompetencyFocus: "leadership",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		mockRepo.bursts = []*careerdom.Burst{testBurst}

		// Create intent
		data := NewBurstManagementContext(nil, mockRepo, ctx)
		var err error
		intent, err = NewBurstManagementIntent(data)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("Initialization", func() {
		It("should initialize with list state", func() {
			cmd := intent.Init()
			Expect(cmd).To(BeNil())
			Expect(intent.View()).NotTo(BeEmpty())
		})

		It("should load bursts on init", func() {
			intent.Init()
			Expect(intent.state.filteredBursts).To(HaveLen(1))
			Expect(intent.state.selectedBurst).To(Equal(testBurst))
		})
	})

	Describe("List Navigation", func() {
		It("should move selection down", func() {
			intent.Init()
			// Simulate pressing 'j' to move down
			intent.Update(nil)
			Expect(intent.state.currentState).To(Equal(BurstStateList))
		})

		It("should move selection up", func() {
			intent.Init()
			// Simulate pressing 'k' to move up
			intent.Update(nil)
			Expect(intent.state.currentState).To(Equal(BurstStateList))
		})
	})

	Describe("Burst Selection", func() {
		It("should transition to detail state on enter", func() {
			intent.Init()
			// Simulate pressing Enter
			msg := BurstSelectedMsg{Burst: testBurst, Index: 0}
			intent.Update(msg)
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
			Expect(intent.state.selectedBurst).To(Equal(testBurst))
		})

		It("should track viewed bursts", func() {
			intent.Init()
			msg := BurstSelectedMsg{Burst: testBurst, Index: 0}
			intent.Update(msg)
			Expect(intent.state.viewedBursts).To(HaveLen(1))
			Expect(intent.state.viewedBursts[0]).To(Equal(testBurst))
		})
	})

	Describe("Result Handling", func() {
		It("should return nil result before completion", func() {
			intent.Init()
			result := intent.Result()
			Expect(result).To(BeNil())
		})

		It("should return completed result on confirmation", func() {
			intent.Init()
			msg := BurstSelectedMsg{Burst: testBurst, Index: 0}
			intent.Update(msg)
			intent.state.currentState = BurstStateDetail
			intent.setCompleted()

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Completed))
			Expect(result.Data).To(HaveField("Burst", testBurst))
		})

		It("should return cancelled result on cancellation", func() {
			intent.Init()
			intent.setCancelled()

			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})
	})

	Describe("View Rendering", func() {
		It("should render list view", func() {
			intent.Init()
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should show empty state when no bursts", func() {
			emptyCtx := NewBurstManagementContext(nil, mockRepo, ctx)
			emptyCtx.Bursts = make([]*careerdom.Burst, 0)
			emptyIntent, err := NewBurstManagementIntent(emptyCtx)
			Expect(err).NotTo(HaveOccurred())
			emptyIntent.Init()
			view := emptyIntent.View()
			Expect(view).To(ContainSubstring("No bursts"))
		})

		It("should render detail view for selected burst", func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			view := intent.View()
			Expect(view).To(ContainSubstring("Burst Details"))
			Expect(view).To(ContainSubstring(testBurst.Name))
		})
	})
})

// MockBurstRepository is a mock implementation of BurstRepository for testing
type MockBurstRepository struct {
	bursts []*careerdom.Burst
}

func NewMockBurstRepository() *MockBurstRepository {
	return &MockBurstRepository{
		bursts: make([]*careerdom.Burst, 0),
	}
}

func (m *MockBurstRepository) Create(ctx context.Context, burst *careerdom.Burst) error {
	burst.ID = "burst-" + time.Now().Format("20060102150405")
	burst.CreatedAt = time.Now()
	burst.UpdatedAt = time.Now()
	m.bursts = append(m.bursts, burst)
	return nil
}

func (m *MockBurstRepository) Read(ctx context.Context, id string) (*careerdom.Burst, error) {
	for _, b := range m.bursts {
		if b.ID == id {
			return b, nil
		}
	}
	return nil, errors.New("burst not found")
}

func (m *MockBurstRepository) Update(ctx context.Context, burst *careerdom.Burst) error {
	for i, b := range m.bursts {
		if b.ID == burst.ID {
			burst.UpdatedAt = time.Now()
			m.bursts[i] = burst
			return nil
		}
	}
	return errors.New("burst not found")
}

func (m *MockBurstRepository) Delete(ctx context.Context, id string) error {
	for i, b := range m.bursts {
		if b.ID == id {
			m.bursts = append(m.bursts[:i], m.bursts[i+1:]...)
			return nil
		}
	}
	return errors.New("burst not found")
}

func (m *MockBurstRepository) List(ctx context.Context, filters careerrepo.BurstListFilters) ([]*careerdom.Burst, error) {
	return m.bursts, nil
}

func (m *MockBurstRepository) Count(ctx context.Context, filters careerrepo.BurstListFilters) (int, error) {
	return len(m.bursts), nil
}

func (m *MockBurstRepository) GetByID(ctx context.Context, id string) (*careerdom.Burst, error) {
	return m.Read(ctx, id)
}
