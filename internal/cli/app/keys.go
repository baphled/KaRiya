package app

// Key input constants for keyboard handling.
// These represent the actual key strings received from tea.KeyMsg.String().
const (
	// Navigation keys.
	keyUp    = "up"
	keyDown  = "down"
	keyLeft  = "left"
	keyRight = "right"
	keyK     = "k" // vim-style up
	keyJ     = "j" // vim-style down
	keyH     = "h" // vim-style left
	keyL     = "l" // vim-style right

	// Action keys.
	keyEnter = "enter"
	keySpace = " "
	keyEsc   = "esc"
	keyTab   = "tab"

	// Global keys.
	keyQuit   = "q"
	keyHelp   = "?"
	keyCtrlC  = "ctrl+c"
	keyCtrlO  = "ctrl+o"
	keyCtrlD  = "ctrl+d"
	keyCtrlU  = "ctrl+u"
	keyPgUp   = "pgup"
	keyPgDown = "pgdown"
)
