package career

import "errors"

// Sentinel errors for optional feature repositories
// These errors indicate that an optional feature (fact/burst storage) is not
// configured, which is a valid application state, not an error condition.
var (
	// ErrFactRepositoryNotConfigured indicates that the fact repository is not set up.
	// This is expected when fact storage features are disabled or initialization failed.
	ErrFactRepositoryNotConfigured = errors.New("fact repository not configured")

	// ErrBurstRepositoryNotConfigured indicates that the burst repository is not set up.
	// This is expected when burst storage features are disabled or initialization failed.
	ErrBurstRepositoryNotConfigured = errors.New("burst repository not configured")

	// ErrSkillRepositoryNotConfigured indicates that the skill repository is not set up.
	// This is expected when skill storage features are disabled or initialisation failed.
	ErrSkillRepositoryNotConfigured = errors.New("skill repository not configured")
)
