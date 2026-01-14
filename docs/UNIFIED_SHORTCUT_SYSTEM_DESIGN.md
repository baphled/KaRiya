---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Unified Shortcut System Design

**Last Updated**: 2026-01-03
**Current State**: 5 separate files with overlapping responsibilities
**Target State**: 1 unified, cohesive shortcut system
**Implementation Status**: Design phase (ready for implementation)

---

## Executive Summary

The current shortcut system is spread across 5 files with overlapping concerns:

1. **shortcut_handler.go** (125 lines) - Basic shortcut registration
2. **shortcut_mapper.go** (255 lines) - Global + context shortcuts
3. **context_shortcut_handler.go** (180 lines) - Context stack management
4. **shortcut_customizer.go** (261 lines) - UI for customization
5. **shortcut_help_system.go** (275 lines) - Help display

**Problem**: Confusing API with duplicate functionality across files

**Solution**: Consolidate into unified system with clear layers:
- **Core API** - Single ShortcutManager (replaces handler + mapper)
- **Context Management** - Integrated context stack
- **Customization** - ShortcutCustomizer (unchanged)
- **Help** - ShortcutHelp (unchanged)

**Result**: 5 files → 4 files, clearer API, reduced confusion

---

## Current System Analysis

### Current Architecture

```
┌─────────────────────────────────────────────────────┐
│            Shortcut System (5 files)                │
├─────────────────────────────────────────────────────┤
│                                                     │
│  ShortcutHandler ────┐                             │
│  (basic registration)│                             │
│                      ├─→ ShortcutMapper            │
│  ContextShortcutHandler │  (global + context)      │
│  (context stack)     │                             │
│                      └─→ ShortcutCustomizer        │
│                           (UI for customization)   │
│                                                     │
│                      ┌─→ ShortcutHelpSystem        │
│                      │   (help display)            │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### Problems

1. **Duplicate Responsibility**: Both Handler and Mapper handle shortcuts
2. **Confusing API**: Multiple entry points for same operation
3. **Context Complexity**: ContextShortcutHandler separate from Mapper
4. **File Count**: 5 files for one feature
5. **Testing Burden**: Multiple test files with overlapping coverage

### Current Usage Pattern

```go
// Create handler
handler := NewShortcutHandler()

// Also create mapper?
mapper := NewShortcutMapper()

// Also create context handler?
contextHandler := NewContextShortcutHandler()
contextHandler.SetMapper(mapper)

// Register shortcuts in multiple places
handler.RegisterShortcut("quit", keyBinding, quitAction)
mapper.RegisterGlobalShortcut("quit", keyBinding, quitAction)
mapper.RegisterContextShortcut("detail", "edit", keyBinding, editAction)

// Get shortcuts from multiple places
shortcuts := handler.GetAllShortcuts()
globalShortcuts := mapper.GetAllGlobalShortcuts()
contextShortcuts := mapper.GetContextShortcuts("detail")

// This is confusing!
```

---

## Proposed Unified System

### New Architecture

```
┌──────────────────────────────────────────────────┐
│         Unified Shortcut System (4 files)        │
├──────────────────────────────────────────────────┤
│                                                  │
│  ShortcutManager                                │
│  ├─ Global shortcuts registry                  │
│  ├─ Context shortcuts registry                 │
│  ├─ Context stack management                   │
│  └─ Unified API for all operations             │
│                                                  │
│  ShortcutCustomizer                            │
│  └─ UI for customization (unchanged)           │
│                                                  │
│  ShortcutHelp                                  │
│  └─ Help display (unchanged)                   │
│                                                  │
│  shortcut_utils.go (NEW)                       │
│  └─ Helper functions, key bindings             │
│                                                  │
└──────────────────────────────────────────────────┘
```

### New Usage Pattern

```go
// Single entry point
manager := NewShortcutManager()

// Register global shortcuts
manager.RegisterGlobal("quit", keyBinding, quitAction)
manager.RegisterGlobal("help", keyBinding, helpAction)

// Register context shortcuts
manager.RegisterContext("detail", "edit", keyBinding, editAction)
manager.RegisterContext("detail", "delete", keyBinding, deleteAction)

// Push/pop context
manager.PushContext("detail")
shortcuts := manager.GetActiveShortcuts() // Merges global + context
manager.PopContext()

// Handle messages
cmd, handled := manager.HandleMessage(msg)

// Get help
help := manager.GetContextHelp("detail")
```

---

## Component Design

### 1. ShortcutManager (New - Core)

**Purpose**: Single unified API for all shortcut operations

**Responsibilities**:
- Register/unregister global shortcuts
- Register/unregister context shortcuts
- Manage context stack
- Handle key messages
- Provide conflict detection
- Merge shortcuts (global + active context)

**API**:

```go
type ShortcutManager struct {
    globalShortcuts   map[string]*ShortcutBinding
    contextShortcuts  map[string]map[string]*ShortcutBinding
    contextStack      []string
    contextMetadata   map[string]map[string]interface{}
}

