# =============================================================================
# KaRiya Capture Event - Sad Path (Cancel at Various Stages)
# =============================================================================
# 
# WORKFLOW: User cancels event capture at different stages
# DURATION: ~35 seconds
# 
# STEPS DEMONSTRATED:
# 1. Start capture, cancel at strategy selection
# 2. Start capture, begin filling form, cancel from form
# 3. Shows that no data is saved when cancelled
#
# CANCELLATION POINTS:
# - Strategy selection screen (Escape → Main menu)
# - Form screen (Escape → Back to strategy)
# - From strategy again (Escape → Main menu)
#
# KEY CONCEPT:
# Escape always takes you back one step without saving.
# No data is persisted until you complete the full workflow.
#
# KEYBOARD SHORTCUTS:
# - Escape: Cancel / Go back (no save)
# - Ctrl+C: Force quit application
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/capture-sad-path-cancel.gif
Output demos/vhs/output/mp4s/capture-sad-path-cancel.mp4

# -----------------------------------------------------------------------------
# SETUP
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
# STEP 1: Launch application
# -----------------------------------------------------------------------------
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 2: Enter Capture Event, then cancel at strategy selection
# -----------------------------------------------------------------------------
# Select Capture Event
Enter
Sleep 1s

# We're at strategy selection - press Escape to cancel
# This returns us to main menu without any action
Escape
Sleep 1s

# Back at main menu - notice no event was created
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 3: Enter Capture Event again, go to form, then cancel
# -----------------------------------------------------------------------------
# Select Capture Event again
Enter
Sleep 1s

# Choose Manual mode to see the full form
Down
Enter
Sleep 1s

# Start typing an event description
Type@50ms "Started working on a new feature but changed my mind about recording this..."
Sleep 1s

# Cancel from the form with Escape
# This goes back to strategy selection, NOT main menu
Escape
Sleep 1s

# We're back at strategy selection
# Cancel again to go to main menu
Escape
Sleep 1s

# Back at main menu - still no event was saved
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 4: Quick demonstration that we can start fresh
# -----------------------------------------------------------------------------
# Select Capture Event one more time
Enter
Sleep 1s

# Notice the form is empty - previous input was not saved
# Choose Quick mode this time
Enter
Sleep 1s

# Type briefly to show fresh state
Type@50ms "Fresh start - nothing was saved from before"
Sleep 1s

# Cancel out completely
Escape
Sleep 1s
Escape
Sleep 1s

# Back at main menu
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
