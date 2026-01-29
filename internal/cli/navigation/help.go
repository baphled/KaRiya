package navigation

import (
	"fmt"
	"strings"
)

// HelpText represents formatted help information for a set of navigation keys.
type HelpText struct {
	Keys     []NavigationKey
	Template string
	Compact  bool
}

// GetHelpText returns a formatted string displaying navigation keys with descriptions.
// Returns full format with one key per line.
func GetHelpText(keys []NavigationKey) string {
	return GetHelpTextCompact(keys, false)
}

// GetHelpTextCompact returns formatted help text with optional compact formatting.
// If compact is true, returns short format (Key: Description | Key: Description).
// If compact is false, returns full format with one key per line.
func GetHelpTextCompact(keys []NavigationKey, compact bool) string {
	if len(keys) == 0 {
		return ""
	}

	var parts []string

	for _, key := range keys {
		description, exists := KeyDescription[key]
		if !exists {
			description = string(key)
		}

		if compact {
			parts = append(parts, fmt.Sprintf("%s: %s", key, description))
		} else {
			parts = append(parts, fmt.Sprintf("%s → %s", key, description))
		}
	}

	if compact {
		return strings.Join(parts, " | ")
	}
	return strings.Join(parts, "\n")
}

// GetContextualHelp returns help text for a specific screen context
// Context can be "form", "list", "metadata_review", "bulk_operations", etc.
func GetContextualHelp(context string) string {
	contextKeyMap := map[string][]NavigationKey{
		"form": {
			KeyUp,
			KeyDown,
			KeySelect,
			KeyToggle,
			KeyBack,
		},
		"list": {
			KeyUp,
			KeyDown,
			KeySelect,
			KeyEdit,
			KeyDelete,
			KeyFilter,
			KeySort,
			KeySearch,
			KeyBack,
		},
		"cv_config_manager": {
			KeyUp,
			KeyDown,
			KeySelect,
			KeyEdit,
			KeyDelete,
			KeyBack,
		},
		"metadata_review": {
			KeyUp,
			KeyDown,
			KeySelect,
			KeyEdit,
			KeyBulk,
			KeyFilter,
			KeySort,
			KeyBack,
		},
		"bulk_operations": {
			KeyUp,
			KeyDown,
			KeyToggle,
			KeySelect,
			KeyBack,
		},
		"home": {
			KeyCapture,
			KeyList,
			KeyMetadata,
			KeyHelp,
			KeyQuit,
		},
		"default": {
			KeyUp,
			KeyDown,
			KeySelect,
			KeyHelp,
			KeyBack,
		},
	}

	keys, exists := contextKeyMap[context]
	if !exists {
		keys = contextKeyMap["default"]
	}

	return GetHelpTextCompact(keys, true)
}

// GetFullHelp returns comprehensive help text for all navigation keys.
func GetFullHelp() string {
	return GetHelpTextCompact(AllNavigationKeys(), false)
}

// GetCompactHelp returns compact help text for all navigation keys (suitable for footers).
func GetCompactHelp() string {
	return GetHelpTextCompact(AllNavigationKeys(), true)
}

// GetGroupedHelp returns help text organized by key groups.
func GetGroupedHelp() string {
	groups := map[string][]NavigationKey{
		"Navigation": {KeyUp, KeyDown, KeyLeft, KeyRight},
		"Actions":    {KeySelect, KeyToggle, KeyAdd, KeyEdit, KeyDelete},
		"Modes":      {KeyCapture, KeyList, KeyMetadata, KeyPending, KeyFacts, KeyGenerate, KeyCV},
		"Tools":      {KeyFilter, KeySort, KeySearch, KeyBulk},
		"Global":     {KeyHelp, KeyQuit, KeyBack},
	}

	var result strings.Builder
	for groupName, keys := range groups {
		result.WriteString(fmt.Sprintf("\n%s:\n", groupName))
		for _, key := range keys {
			description, exists := KeyDescription[key]
			if !exists {
				description = string(key)
			}
			result.WriteString(fmt.Sprintf("  %s → %s\n", key, description))
		}
	}

	return strings.TrimPrefix(result.String(), "\n")
}