// Global shortcuts
func (sm *ShortcutManager) RegisterGlobal(id string, binding key.Binding, action ShortcutAction) error
func (sm *ShortcutManager) UnregisterGlobal(id string) error
func (sm *ShortcutManager) GetGlobal(id string) (*ShortcutBinding, bool)
func (sm *ShortcutManager) GetAllGlobal() map[string]*ShortcutBinding

// Context shortcuts
func (sm *ShortcutManager) RegisterContext(context, id string, binding key.Binding, action ShortcutAction) error
func (sm *ShortcutManager) UnregisterContext(context, id string) error
func (sm *ShortcutManager) GetContext(context, id string) (*ShortcutBinding, bool)
func (sm *ShortcutManager) GetAllContext(context string) map[string]*ShortcutBinding

// Context management
func (sm *ShortcutManager) PushContext(context string)
func (sm *ShortcutManager) PopContext() string
func (sm *ShortcutManager) GetContextStack() []string
func (sm *ShortcutManager) SetCurrentContext(context string)
func (sm *ShortcutManager) GetCurrentContext() string

// Context metadata
func (sm *ShortcutManager) SetMetadata(context, key string, value interface{})
func (sm *ShortcutManager) GetMetadata(context, key string) (interface{}, bool)
func (sm *ShortcutManager) ClearMetadata(context string)

// Active shortcuts (merged)
func (sm *ShortcutManager) GetActiveShortcuts() map[string]*ShortcutBinding
func (sm *ShortcutManager) LookupActive(id string) (*ShortcutBinding, bool)

// Message handling
func (sm *ShortcutManager) HandleMessage(msg tea.Msg) (tea.Cmd, bool)

// Utilities
func (sm *ShortcutManager) CheckConflicts(id string, binding key.Binding) []string
func (sm *ShortcutManager) Clear()
func (sm *ShortcutManager) ClearGlobal()
func (sm *ShortcutManager) ClearContext(context string)
```

**Key Features**:
- Single responsibility: shortcut management
- Clear context hierarchy
- Unified message handling
- Conflict detection
- Metadata storage per context

---

### 2. ShortcutCustomizer (Existing - UI)

**Status**: Keep as-is (no changes)

**Purpose**: UI for customizing shortcuts

**Key Methods**:
- Init(), Update(), View()
- GetCustomShortcuts()
- ApplyCustomizations()

**No changes needed** - this is already well-designed for its purpose

---

### 3. ShortcutHelp (Existing → Simplified)

**Current File**: shortcut_help_system.go

**Changes**:
- Simplify to work with ShortcutManager instead of separate Mapper
- Remove duplicate shortcut retrieval logic
- Consolidate help rendering

**API**:

```go
type ShortcutHelp struct {
    manager *ShortcutManager
    width   int
}

func NewShortcutHelp(manager *ShortcutManager) *ShortcutHelp
func (sh *ShortcutHelp) GetGlobalHelp() string
func (sh *ShortcutHelp) GetContextHelp(context string) string
func (sh *ShortcutHelp) GetAllHelp() string
func (sh *ShortcutHelp) SetWidth(width int)
```

---

### 4. Shortcut Utils (New)

**Purpose**: Helper functions and constants

**Contents**:
- Key binding definitions (CommonShortcuts)
- Key string parsing
- Shortcut validation
- Conflict resolution utilities

**API**:

```go
type ShortcutBinding struct {
    ID       string
    Keys     []key.Binding
    Action   ShortcutAction
    Help     string
    Context  string // "" for global
}

type CommonShortcuts struct {
    Quit   key.Binding
    Help   key.Binding
    Up     key.Binding
    Down   key.Binding
    Enter  key.Binding
    Escape key.Binding
    Back   key.Binding
    Edit   key.Binding
    Delete key.Binding
}

func NewCommonShortcuts() CommonShortcuts
func GetKeyString(binding key.Binding) string
func ParseKeyString(str string) (key.Binding, error)
func ValidateBinding(binding key.Binding) error
```

---

## Migration Path

### Step 1: Create ShortcutManager

1. Create `shortcuts.go` with ShortcutManager
2. Implement all methods
3. Write comprehensive tests
4. Keep existing files unchanged

### Step 2: Create Adapter Layer

1. Create `shortcut_compat.go` with compatibility functions
2. Map old API calls to new ShortcutManager
3. Ensure no breaking changes

### Step 3: Update Customizer

1. Update ShortcutCustomizer to use ShortcutManager
2. Test thoroughly
3. Keep UI unchanged

### Step 4: Update Help System

1. Simplify ShortcutHelpSystem
2. Use ShortcutManager instead of Mapper
3. Remove duplicate logic

### Step 5: Extract Utils

1. Create `shortcut_utils.go`
2. Move helper functions
3. Remove from other files

### Step 6: Deprecate Old Files

1. Mark old files as deprecated
2. Add migration comments
3. Document transition path

### Step 7: Remove Old Files (Future)

1. After all code migrated
2. Keep compatibility layer for 1-2 releases
3. Full removal in major version bump

---

## Data Structures

### ShortcutManager Internal Structure

```go
type ShortcutManager struct {
    // Global shortcuts: id → binding + action
    globalShortcuts map[string]*ShortcutBinding

    // Context shortcuts: context → (id → binding + action)
    contextShortcuts map[string]map[string]*ShortcutBinding

    // Context stack for nested contexts
    contextStack []string

    // Metadata per context: context → (key → value)
    contextMetadata map[string]map[string]interface{}

    // Cached merged shortcuts for active context
    mergedCache map[string]*ShortcutBinding
    cacheDirty  bool
}

