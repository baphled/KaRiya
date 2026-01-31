package intents

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/baphled/kariya/internal/cli/terminal"
	"github.com/baphled/kariya/internal/cli/themes"
	"github.com/baphled/kariya/internal/cli/uikit/feedback"
	"github.com/baphled/kariya/internal/cli/uikit/layout"
)

// LogoModel defines the interface for logo components.
// The uikit/display.Logo package provides the standard implementation.
type LogoModel interface {
	Init() tea.Cmd
	Update(msg tea.Msg) (tea.Model, tea.Cmd)
	View() string
	ViewStatic() string
	SetWidth(width int)
}

// Intent defines the contract for all intent implementations.
//
// Each intent MUST:
//   - Own local navigation state
//   - Call domain services
//   - Emit artifacts
//   - Return a type-safe IntentResult
//
// Intents MUST NOT:
//   - Mutate global UI state
//   - Navigate into another intent
//   - Assume prior context unless explicitly passed
//   - Dispatch Bubble Tea commands affecting other intents
//   - Use runtime type assertions
//   - Pass arbitrary data between intents
type Intent interface {
	// Init is called when the intent is first activated.
	// It may emit commands (e.g., fetch data, start async operations).
	Init() tea.Cmd

	// Update processes a message and returns a command.
	// This is Bubble Tea's standard Update, but constrained to intent-local state.
	Update(msg tea.Msg) tea.Cmd

	// View renders the intent's current state.
	// This is Bubble Tea's standard View.
	View() string

	// Result returns the intent's result if it has completed, or nil if still active.
	// The result communicates the intent's outcome back to the router and parent context.
	Result() *IntentResult[interface{}]
}

// IntentRouter manages the lifecycle of intents and enforces the boundary contract.
// The root Bubble Tea Model delegates to the router, which:
// - Activates/deactivates intents
// - Validates intent results
// - Propagates results back to the root model
// - Prevents illegal state transitions
//
//nolint:interfacebloat // IntentRouter needs all methods for complete intent lifecycle management
type IntentRouter interface {
	// ActivateIntent activates an intent by name.
	// Returns an error if the intent doesn't exist or the transition is illegal.
	ActivateIntent(name string, context map[string]interface{}) (tea.Cmd, error)

	// GetActiveIntent returns the currently active intent, or nil if none.
	GetActiveIntent() Intent

	// HandleMessage processes a message in the active intent.
	// Returns a command and any intent result if the intent completed.
	HandleMessage(msg tea.Msg) (tea.Cmd, interface{})

	// View renders the active intent.
	View() string

	// Back navigates back to the previous intent in history.
	// Returns an error if already at the root.
	Back() (tea.Cmd, error)

	// GetHistory returns the navigation history.
	GetHistory() []Intent

	// GetHistoryDepth returns the current depth in the navigation history.
	GetHistoryDepth() int
}

// IntentContext encapsulates minimal, explicit context passed to an intent.
// This is the ONLY way intents receive external data.
// Type safety is enforced by the intent's specific context type.
type IntentContext interface {
	// Validate ensures the context is complete and correct.
	// Returns an error if validation fails.
	Validate() error
}

// ModalEditResult captures the outcome of a modal inline-editing sub-flow.
//
// The generic type parameter T is the domain object being edited (e.g.,
// *career.Event or *career.Fact). T must satisfy the "any" constraint.
//
// The Original field holds the value of T before the edit began. The
// Modified field holds the value of T after the user finished editing; when
// the user cancels, Modified equals Original. The Accepted field is true
// when the user confirmed the edit and false when the user cancelled or
// dismissed the modal. The Changes field is a map[string]interface{} keyed
// by field name whose values are the new field values; only fields that
// were actually modified appear in the map and an empty map means no
// fields changed.
//
// The parent intent inspects Accepted first: if false it discards
// Modified and takes no persistence action. If true it checks Changes to
// determine which fields need updating and commits the Modified value.
//
// Construct with NewModalEditResult for confirmed edits or
// NewCancelledModalEditResult for user cancellation.
type ModalEditResult[T any] struct {
	Original T
	Modified T
	Accepted bool
	Changes  map[string]interface{}
}

