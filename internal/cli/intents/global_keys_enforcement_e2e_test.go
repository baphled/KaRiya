package intents

import (
	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Global Keys Enforcement E2E Tests
//
// This test suite establishes the contract that ALL intents must follow for global key handling.
// It uses GenerateCV as the reference implementation (gold standard) and tests all other intents
// against the same requirements.
//
// Reference: bugs/bug-001-escape-key-navigation.md - Part 2: Enforcement Strategy

var _ = Describe("Global Keys Enforcement E2E", func() {

	// =========================================================================
	// CONTRACT COMPLIANCE - All Intents Must Follow Same Rules
	// =========================================================================
	//
	// All intents must:
	// 1. Cancel on escape from root state
	// 2. Go back on escape from intermediate states
	// 3. Ignore quit key (q) - quit only works from main menu
	// 4. Respond to help key (?) in all states
	// =========================================================================

	Describe("All Intents Contract Compliance", func() {

		Context("Root State Escape Behavior", func() {
			DescribeTable("should cancel intent on escape from root state",
				func(name string, setupFunc func() (Intent, error)) {
					intent, err := setupFunc()
					Expect(err).NotTo(HaveOccurred(), "%s: Failed to create intent", name)

					intent.Init()
					intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

					result := intent.Result()
					Expect(result).NotTo(BeNil(),
						"%s: Root state escape should produce result", name)
				},
				Entry("GenerateCV", "GenerateCV", func() (Intent, error) {
					return NewGenerateCVIntent(&GenerateCVContext{
						AvailableProfiles: []*CVProfile{{
							ID:             "default",
							Name:           "Default",
							TargetRole:     "staff",
							TargetAudience: "hiring_manager",
						}},
						Events: []*career.Event{
							fixtures.EventWith(uuid.New().String(), "Test", "", ""),
						},
					})
				}),
			)
		})

		Context("Quit Key Behavior", func() {
			DescribeTable("should cancel intent on 'q' key with cancelled status",
				func(name string, setupFunc func() (Intent, error)) {
					intent, err := setupFunc()
					Expect(err).NotTo(HaveOccurred())

					intent.Init()
					intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

					result := intent.Result()
					Expect(result).NotTo(BeNil(),
						"%s: 'q' key should cancel intent", name)
					Expect(result.Status).To(Equal(Cancelled),
						"%s: 'q' key should produce cancelled status", name)
				},
				Entry("GenerateCV", "GenerateCV", func() (Intent, error) {
					return NewGenerateCVIntent(&GenerateCVContext{
						AvailableProfiles: []*CVProfile{{
							ID:             "default",
							Name:           "Default",
							TargetRole:     "staff",
							TargetAudience: "hiring_manager",
						}},
						Events: []*career.Event{
							fixtures.EventWith(uuid.New().String(), "Test", "", ""),
						},
					})
				}),
			)
		})

		Context("Help Key Behavior", func() {
			DescribeTable("should toggle help on '?' in root state",
				func(name string, setupFunc func() (Intent, error)) {
					intent, err := setupFunc()
					Expect(err).NotTo(HaveOccurred())

					intent.Init()

					// Get BaseIntent to check help state.
					var baseIntent *BaseIntent
					switch i := intent.(type) {
					case *GenerateCVIntent:
						baseIntent = i.BaseIntent
					}

					if baseIntent != nil {
						helpBefore := baseIntent.IsHelpVisible()
						intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
						Expect(baseIntent.IsHelpVisible()).NotTo(Equal(helpBefore),
							"%s: Help key should toggle help modal", name)
					}
				},
				Entry("GenerateCV", "GenerateCV", func() (Intent, error) {
					return NewGenerateCVIntent(&GenerateCVContext{
						AvailableProfiles: []*CVProfile{{
							ID:             "default",
							Name:           "Default",
							TargetRole:     "staff",
							TargetAudience: "hiring_manager",
						}},
						Events: []*career.Event{
							fixtures.EventWith(uuid.New().String(), "Test", "", ""),
						},
					})
				}),
			)
		})
	})
})
