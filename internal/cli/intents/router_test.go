package intents_test

import (
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/themes"
	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DefaultIntentRouter", func() {
	var router *intents.DefaultIntentRouter

	BeforeEach(func() {
		router = intents.NewDefaultIntentRouter()
	})

	Describe("RegisterIntent", func() {
		It("should register an intent successfully", func() {
			factory := func() intents.Intent { return intents.NewMockIntent() }
			err := router.RegisterIntent("test_intent", factory)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should return error when registering duplicate intent", func() {
			factory := func() intents.Intent { return intents.NewMockIntent() }
			_ = router.RegisterIntent("test_intent", factory)
			err := router.RegisterIntent("test_intent", factory)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("ActivateIntent", func() {
		BeforeEach(func() {
			factory := func() intents.Intent { return intents.NewMockIntent() }
			_ = router.RegisterIntent("test_intent", factory)
		})

		It("should activate registered intent", func() {
			_, err := router.ActivateIntent("test_intent", nil)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should set active intent", func() {
			_, _ = router.ActivateIntent("test_intent", nil)
			active := router.GetActiveIntent()
			Expect(active).NotTo(BeNil())
		})

		It("should call Init on the intent", func() {
			_, _ = router.ActivateIntent("test_intent", nil)
			active := router.GetActiveIntent()
			mockIntent := active.(*intents.MockIntent)
			Expect(mockIntent.InitCalled).To(BeTrue())
		})

		It("should return error for nonexistent intent", func() {
			_, err := router.ActivateIntent("nonexistent", nil)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("GetActiveIntent", func() {
		It("should return nil when no intent is active", func() {
			Expect(router.GetActiveIntent()).To(BeNil())
		})

		It("should return active intent after activation", func() {
			factory := func() intents.Intent { return intents.NewMockIntent() }
			_ = router.RegisterIntent("test_intent", factory)
			_, _ = router.ActivateIntent("test_intent", nil)
			Expect(router.GetActiveIntent()).NotTo(BeNil())
		})
	})

	Describe("HandleMessage", func() {
		Context("with active intent", func() {
			BeforeEach(func() {
				factory := func() intents.Intent { return intents.NewMockIntent() }
				_ = router.RegisterIntent("test_intent", factory)
				_, _ = router.ActivateIntent("test_intent", nil)
			})

			It("should forward message to active intent", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
				_, _ = router.HandleMessage(msg)
				active := router.GetActiveIntent().(*intents.MockIntent)
				Expect(active.UpdateCalled).To(Equal(1))
			})

			It("should return nil result when intent hasn't completed", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
				_, result := router.HandleMessage(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("without active intent", func() {
			It("should return nil", func() {
				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
				_, result := router.HandleMessage(msg)
				Expect(result).To(BeNil())
			})
		})

		Context("with completed intent", func() {
			It("should return intent result", func() {
				mockIntent := intents.NewMockIntent()
				factory := func() intents.Intent { return mockIntent }
				_ = router.RegisterIntent("test_intent", factory)
				_, _ = router.ActivateIntent("test_intent", nil)

				expectedResult := intents.NewCompletedResult[interface{}]("test data")
				mockIntent.SetResult(expectedResult)

				msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}}
				_, result := router.HandleMessage(msg)

				Expect(result).NotTo(BeNil())
				intentResult, ok := result.(*intents.IntentResult[interface{}])
				Expect(ok).To(BeTrue())
				Expect(intentResult.Status).To(Equal(intents.Completed))
			})
		})
	})

	Describe("View", func() {
		Context("with active intent", func() {
			It("should return view from active intent", func() {
				factory := func() intents.Intent { return intents.NewMockIntent() }
				_ = router.RegisterIntent("test_intent", factory)
				_, _ = router.ActivateIntent("test_intent", nil)

				view := router.View()
				Expect(view).To(Equal("Mock Intent View"))
			})

			It("should call View on the intent", func() {
				factory := func() intents.Intent { return intents.NewMockIntent() }
				_ = router.RegisterIntent("test_intent", factory)
				_, _ = router.ActivateIntent("test_intent", nil)

				_ = router.View()
				active := router.GetActiveIntent().(*intents.MockIntent)
				Expect(active.ViewCalled).To(BeTrue())
			})
		})

		Context("without active intent", func() {
			It("should return default message", func() {
				view := router.View()
				Expect(view).To(Equal("No active intent"))
			})
		})
	})

	Describe("Back", func() {
		Context("with history", func() {
			BeforeEach(func() {
				factory1 := func() intents.Intent { return intents.NewMockIntent() }
				factory2 := func() intents.Intent { return intents.NewMockIntent() }
				_ = router.RegisterIntent("intent1", factory1)
				_ = router.RegisterIntent("intent2", factory2)
				_, _ = router.ActivateIntent("intent1", nil)
				_, _ = router.ActivateIntent("intent2", nil)
			})

			It("should go back without error", func() {
				_, err := router.Back()
				Expect(err).NotTo(HaveOccurred())
			})

			It("should reduce history depth", func() {
				_, _ = router.Back()
				Expect(router.GetHistoryDepth()).To(Equal(1))
			})
		})

		Context("without history", func() {
			It("should return error", func() {
				_, err := router.Back()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetHistory", func() {
		It("should return empty history initially", func() {
			Expect(router.GetHistory()).To(BeEmpty())
		})

		It("should track navigation history", func() {
			factory1 := func() intents.Intent { return intents.NewMockIntent() }
			factory2 := func() intents.Intent { return intents.NewMockIntent() }
			_ = router.RegisterIntent("intent1", factory1)
			_ = router.RegisterIntent("intent2", factory2)
			_, _ = router.ActivateIntent("intent1", nil)
			_, _ = router.ActivateIntent("intent2", nil)

			history := router.GetHistory()
			Expect(history).To(HaveLen(1))
		})
	})

	Describe("GetHistoryDepth", func() {
		It("should return 0 when no intent active", func() {
			Expect(router.GetHistoryDepth()).To(Equal(0))
		})

		It("should return 1 after first activation", func() {
			factory := func() intents.Intent { return intents.NewMockIntent() }
			_ = router.RegisterIntent("intent1", factory)
			_, _ = router.ActivateIntent("intent1", nil)
			Expect(router.GetHistoryDepth()).To(Equal(1))
		})

		It("should return 2 after second activation", func() {
			factory1 := func() intents.Intent { return intents.NewMockIntent() }
			factory2 := func() intents.Intent { return intents.NewMockIntent() }
			_ = router.RegisterIntent("intent1", factory1)
			_ = router.RegisterIntent("intent2", factory2)
			_, _ = router.ActivateIntent("intent1", nil)
			_, _ = router.ActivateIntent("intent2", nil)
			Expect(router.GetHistoryDepth()).To(Equal(2))
		})
	})

	Describe("Theme Management", func() {
		Describe("ThemeManager initialization", func() {
			It("should be initialized by default", func() {
				tm := router.GetThemeManager()
				Expect(tm).NotTo(BeNil())
			})
		})

		Describe("SetThemeManager", func() {
			It("should set custom theme manager", func() {
				customTM := themes.NewThemeManager()
				router.SetThemeManager(customTM)
				Expect(router.GetThemeManager()).To(Equal(customTM))
			})
		})

		Describe("Theme", func() {
			It("should return active theme", func() {
				theme := router.Theme()
				Expect(theme).NotTo(BeNil())
				Expect(theme.Name()).To(Equal("default"))
			})

			It("should return nil when no theme manager", func() {
				router.SetThemeManager(nil)
				theme := router.Theme()
				Expect(theme).To(BeNil())
			})
		})

		Describe("Theme propagation to intents", func() {
			It("should propagate theme to activated intent", func() {
				mockIntent := intents.NewThemeAwareMockIntent()
				factory := func() intents.Intent { return mockIntent }
				_ = router.RegisterIntent("test_intent", factory)
				_, _ = router.ActivateIntent("test_intent", nil)

				active := router.GetActiveIntent()
				themeAware, ok := active.(*intents.ThemeAwareMockIntent)
				Expect(ok).To(BeTrue())
				Expect(themeAware.GetThemeManager()).NotTo(BeNil())
			})

			It("should propagate new theme manager to active intent", func() {
				mockIntent := intents.NewThemeAwareMockIntent()
				factory := func() intents.Intent { return mockIntent }
				_ = router.RegisterIntent("test_intent", factory)
				_, _ = router.ActivateIntent("test_intent", nil)

				newTM := themes.NewThemeManager()
				router.SetThemeManager(newTM)

				active := router.GetActiveIntent()
				themeAware, ok := active.(*intents.ThemeAwareMockIntent)
				Expect(ok).To(BeTrue())
				Expect(themeAware.GetThemeManager()).To(Equal(newTM))
			})
		})
	})

	Describe("Selection Preservation", func() {
		It("should preserve selection when navigating back", func() {
			intent1 := intents.NewMockIntentWithSelection()
			factory1 := func() intents.Intent { return intent1 }
			intent2 := intents.NewMockIntent()
			factory2 := func() intents.Intent { return intent2 }

			_ = router.RegisterIntent("intent1", factory1)
			_ = router.RegisterIntent("intent2", factory2)

			_, _ = router.ActivateIntent("intent1", nil)
			active1 := router.GetActiveIntent().(*intents.MockIntentWithSelection)
			active1.SetSelectedIndex(5)

			_, _ = router.ActivateIntent("intent2", nil)
			_, err := router.Back()
			Expect(err).NotTo(HaveOccurred())

			restoredIntent := router.GetActiveIntent().(*intents.MockIntentWithSelection)
			Expect(restoredIntent.GetSelectedIndex()).To(Equal(5))
			Expect(restoredIntent).To(Equal(intent1))
		})
	})

	Describe("Context-Aware Factory Registration", func() {
		Describe("RegisterIntentWithContext", func() {
			It("should register a context-aware factory", func() {
				factory := func(ctx map[string]interface{}) intents.Intent {
					return intents.NewMockIntentWithContext(ctx)
				}
				err := router.RegisterIntentWithContext("context_intent", factory)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should return error for duplicate registration", func() {
				factory := func(ctx map[string]interface{}) intents.Intent {
					return intents.NewMockIntentWithContext(ctx)
				}
				_ = router.RegisterIntentWithContext("context_intent", factory)
				err := router.RegisterIntentWithContext("context_intent", factory)
				Expect(err).To(HaveOccurred())
			})
		})

		Describe("ActivateIntent with context", func() {
			BeforeEach(func() {
				factory := func(ctx map[string]interface{}) intents.Intent {
					return intents.NewMockIntentWithContext(ctx)
				}
				_ = router.RegisterIntentWithContext("context_intent", factory)
			})

			It("should pass context to factory", func() {
				ctx := map[string]interface{}{
					"editMode": true,
					"eventID":  "event-123",
				}
				_, err := router.ActivateIntent("context_intent", ctx)
				Expect(err).NotTo(HaveOccurred())

				active := router.GetActiveIntent().(*intents.MockIntentWithContext)
				Expect(active.GetContextValue("editMode")).To(BeTrue())
				Expect(active.GetContextValue("eventID")).To(Equal("event-123"))
			})

			It("should work with nil context", func() {
				_, err := router.ActivateIntent("context_intent", nil)
				Expect(err).NotTo(HaveOccurred())

				active := router.GetActiveIntent().(*intents.MockIntentWithContext)
				Expect(active.GetContext()).To(BeNil())
			})

			It("should work with empty context", func() {
				_, err := router.ActivateIntent("context_intent", make(map[string]interface{}))
				Expect(err).NotTo(HaveOccurred())

				active := router.GetActiveIntent().(*intents.MockIntentWithContext)
				Expect(active.GetContext()).To(BeEmpty())
			})
		})

		Describe("Mixed registration", func() {
			It("should support both context-aware and context-less factories", func() {
				// Register context-less factory
				_ = router.RegisterIntent("simple_intent", func() intents.Intent {
					return intents.NewMockIntent()
				})

				// Register context-aware factory
				_ = router.RegisterIntentWithContext("context_intent", func(ctx map[string]interface{}) intents.Intent {
					return intents.NewMockIntentWithContext(ctx)
				})

				// Activate both
				_, err1 := router.ActivateIntent("simple_intent", nil)
				Expect(err1).NotTo(HaveOccurred())

				_, err2 := router.ActivateIntent("context_intent", map[string]interface{}{"key": "value"})
				Expect(err2).NotTo(HaveOccurred())

				active := router.GetActiveIntent().(*intents.MockIntentWithContext)
				Expect(active.GetContextValue("key")).To(Equal("value"))
			})
		})
	})
})