type ShortcutBinding struct {
    ID       string
    Keys     []key.Binding
    Action   ShortcutAction
    Help     string
    Context  string // "" for global, "detail" for context-specific
}
```

---

## API Comparison

### Before (Current - Confusing)

```go
// Multiple entry points
handler := NewShortcutHandler()
mapper := NewShortcutMapper()
contextHandler := NewContextShortcutHandler()

// Register in different places
handler.RegisterShortcut("quit", binding, action)
mapper.RegisterGlobalShortcut("quit", binding, action)
mapper.RegisterContextShortcut("detail", "edit", binding, action)

// Get from different places
h := handler.GetAllShortcuts()
g := mapper.GetAllGlobalShortcuts()
c := mapper.GetContextShortcuts("detail")
```

### After (New - Clear)

```go
// Single entry point
manager := NewShortcutManager()

// Register in one place
manager.RegisterGlobal("quit", binding, action)
manager.RegisterContext("detail", "edit", binding, action)

// Get from one place
active := manager.GetActiveShortcuts()
global := manager.GetAllGlobal()
context := manager.GetAllContext("detail")
```

---

## Benefits

### Clarity
- ✅ Single ShortcutManager handles all operations
- ✅ Clear separation: core vs UI vs utils
- ✅ No duplicate APIs

### Maintainability
- ✅ Fewer files (5 → 4)
- ✅ Clearer responsibilities
- ✅ Easier to test
- ✅ Simpler documentation

### Performance
- ✅ Single registry (no duplication)
- ✅ Efficient context merging
- ✅ Cached merged shortcuts

### Extensibility
- ✅ Easy to add features (e.g., conditional shortcuts)
- ✅ Plugin-friendly API
- ✅ Metadata storage per context

### Backward Compatibility
- ✅ Compatibility layer for old API
- ✅ Gradual migration path
- ✅ No breaking changes immediately

---

## Testing Strategy

### Unit Tests

1. **ShortcutManager**
   - Register/unregister global shortcuts
   - Register/unregister context shortcuts
   - Context stack operations
   - Shortcut lookup (global, context, active)
   - Conflict detection
   - Message handling

2. **ShortcutHelp**
   - Generate help text
   - Format context help
   - Width handling

3. **Utils**
   - Key binding parsing
   - Validation
   - Conflict resolution

### Integration Tests

1. Full workflow (register → push context → lookup → pop)
2. Multiple contexts
3. Metadata management
4. Customizer integration
5. Help system integration

### Coverage Target

- **Core**: 95%+ (ShortcutManager)
- **Overall**: 90%+

---

## Implementation Checklist

- [ ] Design review & approval
- [ ] Create ShortcutManager
- [ ] Write comprehensive tests
- [ ] Create compatibility layer
- [ ] Update ShortcutCustomizer
- [ ] Simplify ShortcutHelpSystem
- [ ] Extract utils
- [ ] Update documentation
- [ ] Deprecate old files
- [ ] Mark for future removal

---

## Timeline

- **Design**: ✅ Complete
- **Implementation**: 4-6 hours
- **Testing**: 2-3 hours
- **Documentation**: 1-2 hours
- **Total**: 7-11 hours

---

## Success Criteria

- ✅ ShortcutManager passes all tests
- ✅ Backward compatibility maintained
- ✅ All existing shortcuts work
- ✅ No performance regression
- ✅ Code coverage ≥90%
- ✅ File count reduced (5 → 4)
- ✅ API clarity improved
- ✅ Documentation updated

---

## Future Enhancements

Once unified system is in place:

1. **Conditional Shortcuts**: Based on state/context
2. **Macro Support**: Chain multiple shortcuts
3. **Shortcut Profiles**: Different profiles per user
4. **Conflict Resolution**: Automatic rebinding
5. **Plugin System**: Dynamic shortcut registration
6. **Analytics**: Track shortcut usage

---

**Status**: 📋 Design Complete - Ready for Implementation
**Generated**: 2026-01-03
**Next Step**: Implement ShortcutManager and run tests

