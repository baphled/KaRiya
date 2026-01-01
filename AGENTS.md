# KaRiya Project Handover Document

## Project Overview
KaRiya is a Career Journal CLI tool designed to help professionals track, manage, and reflect on their career events and progression. It provides an interactive terminal interface for capturing, organizing, and analyzing career milestones.

## Technical Specifications

### Technology Stack
- **Language**: Go (1.24+)
- **CLI Framework**: BubbleTea (Charmbracelet)
- **Testing**: Ginkgo v2
- **Database**: SQLite (modernc.org/sqlite)
- **Version Control**: Semantic Release
- **Commit Management**: Conventional Commits, Commitlint

### Key Dependencies
- github.com/charmbracelet/bubbles
- github.com/charmbracelet/bubbletea
- github.com/onsi/ginkgo/v2
- modernc.org/sqlite

## Development Workflow

### Prerequisites
- Go 1.24 or higher
- Node.js 18+ with npm
- Ginkgo v2 for testing
- Make (for task automation)

### Setup
1. Clone the repository
2. Run `go mod tidy` to install Go dependencies
3. Run `npm install` for Node.js dependencies
4. Run `make install-git-hooks` to setup git hooks

### Key Make Commands
- `make test`: Run all tests
- `make coverage`: Generate code coverage report
- `make install-git-hooks`: Setup git hooks
- `make check-ai-attribution`: Verify AI commit attribution

## Project Structure

### Main Directories
- `cmd/cli/`: CLI entry point and main application
- `internal/cli/`: Core CLI implementation
  - `app/`: Application state and navigation
  - `models/`: Screen models (BubbleTea)
  - `components/`: Reusable UI components
  - `styles/`: Styling and layout
  - `validation/`: Input validation
  - `service/`: Service layer adapters

### Key Configuration Files
- `go.mod`: Go module dependencies
- `package.json`: Node.js dependencies and scripts
- `.commitlintrc.json`: Commit message validation
- `.releaserc.json`: Semantic Release configuration
- `Makefile`: Development task automation

## Development Guidelines

### Commit Message Convention
Use conventional commits format:
```
<type>(<scope>): <subject>

<body>

<footer>
```

### AI Commit Attribution
- All AI-generated code must include attribution
- Format:
  ```
  AI-Generated-By: <Assistant Name> (<Model Version>)
  Reviewed-By: <Your Name>
  ```

### Testing
- 131+ tests across various components
- 100% passing test suite
- Use Ginkgo for testing
- Aim for comprehensive test coverage

## Deployment & Release
- Automated releases via GitHub Actions
- Semantic versioning
- Automatic CHANGELOG generation
- Binaries uploaded with each release

## Troubleshooting
- Refer to README.md for detailed troubleshooting
- Common issues include:
  - Database persistence
  - Terminal compatibility
  - Input navigation

---

## KeyMsg and Navigation Standardization (Latest Update)

### Overview
A standardized key handling system has been implemented to provide a centralized point of truth for keyboard navigation across the entire KaRiya application. This eliminates redundant key handling code and ensures consistent behavior across all views.

### Architecture

#### New Navigation Package Modules

1. **`internal/cli/navigation/key_handler.go`** - Core Key Handler Framework
   - `KeyAction` struct: Represents the result of handling a key press
   - `KeyHandler` interface: Standard interface for handling KeyMsgs
   - Five handler implementations:
     - `ListKeyHandler`: For list-based views (supports up/down, pagination, search, filter, sort, edit, delete)
     - `FormKeyHandler`: For form input views (supports field navigation, form submission)
     - `DialogKeyHandler`: For modal/dialog views (supports yes/no/ok/cancel confirmation)
     - `ViewKeyHandler`: For read-only display views (supports scrolling, navigation)
     - `MenuKeyHandler`: For menu selections (supports up/down navigation, selection)

2. **`internal/cli/navigation/handlers.go`** - Helper Functions and Global Handlers
   - `KeyHandlers` factory struct with methods to create appropriate handlers
   - `GlobalKeyHandler`: Handles application-wide keys (quit, help, home)
   - Helper functions:
     - `IsGlobalKey()`: Check if a key is a global key
     - `IsNavigationKey()`: Check if a key matches a navigation constant
     - `GetNavigationKeyFromString()`: Map key string to NavigationKey constant
     - `MatchesNavigationKey()`: Check if KeyMsg matches a NavigationKey

#### Existing Navigation Constants (Reused)
All handlers map to existing `NavigationKey` constants in `internal/cli/navigation/constants.go`:
- `KeyUp`, `KeyDown`, `KeyLeft`, `KeyRight`: Arrow keys + vim-style keys (k/j/h/l)
- `KeySelect`: Enter key for confirmations
- `KeyToggle`: Space for toggles
- `KeyBack`: Esc for going back
- `KeyFilter`, `KeySort`, `KeySearch`: Action keys
- `KeyEdit`, `KeyDelete`: Item operations
- `KeyQuit`, `KeyHelp`, `KeyHome`: Global actions
- `KeyCapture`, `KeyList`, `KeyBulk`, `KeyMetadata`, `KeyPending`: App navigation

### Implementation Pattern

