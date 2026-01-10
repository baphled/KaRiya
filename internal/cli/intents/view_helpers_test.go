package intents_test

import (
	"errors"
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
)

// CreateStandardView Tests

func TestCreateStandardView(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	view := intents.CreateStandardView(base)

	if view == nil {
		t.Fatal("Expected CreateStandardView to return non-nil")
	}

	if view.TerminalInfo != info {
		t.Error("Expected view to have correct terminal info")
	}

	if !view.UseFullWidth {
		t.Error("Expected view to use full width")
	}
}

func TestCreateStandardView_WithLogo(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Set logo
	logo := components.NewASCIILogo(false, 100)
	base.SetLogo(logo)
	base.SetLogoSpacing(3)

	view := intents.CreateStandardView(base)

	if view.Logo != logo {
		t.Error("Expected view to have logo from BaseIntent")
	}

	if view.LogoSpacing != 3 {
		t.Errorf("Expected logo spacing 3, got %d", view.LogoSpacing)
	}
}

func TestCreateStandardViewWithBreadcrumbs(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	view := intents.CreateStandardViewWithBreadcrumbs(base, "Home", "Settings", "Display")

	if view == nil {
		t.Fatal("Expected CreateStandardViewWithBreadcrumbs to return non-nil")
	}

	if len(view.Breadcrumbs) != 3 {
		t.Errorf("Expected 3 breadcrumbs, got %d", len(view.Breadcrumbs))
	}

	expected := []string{"Home", "Settings", "Display"}
	for i, crumb := range view.Breadcrumbs {
		if crumb != expected[i] {
			t.Errorf("Expected breadcrumb %d to be '%s', got '%s'", i, expected[i], crumb)
		}
	}
}

// Modal Priority Tests

func TestCreateStandardView_ModalPriority_Error(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Set error (highest priority)
	base.SetError(errors.New("test error"))

	view := intents.CreateStandardView(base)

	if !view.ShowModal {
		t.Error("Expected view to show modal when error is set")
	}

	if view.Modal == nil {
		t.Fatal("Expected view to have modal")
	}

	if view.Modal.Type != components.ModalError {
		t.Errorf("Expected error modal, got type %v", view.Modal.Type)
	}
}

func TestCreateStandardView_ModalPriority_Loading(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Set loading (second priority)
	base.SetLoading("Loading...")

	view := intents.CreateStandardView(base)

	if !view.ShowModal {
		t.Error("Expected view to show modal when loading is set")
	}

	if view.Modal == nil {
		t.Fatal("Expected view to have modal")
	}

	if view.Modal.Type != components.ModalLoading {
		t.Errorf("Expected loading modal, got type %v", view.Modal.Type)
	}
}

func TestCreateStandardView_ModalPriority_Progress(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Set progress (third priority)
	base.SetProgress("Processing", "Step 1", 0.5)

	view := intents.CreateStandardView(base)

	if !view.ShowModal {
		t.Error("Expected view to show modal when progress is set")
	}

	if view.Modal == nil {
		t.Fatal("Expected view to have modal")
	}

	if view.Modal.Type != components.ModalProgress {
		t.Errorf("Expected progress modal, got type %v", view.Modal.Type)
	}
}

func TestCreateStandardView_ModalPriority_Success(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Set success (lowest priority)
	base.SetSuccess("Success!")

	view := intents.CreateStandardView(base)

	if !view.ShowModal {
		t.Error("Expected view to show modal when success is set")
	}

	if view.Modal == nil {
		t.Fatal("Expected view to have modal")
	}

	if view.Modal.Type != components.ModalSuccess {
		t.Errorf("Expected success modal, got type %v", view.Modal.Type)
	}
}

func TestCreateStandardView_ModalPriority_ErrorOverLoading(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Set both error and loading
	base.SetLoading("Loading...")
	base.SetError(errors.New("error occurred"))

	view := intents.CreateStandardView(base)

	// Error should take priority
	if view.Modal.Type != components.ModalError {
		t.Errorf("Expected error modal to take priority, got type %v", view.Modal.Type)
	}
}

func TestCreateStandardView_ModalPriority_LoadingOverSuccess(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Set both loading and success
	base.SetLoading("Loading...")
	base.SetSuccess("Success!")

	view := intents.CreateStandardView(base)

	// Loading should take priority
	if view.Modal.Type != components.ModalLoading {
		t.Errorf("Expected loading modal to take priority, got type %v", view.Modal.Type)
	}
}

