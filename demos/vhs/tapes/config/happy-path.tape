# =============================================================================
# KaRiya Configure System - Happy Path
# =============================================================================
#
# WORKFLOW: Browse configuration domains, edit settings, save
# DURATION: ~35 seconds
#
# CONFIGURATION DOMAINS:
# - Profile: User's personal information (name, email, title)
# - AI Settings: LLM model selection, API keys
# - Export: CV export preferences, formats
# - Privacy: Data retention, analytics
#
# STEPS DEMONSTRATED:
# 1. Navigate to Configure System from main menu
# 2. Browse configuration domains
# 3. Select a domain to edit
# 4. Modify settings in form
# 5. Save changes
# 6. Return to main menu
#
# KEYBOARD SHORTCUTS:
# - j/k or arrows: Navigate domains
# - Enter: Edit selected domain
# - Escape: Cancel/back
# - Tab: Move between form fields
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/config-happy-path.gif
Output demos/vhs/output/mp4s/config-happy-path.mp4

# -----------------------------------------------------------------------------
# SETUP: Use minimal config as starting point
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
# STEP 2: Navigate to Configure System (5th menu item)
# Main menu: Capture, Browse, Skills, CV, Configure, Bursts, Facts, Import
# -----------------------------------------------------------------------------
Down
Down
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Browse configuration domains
# -----------------------------------------------------------------------------
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "k"
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: Select Profile domain to edit
# -----------------------------------------------------------------------------
Enter
Sleep 1500ms

# -----------------------------------------------------------------------------
# STEP 5: Modify settings in form
# Navigate through form fields and make changes
# -----------------------------------------------------------------------------
# Name field - clear and update
Ctrl+A
Type "Jane Developer"
Tab
Sleep 300ms

# Email field
Ctrl+A
Type "jane@example.com"
Tab
Sleep 300ms

# Title field
Ctrl+A
Type "Senior Software Engineer"
Tab
Sleep 300ms

# -----------------------------------------------------------------------------
# STEP 6: Save changes (Submit form)
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 7: Verify changes saved (view domain again)
# -----------------------------------------------------------------------------
Enter
Sleep 1500ms

# Confirm changes are persisted
Escape
Sleep 500ms

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
