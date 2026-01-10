package intents_test

import (
	"errors"
	"testing"
	"time"

	"github.com/baphled/kariya/internal/cli/components"
	"github.com/baphled/kariya/internal/cli/intents"
	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
)

func TestNewBaseIntent(t *testing.T) {
	base := intents.NewBaseIntent()

	if base == nil {
		t.Fatal("Expected NewBaseIntent to return non-nil")
	}

	if base.GetTerminalInfo() == nil {
		t.Error("Expected terminal info to be initialized")
	}

	if base.GetLogoSpacing() != 2 {
		t.Errorf("Expected default logo spacing to be 2, got %d", base.GetLogoSpacing())
	}
}

func TestBaseIntent_TerminalInfo(t *testing.T) {
	base := intents.NewBaseIntent()

	// Test UpdateTerminalInfo
	info := terminal.NewInfo()
	info.Width = 120
	info.Height = 40
	info.IsValid = true

	base.UpdateTerminalInfo(info)

	retrieved := base.GetTerminalInfo()
	if retrieved.Width != 120 {
		t.Errorf("Expected width 120, got %d", retrieved.Width)
	}
	if retrieved.Height != 40 {
		t.Errorf("Expected height 40, got %d", retrieved.Height)
	}
	if !retrieved.IsValid {
		t.Error("Expected terminal info to be valid")
	}
}

func TestBaseIntent_GetMinimumSize(t *testing.T) {
	base := intents.NewBaseIntent()

	width, height := base.GetMinimumSize()
	if width <= 0 || height <= 0 {
		t.Errorf("Expected positive minimum size, got width=%d, height=%d", width, height)
	}
}

// Logo Management Tests

func TestBaseIntent_LogoManagement(t *testing.T) {
	base := intents.NewBaseIntent()

	// Initially no logo
	if base.GetLogo() != nil {
		t.Error("Expected GetLogo to return nil initially")
	}

	// Set logo
	logo := components.NewASCIILogo(false, 80)
	base.SetLogo(logo)

	if base.GetLogo() != logo {
		t.Error("Expected GetLogo to return the set logo")
	}
}

func TestBaseIntent_LogoSpacing(t *testing.T) {
	base := intents.NewBaseIntent()

	// Default spacing
	if base.GetLogoSpacing() != 2 {
		t.Errorf("Expected default spacing 2, got %d", base.GetLogoSpacing())
	}

	// Set spacing
	base.SetLogoSpacing(5)
	if base.GetLogoSpacing() != 5 {
		t.Errorf("Expected spacing 5, got %d", base.GetLogoSpacing())
	}
}

// Loading State Tests

func TestBaseIntent_LoadingState(t *testing.T) {
	base := intents.NewBaseIntent()

	// Initially not loading
	if base.IsLoading() {
		t.Error("Expected IsLoading to be false initially")
	}

	// Set loading
	base.SetLoading("Processing...")
	if !base.IsLoading() {
		t.Error("Expected IsLoading to be true after SetLoading")
	}
	if base.GetLoadingMessage() != "Processing..." {
		t.Errorf("Expected loading message 'Processing...', got '%s'", base.GetLoadingMessage())
	}

	// Clear loading
	base.ClearLoading()
	if base.IsLoading() {
		t.Error("Expected IsLoading to be false after ClearLoading")
	}
	if base.GetLoadingMessage() != "" {
		t.Errorf("Expected empty loading message after clear, got '%s'", base.GetLoadingMessage())
	}
}

// Error State Tests

func TestBaseIntent_ErrorState(t *testing.T) {
	base := intents.NewBaseIntent()

	// Initially no error
	if base.HasError() {
		t.Error("Expected HasError to be false initially")
	}
	if base.GetError() != nil {
		t.Error("Expected GetError to return nil initially")
	}

	// Set error
	testErr := errors.New("test error")
	base.SetError(testErr)

	if !base.HasError() {
		t.Error("Expected HasError to be true after SetError")
	}
	if base.GetError() != testErr {
		t.Errorf("Expected GetError to return test error, got %v", base.GetError())
	}

	// Clear error
	base.ClearError()
	if base.HasError() {
		t.Error("Expected HasError to be false after ClearError")
	}
	if base.GetError() != nil {
		t.Error("Expected GetError to return nil after ClearError")
	}
}

// Success State Tests

func TestBaseIntent_SuccessState(t *testing.T) {
	base := intents.NewBaseIntent()

	// Initially no success
	if base.ShouldShowSuccess() {
		t.Error("Expected ShouldShowSuccess to be false initially")
	}

	// Set success
	base.SetSuccess("Operation completed!")
	if !base.ShouldShowSuccess() {
		t.Error("Expected ShouldShowSuccess to be true after SetSuccess")
	}
	if base.GetSuccessMessage() != "Operation completed!" {
		t.Errorf("Expected success message 'Operation completed!', got '%s'", base.GetSuccessMessage())
	}

	// Clear success
	base.ClearSuccess()
	if base.ShouldShowSuccess() {
		t.Error("Expected ShouldShowSuccess to be false after ClearSuccess")
	}
	if base.GetSuccessMessage() != "" {
		t.Errorf("Expected empty success message after clear, got '%s'", base.GetSuccessMessage())
	}
}

func TestBaseIntent_SuccessExpiry(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set success
	base.SetSuccess("Test message")

	// Should show immediately
	if !base.ShouldShowSuccess() {
		t.Error("Expected ShouldShowSuccess to be true immediately after SetSuccess")
	}

	// Wait 3.5 seconds
	time.Sleep(3500 * time.Millisecond)

	// Should not show after 3 seconds
	if base.ShouldShowSuccess() {
		t.Error("Expected ShouldShowSuccess to be false after 3 seconds")
	}
}

