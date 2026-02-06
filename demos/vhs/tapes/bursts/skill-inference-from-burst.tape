# =============================================================================
# KaRiya Burst Management - Skill Inference from Burst
# =============================================================================
# 
# WORKFLOW: User infers skills from a confirmed burst's events
# DURATION: ~50 seconds
# 
# WHAT IS SKILL INFERENCE?
# After a burst is confirmed, the AI can analyze the events within it to
# detect skills, technologies, and methodologies used. This helps build
# your skill profile automatically from your career achievements.
#
# PREREQUISITES:
# - Must have a CONFIRMED burst (facts already extracted)
# - Burst must have events with technology/skill mentions
#
# STEPS DEMONSTRATED:
# 1. Navigate to Burst Management
# 2. Select an existing burst
# 3. View burst detail
# 4. Trigger skill inference (i key)
# 5. Review AI skill suggestions
# 6. Accept/reject skills
# 7. View events that support each skill
#
# SKILL SUGGESTION REVIEW:
# - a: Accept skill (adds to your profile)
# - r: Reject skill (skip this suggestion)
# - v: View events that mention this skill
# - j/k: Navigate between suggestions
# - Escape: Finish review
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/bursts-skill-inference.gif
Output demos/vhs/output/mp4s/bursts-skill-inference.mp4

# -----------------------------------------------------------------------------
# SETUP: Import events and detect bursts first
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
# Import sample events
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
# STEP 2: Navigate to Manage Bursts
# -----------------------------------------------------------------------------
Down
Down
Down
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: First, detect and accept a burst (need confirmed burst for skill inference)
# -----------------------------------------------------------------------------
# Trigger burst detection
Type "s"
Sleep 5s

# Accept a burst
Type "a"
Sleep 3s

# Cancel remaining suggestions
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 4: Select the burst we just created
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 5: Trigger skill inference (i key)
# -----------------------------------------------------------------------------
# This analyzes events in the burst to detect skills
Type "i"

# Wait for AI processing
Sleep 5s

# Skill suggestion modal appears
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 6: Review and accept/reject skill suggestions
# -----------------------------------------------------------------------------
# Accept first skill suggestion
Type "a"
Sleep 2s

# Accept another
Type "a"
Sleep 2s

# View events that support the current skill (v key)
Type "v"
Sleep 2s

# Close events view
Escape
Sleep 500ms

# Accept one more skill
Type "a"
Sleep 2s

# Reject a skill we don't want
Type "r"
Sleep 1s

# Finish skill review
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 7: Close detail and return to main menu
# -----------------------------------------------------------------------------
Escape
Sleep 500ms
Escape
Sleep 1s

# -----------------------------------------------------------------------------
# CLEANUP
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
