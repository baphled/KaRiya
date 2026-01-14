---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# KaRiya TUI - Navigation Guide

**Version**: 1.0  
**Last Updated**: 2026-01-06

---

## Quick Reference Card

### Universal Navigation Keys

| Key | Action | Works From |
|-----|--------|------------|
| **Esc** | Go back one step | All states |
| **m** | Return to main menu | All states |
| **q** | Quit application | All states |
| **Ctrl+C** | Force quit | Anywhere |

**Remember**: You can ALWAYS escape or return to the main menu!

---

## Navigation Patterns by Scenario

### Scenario 1: Browsing Career Events

```
Main Menu
  ↓ (Select "Browse Timeline")
Browse Timeline
  ↓ (Select an event)
Event Details
  
Navigation:
  • Esc → Back to timeline
  • m   → Back to main menu
  • q   → Quit application
```

**Example Flow**:
1. Arrow keys to select event
2. Enter to view details
3. **Esc** to go back to timeline
4. **m** to return to main menu instantly

---

### Scenario 2: Capturing a Career Event

```
Main Menu
  ↓ (Select "Capture Event")
Choose Strategy
  ↓ (Select strategy)
Fill Form
  ↓ (Fill in details)
Review Event
  ↓ (Confirm)
Submit
  
Navigation at each step:
  • Esc → Go back one step
  • m   → Cancel and return to main menu
  • q   → Quit application
```

**Example Flow**:
1. Choose "Quick Capture"
2. Start typing event details
3. **Esc** → Back to strategy selection
4. Choose different strategy
5. Fill form again
6. **m** → Changed mind, return to main menu

**Error Handling**:
- If submit fails, **error stays visible**
- Press **Esc** to go back and fix
- Error message preserved for reference

---

### Scenario 3: Generating a CV

```
Main Menu
  ↓ (Select "Generate CV")
Select Profile
  ↓ (Choose profile)
Select Audience
  ↓ (Choose audience)
⏳ Generating CV...  ← Async operation
Preview CV
  ↓ (Review)
Confirm
  ↓ (Confirm or export)
Export (optional)
  
Special Navigation:
  • During "Generating CV":
    - Esc → CV continues in background
    - m   → Cancel immediately
    - q   → Quit (cancels generation)
```

**Example Flow - Background Generation**:
1. Select profile: "Senior Engineer"
2. Select audience: "Tech Startup"
3. Press Enter to generate
4. See "⏳ Generating CV..."
5. **Press Esc** → CV generation continues in background
6. Navigate to main menu
7. Come back later to see generated CV

**Example Flow - Wait and Review**:
1. Select profile and audience
2. Wait for CV generation
3. Preview shows: "✅ Generated successfully"
4. Review bullets and sections
5. Press **e** to export or **Esc** to go back
6. Press **m** if you want to cancel entirely

---

### Scenario 4: Configuring System Settings

```
Main Menu
  ↓ (Select "Configure System")
Select Domain
  ↓ (e.g., "Database")
Edit Settings
  ↓ (Change settings)
Review Changes
  ↓ (Confirm)
Confirm
  ↓ (Yes)
⏳ Saving...  ← Async operation
Complete
  
Special Navigation:
  • During "Saving":
    - Esc → Save continues in background
    - m   → Cancel immediately
    - q   → Quit (cancels save)
```

**Example Flow - Multi-Domain Configuration**:
1. Select "Database" domain
2. Change connection string
3. Review changes
4. **Esc** → Go back to edit more settings
5. Adjust timeout value
6. Review again
7. Confirm
8. **Esc during saving** → Continues in background
9. **m** → Return to main menu while saving

---

### Scenario 5: Exporting Artifacts

```
Main Menu
  ↓ (Select "Export Artifact")
Select Type
  ↓ (e.g., "CV")
Select Format
  ↓ (e.g., "PDF")
Select Destination
  ↓ (e.g., "File")
Configure
  ↓ (Set file path)
Preview
  ↓ (Review)
Confirm
  ↓ (Yes)
⏳ Exporting...  ← Async operation
Complete
  
Special Navigation:
  • During "Exporting":
    - Esc → Export continues in background
    - m   → Cancel immediately
    - q   → Quit (cancels export)
```

**Example Flow - Quick Export**:
1. Select "Events" artifact
2. Select "JSON" format
3. Select "File" destination
4. **Esc** → Back to format selection
5. Change to "CSV" instead
6. Select destination again
7. Confirm
8. **m during export** → Cancel and return to menu

---

## Special Cases

### Async Operations (Background Work)

When you see these messages:
- "⏳ Generating CV..."
- "⏳ Exporting artifact..."
- "⏳ Saving configuration..."

**You have 3 options**:

| Key | Result |
|-----|--------|
| **Wait** | Operation completes, you see results |
| **Esc** | Operation continues in **background**, you can navigate away |
| **m** | Operation is **cancelled immediately**, return to main menu |
| **q** | Application quits, operation is **cancelled** |

**Pro Tip**: Use **Esc** for long operations (CV generation) so you can do other things while it completes!

---

### Error States

If you see an error (red text):
```
❌ Error: Failed to save configuration
Database connection timeout

Esc: Back (error visible) | m: Main menu | q: Quit
```

