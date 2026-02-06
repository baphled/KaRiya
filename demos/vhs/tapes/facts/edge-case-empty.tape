# =============================================================================
# KaRiya Fact Management - Edge Case: Empty State
# =============================================================================
#
# WORKFLOW: Show empty facts screen behavior
# DURATION: ~15 seconds
#
# WHAT HAPPENS WITH NO FACTS:
# When no facts exist, the screen shows an empty state with hints
# on how to create facts (from bursts, events, or manually).
#
# STEPS DEMONSTRATED:
# 1. Navigate to Manage Facts with empty database
# 2. Observe empty state message
# 3. Show help modal for available actions
# 4. Return to main menu
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/facts-empty.gif
Output demos/vhs/output/mp4s/facts-empty.mp4

# -----------------------------------------------------------------------------
# SETUP: Fresh empty database
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
# STEP 2: Navigate to Manage Facts
# -----------------------------------------------------------------------------
Down
Down
Down
Down
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Observe empty state
# -----------------------------------------------------------------------------
Sleep 2s

# Show help for available actions
Type "?"
Sleep 2s

# Close help
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: Return to main menu
# -----------------------------------------------------------------------------
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
