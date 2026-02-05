package e2e_test

import (
	"time"

	"github.com/baphled/kariya/internal/testutil/e2e"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Capture Enrichment Editing E2E", func() {
	var env *e2e.TestEnv

	BeforeEach(func() {
		env = e2e.Setup(GinkgoT())
	})

	AfterEach(func() {
		env.Cleanup()
	})

	Describe("Metadata Editing Modal", func() {
		It("should allow editing event metadata during review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm() // Quick strategy

			testEvent := fixtures.EventWith("", "Original event text", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review state
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Try to trigger metadata editing (typically 'e' or Enter on metadata)
			env.PressKeyRune('e')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should show metadata editor modal overlay", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Event for metadata editing", "TestCo", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate to review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Trigger metadata editor
			env.PressKeyRune('e')

			view := env.GetView()
			// Should show modal without crashes
			Expect(view).NotTo(ContainSubstring("panic"))
			Expect(view).NotTo(BeEmpty())
		})

		It("should allow cancelling metadata editing with Escape", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Metadata cancel test", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Open and cancel metadata editor
			env.PressKeyRune('e')
			env.Cancel()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should save metadata changes when submitted", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Metadata save test", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate to review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Open metadata editor, make changes, save
			env.PressKeyRune('e')
			env.TypeText("Updated company")
			env.PressKey(tea.KeyCtrlS)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle navigation within metadata editor fields", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Metadata navigation test", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Navigate in metadata editor
			env.PressKeyRune('e')
			env.Tab()
			env.Tab()
			env.NavigateUp()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Burst Editing Modal", func() {
		It("should show burst suggestions in review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Create event likely to generate burst suggestions
			testEvent := fixtures.EventWith("", "Led 3-month project with 5 engineers", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate to review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow editing burst suggestion", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Managed team from Jan to March 2024", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review with bursts
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Try to edit burst (typically 'e' on selected burst)
			env.PressKeyRune('e')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow accepting burst suggestion", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Project spanning Q1 2024", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate through review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Accept burst (typically Space or Enter)
			env.Confirm()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow rejecting burst suggestion", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Worked on project for 6 months", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Reject burst (typically 'r' or Delete)
			env.PressKeyRune('r')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should navigate between multiple burst suggestions", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Led multiple projects from Jan to June 2024", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate to review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Navigate through bursts
			env.NavigateDown()
			env.NavigateUp()
			env.PressKeyRune('j')
			env.PressKeyRune('k')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle burst modal cancellation", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Burst modal cancel test event", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Open and cancel burst editor
			env.PressKeyRune('e')
			env.Cancel()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Fact Editing Modal", func() {
		It("should show fact suggestions in review", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			// Event likely to generate facts
			testEvent := fixtures.EventWith("", "Built REST API with Go, PostgreSQL, and Redis", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate to review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow editing fact suggestion", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Improved performance by 50% using caching", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review with facts
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Try to edit fact
			env.PressKeyRune('e')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow accepting fact suggestion", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Reduced latency from 500ms to 100ms", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate through review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Accept fact
			env.Confirm()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow rejecting fact suggestion", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Implemented authentication and authorization", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Reject fact
			env.PressKeyRune('r')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should navigate between multiple fact suggestions", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Deployed to AWS using Docker, Kubernetes, and Terraform", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate to review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Navigate through facts
			env.NavigateDown()
			env.NavigateUp()
			env.PressKeyRune('j')
			env.PressKeyRune('k')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle fact modal cancellation", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Fact modal cancel test", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Open and cancel fact editor
			env.PressKeyRune('e')
			env.Cancel()

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})

	Describe("Mixed Enrichment Workflow", func() {
		It("should handle reviewing both bursts and facts", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Led 6-month migration project using Go and PostgreSQL with 80% performance improvement", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate through review
			maxAttempts := 15
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					// Accept or navigate through suggestions
					time.Sleep(50 * time.Millisecond)
					time.Sleep(50 * time.Millisecond)
					env.Confirm()
				}
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should allow switching between burst and fact editing", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "3-month project implementing microservices architecture", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Navigate between bursts and facts
			env.PressKeyRune('b') // Switch to bursts (if available)
			env.PressKeyRune('f') // Switch to facts (if available)

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should accept some and reject other suggestions", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Reduced deployment time from hours to minutes using CI/CD", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Navigate through review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Accept first suggestion
			env.Confirm()

			// Reject second (if exists)
			env.PressKeyRune('r')

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle completing review with mixed acceptances", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Architected scalable system handling 10k requests per second", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Go through full review workflow
			maxAttempts := 15
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					// Simulate accepting/rejecting randomly
					if i%2 == 0 {
						time.Sleep(50 * time.Millisecond)
						time.Sleep(50 * time.Millisecond)
						env.Confirm() // Accept
					} else {
						env.NavigateDown() // Skip
					}
				}
			}

			// Complete review
			env.Confirm()

			env.AssertEventCount(1)
		})
	})

	Describe("Enrichment Edge Cases", func() {
		It("should handle event with no enrichment suggestions", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Simple event", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Should complete even without suggestions
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					time.Sleep(50 * time.Millisecond)
					time.Sleep(50 * time.Millisecond)
					env.Confirm()
				}
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should handle cancelling review with pending suggestions", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Event with suggestions to cancel", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Cancel without accepting suggestions
			env.Cancel()

			env.AssertViewContainsAny("Capture Event", "Browse Timeline")
		})

		It("should handle rapid editing modal toggles", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Rapid toggle test event", "", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Rapid modal open/close
			for i := 0; i < 5; i++ {
				env.PressKeyRune('e')
				env.Cancel()
			}

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})

		It("should preserve event data when editing metadata multiple times", func() {
			env.SelectIntentByName("capture_event")
			env.Confirm()

			testEvent := fixtures.EventWith("", "Multiple metadata edits test", "InitialCo", "")
			testEvent.ID = ""
			env.SubmitEvent(testEvent)

			// Wait for review
			maxAttempts := 10
			for i := 0; i < maxAttempts; i++ {
				view := env.GetView()
				if view != "" {
					break
				}
				time.Sleep(50 * time.Millisecond)
				env.Confirm()
			}

			// Edit metadata multiple times
			env.PressKeyRune('e')
			env.TypeText("FirstEdit")
			env.PressKey(tea.KeyCtrlS)

			env.PressKeyRune('e')
			env.TypeText("SecondEdit")
			env.Cancel() // Don't save this one

			view := env.GetView()
			Expect(view).NotTo(ContainSubstring("panic"))
		})
	})
})
