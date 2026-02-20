package e2e_test

import (
	"strings"

	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Skill Inference E2E", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.GetSharedEnv(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	It("infers skills from event text and displays them on review screen", func() {
		env.SelectIntentByName("capture_event")
		env.Confirm() // Select quick strategy

		env.TypeText("Built a backend service using Go and PostgreSQL")
		env.Confirm() // Submit text

		testEvent := fixtures.EventWith("", "Built a backend service using Go and PostgreSQL", "TechCorp", "Backend")
		testEvent.ID = ""
		env.SubmitEvent(testEvent)

		env.AssertViewContains("Skills")
		env.AssertViewContains("Go")
		env.AssertViewContains("PostgreSQL")
	})

	It("allows editing skills via 's' shortcut", func() {
		env.SelectIntentByName("capture_event")
		env.Confirm() // Select quick strategy
		testEvent := fixtures.EventWith("", "Refactoring Java code", "Legacy", "Refactor")
		env.SubmitEvent(testEvent)

		env.PressKeyRune('s') // Open skill editor
		env.AssertViewContains("Review Skill Suggestions")
		env.AssertViewContains("Java")

		env.PressKeyRune('a') // Accept skill

		if strings.Contains(env.GetView(), "Review Skill Suggestions") {
			env.PressKey(tea.KeyEsc) // Close modal
		}
	})

	It("persists accepted skills on submission", func() {
		env.SelectIntentByName("capture_event")
		env.Confirm() // Select quick strategy
		testEvent := fixtures.EventWith("", "Working with Python", "AI", "Scripting")
		env.SubmitEvent(testEvent)

		env.PressKeyRune('s') // Open skill editor
		env.PressKeyRune('a') // Accept Python
		if strings.Contains(env.GetView(), "Review Skill Suggestions") {
			env.PressKey(tea.KeyEsc) // Close modal
		}

		env.Confirm() // Submit review

		events := env.GetEvents()
		Expect(events).To(HaveLen(1))
		Expect(events[0].ID).NotTo(BeEmpty())

		skills := env.GetSkills()
		Expect(skills).To(HaveLen(1))
		Expect(skills[0].Name).To(Equal("Python"))
	})

	It("does not persist rejected skills", func() {
		env.SelectIntentByName("capture_event")
		env.Confirm() // Select quick strategy
		testEvent := fixtures.EventWith("", "Using Rust language", "Systems", "Core")
		env.SubmitEvent(testEvent)

		env.PressKeyRune('s') // Open skill editor
		env.PressKeyRune('r') // Reject Rust
		if strings.Contains(env.GetView(), "Review Skill Suggestions") {
			env.PressKey(tea.KeyEsc) // Close modal
		}

		env.Confirm() // Submit review

		skills := env.GetSkills()
		Expect(skills).To(BeEmpty())
	})
})
