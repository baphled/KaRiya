package export

import (
	"fmt"
	"strings"

	"github.com/baphled/kariya/internal/cli/screens"
	"github.com/baphled/kariya/internal/cli/screens/base"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/types"
	"github.com/baphled/kariya/internal/cli/uikit/containers"
	"github.com/baphled/kariya/internal/cli/uikit/primitives"
	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PreviewStats contains statistics about the preview content.
type PreviewStats struct {
	ItemCount     int   // Number of items being exported (shown in preview)
	TotalCount    int   // Total items available in database
	EstimatedSize int64 // Estimated file size in bytes
	ContentLines  int   // Number of lines in content
	IsTruncated   bool  // Whether preview was truncated
}

// Preview displays a preview of the export content with a scrollable viewport.
type Preview struct {
	*base.BaseScreen

	content      string
	artifactType types.ExportArtifactType
	format       types.ExportFormat
	destination  types.ExportDestination
	breadcrumbs  []string
	viewport     viewport.Model
	ready        bool
	width        int
	height       int
	stats        *PreviewStats
}

// NewPreview creates a new preview screen with viewport support.
func NewPreview(content string, artifactType types.ExportArtifactType, format types.ExportFormat, breadcrumbs []string) *Preview {
	return &Preview{
		BaseScreen:   base.NewBaseScreen(),
		content:      content,
		artifactType: artifactType,
		format:       format,
		destination:  types.ExportDestinationFile, // Default
		breadcrumbs:  breadcrumbs,
		width:        80,
		height:       24,
		ready:        false,
		stats:        nil,
	}
}

// NewPreviewWithStats creates a new preview screen with statistics.
func NewPreviewWithStats(content string, artifactType types.ExportArtifactType, format types.ExportFormat, destination types.ExportDestination, breadcrumbs []string, stats *PreviewStats) *Preview {
	return &Preview{
		BaseScreen:   base.NewBaseScreen(),
		content:      content,
		artifactType: artifactType,
		format:       format,
		destination:  destination,
		breadcrumbs:  breadcrumbs,
		width:        80,
		height:       24,
		ready:        false,
		stats:        stats,
	}
}

// Init initializes the screen.
func (s *Preview) Init() tea.Cmd {
	return nil
}

// Update handles messages.
func (s *Preview) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		s.BaseScreen.HandleWindowSizeMsg(msg)
		s.width = msg.Width
		s.height = msg.Height
		s.ready = false // Force viewport recreation on resize
		return nil, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			// Go back
			return nil, &screens.CancelResult{}

		case "enter":
			// Continue to confirmation
			return nil, &screens.NavigateResult{
				ResultData: "confirm",
			}

		// Vim-style navigation
		case "g":
			// Go to top
			if s.ready {
				s.viewport.GotoTop()
			}
			return nil, nil

		case "G":
			// Go to bottom
			if s.ready {
				s.viewport.GotoBottom()
			}
			return nil, nil

		// Scrolling keys - pass to viewport
		case "up", "k", "down", "j", "pgup", "pgdown", "ctrl+u", "ctrl+d":
			if s.ready {
				s.viewport, cmd = s.viewport.Update(msg)
				return cmd, nil
			}
		}
	}

	return nil, nil
}

// View renders the screen.
func (s *Preview) View() string {
	th := s.getTheme()
	uikitTheme := theme.Default()
	var b strings.Builder

	// Render summary header
	b.WriteString(s.renderSummaryHeader(th))
	b.WriteString("\n")

	// Initialize highlighter and viewport if needed
	if !s.ready {
		s.initializeViewport(th)
	}

	// Wrap viewport in a code block style box
	// Calculate box width: viewport width + padding (2 sides * 1) + border (2 sides * 1)
	boxWidth := s.viewport.Width + 4
	codeBlockTitle := s.getCodeBlockTitle()
	codeBlock := containers.NewBox(uikitTheme).
		Title(codeBlockTitle).
		Content(s.viewport.View()).
		Variant(containers.BoxSubtle).
		Width(boxWidth).
		Padding(1).
		Render()

	b.WriteString(codeBlock)
	b.WriteString("\n")

	// Footer with scroll indicator
	b.WriteString(s.renderFooter(th))

	return b.String()
}

// renderSummaryHeader renders the summary/statistics header.
func (s *Preview) renderSummaryHeader(th themes.Theme) string {
	var b strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(th.PrimaryColor())
	b.WriteString(titleStyle.Render(fmt.Sprintf("Export Preview: %s → %s",
		getArtifactTypeLabel(s.artifactType),
		getFormatLabel(s.format))))
	b.WriteString("\n")

	// Stats line - simple key: value format
	statsLine := s.renderStatsLine(th)
	if statsLine != "" {
		b.WriteString(statsLine)
		b.WriteString("\n")
	}

	// Separator
	sepStyle := lipgloss.NewStyle().Foreground(th.SecondaryColor())
	b.WriteString(sepStyle.Render(strings.Repeat("─", minInt(70, s.width-4))))
	b.WriteString("\n")

	return b.String()
}

