package factmanagement_test

import (
	"context"
	"errors"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	careerrepo "github.com/baphled/kariya/internal/repository/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/intents/factmanagement"
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
	Describe("IntentValidator", func() {
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
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				Expect(intentCtx).NotTo(BeNil())
				Expect(intentCtx.Facts).NotTo(BeNil())
				Expect(intentCtx.Facts).To(BeEmpty())
			})

			It("should set default sort options", func() {
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				Expect(intentCtx.SortBy).To(Equal("date"))
				Expect(intentCtx.SortOrder).To(Equal("desc"))
			})

			It("should set default quality range", func() {
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				Expect(intentCtx.MinQuality).To(Equal(0.0))
				Expect(intentCtx.MaxQuality).To(Equal(1.0))
			})

			It("should initialize empty maps", func() {
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				Expect(intentCtx.FormErrors).NotTo(BeNil())
				Expect(intentCtx.ExpandedRows).NotTo(BeNil())
			})
		})

		Describe("Validate", func() {
			Context("when all fields are valid", func() {
				It("should return no error", func() {
					intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("when Facts is nil", func() {
				It("should initialize to empty slice", func() {
					intentCtx := &factmanagement.IntentValidator{
						Facts: nil,
					}
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Facts).NotTo(BeNil())
				})
			})

			Context("when FormErrors is nil", func() {
				It("should initialize to empty map", func() {
					intentCtx := &factmanagement.IntentValidator{
						FormErrors: nil,
					}
					err := intentCtx.Validate()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.FormErrors).NotTo(BeNil())
				})
			})

			Context("when ExpandedRows is nil", func() {
				It("should initialize to empty map", func() {
					intentCtx := &factmanagement.IntentValidator{
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
				testFact = fixtures.Fact("fact-1", "event-1")
				testFact.Text = "Test fact"
			})

			Context("with nil repository", func() {
				It("should return service not available error", func() {
					intentCtx := &factmanagement.IntentValidator{
						FactRepository: nil,
					}
					err := intentCtx.LoadFacts()
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(factmanagement.ErrServiceNotAvailable))
				})
			})

			Context("with valid repository", func() {
				It("should load facts successfully", func() {
					mockRepo.facts = []*career.Fact{testFact}
					intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
					err := intentCtx.LoadFacts()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Facts).To(HaveLen(1))
				})

				It("should set total facts count", func() {
					mockRepo.facts = []*career.Fact{testFact}
					intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
					_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
					Expect(intentCtx.TotalFacts).To(Equal(1))
				})

				It("should select first fact when available", func() {
					mockRepo.facts = []*career.Fact{testFact}
					intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
					_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
					Expect(intentCtx.SelectedFact).To(Equal(testFact))
					Expect(intentCtx.SelectedFactIndex).To(Equal(0))
				})

				It("should handle empty fact list", func() {
					intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
					err := intentCtx.LoadFacts()
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Facts).To(BeEmpty())
					Expect(intentCtx.SelectedFactIndex).To(Equal(-1))
				})
			})

			Context("when repository returns error", func() {
				It("should propagate the error", func() {
					mockRepo.listErr = errors.New("database error")
					intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
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
				for i := range 25 {
					testFacts[i] = fixtures.FactWith("fact-"+string(rune('a'+i)), "Test fact "+string(rune('a'+i)))
				}
			})

			It("should return empty slice when no facts", func() {
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(BeEmpty())
			})

			It("should return first page of facts", func() {
				mockRepo.facts = testFacts
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(HaveLen(20))
			})

			It("should return partial page at end", func() {
				mockRepo.facts = testFacts
				intentCtx := factmanagement.NewIntentValidator(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
				intentCtx.CurrentPage = 1
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(HaveLen(5))
			})
		})

		Describe("Form Error Handling", func() {
			var intentCtx *factmanagement.IntentValidator

			BeforeEach(func() {
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
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
			var intentCtx *factmanagement.IntentValidator

			BeforeEach(func() {
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
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
				intentCtx *factmanagement.IntentValidator
				testFact  *career.Fact
			)

			BeforeEach(func() {
				testFact = fixtures.Fact("fact-1", "event-1")
				testFact.Text = "Test fact"
				mockRepo.facts = []*career.Fact{testFact}
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
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
						Expect(err).To(Equal(factmanagement.ErrInvalidState))
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
				intentCtx *factmanagement.IntentValidator
				testFact  *career.Fact
			)

			BeforeEach(func() {
				testFact = fixtures.FactWith("fact-1", "Test fact")
				mockRepo.facts = []*career.Fact{testFact}
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
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
					Expect(err).To(Equal(factmanagement.ErrServiceNotAvailable))
				})
			})
		})

		Describe("Selection", func() {
			var (
				intentCtx *factmanagement.IntentValidator
				testFacts []*career.Fact
			)

			BeforeEach(func() {
				testFacts = []*career.Fact{
					fixtures.FactWith("fact-1", "First fact"),
					fixtures.FactWith("fact-2", "Second fact"),
				}
				mockRepo.facts = testFacts
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
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

		Describe("CreateFact", func() {
			var (
				intentCtx *factmanagement.IntentValidator
				testFact  *career.Fact
			)

			BeforeEach(func() {
				testFact = fixtures.FactWith("new-fact", "New fact text")
				testFact.CompetencyCategories = []string{"technical"}
				testFact.RoleFit = "senior_ic"
				testFact.AudienceRelevance = []string{"hiring_manager"}
				testFact.SourceEventID = "event-1"
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
			})

			Context("with valid fact", func() {
				It("should create fact successfully", func() {
					err := intentCtx.CreateFact(testFact)
					Expect(err).NotTo(HaveOccurred())
					Expect(intentCtx.Facts).To(HaveLen(1))
					Expect(intentCtx.TotalFacts).To(Equal(1))
				})

				It("should append to facts list", func() {
					_ = intentCtx.CreateFact(testFact) //nolint:errcheck // test setup
					Expect(intentCtx.Facts[0]).To(Equal(testFact))
				})

				It("should update total facts count", func() {
					_ = intentCtx.CreateFact(testFact) //nolint:errcheck // test setup
					Expect(intentCtx.TotalFacts).To(Equal(1))
				})
			})

			Context("with nil repository", func() {
				It("should return service not available error", func() {
					intentCtx.FactRepository = nil
					err := intentCtx.CreateFact(testFact)
					Expect(err).To(Equal(factmanagement.ErrServiceNotAvailable))
				})
			})

			Context("when repository returns error", func() {
				It("should propagate the error", func() {
					mockRepo.createErr = errors.New("create failed")
					err := intentCtx.CreateFact(testFact)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("create failed"))
				})
			})

			Context("with invalid fact", func() {
				It("should return validation error", func() {
					invalidFact := fixtures.Fact("invalid", "source-event")
					invalidFact.Text = ""
					err := intentCtx.CreateFact(invalidFact)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("UpdateFact", func() {
			var (
				intentCtx *factmanagement.IntentValidator
				testFact  *career.Fact
			)

			BeforeEach(func() {
				testFact = fixtures.FactWith("fact-1", "Original text")
				testFact.CompetencyCategories = []string{"technical"}
				testFact.RoleFit = "senior_ic"
				testFact.AudienceRelevance = []string{"hiring_manager"}
				testFact.SourceEventID = "event-1"
				mockRepo.facts = []*career.Fact{testFact}
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
			})

			Context("with valid fact", func() {
				It("should update fact successfully", func() {
					testFact.Text = "Updated text"
					err := intentCtx.UpdateFact(testFact)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should replace fact in list", func() {
					testFact.Text = "Updated text"
					_ = intentCtx.UpdateFact(testFact) //nolint:errcheck // test setup
					Expect(intentCtx.Facts[0].Text).To(Equal("Updated text"))
				})

				It("should update multiple facts", func() {
					fact2 := fixtures.FactWith("fact-2", "Second fact")
					fact2.CompetencyCategories = []string{"technical"}
					fact2.RoleFit = "senior_ic"
					fact2.AudienceRelevance = []string{"hiring_manager"}
					fact2.SourceEventID = "event-1"
					mockRepo.facts = append(mockRepo.facts, fact2)
					intentCtx.Facts = append(intentCtx.Facts, fact2)
					testFact.Text = "Updated first"
					_ = intentCtx.UpdateFact(testFact) //nolint:errcheck // test setup
					Expect(intentCtx.Facts[0].Text).To(Equal("Updated first"))
					Expect(intentCtx.Facts[1].Text).To(Equal("Second fact"))
				})
			})

			Context("with nil repository", func() {
				It("should return service not available error", func() {
					intentCtx.FactRepository = nil
					err := intentCtx.UpdateFact(testFact)
					Expect(err).To(Equal(factmanagement.ErrServiceNotAvailable))
				})
			})

			Context("when repository returns error", func() {
				It("should propagate the error", func() {
					mockRepo.updateErr = errors.New("update failed")
					err := intentCtx.UpdateFact(testFact)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("update failed"))
				})
			})

			Context("with invalid fact", func() {
				It("should return validation error", func() {
					invalidFact := fixtures.Fact("fact-1", "source-event")
					invalidFact.Text = ""
					err := intentCtx.UpdateFact(invalidFact)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Describe("GetPageFacts Edge Cases", func() {
			var (
				intentCtx *factmanagement.IntentValidator
				testFacts []*career.Fact
			)

			BeforeEach(func() {
				testFacts = make([]*career.Fact, 35)
				for i := range 35 {
					testFacts[i] = fixtures.FactWith("fact-"+string(rune('a'+(i%26))), "Fact "+string(rune('a'+(i%26))))
					testFacts[i].CompetencyCategories = []string{"technical"}
					testFacts[i].RoleFit = "senior_ic"
					testFacts[i].AudienceRelevance = []string{"hiring_manager"}
					testFacts[i].SourceEventID = "event-1"
				}
				mockRepo.facts = testFacts
				intentCtx = factmanagement.NewIntentValidator(ctx, mockRepo)
				_ = intentCtx.LoadFacts() //nolint:errcheck // test setup
			})

			It("should return correct page when at boundary", func() {
				intentCtx.CurrentPage = 1
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(HaveLen(15))
			})

			It("should return empty when page is beyond bounds", func() {
				intentCtx.CurrentPage = 10
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(BeEmpty())
			})

			It("should handle page size changes", func() {
				intentCtx.PageSize = 10
				pageFacts := intentCtx.GetPageFacts()
				Expect(pageFacts).To(HaveLen(10))
			})
		})
	})
})
