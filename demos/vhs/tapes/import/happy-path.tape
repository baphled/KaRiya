# =============================================================================
# KaRiya Import Workflow - Happy Path
# =============================================================================
#
# WORKFLOW: Import career events from CSV file
# DURATION: ~40 seconds
#
# IMPORT PROCESS:
# 1. Select CSV file to import
# 2. Preview import data
# 3. Confirm import
# 4. View import results
# 5. Optionally trigger enrichment (bursts, facts, skills)
#
# CSV FORMAT:
# description,date,company,project,tags,categories
#
# STEPS DEMONSTRATED:
# 1. Show sample CSV contents
# 2. Run import command
# 3. View import progress
# 4. See import summary
#
# KEYBOARD SHORTCUTS:
# - Enter: Confirm
# - Escape: Cancel
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/import-happy-path.gif
Output demos/vhs/output/mp4s/import-happy-path.mp4

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
# STEP 1: Show sample CSV structure
# -----------------------------------------------------------------------------
Type "echo '=== Sample Career Events CSV ==='"
Enter
Sleep 500ms

Type "head -8 demos/vhs/sample-events.csv | column -t -s ','"
Enter
Sleep 3s

# Show event count
Type "echo '' && echo 'Total events to import:' $(tail -n +2 demos/vhs/sample-events.csv | wc -l)"
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 2: Run import command
# Shows import progress and results
# -----------------------------------------------------------------------------
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review"
Enter

# Wait for import to complete
Sleep 5s

# -----------------------------------------------------------------------------
# STEP 3: Show success and verify
# -----------------------------------------------------------------------------
Type "echo '' && echo 'Import completed successfully!'"
Enter
Sleep 2s

# Launch app to verify events
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter
Sleep 2s

# Navigate to Browse Timeline to see imported events
Down
Enter
Sleep 2s

# Scroll through events
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
