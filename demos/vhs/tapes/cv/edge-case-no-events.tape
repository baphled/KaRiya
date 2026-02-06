# =============================================================================
# KaRiya CV Generation - Edge Case: No Events
# =============================================================================
#
# WORKFLOW: Attempt CV generation with empty event database
# DURATION: ~15 seconds
#
# WHAT HAPPENS:
# When no career events exist, the CV generator shows a warning
# that there's insufficient data to generate a meaningful CV.
# It suggests capturing events first.
#
# STEPS DEMONSTRATED:
# 1. Navigate to Generate CV with empty database
# 2. Observe warning message
# 3. Return to main menu (to capture events)
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/cv-no-events.gif
Output demos/vhs/output/mp4s/cv-no-events.mp4

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
# STEP 2: Navigate to Generate CV
# -----------------------------------------------------------------------------
Down
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Observe warning about no events
# The system shows a message explaining CV generation requires events
# -----------------------------------------------------------------------------
Sleep 3s

# -----------------------------------------------------------------------------
# STEP 4: Return to main menu
# Following the suggestion to capture events first
# -----------------------------------------------------------------------------
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
