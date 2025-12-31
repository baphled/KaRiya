package navigation_test

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/navigation"
)

func TestNavigation(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Navigation Registry Suite")
}

// Mock screen types for testing
const (
	HomeScreen         navigation.Screen = "home"
	ListScreen         navigation.Screen = "list"
	DetailScreen       navigation.Screen = "detail"
	EditScreen         navigation.Screen = "edit"
	ConfirmationScreen navigation.Screen = "confirmation"
	SuccessScreen      navigation.Screen = "success"
	ModalScreen        navigation.Screen = "modal"
	NoBackScreen       navigation.Screen = "noback"
)

// Mock model for testing
type mockModel struct {
	id string
}

func (m mockModel) Init() tea.Cmd {
	return nil
}

func (m mockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m mockModel) View() string {
	return "mock view"
}

var _ = Describe("NavigationRegistry", func() {
	var registry *navigation.NavigationRegistry

	BeforeEach(func() {
		registry = navigation.NewNavigationRegistry(nil)
	})

	Describe("NewNavigationRegistry", func() {
		It("creates a registry with default options", func() {
			Expect(registry).NotTo(BeNil())
		})

		It("creates a registry with custom options", func() {
			opts := &navigation.NavigationOptions{
				MaxHistorySize: 50,
				DefaultContext: context.Background(),
			}
			customRegistry := navigation.NewNavigationRegistry(opts)
			Expect(customRegistry).NotTo(BeNil())
		})
	})

	Describe("RegisterScreen", func() {
		It("registers a screen definition", func() {
			def := &navigation.ScreenDefinition{
				ID:          HomeScreen,
				Label:       "Home",
				CanGoBack:   false,
				HelpContext: "home",
			}

			err := registry.RegisterScreen(def)
			Expect(err).NotTo(HaveOccurred())

			retrieved, err := registry.GetScreen(HomeScreen)
			Expect(err).NotTo(HaveOccurred())
			Expect(retrieved.ID).To(Equal(HomeScreen))
			Expect(retrieved.Label).To(Equal("Home"))
		})

		It("returns error for nil definition", func() {
			err := registry.RegisterScreen(nil)
			Expect(err).To(MatchError(navigation.ErrInvalidScreenID))
		})

		It("returns error for empty screen ID", func() {
			def := &navigation.ScreenDefinition{
				Label: "Test",
			}

			err := registry.RegisterScreen(def)
			Expect(err).To(MatchError(navigation.ErrInvalidScreenID))
		})

		It("detects circular parent references", func() {
			// Register first screen
			err := registry.RegisterScreen(&navigation.ScreenDefinition{
				ID:    HomeScreen,
				Label: "Home",
			})
			Expect(err).NotTo(HaveOccurred())

			// Try to register screen with circular reference
			homeParent := HomeScreen
			def := &navigation.ScreenDefinition{
				ID:     HomeScreen,
				Label:  "Home Updated",
				Parent: &homeParent,
			}

			err = registry.RegisterScreen(def)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("circular"))
		})
	})

	Describe("RegisterScreens", func() {
		It("registers multiple screen definitions", func() {
			defs := []*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home"},
				{ID: ListScreen, Label: "List"},
				{ID: DetailScreen, Label: "Detail"},
			}

			err := registry.RegisterScreens(defs)
			Expect(err).NotTo(HaveOccurred())

			allScreens := registry.GetAllScreens()
			Expect(len(allScreens)).To(Equal(3))
		})

		It("stops on first error", func() {
			defs := []*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home"},
				{Label: "Invalid"}, // Missing ID
				{ID: ListScreen, Label: "List"},
			}

			err := registry.RegisterScreens(defs)
			Expect(err).To(HaveOccurred())

			// Only first screen should be registered
			allScreens := registry.GetAllScreens()
			Expect(len(allScreens)).To(Equal(1))
		})
	})

	Describe("GetScreen", func() {
		BeforeEach(func() {
			registry.RegisterScreen(&navigation.ScreenDefinition{
				ID:    HomeScreen,
				Label: "Home",
			})
		})

		It("retrieves a registered screen", func() {
			screen, err := registry.GetScreen(HomeScreen)
			Expect(err).NotTo(HaveOccurred())
			Expect(screen.ID).To(Equal(HomeScreen))
		})

		It("returns error for unregistered screen", func() {
			_, err := registry.GetScreen("nonexistent")
			Expect(err).To(MatchError(ContainSubstring("not found")))
		})
	})

	Describe("Navigate", func() {
		BeforeEach(func() {
			registry.RegisterScreen(&navigation.ScreenDefinition{
				ID:          HomeScreen,
				Label:       "Home",
				CanGoBack:   false,
				HelpContext: "home",
			})
			registry.RegisterScreen(&navigation.ScreenDefinition{
				ID:          ListScreen,
				Label:       "List",
				CanGoBack:   true,
				HelpContext: "list",
			})
		})

		It("navigates to a screen", func() {
			model := mockModel{id: "home"}
			ctx := map[string]interface{}{"key": "value"}

			err := registry.Navigate(HomeScreen, model, ctx)
			Expect(err).NotTo(HaveOccurred())

			current := registry.GetCurrent()
			Expect(current).NotTo(BeNil())
			Expect(current.Screen).To(Equal(HomeScreen))
			Expect(current.Context["key"]).To(Equal("value"))
		})

		It("pushes previous state to history", func() {
			// First navigation
			err := registry.Navigate(HomeScreen, mockModel{id: "home"}, nil)
			Expect(err).NotTo(HaveOccurred())

			// Second navigation
			err = registry.Navigate(ListScreen, mockModel{id: "list"}, nil)
			Expect(err).NotTo(HaveOccurred())

			history := registry.GetHistory()
			Expect(len(history)).To(Equal(1))
			Expect(history[0].Screen).To(Equal(HomeScreen))

			current := registry.GetCurrent()
			Expect(current.Screen).To(Equal(ListScreen))
		})

		It("returns error for unregistered screen", func() {
			err := registry.Navigate("nonexistent", nil, nil)
			Expect(err).To(HaveOccurred())
		})

		It("builds breadcrumbs automatically", func() {
			registry.Navigate(HomeScreen, mockModel{id: "home"}, nil)
			registry.Navigate(ListScreen, mockModel{id: "list"}, nil)

			breadcrumbs := registry.GetBreadcrumbs()
			Expect(len(breadcrumbs)).To(BeNumerically(">=", 1))

			// Check last breadcrumb is current screen
			lastCrumb := breadcrumbs[len(breadcrumbs)-1]
			Expect(lastCrumb.Screen).To(Equal(ListScreen))
			Expect(lastCrumb.IsActive).To(BeTrue())
		})

		It("limits history size", func() {
			// Create registry with small history
			smallRegistry := navigation.NewNavigationRegistry(&navigation.NavigationOptions{
				MaxHistorySize: 3,
			})

			smallRegistry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home"},
				{ID: ListScreen, Label: "List"},
				{ID: DetailScreen, Label: "Detail"},
				{ID: EditScreen, Label: "Edit"},
				{ID: SuccessScreen, Label: "Success"},
			})

			// Navigate through 5 screens
			smallRegistry.Navigate(HomeScreen, nil, nil)
			smallRegistry.Navigate(ListScreen, nil, nil)
			smallRegistry.Navigate(DetailScreen, nil, nil)
			smallRegistry.Navigate(EditScreen, nil, nil)
			smallRegistry.Navigate(SuccessScreen, nil, nil)

			// History should be capped at 3
			history := smallRegistry.GetHistory()
			Expect(len(history)).To(BeNumerically("<=", 3))
		})
	})

	Describe("Back Navigation", func() {
		BeforeEach(func() {
			detailParent := ListScreen
			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{
					ID:          HomeScreen,
					Label:       "Home",
					CanGoBack:   false,
					HelpContext: "home",
				},
				{
					ID:          ListScreen,
					Label:       "List",
					CanGoBack:   true,
					Parent:      nil, // Uses stack
					HelpContext: "list",
				},
				{
					ID:          DetailScreen,
					Label:       "Detail",
					CanGoBack:   true,
					Parent:      &detailParent, // Explicit parent to ListScreen
					HelpContext: "detail",
				},
				{
					ID:          NoBackScreen,
					Label:       "No Back",
					CanGoBack:   false,
					HelpContext: "noback",
				},
			})
		})

		It("navigates back using history stack", func() {
			registry.Navigate(HomeScreen, mockModel{id: "home"}, nil)
			registry.Navigate(ListScreen, mockModel{id: "list"}, nil)

			state, err := registry.Back()
			Expect(err).NotTo(HaveOccurred())
			Expect(state.Screen).To(Equal(HomeScreen))

			current := registry.GetCurrent()
			Expect(current.Screen).To(Equal(HomeScreen))
		})

		It("navigates back using explicit parent", func() {
			registry.Navigate(HomeScreen, mockModel{id: "home"}, nil)
			registry.Navigate(ListScreen, mockModel{id: "list"}, nil)
			registry.Navigate(DetailScreen, mockModel{id: "detail"}, nil)

			state, err := registry.Back()
			Expect(err).NotTo(HaveOccurred())
			// Should go to explicit parent (ListScreen), not previous in history
			Expect(state.Screen).To(Equal(ListScreen))
		})

		It("returns error when no history", func() {
			registry.Navigate(HomeScreen, mockModel{id: "home"}, nil)

			_, err := registry.Back()
			Expect(err).To(HaveOccurred())
		})

		It("returns error when screen disallows back", func() {
			registry.Navigate(HomeScreen, mockModel{id: "home"}, nil)
			registry.Navigate(NoBackScreen, mockModel{id: "noback"}, nil)

			_, err := registry.Back()
			Expect(err).To(MatchError(ContainSubstring("cannot navigate back")))
		})
	})

	Describe("CanGoBack", func() {
		BeforeEach(func() {
			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home", CanGoBack: false},
				{ID: ListScreen, Label: "List", CanGoBack: true},
			})
		})

		It("returns true when history exists", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)

			Expect(registry.CanGoBack()).To(BeTrue())
		})

		It("returns false when no history", func() {
			registry.Navigate(HomeScreen, nil, nil)

			Expect(registry.CanGoBack()).To(BeFalse())
		})

		It("returns false when screen disallows back", func() {
			registry.Navigate(ListScreen, nil, nil)
			registry.Navigate(HomeScreen, nil, nil)

			Expect(registry.CanGoBack()).To(BeFalse())
		})
	})

	Describe("GetBackTarget", func() {
		BeforeEach(func() {
			detailParent := ListScreen
			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home", CanGoBack: false},
				{ID: ListScreen, Label: "List", CanGoBack: true},
				{ID: DetailScreen, Label: "Detail", CanGoBack: true, Parent: &detailParent},
			})
		})

		It("returns history target when no explicit parent", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)

			target, err := registry.GetBackTarget()
			Expect(err).NotTo(HaveOccurred())
			Expect(*target).To(Equal(HomeScreen))
		})

		It("returns explicit parent when defined", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(EditScreen, nil, nil)
			registry.Navigate(DetailScreen, nil, nil)

			target, err := registry.GetBackTarget()
			Expect(err).NotTo(HaveOccurred())
			Expect(*target).To(Equal(ListScreen))
		})

		It("returns error when cannot go back", func() {
			registry.Navigate(HomeScreen, nil, nil)

			_, err := registry.GetBackTarget()
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Breadcrumbs", func() {
		BeforeEach(func() {
			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home"},
				{ID: ListScreen, Label: "Events"},
				{ID: DetailScreen, Label: "Details"},
				{ID: ModalScreen, Label: "Modal", IsModal: true},
			})
		})

		It("generates breadcrumb trail", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)
			registry.Navigate(DetailScreen, nil, nil)

			breadcrumbs := registry.GetBreadcrumbs()
			Expect(len(breadcrumbs)).To(BeNumerically(">=", 2))

			// Check structure
			for i, crumb := range breadcrumbs {
				Expect(crumb.Index).To(Equal(i))
				if i == len(breadcrumbs)-1 {
					Expect(crumb.IsActive).To(BeTrue())
				} else {
					Expect(crumb.IsActive).To(BeFalse())
				}
			}
		})

		It("skips modals in breadcrumbs", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)
			registry.Navigate(ModalScreen, nil, nil)

			breadcrumbs := registry.GetBreadcrumbs()

			// Modal should not appear in breadcrumbs
			for _, crumb := range breadcrumbs {
				Expect(crumb.Screen).NotTo(Equal(ModalScreen))
			}
		})

		It("generates breadcrumb path string", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)
			registry.Navigate(DetailScreen, nil, nil)

			path := registry.GetBreadcrumbPath()
			Expect(path).To(ContainSubstring("Home"))
			Expect(path).To(ContainSubstring("Events"))
			Expect(path).To(ContainSubstring("Details"))
			Expect(path).To(ContainSubstring(">"))
		})
	})

	Describe("NavigateToBreadcrumb", func() {
		BeforeEach(func() {
			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home"},
				{ID: ListScreen, Label: "List"},
				{ID: DetailScreen, Label: "Detail"},
			})
		})

		It("navigates to breadcrumb by index", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)
			registry.Navigate(DetailScreen, nil, nil)

			// Navigate to first breadcrumb (should be Home)
			state, err := registry.NavigateToBreadcrumb(0)
			Expect(err).NotTo(HaveOccurred())
			Expect(state.Screen).To(Equal(HomeScreen))
		})

		It("returns error for invalid index", func() {
			registry.Navigate(HomeScreen, nil, nil)

			_, err := registry.NavigateToBreadcrumb(99)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Context Store", func() {
		It("stores and retrieves context values", func() {
			registry.SetContext("key1", "value1")
			registry.SetContext("key2", 42)

			val1, err := registry.GetContext("key1")
			Expect(err).NotTo(HaveOccurred())
			Expect(val1).To(Equal("value1"))

			val2, err := registry.GetContext("key2")
			Expect(err).NotTo(HaveOccurred())
			Expect(val2).To(Equal(42))
		})

		It("returns error for missing key", func() {
			_, err := registry.GetContext("nonexistent")
			Expect(err).To(MatchError(ContainSubstring("not found")))
		})

		It("clears individual context value", func() {
			registry.SetContext("key1", "value1")
			registry.ClearContext("key1")

			_, err := registry.GetContext("key1")
			Expect(err).To(HaveOccurred())
		})

		It("clears all context values", func() {
			registry.SetContext("key1", "value1")
			registry.SetContext("key2", "value2")

			registry.ClearAllContext()

			_, err := registry.GetContext("key1")
			Expect(err).To(HaveOccurred())

			_, err = registry.GetContext("key2")
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Navigation History", func() {
		BeforeEach(func() {
			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home", CanGoBack: false},
				{ID: ListScreen, Label: "List", CanGoBack: true},
				{ID: DetailScreen, Label: "Detail", CanGoBack: true},
			})
		})

		It("returns navigation history", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)
			registry.Navigate(DetailScreen, nil, nil)

			history := registry.GetHistory()
			Expect(len(history)).To(Equal(2))
			Expect(history[0].Screen).To(Equal(HomeScreen))
			Expect(history[1].Screen).To(Equal(ListScreen))
		})

		It("returns history size", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)

			size := registry.GetHistorySize()
			Expect(size).To(Equal(1))
		})

		It("navigates to specific history index", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.Navigate(ListScreen, nil, nil)
			registry.Navigate(DetailScreen, nil, nil)

			state, err := registry.NavigateToHistoryIndex(0)
			Expect(err).NotTo(HaveOccurred())
			Expect(state.Screen).To(Equal(HomeScreen))

			// History should be trimmed
			Expect(registry.GetHistorySize()).To(Equal(0))
		})

		It("returns error for invalid history index", func() {
			registry.Navigate(HomeScreen, nil, nil)

			_, err := registry.NavigateToHistoryIndex(99)
			Expect(err).To(MatchError(ContainSubstring("invalid")))
		})
	})

	Describe("Reset", func() {
		BeforeEach(func() {
			registry.RegisterScreen(&navigation.ScreenDefinition{
				ID:    HomeScreen,
				Label: "Home",
			})
		})

		It("clears all navigation state", func() {
			registry.Navigate(HomeScreen, nil, nil)
			registry.SetContext("key", "value")

			registry.Reset()

			Expect(registry.GetCurrent()).To(BeNil())
			Expect(registry.GetHistorySize()).To(Equal(0))

			_, err := registry.GetContext("key")
			Expect(err).To(HaveOccurred())
		})

		It("preserves screen definitions", func() {
			registry.Reset()

			// Screen should still be registered
			screen, err := registry.GetScreen(HomeScreen)
			Expect(err).NotTo(HaveOccurred())
			Expect(screen.ID).To(Equal(HomeScreen))
		})
	})

	Describe("Hierarchical Breadcrumbs", func() {
		BeforeEach(func() {
			homeScreen := HomeScreen
			listScreen := ListScreen

			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home", Parent: nil},
				{ID: ListScreen, Label: "List", Parent: &homeScreen},
				{ID: DetailScreen, Label: "Detail", Parent: &listScreen},
			})
		})

		It("builds breadcrumbs from parent hierarchy", func() {
			// Use hierarchical builder
			hierarchicalRegistry := navigation.NewNavigationRegistry(&navigation.NavigationOptions{
				BreadcrumbBuilder: navigation.HierarchicalBreadcrumbBuilder,
			})

			homeScreen := HomeScreen
			listScreen := ListScreen

			hierarchicalRegistry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home", Parent: nil},
				{ID: ListScreen, Label: "List", Parent: &homeScreen},
				{ID: DetailScreen, Label: "Detail", Parent: &listScreen},
			})

			hierarchicalRegistry.Navigate(DetailScreen, nil, nil)

			breadcrumbs := hierarchicalRegistry.GetBreadcrumbs()

			// Should show full parent chain: Home > List > Detail
			Expect(len(breadcrumbs)).To(Equal(3))
			Expect(breadcrumbs[0].Screen).To(Equal(HomeScreen))
			Expect(breadcrumbs[1].Screen).To(Equal(ListScreen))
			Expect(breadcrumbs[2].Screen).To(Equal(DetailScreen))
		})
	})

	Describe("Concurrent Access", func() {
		It("handles concurrent navigation safely", func() {
			registry.RegisterScreens([]*navigation.ScreenDefinition{
				{ID: HomeScreen, Label: "Home"},
				{ID: ListScreen, Label: "List"},
				{ID: DetailScreen, Label: "Detail"},
			})

			done := make(chan bool)

			// Multiple goroutines navigating
			for i := 0; i < 10; i++ {
				go func() {
					registry.Navigate(HomeScreen, nil, nil)
					registry.Navigate(ListScreen, nil, nil)
					registry.GetCurrent()
					registry.GetBreadcrumbs()
					done <- true
				}()
			}

			// Wait for all to complete
			for i := 0; i < 10; i++ {
				<-done
			}

			// Should not panic
			Expect(registry.GetCurrent()).NotTo(BeNil())
		})

		It("handles concurrent context access safely", func() {
			done := make(chan bool)

			for i := 0; i < 10; i++ {
				go func(idx int) {
					key := string(rune('a' + idx))
					registry.SetContext(key, idx)
					registry.GetContext(key)
					done <- true
				}(i)
			}

			for i := 0; i < 10; i++ {
				<-done
			}

			// Should complete without race conditions
		})
	})
})
