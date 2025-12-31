package navigation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/baphled/kariya/internal/cli/navigation"
)

// Integration tests for NavigationRegistry simulating real application scenarios

var _ = Describe("NavigationRegistry Integration Tests", func() {
	var registry *navigation.NavigationRegistry

	// Define a realistic screen hierarchy similar to KaRiya app
	const (
		HomeScreen            navigation.Screen = "home"
		CaptureScreen         navigation.Screen = "capture"
		ListScreen            navigation.Screen = "list"
		ViewScreen            navigation.Screen = "view"
		EditScreen            navigation.Screen = "edit"
		ActionMenuScreen      navigation.Screen = "action_menu"
		ConfirmationScreen    navigation.Screen = "confirmation"
		SuccessScreen         navigation.Screen = "success"
		MetadataReviewScreen  navigation.Screen = "metadata_review"
		MetadataEditorScreen  navigation.Screen = "metadata_editor"
		BulkOperationsScreen  navigation.Screen = "bulk_operations"
		ImportReviewScreen    navigation.Screen = "import_review"
		BurstSuggestionScreen navigation.Screen = "burst_suggestion"
	)

	BeforeEach(func() {
		registry = navigation.NewNavigationRegistry(nil)

		// Register screens with realistic parent relationships
		homeParent := HomeScreen
		listParent := ListScreen
		metadataReviewParent := MetadataReviewScreen

		err := registry.RegisterScreens([]*navigation.ScreenDefinition{
			{
				ID:          HomeScreen,
				Label:       "Home",
				CanGoBack:   false,
				HelpContext: "home",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyCapture, navigation.KeyList, navigation.KeyMetadata},
			},
			{
				ID:          CaptureScreen,
				Label:       "Capture Event",
				CanGoBack:   true,
				Parent:      &homeParent,
				HelpContext: "form",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect},
			},
			{
				ID:          ListScreen,
				Label:       "Events",
				CanGoBack:   true,
				Parent:      &homeParent,
				HelpContext: "list",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect, navigation.KeyFilter, navigation.KeySort},
			},
			{
				ID:          ViewScreen,
				Label:       "Event Details",
				CanGoBack:   true,
				Parent:      &listParent,
				HelpContext: "default",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeyEdit, navigation.KeyDelete},
			},
			{
				ID:          EditScreen,
				Label:       "Edit Event",
				CanGoBack:   true,
				Parent:      &listParent,
				HelpContext: "form",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect},
			},
			{
				ID:          ActionMenuScreen,
				Label:       "Actions",
				CanGoBack:   true,
				IsModal:     true,
				HelpContext: "default",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect},
			},
			{
				ID:          ConfirmationScreen,
				Label:       "Confirm",
				CanGoBack:   true,
				IsModal:     true,
				HelpContext: "default",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect},
			},
			{
				ID:          SuccessScreen,
				Label:       "Success",
				CanGoBack:   false,
				HelpContext: "default",
				Shortcuts:   []navigation.NavigationKey{navigation.KeySelect, navigation.KeyHome},
			},
			{
				ID:          MetadataReviewScreen,
				Label:       "Metadata Review",
				CanGoBack:   true,
				Parent:      &homeParent,
				HelpContext: "metadata_review",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect, navigation.KeyEdit, navigation.KeyBulk},
			},
			{
				ID:          MetadataEditorScreen,
				Label:       "Edit Metadata",
				CanGoBack:   true,
				Parent:      &metadataReviewParent,
				HelpContext: "form",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect},
			},
			{
				ID:          BulkOperationsScreen,
				Label:       "Bulk Operations",
				CanGoBack:   true,
				Parent:      &metadataReviewParent,
				HelpContext: "bulk_operations",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect, navigation.KeyToggle},
			},
			{
				ID:          ImportReviewScreen,
				Label:       "Import Review",
				CanGoBack:   true,
				Parent:      &homeParent,
				HelpContext: "default",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect, navigation.KeyToggle},
			},
			{
				ID:          BurstSuggestionScreen,
				Label:       "Burst Suggestions",
				CanGoBack:   true,
				Parent:      &metadataReviewParent,
				HelpContext: "default",
				Shortcuts:   []navigation.NavigationKey{navigation.KeyBack, navigation.KeySelect},
			},
		})
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Real-world navigation flows", func() {
		Context("Capture workflow", func() {
			It("navigates through capture → success → home", func() {
				// Start at home
				err := registry.Navigate(HomeScreen, mockModel{id: "home"}, nil)
				Expect(err).NotTo(HaveOccurred())

				// Navigate to capture
				err = registry.Navigate(CaptureScreen, mockModel{id: "capture"}, nil)
				Expect(err).NotTo(HaveOccurred())

				breadcrumbs := registry.GetBreadcrumbs()
				Expect(len(breadcrumbs)).To(Equal(2))
				Expect(breadcrumbs[0].Label).To(Equal("Home"))
				Expect(breadcrumbs[1].Label).To(Equal("Capture Event"))

				// After successful capture, go to success screen
				err = registry.Navigate(SuccessScreen, mockModel{id: "success"}, map[string]interface{}{
					"eventID": "evt-123",
				})
				Expect(err).NotTo(HaveOccurred())

				// Context preserved
				current := registry.GetCurrent()
				Expect(current.Context["eventID"]).To(Equal("evt-123"))

				// From success, navigate home (success screen can't go back)
				Expect(registry.CanGoBack()).To(BeFalse())
			})

			It("supports back navigation from capture to home", func() {
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(CaptureScreen, nil, nil)

				// Back to home (using explicit parent)
				state, err := registry.Back()
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(HomeScreen))
			})
		})

		Context("List → View → Edit workflow", func() {
			It("navigates and maintains proper breadcrumbs", func() {
				// Home → List → View
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)

				// Store event context
				eventCtx := map[string]interface{}{
					"eventID": "evt-456",
					"title":   "Career Milestone",
				}
				registry.Navigate(ViewScreen, nil, eventCtx)

				breadcrumbs := registry.GetBreadcrumbs()
				Expect(len(breadcrumbs)).To(Equal(3))
				Expect(breadcrumbs[0].Label).To(Equal("Home"))
				Expect(breadcrumbs[1].Label).To(Equal("Events"))
				Expect(breadcrumbs[2].Label).To(Equal("Event Details"))

				// Context preserved
				current := registry.GetCurrent()
				Expect(current.Context["eventID"]).To(Equal("evt-456"))

				// Back goes to List (explicit parent)
				state, err := registry.Back()
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(ListScreen))
			})

			It("handles breadcrumb click navigation", func() {
				// Build navigation trail
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.Navigate(ViewScreen, nil, nil)

				_ = registry.GetBreadcrumbs()

				// Click on first breadcrumb (Home)
				state, err := registry.NavigateToBreadcrumb(0)
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(HomeScreen))

				// History should be cleared
				Expect(registry.GetHistorySize()).To(Equal(0))
			})
		})

		Context("Metadata review workflow", func() {
			It("navigates through metadata review → editor → back", func() {
				// Home → Metadata Review
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(MetadataReviewScreen, nil, map[string]interface{}{
					"importedEventIDs": []string{"evt-1", "evt-2", "evt-3"},
				})

				// Store event to edit in context store
				registry.SetContext("currentEditEventID", "evt-2")

				// Navigate to editor
				registry.Navigate(MetadataEditorScreen, nil, map[string]interface{}{
					"eventID": "evt-2",
					"field":   "description",
				})

				breadcrumbs := registry.GetBreadcrumbs()
				Expect(len(breadcrumbs)).To(Equal(3))
				Expect(breadcrumbs[2].Label).To(Equal("Edit Metadata"))

				// Back to metadata review
				state, err := registry.Back()
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(MetadataReviewScreen))

				// Context still available
				eventID, err := registry.GetContext("currentEditEventID")
				Expect(err).NotTo(HaveOccurred())
				Expect(eventID).To(Equal("evt-2"))
			})

			It("handles bulk operations flow", func() {
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(MetadataReviewScreen, nil, nil)

				// Navigate to bulk operations
				selectedEvents := []string{"evt-1", "evt-2", "evt-3"}
				registry.Navigate(BulkOperationsScreen, nil, map[string]interface{}{
					"selectedEvents": selectedEvents,
				})

				current := registry.GetCurrent()
				Expect(current.Context["selectedEvents"]).To(Equal(selectedEvents))

				// Back to metadata review
				state, err := registry.Back()
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(MetadataReviewScreen))
			})
		})

		Context("Modal screens", func() {
			It("excludes modals from breadcrumbs", func() {
				// Navigate with action menu modal
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.Navigate(ViewScreen, nil, nil)
				registry.Navigate(ActionMenuScreen, nil, nil)

				breadcrumbs := registry.GetBreadcrumbs()

				// Modal should not appear
				for _, crumb := range breadcrumbs {
					Expect(crumb.Screen).NotTo(Equal(ActionMenuScreen))
				}

				// But we should be able to go back
				Expect(registry.CanGoBack()).To(BeTrue())
			})

			It("handles confirmation dialogs", func() {
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.Navigate(ViewScreen, nil, map[string]interface{}{
					"eventID": "evt-to-delete",
				})

				// Open confirmation
				registry.Navigate(ConfirmationScreen, nil, map[string]interface{}{
					"action":  "delete",
					"eventID": "evt-to-delete",
				})

				// Confirmation is modal, not in breadcrumbs
				breadcrumbs := registry.GetBreadcrumbs()
				for _, crumb := range breadcrumbs {
					Expect(crumb.Screen).NotTo(Equal(ConfirmationScreen))
				}

				// Can navigate back
				state, err := registry.Back()
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(ViewScreen))
			})
		})

		Context("Complex navigation patterns", func() {
			It("handles deep navigation with mixed parent strategies", func() {
				// Home → List → View → Edit
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.Navigate(ViewScreen, nil, nil)
				registry.Navigate(EditScreen, nil, nil)

				// Edit has explicit parent (List), not View
				target, err := registry.GetBackTarget()
				Expect(err).NotTo(HaveOccurred())
				Expect(*target).To(Equal(ListScreen))

				// Navigate back
				state, err := registry.Back()
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(ListScreen))
			})

			It("preserves context across multiple navigations", func() {
				// Start import workflow
				registry.SetContext("importFile", "/path/to/data.csv")

				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ImportReviewScreen, nil, map[string]interface{}{
					"rows":       []string{"row1", "row2"},
					"validRows":  2,
					"errorCount": 0,
				})

				// After import, go to metadata review
				registry.Navigate(MetadataReviewScreen, nil, map[string]interface{}{
					"importedEventIDs": []string{"evt-1", "evt-2"},
				})

				// Global context still available
				importFile, err := registry.GetContext("importFile")
				Expect(err).NotTo(HaveOccurred())
				Expect(importFile).To(Equal("/path/to/data.csv"))

				// Local context preserved
				current := registry.GetCurrent()
				Expect(current.Context["importedEventIDs"]).To(HaveLen(2))
			})

			It("maintains history size limits", func() {
				// Create registry with small history
				smallRegistry := navigation.NewNavigationRegistry(&navigation.NavigationOptions{
					MaxHistorySize: 5,
				})

				// Register just a few screens
				smallRegistry.RegisterScreens([]*navigation.ScreenDefinition{
					{ID: HomeScreen, Label: "Home"},
					{ID: CaptureScreen, Label: "Capture"},
					{ID: ListScreen, Label: "List"},
					{ID: ViewScreen, Label: "View"},
					{ID: EditScreen, Label: "Edit"},
					{ID: SuccessScreen, Label: "Success"},
					{ID: MetadataReviewScreen, Label: "Metadata"},
				})

				// Navigate through many screens
				smallRegistry.Navigate(HomeScreen, nil, nil)
				smallRegistry.Navigate(CaptureScreen, nil, nil)
				smallRegistry.Navigate(ListScreen, nil, nil)
				smallRegistry.Navigate(ViewScreen, nil, nil)
				smallRegistry.Navigate(EditScreen, nil, nil)
				smallRegistry.Navigate(SuccessScreen, nil, nil)
				smallRegistry.Navigate(MetadataReviewScreen, nil, nil)

				// History should be capped
				history := smallRegistry.GetHistory()
				Expect(len(history)).To(BeNumerically("<=", 5))
			})
		})

		Context("Navigation state recovery", func() {
			It("recovers from history at any point", func() {
				// Build history
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.Navigate(ViewScreen, nil, nil)
				registry.Navigate(EditScreen, nil, nil)

				// Get history
				history := registry.GetHistory()
				Expect(len(history)).To(Equal(3))

				// Jump to second item in history
				state, err := registry.NavigateToHistoryIndex(1)
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(ListScreen))

				// History trimmed
				newHistory := registry.GetHistory()
				Expect(len(newHistory)).To(Equal(1))
				Expect(newHistory[0].Screen).To(Equal(HomeScreen))
			})

			It("resets cleanly", func() {
				// Build complex state
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.SetContext("key1", "value1")
				registry.SetContext("key2", "value2")

				// Reset
				registry.Reset()

				// Everything cleared
				Expect(registry.GetCurrent()).To(BeNil())
				Expect(registry.GetHistorySize()).To(Equal(0))
				_, err := registry.GetContext("key1")
				Expect(err).To(HaveOccurred())
			})
		})

		Context("Hierarchical breadcrumbs", func() {
			It("builds breadcrumbs from parent relationships", func() {
				// Use hierarchical builder
				hierarchicalRegistry := navigation.NewNavigationRegistry(&navigation.NavigationOptions{
					BreadcrumbBuilder: navigation.HierarchicalBreadcrumbBuilder,
				})

				homeParent := HomeScreen
				metadataReviewParent := MetadataReviewScreen

				hierarchicalRegistry.RegisterScreens([]*navigation.ScreenDefinition{
					{ID: HomeScreen, Label: "Home", Parent: nil},
					{ID: MetadataReviewScreen, Label: "Metadata Review", Parent: &homeParent},
					{ID: MetadataEditorScreen, Label: "Edit Metadata", Parent: &metadataReviewParent},
				})

				// Navigate directly to deep screen
				hierarchicalRegistry.Navigate(MetadataEditorScreen, nil, nil)

				// Breadcrumbs should show full parent chain
				breadcrumbs := hierarchicalRegistry.GetBreadcrumbs()
				Expect(len(breadcrumbs)).To(Equal(3))
				Expect(breadcrumbs[0].Screen).To(Equal(HomeScreen))
				Expect(breadcrumbs[1].Screen).To(Equal(MetadataReviewScreen))
				Expect(breadcrumbs[2].Screen).To(Equal(MetadataEditorScreen))
			})
		})

		Context("Concurrent navigation", func() {
			It("handles concurrent access safely", func() {
				done := make(chan bool)
				errors := make(chan error, 10)

				// Multiple goroutines performing navigation
				for i := 0; i < 10; i++ {
					go func(idx int) {
						defer func() { done <- true }()

						// Navigate
						if err := registry.Navigate(HomeScreen, nil, nil); err != nil {
							errors <- err
							return
						}

						if err := registry.Navigate(ListScreen, nil, nil); err != nil {
							errors <- err
							return
						}

						// Set context
						registry.SetContext(string(rune('a'+idx)), idx)

						// Get breadcrumbs
						_ = registry.GetBreadcrumbs()

						// Try back
						if registry.CanGoBack() {
							_, _ = registry.Back()
						}
					}(i)
				}

				// Wait for all to complete
				for i := 0; i < 10; i++ {
					<-done
				}

				close(errors)

				// Check for errors
				errList := make([]error, 0)
				for err := range errors {
					errList = append(errList, err)
				}

				Expect(errList).To(BeEmpty())
			})
		})
	})

	Describe("Real-world error scenarios", func() {
		It("handles navigation to unregistered screen gracefully", func() {
			registry.Navigate(HomeScreen, nil, nil)

			err := registry.Navigate("nonexistent", nil, nil)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))

			// Current state unchanged
			current := registry.GetCurrent()
			Expect(current.Screen).To(Equal(HomeScreen))
		})

		It("handles back navigation when not allowed", func() {
			registry.Navigate(SuccessScreen, nil, nil)

			_, err := registry.Back()
			Expect(err).To(HaveOccurred())

			Expect(registry.CanGoBack()).To(BeFalse())
		})

		It("handles breadcrumb navigation to invalid index", func() {
			registry.Navigate(HomeScreen, nil, nil)

			_, err := registry.NavigateToBreadcrumb(99)
			Expect(err).To(HaveOccurred())
		})

		It("handles context access for missing keys", func() {
			_, err := registry.GetContext("nonexistent")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("not found"))
		})
	})

	Describe("Navigation patterns from app.go", func() {
		Context("Screen transition messages", func() {
			It("simulates FormSubmittedMsg → Success → Home", func() {
				// Capture form submission
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(CaptureScreen, nil, nil)

				// Store form data in context
				registry.SetContext("submittedEvent", map[string]interface{}{
					"title":       "New Event",
					"date":        "2025-12-31",
					"description": "Test event",
				})

				// Navigate to success
				registry.Navigate(SuccessScreen, nil, map[string]interface{}{
					"eventID": "evt-789",
					"message": "Event captured successfully",
				})

				// Verify context
				eventData, err := registry.GetContext("submittedEvent")
				Expect(err).NotTo(HaveOccurred())
				Expect(eventData).NotTo(BeNil())
			})

			It("simulates ActionMenu → ViewScreen flow", func() {
				// View event
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.Navigate(ViewScreen, nil, map[string]interface{}{
					"eventID": "evt-123",
				})

				// Open action menu
				registry.Navigate(ActionMenuScreen, nil, map[string]interface{}{
					"actions": []string{"Edit", "Delete", "Export"},
				})

				// Back skips modal
				state, err := registry.Back()
				Expect(err).NotTo(HaveOccurred())
				Expect(state.Screen).To(Equal(ViewScreen))
			})
		})

		Context("Global navigation shortcuts", func() {
			It("simulates 'h' key (home) from any screen", func() {
				// Navigate deep
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(ListScreen, nil, nil)
				registry.Navigate(ViewScreen, nil, nil)
				registry.Navigate(EditScreen, nil, nil)

				// Jump to home (direct navigation)
				err := registry.Navigate(HomeScreen, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				// History preserved
				Expect(registry.GetHistorySize()).To(Equal(4))

				current := registry.GetCurrent()
				Expect(current.Screen).To(Equal(HomeScreen))
			})

			It("simulates 'l' key (list) from any screen", func() {
				registry.Navigate(HomeScreen, nil, nil)
				registry.Navigate(CaptureScreen, nil, nil)

				// Jump to list
				err := registry.Navigate(ListScreen, nil, nil)
				Expect(err).NotTo(HaveOccurred())

				current := registry.GetCurrent()
				Expect(current.Screen).To(Equal(ListScreen))
			})
		})
	})
})
