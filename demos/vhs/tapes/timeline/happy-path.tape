# =============================================================================
# KaRiya Browse Timeline - Happy Path
# =============================================================================
# 
# WORKFLOW: User browses career events, views details, uses filters
# DURATION: ~45 seconds
# 
# STEPS DEMONSTRATED:
# 1. Navigate to Browse Timeline from main menu
# 2. See list of events (sorted by date, newest first)
# 3. Navigate through events using vim-style keys
# 4. View event detail modal
# 5. Open filter modal and apply a filter
# 6. Open sort modal
# 7. Use search functionality
# 8. Clear filters
#
# TABLE NAVIGATION:
# - j/Down: Move down
# - k/Up: Move up
# - g: Go to first item
# - G: Go to last item
# - PgDn/Ctrl+d: Page down
# - PgUp/Ctrl+u: Page up
#
# ACTIONS:
# - Enter: View event detail
# - f: Open filter modal
# - s: Open sort modal
# - /: Open search modal
# - x: Clear all filters
# - e: Edit selected event
# - d: Delete selected event
# - Escape: Go back / Close modal
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/timeline-happy-path.gif
Output demos/vhs/output/mp4s/timeline-happy-path.mp4

# -----------------------------------------------------------------------------
# SETUP: Import sample events to have data to browse
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
# Import sample events for a populated timeline
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
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
# STEP 3: Navigate through the timeline using vim keys
# -----------------------------------------------------------------------------
# Move down through events
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms

# Move back up
Type "k"
Sleep 200ms
Type "k"
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: View event detail (Enter)
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# Close detail modal (Escape)
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: Open filter modal (f)
# -----------------------------------------------------------------------------
Type "f"
Sleep 1500ms

# Close filter without applying (Escape)
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 6: Open sort modal (s)
# -----------------------------------------------------------------------------
Type "s"
Sleep 1500ms

# Close sort modal
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 7: Open search modal (/)
# -----------------------------------------------------------------------------
Type "/"
Sleep 1s

# Type a search term
Type@80ms "architecture"
Sleep 500ms

# Submit search (if the modal supports it)
Enter
Sleep 1s

# Clear search/filters (x)
Type "x"
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 8: Return to main menu
# -----------------------------------------------------------------------------
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
