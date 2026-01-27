package bootstrap

import (
	"strings"

	"github.com/baphled/kariya/internal/config"
)

// IsProfileComplete checks if the profile has all required fields.
// Required fields are Name and Email - without these, CV generation cannot work.
func IsProfileComplete(cfg *config.Config) bool {
	if cfg == nil {
		return false
	}
	return strings.TrimSpace(cfg.Profile.Name) != "" &&
		strings.TrimSpace(cfg.Profile.Email) != ""
}
