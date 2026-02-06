# =============================================================================
# KaRiya CV Generation - Happy Path (Full Wizard)
# =============================================================================
#
# WORKFLOW: Complete CV generation wizard from start to finish
# DURATION: ~60 seconds
#
# CV WIZARD STEPS:
# 1. WHO - Select profile and target audience
# 2. TECH - Select technology focus (if applicable)
# 3. FORMAT - Choose skills format and CV length
# 4. REVIEW - Preview generated CV
# 5. EXPORT - Export to file or clipboard
#
# STEPS DEMONSTRATED:
# 1. Navigate to Generate CV from main menu
# 2. Complete WHO step (profile + audience)
# 3. Complete TECH step (technology focus)
# 4. Complete FORMAT step (skills format + length)
# 5. Review generated CV
# 6. Scroll through preview
# 7. Export to Markdown file
#
# KEYBOARD SHORTCUTS:
# - Enter: Confirm selection / proceed
# - Tab: Move between fields
# - j/k: Navigate options
# - g/G: Jump to top/bottom of preview
# - p: Preview CV
# - x: Export CV
# - Escape: Cancel/back
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/cv-happy-path.gif
Output demos/vhs/output/mp4s/cv-happy-path.mp4

# -----------------------------------------------------------------------------
# SETUP: Import comprehensive events for better CV
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
# Import sample events for rich CV content
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
# STEP 2: Navigate to Generate CV (4th menu item)
# -----------------------------------------------------------------------------
Down
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: WHO step - Profile selection
# Accept default profile
# -----------------------------------------------------------------------------
Enter
Sleep 500ms

# WHO step - Audience selection
# Select "Hiring Manager" (default)
Enter
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: TECH step - Technology focus
# Select "Language Agnostic" or first option
# -----------------------------------------------------------------------------
Sleep 1s
Enter
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: FORMAT step - Skills format
# Select skills display format
# -----------------------------------------------------------------------------
Enter
Sleep 500ms

# FORMAT step - CV length
# Accept default length
Enter
Sleep 500ms

# Confirm final settings
Enter
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 6: Wait for CV generation
# Shows progress indicator during AI processing
# -----------------------------------------------------------------------------
Sleep 4s

# -----------------------------------------------------------------------------
# STEP 7: Review generated CV metadata
# -----------------------------------------------------------------------------
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 8: Preview full CV (p key)
# -----------------------------------------------------------------------------
Type "p"
Sleep 2s

# Scroll through CV content
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 500ms

# Jump to bottom
Type "G"
Sleep 1s

# Jump back to top
Type "g"
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 9: Export CV (x key)
# -----------------------------------------------------------------------------
Type "x"
Sleep 1s

# Select Markdown format
Down
Sleep 200ms
Enter
Sleep 500ms

# Select "Save to file"
Enter
Sleep 2s

# Confirm export
Enter
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 10: Exit cleanly
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
