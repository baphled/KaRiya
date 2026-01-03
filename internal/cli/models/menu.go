package models

import (
	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/navigation"
	"github.com/baphled/kariya/internal/cli/styles"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MenuItem represents a menu option with its metadata
type MenuItem struct {
	Key         navigation.NavigationKey // Keyboard shortcut using standardized navigation key
	Title       string                   // Display name
	Description string                   // Detailed description
	Category    string                   // Menu category (e.g., "Core Actions", "Data Management")
}

// MenuModel represents the main menu screen
type MenuModel struct {
	*BaseStandardModel
	items       []MenuItem
	selectedIdx int
	width       int
	height      int
	header      components.HeaderModel
	helpFooter  components.HelpFooterModel
	hasPending  bool   // Whether there are pending items to review
	pendingInfo string // Information about pending items
	keyHandler  navigation.KeyHandler
}

// NewMenuModel creates a new menu model with all available menu items
func NewMenuModel(hasPending bool, pendingInfo string) *MenuModel {
	m := &MenuModel{
		BaseStandardModel: NewBaseStandardModel(),
		selectedIdx:       0,
		width:             80,
		height:            24,
		header:            components.NewHeader("KaRiya - Career Journal CLI", 80),
		helpFooter:        components.NewHelpFooter("menu", 80),
		hasPending:        hasPending,
		pendingInfo:       pendingInfo,
		keyHandler:        navigation.NewMenuKeyHandler(),
	}

	// Initialize menu items organized by category
	// Using standardized NavigationKey constants for consistency
	m.items = []MenuItem{
		// Core Actions
		{
			Key:         navigation.KeyCapture,
			Title:       "Capture Career Event",
			Description: "Record a new career event, achievement, or milestone",
			Category:    "Core Actions",
		},
		{
			Key:         navigation.KeyList,
			Title:       "List Events",
			Description: "View all recorded career events with filtering and sorting",
			Category:    "Core Actions",
		},
		{
			Key:         navigation.KeyBulk,
			Title:       "View Bursts",
			Description: "Analyze career bursts and activity patterns",
			Category:    "Core Actions",
		},
		// Data Management
		{
			Key:         navigation.KeyMetadata,
			Title:       "Metadata Review",
			Description: "Review and enrich event metadata, extract facts and bursts",
			Category:    "Data Management",
		},
		{
			Key:         navigation.KeyFacts,
			Title:       "View All Facts",
			Description: "Browse and manage all extracted facts from your career events",
			Category:    "Data Management",
		},
		// CV Management
		{
			Key:         navigation.KeyCV,
			Title:       "Manage CV Configurations",
			Description: "Create, edit, and manage CV configuration profiles",
			Category:    "CV Management",
		},
		{
			Key:         navigation.KeyGenerate,
			Title:       "Generate CV",
			Description: "Generate a CV from a saved configuration profile",
			Category:    "CV Management",
		},
		// Navigation
		{
			Key:         navigation.KeyHome,
			Title:       "Home",
			Description: "Return to the home screen",
			Category:    "Navigation",
		},
		{
			Key:         navigation.KeyHelp,
			Title:       "Help",
			Description: "Show keyboard shortcuts and help information",
			Category:    "Navigation",
		},
		{
			Key:         navigation.KeyQuit,
			Title:       "Quit",
			Description: "Exit the application",
			Category:    "Navigation",
		},
	}

	// Add pending items option if applicable
	if hasPending {
		m.items = append([]MenuItem{
			{
				Key:         navigation.KeyPending,
				Title:       "Review Pending Items",
				Description: pendingInfo,
				Category:    "Pending",
			},
		}, m.items...)
	}

	return m
}

// Init initializes the model
func (m *MenuModel) Init() tea.Cmd {
	return nil
}

// Update handles messages for the menu
func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.header.SetWidth(msg.Width)
		m.helpFooter.SetWidth(msg.Width)
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	case tea.KeyMsg:
		// Use centralized key handler for menu navigation
		action := m.keyHandler.HandleKey(msg)
		if !action.IsHandled {
			return m, nil
		}

		// Handle the action based on navigation key
		switch {
		case action.NavigationKey != nil:
			return m.handleNavigationKey(*action.NavigationKey)
		}
	}
	return m, nil
}

