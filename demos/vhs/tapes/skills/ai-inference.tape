# =============================================================================
# KaRiya Skills Management - AI Skill Inference
# =============================================================================
#
# WORKFLOW: Use AI to automatically infer skills from career events
# DURATION: ~45 seconds
#
# HOW AI INFERENCE WORKS:
# KaRiya analyzes all your career events and extracts skills mentioned
# in descriptions, technologies used, and achievements. The AI suggests
# skills with confidence scores.
#
# STEPS DEMONSTRATED:
# 1. Navigate to Manage Skills (empty initially)
# 2. Trigger AI skill inference (i key)
# 3. Wait for AI processing
# 4. Review suggested skills
# 5. Accept/reject suggestions
# 6. View accepted skills in list
#
# AI SUGGESTION REVIEW:
# - a: Accept current suggestion
# - r: Reject current suggestion
# - j/k: Navigate between suggestions
# - A: Accept all remaining
# - Escape: Cancel (reject remaining)
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/skills-ai-inference.gif
Output demos/vhs/output/mp4s/skills-ai-inference.mp4

# -----------------------------------------------------------------------------
# SETUP: Import events WITHOUT skill inference (for demo purposes)
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
# Import events - skills will be inferred manually in the demo
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review --skip-skill-inference > /dev/null 2>&1`
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
# STEP 2: Navigate to Manage Skills
# -----------------------------------------------------------------------------
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Trigger AI skill inference (i key)
# This analyzes ALL events and suggests skills
# -----------------------------------------------------------------------------
Type "i"

# Wait for AI processing (shows progress indicator)
Sleep 5s

# -----------------------------------------------------------------------------
# STEP 4: Review suggested skills
# AI presents each suggestion with confidence score
# -----------------------------------------------------------------------------
Sleep 2s

# Navigate through suggestions
Type "j"
Sleep 500ms

# Accept first suggestion
Type "a"
Sleep 1s

# Accept another
Type "a"
Sleep 1s

# Reject one we don't want
Type "r"
Sleep 500ms

# Accept another
Type "a"
Sleep 1s

# Cancel review (reject remaining)
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 5: View accepted skills in list
# Skills now appear in the main list
# -----------------------------------------------------------------------------
Sleep 1s

# Navigate to see skills
Type "j"
Sleep 200ms
Type "j"
Sleep 500ms

# View a skill detail
Enter
Sleep 2s

# Close detail
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