// Error Title Extraction Tests

func TestExtractErrorTitle_Validation(t *testing.T) {
	base := intents.NewBaseIntent()
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	base.SetError(errors.New("validation failed: field is required"))
	view := intents.CreateStandardView(base)

	if !strings.Contains(view.Modal.Title, "Validation") {
		t.Errorf("Expected validation error title, got '%s'", view.Modal.Title)
	}
}

func TestExtractErrorTitle_Database(t *testing.T) {
	base := intents.NewBaseIntent()
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	base.SetError(errors.New("database connection failed"))
	view := intents.CreateStandardView(base)

	if !strings.Contains(view.Modal.Title, "Database") {
		t.Errorf("Expected database error title, got '%s'", view.Modal.Title)
	}
}

func TestExtractErrorTitle_Default(t *testing.T) {
	base := intents.NewBaseIntent()
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	base.SetError(errors.New("something went wrong"))
	view := intents.CreateStandardView(base)

	if view.Modal.Title != "Error" {
		t.Errorf("Expected default error title 'Error', got '%s'", view.Modal.Title)
	}
}

// Manual Modal Helper Tests

func TestShowErrorModal(t *testing.T) {
	view := components.NewStandardView(terminal.NewInfo())
	err := errors.New("test error")

	result := intents.ShowErrorModal(view, err)

	if result != view {
		t.Error("Expected ShowErrorModal to return the same view")
	}

	if !view.ShowModal {
		t.Error("Expected modal to be shown")
	}

	if view.Modal.Type != components.ModalError {
		t.Errorf("Expected error modal, got type %v", view.Modal.Type)
	}
}

func TestShowErrorModal_NilError(t *testing.T) {
	view := components.NewStandardView(terminal.NewInfo())

	result := intents.ShowErrorModal(view, nil)

	if result != view {
		t.Error("Expected ShowErrorModal to return the same view")
	}

	// Should not show modal for nil error
	if view.ShowModal {
		t.Error("Expected modal not to be shown for nil error")
	}
}

func TestShowLoadingModal(t *testing.T) {
	view := components.NewStandardView(terminal.NewInfo())

	result := intents.ShowLoadingModal(view, "Processing...", true)

	if result != view {
		t.Error("Expected ShowLoadingModal to return the same view")
	}

	if !view.ShowModal {
		t.Error("Expected modal to be shown")
	}

	if view.Modal.Type != components.ModalLoading {
		t.Errorf("Expected loading modal, got type %v", view.Modal.Type)
	}

	if view.Modal.Message != "Processing..." {
		t.Errorf("Expected message 'Processing...', got '%s'", view.Modal.Message)
	}
}

func TestShowProgressModal(t *testing.T) {
	view := components.NewStandardView(terminal.NewInfo())

	result := intents.ShowProgressModal(view, "Uploading", "50% complete", 0.5)

	if result != view {
		t.Error("Expected ShowProgressModal to return the same view")
	}

	if !view.ShowModal {
		t.Error("Expected modal to be shown")
	}

	if view.Modal.Type != components.ModalProgress {
		t.Errorf("Expected progress modal, got type %v", view.Modal.Type)
	}

	if view.Modal.Progress != 0.5 {
		t.Errorf("Expected progress 0.5, got %f", view.Modal.Progress)
	}
}

func TestShowSuccessModal(t *testing.T) {
	view := components.NewStandardView(terminal.NewInfo())

	result := intents.ShowSuccessModal(view, "Operation completed!")

	if result != view {
		t.Error("Expected ShowSuccessModal to return the same view")
	}

	if !view.ShowModal {
		t.Error("Expected modal to be shown")
	}

	if view.Modal.Type != components.ModalSuccess {
		t.Errorf("Expected success modal, got type %v", view.Modal.Type)
	}

	if view.Modal.Message != "Operation completed!" {
		t.Errorf("Expected message 'Operation completed!', got '%s'", view.Modal.Message)
	}
}

// Footer Helper Tests