// renderStatsLine renders the statistics as a simple formatted line.
func (s *Preview) renderStatsLine(th themes.Theme) string {
	labelStyle := lipgloss.NewStyle().Foreground(th.SecondaryColor())
	valueStyle := lipgloss.NewStyle().Foreground(th.PrimaryColor())

	var parts []string

	// Destination
	destLabel := "File"
	if s.destination == types.ExportDestinationClipboard {
		destLabel = "Clipboard"
	}
	parts = append(parts, labelStyle.Render("To: ")+valueStyle.Render(destLabel))

	// Stats if available
	if s.stats != nil {
		if s.stats.ItemCount > 0 {
			itemsLabel := fmt.Sprintf("%d items", s.stats.ItemCount)
			if s.stats.TotalCount > s.stats.ItemCount {
				itemsLabel = fmt.Sprintf("%d of %d items", s.stats.ItemCount, s.stats.TotalCount)
			}
			parts = append(parts, labelStyle.Render("Items: ")+valueStyle.Render(itemsLabel))
		}

		if s.stats.EstimatedSize > 0 {
			parts = append(parts, labelStyle.Render("Size: ")+valueStyle.Render("~"+formatFileSize(s.stats.EstimatedSize)))
		}

		if s.stats.IsTruncated {
			warnStyle := lipgloss.NewStyle().Foreground(th.WarningColor())
			parts = append(parts, warnStyle.Render("(preview truncated)"))
		}
	}

	return strings.Join(parts, "  │  ")
}

// initializeViewport sets up the viewport with content.
func (s *Preview) initializeViewport(_ themes.Theme) {
	// Calculate viewport dimensions
	// Account for:
	// - Header: ~5 lines (title, stats, separator, blank line)
	// - Footer: 2 lines (separator, help badges)
	// - Box wrapper: 2 (border) + 2 (padding) = 4 lines vertical, 6 chars horizontal
	//   (border=1 each side, padding=1 each side, plus padding internal spacing)
	viewportHeight := s.height - 11 // 5 header + 2 footer + 4 box
	if viewportHeight < 5 {
		viewportHeight = 5
	}

	// Width: screen width minus margins (2) minus box border (2) minus box padding (2) minus content padding (2)
	viewportWidth := s.width - 8
	if viewportWidth < 40 {
		viewportWidth = 40
	}

	s.viewport = viewport.New(viewportWidth, viewportHeight)
	s.viewport.SetContent(s.content)
	s.ready = true
}

// RenderContent returns just the content without StandardView wrapper.
func (s *Preview) RenderContent() string {
	return s.View()
}

// renderFooter renders the footer with help text and scroll indicator using UIKit primitives.
func (s *Preview) renderFooter(th themes.Theme) string {
	// Convert themes.Theme to uikit theme.Theme for primitives
	uikitTheme := theme.Default()

	// Build badges for footer
	badges := []*primitives.Badge{
		primitives.HelpKeyBadge("↑↓/jk", "Scroll", uikitTheme),
		primitives.HelpKeyBadge("g/G", "Top/Bottom", uikitTheme),
	}

	// Add scroll percentage if scrollable
	if s.ready && s.viewport.TotalLineCount() > s.viewport.Height {
		pct := int(s.viewport.ScrollPercent() * 100)
		badges = append(badges, primitives.HelpKeyBadge(fmt.Sprintf("[%d%%]", pct), "", uikitTheme))
	}

	badges = append(badges,
		primitives.HelpKeyBadge("Enter", "Export", uikitTheme),
		primitives.BackBadge(uikitTheme),
	)

	// Render separator line and badges
	footerStyle := lipgloss.NewStyle().Foreground(th.SecondaryColor())
	separator := footerStyle.Render(strings.Repeat("─", minInt(70, s.width-4)))

	return separator + "\n" + primitives.RenderHelpFooter(uikitTheme, badges...)
}

// getTheme returns the theme from BaseScreen or a default theme.
func (s *Preview) getTheme() themes.Theme {
	if t := s.BaseScreen.Theme(); t != nil {
		if theme, ok := t.(themes.Theme); ok {
			return theme
		}
	}
	return themes.NewDefaultTheme()
}

// SetStats sets the preview statistics.
func (s *Preview) SetStats(stats *PreviewStats) {
	s.stats = stats
	s.ready = false // Force re-render
}

// SetDestination sets the export destination.
func (s *Preview) SetDestination(dest types.ExportDestination) {
	s.destination = dest
}

// getCodeBlockTitle returns a code block title indicating the format (like markdown code fences).
func (s *Preview) getCodeBlockTitle() string {
	switch s.format {
	case types.ExportFormatJSON:
		return "json"
	case types.ExportFormatYAML:
		return "yaml"
	case types.ExportFormatCSV:
		return "csv"
	case types.ExportFormatMD:
		return "markdown"
	case types.ExportFormatTXT:
		return "text"
	default:
		return string(s.format)
	}
}

// Helper functions

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func getArtifactTypeLabel(t types.ExportArtifactType) string {
	switch t {
	case types.ExportTypeEvents:
		return "Career Events"
	case types.ExportTypeFacts:
		return "Facts"
	case types.ExportTypeBursts:
		return "Bursts"
	case types.ExportTypeCV:
		return "CV"
	case types.ExportTypeProfile:
		return "Profile"
	default:
		return string(t)
	}
}

func getFormatLabel(f types.ExportFormat) string {
	switch f {
	case types.ExportFormatJSON:
		return "JSON"
	case types.ExportFormatYAML:
		return "YAML"
	case types.ExportFormatCSV:
		return "CSV"
	case types.ExportFormatTXT:
		return "Text"
	case types.ExportFormatMD:
		return "Markdown"
	default:
		return string(f)
	}
}

func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
	)

	switch {
	case bytes >= MB:
		return fmt.Sprintf("%.1f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.1f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
