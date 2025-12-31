package app

import (
	"context"
	"os"
	"strings"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/importer"
	"github.com/baphled/kariya/internal/cli/models"
	"github.com/baphled/kariya/internal/cli/service"
	"github.com/baphled/kariya/internal/cli/styles"
	"github.com/baphled/kariya/internal/domain/career"
	careerservice "github.com/baphled/kariya/internal/service/career"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Screen represents the different screens in the application
type Screen string

const (
	HomeScreen            Screen = "home"
	CaptureScreen         Screen = "capture"
	ListScreen            Screen = "list"
	ViewScreen            Screen = "view"
	QuitScreen            Screen = "quit"
	SuccessScreen         Screen = "success"
	ActionMenuScreen      Screen = "action_menu"
	ConfirmationScreen    Screen = "confirmation"
	ImportReviewScreen    Screen = "import_review"
	ImportProgressScreen  Screen = "import_progress"
	MetadataReviewScreen  Screen = "metadata_review"
	MetadataEditorScreen  Screen = "metadata_editor"
	BulkOperationsScreen  Screen = "bulk_operations"
	BurstSuggestionScreen Screen = "burst_suggestion"
)

// Model represents the main application state
type Model struct {
	cliService             *service.CLIEventService
	service                *careerservice.Service
	currentScreen          Screen
	previousScreen         Screen
	screenBeforeActionMenu Screen   // Track screen before action menu for proper back navigation
	breadcrumbs            []string // Navigation breadcrumb trail
	width                  int
	height                 int
	formModel              *models.FormModel
	successModel           *models.SuccessModel
	listModel              *models.ListModel
	detailsModel           *models.DetailsModel
	actionMenuModel        *models.ActionMenuModel
	confirmationDialog     *models.ConfirmationDialog
	deleteEventID          string // Track the event being deleted
	importService          *importer.ImportService
	importReviewModel      *models.ImportReviewModel
	importProgressModel    *models.ImportProgressModel
	importFilePath         string // Path to CSV file being imported
	metadataReviewModel    *models.MetadataReviewModel
	metadataEditorModel    *models.MetadataEditorModel
	bulkOperationsModel    *models.BulkOperationsModel
	burstSuggestionModel   *models.BurstSuggestionModel
}

// NewModel creates a new application model
func NewModel(cliService *service.CLIEventService, careerService *careerservice.Service) *Model {
	ctx := context.Background()
	return &Model{
		cliService:             cliService,
		service:                careerService,
		currentScreen:          HomeScreen,
		previousScreen:         HomeScreen,
		screenBeforeActionMenu: HomeScreen,
		breadcrumbs:            []string{"Home"},
		width:                  80,
		height:                 24,
		formModel:              models.NewFormModel(cliService),
		successModel:           nil,
		listModel:              models.NewListModel(careerService, ctx),
		detailsModel:           nil,
		actionMenuModel:        nil,
		confirmationDialog:     nil,
		deleteEventID:          "",
		importService:          importer.NewImportService(careerService),
		importReviewModel:      nil,
		importProgressModel:    nil,
		importFilePath:         "",
		metadataReviewModel:    models.NewMetadataReviewModel(careerService, ctx),
		metadataEditorModel:    nil,
		bulkOperationsModel:    nil,
		burstSuggestionModel:   nil,
	}
}

// Init initializes the application
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model state
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle quit and back messages first
	switch msg.(type) {
	case models.BackMsg:
		// Special handling for ListScreen - always go back to HomeScreen
		if m.currentScreen == ListScreen {
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
			return m, nil
		}
		// Special handling for ConfirmationScreen - go back to action menu or previous screen
		if m.currentScreen == ConfirmationScreen {
			m.previousScreen = m.currentScreen
			m.currentScreen = ActionMenuScreen
			m.confirmationDialog = nil
			m.deleteEventID = ""
			return m, nil
		}
		// Special handling for ImportReviewScreen - go back to home
		if m.currentScreen == ImportReviewScreen {
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
			m.importReviewModel = nil
			m.importFilePath = ""
			return m, nil
		}
		// Special handling for ViewScreen that came from ActionMenuScreen
		if m.currentScreen == ViewScreen && m.previousScreen == ActionMenuScreen {
			// Skip ActionMenuScreen and go directly to the screen before it
			targetScreen := m.screenBeforeActionMenu
			m.previousScreen = m.currentScreen
			m.currentScreen = targetScreen
		} else {
			targetScreen := m.previousScreen
			m.previousScreen = m.currentScreen
			m.currentScreen = targetScreen
		}
		return m, nil
	case models.QuitMsg:
		return m, tea.Quit
	}

	// Handle models.ViewEventMsg
	if viewMsg, ok := msg.(models.ViewEventMsg); ok {
		m.detailsModel = models.NewDetailsModel(viewMsg.Event)
		m.previousScreen = m.currentScreen
		m.currentScreen = ViewScreen
		return m, nil
	}

	// Handle EditEventMsg - open metadata editor
	if editMsg, ok := msg.(EditEventMsg); ok {
		ctx := context.Background()
		m.metadataEditorModel = models.NewMetadataEditorModel(editMsg.Event, m.service, m.cliService, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = MetadataEditorScreen
		return m, nil
	}

	// Handle BulkOperationsMsg - open bulk operations
	if bulkMsg, ok := msg.(BulkOperationsMsg); ok {
		ctx := context.Background()
		m.bulkOperationsModel = models.NewBulkOperationsModel(bulkMsg.Events, m.service, m.cliService, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = BulkOperationsScreen
		return m, nil
	}

	// Handle BurstSuggestionsTriggeredMsg - generate and show burst suggestions
	if burstMsg, ok := msg.(models.BurstSuggestionsTriggeredMsg); ok {
		ctx := context.Background()

		// Generate burst suggestions from event IDs
		suggestions, err := m.service.SuggestBursts(ctx, burstMsg.EventIDs)
		if err != nil {
			// Continue without suggestions if error occurs
			return m, nil
		}

		// If no suggestions, return to previous screen
		if len(suggestions) == 0 {
			return m, nil
		}

		// Create burst suggestion model with generated suggestions
		m.burstSuggestionModel = models.NewBurstSuggestionModel(m.service, suggestions, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = BurstSuggestionScreen
		return m, nil
	}

	// Handle ConfirmBurstMsg - persist confirmed burst
	if confirmMsg, ok := msg.(models.ConfirmBurstMsg); ok {
		ctx := context.Background()
		err := m.service.ConfirmBurst(ctx, confirmMsg.Burst)
		if err != nil {
			// Log error but continue processing
			// Could show error to user in future
		}
		return m, nil
	}

	// Handle RejectBurstSuggestionMsg - record rejection
	if rejectMsg, ok := msg.(models.RejectBurstSuggestionMsg); ok {
		ctx := context.Background()
		err := m.service.RejectBurstSuggestion(ctx, rejectMsg.Suggestion.EventIDs)
		if err != nil {
			// Log error but continue processing
		}
		return m, nil
	}

	// Handle BurstProcessingCompleteMsg - navigate back to metadata review or home
	if _, ok := msg.(models.BurstProcessingCompleteMsg); ok {
		// Return to metadata review screen if it exists, otherwise home
		m.previousScreen = m.currentScreen
		if m.metadataReviewModel != nil {
			m.currentScreen = MetadataReviewScreen
		} else {
			m.currentScreen = HomeScreen
		}
		return m, nil
	}

	// Handle models.EventActionMenuMsg
	if actionMenuMsg, ok := msg.(models.EventActionMenuMsg); ok {
		m.actionMenuModel = models.NewActionMenuModel(actionMenuMsg.Event)
		m.previousScreen = m.currentScreen
		m.currentScreen = ActionMenuScreen
		return m, nil
	}

	// Handle EventActionSelectedMsg
	if actionMsg, ok := msg.(models.EventActionSelectedMsg); ok {
		switch actionMsg.Action {
		case models.EventActionView:
			m.detailsModel = models.NewDetailsModel(actionMsg.Event)
			// When transitioning from ActionMenuScreen to ViewScreen,
			// preserve the screen before the action menu for back navigation
			if m.currentScreen == ActionMenuScreen {
				m.screenBeforeActionMenu = m.previousScreen
			}
			m.previousScreen = m.currentScreen
			m.currentScreen = ViewScreen
		case models.EventActionEdit:
			m.previousScreen = m.currentScreen
			m.currentScreen = CaptureScreen
			m.formModel = models.NewFormModel(m.cliService)
			m.formModel.LoadEventForEditing(actionMsg.Event)
		case models.EventActionDelete:
			// Show confirmation dialog for deletion
			m.deleteEventID = actionMsg.Event.ID
			m.confirmationDialog = models.NewConfirmationDialog(
				"Delete Event",
				"Are you sure you want to delete this event? This action cannot be undone.",
			)
			m.previousScreen = m.currentScreen
			m.currentScreen = ConfirmationScreen
		}
		return m, nil
	}

	// Handle confirmation dialog messages
	if m.currentScreen == ConfirmationScreen && m.confirmationDialog != nil {
		updatedDialog, cmd := m.confirmationDialog.Update(msg)
		m.confirmationDialog = updatedDialog

		if m.confirmationDialog.IsConfirmed() {
			// User confirmed deletion - delete the event
			ctx := context.Background()
			err := m.service.DeleteEvent(ctx, m.deleteEventID)
			if err != nil {
				// Handle error - could show error message
				m.previousScreen = m.currentScreen
				m.currentScreen = ListScreen
				m.confirmationDialog = nil
				m.deleteEventID = ""
				return m, nil
			}

			// Refresh the list after deletion
			m.listModel = models.NewListModel(m.service, ctx)
			m.previousScreen = m.currentScreen
			m.currentScreen = ListScreen
			m.confirmationDialog = nil
			m.deleteEventID = ""
			return m, nil
		}

		if m.confirmationDialog.IsCancelled() {
			// User cancelled deletion - go back to action menu
			m.previousScreen = m.currentScreen
			m.currentScreen = ActionMenuScreen
			m.confirmationDialog = nil
			m.deleteEventID = ""
			return m, nil
		}

		return m, cmd
	}

	// Handle ImportReviewMsg
	if importMsg, ok := msg.(models.ImportReviewMsg); ok {
		if importMsg.Action == "cancel" {
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
			m.importReviewModel = nil
			m.importFilePath = ""
			return m, nil
		}

		if importMsg.Action == "import" {
			// Start import process
			m.importProgressModel = models.NewImportProgressModel(len(importMsg.SelectedRows))
			m.previousScreen = m.currentScreen
			m.currentScreen = ImportProgressScreen

			// Return command to perform the import
			return m, func() tea.Msg {
				ctx := context.Background()
				result, err := m.importService.ImportRows(ctx, m.importReviewModel.ParsedRows, importMsg.SelectedRows)
				return models.ImportResultMsg{
					Result: result,
					Error:  err,
				}
			}
		}
	}

	// Handle ImportResultMsg
	if _, ok := msg.(models.ImportResultMsg); ok {
		if m.importProgressModel != nil {
			updatedProgress, cmd := m.importProgressModel.Update(msg)
			m.importProgressModel = updatedProgress.(*models.ImportProgressModel)

			// If import completed, navigate to metadata review
			if m.importProgressModel.Completed {
				switch msg := msg.(type) {
				case models.ImportResultMsg:
					m.previousScreen = m.currentScreen
					m.currentScreen = MetadataReviewScreen
					m.importReviewModel = nil
					m.importProgressModel = nil
					m.importFilePath = ""
					// Initialize metadata review model with imported events
					ctx := context.Background()
					if msg.Result != nil && len(msg.Result.CreatedEvents) > 0 {
						// Extract IDs from created events
						importedIDs := make([]string, 0, len(msg.Result.CreatedEvents))
						for _, event := range msg.Result.CreatedEvents {
							importedIDs = append(importedIDs, event.ID)
						}
						m.metadataReviewModel = models.NewMetadataReviewModelForImport(m.service, ctx, importedIDs)
					} else {
						m.metadataReviewModel = models.NewMetadataReviewModel(m.service, ctx)
					}
					return m, nil
				}
			}

			return m, cmd
		}
	}

	// Handle FormSubmittedMsg
	if submitMsg, ok := msg.(FormSubmittedMsg); ok {
		m.successModel = models.NewSuccessModel(submitMsg.Event)
		m.previousScreen = m.currentScreen
		m.currentScreen = SuccessScreen
		return m, nil
	}

	// Handle SuccessModel messages
	switch msg := msg.(type) {
	case models.CaptureAnotherMsg:
		m.previousScreen = m.currentScreen
		m.currentScreen = CaptureScreen
		m.formModel = models.NewFormModel(m.cliService)
		m.successModel = nil
		return m, nil
	case models.ReviewMetadataMsg:
		ctx := context.Background()
		// Navigate to metadata review screen with just the captured event
		m.metadataReviewModel = models.NewMetadataReviewModelForImport(m.service, ctx, []string{msg.EventID})
		m.previousScreen = m.currentScreen
		m.currentScreen = MetadataReviewScreen
		m.successModel = nil
		return m, nil
	case models.ViewRecentMsg:
		ctx := context.Background()
		m.listModel = models.NewListModel(m.service, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = ListScreen
		m.successModel = nil
		return m, nil
	}

	// Delegate to active model based on current screen
	switch m.currentScreen {
	case ViewScreen:
		if m.detailsModel != nil {
			updatedDetailsModel, cmd := m.detailsModel.Update(msg)
			m.detailsModel = updatedDetailsModel.(*models.DetailsModel)
			return m, cmd
		}

	case CaptureScreen:
		if m.formModel != nil {
			updatedFormModel, cmd := m.formModel.Update(msg)
			m.formModel = updatedFormModel.(*models.FormModel)
			if m.formModel.Submitted() {
				m.successModel = models.NewSuccessModel(m.formModel.Event())
				m.currentScreen = SuccessScreen
			}
			return m, cmd
		}

	case SuccessScreen:
		if m.successModel != nil {
			updatedSuccessModel, cmd := m.successModel.Update(msg)
			m.successModel = updatedSuccessModel.(*models.SuccessModel)
			return m, cmd
		}

	case ListScreen:
		if m.listModel != nil {
			updatedListModel, cmd := m.listModel.Update(msg)
			m.listModel = updatedListModel.(*models.ListModel)
			return m, cmd
		}
	case MetadataReviewScreen:
		if m.metadataReviewModel != nil {
			updatedModel, cmd := m.metadataReviewModel.Update(msg)
			m.metadataReviewModel = updatedModel.(*models.MetadataReviewModel)
			return m, cmd
		}

	case MetadataEditorScreen:
		if m.metadataEditorModel != nil {
			updatedModel, cmd := m.metadataEditorModel.Update(msg)
			m.metadataEditorModel = updatedModel.(*models.MetadataEditorModel)
			if m.metadataEditorModel.IsSubmitted() {
				// Editor submitted - update event in service and return to metadata review
				m.metadataReviewModel.Refresh()
				m.currentScreen = MetadataReviewScreen
				m.previousScreen = MetadataEditorScreen
			} else if m.metadataEditorModel.IsCancelled() {
				// Editor cancelled - return to metadata review
				m.currentScreen = MetadataReviewScreen
				m.previousScreen = MetadataEditorScreen
			}
			return m, cmd
		}

	case BulkOperationsScreen:
		if m.bulkOperationsModel != nil {
			updatedModel, cmd := m.bulkOperationsModel.Update(msg)
			m.bulkOperationsModel = updatedModel.(*models.BulkOperationsModel)
			if m.bulkOperationsModel.ChangesApplied() {
				// Changes applied - return to metadata review
				m.metadataReviewModel.Refresh()
				m.currentScreen = MetadataReviewScreen
				m.previousScreen = BulkOperationsScreen
			} else if m.bulkOperationsModel.WasCancelled() {
				// Cancelled - return to metadata review
				m.currentScreen = MetadataReviewScreen
				m.previousScreen = BulkOperationsScreen
			}
			return m, cmd
		}

	case BurstSuggestionScreen:
		if m.burstSuggestionModel != nil {
			updatedModel, cmd := m.burstSuggestionModel.Update(msg)
			m.burstSuggestionModel = updatedModel.(*models.BurstSuggestionModel)

			if m.burstSuggestionModel.IsDone() {
				// Burst suggestions processing is done
				ctx := context.Background()

				// Persist confirmed bursts
				for _, suggestion := range m.burstSuggestionModel.GetConfirmed() {
					// Convert BurstSuggestion to Burst domain object
					burst := &career.Burst{
						Name:        suggestion.Name,
						Description: suggestion.Description,
						EventIDs:    suggestion.EventIDs,
					}
					if err := m.service.ConfirmBurst(ctx, burst); err != nil {
						// Continue even if one fails
					}
				}

				// Record rejected suggestions to prevent re-suggesting
				for _, suggestion := range m.burstSuggestionModel.GetRejected() {
					m.service.RejectBurstSuggestion(ctx, suggestion.EventIDs)
				}

				// Return to metadata review screen
				m.currentScreen = MetadataReviewScreen
				m.previousScreen = BurstSuggestionScreen
			}
			return m, cmd
		}

	case ActionMenuScreen:
		if m.actionMenuModel != nil {
			updatedActionMenuModel, cmd := m.actionMenuModel.Update(msg)
			m.actionMenuModel = updatedActionMenuModel.(*models.ActionMenuModel)
			return m, cmd
		}

	case ImportReviewScreen:
		if m.importReviewModel != nil {
			updatedImportModel, cmd := m.importReviewModel.Update(msg)
			m.importReviewModel = updatedImportModel.(*models.ImportReviewModel)
			return m, cmd
		}

	case ImportProgressScreen:
		if m.importProgressModel != nil {
			updatedProgressModel, cmd := m.importProgressModel.Update(msg)
			m.importProgressModel = updatedProgressModel.(*models.ImportProgressModel)
			return m, cmd
		}
	}

	// Handle breadcrumb click messages
	if breadcrumbMsg, ok := msg.(BreadcrumbClickedMsg); ok {
		return m.handleBreadcrumbClick(breadcrumbMsg.Index)
	}

	// Handle global navigation shortcuts
	switch msg := msg.(type) {
	case tea.MouseMsg:
		// Handle mouse clicks on breadcrumbs
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			// Check if click was on breadcrumb area (first line of screen)
			clickedIndex := m.getBreadcrumbIndexFromClick(msg.X, msg.Y)
			if clickedIndex >= 0 {
				return m.handleBreadcrumbClick(clickedIndex)
			}
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			m.previousScreen = m.currentScreen
			m.currentScreen = HomeScreen
		case "c":
			m.previousScreen = m.currentScreen
			m.currentScreen = CaptureScreen
			m.formModel = models.NewFormModel(m.cliService)
		case "l":
			ctx := context.Background()
			m.listModel = models.NewListModel(m.service, ctx)
			m.previousScreen = m.currentScreen
			m.currentScreen = ListScreen
		case "m":
			m.metadataReviewModel.Refresh()
			m.previousScreen = m.currentScreen
			m.currentScreen = MetadataReviewScreen
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}

	return m, nil
}

// View renders the current screen
func (m *Model) View() string {
	switch m.currentScreen {
	case HomeScreen:
		return m.renderHome()
	case CaptureScreen:
		if m.formModel != nil {
			m.formModel.SetBreadcrumbs(m.breadcrumbs)
			return m.formModel.View()
		}
		return "Error: Form model not initialized\n"
	case ListScreen:
		if m.listModel != nil {
			m.listModel.SetBreadcrumbs(m.breadcrumbs)
			return m.listModel.View()
		}
		return "Error: List model not initialized\n"
	case MetadataReviewScreen:
		if m.metadataReviewModel != nil {
			m.metadataReviewModel.SetBreadcrumbs(m.breadcrumbs)
			// Ensure width/height are propagated
			m.metadataReviewModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			return m.metadataReviewModel.View()
		}
		return "Error: Metadata Review model not initialized\n"
	case MetadataEditorScreen:
		if m.metadataEditorModel != nil {
			return m.metadataEditorModel.View()
		}
		return "Error: Metadata Editor model not initialized\n"
	case BulkOperationsScreen:
		if m.bulkOperationsModel != nil {
			return m.bulkOperationsModel.View()
		}
		return "Error: Bulk Operations model not initialized\n"
	case BurstSuggestionScreen:
		if m.burstSuggestionModel != nil {
			return m.burstSuggestionModel.View()
		}
		return "Error: Burst Suggestion model not initialized\n"
	case ViewScreen:
		if m.detailsModel != nil {
			return m.detailsModel.View()
		}
		return "Error: Details model not initialized\n"
	case SuccessScreen:
		if m.successModel != nil {
			m.successModel.SetBreadcrumbs(m.breadcrumbs)
			// Ensure width/height are propagated
			m.successModel.Update(tea.WindowSizeMsg{Width: m.width, Height: m.height})
			return m.successModel.View()
		}
		return "Success!\n"
	case ActionMenuScreen:
		if m.actionMenuModel != nil {
			return m.actionMenuModel.View()
		}
		return "Error: Action menu not initialized\n"
	case ConfirmationScreen:
		if m.confirmationDialog != nil {
			return m.confirmationDialog.View()
		}
		return "Error: Confirmation dialog not initialized\n"
	case ImportReviewScreen:
		if m.importReviewModel != nil {
			return m.importReviewModel.View()
		}
		return "Error: Import review model not initialized\n"
	case ImportProgressScreen:
		if m.importProgressModel != nil {
			return m.importProgressModel.View()
		}
		return "Error: Import progress model not initialized\n"
	case QuitScreen:
		return "Goodbye!\n"
	default:
		return m.renderHome()
	}
}

// renderHome renders the home screen
func (m *Model) renderHome() string {
	title := styles.HeaderMain.Render("KaRiya - Career Journal CLI")
	commandsHeader := styles.HeaderSection.Render("Commands:")
	commands := []string{
		styles.InfoText.Render("c") + " - Capture Career Event",
		styles.InfoText.Render("l") + " - List Events",
		styles.InfoText.Render("h") + " - Home",
		styles.InfoText.Render("q") + " - Quit",
	}
	commandsList := strings.Join(commands, "\n")
	cta := styles.SuccessBox.Width(styles.MaxWidth(m.width) - 4).Render("Press 'c' to get started!")
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		commandsHeader,
		commandsList,
		"",
		cta,
	)
	card := styles.ResponsiveCard(m.width).Render(content)
	return styles.Center(card, m.width, m.height)
}

// SetInitialScreen sets the initial screen to display on startup
func (m *Model) SetInitialScreen(screen Screen) {
	m.currentScreen = screen
	m.previousScreen = screen
}

// SetInitialCaptureMode sets the initial capture mode for the form
func (m *Model) SetInitialCaptureMode(mode string) {
	if m.formModel != nil {
		m.formModel.SetInitialMode(mode)
	}
}

// StartImport initiates the import process with a CSV file
func (m *Model) StartImport(filePath string) tea.Cmd {
	m.importFilePath = filePath
	m.previousScreen = m.currentScreen
	m.currentScreen = ImportReviewScreen

	return func() tea.Msg {
		// Open and parse the CSV file
		file, err := os.Open(filePath)
		if err != nil {
			return models.ImportReviewMsg{Action: "error"}
		}
		defer file.Close()

		ctx := context.Background()
		parsedRows, err := m.importService.PrepareImport(ctx, file)
		if err != nil {
			return models.ImportReviewMsg{Action: "error"}
		}

		// Create import review model
		m.importReviewModel = models.NewImportReviewModel(parsedRows)
		return models.ImportReviewMsg{Action: "prepared"}
	}
}

// GetBreadcrumbs returns current breadcrumb trail
func (m *Model) GetBreadcrumbs() []string {
	return m.breadcrumbs
}

// updateBreadcrumbs updates breadcrumb trail based on current screen and navigation history
func (m *Model) updateBreadcrumbs() {
	switch m.currentScreen {
	case HomeScreen:
		m.breadcrumbs = []string{"Home"}
	case CaptureScreen:
		m.breadcrumbs = []string{"Home", "Capture Event"}
	case ListScreen:
		m.breadcrumbs = []string{"Home", "Events"}
	case ViewScreen:
		// Build breadcrumb based on previous screen
		if m.previousScreen == ListScreen {
			m.breadcrumbs = []string{"Home", "Events", "Details"}
		} else if m.previousScreen == ActionMenuScreen {
			m.breadcrumbs = []string{"Home", "Events", "Details"}
		} else {
			m.breadcrumbs = []string{"Home", "Details"}
		}
	case MetadataReviewScreen:
		m.breadcrumbs = []string{"Home", "Metadata Review"}
	case MetadataEditorScreen:
		m.breadcrumbs = []string{"Home", "Metadata Review", "Edit Event"}
	case BulkOperationsScreen:
		m.breadcrumbs = []string{"Home", "Metadata Review", "Bulk Operations"}
	case BurstSuggestionScreen:
		m.breadcrumbs = []string{"Home", "Metadata Review", "Burst Suggestions"}
	case SuccessScreen:
		// Build breadcrumb based on previous screen
		if m.previousScreen == CaptureScreen {
			m.breadcrumbs = []string{"Home", "Capture Event", "Success"}
		} else {
			m.breadcrumbs = []string{"Home", "Success"}
		}
	case ActionMenuScreen:
		m.breadcrumbs = []string{"Home", "Events", "Actions"}
	case ConfirmationScreen:
		m.breadcrumbs = []string{"Home", "Events", "Confirm Delete"}
	case ImportReviewScreen:
		m.breadcrumbs = []string{"Home", "Import Review"}
	case ImportProgressScreen:
		m.breadcrumbs = []string{"Home", "Import", "Progress"}
	default:
		m.breadcrumbs = []string{"Home"}
	}
}

// getBreadcrumbIndexFromClick determines which breadcrumb was clicked based on mouse coordinates
func (m *Model) getBreadcrumbIndexFromClick(x, y int) int {
	// Breadcrumbs are rendered at the top of the screen
	// We need to account for any offset and use the header's click detection

	// Create a temporary header with current breadcrumbs
	header := components.NewHeader("", m.width)
	header.SetBreadcrumbs(m.breadcrumbs)

	return header.GetClickedBreadcrumbIndex(x, y)
}

// handleBreadcrumbClick navigates based on which breadcrumb was clicked
func (m *Model) handleBreadcrumbClick(index int) (tea.Model, tea.Cmd) {
	if index < 0 || index >= len(m.breadcrumbs) {
		return m, nil
	}

	// Map breadcrumb index to screen navigation
	// Index 0 is always "Home", index 1 is the parent screen, etc.
	if index == 0 {
		// Navigate to Home
		m.previousScreen = m.currentScreen
		m.currentScreen = HomeScreen
		m.updateBreadcrumbs()
		return m, nil
	}

	// For other breadcrumbs, we need to navigate based on the breadcrumb trail
	// This is more complex as we need to reverse-engineer the screen from the breadcrumb
	breadcrumb := m.breadcrumbs[index]

	switch breadcrumb {
	case "Capture Event":
		m.previousScreen = m.currentScreen
		m.currentScreen = CaptureScreen
		m.formModel = models.NewFormModel(m.cliService)
	case "Events":
		ctx := context.Background()
		m.listModel = models.NewListModel(m.service, ctx)
		m.previousScreen = m.currentScreen
		m.currentScreen = ListScreen
	case "Metadata Review":
		m.metadataReviewModel.Refresh()
		m.previousScreen = m.currentScreen
		m.currentScreen = MetadataReviewScreen
	case "Import Review":
		m.previousScreen = m.currentScreen
		m.currentScreen = ImportReviewScreen
	}

	m.updateBreadcrumbs()
	return m, nil
}
