package intents

import (
	"time"

	"github.com/baphled/kariya/internal/domain/career"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EditMetadataModal", func() {
	var modal *EditMetadataModal

	BeforeEach(func() {
		modal = NewEditMetadataModal("Acme Corp", "Project X", []string{"golang", "tui"}, []string{"backend", "devops"})
	})

	Describe("Creation", func() {
		It("should create a new metadata modal with initial values", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.original).NotTo(BeNil())
			Expect(modal.modified).NotTo(BeNil())
			Expect(modal.original.Company).To(Equal("Acme Corp"))
			Expect(modal.original.Project).To(Equal("Project X"))
			Expect(modal.original.Tags).To(Equal([]string{"golang", "tui"}))
			Expect(modal.original.Categories).To(Equal([]string{"backend", "devops"}))
		})

		It("should initialize with focused on first field", func() {
			Expect(modal.focused).To(Equal(0))
		})

		It("should not be complete initially", func() {
			Expect(modal.IsComplete()).To(BeFalse())
		})

		It("should initialize all inputs", func() {
			Expect(modal.inputs).To(HaveLen(4))
			Expect(modal.inputs[0].Focused()).To(BeTrue())
			Expect(modal.inputs[1].Focused()).To(BeFalse())
			Expect(modal.inputs[2].Focused()).To(BeFalse())
			Expect(modal.inputs[3].Focused()).To(BeFalse())
		})

		It("should set input values correctly", func() {
			Expect(modal.inputs[0].Value()).To(Equal("Acme Corp"))
			Expect(modal.inputs[1].Value()).To(Equal("Project X"))
			Expect(modal.inputs[2].Value()).To(Equal("golang, tui"))
			Expect(modal.inputs[3].Value()).To(Equal("backend, devops"))
		})

		It("should handle empty metadata", func() {
			emptyModal := NewEditMetadataModal("", "", []string{}, []string{})
			Expect(emptyModal.original.Company).To(Equal(""))
			Expect(emptyModal.original.Project).To(Equal(""))
			Expect(emptyModal.original.Tags).To(BeEmpty())
			Expect(emptyModal.original.Categories).To(BeEmpty())
		})
	})

	Describe("Navigation", func() {
		It("should move focus forward with Tab", func() {
			Expect(modal.focused).To(Equal(0))
			Expect(modal.inputs[0].Focused()).To(BeTrue())

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(1))
			Expect(modal.inputs[0].Focused()).To(BeFalse())
			Expect(modal.inputs[1].Focused()).To(BeTrue())

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(2))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(3))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(0)) // Wraps around
		})

		It("should move focus backward with Shift+Tab", func() {
			modal.focused = 0
			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(3)) // Wraps around

			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(2))

			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(1))

			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(0))
		})

		It("should handle Tab wrapping from last field to first", func() {
			modal.focused = 3
			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(0))
		})

		It("should handle Shift+Tab wrapping from first field to last", func() {
			modal.focused = 0
			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(3))
		})
	})

	Describe("Confirmation", func() {
		It("should create result on Enter with changes", func() {
			modal.inputs[0].SetValue("NewCorp")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			Expect(modal.IsComplete()).To(BeTrue())
			result := modal.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Accepted).To(BeTrue())
			Expect(result.Modified.Company).To(Equal("NewCorp"))
		})

		It("should track changes in result", func() {
			modal.inputs[0].SetValue("NewCorp")
			modal.inputs[1].SetValue("NewProject")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := modal.Result()
			Expect(result.Changes).To(HaveKey("company"))
			Expect(result.Changes).To(HaveKey("project"))
			Expect(result.Changes["company"]).To(Equal("NewCorp"))
			Expect(result.Changes["project"]).To(Equal("NewProject"))
		})

		It("should not track unchanged fields", func() {
			modal.inputs[0].SetValue("Acme Corp") // Same as original
			modal.inputs[1].SetValue("NewProject")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := modal.Result()
			Expect(result.Changes).NotTo(HaveKey("company"))
			Expect(result.Changes).To(HaveKey("project"))
		})

		It("should parse tags correctly on confirmation", func() {
			modal.inputs[2].SetValue("go, rust, python")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := modal.Result()
			Expect(result.Modified.Tags).To(Equal([]string{"go", "rust", "python"}))
		})

		It("should parse categories correctly on confirmation", func() {
			modal.inputs[3].SetValue("backend, devops, security")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := modal.Result()
			Expect(result.Modified.Categories).To(Equal([]string{"backend", "devops", "security"}))
		})

		It("should not return nil result after Enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.result).NotTo(BeNil())
		})
	})

	Describe("Cancellation", func() {
		It("should create cancelled result on Escape", func() {
			modal.inputs[0].SetValue("NewCorp")
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			Expect(modal.IsComplete()).To(BeTrue())
			result := modal.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Accepted).To(BeFalse())
			Expect(result.Modified).To(Equal(result.Original))
		})

		It("should preserve original data on cancel", func() {
			modal.inputs[0].SetValue("NewCorp")
			modal.inputs[1].SetValue("NewProject")
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := modal.Result()
			Expect(result.Original.Company).To(Equal("Acme Corp"))
			Expect(result.Original.Project).To(Equal("Project X"))
		})

		It("should not track changes on cancel", func() {
			modal.inputs[0].SetValue("NewCorp")
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := modal.Result()
			Expect(result.Changes).To(BeEmpty())
		})

		It("should restore original values on cancel", func() {
			modal.inputs[0].SetValue("NewCorp")
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := modal.Result()
			Expect(result.Modified.Company).To(Equal("Acme Corp"))
		})
	})

	Describe("Window Sizing", func() {
		It("should respond to window size changes", func() {
			Expect(modal.width).To(Equal(80))
			Expect(modal.height).To(Equal(24))

			modal.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
			Expect(modal.width).To(Equal(120))
			Expect(modal.height).To(Equal(40))
		})

		It("should handle small window sizes", func() {
			modal.Update(tea.WindowSizeMsg{Width: 30, Height: 10})
			Expect(modal.width).To(Equal(30))
			Expect(modal.height).To(Equal(10))
		})

		It("should handle large window sizes", func() {
			modal.Update(tea.WindowSizeMsg{Width: 200, Height: 100})
			Expect(modal.width).To(Equal(200))
			Expect(modal.height).To(Equal(100))
		})
	})

	Describe("View Rendering", func() {
		It("should render modal view", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Edit Event Metadata"))
		})

		It("should contain field labels in view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Company"))
			Expect(view).To(ContainSubstring("Project"))
		})

		It("should return empty view after completion", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := modal.View()
			Expect(view).To(BeEmpty())
		})

		It("should show focus indicator in view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("→")) // Focus indicator
		})

		It("should show keyboard shortcuts in view", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Enter"))
			Expect(view).To(ContainSubstring("Esc"))
			Expect(view).To(ContainSubstring("Tab"))
		})
	})

	Describe("Result Handling", func() {
		It("should return nil result before completion", func() {
			Expect(modal.Result()).To(BeNil())
		})

		It("should return result after Enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.Result()).NotTo(BeNil())
		})

		It("should return result after Escape", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.Result()).NotTo(BeNil())
		})

		It("should mark result as accepted on Enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.Result().Accepted).To(BeTrue())
		})

		It("should mark result as not accepted on Escape", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.Result().Accepted).To(BeFalse())
		})
	})
})

