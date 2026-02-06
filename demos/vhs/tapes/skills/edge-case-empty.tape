# =============================================================================
# KaRiya Skills Management - Edge Case: Empty State
# =============================================================================
#
# WORKFLOW: Show empty skills screen behavior
# DURATION: ~15 seconds
#
# WHAT HAPPENS WITH NO SKILLS:
# When no skills exist, the screen shows an empty state message
# with hints on how to add skills (manual or AI inference).
#
# STEPS DEMONSTRATED:
# 1. Navigate to Manage Skills with empty database
# 2. Observe empty state message
# 3. Show available actions (n for new, i for infer)
# 4. Return to main menu
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/skills-empty.gif
Output demos/vhs/output/mp4s/skills-empty.mp4

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
# STEP 2: Navigate to Manage Skills
# -----------------------------------------------------------------------------
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Observe empty state
# The screen shows "No skills found" with action hints
# -----------------------------------------------------------------------------
Sleep 2s

# Show help modal to see available actions
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
