# BDD Helper Layer - Page Object Pattern

## Purpose

This helper layer abstracts UI navigation from BDD step definitions, making tests **maintainable** and **resilient to UI changes**.

**Problem**: Direct navigation calls (`env.Tab()`, `env.Confirm()`) in step definitions mean UI changes break multiple test files.

**Solution**: Centralize navigation in helper classes. Change navigation logic in **ONE place**.

## Available Helpers

### FormHelper

Form field navigation and interaction.

```go
form, _ := helpers.NewFormHelper(ctx)

// Navigate to field
form.NavigateToField("Log Level")

// Set field value
form.SetFieldValue("Data Directory", "/custom/path")

// Select from dropdown
form.SelectFieldOption("Log Level", "debug")

// Toggle boolean
form.ToggleBoolean("Auto Backup")

// Submit/Cancel
form.SubmitForm()  // Ctrl+S
form.CancelForm()  // Escape
```

### ModalHelper

Modal interaction and verification.

```go
modal, _ := helpers.NewModalHelper(ctx)

// Open modal
modal.OpenModalWithKey('x')

// Wait for modal
modal.WaitForModalToAppear("Export Options")

// Interact
modal.ConfirmModal()
modal.CloseModal()

// Verify
modal.IsModalVisible("CV Review")
modal.AssertModalContains("success")
```

### ListHelper

List navigation and selection.

```go
list, _ := helpers.NewListHelper(ctx)

// Navigate to item
list.NavigateToItem("Generate CV")

// Select item
list.SelectItemByName("System Settings")

// Check visibility
list.IsItemVisible("Manage Facts")
```

## Usage in Step Definitions

### Before (Brittle)

```go
func iNavigateToField(ctx, label) {
    env := GetEnv(ctx)
    for i := 0; i < 20; i++ {
        if strings.Contains(env.GetView(), label) {
            return nil
        }
        env.Tab()  // Direct call - breaks if navigation changes
    }
}
```

### After (Maintainable)

```go
func iNavigateToField(ctx, label) {
    form, err := helpers.NewFormHelper(ctx)
    if err != nil {
        return err
    }
    return form.NavigateToField(label)  // Single place to update
}
```

## Benefits

1. **Single Point of Change**: UI navigation logic in ONE place
2. **Consistent Behavior**: All steps use same navigation
3. **Easy Testing**: Helpers can be unit tested
4. **Clear Intent**: `form.SetFieldValue()` vs raw key presses
5. **Reduced Duplication**: Reuse helpers across step files

## When UI Changes

**Example**: Navigation changes from Tab to Arrow keys

**Without helpers**: Update 50+ step definitions across 10 files  
**With helpers**: Update `FormHelper.NavigateToField()` only

## Adding New Helpers

1. Create helper in `features/support/helpers/`
2. Follow existing patterns (FormHelper, ModalHelper)
3. Document public methods
4. Update this README

## Guidelines

- **Keep helpers UI-agnostic**: No business logic
- **One responsibility per helper**: Forms, Modals, Lists
- **Clear method names**: `NavigateToField` not `DoStuff`
- **Error handling**: Return errors, let steps decide handling
- **Documentation**: Comment all public methods
