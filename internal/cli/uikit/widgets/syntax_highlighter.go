// Package widgets provides higher-level composite components that combine
// primitives and behaviors into reusable UI patterns.
package widgets

import (
	"regexp"
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

// Pre-compiled regex patterns for syntax highlighting (avoids recompilation in loops)
var (
	jsonNumberPattern = regexp.MustCompile(`^-?\d+\.?\d*(?:[eE][+-]?\d+)?$`)
)

// SyntaxHighlighter provides format-aware syntax highlighting for export content.
type SyntaxHighlighter struct {
	theme.Aware
}

// NewSyntaxHighlighter creates a new syntax highlighter with the given theme.
func NewSyntaxHighlighter(th theme.Theme) *SyntaxHighlighter {
	h := &SyntaxHighlighter{}
	if th != nil {
		h.SetTheme(th)
	}
	return h
}

// HighlightJSON adds syntax highlighting to JSON content.
// It processes each line token by token to avoid regex interference.
func (h *SyntaxHighlighter) HighlightJSON(content string) string {
	// Define styles
	keyStyle := lipgloss.NewStyle().Foreground(h.AccentColor())
	stringStyle := lipgloss.NewStyle().Foreground(h.SuccessColor())
	numberStyle := lipgloss.NewStyle().Foreground(h.WarningColor())
	boolStyle := lipgloss.NewStyle().Foreground(h.InfoColor())
	nullStyle := lipgloss.NewStyle().Foreground(h.SecondaryColor())
	bracketStyle := lipgloss.NewStyle().Foreground(h.PrimaryColor())
	punctStyle := lipgloss.NewStyle().Foreground(h.SecondaryColor())

	lines := strings.Split(content, "\n")
	highlighted := make([]string, 0, len(lines))

	// Single comprehensive pattern that captures JSON tokens
	// Groups: 1=key, 2=string value, 3=number, 4=bool, 5=null, 6=bracket/punct
	tokenPattern := regexp.MustCompile(
		`("([^"\\]|\\.)*")\s*:` + // Key followed by colon
			`|:\s*("([^"\\]|\\.)*")` + // String value after colon
			`|:\s*(-?\d+\.?\d*(?:[eE][+-]?\d+)?)` + // Number after colon
			`|:\s*(true|false)` + // Boolean after colon
			`|:\s*(null)` + // Null after colon
			`|([{}\[\],])`, // Brackets and punctuation
	)

	for _, line := range lines {
		result := tokenPattern.ReplaceAllStringFunc(line, func(match string) string {
			// Determine what type of token this is
			trimmed := strings.TrimSpace(match)

			// Key: "something":
			if strings.HasSuffix(trimmed, ":") && strings.HasPrefix(trimmed, "\"") {
				// Extract key and preserve spacing
				keyEnd := strings.LastIndex(match, "\":")
				if keyEnd > 0 {
					key := match[:keyEnd+1]
					rest := match[keyEnd+1:]
					return keyStyle.Render(key) + rest
				}
				return keyStyle.Render(strings.TrimSuffix(trimmed, ":")) + ":"
			}

			// String value: : "something"
			if strings.HasPrefix(trimmed, ":") {
				valueStart := strings.Index(match, ":")
				prefix := match[:valueStart+1]
				value := strings.TrimSpace(match[valueStart+1:])

				if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
					return prefix + " " + stringStyle.Render(value)
				}

				// Number
				if jsonNumberPattern.MatchString(value) {
					return prefix + " " + numberStyle.Render(value)
				}

				// Boolean
				if value == "true" || value == "false" {
					return prefix + " " + boolStyle.Render(value)
				}

				// Null
				if value == "null" {
					return prefix + " " + nullStyle.Render(value)
				}
			}

			// Brackets and punctuation
			switch trimmed {
			case "{", "}":
				return bracketStyle.Render(trimmed)
			case "[", "]":
				return bracketStyle.Render(trimmed)
			case ",":
				return punctStyle.Render(trimmed)
			}

			return match
		})

		highlighted = append(highlighted, result)
	}

	return strings.Join(highlighted, "\n")
}

// HighlightYAML adds syntax highlighting to YAML content.
func (h *SyntaxHighlighter) HighlightYAML(content string) string {
	keyStyle := lipgloss.NewStyle().Foreground(h.AccentColor())
	stringStyle := lipgloss.NewStyle().Foreground(h.SuccessColor())
	numberStyle := lipgloss.NewStyle().Foreground(h.WarningColor())
	boolStyle := lipgloss.NewStyle().Foreground(h.InfoColor())
	commentStyle := lipgloss.NewStyle().Foreground(h.SecondaryColor()).Italic(true)

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
		Foreground(h.AccentColor())
	cellStyle := lipgloss.NewStyle().Foreground(h.PrimaryColor())
	separatorStyle := lipgloss.NewStyle().Foreground(h.SecondaryColor())

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
		Foreground(h.AccentColor())
	separatorStyle := lipgloss.NewStyle().Foreground(h.SecondaryColor())

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
