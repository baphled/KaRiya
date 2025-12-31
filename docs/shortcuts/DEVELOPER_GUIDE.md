# KaRiya Keyboard Shortcut System - Developer Guide

## Architecture Overview

The KaRiya shortcut system consists of four main components:

### 1. ShortcutMapper
Manages global and context-specific shortcuts with conflict detection.

```go
// Create a new mapper
mapper := models.NewShortcutMapper()

// Register global shortcuts
binding := key.NewBinding(
  key.WithKeys("ctrl+s"),
  key.WithHelp("ctrl+s", "save"),
)
mapper.RegisterGlobalShortcut("save", binding, action)

// Register context shortcuts
mapper.RegisterContextShortcut("list", "edit", binding, action)

// Get merged shortcuts (context overrides global)
shortcuts := mapper.GetMergedShortcuts("list")
```

### 2. ContextShortcutHandler
Manages context stacks and context-aware shortcut resolution.

```go
// Create handler
handler := models.NewContextShortcutHandler()
handler.SetMapper(mapper)

// Switch contexts
handler.SetCurrentContext("list")

// Get active shortcuts for current context
activeShortcuts := handler.GetActiveShortcuts()

// Use context stack for modal dialogs
handler.PushContext("modal")
// ... modal operations ...
handler.PopContext() // Back to previous context
```

### 3. ShortcutHelpSystem
Provides help information and discoverability for shortcuts.

```go
// Create help system
help := models.NewShortcutHelpSystem()

// Register shortcuts with help text
help.RegisterShortcut("save", binding, "Save the file", []string{"editor"})

// Categorize for organization
help.CategorizeShortcut("save", "File Operations")

// Get context-specific help
contextHelp := help.GenerateContextHelp("editor")

// Search for shortcuts
results := help.SearchShortcuts("save")
```

### 4. ShortcutCustomizer
Allows users to customize shortcuts and manage profiles.

```go
// Create customizer
customizer := models.NewShortcutCustomizer()

// Override a shortcut
newBinding := key.NewBinding(key.WithKeys("ctrl+shift+s"))
customizer.OverrideShortcut("save", newBinding)

// Create profiles
profile := customizer.CreateProfile("work")
profile.OverrideShortcut("save", binding)

// Apply profile
customizer.ApplyProfile("work")
```

## Integration Guide

### Adding Shortcuts to a Model

1. **Implement the StandardModel interface** in your model
2. **Register shortcuts** in your model's initialization:

```go
type MyModel struct {
  *models.BaseStandardModel
  shortcuts *models.ShortcutMapper
}

func NewMyModel() *MyModel {
  m := &MyModel{
    BaseStandardModel: models.NewBaseStandardModel(),
    shortcuts: models.NewShortcutMapper(),
  }

  // Register shortcuts
  binding := key.NewBinding(key.WithKeys("e"), key.WithHelp("e", "edit"))
  m.shortcuts.RegisterGlobalShortcut("edit", binding, m.handleEdit)

  return m
}
```

3. **Handle shortcuts in Update method**:

```go
func (m *MyModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  switch msg := msg.(type) {
  case tea.KeyMsg:
    cmd, handled := m.shortcuts.HandleMsg(msg)
    if handled {
      return m, cmd
    }
  }
  // ... rest of update logic
  return m, nil
}
```

### Context Management

When navigating between screens, use the ContextShortcutHandler:

```go
// In your app model
func (a *App) SwitchScreen(screenName string) {
  a.contextHandler.SetCurrentContext(screenName)
  a.currentScreen = screenName
}

// For modal dialogs
func (a *App) ShowModal() {
  a.contextHandler.PushContext("modal")
}

func (a *App) CloseModal() {
  a.contextHandler.PopContext()
}
```

## Testing Shortcuts

### Unit Testing

```go
func TestMyShortcuts(t *testing.T) {
  var _ = Describe("MyModel Shortcuts", func() {
    It("should handle edit shortcut", func() {
      model := NewMyModel()
      binding := key.NewBinding(key.WithKeys("e"))

      cmd, exists := model.shortcuts.InvokeShortcut("edit")

      Expect(exists).To(BeTrue())
      Expect(cmd).NotTo(BeNil())
    })
  })
}
```

### Integration Testing

```go
func TestShortcutFlow(t *testing.T) {
  // Test complete shortcut flow through application
  app := NewApp()

  // Switch context
  app.SwitchScreen("list")

  // Verify context shortcuts work
  shortcuts := app.contextHandler.GetActiveShortcuts()
  Expect(shortcuts).To(HaveKey("edit"))
}
```

## Best Practices

### 1. Use StandardModel Interface
All models should implement the StandardModel interface to ensure consistent shortcut handling across the application.

### 2. Document Shortcuts
Always register shortcuts with help text:

