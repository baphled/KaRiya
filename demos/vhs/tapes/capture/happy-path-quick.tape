# =============================================================================
# KaRiya Capture Event - Happy Path (Quick Mode)
# =============================================================================
# 
# WORKFLOW: User captures a career event using quick mode (minimal fields)
# DURATION: ~30 seconds
# 
# STEPS DEMONSTRATED:
# 1. Select "Capture Event" from main menu
# 2. Choose "Quick" capture strategy (minimal form)
# 3. Fill in only required fields:
#    - Event description (required)
#    - Date (optional, defaults to today)
# 4. Submit form
# 5. See success modal
# 6. Confirm on review screen
#
# QUICK MODE is ideal for:
# - Rapid event capture
# - When you want to add details later
# - Brief notes about achievements
#
# KEYBOARD SHORTCUTS:
# - Tab: Next field
# - Enter: Submit form
# - Escape: Cancel / Go back
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/capture-happy-path-quick.gif
Output demos/vhs/output/mp4s/capture-happy-path-quick.mp4

# -----------------------------------------------------------------------------
# SETUP: Create environment with profile (skip onboarding)
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
# STEP 2: Select "Capture Event" (first menu item)
# -----------------------------------------------------------------------------
Enter
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 3: Choose "Quick" capture strategy (first option - already selected)
# -----------------------------------------------------------------------------
# Quick mode is the first option, just press Enter
Enter
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 4: Fill in the quick capture form
# -----------------------------------------------------------------------------
# Field 1: Event Description (text area - required)
Type@50ms "Led team migration from monolith to microservices architecture, improving deployment frequency by 300%"
Tab
Sleep 500ms

# Field 2: Date (text input - optional, using relative date)
Type@80ms "2 weeks ago"
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: Submit the form
# -----------------------------------------------------------------------------
Enter

# Wait for save and success modals
Sleep 3s

# -----------------------------------------------------------------------------
# STEP 6: Review screen - confirm to complete
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# Back at main menu
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP: Exit application
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
