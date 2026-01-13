package intents_test

import (
	"github.com/baphled/kariya/internal/testutil/e2e"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Note: These tests focus on navigation paths and UI rendering with empty state.
// Tests with populated skills data can be added once e2e.TestEnv supports SkillRepository.

var _ = Describe("ManageSkills Navigation", func() {
	var env *e2e.TestEnv

	Describe("Menu Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show ManageSkills as menu item", func() {
			env.AssertViewContains("Manage Skills")
		})

		It("should navigate to ManageSkills when selected", func() {
			env.SelectIntentByName("manage_skills")
			env.AssertViewContainsAny("Skills", "Manage Skills")
		})

		It("should show context help for skills list", func() {
			env.SelectIntentByName("manage_skills")
			env.AssertViewContainsAny("Enter", "Esc", "n", "New skill")
		})

		It("should allow navigation back from skills list", func() {
			env.SelectIntentByName("manage_skills")
			env.Cancel()
			env.AssertViewContains("Capture Event")
		})
	})

	Describe("Empty List Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should show empty message when no skills exist", func() {
			env.SelectIntentByName("manage_skills")
			env.AssertViewContainsAny("No skills", "empty", "0 skills")
		})

		It("should not panic on empty list", func() {
			env.SelectIntentByName("manage_skills")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow adding skill from empty list", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.AssertViewContainsAny("Add", "Skill Name")
		})
	})

	Describe("Filter Menu Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open filter menu with 'f'", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			env.AssertViewContainsAny("Filter", "Category")
		})

		It("should navigate back from filter menu", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			env.Cancel()
			env.AssertViewContains("Skills")
		})

		It("should show filter options", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			env.AssertViewContainsAny("Category", "Level", "All")
		})

		It("should navigate through filter options", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should open filter modal on 'f' key", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			// Modal should be visible - check for any filter-related content
			env.AssertViewContainsAny("Category", "Level", "Sort", "Filter")
		})

		It("should show filter modal with form fields", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			// Modal should show filter options (huh forms show fields progressively)
			env.AssertViewContainsAny("Category", "Level", "Sort")
		})
	})

	Describe("Sort Menu Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open sort menu with 's'", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			env.AssertViewContainsAny("Sort", "Name", "Category")
		})

		It("should navigate back from sort menu", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			env.Cancel()
			env.AssertViewContains("Skills")
		})

		It("should show sort options", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			env.AssertViewContainsAny("Name", "Category", "Level")
		})

		It("should navigate through sort options", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should open sort modal on 's' key", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			// Modal should be visible with sort options
			env.AssertViewContainsAny("Sort", "Name", "Category", "Level")
		})

		It("should show sort menu title", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			env.AssertViewContains("Sort")
		})
	})

	Describe("Add Form Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should open add form with 'n'", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.AssertViewContainsAny("Add", "Skill Name")
		})

		It("should navigate back from add form", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.Cancel()
			env.AssertViewContains("Skills")
		})

		It("should show field labels", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.AssertViewContainsAny("Name", "Category", "Level")
		})

		It("should navigate through form fields with Tab", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.Tab()
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow typing in name field", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.TypeText("Go")
			view := env.GetView()
			Expect(view).To(ContainSubstring("Go"))
		})

		It("should show category options", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.Tab()
			env.AssertViewContainsAny("backend", "frontend", "devops")
		})

		It("should show level options", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.Tab()
			env.Tab()
			env.AssertViewContainsAny("beginner", "intermediate", "advanced", "expert")
		})

		It("should preserve data when navigating between fields", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.TypeText("Go")
			env.Tab()
			env.Tab()
			// Data should still be there
			view := env.GetView()
			Expect(view).To(ContainSubstring("Go"))
		})

		It("should show form breadcrumbs", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.AssertViewContainsAny("Skills", "Add", "▸")
		})

		It("should show Esc to cancel in footer", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.AssertViewContains("Esc")
		})
	})

	Describe("Workflow Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should navigate from menu to list", func() {
			env.SelectIntentByName("manage_skills")
			env.AssertViewContains("Skills")
		})

		It("should navigate from list to add form", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.AssertViewContainsAny("Add", "Skill Name")
		})

		It("should navigate from add form back to list", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			env.Cancel()
			env.AssertViewContains("Skills")
		})

		It("should navigate to filter and back", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			env.AssertViewContains("Filter")
			env.Cancel()
			env.AssertViewContains("Skills")
		})

		It("should navigate to sort and back", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			env.AssertViewContains("Sort")
			env.Cancel()
			env.AssertViewContains("Skills")
		})
	})

	Describe("Vim-style Navigation", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should support 'j' for down navigation in filter", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			env.PressKeyRune('j')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should support 'k' for up navigation in filter", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			env.PressKeyRune('j')
			env.PressKeyRune('k')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
		})

		It("should close filter modal with Esc key", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			// Modal should be visible
			env.AssertViewContainsAny("Category", "Level", "Sort")
			env.Cancel() // Esc key
			// Should return to skills list
			env.AssertViewContains("Skills")
		})

		It("should close sort modal with Esc key", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			// Modal should be visible
			env.AssertViewContainsAny("Sort", "Name", "Category")
			env.Cancel() // Esc key
			// Should return to skills list
			env.AssertViewContains("Skills")
		})
	})

	Describe("View Rendering", func() {
		BeforeEach(func() {
			env = e2e.SetupWithMemory(GinkgoT())
		})

		AfterEach(func() {
			env.Cleanup()
		})

		It("should render list view without errors", func() {
			env.SelectIntentByName("manage_skills")
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("error"))
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should render filter menu without errors", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('f')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("error"))
		})

		It("should render sort menu without errors", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('s')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("error"))
		})

		It("should render add form without errors", func() {
			env.SelectIntentByName("manage_skills")
			env.PressKeyRune('n')
			view := env.GetView()
			Expect(view).NotTo(BeEmpty())
			Expect(view).NotTo(ContainSubstring("error"))
		})

		It("should show footer in all views", func() {
			env.SelectIntentByName("manage_skills")
			env.AssertViewContainsAny("Esc", "Quit")
			env.PressKeyRune('f')
			env.AssertViewContainsAny("Esc", "Quit")
			env.Cancel()
			env.PressKeyRune('s')
			env.AssertViewContainsAny("Esc", "Quit")
		})
	})
})
