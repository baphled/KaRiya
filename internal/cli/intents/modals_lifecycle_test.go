package intents

import (
	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Modal Lifecycle Tests", func() {
	// ============================================================================
	// EditMetadataModal Lifecycle Tests
	// ============================================================================
	Describe("EditMetadataModal", func() {
		var modal *EditMetadataModal

		BeforeEach(func() {
			modal = NewEditMetadataModal(
				"Test Company",
				"Test Project",
				[]string{"leadership", "technical"},
				[]string{"achievement"},
			)
		})

		Describe("Update", func() {
			It("handles WindowSizeMsg and updates dimensions", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				modal.Update(msg)

				Expect(modal.width).To(Equal(120))
				Expect(modal.height).To(Equal(40))
			})

			It("preserves form state after resize", func() {
				// Get initial form view
				initialView := modal.View()
				Expect(initialView).NotTo(BeEmpty())

				// Resize
				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				modal.Update(msg)

				// Form should still render
				Expect(modal.View()).NotTo(BeEmpty())
			})

			It("returns nil command for WindowSizeMsg", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				cmd := modal.Update(msg)

				Expect(cmd).To(BeNil())
			})
		})

		Describe("View", func() {
			It("renders form content when active", func() {
				view := modal.View()

				Expect(view).To(ContainSubstring("Edit Event Metadata"))
			})

			It("returns empty string when result is accepted", func() {
				modal.result = &ModalEditResult[*MetadataSnapshot]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				Expect(modal.View()).To(BeEmpty())
			})

			It("still renders when result is cancelled", func() {
				modal.result = &ModalEditResult[*MetadataSnapshot]{
					Original: modal.original,
					Modified: modal.original,
					Accepted: false,
				}

				// When result exists but Accepted is false, view still renders
				// The condition is `result != nil && result.Accepted` so cancelled results still show form
				// This allows showing the cancellation state before closing
				Expect(modal.View()).NotTo(BeEmpty())
			})
		})

		Describe("Result", func() {
			It("returns nil initially", func() {
				Expect(modal.Result()).To(BeNil())
			})

			It("returns result after completion", func() {
				modal.result = &ModalEditResult[*MetadataSnapshot]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Accepted).To(BeTrue())
			})
		})

		Describe("IsComplete", func() {
			It("returns false initially", func() {
				Expect(modal.IsComplete()).To(BeFalse())
			})

			It("returns true when result is set", func() {
				modal.result = &ModalEditResult[*MetadataSnapshot]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				Expect(modal.IsComplete()).To(BeTrue())
			})

			It("returns true for cancelled result", func() {
				modal.result = &ModalEditResult[*MetadataSnapshot]{
					Original: modal.original,
					Modified: modal.original,
					Accepted: false,
				}

				Expect(modal.IsComplete()).To(BeTrue())
			})
		})

		Describe("computeChanges", func() {
			It("detects company change", func() {
				modal.modified = &MetadataSnapshot{
					Company:    "New Company",
					Project:    modal.original.Project,
					Tags:       modal.original.Tags,
					Categories: modal.original.Categories,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("company"))
				Expect(changes["company"]).To(Equal("New Company"))
			})

			It("detects project change", func() {
				modal.modified = &MetadataSnapshot{
					Company:    modal.original.Company,
					Project:    "New Project",
					Tags:       modal.original.Tags,
					Categories: modal.original.Categories,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("project"))
				Expect(changes["project"]).To(Equal("New Project"))
			})

			It("detects tag changes", func() {
				modal.modified = &MetadataSnapshot{
					Company:    modal.original.Company,
					Project:    modal.original.Project,
					Tags:       []string{"new-tag"},
					Categories: modal.original.Categories,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("tags"))
			})

			It("detects category changes", func() {
				modal.modified = &MetadataSnapshot{
					Company:    modal.original.Company,
					Project:    modal.original.Project,
					Tags:       modal.original.Tags,
					Categories: []string{"new-category"},
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("categories"))
			})

			It("returns empty map when nothing changed", func() {
				modal.modified = copyMetadataSnapshot(modal.original)

				changes := modal.computeChanges()
				Expect(changes).To(BeEmpty())
			})

			It("detects multiple changes", func() {
				modal.modified = &MetadataSnapshot{
					Company:    "New Company",
					Project:    "New Project",
					Tags:       []string{"new-tag"},
					Categories: []string{"new-category"},
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveLen(4))
				Expect(changes).To(HaveKey("company"))
				Expect(changes).To(HaveKey("project"))
				Expect(changes).To(HaveKey("tags"))
				Expect(changes).To(HaveKey("categories"))
			})
		})

		Describe("createCancelledResult", func() {
			It("sets result with original data", func() {
				modal.createCancelledResult()

				Expect(modal.result).NotTo(BeNil())
				Expect(modal.result.Accepted).To(BeFalse())
				Expect(modal.result.Modified).To(Equal(modal.original))
			})

			It("has empty changes map", func() {
				modal.createCancelledResult()

				Expect(modal.result.Changes).To(BeEmpty())
			})
		})

		Describe("syncModified", func() {
			It("syncs form field values to modified", func() {
				// Set form field values directly
				*modal.company = "Synced Company"
				*modal.project = "Synced Project"
				modal.tags = []string{"synced-tag"}
				modal.categories = []string{"synced-category"}

				modal.syncModified()

				Expect(modal.modified.Company).To(Equal("Synced Company"))
				Expect(modal.modified.Project).To(Equal("Synced Project"))
				Expect(modal.modified.Tags).To(ContainElement("synced-tag"))
				Expect(modal.modified.Categories).To(ContainElement("synced-category"))
			})
		})
	})

	// ============================================================================
	// EditBurstModal Lifecycle Tests
	// ============================================================================
	Describe("EditBurstModal", func() {
		var (
			modal *EditBurstModal
			burst *career.Burst
		)

		BeforeEach(func() {
			burst = &career.Burst{
				ID:          "burst-1",
				Name:        "Test Burst",
				Description: "Test Description",
			}
			modal = NewEditBurstModal(burst)
		})

		Describe("Update", func() {
			It("handles WindowSizeMsg and updates dimensions", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				modal.Update(msg)

				Expect(modal.width).To(Equal(120))
				Expect(modal.height).To(Equal(40))
			})

			It("returns nil command for WindowSizeMsg", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				cmd := modal.Update(msg)

				Expect(cmd).To(BeNil())
			})
		})

		Describe("View", func() {
			It("renders form content when active", func() {
				view := modal.View()

				Expect(view).To(ContainSubstring("Edit Burst"))
			})

			It("returns empty string when result is accepted", func() {
				modal.result = &ModalEditResult[*career.Burst]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				Expect(modal.View()).To(BeEmpty())
			})
		})

		Describe("Result", func() {
			It("returns nil initially", func() {
				Expect(modal.Result()).To(BeNil())
			})

			It("returns result after completion", func() {
				modal.result = &ModalEditResult[*career.Burst]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Accepted).To(BeTrue())
			})
		})

		Describe("IsComplete", func() {
			It("returns false initially", func() {
				Expect(modal.IsComplete()).To(BeFalse())
			})

			It("returns true when result is set", func() {
				modal.result = &ModalEditResult[*career.Burst]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				Expect(modal.IsComplete()).To(BeTrue())
			})
		})

		Describe("SetTestResult", func() {
			It("sets result directly for testing", func() {
				testResult := &ModalEditResult[*career.Burst]{
					Original: modal.original,
					Modified: &career.Burst{
						ID:          "burst-1",
						Name:        "Modified Burst",
						Description: "Modified Description",
					},
					Accepted: true,
				}

				modal.SetTestResult(testResult)

				Expect(modal.result).To(Equal(testResult))
				Expect(modal.modified.Name).To(Equal("Modified Burst"))
			})

			It("handles nil result", func() {
				modal.SetTestResult(nil)

				Expect(modal.result).To(BeNil())
			})
		})

		Describe("computeChanges", func() {
			It("detects name change", func() {
				modal.modified = &career.Burst{
					ID:          modal.original.ID,
					Name:        "New Name",
					Description: modal.original.Description,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("name"))
				Expect(changes["name"]).To(Equal("New Name"))
			})

			It("detects description change", func() {
				modal.modified = &career.Burst{
					ID:          modal.original.ID,
					Name:        modal.original.Name,
					Description: "New Description",
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("description"))
				Expect(changes["description"]).To(Equal("New Description"))
			})

			It("returns empty map when nothing changed", func() {
				modal.modified = &career.Burst{
					ID:          modal.original.ID,
					Name:        modal.original.Name,
					Description: modal.original.Description,
				}

				changes := modal.computeChanges()
				Expect(changes).To(BeEmpty())
			})

			It("detects both changes", func() {
				modal.modified = &career.Burst{
					ID:          modal.original.ID,
					Name:        "New Name",
					Description: "New Description",
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveLen(2))
			})
		})

		Describe("createCancelledResult", func() {
			It("sets result with original data", func() {
				modal.createCancelledResult()

				Expect(modal.result).NotTo(BeNil())
				Expect(modal.result.Accepted).To(BeFalse())
				Expect(modal.result.Modified).To(Equal(modal.original))
			})

			It("has empty changes map", func() {
				modal.createCancelledResult()

				Expect(modal.result.Changes).To(BeEmpty())
			})
		})

		Describe("syncModified", func() {
			It("syncs form data to modified burst", func() {
				modal.formData.Name = "Synced Name"
				modal.formData.Description = "Synced Description"

				modal.syncModified()

				Expect(modal.modified.Name).To(Equal("Synced Name"))
				Expect(modal.modified.Description).To(Equal("Synced Description"))
				// ID should be preserved from original
				Expect(modal.modified.ID).To(Equal(modal.original.ID))
			})
		})
	})

	// ============================================================================
	// EditFactModal Lifecycle Tests
	// ============================================================================
	Describe("EditFactModal", func() {
		var (
			modal *EditFactModal
			fact  *career.Fact
		)

		BeforeEach(func() {
			fact = &career.Fact{
				ID:                   "fact-1",
				Text:                 "Test fact text",
				CompetencyCategories: []string{"leadership"},
				RoleFit:              career.RoleFitPrincipal,
				AudienceRelevance:    []string{"enterprise"},
				StrengthSignal:       "high",
				SourceEventID:        "event-1",
			}
			modal = NewEditFactModal(fact)
		})

		Describe("Update", func() {
			It("handles WindowSizeMsg and updates dimensions", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				modal.Update(msg)

				Expect(modal.width).To(Equal(120))
				Expect(modal.height).To(Equal(40))
			})

			It("returns nil command for WindowSizeMsg", func() {
				msg := tea.WindowSizeMsg{Width: 100, Height: 30}
				cmd := modal.Update(msg)

				Expect(cmd).To(BeNil())
			})
		})

		Describe("View", func() {
			It("renders form content when active", func() {
				view := modal.View()

				Expect(view).To(ContainSubstring("Edit Fact"))
			})

			It("returns empty string when result is accepted", func() {
				modal.result = &ModalEditResult[*career.Fact]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				Expect(modal.View()).To(BeEmpty())
			})
		})

		Describe("Result", func() {
			It("returns nil initially", func() {
				Expect(modal.Result()).To(BeNil())
			})

			It("returns result after completion", func() {
				modal.result = &ModalEditResult[*career.Fact]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				result := modal.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Accepted).To(BeTrue())
			})
		})

		Describe("IsComplete", func() {
			It("returns false initially", func() {
				Expect(modal.IsComplete()).To(BeFalse())
			})

			It("returns true when result is set", func() {
				modal.result = &ModalEditResult[*career.Fact]{
					Original: modal.original,
					Modified: modal.modified,
					Accepted: true,
				}

				Expect(modal.IsComplete()).To(BeTrue())
			})
		})

		Describe("computeChanges", func() {
			It("detects text change", func() {
				modal.modified = &career.Fact{
					ID:                   modal.original.ID,
					Text:                 "New text",
					CompetencyCategories: modal.original.CompetencyCategories,
					RoleFit:              modal.original.RoleFit,
					AudienceRelevance:    modal.original.AudienceRelevance,
					StrengthSignal:       modal.original.StrengthSignal,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("text"))
				Expect(changes["text"]).To(Equal("New text"))
			})

			It("detects competency categories change", func() {
				modal.modified = &career.Fact{
					ID:                   modal.original.ID,
					Text:                 modal.original.Text,
					CompetencyCategories: []string{"technical", "innovation"},
					RoleFit:              modal.original.RoleFit,
					AudienceRelevance:    modal.original.AudienceRelevance,
					StrengthSignal:       modal.original.StrengthSignal,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("competency_categories"))
			})

			It("detects role fit change", func() {
				modal.modified = &career.Fact{
					ID:                   modal.original.ID,
					Text:                 modal.original.Text,
					CompetencyCategories: modal.original.CompetencyCategories,
					RoleFit:              career.RoleFitEM,
					AudienceRelevance:    modal.original.AudienceRelevance,
					StrengthSignal:       modal.original.StrengthSignal,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("role_fit"))
			})

			It("detects audience relevance change", func() {
				modal.modified = &career.Fact{
					ID:                   modal.original.ID,
					Text:                 modal.original.Text,
					CompetencyCategories: modal.original.CompetencyCategories,
					RoleFit:              modal.original.RoleFit,
					AudienceRelevance:    []string{"startup"},
					StrengthSignal:       modal.original.StrengthSignal,
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("audience_relevance"))
			})

			It("detects strength signal change", func() {
				modal.modified = &career.Fact{
					ID:                   modal.original.ID,
					Text:                 modal.original.Text,
					CompetencyCategories: modal.original.CompetencyCategories,
					RoleFit:              modal.original.RoleFit,
					AudienceRelevance:    modal.original.AudienceRelevance,
					StrengthSignal:       "medium",
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveKey("strength_signal"))
			})

			It("returns empty map when nothing changed", func() {
				modal.modified = &career.Fact{
					ID:                   modal.original.ID,
					Text:                 modal.original.Text,
					CompetencyCategories: modal.original.CompetencyCategories,
					RoleFit:              modal.original.RoleFit,
					AudienceRelevance:    modal.original.AudienceRelevance,
					StrengthSignal:       modal.original.StrengthSignal,
				}

				changes := modal.computeChanges()
				Expect(changes).To(BeEmpty())
			})

			It("detects multiple changes", func() {
				modal.modified = &career.Fact{
					ID:                   modal.original.ID,
					Text:                 "New text",
					CompetencyCategories: []string{"new-category"},
					RoleFit:              career.RoleFitStaff,
					AudienceRelevance:    []string{"startup"},
					StrengthSignal:       "low",
				}

				changes := modal.computeChanges()
				Expect(changes).To(HaveLen(5))
			})
		})

		Describe("createCancelledResult", func() {
			It("sets result with original data", func() {
				modal.createCancelledResult()

				Expect(modal.result).NotTo(BeNil())
				Expect(modal.result.Accepted).To(BeFalse())
				Expect(modal.result.Modified).To(Equal(modal.original))
			})

			It("has empty changes map", func() {
				modal.createCancelledResult()

				Expect(modal.result.Changes).To(BeEmpty())
			})
		})

		Describe("syncModified", func() {
			It("syncs form data to modified fact", func() {
				modal.formData.Text = "Synced Text"
				modal.formData.CompetencyCategories = []string{"synced-category"}
				modal.formData.RoleFit = string(career.RoleFitStaff)
				modal.formData.AudienceRelevance = []string{"synced-audience"}

				modal.syncModified()

				Expect(modal.modified.Text).To(Equal("Synced Text"))
				Expect(modal.modified.CompetencyCategories).To(ContainElement("synced-category"))
				Expect(modal.modified.RoleFit).To(Equal(career.RoleFitStaff))
				Expect(modal.modified.AudienceRelevance).To(ContainElement("synced-audience"))
				// StrengthSignal should be preserved from original
				Expect(modal.modified.StrengthSignal).To(Equal(modal.original.StrengthSignal))
				// ID should be preserved
				Expect(modal.modified.ID).To(Equal(modal.original.ID))
			})

			It("preserves source event ID", func() {
				modal.syncModified()

				Expect(modal.modified.SourceEventID).To(Equal(modal.original.SourceEventID))
			})
		})
	})

	// ============================================================================
	// Helper Function Tests
	// ============================================================================
	Describe("copyMetadataSnapshot", func() {
		It("returns nil for nil input", func() {
			result := copyMetadataSnapshot(nil)
			Expect(result).To(BeNil())
		})

		It("creates a deep copy", func() {
			original := &MetadataSnapshot{
				Company:    "Company",
				Project:    "Project",
				Tags:       []string{"tag1", "tag2"},
				Categories: []string{"cat1"},
			}

			copy := copyMetadataSnapshot(original)

			// Verify values are equal
			Expect(copy.Company).To(Equal(original.Company))
			Expect(copy.Project).To(Equal(original.Project))
			Expect(copy.Tags).To(Equal(original.Tags))
			Expect(copy.Categories).To(Equal(original.Categories))

			// Verify slices are independent
			copy.Tags[0] = "modified"
			Expect(original.Tags[0]).To(Equal("tag1"))
		})
	})

	Describe("slicesEqual", func() {
		It("returns true for equal slices", func() {
			a := []string{"a", "b", "c"}
			b := []string{"a", "b", "c"}

			Expect(slicesEqual(a, b)).To(BeTrue())
		})

		It("returns false for different lengths", func() {
			a := []string{"a", "b"}
			b := []string{"a", "b", "c"}

			Expect(slicesEqual(a, b)).To(BeFalse())
		})

		It("returns false for different content", func() {
			a := []string{"a", "b", "c"}
			b := []string{"a", "x", "c"}

			Expect(slicesEqual(a, b)).To(BeFalse())
		})

		It("returns true for empty slices", func() {
			a := []string{}
			b := []string{}

			Expect(slicesEqual(a, b)).To(BeTrue())
		})

		It("returns true for nil slices", func() {
			var a, b []string

			Expect(slicesEqual(a, b)).To(BeTrue())
		})
	})
})