```go
// Good - with help text
help.RegisterShortcut("edit", binding, "Edit selected item", []string{"list"})

// Avoid - no help text
mapper.RegisterGlobalShortcut("e", binding, action)
```

### 3. Avoid Conflicts
Check for conflicts before registering:

```go
conflicts := mapper.CheckConflicts("edit", newBinding)
if len(conflicts) > 0 {
  // Handle conflict
  log.Printf("Conflict with: %v", conflicts)
}
```

### 4. Use Context-Sensitive Shortcuts
Only enable shortcuts that make sense in the current context:

```go
// List context - enable editing shortcuts
mapper.RegisterContextShortcut("list", "edit", binding, editAction)

// Form context - don't enable list shortcuts
// mapper.RegisterContextShortcut("form", "edit", ...) - skip this
```

### 5. Provide Help
Always implement help for your shortcuts:

```go
// Register with all details
help.RegisterShortcut(
  "bulk-delete",
  binding,
  "Delete all selected events",
  []string{"list"},
)
help.CategorizeShortcut("bulk-delete", "Bulk Operations")
```

## Common Patterns

### Navigation Shortcuts

```go
// Create consistent navigation across all screens
for _, direction := range []string{"up", "down", "left", "right"} {
  binding := commonShortcuts[direction]
  mapper.RegisterGlobalShortcut(direction, binding,
    getNavigationHandler(direction))
}
```

### Modal Dialogs

```go
// Modal uses context stack
originalContext := handler.GetCurrentContext()
handler.PushContext("confirmation-modal")

// Show modal...

// Close modal and restore context
handler.PopContext()
Expect(handler.GetCurrentContext()).To(Equal(originalContext))
```

### Shortcut Customization

```go
// Load user customizations
customizer := models.NewShortcutCustomizer()
customizer.SetDefaults(defaultShortcuts)

// Apply user profile
if profile := customizer.GetActiveProfile(); profile != "" {
  customizer.ApplyProfile(profile)
}

// Use customized shortcuts
for id, binding := range customizer.GetAllOverrides() {
  mapper.RegisterGlobalShortcut(id, binding, getAction(id))
}
```

## Performance Considerations

### Shortcut Lookup
- **O(1)** for GetGlobalShortcut/GetContextShortcut
- **O(n)** for CheckConflicts where n = number of shortcuts
- **O(1)** for HandleMsg with hashing

### Context Switching
- **O(1)** for SetCurrentContext
- **O(n)** for PushContext (where n = stack depth)
- **O(1)** for PopContext

### Profile Application
- **O(m)** for ApplyProfile where m = overrides in profile

### Optimization Tips
1. Cache frequently accessed shortcuts
2. Lazy-load profiles on demand
3. Pre-compute conflict maps during initialization
4. Use connection pooling for persistence

## Migration Guide

### From Existing System to New System

If migrating from a different shortcut system:

1. **Map old shortcuts to new structure**:
```go
// Old system
oldShortcuts := map[string]string{"save": "ctrl+s"}

// New system
for id, key := range oldShortcuts {
  binding := key.NewBinding(key.WithKeys(key))
  mapper.RegisterGlobalShortcut(id, binding, getAction(id))
}
```

2. **Update models to use StandardModel**:
```go
// Implement StandardModel interface
func (m *MyModel) RegisterShortcuts(shortcuts map[string]key.Binding) {
  for id, binding := range shortcuts {
    m.shortcuts.RegisterGlobalShortcut(id, binding, m.getAction(id))
  }
}
```

3. **Test thoroughly**:
```go
// Ensure all shortcuts work in new system
func TestMigration(t *testing.T) {
  oldApp := setupOldSystem()
  newApp := setupNewSystem()

  compareShortcuts(oldApp, newApp)
}
```

## Troubleshooting

### Shortcuts Not Registering

Check:
1. Model implements StandardModel interface
2. Shortcuts registered in correct context
3. Binding created with both `WithKeys()` and `WithHelp()`

### Conflicts During Initialization

Use conflict detection:
```go
conflicts := mapper.CheckConflicts("my-action", binding)
if len(conflicts) > 0 {
  log.Fatalf("Shortcut conflict detected: %v", conflicts)
}
```

### Context Stack Issues

Verify stack integrity:
```go
depth := handler.GetContextDepth()
stack := handler.GetContextStack()
log.Printf("Context depth: %d, Stack: %v", depth, stack)
```

## Future Enhancements

Planned improvements to the shortcut system:

1. **Key Macro Recording** - Record and replay key sequences
2. **Shortcut Analytics** - Track which shortcuts users use most
3. **Smart Recommendations** - Suggest shortcuts based on usage
4. **Accessibility** - Voice command shortcuts
5. **Internationalization** - Shortcuts for different languages/layouts

---

**Last Updated**: 2025-12-31
**Version**: 1.0
**Maintainer**: KaRiya Team

