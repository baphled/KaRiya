package intents

import (
	"context"
	"errors"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
			Expect(intent.state.currentState).To(Equal(BurstStateList))
		})

		It("should have bursts in filtered list", func() {
			intent.Init()
			Expect(len(intent.state.filteredBursts)).To(Equal(1))
		})
	})

	Describe("Update - List View", func() {
		BeforeEach(func() {
			intent.Init()
		})

		It("should move selection down", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyDown})
			Expect(intent.state.selectedIndex).To(Equal(0))
		})

		It("should transition to detail view on enter", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
		})

		It("should cancel on q key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})
	})

	Describe("Result Handling", func() {
		It("should return completed result on success", func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
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
			emptyRepo := NewMockBurstRepository()
			emptyCtx := NewBurstManagementContext(nil, emptyRepo, ctx)
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

	Describe("Enhanced Detail View - Events", func() {
		var (
			testEvent1 *careerdom.CareerEvent
			testEvent2 *careerdom.CareerEvent
		)

		BeforeEach(func() {
			// Create test events
			testEvent1 = &careerdom.CareerEvent{
				ID:   "event-1",
				Text: "First event in burst",
				Date: time.Now().AddDate(0, -1, 0),
				Tags: []string{"tag1", "tag2"},
			}
			testEvent2 = &careerdom.CareerEvent{
				ID:   "event-2",
				Text: "Second event in burst",
				Date: time.Now(),
				Tags: []string{"tag3"},
			}

			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
		})

		It("should transition to events view on 'e' key", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'e'}})
			Expect(intent.state.currentState).To(Equal(BurstStateDetailEvents))
			Expect(intent.state.loadingEvents).To(BeTrue())
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle events loaded message", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.loadingEvents = true

			msg := BurstEventsLoadedMsg{
				Events: []*careerdom.CareerEvent{testEvent1, testEvent2},
			}

			intent.Update(msg)
			Expect(intent.state.loadingEvents).To(BeFalse())
			Expect(len(intent.state.burstEvents)).To(Equal(2))
			Expect(intent.state.burstEvents[0].ID).To(Equal("event-1"))
		})

		It("should render loading state while loading events", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.loadingEvents = true

			view := intent.View()
			Expect(view).To(ContainSubstring("Loading events"))
		})

		It("should render events view with event details", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.burstEvents = []*careerdom.CareerEvent{testEvent1, testEvent2}
			intent.state.loadingEvents = false

			view := intent.View()
			Expect(view).To(ContainSubstring("Events in Burst"))
			Expect(view).To(ContainSubstring(testBurst.Name))
			Expect(view).To(ContainSubstring(testEvent1.Text))
			Expect(view).To(ContainSubstring(testEvent2.Text))
			Expect(view).To(ContainSubstring("tag1"))
		})

		It("should show empty state when no events found", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.burstEvents = []*careerdom.CareerEvent{}
			intent.state.loadingEvents = false

			view := intent.View()
			Expect(view).To(ContainSubstring("No events found"))
		})

		It("should return to detail view on escape key", func() {
			intent.state.currentState = BurstStateDetailEvents

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
		})

		It("should cancel on q key from events view", func() {
			intent.state.currentState = BurstStateDetailEvents

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should handle no burst selected error when loading events", func() {
			intent.state.selectedBurst = nil

			cmd := intent.loadEventsForBurst()
			msg := cmd().(BurstEventsLoadedMsg)

			Expect(msg.Error).To(HaveOccurred())
			Expect(msg.Error.Error()).To(ContainSubstring("no burst selected"))
		})

		It("should handle error in events loaded message", func() {
			intent.state.currentState = BurstStateDetailEvents
			intent.state.loadingEvents = true

			msg := BurstEventsLoadedMsg{
				Error: errors.New("failed to load events"),
			}

			intent.Update(msg)
			Expect(intent.state.loadingEvents).To(BeFalse())
			Expect(len(intent.state.burstEvents)).To(Equal(0))
		})
	})

	Describe("Enhanced Detail View - Facts", func() {
		var (
			testFact1 *careerdom.Fact
			testFact2 *careerdom.Fact
		)

		BeforeEach(func() {
			// Create test facts
			testFact1 = &careerdom.Fact{
				ID:                   "fact-1",
				Text:                 "Led team of 5 engineers",
				CompetencyCategories: []string{"leadership", "technical"},
				StrengthSignal:       "strong",
				SourceBurstID:        "burst-1",
			}
			testFact2 = &careerdom.Fact{
				ID:                   "fact-2",
				Text:                 "Improved system performance by 50%",
				CompetencyCategories: []string{"technical"},
				StrengthSignal:       "moderate",
				SourceBurstID:        "burst-1",
			}

			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst
		})

		It("should transition to facts view on 'f' key", func() {
			cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'f'}})
			Expect(intent.state.currentState).To(Equal(BurstStateDetailFacts))
			Expect(intent.state.loadingFacts).To(BeTrue())
			Expect(cmd).NotTo(BeNil())
		})

		It("should handle facts loaded message", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.loadingFacts = true

			msg := BurstFactsLoadedMsg{
				Facts: []*careerdom.Fact{testFact1, testFact2},
			}

			intent.Update(msg)
			Expect(intent.state.loadingFacts).To(BeFalse())
			Expect(len(intent.state.burstFacts)).To(Equal(2))
			Expect(intent.state.burstFacts[0].ID).To(Equal("fact-1"))
		})

		It("should render loading state while loading facts", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.loadingFacts = true

			view := intent.View()
			Expect(view).To(ContainSubstring("Loading facts"))
		})

		It("should render facts view with fact details", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.burstFacts = []*careerdom.Fact{testFact1, testFact2}
			intent.state.loadingFacts = false

			view := intent.View()
			Expect(view).To(ContainSubstring("Facts from Burst"))
			Expect(view).To(ContainSubstring(testBurst.Name))
			Expect(view).To(ContainSubstring(testFact1.Text))
			Expect(view).To(ContainSubstring(testFact2.Text))
			Expect(view).To(ContainSubstring("leadership"))
			Expect(view).To(ContainSubstring("strong"))
		})

		It("should show helpful message when no facts exist", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.burstFacts = []*careerdom.Fact{}
			intent.state.loadingFacts = false

			view := intent.View()
			Expect(view).To(ContainSubstring("No facts extracted yet"))
			Expect(view).To(ContainSubstring("Confirm the burst"))
		})

		It("should return to detail view on escape key", func() {
			intent.state.currentState = BurstStateDetailFacts

			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.state.currentState).To(Equal(BurstStateDetail))
		})

		It("should cancel on q key from facts view", func() {
			intent.state.currentState = BurstStateDetailFacts

			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
			Expect(intent.result.Status).To(Equal(Cancelled))
		})

		It("should handle no burst selected error when loading facts", func() {
			intent.state.selectedBurst = nil

			cmd := intent.loadFactsForBurst()
			msg := cmd().(BurstFactsLoadedMsg)

			Expect(msg.Error).To(HaveOccurred())
			Expect(msg.Error.Error()).To(ContainSubstring("no burst selected"))
		})

		It("should handle error in facts loaded message", func() {
			intent.state.currentState = BurstStateDetailFacts
			intent.state.loadingFacts = true

			msg := BurstFactsLoadedMsg{
				Error: errors.New("failed to load facts"),
			}

			intent.Update(msg)
			Expect(intent.state.loadingFacts).To(BeFalse())
			Expect(len(intent.state.burstFacts)).To(Equal(0))
		})
	})

	Describe("Detail View Enhanced Footer", func() {
		It("should show enhanced keyboard shortcuts in detail view", func() {
			intent.Init()
			intent.state.currentState = BurstStateDetail
			intent.state.selectedBurst = testBurst

			view := intent.View()
			Expect(view).To(ContainSubstring("e=events"))
			Expect(view).To(ContainSubstring("f=facts"))
		})
	})

	Describe("Pagination", func() {
		var (
			manyBurstsIntent *BurstManagementIntent
			manyBurstsRepo   *MockBurstRepository
		)

		BeforeEach(func() {
			// Create 35 bursts to span multiple pages (pageSize = 15)
			manyBurstsRepo = NewMockBurstRepository()
			for i := 0; i < 35; i++ {
				burst := &careerdom.Burst{
					ID:              "burst-" + fmt.Sprintf("%02d", i),
					Name:            fmt.Sprintf("Burst %02d", i+1),
					Description:     fmt.Sprintf("Burst description %d", i),
					EventIDs:        []string{"event-1"},
					CompetencyFocus: fmt.Sprintf("competency-%d", i%5),
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				}
				manyBurstsRepo.bursts = append(manyBurstsRepo.bursts, burst)
			}

			manyBurstsCtx := NewBurstManagementContext(nil, manyBurstsRepo, ctx)
			var err error
			manyBurstsIntent, err = NewBurstManagementIntent(manyBurstsCtx)
			Expect(err).NotTo(HaveOccurred())
			manyBurstsIntent.Init()
		})

		It("should display correct bursts on first page", func() {
			// Verify we're on page 1
			view := manyBurstsIntent.View()
			Expect(view).To(ContainSubstring("Page 1 of 3"))

			// Verify table shows first 15 bursts
			rows := manyBurstsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15))
		})

		It("should update table rows when navigating to next page", func() {
			// Navigate to page 2 using ctrl+d key (pgdn)
			manyBurstsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

			// Verify we're on page 2
			view := manyBurstsIntent.View()
			Expect(view).To(ContainSubstring("Page 2 of 3"))

			// Verify the selected index is now in the second page range
			Expect(manyBurstsIntent.state.selectedIndex).To(Equal(15))

			// FAILING TEST: Verify table shows the correct bursts for page 2
			// Currently the table shows ALL bursts (all 35 rows) instead of just the current page (15 rows)
			rows := manyBurstsIntent.table.Rows()
			Expect(len(rows)).To(Equal(15), "Table should show only 15 bursts for page 2, but shows %d bursts", len(rows))
		})

		It("should update table rows when navigating to last page", func() {
			// Navigate to page 3 using ctrl+d key twice
			manyBurstsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
			manyBurstsIntent.Update(tea.KeyMsg{Type: tea.KeyCtrlD})

			// Verify we're on page 3
			view := manyBurstsIntent.View()
			Expect(view).To(ContainSubstring("Page 3 of 3"))

			// Verify the selected index is now in the third page range
			Expect(manyBurstsIntent.state.selectedIndex).To(Equal(30))

			// FAILING TEST: Verify table shows the correct bursts for page 3
			// Currently the table shows ALL bursts (all 35 rows) instead of just the current page (5 rows)
			rows := manyBurstsIntent.table.Rows()
			Expect(len(rows)).To(Equal(5), "Table should show only 5 bursts for page 3, but shows %d bursts", len(rows))
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