// handleNavigationKey processes navigation key actions for menu
func (m *MenuModel) handleNavigationKey(key navigation.NavigationKey) (tea.Model, tea.Cmd) {
	switch key {
	case navigation.KeyUp:
		if m.selectedIdx > 0 {
			m.selectedIdx--
		}
		return m, nil
	case navigation.KeyDown:
		if m.selectedIdx < len(m.items)-1 {
			m.selectedIdx++
		}
		return m, nil
	case navigation.KeySelect:
		// Return the selected menu item's key
		selectedItem := m.items[m.selectedIdx]
		return m, func() tea.Msg {
			return MenuItemSelectedMsg{Key: string(selectedItem.Key)}
		}
	case navigation.KeyBack:
		// Escape from menu - go back or quit
		return m, func() tea.Msg { return BackMsg{} }
	}
	return m, nil
}

// View renders the menu
func (m *MenuModel) View() string {
	headerContent := m.renderHeader()
	contentArea := m.renderMenuContent()
	footerContent := m.renderFooter()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		headerContent,
		contentArea,
		footerContent,
	)
}

// renderMenuContent renders the menu items with categories
func (m *MenuModel) renderMenuContent() string {
	var sections []string

	// Group items by category
	categories := make(map[string][]MenuItem)
	categoryOrder := []string{}

	for _, item := range m.items {
		if _, exists := categories[item.Category]; !exists {
			categoryOrder = append(categoryOrder, item.Category)
		}
		categories[item.Category] = append(categories[item.Category], item)
	}

	// Render each category
	for _, category := range categoryOrder {
		items := categories[category]

		// Category header
		categoryHeader := styles.HeaderSection.Render(category)
		sections = append(sections, categoryHeader)

		// Render items in this category
		for _, item := range items {
			// Find the index in the overall items list
			globalIdx := 0
			for j, it := range m.items {
				if it.Key == item.Key && it.Title == item.Title {
					globalIdx = j
					break
				}
			}

			itemView := m.renderMenuItem(item, globalIdx == m.selectedIdx)
			sections = append(sections, itemView)
		}

		// Add spacing between categories
		sections = append(sections, "")
	}

	// Pending items notification if applicable
	if m.hasPending {
		notification := styles.WarningBox.Width(styles.MaxWidth(m.width) - 4).Render(
			"⚠️  " + m.pendingInfo,
		)
		sections = append(sections, "", notification)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(content)
}

// renderMenuItem renders a single menu item
func (m *MenuModel) renderMenuItem(item MenuItem, isSelected bool) string {
	// Format the key shortcut - convert NavigationKey to string for display
	keyDisplay := styles.InfoText.Render(string(item.Key))

	// Format the title
	titleStyle := styles.InfoText
	if isSelected {
		titleStyle = styles.ListItemSelected
	}
	titleDisplay := titleStyle.Render(item.Title)

	// Format the description
	descriptionDisplay := styles.InfoHint.Render("  " + item.Description)

	// Build the menu item with selection indicator
	var indicator string
	if isSelected {
		indicator = styles.ListItemFocused.Render("▶ ")
	} else {
		indicator = "  "
	}

	// Combine key and title on first line
	firstLine := indicator + keyDisplay + " - " + titleDisplay

	// Return both lines
	return lipgloss.JoinVertical(
		lipgloss.Left,
		firstLine,
		descriptionDisplay,
	)
}

// renderHeader renders the header component
func (m *MenuModel) renderHeader() string {
	return m.header.View()
}

// renderFooter renders the footer component
func (m *MenuModel) renderFooter() string {
	return m.helpFooter.View()
}

// MenuItemSelectedMsg is sent when a menu item is selected
type MenuItemSelectedMsg struct {
	Key string // Keyboard shortcut that was selected (as string for compatibility)
}
