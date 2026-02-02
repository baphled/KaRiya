// Package screens provides UI views for the TUI application.
//
// # Overview
//
// The screens package contains all screen implementations used by intents.
// Each screen is a stateless view component that handles rendering and
// user input for a specific UI state.
//
// # Architecture
//
// Screens follow the naming convention:
//   - Package: screens/{feature}/ (e.g., screens/facts/)
//   - Files: list.go, detail.go, form.go, modals/
//   - Structs: {Entity}{Type}Screen (e.g., FactListScreen)
//
// # Screen Types
//
//   - ListScreen: Table-based list views
//   - DetailScreen: Single item detail views
//   - FormScreen: Input forms
//   - DeleteScreen: Confirmation dialogs
//   - SelectScreen: Item selection
//
// # Usage
//
// Screens return ScreenResult to communicate with intents:
//
//	func (s *MyScreen) Update(msg tea.Msg) (tea.Cmd, screens.ScreenResult) {
//	    if key == "esc" {
//	        return nil, screens.NewCancelResult("")
//	    }
//	    return nil, nil
//	}
//
// For detailed documentation, see docs/conventions/SCREEN_NAMING.md.
package screens
