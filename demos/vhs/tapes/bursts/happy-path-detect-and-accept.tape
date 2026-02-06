# =============================================================================
# KaRiya Burst Management - Happy Path (Detect and Accept Bursts)
# =============================================================================
# 
# WORKFLOW: User triggers AI burst detection and accepts suggestions
# DURATION: ~60 seconds
# 
# WHAT IS A BURST?
# A "burst" is a period of focused career activity around a specific theme,
# technology, or project. KaRiya's AI analyzes your events to identify these
# patterns automatically.
#
# STEPS DEMONSTRATED:
# 1. Navigate to Burst Management
# 2. Trigger AI burst detection (s key)
# 3. Review AI suggestions
# 4. Accept some bursts (a key) - facts are auto-extracted
# 5. Reject others (r key)
# 6. View accepted burst details
# 7. View events within the burst
# 8. View extracted facts
#
# AI SUGGESTION REVIEW:
# - a: Accept current suggestion (saves burst, extracts facts)
# - r: Reject current suggestion
# - j/k: Navigate between suggestions
# - Escape: Cancel review
#
# BURST DETAIL ACTIONS:
# - v: View events in this burst
# - f: View facts extracted from burst
# - s: View skills associated
# - c: Confirm burst (triggers fact extraction if not done)
# - i: Infer skills from burst events
# - e: Edit burst name/description
# - d: Delete burst
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/bursts-happy-path-detect.gif
Output demos/vhs/output/mp4s/bursts-happy-path-detect.mp4

# -----------------------------------------------------------------------------
# SETUP: Import sample events to have data for burst detection
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
# Import sample events - AI will detect patterns in these
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
# STEP 2: Navigate to Manage Bursts (6th menu item)
# -----------------------------------------------------------------------------
Down
Down
Down
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Trigger AI burst detection (s key)
# -----------------------------------------------------------------------------
# This analyzes all unassigned events and suggests career bursts
Type "s"

# Wait for AI processing - shows loading modal
Sleep 5s

# Suggestion review modal appears
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 4: Review and accept/reject suggestions
# -----------------------------------------------------------------------------
# Navigate through suggestions
Type "j"
Sleep 500ms

# Accept first suggestion (a key)
# This creates the burst AND automatically extracts facts
Type "a"
Sleep 3s

# Accept another suggestion
Type "a"
Sleep 3s

# Reject one (r key) - we don't want this burst
Type "r"
Sleep 1s

# Finish review - cancel remaining suggestions
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 5: View the accepted burst details
# -----------------------------------------------------------------------------
# Select first burst in list
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 6: View events in this burst (v key)
# -----------------------------------------------------------------------------
Type "v"
Sleep 2s

# Close events modal
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 7: View facts extracted from this burst (f key)
# -----------------------------------------------------------------------------
Type "f"
Sleep 2s

# Close facts modal
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 8: Close detail and return to main menu
# -----------------------------------------------------------------------------
# Close detail modal
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
