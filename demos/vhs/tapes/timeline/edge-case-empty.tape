# =============================================================================
# KaRiya Browse Timeline - Edge Case (Empty Timeline)
# =============================================================================
# 
# WORKFLOW: User opens browse timeline with no events
# DURATION: ~20 seconds
# 
# STEPS DEMONSTRATED:
# 1. Open Browse Timeline with empty database
# 2. See empty state message
# 3. Return to main menu
#
# EDGE CASE:
# When there are no events, the timeline shows a helpful message
# guiding the user to capture their first event.
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/timeline-edge-case-empty.gif
Output demos/vhs/output/mp4s/timeline-edge-case-empty.mp4

# -----------------------------------------------------------------------------
# SETUP: Create environment with NO events (empty database)
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
# STEP 2: Navigate to Browse Timeline (second menu item)
# -----------------------------------------------------------------------------
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Observe empty state
# -----------------------------------------------------------------------------
# The screen shows "No events found" or similar message
# Navigation keys have no effect (nothing to navigate)
Type "j"
Sleep 300ms
Type "k"
Sleep 300ms

# Trying to view detail does nothing (no events)
Enter
Sleep 500ms

# Show that we can still access filter/sort (though they won't help)
Type "f"
Sleep 1s
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: Return to main menu
# -----------------------------------------------------------------------------
Escape
Sleep 1s

# The user would typically go to Capture Event next
Sleep 500ms

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
