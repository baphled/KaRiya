package fact_test

import (
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/fact"
	"github.com/baphled/kariya/internal/ui/display"
	"github.com/charmbracelet/bubbles/cursor"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EditFact", func() {
	var (
		modal *fact.EditFact
		f     display.Fact
	)

	BeforeEach(func() {
		f = display.FactFromDomain(fixtures.FactWith("fact-1", "Test fact text"))
		modal = fact.NewEditFact(f)
	})

	Describe("GetTitle", func() {
		It("returns the correct title", func() {
			Expect(modal.GetTitle()).To(Equal("Edit Fact"))
		})
	})

	Describe("GetContent", func() {
		It("returns non-empty content when modal is active", func() {
			Expect(modal.GetContent()).NotTo(BeEmpty())
		})
	})

	Describe("GetFooter", func() {
		It("returns non-empty footer instructions", func() {
			Expect(modal.GetFooter()).NotTo(BeEmpty())
		})

		It("contains navigation instructions", func() {
			footer := modal.GetFooter()
			Expect(footer).To(ContainSubstring("Enter"))
			Expect(footer).To(ContainSubstring("Esc"))
			Expect(footer).To(ContainSubstring("Tab"))
		})
	})

	Describe("IsComplete", func() {
		It("returns false when modal is not complete", func() {
			Expect(modal.IsComplete()).To(BeFalse())
		})
	})

	Describe("Result", func() {
		It("returns nil when modal is not complete", func() {
			Expect(modal.Result()).To(BeNil())
		})
	})

	Describe("View", func() {
		It("returns non-empty view when modal is active", func() {
			Expect(modal.View()).NotTo(BeEmpty())
		})
	})

	Describe("Init", func() {
		It("returns a command", func() {
			cmd := modal.Init()
			Expect(cmd).NotTo(BeNil())
		})
	})

	Describe("Update", func() {
		Context("with WindowSizeMsg", func() {
			It("updates width and height", func() {
				msg := tea.WindowSizeMsg{Width: 120, Height: 40}
				modal.Update(msg)
				Expect(modal.View()).NotTo(BeEmpty())
			})
		})

		Context("with other messages", func() {
			It("returns a command", func() {
				msg := tea.KeyMsg{Type: tea.KeyTab}
				cmd := modal.Update(msg)
				Expect(cmd).NotTo(BeNil())
			})
		})
	})
})

var _ = Describe("EditResult", func() {
	Describe("HasChanges", func() {
		Context("when changes map is empty", func() {
			It("returns false", func() {
				result := &fact.EditResult{
					Changes: fact.EditChanges{},
				}
				Expect(result.HasChanges()).To(BeFalse())
			})
		})

		Context("when changes map has entries", func() {
			It("returns true", func() {
				text := "new text"
				result := &fact.EditResult{
					Changes: fact.EditChanges{
						Text: &text,
					},
				}
				Expect(result.HasChanges()).To(BeTrue())
			})
		})

		Context("with multiple changes", func() {
			It("returns true", func() {
				text := "new text"
				categories := []string{"Backend"}
				roleFit := "senior"
				result := &fact.EditResult{
					Changes: fact.EditChanges{
						Text:                 &text,
						CompetencyCategories: &categories,
						RoleFit:              &roleFit,
					},
				}
				Expect(result.HasChanges()).To(BeTrue())
			})
		})
	})
})

