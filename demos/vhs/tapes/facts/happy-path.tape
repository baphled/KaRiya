# =============================================================================
# KaRiya Fact Management - Happy Path (Browse and View)
# =============================================================================
#
# WORKFLOW: Browse facts list, view detail, edit fact
# DURATION: ~30 seconds
#
# WHAT ARE FACTS?
# Facts are key achievements, metrics, and accomplishments extracted from
# career events and bursts. They become bullet points in your CV.
# Examples: "Reduced API latency by 40%", "Led team of 5 engineers"
#
# STEPS DEMONSTRATED:
# 1. Navigate to Manage Facts from main menu
# 2. Browse facts list with vim navigation
# 3. View fact detail
# 4. Open edit mode
# 5. Return to main menu
#
# KEYBOARD SHORTCUTS:
# - j/k or arrows: Navigate list
# - g/G: Jump to first/last
# - Enter: View fact detail
# - e: Edit fact
# - d: Delete fact
# - n: Add new fact manually
# - Escape: Back/cancel
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/facts-happy-path.gif
Output demos/vhs/output/mp4s/facts-happy-path.mp4

# -----------------------------------------------------------------------------
# SETUP: Import events and extract facts
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
# Import events - facts get extracted from event descriptions
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
# STEP 2: Navigate to Manage Facts (7th menu item)
# Main menu: Capture, Browse, Skills, CV, Configure, Bursts, Facts, Import
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
# STEP 3: Browse facts list
# -----------------------------------------------------------------------------
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "k"
Sleep 500ms

# Jump to last
Type "G"
Sleep 500ms

# Jump to first
Type "g"
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: View fact detail
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# Close detail
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: Open edit mode (e key)
# -----------------------------------------------------------------------------
Type "e"
Sleep 1500ms

# Cancel edit without changes
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 6: Return to main menu
# -----------------------------------------------------------------------------
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
