package app

import (
	"github.com/baphled/kariya/internal/domain/career"
)

// FormSubmittedMsg is sent when a form is successfully submitted
type FormSubmittedMsg struct {
	Event *career.CareerEvent
	Err   error
}

// NavigateMsg is sent to navigate to a specific screen
type NavigateMsg struct {
	Screen Screen
}

// SuccessNavigateMsg is sent from success screen for navigation
type SuccessNavigateMsg struct {
	Action string // "capture_another", "view_list", "exit"
}

