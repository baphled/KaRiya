# =============================================================================
# KaRiya Fact Management - Create and Delete Facts
# =============================================================================
#
# WORKFLOW: Create a new fact manually, then delete a fact
# DURATION: ~35 seconds
#
# STEPS DEMONSTRATED:
# 1. Navigate to Manage Facts
# 2. Create new fact manually (n key)
# 3. Fill in fact form
# 4. View created fact
# 5. Delete a fact (d key)
# 6. Confirm deletion
#
# FACT FORM FIELDS:
# - Statement: The achievement/metric (required)
# - Context: Additional context/details (optional)
# - Category: Type of fact (Achievement, Metric, Responsibility)
# - Source: Which event/burst this relates to (optional)
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/facts-create-delete.gif
Output demos/vhs/output/mp4s/facts-create-delete.mp4

# -----------------------------------------------------------------------------
# SETUP: Import events for context
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
# STEP 2: Navigate to Manage Facts
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
# STEP 3: Create new fact (n key)
# -----------------------------------------------------------------------------
Type "n"
Sleep 1500ms

# -----------------------------------------------------------------------------
# STEP 4: Fill in fact form
# -----------------------------------------------------------------------------
# Statement (required)
Type "Reduced deployment time from 2 hours to 15 minutes using CI/CD automation"
Tab
Sleep 300ms

# Context (optional)
Type "Implemented GitHub Actions with automated testing and staged deployments"
Tab
Sleep 300ms

# Category - select "Achievement"
Enter
Tab
Sleep 300ms

# -----------------------------------------------------------------------------
# STEP 5: Submit form
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 6: View the created fact
# -----------------------------------------------------------------------------
# Navigate to the new fact (should be at top or bottom)
Type "G"
Sleep 500ms

Enter
Sleep 2s

# Close detail
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 7: Delete a fact (d key)
# -----------------------------------------------------------------------------
Type "d"
Sleep 1s

# Confirm deletion
Enter
Sleep 1s

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
