# =============================================================================
# KaRiya Configure System - Sad Path (Cancel Editing)
# =============================================================================
#
# WORKFLOW: Start editing, then cancel without saving
# DURATION: ~20 seconds
#
# STEPS DEMONSTRATED:
# 1. Navigate to Configure System
# 2. Select a domain to edit
# 3. Make some changes in the form
# 4. Cancel without saving (Escape)
# 5. Verify original values preserved
#
# CANCEL BEHAVIOR:
# - Pressing Escape cancels the edit
# - No changes are saved
# - Returns to domain list
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/config-sad-path-cancel.gif
Output demos/vhs/output/mp4s/config-sad-path-cancel.mp4

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
# STEP 2: Navigate to Configure System
# -----------------------------------------------------------------------------
Down
Down
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Select Profile domain to edit
# -----------------------------------------------------------------------------
Enter
Sleep 1500ms

# -----------------------------------------------------------------------------
# STEP 4: Make changes (that we'll cancel)
# -----------------------------------------------------------------------------
Ctrl+A
Type "THIS WILL BE CANCELLED"
Tab
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: Cancel without saving (Escape)
# -----------------------------------------------------------------------------
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 6: Verify original values preserved
# Re-open the domain to see original values
# -----------------------------------------------------------------------------
Enter
Sleep 1500ms

# Original value should still be there (not "THIS WILL BE CANCELLED")
Sleep 1s

# Close without changes
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 7: Return to main menu
# -----------------------------------------------------------------------------
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