var _ = Describe("EditBurstModal", func() {
	var modal *EditBurstModal
	var burst *career.Burst

	BeforeEach(func() {
		burst = &career.Burst{
			ID:              "test-burst-id",
			Name:            "Test Burst",
			Description:     "Test Description",
			CompetencyFocus: "Test Focus",
			EventIDs:        []string{},
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		modal = NewEditBurstModal(burst)
	})

	Describe("Creation", func() {
		It("should create a new burst modal with initial values", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.original).NotTo(BeNil())
			Expect(modal.original.Name).To(Equal("Test Burst"))
			Expect(modal.original.Description).To(Equal("Test Description"))
			Expect(modal.original.CompetencyFocus).To(Equal("Test Focus"))
		})

		It("should initialize with focused on first field", func() {
			Expect(modal.focused).To(Equal(0))
		})

		It("should not be complete initially", func() {
			Expect(modal.IsComplete()).To(BeFalse())
		})

		It("should initialize all inputs", func() {
			Expect(modal.inputs).To(HaveLen(3))
			Expect(modal.inputs[0].Focused()).To(BeTrue())
		})
	})

	Describe("Navigation", func() {
		It("should move focus forward with Tab", func() {
			Expect(modal.focused).To(Equal(0))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(1))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(2))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(0)) // Wrap around
		})

		It("should move focus backward with Shift+Tab", func() {
			modal.focused = 0
			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(2)) // Wraps around

			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(1))
		})
	})

	Describe("Confirmation", func() {
		It("should create result on Enter with changes", func() {
			modal.inputs[0].SetValue("Modified Name")

			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(modal.accepted).To(BeTrue())
			Expect(modal.result).NotTo(BeNil())
			Expect(modal.result.Accepted).To(BeTrue())
			Expect(modal.result.Modified.Name).To(Equal("Modified Name"))
		})

		It("should track changes in result", func() {
			modal.inputs[0].SetValue("New Name")
			modal.inputs[1].SetValue("New Description")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := modal.Result()
			Expect(result.Changes).To(HaveKey("name"))
			Expect(result.Changes).To(HaveKey("description"))
		})
	})

	Describe("Cancellation", func() {
		It("should create cancelled result on Escape", func() {
			modal.inputs[0].SetValue("Modified Name")

			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(modal.accepted).To(BeFalse())
			Expect(modal.result).NotTo(BeNil())
			Expect(modal.result.Accepted).To(BeFalse())
			Expect(modal.result.Modified.Name).To(Equal(burst.Name))
		})

		It("should preserve original data on cancel", func() {
			modal.inputs[0].SetValue("Modified Name")
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := modal.Result()
			Expect(result.Original.Name).To(Equal("Test Burst"))
		})
	})

	Describe("View Rendering", func() {
		It("should render modal view", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Edit Burst"))
		})

		It("should contain field labels", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Name"))
			Expect(view).To(ContainSubstring("Description"))
		})

		It("should return empty view after completion", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := modal.View()
			Expect(view).To(BeEmpty())
		})
	})

	Describe("Result Handling", func() {
		It("should return nil result before completion", func() {
			Expect(modal.Result()).To(BeNil())
		})

		It("should return result after Enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.Result()).NotTo(BeNil())
		})

		It("should return result after Escape", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.Result()).NotTo(BeNil())
		})

		It("should track IsComplete correctly", func() {
			Expect(modal.IsComplete()).To(BeFalse())
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.IsComplete()).To(BeTrue())
		})
	})
})

