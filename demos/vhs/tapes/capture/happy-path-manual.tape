# =============================================================================
# KaRiya Capture Event - Happy Path (Manual Mode)
# =============================================================================
# 
# WORKFLOW: User captures a career event using manual mode with all fields
# DURATION: ~45 seconds
# 
# STEPS DEMONSTRATED:
# 1. Select "Capture Event" from main menu
# 2. Choose "Manual" capture strategy (full form)
# 3. Fill in all fields:
#    - Event description (required)
#    - Date
#    - Company OR Project (not both)
#    - Tags (multi-select)
#    - Categories (multi-select)
# 4. Submit form
# 5. See success modal
# 6. Review screen with event details
# 7. Confirm to complete
#
# FORM FIELDS (Manual Mode):
# 1. Event Description (text area, required, 10-2000 chars)
# 2. Date (text input, YYYY-MM-DD or "today" or relative like "1 week ago")
# 3. Company (text input) - OR Project, never both
# 4. Project (text input) - OR Company, never both
# 5. Tags (multi-select, max 6)
# 6. Categories (multi-select, max 6)
# 7. Save Changes (confirm button)
#
# KEYBOARD SHORTCUTS:
# - Tab: Next field
# - Shift+Tab: Previous field
# - Space: Toggle selection in multi-select
# - Enter: Submit form
# - Escape: Cancel / Go back
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/capture-happy-path-manual.gif
Output demos/vhs/output/mp4s/capture-happy-path-manual.mp4

# -----------------------------------------------------------------------------
# SETUP: Create environment with profile (skip onboarding)
# -----------------------------------------------------------------------------
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

# -----------------------------------------------------------------------------
# STEP 1: Launch application and see main menu
# -----------------------------------------------------------------------------
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 2: Select "Capture Event" (first menu item)
# -----------------------------------------------------------------------------
# Capture Event is the first item, so just press Enter
Enter
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 3: Choose capture strategy - select "Manual" (second option)
# -----------------------------------------------------------------------------
# Strategy screen shows: Quick (first) and Manual (second)
Down
Enter
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 4: Fill in the event form
# -----------------------------------------------------------------------------
# Field 1: Event Description (text area - required)
# Write a detailed achievement with metrics
Type@50ms "Architected and implemented a distributed caching solution using Redis clusters, reducing database load by 75% and improving API response times from 500ms to under 50ms across 12 microservices"
Tab
Sleep 500ms

# Field 2: Date (text input)
Type@80ms "2025-01-15"
Tab
Sleep 500ms

# Field 3: Company (text input) - using company for this event
Type@80ms "TechCorp Inc"
Tab
Sleep 500ms

# Field 4: Project (text input) - SKIP because we used Company
# Events should have Company OR Project, not both
Tab
Sleep 500ms

# Field 5: Tags (multi-select)
# Navigate and select relevant tags with Space key
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

# Field 6: Categories (multi-select)
# Select one or two categories
Down
Space
Sleep 200ms
Tab
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: Submit the form
# -----------------------------------------------------------------------------
# Press Enter on "Save Changes" button
Enter

# Wait for save modal ("Saving event...")
Sleep 1s

# Wait for success modal ("Event saved!")
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 6: Review screen appears
# -----------------------------------------------------------------------------
# Shows: Event details, Inferred Bursts, Inferred Facts
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 7: Confirm on review screen to complete the workflow
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# Back at main menu - capture complete!
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP: Exit application
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