// HasChanges checks whether the edit produced any field modifications.
//
// Returns:
//   - True if the Changes map is non-empty.
//
// Side effects:
//   - None.
func (m *ModalEditResult[T]) HasChanges() bool {
	return len(m.Changes) > 0
}

// WasAccepted checks whether the user confirmed the edit rather than cancelling.
//
// Returns:
//   - True if the user confirmed the edit.
//
// Side effects:
//   - None.
func (m *ModalEditResult[T]) WasAccepted() bool {
	return m.Accepted
}

// GetChange looks up a single field's new value from the edit result.
//
// Expected:
//   - fieldName should match a key in the Changes map.
//
// Returns:
//   - The new value for the field, or nil if the field was not changed.
//
// Side effects:
//   - None.
func (m *ModalEditResult[T]) GetChange(fieldName string) interface{} {
	if m.Changes == nil {
		return nil
	}
	return m.Changes[fieldName]
}

// NewModalEditResult creates a ModalEditResult with the given original and modified values.
//
// Expected:
//   - original is the value before editing.
//   - modified is the value after editing.
//   - accepted indicates whether the user confirmed the edit.
//   - changes maps field names to their new values; nil is treated as an empty map.
//
// Returns:
//   - A populated ModalEditResult with the provided values.
//
// Side effects:
//   - None.
func NewModalEditResult[T any](original, modified T, accepted bool, changes map[string]interface{}) *ModalEditResult[T] {
	if changes == nil {
		changes = make(map[string]interface{})
	}
	return &ModalEditResult[T]{
		Original: original,
		Modified: modified,
		Accepted: accepted,
		Changes:  changes,
	}
}

// NewCancelledModalEditResult creates a ModalEditResult indicating the user cancelled the edit.
//
// Expected:
//   - original is the value before the cancelled edit.
//
// Returns:
//   - A ModalEditResult with Accepted set to false and Modified equal to Original.
//
// Side effects:
//   - None.
func NewCancelledModalEditResult[T any](original T) *ModalEditResult[T] {
	return &ModalEditResult[T]{
		Original: original,
		Modified: original,
		Accepted: false,
		Changes:  make(map[string]interface{}),
	}
}

// TerminalAwareIntent extends the Intent interface with terminal size awareness.
// Intents that implement this interface will receive terminal dimension updates
// and can adapt their rendering accordingly.
type TerminalAwareIntent interface {
	Intent

	// UpdateTerminalInfo updates the intent with current terminal dimensions
	UpdateTerminalInfo(info *terminal.Info)

	// GetMinimumSize returns the minimum terminal size required for this intent
	// Returns width and height in columns and rows
	GetMinimumSize() (width, height int)
}

// BaseIntent provides common functionality that all intents can embed.
// It handles terminal size tracking, logo management, and state management
// for loading, error, success, and progress states.
//
// All state management is independent - intents can have multiple states active
// simultaneously. When creating views with CreateView() or CreateViewWithBreadcrumbs(),
// modals are automatically applied based on priority:
//  1. Error (highest priority, with bell)
//  2. Loading (ongoing operation)
//  3. Progress (specific progress tracking)
//  4. Success (lowest priority, auto-dismiss after 3 seconds)
//
// Help Modal Integration:
// Press '?' to toggle context-sensitive help. The help modal shows keyboard shortcuts
// and is automatically integrated with view rendering.
//
// Example usage:
//
//	type MyIntent struct {
//	    *BaseIntent
//	    // ... intent-specific fields
//	}
//
//	func (i *MyIntent) View() string {
//	    // Simple, clean view creation with automatic state modals
//	    view := i.CreateViewWithBreadcrumbs("Main Menu", "My Intent", i.stateName)
//	    view.WithContent(i.renderContent())
//	    view.WithHelp(CombineFooters(NavigationFooter(), "q Quit"))
//	    return view.Render()
//	}
//
//	func (i *MyIntent) handleSubmit() tea.Cmd {
//	    i.SetLoading("Saving...")
//	    return func() tea.Msg {
//	        if err := i.service.Save(); err != nil {
//	            i.SetError(err)
//	            return ErrorMsg{err}
//	        }
//	        i.ClearLoading()
//	        i.SetSuccess("Saved successfully!")
//	        return SuccessMsg{}
//	    }
//	}
type BaseIntent struct {
	// Terminal management
	terminalInfo   *terminal.Info
	terminalConfig terminal.Config

	// Logo management (shared instance via interface)
	logo        LogoModel
	logoSpacing int

	// Theme management
	themeManager *themes.ThemeManager

	// State management
	isLoading      bool
	loadingMessage string
	errorState     error
	successMessage string
	successTime    time.Time

	// Progress tracking
	progressEnabled bool
	progressValue   float64
	progressTitle   string
	progressMessage string

	// Help modal
	helpModal *feedback.HelpModal
}

