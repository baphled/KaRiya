package fact_management_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/intents/fact_management"
	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
)

// MockFactRepository implements FactRepository for testing.
type MockFactRepository struct {
	facts     []*career.Fact
	createErr error
	updateErr error
	deleteErr error
	listErr   error
}

func NewMockFactRepository() *MockFactRepository {
	return &MockFactRepository{
		facts: make([]*career.Fact, 0),
	}
}

func (m *MockFactRepository) Create(_ context.Context, fact *career.Fact) error {
	if m.createErr != nil {
		return m.createErr
	}
	if fact.ID == "" {
		fact.ID = "fact-" + time.Now().Format("20060102150405")
	}
	m.facts = append(m.facts, fact)
	return nil
}

func (m *MockFactRepository) GetByID(_ context.Context, id string) (*career.Fact, error) {
	for _, f := range m.facts {
		if f.ID == id {
			return f, nil
		}
	}
	return nil, errors.New("not found")
}

func (m *MockFactRepository) Update(_ context.Context, fact *career.Fact) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	for i, f := range m.facts {
		if f.ID == fact.ID {
			m.facts[i] = fact
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockFactRepository) Delete(_ context.Context, id string) error {
	if m.deleteErr != nil {
		return m.deleteErr
	}
	for i, f := range m.facts {
		if f.ID == id {
			m.facts = append(m.facts[:i], m.facts[i+1:]...)
			return nil
		}
	}
	return errors.New("not found")
}

func (m *MockFactRepository) List(_ context.Context, _ careerrepo.FactListFilters) ([]*career.Fact, error) {
	if m.listErr != nil {
		return nil, m.listErr
	}
	return m.facts, nil
}

func (m *MockFactRepository) Count(_ context.Context, _ careerrepo.FactListFilters) (int, error) {
	return len(m.facts), nil
}

func (m *MockFactRepository) GetBySourceEventID(_ context.Context, _ string) ([]*career.Fact, error) {
	return make([]*career.Fact, 0), nil
}

func (m *MockFactRepository) GetBySourceBurstID(_ context.Context, _ string) ([]*career.Fact, error) {
	return make([]*career.Fact, 0), nil
}

var _ = Describe("Context", func() {
	Describe("IntentContext", func() {
		var (
			ctx      context.Context
			mockRepo *MockFactRepository
		)

		BeforeEach(func() {
			ctx = context.Background()
			mockRepo = NewMockFactRepository()
		})

		Describe("Construction", func() {
			It("should create a context with default values", func() {
				intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx).NotTo(BeNil())
				Expect(intentCtx.Facts).NotTo(BeNil())
				Expect(intentCtx.Facts).To(BeEmpty())
			})

			It("should set default sort options", func() {
				intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx.SortBy).To(Equal("date"))
				Expect(intentCtx.SortOrder).To(Equal("desc"))
			})

			It("should set default quality range", func() {
				intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx.MinQuality).To(Equal(0.0))
				Expect(intentCtx.MaxQuality).To(Equal(1.0))
			})

			It("should initialize empty maps", func() {
				intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
				Expect(intentCtx.FormErrors).NotTo(BeNil())
				Expect(intentCtx.ExpandedRows).NotTo(BeNil())
			})
		})

		Describe("Validate", func() {
			Context("when all fields are valid", func() {
				It("should return no error", func() {
					intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when Facts is nil", func() {
				It("should initialize to empty slice", func() {
					intentCtx := &fact_management.IntentContext{
						Facts: nil,
					}
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Facts).NotTo(BeNil())
				})
			})

			Context("when FormErrors is nil", func() {
				It("should initialize to empty map", func() {
					intentCtx := &fact_management.IntentContext{
						FormErrors: nil,
					}
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.FormErrors).NotTo(BeNil())
				})
			})

			Context("when ExpandedRows is nil", func() {
				It("should initialize to empty map", func() {
					intentCtx := &fact_management.IntentContext{
						ExpandedRows: nil,
					}
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.ExpandedRows).NotTo(BeNil())
				})
			})
		})

		Describe("LoadFacts", func() {
			var testFact *career.Fact

			BeforeEach(func() {
				testFact = &career.Fact{
					ID:                   "fact-1",
					Text:                 "Test fact",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "high",
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}
			})

			Context("with nil repository", func() {
				It("should return service not available error", func() {
					intentCtx := &fact_management.IntentContext{
						FactRepository: nil,
						Context:        ctx,
					}
					err := intentCtx.LoadFacts()
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(fact_management.ErrServiceNotAvailable))
				})
			})

			Context("with valid repository", func() {
				It("should load facts successfully", func() {
					mockRepo.facts = []*career.Fact{testFact}
					intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
					err := intentCtx.LoadFacts()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Facts).To(HaveLen(1))
				})

				It("should set total facts count", func() {
					mockRepo.facts = []*career.Fact{testFact}
					intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
					_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
					Expect(intentCtx.TotalFacts).To(Equal(1))
				})

				It("should select first fact when available", func() {
					mockRepo.facts = []*career.Fact{testFact}
					intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
					_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
					Expect(intentCtx.SelectedFact).To(Equal(testFact))
					Expect(intentCtx.SelectedFactIndex).To(Equal(0))
				})

				It("should handle empty fact list", func() {
					intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
					err := intentCtx.LoadFacts()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Facts).To(BeEmpty())
					Expect(intentCtx.SelectedFactIndex).To(Equal(-1))
				})
			})

			Context("when repository returns error", func() {
				It("should propagate the error", func() {
					mockRepo.listErr = errors.New("database error")
					intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
					err := intentCtx.LoadFacts()
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("database error"))
				})
			})
		})

		Describe("GetPageFacts", func() {
			var testFacts []*career.Fact

			BeforeEach(func() {
				testFacts = make([]*career.Fact, 25)
				for i := 0; i < 25; i++ {
					testFacts[i] = &career.Fact{
						ID:   "fact-" + string(rune('a'+i)),
						Text: "Test fact " + string(rune('a'+i)),
					}
				}
			})

			It("should return empty slice when no facts", func() {
				intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(BeEmpty())
			})

			It("should return first page of facts", func() {
				mockRepo.facts = testFacts
				intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(HaveLen(20))
			})

			It("should return partial page at end", func() {
				mockRepo.facts = testFacts
				intentCtx := fact_management.NewIntentContext(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
				intentCtx.CurrentPage = 1
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(HaveLen(5))
			})
		})

		Describe("Form Error Handling", func() {
			var intentCtx *fact_management.IntentContext

			BeforeEach(func() {
				intentCtx = fact_management.NewIntentContext(ctx, mockRepo)
			})

			It("should set form error", func() {
				intentCtx.SetFormError("text", "Text is required")
				Expect(intentCtx.HasFormErrors()).To(BeTrue())
			})

			It("should clear form errors", func() {
				intentCtx.SetFormError("text", "Text is required")
				intentCtx.ClearFormErrors()
				Expect(intentCtx.HasFormErrors()).To(BeFalse())
			})

			It("should report no errors initially", func() {
				Expect(intentCtx.HasFormErrors()).To(BeFalse())
			})
		})

		Describe("Row Expansion", func() {
			var intentCtx *fact_management.IntentContext

			BeforeEach(func() {
				intentCtx = fact_management.NewIntentContext(ctx, mockRepo)
			})

			It("should toggle row expansion", func() {
				intentCtx.ToggleRowExpansion(0)
				Expect(intentCtx.IsRowExpanded(0)).To(BeTrue())
			})

			It("should toggle off when expanded", func() {
				intentCtx.ToggleRowExpansion(0)
				intentCtx.ToggleRowExpansion(0)
				Expect(intentCtx.IsRowExpanded(0)).To(BeFalse())
			})
		})

		Describe("Edit Operations", func() {
			var (
				intentCtx *fact_management.IntentContext
				testFact  *career.Fact
			)

			BeforeEach(func() {
				testFact = &career.Fact{
					ID:                   "fact-1",
					Text:                 "Test fact",
					CompetencyCategories: []string{"technical"},
					StrengthSignal:       "high",
					CreatedAt:            time.Now(),
					UpdatedAt:            time.Now(),
				}
				mockRepo.facts = []*career.Fact{testFact}
				intentCtx = fact_management.NewIntentContext(ctx, mockRepo)
			})

			Describe("StartNewFact", func() {
				It("should create new editing fact", func() {
					intentCtx.StartNewFact()
					Expect(intentCtx.EditingFact).NotTo(BeNil())
					Expect(intentCtx.IsNewFact).To(BeTrue())
				})

				It("should clear form errors", func() {
					intentCtx.SetFormError("text", "error")
					intentCtx.StartNewFact()
					Expect(intentCtx.HasFormErrors()).To(BeFalse())
				})
			})

			Describe("StartEditFact", func() {
				It("should copy fact for editing", func() {
					intentCtx.StartEditFact(testFact)
					Expect(intentCtx.EditingFact).NotTo(BeNil())
					Expect(intentCtx.EditingFact.ID).To(Equal(testFact.ID))
					Expect(intentCtx.IsNewFact).To(BeFalse())
				})

				It("should not modify original fact", func() {
					intentCtx.StartEditFact(testFact)
					intentCtx.EditingFact.Text = "Modified"
					Expect(testFact.Text).To(Equal("Test fact"))
				})
			})

			Describe("CancelEdit", func() {
				It("should clear editing fact", func() {
					intentCtx.StartNewFact()
					intentCtx.CancelEdit()
					Expect(intentCtx.EditingFact).To(BeNil())
				})

				It("should reset new fact flag", func() {
					intentCtx.StartNewFact()
					intentCtx.CancelEdit()
					Expect(intentCtx.IsNewFact).To(BeFalse())
				})
			})

			Describe("SaveEdit", func() {
				Context("with no editing fact", func() {
					It("should return invalid state error", func() {
						err := intentCtx.SaveEdit()
						Expect(err).To(Equal(fact_management.ErrInvalidState))
					})
				})

				Context("with new fact", func() {
					It("should create fact in repository", func() {
						intentCtx.StartNewFact()
						intentCtx.EditingFact.ID = "new-fact-id"
						intentCtx.EditingFact.Text = "New fact text"
						intentCtx.EditingFact.CompetencyCategories = []string{"technical"}
						intentCtx.EditingFact.RoleFit = "senior_ic"
						intentCtx.EditingFact.AudienceRelevance = []string{"hiring_manager"}
						intentCtx.EditingFact.SourceEventID = "event-1"
						err := intentCtx.SaveEdit()
						Expect(err).NotTo(HaveOccurred())
						Expect(mockRepo.facts).To(HaveLen(2))
					})
				})

				Context("with existing fact", func() {
					It("should update fact in repository", func() {
						intentCtx.StartEditFact(testFact)
						intentCtx.EditingFact.Text = "Updated text"
						intentCtx.EditingFact.RoleFit = "senior_ic"
						intentCtx.EditingFact.AudienceRelevance = []string{"hiring_manager"}
						intentCtx.EditingFact.SourceEventID = "event-1"
						err := intentCtx.SaveEdit()
						Expect(err).NotTo(HaveOccurred())
					})
				})
			})
		})

		Describe("Delete Operations", func() {
			var (
				intentCtx *fact_management.IntentContext
				testFact  *career.Fact
			)

			BeforeEach(func() {
				testFact = &career.Fact{
					ID:   "fact-1",
					Text: "Test fact",
				}
				mockRepo.facts = []*career.Fact{testFact}
				intentCtx = fact_management.NewIntentContext(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
			})

			It("should delete fact from repository", func() {
				err := intentCtx.DeleteFact(testFact.ID)
				Expect(err).NotTo(HaveOccurred())
				Expect(intentCtx.Facts).To(BeEmpty())
			})

			It("should update total count after delete", func() {
				_ = intentCtx.DeleteFact(testFact.ID) //nolint:errcheck // test setup
				Expect(intentCtx.TotalFacts).To(Equal(0))
			})

			It("should clear selection after delete", func() {
				_ = intentCtx.DeleteFact(testFact.ID) //nolint:errcheck // test setup
				Expect(intentCtx.SelectedFact).To(BeNil())
				Expect(intentCtx.SelectedFactIndex).To(Equal(-1))
			})

			Context("with nil repository", func() {
				It("should return service not available error", func() {
					intentCtx.FactRepository = nil
					err := intentCtx.DeleteFact(testFact.ID)
					Expect(err).To(Equal(fact_management.ErrServiceNotAvailable))
				})
			})
		})

		Describe("Selection", func() {
			var (
				intentCtx *fact_management.IntentContext
				testFacts []*career.Fact
			)

			BeforeEach(func() {
				testFacts = []*career.Fact{
					{ID: "fact-1", Text: "First fact"},
					{ID: "fact-2", Text: "Second fact"},
				}
				mockRepo.facts = testFacts
				intentCtx = fact_management.NewIntentContext(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
			})

			It("should select fact by index", func() {
				intentCtx.SelectFact(1)
				Expect(intentCtx.SelectedFact).To(Equal(testFacts[1]))
				Expect(intentCtx.SelectedFactIndex).To(Equal(1))
			})

			It("should return selected fact", func() {
				intentCtx.SelectFact(0)
				Expect(intentCtx.GetSelectedFact()).To(Equal(testFacts[0]))
			})

			It("should ignore invalid index", func() {
				intentCtx.SelectFact(100)
				Expect(intentCtx.SelectedFact).To(Equal(testFacts[0]))
			})
		})
	})
})
