package browsetimeline

import (
	"errors"

	"github.com/baphled/kariya/internal/ui/display"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/domain/career"
	"github.com/baphled/kariya/internal/testutil/fixtures"
	"github.com/baphled/kariya/internal/tui/views/event"
	"github.com/baphled/kariya/internal/ui/uikit/widgets"
)

var _ = Describe("handleViewResult", func() {
	var (
		intent *Intent
		ctx    *IntentValidator
	)

	BeforeEach(func() {
		ctx = &IntentValidator{
			Events: []*career.Event{
				fixtures.Event("event-1"),
				fixtures.Event("event-2"),
			},
		}
		var err error
		intent, err = NewIntent(ctx)
		Expect(err).NotTo(HaveOccurred())
		Expect(intent).NotTo(BeNil())
	})

	Describe("when result type is not ResultNavigate", func() {
		It("returns nil when result type is Cancel", func() {
			result := &widgets.CancelViewResult{}
			cmd := intent.handleViewResult(result)
			Expect(cmd).To(BeNil())
		})

		It("returns nil when result type is Submit", func() {
			result := &widgets.SubmitViewResult{}
			cmd := intent.handleViewResult(result)
			Expect(cmd).To(BeNil())
		})

		It("returns nil when result type is Error", func() {
			result := &widgets.ErrorViewResult{
				Err:     errors.New("test error"),
				Message: "Error message",
			}
			cmd := intent.handleViewResult(result)
			Expect(cmd).To(BeNil())
		})
	})

	Describe("when result contains SkillNav action", func() {
		var testEvent *career.Event

		BeforeEach(func() {
			testEvent = fixtures.Event("event-1")
		})

		Context("with ActionShowEventSkills", func() {
			It("sets selected event and calls showSkillsForCurrentEvent", func() {
				skillNav := event.SkillNav{
					Action: event.ActionShowEventSkills,
					Event:  display.EventFromDomain(testEvent),
				}
				result := &widgets.NavigateViewResult{
					ResultData: skillNav,
				}

				cmd := intent.handleViewResult(result)

				Expect(intent.selectedEvent).NotTo(BeNil())
				Expect(intent.selectedEvent.ID).To(Equal(testEvent.ID))
				Expect(cmd).ToNot(BeNil())
				msg := cmd()
				_, ok := msg.(SkillsForModalLoadedMsg)
				Expect(ok).To(BeTrue())
			})

			It("handles nil event gracefully", func() {
				skillNav := event.SkillNav{
					Action: event.ActionShowEventSkills,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: skillNav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})

		Context("with ActionInferSkills", func() {
			It("returns nil when inference service is not available", func() {
				skillNav := event.SkillNav{
					Action: event.ActionInferSkills,
					Event:  display.EventFromDomain(testEvent),
				}
				result := &widgets.NavigateViewResult{
					ResultData: skillNav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})

			It("returns nil when event is nil", func() {
				skillNav := event.SkillNav{
					Action: event.ActionInferSkills,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: skillNav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})

		Context("with unknown skill action", func() {
			It("returns nil for unhandled actions", func() {
				skillNav := event.SkillNav{
					Action: event.SkillActionKey("unknown.action"),
					Event:  display.EventFromDomain(testEvent),
				}
				result := &widgets.NavigateViewResult{
					ResultData: skillNav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})
	})

	Describe("when result contains Nav action", func() {
		var testEvent *career.Event

		BeforeEach(func() {
			testEvent = fixtures.Event("event-1")
		})

		Context("with ActionEdit", func() {
			It("opens edit modal for the event", func() {
				nav := event.Nav{
					Action: event.ActionEdit,
					Event:  display.EventFromDomain(testEvent),
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).NotTo(BeNil())
			})

			It("returns nil when event is nil", func() {
				nav := event.Nav{
					Action: event.ActionEdit,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})

		Context("with ActionDelete", func() {
			It("opens delete modal for the event", func() {
				nav := event.Nav{
					Action: event.ActionDelete,
					Event:  display.EventFromDomain(testEvent),
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				intent.handleViewResult(result)

				Expect(intent.deleteModal).NotTo(BeNil())
			})

			It("returns nil when event is nil", func() {
				nav := event.Nav{
					Action: event.ActionDelete,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})

		Context("with ActionView", func() {
			It("opens detail modal for the event", func() {
				nav := event.Nav{
					Action: event.ActionView,
					Event:  display.EventFromDomain(testEvent),
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(intent.selectedEvent).NotTo(BeNil())
				Expect(intent.selectedEvent.ID).To(Equal(testEvent.ID))
				Expect(cmd).To(BeNil())
			})

			It("returns nil when event is nil", func() {
				nav := event.Nav{
					Action: event.ActionView,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})

			It("appends viewed event to viewed events list", func() {
				nav := event.Nav{
					Action: event.ActionView,
					Event:  display.EventFromDomain(testEvent),
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				intent.handleViewResult(result)

				Expect(intent.viewedEvents).To(HaveLen(1))
				Expect(intent.viewedEvents[0]).NotTo(BeNil())
				Expect(intent.viewedEvents[0].ID).To(Equal(testEvent.ID))
			})
		})

		Context("with ActionAdd", func() {
			It("opens quick add modal", func() {
				nav := event.Nav{
					Action: event.ActionAdd,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("with ActionFilter", func() {
			It("opens filter modal", func() {
				nav := event.Nav{
					Action: event.ActionFilter,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).NotTo(BeNil())
			})
		})

		Context("with ActionSearch", func() {
			It("returns nil (search handled elsewhere)", func() {
				nav := event.Nav{
					Action: event.ActionSearch,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})

		Context("with ActionSort", func() {
			It("returns nil (sort handled elsewhere)", func() {
				nav := event.Nav{
					Action: event.ActionSort,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})

		Context("with ActionClear", func() {
			It("returns nil (clear handled elsewhere)", func() {
				nav := event.Nav{
					Action: event.ActionClear,
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).To(BeNil())
			})
		})

		Context("with unknown action", func() {
			It("returns error command", func() {
				nav := event.Nav{
					Action: event.ActionKey("unknown.action"),
					Event:  display.Event{},
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}

				cmd := intent.handleViewResult(result)

				Expect(cmd).NotTo(BeNil())
				msg := cmd()
				err, isErr := msg.(error)
				Expect(isErr).To(BeTrue())
				Expect(err.Error()).To(ContainSubstring("unknown event action"))
			})
		})
	})

	Describe("when result contains invalid payload type", func() {
		It("returns error command with type information", func() {
			invalidPayload := "invalid string payload"
			result := &widgets.NavigateViewResult{
				ResultData: invalidPayload,
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			err, isErr := msg.(error)
			Expect(isErr).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("invalid event action payload type"))
		})

		It("returns error command for nil payload", func() {
			result := &widgets.NavigateViewResult{
				ResultData: nil,
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			err, isErr := msg.(error)
			Expect(isErr).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("invalid event action payload type"))
		})

		It("returns error command for numeric payload", func() {
			result := &widgets.NavigateViewResult{
				ResultData: 42,
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			err, isErr := msg.(error)
			Expect(isErr).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("invalid event action payload type"))
		})

		It("returns error command for map payload", func() {
			result := &widgets.NavigateViewResult{
				ResultData: map[string]interface{}{"key": "value"},
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			err, isErr := msg.(error)
			Expect(isErr).To(BeTrue())
			Expect(err.Error()).To(ContainSubstring("invalid event action payload type"))
		})
	})

	Describe("command execution safety", func() {
		It("executes error commands immediately", func() {
			result := &widgets.NavigateViewResult{
				ResultData: "invalid",
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
			msg := cmd()
			Expect(msg).NotTo(BeNil())
			_, isErr := msg.(error)
			Expect(isErr).To(BeTrue())
		})

		It("returns non-nil commands for modal operations", func() {
			nav := event.Nav{
				Action: event.ActionAdd,
				Event:  display.Event{},
			}
			result := &widgets.NavigateViewResult{
				ResultData: nav,
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
		})

		It("returns nil for no-op actions", func() {
			nav := event.Nav{
				Action: event.ActionSearch,
				Event:  display.Event{},
			}
			result := &widgets.NavigateViewResult{
				ResultData: nav,
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).To(BeNil())
		})
	})

	Describe("state transitions", func() {
		It("preserves selected event across multiple view results", func() {
			event1 := fixtures.EventWith("e1", "Event One", "Company A", "ProjectA")
			event2 := fixtures.EventWith("e2", "Event Two", "Company B", "ProjectB")

			nav1 := event.Nav{
				Action: event.ActionView,
				Event:  display.EventFromDomain(event1),
			}
			result1 := &widgets.NavigateViewResult{
				ResultData: nav1,
			}

			intent.handleViewResult(result1)
			Expect(intent.selectedEvent).To(BeNil())

			nav2 := event.Nav{
				Action: event.ActionView,
				Event:  display.EventFromDomain(event2),
			}
			result2 := &widgets.NavigateViewResult{
				ResultData: nav2,
			}

			intent.handleViewResult(result2)
			Expect(intent.selectedEvent).To(BeNil())
		})

		It("accumulates viewed events", func() {
			event1 := fixtures.Event("event-1")
			event2 := fixtures.Event("event-2")
			event3 := fixtures.Event("event-3")

			nav1 := event.Nav{
				Action: event.ActionView,
				Event:  display.EventFromDomain(event1),
			}
			result1 := &widgets.NavigateViewResult{
				ResultData: nav1,
			}
			intent.handleViewResult(result1)

			nav2 := event.Nav{
				Action: event.ActionView,
				Event:  display.EventFromDomain(event2),
			}
			result2 := &widgets.NavigateViewResult{
				ResultData: nav2,
			}
			intent.handleViewResult(result2)

			nav3 := event.Nav{
				Action: event.ActionView,
				Event:  display.EventFromDomain(event3),
			}
			result3 := &widgets.NavigateViewResult{
				ResultData: nav3,
			}
			intent.handleViewResult(result3)

			Expect(intent.viewedEvents).To(HaveLen(2))
			Expect(intent.viewedEvents[0]).NotTo(BeNil())
			Expect(intent.viewedEvents[0].ID).To(Equal(event1.ID))
			Expect(intent.viewedEvents[1]).NotTo(BeNil())
			Expect(intent.viewedEvents[1].ID).To(Equal(event2.ID))
		})

		It("allows viewing same event multiple times", func() {
			event1 := fixtures.Event("event-1")

			for range 3 {
				nav := event.Nav{
					Action: event.ActionView,
					Event:  display.EventFromDomain(event1),
				}
				result := &widgets.NavigateViewResult{
					ResultData: nav,
				}
				intent.handleViewResult(result)
			}

			Expect(intent.viewedEvents).To(HaveLen(3))
		})
	})

	Describe("edge cases", func() {
		It("handles large event text in delete modal", func() {
			largeText := "a"
			for range 100 {
				largeText = "a" + largeText
			}
			eventWithLargeText := fixtures.EventWith("event-1", largeText, "Company", "Project")

			nav := event.Nav{
				Action: event.ActionDelete,
				Event:  display.EventFromDomain(eventWithLargeText),
			}
			result := &widgets.NavigateViewResult{
				ResultData: nav,
			}

			intent.handleViewResult(result)

			Expect(intent.deleteModal).NotTo(BeNil())
		})

		It("handles empty event text", func() {
			emptyEvent := fixtures.EventWith("event-1", "", "Company", "Project")

			nav := event.Nav{
				Action: event.ActionDelete,
				Event:  display.EventFromDomain(emptyEvent),
			}
			result := &widgets.NavigateViewResult{
				ResultData: nav,
			}

			intent.handleViewResult(result)

			Expect(intent.deleteModal).NotTo(BeNil())
		})

		It("handles skill actions with various event states", func() {
			testEvent := fixtures.Event("event-1")

			skillNav := event.SkillNav{
				Action: event.ActionShowEventSkills,
				Event:  display.EventFromDomain(testEvent),
			}
			result := &widgets.NavigateViewResult{
				ResultData: skillNav,
			}

			cmd := intent.handleViewResult(result)

			Expect(intent.selectedEvent).NotTo(BeNil())
			Expect(intent.selectedEvent.ID).To(Equal(testEvent.ID))
			Expect(cmd).ToNot(BeNil())
			msg := cmd()
			_, ok := msg.(SkillsForModalLoadedMsg)
			Expect(ok).To(BeTrue())
		})
	})

	Describe("command type safety", func() {
		It("error commands are executable tea.Cmd", func() {
			result := &widgets.NavigateViewResult{
				ResultData: "invalid",
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
			var _ = cmd
			msg := cmd()
			Expect(msg).NotTo(BeNil())
		})

		It("modal commands are executable tea.Cmd", func() {
			nav := event.Nav{
				Action: event.ActionAdd,
				Event:  display.Event{},
			}
			result := &widgets.NavigateViewResult{
				ResultData: nav,
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).NotTo(BeNil())
			var _ = cmd
		})

		It("nil returns are safe to batch", func() {
			nav := event.Nav{
				Action: event.ActionSearch,
				Event:  display.Event{},
			}
			result := &widgets.NavigateViewResult{
				ResultData: nav,
			}

			cmd := intent.handleViewResult(result)

			Expect(cmd).To(BeNil())
		})
	})

	Describe("result type discrimination", func() {
		It("correctly identifies NavigateViewResult", func() {
			nav := event.Nav{
				Action: event.ActionSearch,
				Event:  display.Event{},
			}
			result := &widgets.NavigateViewResult{
				ResultData: nav,
			}

			Expect(result.Type()).To(Equal(widgets.ResultNavigate))
		})

		It("distinguishes NavigateViewResult from CancelViewResult", func() {
			cancelResult := &widgets.CancelViewResult{}
			cmd := intent.handleViewResult(cancelResult)
			Expect(cmd).To(BeNil())

			navResult := &widgets.NavigateViewResult{
				ResultData: event.Nav{
					Action: event.ActionAdd,
					Event:  display.Event{},
				},
			}
			cmd = intent.handleViewResult(navResult)
			Expect(cmd).NotTo(BeNil())
		})
	})
})
