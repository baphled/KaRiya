package configure

import (
	"github.com/baphled/kariya/internal/cli/configtypes"
)

// ConfigurationDomain re-exports configtypes.ConfigurationDomain for the configure package.
type ConfigurationDomain = configtypes.ConfigurationDomain

// ConfigurationSetting re-exports configtypes.ConfigurationSetting for the configure package.
type ConfigurationSetting = configtypes.ConfigurationSetting

// ConfigurationState re-exports configtypes.ConfigurationState for the configure package.
type ConfigurationState = configtypes.ConfigurationState

// ConfigurationChanges re-exports configtypes.ConfigurationChanges for the configure package.
type ConfigurationChanges = configtypes.ConfigurationChanges

// DomainSystem, DomainProfile, DomainExport, and DomainUI re-export configtypes domain constants.
const (
	DomainSystem  = configtypes.DomainSystem
	DomainProfile = configtypes.DomainProfile
	DomainExport  = configtypes.DomainExport
	DomainUI      = configtypes.DomainUI
)

// ConfigStateSelectDomain and related constants re-export configtypes workflow state values.
const (
	ConfigStateSelectDomain  = configtypes.ConfigStateSelectDomain
	ConfigStateEditSettings  = configtypes.ConfigStateEditSettings
	ConfigStateReviewChanges = configtypes.ConfigStateReviewChanges
	ConfigStateConfirm       = configtypes.ConfigStateConfirm
	ConfigStateSaving        = configtypes.ConfigStateSaving
	ConfigStateComplete      = configtypes.ConfigStateComplete
	ConfigStateFailed        = configtypes.ConfigStateFailed
)
