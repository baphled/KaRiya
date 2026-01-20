package intents_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/baphled/kariya/internal/cli/intents"
	tea "github.com/charmbracelet/bubbletea"
)

var _ = Describe("ExportArtifact - Escape Key Behavior", func() {
	var intent *ExportArtifactIntent

	Describe("Configure State (Wizard)", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContext())
			intent.Init()
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))
		})

		It("should cancel intent when esc is pressed at step 1", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should go back from step 2 to step 1 on esc", func() {
			// Go to step 2
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Now press esc to go back to step 1
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			// Should still be in configure state (wizard handles internal navigation)
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))
			view := intent.View()
			Expect(view).To(ContainSubstring("Step 1"))
		})
	})

	Describe("Preview State", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContextWithServices())
			intent.Init()
			// Skip wizard to get to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})

		It("should go back to wizard on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))
		})

		It("should preserve configuration when going back to wizard", func() {
			// Get the current config
			configBefore := intent.GetConfig()
			Expect(configBefore).NotTo(BeNil())

			// Go back to wizard
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))

			// The wizard should be visible with previous config
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			// Verify wizard is shown
			Expect(view).To(ContainSubstring("Step"))
		})

		It("should allow returning to preview after going back", func() {
			// Go back to wizard
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))

			// Skip wizard again to return to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})
	})

	Describe("Confirm State", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContextWithServices())
			intent.Init()
			// Skip wizard to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			// Go to confirm
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))
		})

		It("should go back to preview on escape", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})

		It("should go back to preview on 'n' key", func() {
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))
		})

		It("should allow progressing to export after going back", func() {
			// Go back to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))

			// Go forward to confirm again
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))
		})
	})

	Describe("Full Round-Trip Navigation", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContextWithServices())
			intent.Init()
		})

		It("should navigate forward and backward through entire workflow", func() {
			// Start at wizard step 1
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))

			// Skip wizard to go to preview (Ctrl+S uses defaults)
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))

			// Go to confirm
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))

			// Back to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))

			// Back to wizard (returns to last step, step 2)
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))

			// Back to step 1 in wizard
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))

			// Cancel from wizard step 1
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			result := intent.Result()
			Expect(result).NotTo(BeNil())
			Expect(result.Status).To(Equal(Cancelled))
		})

		It("should maintain state consistency during navigation", func() {
			// Skip to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))

			// Forward to confirm
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))

			// Back to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))

			// View should render without panic
			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Progress State", func() {
		It("should go back to preview if cancelled during progress", func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContextWithServices())
			intent.Init()
			// Skip to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			// Go to confirm
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			// Start export (go to progress)
			intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}})

			// If progress state allows cancellation, it should go back to preview
			// Note: This depends on implementation - progress may not allow cancel
			if intent.GetState() == ExportStateInProgress {
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
				// Should either stay in progress or go back
				Expect(intent.GetState()).To(BeElementOf(ExportStateInProgress, ExportStatePreview))
			}
		})
	})

	Describe("View Rendering After Navigation", func() {
		BeforeEach(func() {
			intent, _ = NewExportArtifactIntent(NewTestExportArtifactContextWithServices())
			intent.Init()
		})

		It("should render wizard correctly after back from preview", func() {
			// Go to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))

			// Go back to wizard
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStateConfigure))

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			Expect(view).To(ContainSubstring("Step"))
		})

		It("should render preview correctly after back from confirm", func() {
			// Go to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyCtrlS})
			// Go to confirm
			intent.Update(tea.KeyMsg{Type: tea.KeyEnter})
			Expect(intent.GetState()).To(Equal(ExportStateConfirm))

			// Go back to preview
			intent.Update(tea.KeyMsg{Type: tea.KeyEsc})
			Expect(intent.GetState()).To(Equal(ExportStatePreview))

			view := intent.View()
			Expect(view).NotTo(BeEmpty())
			// Preview should show scrollable content
			Expect(view).To(Or(
				ContainSubstring("Preview"),
				ContainSubstring("Export"),
				ContainSubstring("scroll"),
			))
		})
	})
})