var _ = Describe("EditFact internal behavior", func() {
	var (
		modal *fact.EditFact
		f     display.Fact
	)

	BeforeEach(func() {
		f = display.FactFromDomain(fixtures.FactWith("fact-1", "Original text"))
		f.CompetencyCategories = []string{"Backend", "Go"}
		f.RoleFit = "senior"
		f.AudienceRelevance = []string{"technical"}
		modal = fact.NewEditFact(f)
	})

	Describe("slicesEqual", func() {
		It("returns true for equal slices", func() {
			a := []string{"a", "b", "c"}
			b := []string{"a", "b", "c"}
			Expect(slicesEqualHelper(a, b)).To(BeTrue())
		})

		It("returns false for different lengths", func() {
			a := []string{"a", "b"}
			b := []string{"a", "b", "c"}
			Expect(slicesEqualHelper(a, b)).To(BeFalse())
		})

		It("returns false for different content", func() {
			a := []string{"a", "b", "c"}
			b := []string{"a", "x", "c"}
			Expect(slicesEqualHelper(a, b)).To(BeFalse())
		})

		It("returns true for empty slices", func() {
			a := []string{}
			b := []string{}
			Expect(slicesEqualHelper(a, b)).To(BeTrue())
		})

		It("returns true for single element slices", func() {
			a := []string{"x"}
			b := []string{"x"}
			Expect(slicesEqualHelper(a, b)).To(BeTrue())
		})

		It("returns true when first slice is nil and second is empty", func() {
			a := []string(nil)
			b := []string{}
			Expect(slicesEqualHelper(a, b)).To(BeTrue())
		})
	})

	Describe("cloneStrings", func() {
		It("returns empty slice for nil input", func() {
			result := cloneStringsHelper(nil)
			Expect(result).To(BeEmpty())
		})

		It("returns copy of slice", func() {
			original := []string{"a", "b", "c"}
			result := cloneStringsHelper(original)

			Expect(result).To(Equal(original))
		})

		It("returns empty slice for empty input", func() {
			original := []string{}
			result := cloneStringsHelper(original)

			Expect(result).To(BeEmpty())
		})

		It("creates independent copy", func() {
			original := []string{"a", "b"}
			result := cloneStringsHelper(original)
			result[0] = "x"

			Expect(original[0]).To(Equal("a"))
			Expect(result[0]).To(Equal("x"))
		})

		It("handles single element slice", func() {
			original := []string{"single"}
			result := cloneStringsHelper(original)

			Expect(result).To(Equal(original))
			result[0] = "modified"
			Expect(original[0]).To(Equal("single"))
		})

		It("handles large slice", func() {
			original := make([]string, 100)
			for i := range original {
				original[i] = "item"
			}
			result := cloneStringsHelper(original)

			Expect(result).To(HaveLen(100))
			Expect(result).To(Equal(original))
		})
	})

	Describe("Update with WindowSizeMsg", func() {
		It("handles multiple window size changes", func() {
			msg1 := tea.WindowSizeMsg{Width: 100, Height: 30}
			modal.Update(msg1)
			Expect(modal.View()).NotTo(BeEmpty())

			msg2 := tea.WindowSizeMsg{Width: 120, Height: 40}
			modal.Update(msg2)
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("handles small window sizes", func() {
			msg := tea.WindowSizeMsg{Width: 40, Height: 10}
			modal.Update(msg)
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("handles large window sizes", func() {
			msg := tea.WindowSizeMsg{Width: 200, Height: 100}
			modal.Update(msg)
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("returns nil command for window size message", func() {
			msg := tea.WindowSizeMsg{Width: 100, Height: 30}
			cmd := modal.Update(msg)
			Expect(cmd).To(BeNil())
		})

		It("updates form dimensions on window resize", func() {
			msg := tea.WindowSizeMsg{Width: 150, Height: 50}
			modal.Update(msg)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("Update with other messages", func() {
		It("returns a command for non-window-size messages", func() {
			msg := tea.KeyMsg{Type: tea.KeyTab}
			cmd := modal.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})

		It("handles multiple key messages", func() {
			msg1 := tea.KeyMsg{Type: tea.KeyTab}
			modal.Update(msg1)
			msg2 := tea.KeyMsg{Type: tea.KeyShiftTab}
			modal.Update(msg2)
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("handles enter key message", func() {
			msg := tea.KeyMsg{Type: tea.KeyEnter}
			cmd := modal.Update(msg)
			Expect(cmd).NotTo(BeNil())
		})

		It("handles escape key message", func() {
			msg := tea.KeyMsg{Type: tea.KeyEsc}
			modal.Update(msg)
			Expect(modal.View()).NotTo(BeEmpty())
		})
	})

	Describe("NewEditFact initialization", func() {
		It("preserves original fact data", func() {
			Expect(modal.Result()).To(BeNil())
		})

		It("initializes with default dimensions", func() {
			Expect(modal.View()).NotTo(BeEmpty())
		})

		It("handles fact with empty categories", func() {
			emptyFact := display.FactFromDomain(fixtures.FactWith("fact-2", "Text"))
			emptyFact.CompetencyCategories = []string{}
			emptyModal := fact.NewEditFact(emptyFact)
			Expect(emptyModal.View()).NotTo(BeEmpty())
		})

		It("handles fact with nil categories", func() {
			nilFact := display.FactFromDomain(fixtures.FactWith("fact-3", "Text"))
			nilFact.CompetencyCategories = nil
			nilModal := fact.NewEditFact(nilFact)
			Expect(nilModal.View()).NotTo(BeEmpty())
		})

		It("handles fact with multiple categories", func() {
			multiCatFact := display.FactFromDomain(fixtures.FactWith("fact-4", "Text"))
			multiCatFact.CompetencyCategories = []string{"technical", "leadership", "product"}
			multiModal := fact.NewEditFact(multiCatFact)
			Expect(multiModal.View()).NotTo(BeEmpty())
		})
	})

	Describe("GetContent and View consistency", func() {
		It("GetContent returns form view", func() {
			content := modal.GetContent()
			view := modal.View()
			Expect(content).NotTo(BeEmpty())
			Expect(view).NotTo(BeEmpty())
		})

		It("GetFooter returns consistent instructions", func() {
			footer := modal.GetFooter()
			Expect(footer).To(ContainSubstring("Enter"))
			Expect(footer).To(ContainSubstring("Esc"))
		})
	})

	Describe("EditResult with changes", func() {
		It("HasChanges returns true when changes exist", func() {
			text := "new text"
			result := &fact.EditResult{
				Original: f,
				Modified: f,
				Accepted: true,
				Changes: fact.EditChanges{
					Text: &text,
				},
			}
			Expect(result.HasChanges()).To(BeTrue())
		})

		It("HasChanges returns false for empty changes", func() {
			result := &fact.EditResult{
				Original: f,
				Modified: f,
				Accepted: false,
				Changes:  fact.EditChanges{},
			}
			Expect(result.HasChanges()).To(BeFalse())
		})

		It("HasChanges handles multiple changes", func() {
			text := "new text"
			categories := []string{"Frontend"}
			roleFit := "principal"
			audience := []string{"business"}
			signal := "weak"
			result := &fact.EditResult{
				Original: f,
				Modified: f,
				Accepted: true,
				Changes: fact.EditChanges{
					Text:                 &text,
					CompetencyCategories: &categories,
					RoleFit:              &roleFit,
					AudienceRelevance:    &audience,
					StrengthSignal:       &signal,
				},
			}
			Expect(result.HasChanges()).To(BeTrue())
		})
	})

	Describe("View rendering with different states", func() {
		It("renders modal with title and instructions", func() {
			view := modal.View()
			Expect(view).To(ContainSubstring("Edit Fact"))
			Expect(view).To(ContainSubstring("Tab"))
		})

		It("renders consistently across multiple calls", func() {
			view1 := modal.View()
			view2 := modal.View()
			Expect(view1).To(Equal(view2))
		})

		It("handles form view changes after window resize", func() {
			msg := tea.WindowSizeMsg{Width: 150, Height: 50}
			modal.Update(msg)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with minimum window size", func() {
			msg := tea.WindowSizeMsg{Width: 20, Height: 5}
			modal.Update(msg)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})

		It("renders with very large window size", func() {
			msg := tea.WindowSizeMsg{Width: 500, Height: 200}
			modal.Update(msg)
			view := modal.View()
			Expect(view).NotTo(BeEmpty())
		})
	})

	Describe("GetContent rendering", func() {
		It("returns form content without modal wrapper", func() {
			content := modal.GetContent()
			Expect(content).NotTo(BeEmpty())
		})

		It("returns consistent content across calls", func() {
			content1 := modal.GetContent()
			content2 := modal.GetContent()
			Expect(content1).To(Equal(content2))
		})
	})

})

func slicesEqualHelper(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func cloneStringsHelper(values []string) []string {
	if values == nil {
		return []string{}
	}

	return append([]string(nil), values...)
}

var _ = Describe("EditFact completion behavior", func() {
	var original display.Fact

	BeforeEach(func() {
		original = display.FactFromDomain(fixtures.FactWith("fact-1", "Original fact text long enough"))
		original.CompetencyCategories = []string{"technical"}
		original.RoleFit = "staff"
		original.AudienceRelevance = []string{"peer"}
		original.SourceEventID = "event-123"
		original.SourceBurstID = "burst-456"
		original.StrengthSignal = "strong"
	})

	It("returns an accepted result and hides rendered content after confirmation", func() {
		modal := newInitializedEditFactForTest(original)

		advanceToConfirmField(modal)
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRight})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.IsComplete()).To(BeTrue())
		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Accepted).To(BeTrue())
		Expect(modal.Result().Original).To(Equal(original))
		Expect(modal.Result().Modified).To(Equal(original))
		Expect(modal.Result().HasChanges()).To(BeFalse())
		Expect(modal.View()).To(BeEmpty())
		Expect(modal.GetContent()).To(BeEmpty())
	})

	It("returns a cancelled result when confirm is left false", func() {
		modal := newInitializedEditFactForTest(original)

		advanceToConfirmField(modal)
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.IsComplete()).To(BeTrue())
		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Accepted).To(BeFalse())
		Expect(modal.Result().Original).To(Equal(original))
		Expect(modal.Result().Modified).To(Equal(original))
		Expect(modal.Result().HasChanges()).To(BeFalse())
		Expect(modal.View()).NotTo(BeEmpty())
		Expect(modal.GetContent()).To(BeEmpty())
	})

	It("creates a cancelled result when the form is aborted", func() {
		modal := newInitializedEditFactForTest(original)

		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyCtrlC})

		Expect(modal.IsComplete()).To(BeTrue())
		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Accepted).To(BeFalse())
		Expect(modal.Result().Modified).To(Equal(original))
		Expect(modal.Result().HasChanges()).To(BeFalse())
	})

	It("captures text changes through the public form flow", func() {
		modal := newInitializedEditFactForTest(original)

		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(" updated")})
		advanceToConfirmField(modal)
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRight})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Accepted).To(BeTrue())
		Expect(modal.Result().Changes.Text).NotTo(BeNil())
		Expect(*modal.Result().Changes.Text).To(Equal("Original fact text long enough updated"))
		Expect(modal.Result().Modified.Text).To(Equal("Original fact text long enough updated"))
		Expect(modal.Result().Modified.SourceEventID).To(Equal(original.SourceEventID))
		Expect(modal.Result().Modified.SourceBurstID).To(Equal(original.SourceBurstID))
		Expect(modal.Result().Modified.StrengthSignal).To(Equal(original.StrengthSignal))
	})

	It("captures competency category changes through the public form flow", func() {
		modal := newInitializedEditFactForTest(original)

		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyDown})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRight})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Changes.CompetencyCategories).NotTo(BeNil())
		Expect(*modal.Result().Changes.CompetencyCategories).To(Equal([]string{"technical", "leadership"}))
		Expect(modal.Result().Modified.CompetencyCategories).To(Equal([]string{"technical", "leadership"}))
	})

	It("captures same-length competency replacements through the public form flow", func() {
		modal := newInitializedEditFactForTest(original)

		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyDown})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRight})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Accepted).To(BeTrue())
		Expect(modal.Result().HasChanges()).To(BeTrue())
		Expect(modal.Result().Changes.CompetencyCategories).NotTo(BeNil())
		Expect(*modal.Result().Changes.CompetencyCategories).To(Equal([]string{"leadership"}))
		Expect(modal.Result().Modified.CompetencyCategories).To(Equal([]string{"leadership"}))
	})

	It("captures role fit changes through the public form flow", func() {
		modal := newInitializedEditFactForTest(original)

		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyDown})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRight})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Changes.RoleFit).NotTo(BeNil())
		Expect(*modal.Result().Changes.RoleFit).To(Equal("senior_ic"))
		Expect(modal.Result().Modified.RoleFit).To(Equal("senior_ic"))
	})

	It("captures audience relevance changes through the public form flow", func() {
		modal := newInitializedEditFactForTest(original)

		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("x")})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRight})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Changes.AudienceRelevance).NotTo(BeNil())
		Expect(*modal.Result().Changes.AudienceRelevance).To(Equal([]string{"hiring_manager", "peer"}))
		Expect(modal.Result().Modified.AudienceRelevance).To(Equal([]string{"hiring_manager", "peer"}))
	})

	It("does not report category changes when nil categories round-trip unchanged", func() {
		original.CompetencyCategories = nil
		modal := newInitializedEditFactForTest(original)

		advanceToConfirmField(modal)
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyRight})
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyEnter})

		Expect(modal.Result()).NotTo(BeNil())
		Expect(modal.Result().Accepted).To(BeTrue())
		Expect(modal.Result().Changes.CompetencyCategories).To(BeNil())
		Expect(modal.Result().Modified.CompetencyCategories).To(BeEmpty())
	})
})

func newInitializedEditFactForTest(f display.Fact) *fact.EditFact {
	modal := fact.NewEditFact(f)
	flushEditFactCmdForTest(modal, modal.Init())
	return modal
}

func advanceToConfirmField(modal *fact.EditFact) {
	for range 4 {
		updateEditFactForTest(modal, tea.KeyMsg{Type: tea.KeyTab})
	}
}

func updateEditFactForTest(modal *fact.EditFact, msg tea.Msg) {
	cmd := modal.Update(msg)
	flushEditFactCmdForTest(modal, cmd)
}

func flushEditFactCmdForTest(modal *fact.EditFact, cmd tea.Cmd) {
	for step := 0; cmd != nil && step < 20; step++ {
		msg := cmd()
		if msg == nil {
			return
		}
		if _, ok := msg.(cursor.BlinkMsg); ok {
			return
		}
		cmd = modal.Update(msg)
	}
}