func TestStandardHelpFooter(t *testing.T) {
	shortcuts := map[string]string{
		"Enter": "Select",
		"Esc":   "Back",
	}

	footer := intents.StandardHelpFooter(shortcuts)

	if footer == "" {
		t.Error("Expected non-empty footer")
	}

	// Should contain both shortcuts
	if !strings.Contains(footer, "Enter") || !strings.Contains(footer, "Select") {
		t.Errorf("Expected footer to contain shortcuts, got '%s'", footer)
	}
}

func TestStandardHelpFooter_Empty(t *testing.T) {
	footer := intents.StandardHelpFooter(map[string]string{})

	if footer != "" {
		t.Errorf("Expected empty footer, got '%s'", footer)
	}
}

func TestNavigationFooter(t *testing.T) {
	footer := intents.NavigationFooter()

	if footer == "" {
		t.Error("Expected non-empty navigation footer")
	}

	// Should contain standard navigation shortcuts
	expectedParts := []string{"Up", "Down", "Enter", "Select", "Esc", "Back"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected navigation footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestFormFooter(t *testing.T) {
	footer := intents.FormFooter()

	if footer == "" {
		t.Error("Expected non-empty form footer")
	}

	// Should contain form navigation shortcuts
	expectedParts := []string{"Tab", "Next", "Enter", "Submit", "Esc", "Cancel"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected form footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestListFooter(t *testing.T) {
	footer := intents.ListFooter()

	if footer == "" {
		t.Error("Expected non-empty list footer")
	}

	// Should contain list shortcuts including search
	expectedParts := []string{"Up", "Down", "Search", "Top", "Bottom"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected list footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestDetailViewFooter(t *testing.T) {
	footer := intents.DetailViewFooter()

	if footer == "" {
		t.Error("Expected non-empty detail view footer")
	}

	// Should contain scroll shortcuts
	expectedParts := []string{"Scroll", "Up", "Down", "Esc", "Back"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected detail view footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestModalFooter(t *testing.T) {
	actions := []string{"Enter Confirm", "Esc Cancel"}
	footer := intents.ModalFooter(actions)

	if footer == "" {
		t.Error("Expected non-empty modal footer")
	}

	for _, action := range actions {
		if !strings.Contains(footer, action) {
			t.Errorf("Expected modal footer to contain '%s', got '%s'", action, footer)
		}
	}
}

func TestModalFooter_Empty(t *testing.T) {
	footer := intents.ModalFooter([]string{})

	if !strings.Contains(footer, "Close") {
		t.Errorf("Expected default 'Esc Close', got '%s'", footer)
	}
}

func TestCombineFooters(t *testing.T) {
	footer1 := "↑/k Up  ↓/j Down"
	footer2 := "Enter Select"
	footer3 := "Esc Back"

	combined := intents.CombineFooters(footer1, footer2, footer3)

	if combined == "" {
		t.Error("Expected non-empty combined footer")
	}

	// Should contain all parts separated by |
	if !strings.Contains(combined, footer1) {
		t.Errorf("Expected combined footer to contain '%s'", footer1)
	}
	if !strings.Contains(combined, footer2) {
		t.Errorf("Expected combined footer to contain '%s'", footer2)
	}
	if !strings.Contains(combined, footer3) {
		t.Errorf("Expected combined footer to contain '%s'", footer3)
	}

	// Should use | separator
	if !strings.Contains(combined, "|") {
		t.Error("Expected combined footer to use | separator")
	}
}

func TestCombineFooters_Empty(t *testing.T) {
	combined := intents.CombineFooters()

	if combined != "" {
		t.Errorf("Expected empty combined footer, got '%s'", combined)
	}
}

func TestCombineFooters_WithEmptyStrings(t *testing.T) {
	combined := intents.CombineFooters("Footer 1", "", "Footer 2", "   ")

	// Should skip empty strings
	if strings.Contains(combined, "||") {
		t.Error("Expected combined footer to skip empty strings")
	}

	// Should still contain non-empty parts
	if !strings.Contains(combined, "Footer 1") || !strings.Contains(combined, "Footer 2") {
		t.Errorf("Expected combined footer to contain non-empty parts, got '%s'", combined)
	}
}

// =============================================================================
// Theme-Aware KeyBadge Footer Helper Tests
// =============================================================================

func TestThemedNavigationFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedNavigationFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed navigation footer")
	}

	// Should contain navigation elements
	expectedParts := []string{"Navigate", "Select", "Back"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed navigation footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedNavigationFooter_NilTheme(t *testing.T) {
	footer := intents.ThemedNavigationFooter(nil)

	if footer == "" {
		t.Error("Expected non-empty footer even with nil theme")
	}

	// Should still contain navigation elements (fallback styling)
	if !strings.Contains(footer, "Navigate") {
		t.Errorf("Expected footer to contain 'Navigate', got '%s'", footer)
	}
}

func TestThemedFormFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedFormFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed form footer")
	}

	// Should contain form elements
	expectedParts := []string{"Next", "Previous", "Submit", "Cancel"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed form footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedListFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedListFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed list footer")
	}

	// Should contain list elements including search
	expectedParts := []string{"Navigate", "Select", "Search", "Back"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed list footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedDetailViewFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedDetailViewFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed detail view footer")
	}

	// Should contain scroll and back elements
	expectedParts := []string{"Scroll", "Back"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed detail view footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedConfirmFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedConfirmFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed confirm footer")
	}

	// Should contain confirm and cancel
	expectedParts := []string{"Confirm", "Cancel"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed confirm footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedEditFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedEditFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed edit footer")
	}

	// Should contain save and cancel
	expectedParts := []string{"Save", "Cancel"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed edit footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedExportFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedExportFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed export footer")
	}

	// Should contain navigate, select, confirm, back
	expectedParts := []string{"Navigate", "Select", "Confirm", "Back"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed export footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedMenuFooter(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedMenuFooter(theme)

	if footer == "" {
		t.Error("Expected non-empty themed menu footer")
	}

	// Should contain menu elements
	expectedParts := []string{"Navigate", "Select", "Help", "Quit"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed menu footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedCustomFooter(t *testing.T) {
	theme := createTestTheme()

	footer := intents.ThemedCustomFooter(theme,
		components.NewKeyBadge("x", "Custom1"),
		components.NewKeyBadge("y", "Custom2"),
	)

	if footer == "" {
		t.Error("Expected non-empty themed custom footer")
	}

	// Should contain custom badges
	expectedParts := []string{"Custom1", "Custom2", "x", "y"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed custom footer to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestThemedCustomFooter_Empty(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedCustomFooter(theme)

	if footer != "" {
		t.Errorf("Expected empty footer with no badges, got '%s'", footer)
	}
}

func TestThemedGlobalBadges(t *testing.T) {
	theme := createTestTheme()
	footer := intents.ThemedGlobalBadges(theme)

	if footer == "" {
		t.Error("Expected non-empty themed global badges")
	}

	// Should contain quit and main menu
	expectedParts := []string{"Quit", "Main Menu"}
	for _, part := range expectedParts {
		if !strings.Contains(footer, part) {
			t.Errorf("Expected themed global badges to contain '%s', got '%s'", part, footer)
		}
	}
}

func TestCombineThemedFooters(t *testing.T) {
	theme := createTestTheme()

	footer1 := intents.ThemedNavigationFooter(theme)
	footer2 := intents.ThemedGlobalBadges(theme)

	combined := intents.CombineThemedFooters(footer1, footer2)

	if combined == "" {
		t.Error("Expected non-empty combined themed footer")
	}

	// Should contain elements from both footers
	expectedParts := []string{"Navigate", "Quit"}
	for _, part := range expectedParts {
		if !strings.Contains(combined, part) {
			t.Errorf("Expected combined themed footer to contain '%s', got '%s'", part, combined)
		}
	}
}

func TestCombineThemedFooters_Empty(t *testing.T) {
	combined := intents.CombineThemedFooters()

	if combined != "" {
		t.Errorf("Expected empty combined themed footer, got '%s'", combined)
	}
}

func TestCombineThemedFooters_WithEmptyStrings(t *testing.T) {
	theme := createTestTheme()

	footer1 := intents.ThemedNavigationFooter(theme)

	combined := intents.CombineThemedFooters(footer1, "", "   ")

	// Should skip empty strings
	if combined == "" {
		t.Error("Expected non-empty combined footer")
	}

	// Should contain elements from non-empty footer
	if !strings.Contains(combined, "Navigate") {
		t.Errorf("Expected combined footer to contain 'Navigate', got '%s'", combined)
	}
}

// Helper function to create a test theme
func createTestTheme() themes.Theme {
	return themes.NewDefaultTheme()
}