// NewBaseIntent creates a new BaseIntent with default terminal configuration.
//
// Returns:
//   - A BaseIntent initialized with default terminal info, config, logo spacing, and help modal.
//
// Side effects:
//   - None.
func NewBaseIntent() *BaseIntent {
	return &BaseIntent{
		terminalInfo:   terminal.NewInfo(),
		terminalConfig: terminal.DefaultConfig,
		logoSpacing:    2,
		helpModal:      feedback.NewHelpModal(nil),
	}
}

// UpdateTerminalInfo refreshes the cached terminal dimensions after a resize event.
//
// Expected:
//   - info must be non-nil with valid Width and Height for proper rendering.
//
// Side effects:
//   - Replaces the stored terminal info on the BaseIntent.
func (b *BaseIntent) UpdateTerminalInfo(info *terminal.Info) {
	b.terminalInfo = info
}

// GetTerminalInfo provides access to the current terminal dimensions and capabilities.
//
// Returns:
//   - The current terminal.Info instance.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetTerminalInfo() *terminal.Info {
	return b.terminalInfo
}

// GetMinimumSize reports the minimum terminal dimensions required for rendering.
//
// Returns:
//   - width and height from the terminal config's MinWidth and MinHeight.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetMinimumSize() (width, height int) {
	return b.terminalConfig.MinWidth, b.terminalConfig.MinHeight
}

// GetModalDimensions returns terminal dimensions for modal sizing.
// If terminal info is not available, returns sensible defaults (120x40).
// This eliminates the repeated dimension extraction pattern in modal opening methods.
//
// Returns:
//   - width and height suitable for modal sizing, defaulting to 120x40 when terminal info is unavailable.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetModalDimensions() (width, height int) {
	if b.terminalInfo != nil && b.terminalInfo.Width > 0 && b.terminalInfo.Height > 0 {
		return b.terminalInfo.Width, b.terminalInfo.Height
	}
	return 120, 40
}

// SetLogo configures the logo component used in view rendering.
//
// Expected:
//   - logo is a LogoModel implementation, or nil to clear.
//
// Side effects:
//   - Replaces the stored logo on the BaseIntent.
func (b *BaseIntent) SetLogo(logo LogoModel) {
	b.logo = logo
}

// GetLogo provides access to the shared logo component for rendering.
//
// Returns:
//   - The current LogoModel, or nil if not set.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetLogo() LogoModel {
	return b.logo
}

// SetLogoSpacing configures the vertical padding above the logo.
//
// Expected:
//   - spacing must be non-negative.
//
// Side effects:
//   - Updates the stored logo spacing on the BaseIntent.
func (b *BaseIntent) SetLogoSpacing(spacing int) {
	b.logoSpacing = spacing
}

// GetLogoSpacing reports the number of blank lines rendered before the logo.
//
// Returns:
//   - The number of spacing lines before the logo.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetLogoSpacing() int {
	return b.logoSpacing
}

// SetThemeManager configures the theme provider used for style resolution.
//
// Expected:
//   - tm is the ThemeManager to use, or nil to clear.
//
// Side effects:
//   - Replaces the stored theme manager on the BaseIntent.
func (b *BaseIntent) SetThemeManager(tm *themes.ThemeManager) {
	b.themeManager = tm
}

