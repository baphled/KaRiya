package components

import (
	"regexp"
	"strings"

	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/charmbracelet/lipgloss"
)

// SyntaxHighlighter provides format-aware syntax highlighting for export content.
type SyntaxHighlighter struct {
	theme themes.Theme
}

// NewSyntaxHighlighter creates a new syntax highlighter with the given theme.
func NewSyntaxHighlighter(theme themes.Theme) *SyntaxHighlighter {
	if theme == nil {
		theme = themes.NewDefaultTheme()
	}
	return &SyntaxHighlighter{theme: theme}
}

// HighlightJSON adds syntax highlighting to JSON content.
func (h *SyntaxHighlighter) HighlightJSON(content string) string {
	// Define styles
	keyStyle := lipgloss.NewStyle().Foreground(h.theme.AccentColor())
	stringStyle := lipgloss.NewStyle().Foreground(h.theme.SuccessColor())
	numberStyle := lipgloss.NewStyle().Foreground(h.theme.WarningColor())
	boolStyle := lipgloss.NewStyle().Foreground(h.theme.InfoColor())
	nullStyle := lipgloss.NewStyle().Foreground(h.theme.SecondaryColor())
	bracketStyle := lipgloss.NewStyle().Foreground(h.theme.PrimaryColor())

	lines := strings.Split(content, "\n")
	highlighted := make([]string, 0, len(lines))

	// Regex patterns
	keyPattern := regexp.MustCompile(`"([^"]+)"(\s*):`)
	stringPattern := regexp.MustCompile(`:\s*"([^"]*)"`)
	numberPattern := regexp.MustCompile(`:\s*(-?\d+\.?\d*)`)
	boolPattern := regexp.MustCompile(`:\s*(true|false)`)
	nullPattern := regexp.MustCompile(`:\s*(null)`)

	for _, line := range lines {
		// Highlight keys
		line = keyPattern.ReplaceAllStringFunc(line, func(match string) string {
			// Extract key name and colon
			parts := keyPattern.FindStringSubmatch(match)
			if len(parts) >= 3 {
				return keyStyle.Render("\""+parts[1]+"\"") + parts[2] + ":"
			}
			return match
		})

		// Highlight string values (after colon)
		line = stringPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := stringPattern.FindStringSubmatch(match)
			if len(parts) >= 2 {
				return ": " + stringStyle.Render("\""+parts[1]+"\"")
			}
			return match
		})

		// Highlight numbers
		line = numberPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := numberPattern.FindStringSubmatch(match)
			if len(parts) >= 2 {
				return ": " + numberStyle.Render(parts[1])
			}
			return match
		})

		// Highlight booleans
		line = boolPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := boolPattern.FindStringSubmatch(match)
			if len(parts) >= 2 {
				return ": " + boolStyle.Render(parts[1])
			}
			return match
		})

		// Highlight null
		line = nullPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := nullPattern.FindStringSubmatch(match)
			if len(parts) >= 2 {
				return ": " + nullStyle.Render(parts[1])
			}
			return match
		})

		// Highlight brackets
		line = strings.ReplaceAll(line, "{", bracketStyle.Render("{"))
		line = strings.ReplaceAll(line, "}", bracketStyle.Render("}"))
		line = strings.ReplaceAll(line, "[", bracketStyle.Render("["))
		line = strings.ReplaceAll(line, "]", bracketStyle.Render("]"))

		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