#### For List Views (e.g., list.go, burst_list.go)
```go
// In struct definition:
type ListModel struct {
    // ... existing fields ...
    keyHandler  navigation.KeyHandler
}

// In constructor:
keyHandler: navigation.NewListKeyHandler()

// In Update() method:
case tea.KeyMsg:
    action := m.keyHandler.HandleKey(msg)
    if !action.IsHandled {
        return m, nil
    }
    switch {
    case action.NavigationKey != nil:
        return m.handleNavigationKey(*action.NavigationKey)
    case action.ActionType != "":
        return m.handleActionType(action.ActionType)
    }

// Add handler methods:
func (m *ListModel) handleNavigationKey(key navigation.NavigationKey) (tea.Model, tea.Cmd) {
    switch key {
    case navigation.KeyUp:
        m.prevItem()
    case navigation.KeyDown:
        m.nextItem()
    case navigation.KeySelect:
        return m, m.viewSelectedEvent()
    // ... handle other navigation keys ...
    }
    return m, nil
}

func (m *ListModel) handleActionType(actionType string) (tea.Model, tea.Cmd) {
    switch actionType {
    case "navigate:first":
        m.goToFirstItem()
    case "navigate:last":
        m.goToLastItem()
    case "navigate:page_up":
        m.prevPage()
    case "navigate:page_down":
        m.nextPage()
    // ... handle other action types ...
    }
    return m, nil
}
```

#### For Form Views
```go
// Same pattern as above, but:
keyHandler: navigation.NewFormKeyHandler()

// Handle form-specific action types like:
// - "field:next": Move to next form field
// - "field:previous": Move to previous form field
// - "form:submit": Submit the form
// - "form:cancel": Cancel the form
```

#### For Dialog Views
```go
// Same pattern as above, but:
keyHandler: navigation.NewDialogKeyHandler()

// Handle dialog-specific action types like:
// - "dialog:yes": Confirm/Yes button
// - "dialog:no": Cancel/No button
// - "dialog:ok": OK confirmation
// - "dialog:cancel": Cancel action
```

#### For View/Display Screens
```go
// Same pattern as above, but:
keyHandler: navigation.NewViewKeyHandler()

// Handle view-specific action types like:
// - "view:scroll_up": Scroll content up
// - "view:scroll_down": Scroll content down
// - "view:page_up": Page up
// - "view:page_down": Page down
// - "view:first": Jump to start
// - "view:last": Jump to end
```

### Files Refactored So Far

1. **menu.go** ✅ - Completed
   - Uses `MenuKeyHandler`
   - Handles up/down navigation and selection
   
2. **list.go** ✅ - Completed
   - Uses `ListKeyHandler`
   - Handles pagination, navigation, search, filter, sort
   
3. **confirmation_dialog.go** ✅ - Completed
   - Uses `DialogKeyHandler`
   - Handles yes/no confirmation with left/right switching

### Files Pending Refactoring

**High Priority (Core Functionality):**
- form.go - Complex form with field navigation
- help.go - Help screen with section navigation
- details.go - Details view with scrolling
- tutorial.go - Interactive tutorial

**Medium Priority (List-Based Views):**
- burst_list.go - Burst events list
- fact_list.go - Facts list
- facts_results.go - Fact search results
- import_review.go - CSV import review
- metadata_review.go - Metadata review workflow
- view_event.go - Event details view
- view_event_with_facts.go - Event with facts display
- burst_suggestion.go - Burst suggestions
- action_menu.go - Action selection menu

**Lower Priority (Specialized):**
- bulk_operations.go - Bulk operations
- search.go - Search interface
- metadata_editor.go - Metadata editor
- fact_editor.go - Fact editor

### Key Benefits

1. **Centralized Key Handling**: Single point of truth for all keyboard shortcuts
2. **Consistency**: All views behave consistently for common keys (up/down, esc, q)
3. **Maintainability**: Changes to key bindings only need to be made in one place
4. **Testability**: Key handlers can be tested independently of view logic
5. **Reusability**: Different view types use the same patterns
6. **Extensibility**: New action types can be added without changing the core framework
7. **Vim Key Support**: All list and form views support vim-style keys (j/k/h/l)
8. **Existing Constants**: Leverages existing `NavigationKey` constants for consistency

### Recommended Refactoring Order

1. Simple dialogs first (confirmation_dialog ✅)
2. Simple view screens (help, tutorial)
3. List-based views (fact_list, burst_list, etc.)
4. Complex forms (form, metadata_editor)
5. Specialized screens (editors, menus)

### Testing Strategy

- All key handler implementations should be tested in isolation
- Views should be tested to ensure they correctly handle KeyActions
- Integration tests should verify key mappings work end-to-end
- Existing test suite (131+ tests) should continue to pass

### Migration Path

Old Code Pattern:
```go
case tea.KeyMsg:
    switch msg.String() {
    case "up", "k":
        m.prevItem()
    case "down", "j":
        m.nextItem()
    // ... repeated patterns across 30+ files ...
    }
```

New Code Pattern:
```go
case tea.KeyMsg:
    action := m.keyHandler.HandleKey(msg)
    if !action.IsHandled {
        return m, nil
    }
    switch {
    case action.NavigationKey != nil:
        return m.handleNavigationKey(*action.NavigationKey)
    case action.ActionType != "":
        return m.handleActionType(action.ActionType)
    }
```

### Next Steps for Future Development

1. Complete refactoring of remaining models using the established pattern
2. Add tests for key handler implementations
3. Consider adding configuration-based key binding customization
4. Add keyboard shortcut customization UI (may already exist in shortcut_customizer.go)
5. Document all action types in a central location for developer reference

