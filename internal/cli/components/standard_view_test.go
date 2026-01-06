package components

import (
	"strings"
	"testing"

	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/charmbracelet/lipgloss"
)

func TestNewStandardView(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo)

	if view == nil {
		t.Fatal("Expected NewStandardView to return non-nil view")
	}

	if view.ShowLogo {
		t.Error("Expected ShowLogo to be false by default")
	}

	if view.LogoSpacing != 2 {
		t.Errorf("Expected LogoSpacing to be 2, got %d", view.LogoSpacing)
	}

	if !view.ShowFooter {
		t.Error("Expected ShowFooter to be true by default")
	}

	if view.ShowFooterSeparator {
		t.Error("Expected ShowFooterSeparator to be false by default")
	}

	if !view.UseFullWidth {
		t.Error("Expected UseFullWidth to be true by default")
	}
}

func TestStandardView_WithLogo(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)
	view := NewStandardView(termInfo).WithLogo(logo, 3)

	if !view.ShowLogo {
		t.Error("Expected ShowLogo to be true after WithLogo")
	}

	if view.Logo != logo {
		t.Error("Expected Logo to be set")
	}

	if view.LogoSpacing != 3 {
		t.Errorf("Expected LogoSpacing to be 3, got %d", view.LogoSpacing)
	}
}

func TestStandardView_WithBreadcrumbs(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo).WithBreadcrumbs("Home", "Settings", "Profile")

	if !view.ShowHeader {
		t.Error("Expected ShowHeader to be true after WithBreadcrumbs")
	}

	if len(view.Breadcrumbs) != 3 {
		t.Errorf("Expected 3 breadcrumbs, got %d", len(view.Breadcrumbs))
	}

	if view.Breadcrumbs[0] != "Home" {
		t.Errorf("Expected first breadcrumb to be 'Home', got '%s'", view.Breadcrumbs[0])
	}
}

func TestStandardView_WithTitle(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo).WithTitle("My Title", "My Subtitle")

	if !view.ShowHeader {
		t.Error("Expected ShowHeader to be true after WithTitle")
	}

	if view.Title != "My Title" {
		t.Errorf("Expected Title to be 'My Title', got '%s'", view.Title)
	}

	if view.Subtitle != "My Subtitle" {
		t.Errorf("Expected Subtitle to be 'My Subtitle', got '%s'", view.Subtitle)
	}
}

func TestStandardView_WithContent(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	content := "Test content\nLine 2\nLine 3"
	view := NewStandardView(termInfo).WithContent(content)

	if view.Content != content {
		t.Errorf("Expected Content to be '%s', got '%s'", content, view.Content)
	}
}

func TestStandardView_WithContentStyle(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	style := lipgloss.NewStyle().Bold(true)
	view := NewStandardView(termInfo).WithContentStyle(style)

	if !view.ContentStyle.GetBold() {
		t.Error("Expected ContentStyle to have Bold set to true")
	}
}

func TestStandardView_WithHelp(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	helpText := "Press q to quit"
	view := NewStandardView(termInfo).WithHelp(helpText)

	if !view.ShowFooter {
		t.Error("Expected ShowFooter to be true after WithHelp")
	}

	if view.HelpText != helpText {
		t.Errorf("Expected HelpText to be '%s', got '%s'", helpText, view.HelpText)
	}
}

func TestStandardView_WithFooterSeparator(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo).WithFooterSeparator(true)

	if !view.ShowFooterSeparator {
		t.Error("Expected ShowFooterSeparator to be true")
	}

	view2 := NewStandardView(termInfo).WithFooterSeparator(false)
	if view2.ShowFooterSeparator {
		t.Error("Expected ShowFooterSeparator to be false")
	}
}

func TestStandardView_ShowModalOverlay(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	modal := NewErrorModal("Error", "Test error")
	view := NewStandardView(termInfo).ShowModalOverlay(modal)

	if !view.ShowModal {
		t.Error("Expected ShowModal to be true after ShowModalOverlay")
	}

	if view.Modal != modal {
		t.Error("Expected Modal to be set")
	}
}

func TestStandardView_SetUseFullWidth(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo).SetUseFullWidth(false)

	if view.UseFullWidth {
		t.Error("Expected UseFullWidth to be false")
	}

	view2 := NewStandardView(termInfo).SetUseFullWidth(true)
	if !view2.UseFullWidth {
		t.Error("Expected UseFullWidth to be true")
	}
}

func TestStandardView_Render_WithLogo(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)
	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithContent("Test content")

	output := view.Render()

	if output == "" {
		t.Error("Expected Render to return non-empty string")
	}

	// Should contain logo art
	if !strings.Contains(output, "KARIYA") && !strings.Contains(output, "██") {
		t.Error("Expected output to contain logo art")
	}
}

func TestStandardView_Render_WithBreadcrumbs(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo).
		WithBreadcrumbs("Home", "Settings").
		WithContent("Test content")

	output := view.Render()

	if !strings.Contains(output, "Home > Settings") {
		t.Error("Expected output to contain breadcrumbs with '>' separator")
	}
}

func TestStandardView_Render_WithFooterSeparator(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo).
		WithContent("Test content").
		WithHelp("Press q to quit").
		WithFooterSeparator(true)

	output := view.Render()

	// Should contain separator character
	if !strings.Contains(output, "─") {
		t.Error("Expected output to contain footer separator")
	}
}

func TestStandardView_Render_WithHelp(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	helpText := "Press q to quit, h for help"
	view := NewStandardView(termInfo).
		WithContent("Test content").
		WithHelp(helpText)

	output := view.Render()

	if !strings.Contains(output, helpText) {
		t.Error("Expected output to contain help text")
	}
}

func TestStandardView_BuilderChaining(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	logo := NewASCIILogo(false, termInfo.Width)

	// Test that all builder methods return *StandardView for chaining
	view := NewStandardView(termInfo).
		WithLogo(logo, 2).
		WithBreadcrumbs("Home", "Settings").
		WithTitle("Title", "Subtitle").
		WithContent("Content").
		WithHelp("Help").
		WithFooterSeparator(true).
		SetUseFullWidth(true)

	if view == nil {
		t.Error("Expected chained builder to return non-nil view")
	}

	// Verify all settings were applied
	if !view.ShowLogo {
		t.Error("Expected ShowLogo to be true")
	}
	if len(view.Breadcrumbs) == 0 {
		t.Error("Expected breadcrumbs to be set")
	}
	if view.Title == "" {
		t.Error("Expected title to be set")
	}
	if view.Content == "" {
		t.Error("Expected content to be set")
	}
	if view.HelpText == "" {
		t.Error("Expected help text to be set")
	}
	if !view.ShowFooterSeparator {
		t.Error("Expected footer separator to be enabled")
	}
	if !view.UseFullWidth {
		t.Error("Expected UseFullWidth to be true")
	}
}

func TestStandardView_Render_EmptyContent(t *testing.T) {
	termInfo := &terminal.Info{Width: 120, Height: 40}
	view := NewStandardView(termInfo)

	output := view.Render()

	// Empty content may result in empty or whitespace-only output
	// This is acceptable behavior
	_ = output // No assertion needed - just ensure it doesn't panic
}