// Progress State Tests

func TestBaseIntent_ProgressState(t *testing.T) {
	base := intents.NewBaseIntent()

	// Initially no progress
	if base.IsProgressEnabled() {
		t.Error("Expected IsProgressEnabled to be false initially")
	}

	// Set progress
	base.SetProgress("Processing", "Step 1 of 3", 0.33)

	if !base.IsProgressEnabled() {
		t.Error("Expected IsProgressEnabled to be true after SetProgress")
	}

	title, message, value := base.GetProgress()
	if title != "Processing" {
		t.Errorf("Expected progress title 'Processing', got '%s'", title)
	}
	if message != "Step 1 of 3" {
		t.Errorf("Expected progress message 'Step 1 of 3', got '%s'", message)
	}
	if value != 0.33 {
		t.Errorf("Expected progress value 0.33, got %f", value)
	}

	// Clear progress
	base.ClearProgress()
	if base.IsProgressEnabled() {
		t.Error("Expected IsProgressEnabled to be false after ClearProgress")
	}

	title, message, value = base.GetProgress()
	if title != "" || message != "" || value != 0.0 {
		t.Errorf("Expected cleared progress state, got title='%s', message='%s', value=%f", title, message, value)
	}
}

// View Creation Tests

func TestBaseIntent_CreateView(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Create view
	view := base.CreateView()

	if view == nil {
		t.Fatal("Expected CreateView to return non-nil")
	}

	if view.TerminalInfo != info {
		t.Error("Expected view to have terminal info from BaseIntent")
	}
}

func TestBaseIntent_CreateViewWithBreadcrumbs(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Create view with breadcrumbs
	view := base.CreateViewWithBreadcrumbs("Main", "Settings", "Display")

	if view == nil {
		t.Fatal("Expected CreateViewWithBreadcrumbs to return non-nil")
	}

	if len(view.Breadcrumbs) != 3 {
		t.Errorf("Expected 3 breadcrumbs, got %d", len(view.Breadcrumbs))
	}

	if view.Breadcrumbs[0] != "Main" || view.Breadcrumbs[1] != "Settings" || view.Breadcrumbs[2] != "Display" {
		t.Errorf("Expected breadcrumbs ['Main', 'Settings', 'Display'], got %v", view.Breadcrumbs)
	}
}

func TestBaseIntent_CreateView_WithLogo(t *testing.T) {
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

	// Create view
	view := base.CreateView()

	if view == nil {
		t.Fatal("Expected CreateView to return non-nil")
	}

	if view.Logo != logo {
		t.Error("Expected view to have logo from BaseIntent")
	}

	if view.LogoSpacing != 2 {
		t.Errorf("Expected logo spacing 2, got %d", view.LogoSpacing)
	}
}

func TestBaseIntent_CreateView_WithStates(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set up terminal info
	info := terminal.NewInfo()
	info.Width = 100
	info.Height = 30
	info.IsValid = true
	base.UpdateTerminalInfo(info)

	// Test error state triggers modal
	base.SetError(errors.New("test error"))
	view := base.CreateView()

	if view == nil {
		t.Fatal("Expected CreateView to return non-nil")
	}

	if !view.ShowModal {
		t.Error("Expected view to show modal when error state is set")
	}

	if view.Modal == nil {
		t.Error("Expected view to have modal content when error state is set")
	}
}

// State Independence Tests

func TestBaseIntent_StateIndependence(t *testing.T) {
	base := intents.NewBaseIntent()

	// Set multiple states independently
	base.SetLoading("Loading...")
	base.SetError(errors.New("error occurred"))
	base.SetSuccess("Success!")
	base.SetProgress("Processing", "Step 1", 0.5)

	// All states should be independent
	if !base.IsLoading() {
		t.Error("Expected loading state to remain set")
	}
	if !base.HasError() {
		t.Error("Expected error state to remain set")
	}
	if !base.ShouldShowSuccess() {
		t.Error("Expected success state to remain set")
	}
	if !base.IsProgressEnabled() {
		t.Error("Expected progress state to remain set")
	}

	// Clear one state should not affect others
	base.ClearLoading()
	if !base.HasError() || !base.ShouldShowSuccess() || !base.IsProgressEnabled() {
		t.Error("Expected other states to remain unaffected when clearing loading")
	}
}

// Theme Management Tests

func TestBaseIntent_ThemeManagement(t *testing.T) {
	base := intents.NewBaseIntent()

	// Initially no theme manager
	if base.GetThemeManager() != nil {
		t.Error("Expected GetThemeManager to return nil initially")
	}

	// Theme should return nil when no manager is set
	if base.Theme() != nil {
		t.Error("Expected Theme to return nil when no manager is set")
	}

	// Set theme manager
	tm := themes.NewThemeManager()
	base.SetThemeManager(tm)

	if base.GetThemeManager() != tm {
		t.Error("Expected GetThemeManager to return the set theme manager")
	}

	// Theme should return the active theme
	theme := base.Theme()
	if theme == nil {
		t.Error("Expected Theme to return non-nil when manager is set")
	}

	if theme.Name() != "default" {
		t.Errorf("Expected theme name 'default', got '%s'", theme.Name())
	}
}

func TestBaseIntent_ThemeStyles(t *testing.T) {
	base := intents.NewBaseIntent()
	tm := themes.NewThemeManager()
	base.SetThemeManager(tm)

	theme := base.Theme()
	if theme == nil {
		t.Fatal("Expected theme to be non-nil")
	}

	// Verify we can access palette
	palette := theme.Palette()
	if palette == nil {
		t.Error("Expected palette to be non-nil")
	}

	// Verify we can access styles
	styles := theme.Styles()
	if styles == nil {
		t.Error("Expected styles to be non-nil")
	}
}
