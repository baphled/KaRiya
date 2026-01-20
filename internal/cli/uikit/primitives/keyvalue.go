package primitives

import (
	"strings"

	"github.com/baphled/kariya/internal/cli/uikit/theme"
	"github.com/charmbracelet/lipgloss"
)

// KeyValue is a component for displaying label-value pairs with consistent alignment.
// It renders a muted label followed by a styled value, useful for detail views.
//
// Example:
//
//	kv := primitives.NewKeyValue(theme).
//	    LabelWidth(15).
//	    Add("Name:", "John Doe").
//	    Add("Email:", "john@example.com").
//	    Add("Role:", "Developer")
//	rendered := kv.Render()
type KeyValue struct {
	theme.Aware
	labelWidth int
	pairs      []keyValuePair
	separator  string
}

type keyValuePair struct {
	label string
	value string
	muted bool // Whether value should be muted (like timestamps)
}

// NewKeyValue creates a new key-value component with the given theme.
func NewKeyValue(th theme.Theme) *KeyValue {
	kv := &KeyValue{
		labelWidth: 12,
		pairs:      make([]keyValuePair, 0),
		separator:  "",
	}
	if th != nil {
		kv.SetTheme(th)
	}
	return kv
}

// LabelWidth sets the fixed width for labels.
// Returns the KeyValue for method chaining.
func (kv *KeyValue) LabelWidth(width int) *KeyValue {
	kv.labelWidth = width
	return kv
}

// Separator sets the separator string between pairs (e.g., newline for vertical list).
// Returns the KeyValue for method chaining.
func (kv *KeyValue) Separator(sep string) *KeyValue {
	kv.separator = sep
	return kv
}

// Add adds a label-value pair with the value in primary color.
// Returns the KeyValue for method chaining.
func (kv *KeyValue) Add(label, value string) *KeyValue {
	kv.pairs = append(kv.pairs, keyValuePair{label: label, value: value, muted: false})
	return kv
}

// AddMuted adds a label-value pair with the value in muted color.
// Useful for timestamps and less important information.
// Returns the KeyValue for method chaining.
func (kv *KeyValue) AddMuted(label, value string) *KeyValue {
	kv.pairs = append(kv.pairs, keyValuePair{label: label, value: value, muted: true})
	return kv
}

// AddBlank adds a blank line (empty pair) for spacing.
// Returns the KeyValue for method chaining.
func (kv *KeyValue) AddBlank() *KeyValue {
	kv.pairs = append(kv.pairs, keyValuePair{label: "", value: "", muted: false})
	return kv
}

// Clear removes all pairs.
// Returns the KeyValue for method chaining.
func (kv *KeyValue) Clear() *KeyValue {
	kv.pairs = make([]keyValuePair, 0)
	return kv
}

// Render returns the styled key-value pairs as a string.
func (kv *KeyValue) Render() string {
	if len(kv.pairs) == 0 {
		return ""
	}

	th := kv.Theme()
	if th == nil {
		th = theme.Default()
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(th.MutedColor()).
		Width(kv.labelWidth)

	valueStyle := lipgloss.NewStyle().
		Foreground(th.PrimaryColor()).
		Bold(true)

	mutedValueStyle := lipgloss.NewStyle().
		Foreground(th.MutedColor())

	lines := make([]string, 0, len(kv.pairs))
	for _, pair := range kv.pairs {
		if pair.label == "" && pair.value == "" {
			// Blank line for spacing
			lines = append(lines, "")
			continue
		}

		var valuePart string
		if pair.muted {
			valuePart = mutedValueStyle.Render(pair.value)
		} else {
			valuePart = valueStyle.Render(pair.value)
		}

		line := labelStyle.Render(pair.label) + valuePart
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}
