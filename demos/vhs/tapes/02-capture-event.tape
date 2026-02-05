# KaRiya Capture Event Demo
# Duration: ~30 seconds
# Shows capturing a new career event using Manual mode with proper save flow
#
# This demo shows how to capture a new career event, demonstrating:
# - Strategy selection (Quick vs Manual)
# - Form fields and navigation
# - Multi-select fields for tags/categories
# - Submit and review flow
# - Final confirmation
#
# Form fields (Manual mode):
# 1. Event Description (text area, required)
# 2. Date (text input, optional)
# 3. Company (text input, optional) - OR Project, not both
# 4. Project (text input, optional) - OR Company, not both
# 5. Tags (multi-select, use Space to select)
# 6. Categories (multi-select, use Space to select)
# 7. Save Changes (confirm button)

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/02-capture-event.gif
Output demos/vhs/output/mp4s/02-capture-event.mp4

# Setup with existing config (hidden)
Hide
Type `rm -rf /tmp/kariya-demo && mkdir -p /tmp/kariya-demo`
Enter
Sleep 300ms
Type `cp demos/vhs/minimal-config.yaml /tmp/kariya-demo/config.yaml`
Enter
Sleep 300ms
Type "cd /home/baphled/Projects/KoRiya"
Enter
Sleep 300ms
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Select Capture Event (first option)
Enter

# Wait for strategy selection
Sleep 1s

# Select Manual mode (second option) for full form
Down
Enter

# Wait for event form
Sleep 1s

# Field 1: Event Description (text area)
Type@60ms "Designed and implemented distributed caching layer reducing database load by 75% and improving API response times from 500ms to 50ms"
Tab
Sleep 500ms

# Field 2: Date (text input)
Type@80ms "2025-04-15"
Tab
Sleep 500ms

# Field 3: Company (text input) - using company for this event
Type@80ms "TechCorp"
Tab
Sleep 500ms

# Field 4: Project (text input) - skip, already have company
Tab
Sleep 500ms

# Field 5: Tags (multi-select) - navigate and select with Space
Down
Space
Sleep 200ms
Down
Space
Sleep 200ms
Down
Space
Sleep 200ms
Tab
Sleep 500ms

# Field 6: Categories (multi-select) - select one category
Down
Space
Sleep 200ms
Tab
Sleep 500ms

# Field 7: Submit - press Enter on Save Changes button
Enter

# Wait for save modal ("Saving event...")
Sleep 1s

# Success modal appears ("Event saved!")
Sleep 2s

# Review screen appears - shows event details
Sleep 1s

# Confirm on review screen to complete
Enter

# Wait for intent to complete and return to menu
Sleep 2s

# Show we're back at main menu
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