**The error message stays visible when you press Esc!**

This lets you:
1. Read the full error message
2. Press **Esc** to go back
3. Fix the issue
4. Try again with the error still visible for reference

---

## Keyboard Shortcuts Summary

### In Any List/Table View
| Key | Action |
|-----|--------|
| ↑/k | Move up |
| ↓/j | Move down |
| PgUp/b | Page up |
| PgDn/f | Page down |
| Home/g | Go to first item |
| End/G | Go to last item |
| Enter | Select/confirm |
| Esc | Go back |
| m | Main menu |
| q | Quit |

### In Any Form/Input
| Key | Action |
|-----|--------|
| Tab | Next field |
| Shift+Tab | Previous field |
| Enter | Submit/next |
| Esc | Cancel/go back |
| m | Main menu |
| q | Quit |

### In Any Preview/Review
| Key | Action |
|-----|--------|
| ↑/↓ | Scroll |
| PgUp/PgDn | Page scroll |
| e | Edit |
| c | Confirm |
| Esc | Go back |
| m | Main menu |
| q | Quit |

---

## Common Navigation Workflows

### "I made a mistake, go back!"
1. Press **Esc** once → Previous screen
2. Press **Esc** again → One more step back
3. Or press **m** → Instantly to main menu

### "Cancel everything and start over"
1. Press **m** from anywhere → Main menu
2. Start fresh

### "I want to quit"
1. Press **q** from anywhere → Application quits
2. Or **Ctrl+C** → Force quit

### "Long operation, come back later"
1. Start operation (e.g., Generate CV)
2. Press **Esc** → Operation continues in background
3. Navigate to main menu
4. Do other things
5. Come back later to see results

### "Oops, I don't want to wait"
1. Operation started (e.g., Exporting)
2. Press **m** → Cancel immediately
3. Back to main menu

---

## Tips and Tricks

### 💡 Tip 1: Escape is Your Friend
Never feel stuck! **Esc** always gets you out.

### 💡 Tip 2: Main Menu is One Key Away
Lost in a deep workflow? Press **m** to reset.

### 💡 Tip 3: Background Operations
For long CV generations:
- Start generation
- Press **Esc** (continues in background)
- Browse timeline or edit events
- Come back later

### 💡 Tip 4: Error Recovery
Errors don't disappear:
- Read the error
- Press **Esc**
- Fix the issue
- Error still visible for reference

### 💡 Tip 5: Quick Exit Levels
- **Esc** → Back one step
- **Esc + Esc** → Back two steps
- **m** → Main menu instantly
- **q** → Quit application

---

## Troubleshooting

### Q: I pressed a key and nothing happened
**A**: You might be in a state that's processing. Look for:
- "⏳ Generating..."
- "⏳ Saving..."
- "⏳ Exporting..."

Wait for completion or press **m** to cancel.

### Q: Can I cancel during save/export?
**A**: Yes!
- **Esc** → Continues in background (safe)
- **m** → Cancels immediately
- **q** → Quits (cancels operation)

### Q: I want to see the error again
**A**: Errors stay visible when you press **Esc** to go back. They only disappear when you:
- Move forward to next state
- Press **m** to go to main menu
- Press **q** to quit

### Q: How do I know what keys work in current screen?
**A**: Look at the **footer** at the bottom of every screen:
```
Enter: Confirm | Esc: Back | m: Main menu | q: Quit
```

The footer always shows available keys!

### Q: What if I'm stuck?
**A**: You can ALWAYS:
1. Press **Esc** to go back
2. Press **m** to go to main menu
3. Press **q** to quit
4. Press **Ctrl+C** to force quit

**You can never get truly stuck!**

---

## Advanced Usage

### Multi-Step Workflows with Quick Escapes

**Scenario**: Generate multiple CVs for different audiences

```
1. Main Menu → Generate CV
2. Select Profile: "Senior Engineer"
3. Select Audience: "Tech Startup"
4. Wait for generation
5. Press 'm' → Main menu (skip preview)
6. Generate CV again
7. Same profile: "Senior Engineer"
8. Different audience: "Enterprise"
9. Repeat
```

### Background CV Generation

**Scenario**: Generate CV while browsing events

```
1. Generate CV → Select profile → Select audience
2. See "⏳ Generating CV..."
3. Press Esc → CV continues in background
4. Main Menu → Browse Timeline
5. Review recent events
6. Main Menu → Check CV (if done)
```

---

## Summary: The Golden Rules

1. ✅ **Esc always goes back** (one step)
2. ✅ **'m' always goes to main menu** (any depth)
3. ✅ **'q' always quits** (immediate)
4. ✅ **Async operations can continue in background** (Esc during "⏳")
5. ✅ **Errors stay visible when going back** (Esc preserves errors)
6. ✅ **Footer shows available keys** (always read the footer)
7. ✅ **You can never get stuck** (multiple escape options)

---

## Need Help?

- **In-app help**: Press **?** (if available in current screen)
- **View this guide**: `docs/USER_GUIDE_NAVIGATION.md`
- **Report issues**: Check `docs/ESCAPE_KEY_STANDARDIZATION_COMPLETE.md`

**Happy navigating! 🚀**