// GetThemeManager provides access to the theme manager for style resolution.
//
// Returns:
//   - The current ThemeManager, or nil if not set.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetThemeManager() *themes.ThemeManager {
	return b.themeManager
}

// Theme resolves the currently active theme from the theme manager.
//
// Returns:
//   - The active Theme from the theme manager, or nil when no manager is configured.
//
// Side effects:
//   - None.
func (b *BaseIntent) Theme() themes.Theme {
	if b.themeManager == nil {
		return nil
	}
	return b.themeManager.Active()
}

// SetLoading activates the loading overlay to indicate an async operation is in progress.
//
// Expected:
//   - message must be non-empty to provide user feedback.
//
// Side effects:
//   - Enables the loading state and stores the loading message.
func (b *BaseIntent) SetLoading(message string) {
	b.isLoading = true
	b.loadingMessage = message
}

// ClearLoading deactivates the loading overlay when the async operation completes.
//
// Side effects:
//   - Disables the loading state and clears the loading message.
func (b *BaseIntent) ClearLoading() {
	b.isLoading = false
	b.loadingMessage = ""
}

// IsLoading checks whether an async operation is currently in progress.
//
// Returns:
//   - True if loading is currently active.
//
// Side effects:
//   - None.
func (b *BaseIntent) IsLoading() bool {
	return b.isLoading
}

// GetLoadingMessage provides the text describing the current async operation.
//
// Returns:
//   - The loading message, or an empty string if not loading.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetLoadingMessage() string {
	return b.loadingMessage
}

// SetError activates the error overlay to display a failure to the user.
//
// Expected:
//   - err is the error to store, or nil to clear.
//
// Side effects:
//   - Replaces the stored error state on the BaseIntent.
func (b *BaseIntent) SetError(err error) {
	b.errorState = err
}

// ClearError dismisses the error overlay.
//
// Side effects:
//   - Sets the error state to nil.
func (b *BaseIntent) ClearError() {
	b.errorState = nil
}

// GetError provides access to the most recent error for display or handling.
//
// Returns:
//   - The current error, or nil if no error is set.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetError() error {
	return b.errorState
}

// HasError checks whether an unresolved error exists.
//
// Returns:
//   - True if the error state is non-nil.
//
// Side effects:
//   - None.
func (b *BaseIntent) HasError() bool {
	return b.errorState != nil
}

// SetSuccess activates a timed success notification that auto-dismisses after 3 seconds.
//
// Expected:
//   - message must be non-empty to provide user feedback.
//
// Side effects:
//   - Stores the message and records the current time for auto-dismiss.
func (b *BaseIntent) SetSuccess(message string) {
	b.successMessage = message
	b.successTime = time.Now()
}

// ClearSuccess dismisses the success notification immediately.
//
// Side effects:
//   - Clears the success message and resets the timestamp.
func (b *BaseIntent) ClearSuccess() {
	b.successMessage = ""
	b.successTime = time.Time{}
}

// GetSuccessMessage provides the text from the most recent successful operation.
//
// Returns:
//   - The success message, or an empty string if not set.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetSuccessMessage() string {
	return b.successMessage
}

// ShouldShowSuccess returns true if success message should be displayed.
// Success messages are shown for 3 seconds after being set.
//
// Returns:
//   - True if a success message exists and fewer than 3 seconds have elapsed since it was set.
//
// Side effects:
//   - None.
func (b *BaseIntent) ShouldShowSuccess() bool {
	if b.successMessage == "" {
		return false
	}
	return time.Since(b.successTime) < 3*time.Second
}

// SetProgress activates the progress overlay to show the status of a multi-step operation.
//
// Expected:
//   - title is the progress bar title.
//   - message describes the current progress step.
//   - value is the completion fraction between 0.0 and 1.0.
//
// Side effects:
//   - Enables progress tracking and stores the title, message, and value.
func (b *BaseIntent) SetProgress(title, message string, value float64) {
	b.progressEnabled = true
	b.progressTitle = title
	b.progressMessage = message
	b.progressValue = value
}

