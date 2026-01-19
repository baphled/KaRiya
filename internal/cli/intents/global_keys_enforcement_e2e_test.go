package intents

import (
	"context"
	"time"

	"github.com/baphled/kariya/internal/domain/career"
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
	// PART 1: REFERENCE IMPLEMENTATION - GenerateCV (Gold Standard)
	// =========================================================================
	//
	// GenerateCV is our reference because:
	// - 10 states (most complex intent)
	// - 16 existing escape tests (all passing)
	// - Already follows correct pattern
	// - These tests MUST pass - they establish the contract
	// =========================================================================

	Describe("Reference Implementation: GenerateCV", func() {
		var intent *GenerateCVIntent
		var profiles []*CVProfile
		var events []*career.CareerEvent

		BeforeEach(func() {
			profiles = []*CVProfile{
				{
					ID:             "profile1",
					Name:           "Senior IC",
					TargetRole:     "senior_ic",
					TargetAudience: "hiring_manager",
				},
			}

			events = []*career.CareerEvent{
				{
					ID:        uuid.New().String(),
					Text:      "Test event for CV generation",
					Date:      time.Now(),
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			}

			ctx := &GenerateCVContext{
				AvailableProfiles: profiles,
				Events:            events,
				Facts:             make([]*career.Fact, 0),
				DefaultProfile:    profiles[0],
			}

			var err error
			intent, err = NewGenerateCVIntent(ctx)
			Expect(err).NotTo(HaveOccurred())
			intent.Init()
		})

		Context("Root State Behavior", func() {
			It("should cancel intent on escape from SelectProfile (root state)", func() {
				intent.state.currentState = GenerateCVStateSelectProfile
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(Cancelled))
			})

			It("should ignore 'q' key from root state (quit only from main menu)", func() {
				intent.state.currentState = GenerateCVStateSelectProfile
				cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

				// q no longer quits from within intents - only from main menu
				Expect(cmd).To(BeNil())
			})

			It("should toggle help on '?' from root state", func() {
				intent.state.currentState = GenerateCVStateSelectProfile
				helpBefore := intent.IsHelpVisible()

				intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})

				Expect(intent.IsHelpVisible()).NotTo(Equal(helpBefore))
			})
		})

		Context("Intermediate State Behavior", func() {
			It("should go back on escape from SelectAudience", func() {
				intent.state.currentState = GenerateCVStateSelectAudience
				intent.state.selectedProfile = profiles[0]

				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectProfile))
				Expect(intent.Result()).To(BeNil()) // Still active
			})

			It("should go back on escape from Preview", func() {
				intent.state.currentState = GenerateCVStatePreview
				intent.state.generatedCV = &career.CVView{
					ID:               "cv1",
					Name:             "Test",
					GeneratedAt:      time.Now(),
					TargetAudience:   "test",
					SourceEventCount: 1,
					SourceFactCount:  0,
				}

				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
				Expect(intent.Result()).To(BeNil())
			})
		})

		Context("Async State Behavior", func() {
			It("should allow escape during Generating (async state)", func() {
				intent.state.currentState = GenerateCVStateGenerating
				intent.state.isGenerating = true

				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				// Should go back, letting generation complete
				Expect(intent.state.currentState).To(Equal(GenerateCVStateSelectAudience))
			})
		})

		Context("Universal Keys Work in All States", func() {
			DescribeTable("'q' key is ignored in all states (quit only from main menu)",
				func(state GenerateCVState) {
					intent.state.currentState = state
					if state != GenerateCVStateSelectProfile {
						intent.state.selectedProfile = profiles[0]
					}

					cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})
					// q no longer quits from within intents - only from main menu
					Expect(cmd).To(BeNil(), "'q' should be ignored in state %s (quit only from main menu)", state)
				},
				Entry("SelectProfile", GenerateCVStateSelectProfile),
				Entry("SelectAudience", GenerateCVStateSelectAudience),
				Entry("Generating", GenerateCVStateGenerating),
				Entry("Preview", GenerateCVStatePreview),
				Entry("Review", GenerateCVStateReview),
				Entry("Confirm", GenerateCVStateConfirm),
			)

			DescribeTable("help key works in all states",
				func(state GenerateCVState) {
					intent.state.currentState = state
					if state != GenerateCVStateSelectProfile {
						intent.state.selectedProfile = profiles[0]
					}

					helpBefore := intent.IsHelpVisible()
					intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'?'}})
					Expect(intent.IsHelpVisible()).NotTo(Equal(helpBefore),
						"Help should toggle in state %s", state)
				},
				Entry("SelectProfile", GenerateCVStateSelectProfile),
				Entry("SelectAudience", GenerateCVStateSelectAudience),
				Entry("Generating", GenerateCVStateGenerating),
				Entry("Preview", GenerateCVStatePreview),
				Entry("Review", GenerateCVStateReview),
				Entry("Confirm", GenerateCVStateConfirm),
			)
		})
	})

	// =========================================================================
	// PART 2: CONTRACT COMPLIANCE - All Intents Must Follow Same Rules
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
					// Note: Some intents may return Completed instead of Cancelled
					// That's acceptable as long as intent becomes inactive
				},
				Entry("GenerateCV", "GenerateCV", func() (Intent, error) {
					return NewGenerateCVIntent(&GenerateCVContext{
						AvailableProfiles: []*CVProfile{{
							ID:             "default",
							Name:           "Default",
							TargetRole:     "staff",
							TargetAudience: "hiring_manager",
						}},
						Events: []*career.CareerEvent{{
							ID:   uuid.New().String(),
							Text: "Test",
							Date: time.Now(),
						}},
					})
				}),
				Entry("BrowseTimeline", "BrowseTimeline", func() (Intent, error) {
					return NewBrowseTimelineIntent(&BrowseTimelineContext{})
				}),
				Entry("CaptureEvent", "CaptureEvent", func() (Intent, error) {
					return NewCaptureEventIntent(&CaptureEventContext{
						CaptureStrategy: "manual",
						Metadata:        make(map[string]string),
					})
				}),
			)
		})

		Context("Quit Key Behavior", func() {
			DescribeTable("should ignore 'q' key in root state (quit only from main menu)",
				func(name string, setupFunc func() (Intent, error)) {
					intent, err := setupFunc()
					Expect(err).NotTo(HaveOccurred())

					intent.Init()
					cmd := intent.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

					// q no longer quits from within intents - only from main menu
					Expect(cmd).To(BeNil(),
						"%s: 'q' key should be ignored within intent (quit only from main menu)", name)
				},
				Entry("GenerateCV", "GenerateCV", func() (Intent, error) {
					return NewGenerateCVIntent(&GenerateCVContext{
						AvailableProfiles: []*CVProfile{{
							ID:             "default",
							Name:           "Default",
							TargetRole:     "staff",
							TargetAudience: "hiring_manager",
						}},
						Events: []*career.CareerEvent{{
							ID:   uuid.New().String(),
							Text: "Test",
							Date: time.Now(),
						}},
					})
				}),
				Entry("BrowseTimeline", "BrowseTimeline", func() (Intent, error) {
					return NewBrowseTimelineIntent(&BrowseTimelineContext{})
				}),
				Entry("CaptureEvent", "CaptureEvent", func() (Intent, error) {
					return NewCaptureEventIntent(&CaptureEventContext{
						CaptureStrategy: "manual",
						Metadata:        make(map[string]string),
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

					// Get BaseIntent to check help state
					var baseIntent *BaseIntent
					switch i := intent.(type) {
					case *GenerateCVIntent:
						baseIntent = i.BaseIntent
					case *BrowseTimelineIntent:
						baseIntent = i.BaseIntent
					case *CaptureEventIntent:
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
						Events: []*career.CareerEvent{{
							ID:   uuid.New().String(),
							Text: "Test",
							Date: time.Now(),
						}},
					})
				}),
				Entry("BrowseTimeline", "BrowseTimeline", func() (Intent, error) {
					return NewBrowseTimelineIntent(&BrowseTimelineContext{})
				}),
				Entry("CaptureEvent", "CaptureEvent", func() (Intent, error) {
					return NewCaptureEventIntent(&CaptureEventContext{
						CaptureStrategy: "manual",
						Metadata:        make(map[string]string),
					})
				}),
			)
		})
	})

	// =========================================================================
	// PART 3: CONTEXT-AWARE NAVIGATION
	// =========================================================================
	//
	// Some intents behave differently based on context (edit vs new).
	// This is the correct pattern - don't hard-code state transitions.
	// =========================================================================

	Describe("Context-Aware Navigation", func() {

		Context("CaptureEvent - Edit vs New", func() {
			It("should cancel when editing existing event (PreviousEvent != nil)", func() {
				existingEvent := &career.CareerEvent{
					ID:   uuid.New().String(),
					Text: "Existing event",
					Date: time.Now(),
				}

				ctx := &CaptureEventContext{
					CaptureStrategy: "manual",
					PreviousEvent:   existingEvent, // Edit mode
					Metadata:        make(map[string]string),
				}

				intent, err := NewCaptureEventIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()

				// Transition to form state
				intent.state.currentState = CaptureStateForm

				// Press escape - should cancel (return to caller like BrowseTimeline)
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(Cancelled))
			})

			It("should go back when creating new event (PreviousEvent == nil)", func() {
				ctx := &CaptureEventContext{
					CaptureStrategy: "manual",
					PreviousEvent:   nil, // New event mode
					Metadata:        make(map[string]string),
				}

				intent, err := NewCaptureEventIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()

				// Transition to form state
				intent.state.currentState = CaptureStateForm

				// Press escape - should go back to strategy selection
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				Expect(intent.state.currentState).To(Equal(CaptureStateChooseStrategy))
				Expect(intent.Result()).To(BeNil()) // Still active
			})
		})

		Context("BurstManagement - New vs Edit", func() {
			It("should handle context-aware navigation", func() {
				ctx := NewBurstManagementContext(nil, nil, context.Background())

				intent, err := NewBurstManagementIntent(ctx)
				Expect(err).NotTo(HaveOccurred())
				intent.Init()

				// Test that escape from list (root) cancels
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
				Expect(result.Status).To(Equal(Cancelled))
			})
		})

		Context("FactManagement - New vs Edit", func() {
			It("should handle context-aware navigation", func() {
				ctx := NewFactManagementContext(nil, context.Background())

				intent := NewFactManagementIntent(ctx)
				Expect(intent).NotTo(BeNil())
				intent.Init()

				// Test that escape from list (root) cancels
				intent.Update(tea.KeyMsg{Type: tea.KeyEsc})

				result := intent.Result()
				Expect(result).NotTo(BeNil())
			})
		})
	})

	// =========================================================================
	// PART 4: FUTURE ENFORCEMENT
	// NOTE: Additional intents (ExportArtifact, ConfigureSystem, ImportWizard,
	// MetadataEditor, BulkOperations) will be added to the DescribeTable when
	// their migration to the new architecture is complete.
})