var _ = Describe("EditFactModal", func() {
	var modal *EditFactModal
	var fact *career.Fact

	BeforeEach(func() {
		fact = &career.Fact{
			ID:                   "test-fact-id",
			Text:                 "Test Fact",
			CompetencyCategories: []string{"category1"},
			RoleFit:              career.RoleFitPrincipal,
			AudienceRelevance:    []string{"audience1"},
			StrengthSignal:       "Test Signal",
			SourceEventID:        "event-id",
			SourceBurstID:        "burst-id",
			CreatedAt:            time.Now(),
			UpdatedAt:            time.Now(),
		}
		modal = NewEditFactModal(fact)
	})

	Describe("Creation", func() {
		It("should create a new fact modal with initial values", func() {
			Expect(modal).NotTo(BeNil())
			Expect(modal.original).NotTo(BeNil())
			Expect(modal.original.Text).To(Equal("Test Fact"))
			Expect(modal.original.CompetencyCategories).To(Equal([]string{"category1"}))
			Expect(modal.original.RoleFit).To(Equal(career.RoleFitPrincipal))
		})

		It("should initialize with focused on first field", func() {
			Expect(modal.focused).To(Equal(0))
		})

		It("should not be complete initially", func() {
			Expect(modal.IsComplete()).To(BeFalse())
		})

		It("should initialize all inputs", func() {
			Expect(modal.inputs).To(HaveLen(5))
			Expect(modal.inputs[0].Focused()).To(BeTrue())
		})
	})

	Describe("Navigation", func() {
		It("should move focus forward with Tab", func() {
			Expect(modal.focused).To(Equal(0))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(1))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(2))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(3))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(4))

			modal.Update(tea.KeyMsg{Type: tea.KeyTab})
			Expect(modal.focused).To(Equal(0)) // Wrap around
		})

		It("should move focus backward with Shift+Tab", func() {
			modal.focused = 0
			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(4)) // Wraps around

			modal.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
			Expect(modal.focused).To(Equal(3))
		})
	})

	Describe("Confirmation", func() {
		It("should create result on Enter with changes", func() {
			modal.inputs[0].SetValue("Modified Fact")

			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(cmd).To(BeNil())
			Expect(modal.accepted).To(BeTrue())
			Expect(modal.result).NotTo(BeNil())
			Expect(modal.result.Accepted).To(BeTrue())
			Expect(modal.result.Modified.Text).To(Equal("Modified Fact"))
		})

		It("should track changes in result", func() {
			modal.inputs[0].SetValue("New Text")
			modal.inputs[1].SetValue("new, categories")
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})

			result := modal.Result()
			Expect(result.Changes).To(HaveKey("text"))
			Expect(result.Changes).To(HaveKey("competency_categories"))
		})
	})

	Describe("Cancellation", func() {
		It("should create cancelled result on Escape", func() {
			modal.inputs[0].SetValue("Modified Fact")

			cmd := modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(cmd).To(BeNil())
			Expect(modal.accepted).To(BeFalse())
			Expect(modal.result).NotTo(BeNil())
			Expect(modal.result.Accepted).To(BeFalse())
			Expect(modal.result.Modified.Text).To(Equal(fact.Text))
		})

		It("should preserve original data on cancel", func() {
			modal.inputs[0].SetValue("Modified Fact")
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})

			result := modal.Result()
			Expect(result.Original.Text).To(Equal("Test Fact"))
		})
	})

	Describe("View Rendering", func() {
		It("should render modal view", func() {
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Edit Fact"))
		})

		It("should contain field labels", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Text"))
			Expect(view).To(ContainSubstring("Categories"))
		})

		It("should return empty view after completion", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			view := modal.View()
			Expect(view).To(BeEmpty())
		})
	})

	Describe("Result Handling", func() {
		It("should return nil result before completion", func() {
			Expect(modal.Result()).To(BeNil())
		})

		It("should return result after Enter", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.Result()).NotTo(BeNil())
		})

		It("should return result after Escape", func() {
			modal.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(modal.Result()).NotTo(BeNil())
		})

		It("should track IsComplete correctly", func() {
			Expect(modal.IsComplete()).To(BeFalse())
			modal.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(modal.IsComplete()).To(BeTrue())
		})
	})
})

