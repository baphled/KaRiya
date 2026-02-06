# =============================================================================
# KaRiya Skills Management - Happy Path (Browse and View)
# =============================================================================
#
# WORKFLOW: Browse skills list, navigate, view detail, filter, sort
# DURATION: ~30 seconds
#
# WHAT ARE SKILLS?
# Skills are technical and soft skills extracted from your career events.
# They can be manually added or automatically inferred using AI.
#
# STEPS DEMONSTRATED:
# 1. Navigate to Manage Skills from main menu
# 2. Browse skills list with vim navigation (j/k)
# 3. View skill detail (Enter)
# 4. Use filter modal (f key)
# 5. Use sort modal (s key)
# 6. Return to main menu
#
# KEYBOARD SHORTCUTS:
# - j/k or arrows: Navigate list
# - g/G: Jump to first/last
# - Enter: View skill detail
# - f: Filter skills
# - s: Sort skills
# - n: Add new skill manually
# - i: Infer skills from all events (AI)
# - Escape: Back/cancel
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/skills-happy-path.gif
Output demos/vhs/output/mp4s/skills-happy-path.mp4

# -----------------------------------------------------------------------------
# SETUP: Import sample events (skills are inferred during import)
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
# Import events - skills get created from event data
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
# STEP 2: Navigate to Manage Skills (3rd menu item)
# Main menu order: Capture, Browse, Skills, CV, Configure, Bursts, Facts, Import
# -----------------------------------------------------------------------------
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Browse skills list
# Use vim-style navigation (j/k) to move through list
# -----------------------------------------------------------------------------
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "k"
Sleep 500ms

# Jump to last item
Type "G"
Sleep 500ms

# Jump back to first
Type "g"
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: View skill detail (Enter)
# Shows skill name, description, source, proficiency
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# Close detail modal
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: Open filter modal (f key)
# Filter by skill type, proficiency level, source
# -----------------------------------------------------------------------------
Type "f"
Sleep 1500ms

# Close filter without applying
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 6: Open sort modal (s key)
# Sort by name, proficiency, date added
# -----------------------------------------------------------------------------
Type "s"
Sleep 1500ms

# Close sort without applying
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
