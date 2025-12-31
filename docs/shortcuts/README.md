# KaRiya Keyboard Shortcuts Documentation

## Overview

KaRiya provides a comprehensive keyboard shortcut system designed to improve user productivity and accessibility. This document covers all available shortcuts, how to customize them, and best practices for using the application.

## Table of Contents

1. [Default Shortcuts](#default-shortcuts)
2. [Context-Specific Shortcuts](#context-specific-shortcuts)
3. [Customizing Shortcuts](#customizing-shortcuts)
4. [Shortcut Profiles](#shortcut-profiles)
5. [Discovering Shortcuts](#discovering-shortcuts)
6. [Best Practices](#best-practices)
7. [Troubleshooting](#troubleshooting)

## Default Shortcuts

### Navigation Shortcuts

| Key | Action | Context | Description |
|-----|--------|---------|-------------|
| `Esc` | Back | All | Return to previous screen |
| `↑` / `k` | Up | Lists, Forms | Navigate up |
| `↓` / `j` | Down | Lists, Forms | Navigate down |
| `←` / `h` | Left | Forms, Menus | Navigate left |
| `→` / `l` | Right | Forms, Menus | Navigate right |
| `PgUp` | Page Up | Lists | Scroll up one page |
| `PgDn` | Page Down | Lists | Scroll down one page |

### Edit Shortcuts

| Key | Action | Context | Description |
|-----|--------|---------|-------------|
| `e` | Edit | Lists | Edit selected event |
| `d` | Delete | Lists | Delete selected event |
| `c` | Capture | All | Create new career event |
| `Enter` | Confirm | Forms | Submit form or confirm action |
| `Space` | Toggle | Forms, Lists | Toggle checkbox or expand item |

### View & Filter Shortcuts

| Key | Action | Context | Description |
|-----|--------|---------|-------------|
| `/` | Search | Lists | Open search/filter |
| `f` | Filter | Lists | Show filter options |
| `s` | Sort | Lists | Show sort options |
| `m` | Metadata | All | Open metadata review |
| `b` | Bulk | Lists | Enter bulk operations mode |

### System Shortcuts

| Key | Action | Context | Description |
|-----|--------|---------|-------------|
| `?` / `h` | Help | All | Display help information |
| `q` | Quit | All | Exit application |
| `ctrl+c` | Quit | All | Force quit |

## Context-Specific Shortcuts

Shortcuts may behave differently depending on your current screen or context. The system automatically manages shortcut availability based on where you are in the application.

### List Context Shortcuts
Available when viewing lists of events:
- `e` - Edit selected event
- `d` - Delete selected event
- `Space` - Select/deselect event
- `Enter` - View event details

### Form Context Shortcuts
Available when editing or creating events:
- `Tab` / `Shift+Tab` - Move between form fields
- `Enter` - Submit form
- `Esc` - Cancel and return to list

### Search Context Shortcuts
Available when searching or filtering:
- `Enter` - Apply search/filter
- `Esc` - Cancel search
- `↑`/`↓` - Navigate search results

## Customizing Shortcuts

### View Current Customizations

To see your custom shortcuts and disabled keys:
```
Press ? to open help → View custom shortcuts
```

### Change a Shortcut

1. Open the settings menu
2. Select "Keyboard Shortcuts"
3. Find the shortcut you want to change
4. Press the new key combination
5. Confirm the change

### Disable a Shortcut

1. Open the settings menu
2. Select "Keyboard Shortcuts"
3. Find the shortcut to disable
4. Press spacebar to toggle it off
5. Press Enter to confirm

### Reset to Defaults

To reset all shortcuts to their default configuration:
1. Open the settings menu
2. Select "Keyboard Shortcuts"
3. Select "Reset to Defaults"
4. Confirm the action

## Shortcut Profiles

Shortcut profiles allow you to save and switch between different keyboard configurations.

### Create a Profile

1. Open the settings menu
2. Select "Shortcut Profiles"
3. Select "New Profile"
4. Enter a name (e.g., "Vim-style", "Work", "Gaming")
5. Configure shortcuts as desired
6. Press Enter to save

### Apply a Profile

1. Open the settings menu
2. Select "Shortcut Profiles"
3. Select the profile you want to use
4. Press Enter to apply
5. The new profile will be active immediately

### Common Profiles

#### Vim-Style Profile
Optimized for users familiar with Vim editor:
```
Navigation:
  h = left
  j = down
  k = up
  l = right
  G = end of list
  g = start of list

Editing:
  d = delete
  y = copy
  p = paste
```

#### Emacs-Style Profile
Optimized for Emacs users:
```
Navigation:
  ctrl+n = down
  ctrl+p = up
  ctrl+f = right
  ctrl+b = left
  ctrl+a = start of line
  ctrl+e = end of line
```

## Discovering Shortcuts

### Built-in Help System

Press `?` at any time to see available shortcuts for your current context:
- Shortcuts are grouped by category
- Descriptions explain what each shortcut does
- Context availability is clearly marked

### Search Shortcuts

You can search for specific shortcuts by:
1. Pressing `?` to open help
2. Typing a keyword (e.g., "delete", "ctrl+s")
3. Available shortcuts matching your search will be highlighted

### View All Shortcuts

To see all available shortcuts in the application:
1. Press `?` to open help
2. Select "All Shortcuts"
3. Browse through all available shortcuts
4. Use categories to find related shortcuts

## Best Practices

### 1. Learn Context-Aware Shortcuts First
Start by learning the most common shortcuts for the screens you use most frequently. The application will remind you of available shortcuts on each screen.

### 2. Avoid Conflicts
When customizing shortcuts, avoid:
- Shortcuts already used by your operating system
- Similar shortcuts that are easy to confuse
- Too many modifier keys (Ctrl+Shift+Alt+...) which are hard to remember

### 3. Use Profiles for Different Workflows
Create different profiles for:
- **Daily work** - Optimized for event capture and review
- **Data entry** - Optimized for form navigation
- **Browsing** - Optimized for list navigation and searching

### 4. Keep Related Shortcuts Together
If you customize shortcuts, try to keep related actions close together:
- Place delete next to select
- Place edit next to view
- Keep navigation keys consistent

### 5. Document Your Customizations
If you create custom profiles:
- Give them descriptive names
- Document your reason for the changes
- Share useful profiles with team members

## Troubleshooting

### Shortcut Not Working

**Problem**: A shortcut key doesn't seem to work
**Solution**:
1. Verify the context - the shortcut may only work in certain screens
2. Check if the shortcut is disabled in your custom settings
3. Press `?` to see available shortcuts for your current context
4. Ensure no other application is intercepting the key

### Conflicting Shortcuts

**Problem**: Two shortcuts use the same key
**Solution**:
1. Open settings → Keyboard Shortcuts
2. Find the conflicting shortcuts
3. Reassign one to a different key
4. Consider creating a profile for specific workflows

### Forgotten Shortcuts

**Problem**: Can't remember a shortcut
**Solution**:
1. Press `?` at any time to open help
2. Use the search feature to find the action you want
3. Read the shortcut definition
4. Create a custom profile with more intuitive shortcuts for you

### Reset Not Working

**Problem**: Reset to defaults didn't work
**Solution**:
1. Go to your KaRiya config directory
2. Find the shortcuts configuration file
3. Delete it and restart the application
4. Shortcuts will be restored to factory defaults

## Configuration Files

Shortcut customizations are stored in:
- **Linux/Mac**: `~/.config/kariya/shortcuts.json`
- **Windows**: `%APPDATA%\KaRiya\shortcuts.json`

You can:
- Manually edit the JSON file for advanced customization
- Share configuration files between computers
- Backup your customizations regularly

## Advanced Topics

### Creating Custom Shortcut Profiles via Configuration

You can create shortcut profiles programmatically:

```json
{
  "profiles": {
    "vim-style": {
      "name": "Vim Style",
      "shortcuts": {
        "navigate_up": "k",
        "navigate_down": "j",
        "navigate_left": "h",
        "navigate_right": "l"
      },
      "disabled": ["ctrl+s"]
    }
  },
  "active_profile": "vim-style"
}
```

### Keyboard Layout Considerations

If you use a non-QWERTY keyboard layout:
1. KaRiya detects your keyboard layout automatically
2. Shortcuts work based on key position, not character
3. You can override this in settings if needed

## Support

For issues with shortcuts or to request new customization features:
- Visit the [KaRiya GitHub Issues](https://github.com/baphled/kariya/issues)
- Include your operating system and keyboard layout
- Describe the shortcut problem with steps to reproduce

---

**Last Updated**: 2025-12-31
**Version**: 1.0

