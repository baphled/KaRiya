package bootstrap

import (
	"strings"

	"github.com/baphled/kariya/internal/config"
)

// IsProfileComplete checks if the profile has all required fields.
//
// Expected:
//   - config must be a valid configuration object.
//
// Returns:
//   - A bool value.
//
// Side effects:
//   - None.
func IsProfileComplete(cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	return strings.TrimSpace(cfg.Profile.Name) != "" &&
		strings.TrimSpace(cfg.Profile.Email) != ""
}