var _ = Describe("Helper Functions", func() {
	Describe("formatStringSlice", func() {
		It("should handle empty slice", func() {
			result := formatStringSlice([]string{})
			Expect(result).To(Equal(""))
		})

		It("should handle single item", func() {
			result := formatStringSlice([]string{"item1"})
			Expect(result).To(Equal("item1"))
		})

		It("should handle multiple items", func() {
			result := formatStringSlice([]string{"item1", "item2", "item3"})
			Expect(result).To(Equal("item1, item2, item3"))
		})
	})

	Describe("parseStringSlice", func() {
		It("should handle empty string", func() {
			result := parseStringSlice("")
			Expect(result).To(BeEmpty())
		})

		It("should handle single item", func() {
			result := parseStringSlice("item1")
			Expect(result).To(Equal([]string{"item1"}))
		})

		It("should handle multiple items", func() {
			result := parseStringSlice("item1, item2, item3")
			Expect(result).To(Equal([]string{"item1", "item2", "item3"}))
		})

		It("should handle items with extra spaces", func() {
			result := parseStringSlice("  item1  ,  item2  ")
			Expect(result).To(Equal([]string{"item1", "item2"}))
		})

		It("should handle items with trailing comma", func() {
			result := parseStringSlice("item1, item2,")
			Expect(result).To(Equal([]string{"item1", "item2"}))
		})
	})

	Describe("slicesEqual", func() {
		It("should handle both empty slices", func() {
			result := slicesEqual([]string{}, []string{})
			Expect(result).To(BeTrue())
		})

		It("should handle different lengths", func() {
			result := slicesEqual([]string{"a"}, []string{"a", "b"})
			Expect(result).To(BeFalse())
		})

		It("should handle same content", func() {
			result := slicesEqual([]string{"a", "b"}, []string{"a", "b"})
			Expect(result).To(BeTrue())
		})

		It("should handle different content", func() {
			result := slicesEqual([]string{"a", "b"}, []string{"b", "a"})
			Expect(result).To(BeFalse())
		})
	})

	Describe("truncateString", func() {
		It("should not truncate when not needed", func() {
			result := truncateString("hello", 10)
			Expect(result).To(Equal("hello"))
		})

		It("should handle exact length", func() {
			result := truncateString("hello", 5)
			Expect(result).To(Equal("hello"))
		})

		It("should truncate when needed", func() {
			result := truncateString("hello world", 8)
			Expect(result).To(Equal("hello..."))
		})

		It("should handle single character", func() {
			result := truncateString("a", 1)
			Expect(result).To(Equal("a"))
		})
	})

	Describe("copyMetadataSnapshot", func() {
		It("should create a deep copy of metadata snapshot", func() {
			original := &MetadataSnapshot{
				Company:    "Corp",
				Project:    "Proj",
				Tags:       []string{"a", "b"},
				Categories: []string{"x", "y"},
			}

			copy := copyMetadataSnapshot(original)
			copy.Company = "NewCorp"
			copy.Tags[0] = "z"

			Expect(original.Company).To(Equal("Corp"))
			Expect(original.Tags[0]).To(Equal("a"))
		})

		It("should handle nil input", func() {
			copy := copyMetadataSnapshot(nil)
			Expect(copy).To(BeNil())
		})
	})
})