// HighlightYAML adds syntax highlighting to YAML content.
func (h *SyntaxHighlighter) HighlightYAML(content string) string {
	keyStyle := lipgloss.NewStyle().Foreground(h.theme.AccentColor())
	stringStyle := lipgloss.NewStyle().Foreground(h.theme.SuccessColor())
	numberStyle := lipgloss.NewStyle().Foreground(h.theme.WarningColor())
	boolStyle := lipgloss.NewStyle().Foreground(h.theme.InfoColor())
	commentStyle := lipgloss.NewStyle().Foreground(h.theme.SecondaryColor()).Italic(true)

	lines := strings.Split(content, "\n")
	highlighted := make([]string, 0, len(lines))

	keyPattern := regexp.MustCompile(`^(\s*)([a-zA-Z_][a-zA-Z0-9_]*)(\s*):`)
	listItemPattern := regexp.MustCompile(`^(\s*)-\s*`)
	commentPattern := regexp.MustCompile(`#.*$`)
	numberPattern := regexp.MustCompile(`:\s*(-?\d+\.?\d*)$`)
	boolPattern := regexp.MustCompile(`:\s*(true|false)$`)

	for _, line := range lines {
		// Highlight comments first
		if commentPattern.MatchString(line) {
			line = commentPattern.ReplaceAllStringFunc(line, func(match string) string {
				return commentStyle.Render(match)
			})
		}

		// Highlight keys
		line = keyPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := keyPattern.FindStringSubmatch(match)
			if len(parts) >= 4 {
				return parts[1] + keyStyle.Render(parts[2]) + parts[3] + ":"
			}
			return match
		})

		// Highlight list items
		line = listItemPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := listItemPattern.FindStringSubmatch(match)
			if len(parts) >= 2 {
				return parts[1] + stringStyle.Render("-") + " "
			}
			return match
		})

		// Highlight numbers
		line = numberPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := numberPattern.FindStringSubmatch(match)
			if len(parts) >= 2 {
				return ": " + numberStyle.Render(parts[1])
			}
			return match
		})

		// Highlight booleans
		line = boolPattern.ReplaceAllStringFunc(line, func(match string) string {
			parts := boolPattern.FindStringSubmatch(match)
			if len(parts) >= 2 {
				return ": " + boolStyle.Render(parts[1])
			}
			return match
		})

		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

// HighlightCSV adds syntax highlighting to CSV content.
func (h *SyntaxHighlighter) HighlightCSV(content string) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(h.theme.AccentColor())
	cellStyle := lipgloss.NewStyle().Foreground(h.theme.PrimaryColor())
	separatorStyle := lipgloss.NewStyle().Foreground(h.theme.SecondaryColor())

	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		return content
	}

	highlighted := make([]string, 0, len(lines))

	for i, line := range lines {
		if line == "" {
			highlighted = append(highlighted, line)
			continue
		}

		cells := strings.Split(line, ",")
		styledCells := make([]string, 0, len(cells))

		for _, cell := range cells {
			cell = strings.TrimSpace(cell)
			if i == 0 {
				// Header row
				styledCells = append(styledCells, headerStyle.Render(cell))
			} else {
				styledCells = append(styledCells, cellStyle.Render(cell))
			}
		}

		highlighted = append(highlighted, strings.Join(styledCells, separatorStyle.Render(", ")))
	}

	return strings.Join(highlighted, "\n")
}

// HighlightPlainText adds minimal highlighting to plain text.
func (h *SyntaxHighlighter) HighlightPlainText(content string) string {
	// For plain text, just highlight section headers (lines ending with :)
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(h.theme.AccentColor())
	separatorStyle := lipgloss.NewStyle().Foreground(h.theme.SecondaryColor())

	lines := strings.Split(content, "\n")
	highlighted := make([]string, 0, len(lines))

	for _, line := range lines {
		// Highlight separator lines
		if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "===") || strings.HasPrefix(line, "───") {
			highlighted = append(highlighted, separatorStyle.Render(line))
			continue
		}

		// Highlight section headers (lines that look like headers)
		if strings.HasSuffix(strings.TrimSpace(line), ":") && !strings.Contains(line, "  ") {
			highlighted = append(highlighted, headerStyle.Render(line))
			continue
		}

		highlighted = append(highlighted, line)
	}

	return strings.Join(highlighted, "\n")
}

// Highlight applies format-appropriate highlighting to content.
func (h *SyntaxHighlighter) Highlight(content, format string) string {
	switch strings.ToLower(format) {
	case "json":
		return h.HighlightJSON(content)
	case "yaml":
		return h.HighlightYAML(content)
	case "csv":
		return h.HighlightCSV(content)
	case "txt", "text", "markdown", "md":
		return h.HighlightPlainText(content)
	default:
		return content
	}
}
