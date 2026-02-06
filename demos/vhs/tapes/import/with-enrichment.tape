# =============================================================================
# KaRiya Import Workflow - With Full Enrichment
# =============================================================================
#
# WORKFLOW: Import events and trigger AI enrichment (bursts, facts, skills)
# DURATION: ~60 seconds
#
# ENRICHMENT PIPELINE:
# 1. Import CSV events
# 2. AI detects career bursts from event patterns
# 3. AI extracts facts from events and bursts
# 4. AI infers skills from all content
#
# STEPS DEMONSTRATED:
# 1. Import sample events
# 2. Launch app and browse timeline
# 3. Navigate to bursts - see AI-detected bursts
# 4. Navigate to facts - see extracted facts
# 5. Navigate to skills - see inferred skills
#
# This demonstrates the full data enrichment pipeline
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/import-with-enrichment.gif
Output demos/vhs/output/mp4s/import-with-enrichment.mp4

# -----------------------------------------------------------------------------
# SETUP: Run full import with enrichment
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
# STEP 1: Run import with full enrichment
# This triggers the complete AI pipeline
# -----------------------------------------------------------------------------
Type "echo 'Starting import with full AI enrichment...'"
Enter
Sleep 1s

Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv"
Enter

# Wait for import and enrichment (this takes longer)
Sleep 8s

# -----------------------------------------------------------------------------
# STEP 2: Launch app to see results
# -----------------------------------------------------------------------------
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Browse Timeline - see imported events
# -----------------------------------------------------------------------------
Down
Enter
Sleep 2s

Type "j"
Sleep 200ms
Type "j"
Sleep 500ms

# Return to main menu
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 4: Check Bursts - see AI-detected bursts
# -----------------------------------------------------------------------------
Down
Down
Down
Down
Down
Enter
Sleep 2s

Type "j"
Sleep 500ms

# Return to main menu
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 5: Check Facts - see extracted facts
# -----------------------------------------------------------------------------
Down
Down
Down
Down
Down
Down
Enter
Sleep 2s

Type "j"
Sleep 200ms
Type "j"
Sleep 500ms

# Return to main menu
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# STEP 6: Check Skills - see inferred skills
# -----------------------------------------------------------------------------
Down
Down
Enter
Sleep 2s

Type "j"
Sleep 200ms
Type "j"
Sleep 500ms

# Return to main menu
Escape
Sleep 500ms

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