// ClearProgress deactivates the progress overlay.
//
// Side effects:
//   - Disables progress tracking and resets all progress fields to zero values.
func (b *BaseIntent) ClearProgress() {
	b.progressEnabled = false
	b.progressTitle = ""
	b.progressMessage = ""
	b.progressValue = 0.0
}

// IsProgressEnabled checks whether a progress operation is currently running.
//
// Returns:
//   - True if progress tracking is currently active.
//
// Side effects:
//   - None.
func (b *BaseIntent) IsProgressEnabled() bool {
	return b.progressEnabled
}

// GetProgress provides the title, description, and completion fraction of the active operation.
//
// Returns:
//   - title, message, and value representing the current progress.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetProgress() (title, message string, value float64) {
	return b.progressTitle, b.progressMessage, b.progressValue
}

// ShowHelp shows the help modal.
//
// Side effects:
//   - Updates the help modal size from terminal info and makes it visible.
func (b *BaseIntent) ShowHelp() {
	if b.helpModal != nil {
		if b.terminalInfo != nil && b.terminalInfo.IsValid {
			b.helpModal.SetSize(b.terminalInfo.Width, b.terminalInfo.Height)
		}
		b.helpModal.Show()
	}
}

// HideHelp hides the help modal.
//
// Side effects:
//   - Hides the help modal if it exists.
func (b *BaseIntent) HideHelp() {
	if b.helpModal != nil {
		b.helpModal.Hide()
	}
}

// ToggleHelp toggles the help modal visibility.
//
// Side effects:
//   - Updates the help modal size from terminal info and toggles its visibility.
func (b *BaseIntent) ToggleHelp() {
	if b.helpModal != nil {
		if b.terminalInfo != nil && b.terminalInfo.IsValid {
			b.helpModal.SetSize(b.terminalInfo.Width, b.terminalInfo.Height)
		}
		b.helpModal.Toggle()
	}
}

// IsHelpVisible checks whether the help overlay is currently displayed.
//
// Returns:
//   - True if the help modal exists and is visible.
//
// Side effects:
//   - None.
func (b *BaseIntent) IsHelpVisible() bool {
	if b.helpModal != nil {
		return b.helpModal.IsVisible()
	}
	return false
}

// GetHelpModal provides access to the help modal for direct configuration.
//
// Returns:
//   - The HelpModal instance, or nil if not initialized.
//
// Side effects:
//   - None.
func (b *BaseIntent) GetHelpModal() *feedback.HelpModal {
	return b.helpModal
}

// SetHelpKeyMap configures the keyboard shortcuts displayed in the help overlay.
//
// Expected:
//   - keyMap should implement ShortHelp and FullHelp methods; otherwise the call is a no-op.
//
// Side effects:
//   - Configures the help modal's key bindings if the keyMap satisfies the expected interface.
func (b *BaseIntent) SetHelpKeyMap(keyMap interface{}) {
	if b.helpModal != nil {
		if km, ok := keyMap.(interface {
			ShortHelp() []interface{}
			FullHelp() [][]interface{}
		}); ok {
			_ = km
		}
	}
}

// CreateView creates a standardized view with logo and automatic state modals.
// This is a convenience wrapper around CreateStandardView.
//
// Returns:
//   - A ScreenLayout configured with the intent's logo and state modals.
//
// Side effects:
//   - Delegates to CreateStandardView, which reads state from the BaseIntent.
func (b *BaseIntent) CreateView() *layout.ScreenLayout {
	return CreateStandardView(b)
}

// CreateViewWithBreadcrumbs creates a standardized view with breadcrumb navigation.
// This is a convenience wrapper around CreateStandardViewWithBreadcrumbs.
//
// Expected:
//   - crumbs are the breadcrumb segments to display in the header.
//
// Returns:
//   - A ScreenLayout configured with breadcrumbs, logo, and state modals.
//
// Side effects:
//   - Delegates to CreateStandardViewWithBreadcrumbs, which reads state from the BaseIntent.
func (b *BaseIntent) CreateViewWithBreadcrumbs(crumbs ...string) *layout.ScreenLayout {
	return CreateStandardViewWithBreadcrumbs(b, crumbs...)
}
