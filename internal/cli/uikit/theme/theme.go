// Package theme provides theme infrastructure for UIKit components.
// It re-exports the themes.Theme interface and provides a default theme.
package theme

import (
	"github.com/baphled/kariya/internal/cli/themes"
)

// Theme is re-exported from internal/cli/themes for convenience.
// It defines the contract for all themes used in UIKit components.
type Theme = themes.Theme

// Default returns the default KaRiya theme.
// This is the theme used when no theme is explicitly provided.
func Default() Theme {
	return themes.NewDefaultTheme()
}
