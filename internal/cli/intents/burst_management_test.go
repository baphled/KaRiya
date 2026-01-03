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
		model     *BurstManagementModel
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
		model = NewBurstManagementIntent(data)
	})

	Describe("Initialization", func() {
		It("should initialize with list state", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
			Expect(model.View()).NotTo(BeEmpty())
		})

		It("should load bursts on init", func() {
			cmd := model.Init()
			Expect(cmd).To(BeNil())
			view := model.View()
			Expect(view).To(ContainSubstring("Test Burst"))
		})

		It("should render without error", func() {
			model.Init()
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("List State", func() {
		BeforeEach(func() {
			model.Init()
		})

		It("should render list view", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Bursts"))
			Expect(view).To(ContainSubstring("Test Burst"))
		})

		It("should show empty state when no bursts", func() {
			mockRepo.bursts = []*careerdom.Burst{}
			model.Init()
			view := model.View()
			Expect(view).To(ContainSubstring("No bursts found"))
		})

		It("should display pagination info", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("Page"))
			Expect(view).To(ContainSubstring("Total"))
		})

		It("should show help text", func() {
			view := model.View()
			Expect(view).To(ContainSubstring("navigate"))
			Expect(view).To(ContainSubstring("quit"))
		})
	})

	Describe("Context Operations", func() {
		It("should load bursts", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			err := data.LoadBursts()
			Expect(err).To(BeNil())
			Expect(data.Bursts).To(HaveLen(1))
		})

		It("should create burst with valid data", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			newBurst := &careerdom.Burst{
				ID:              "new-burst",
				Name:            "New Burst",
				Description:     "New description",
				CompetencyFocus: "leadership",
				EventIDs:        []string{"event-1", "event-2"},
			}
			err := data.CreateBurst(newBurst)
			Expect(err).To(BeNil())
			Expect(data.Bursts).To(HaveLen(1))
		})

		It("should update burst", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.LoadBursts()
			burst := data.Bursts[0]
			burst.Name = "Updated Name"
			err := data.UpdateBurst(burst)
			Expect(err).To(BeNil())
		})

		It("should delete burst", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.LoadBursts()
			initialCount := len(data.Bursts)
			err := data.DeleteBurst(testBurst.ID)
			Expect(err).To(BeNil())
			Expect(data.Bursts).To(HaveLen(initialCount - 1))
		})

		It("should handle pagination", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.LoadBursts()
			pageBursts := data.GetPageBursts()
			Expect(pageBursts).NotTo(BeNil())
		})

		It("should track form errors", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.SetFormError("name", "Name is required")
			Expect(data.HasFormErrors()).To(BeTrue())
			Expect(data.FormErrors["name"]).To(Equal("Name is required"))
		})

		It("should clear form errors", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.SetFormError("name", "Name is required")
			data.ClearFormErrors()
			Expect(data.HasFormErrors()).To(BeFalse())
		})

		It("should toggle row expansion", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.ToggleRowExpansion(0)
			Expect(data.IsRowExpanded(0)).To(BeTrue())
			data.ToggleRowExpansion(0)
			Expect(data.IsRowExpanded(0)).To(BeFalse())
		})

		It("should start new burst", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.StartNewBurst()
			Expect(data.IsNewBurst).To(BeTrue())
			Expect(data.EditingBurst).NotTo(BeNil())
		})

		It("should start edit burst", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.StartEditBurst(testBurst)
			Expect(data.IsNewBurst).To(BeFalse())
			Expect(data.EditingBurst.ID).To(Equal(testBurst.ID))
		})

		It("should cancel edit", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.StartNewBurst()
			data.CancelEdit()
			Expect(data.EditingBurst).To(BeNil())
			Expect(data.IsNewBurst).To(BeFalse())
		})

		It("should select burst safely", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.LoadBursts()
			data.SelectBurst(0)
			Expect(data.GetSelectedBurst()).NotTo(BeNil())
		})

		It("should handle nil repository", func() {
			data := NewBurstManagementContext(nil, nil, ctx)
			err := data.LoadBursts()
			Expect(err).NotTo(BeNil())
		})

		It("should handle empty burst list", func() {
			mockRepo.bursts = []*careerdom.Burst{}
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			err := data.LoadBursts()
			Expect(err).To(BeNil())
			Expect(data.Bursts).To(HaveLen(0))
		})

		It("should get page bursts", func() {
			data := NewBurstManagementContext(nil, mockRepo, ctx)
			data.LoadBursts()
			pageBursts := data.GetPageBursts()
			Expect(len(pageBursts) > 0).To(BeTrue())
		})
	})

	Describe("View Methods", func() {
		It("should render view with bursts", func() {
			model.Init()
			view := model.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("should handle view updates", func() {
			model.Init()
			view := model.View()
			Expect(view).To(ContainSubstring("Bursts"))
		})
	})

	Describe("Result Handling", func() {
		It("should return result", func() {
			result := model.Result()
			Expect(result).NotTo(BeNil())
		})
	})
})

var (
	ErrMockError = errors.New("mock error")
	ErrNotFound  = errors.New("not found")
)

// Mock BurstRepository for testing
type MockBurstRepository struct {
	bursts      []*careerdom.Burst
	shouldError bool
}

func NewMockBurstRepository() *MockBurstRepository {
	return &MockBurstRepository{
		bursts:      make([]*careerdom.Burst, 0),
		shouldError: false,
	}
}

func (m *MockBurstRepository) Create(ctx context.Context, burst *careerdom.Burst) error {
	if m.shouldError {
		return ErrMockError
	}
	if burst.ID == "" {
		burst.ID = "burst-" + time.Now().Format("20060102150405")
	}
	m.bursts = append(m.bursts, burst)
	return nil
}

func (m *MockBurstRepository) GetByID(ctx context.Context, id string) (*careerdom.Burst, error) {
	if m.shouldError {
		return nil, ErrMockError
	}
	for _, b := range m.bursts {
		if b.ID == id {
			return b, nil
		}
	}
	return nil, ErrNotFound
}

func (m *MockBurstRepository) Update(ctx context.Context, burst *careerdom.Burst) error {
	if m.shouldError {
		return ErrMockError
	}
	for i, b := range m.bursts {
		if b.ID == burst.ID {
			m.bursts[i] = burst
			return nil
		}
	}
	return ErrNotFound
}

func (m *MockBurstRepository) Delete(ctx context.Context, id string) error {
	if m.shouldError {
		return ErrMockError
	}
	for i, b := range m.bursts {
		if b.ID == id {
			m.bursts = append(m.bursts[:i], m.bursts[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}

func (m *MockBurstRepository) List(ctx context.Context, filters careerrepo.BurstListFilters) ([]*careerdom.Burst, error) {
	if m.shouldError {
		return nil, ErrMockError
	}
	return m.bursts, nil
}

func (m *MockBurstRepository) Count(ctx context.Context, filters careerrepo.BurstListFilters) (int, error) {
	if m.shouldError {
		return 0, ErrMockError
	}
	return len(m.bursts), nil
}
