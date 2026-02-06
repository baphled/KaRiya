# =============================================================================
# KaRiya Skills Management - Add Skill Manually
# =============================================================================
#
# WORKFLOW: Create a new skill manually via form
# DURATION: ~30 seconds
#
# STEPS DEMONSTRATED:
# 1. Navigate to Manage Skills
# 2. Press 'n' to open new skill form
# 3. Fill in skill details (name, description, proficiency)
# 4. Submit form to create skill
# 5. Verify skill appears in list
#
# FORM FIELDS:
# - Skill Name: Text input (required)
# - Description: Text area (optional)
# - Proficiency: Select (Beginner, Intermediate, Advanced, Expert)
# - Type: Select (Technical, Soft Skill, Domain Knowledge)
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/skills-add-skill.gif
Output demos/vhs/output/mp4s/skills-add-skill.mp4

# -----------------------------------------------------------------------------
# SETUP: Fresh database (no skills yet)
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
# STEP 2: Navigate to Manage Skills
# -----------------------------------------------------------------------------
Down
Down
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 3: Open new skill form (n key)
# -----------------------------------------------------------------------------
Type "n"
Sleep 1500ms

# -----------------------------------------------------------------------------
# STEP 4: Fill in skill form
# -----------------------------------------------------------------------------
# Skill Name (text input)
Type "Go Programming"
Tab
Sleep 300ms

# Description (text area)
Type "Backend development with Go, including microservices, APIs, and CLI tools"
Tab
Sleep 300ms

# Proficiency (select) - Navigate to "Advanced"
Down
Down
Sleep 200ms
Enter
Tab
Sleep 300ms

# Type (select) - "Technical" is default, keep it
Enter
Tab
Sleep 300ms

# -----------------------------------------------------------------------------
# STEP 5: Submit form
# -----------------------------------------------------------------------------
Enter
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 6: Verify skill appears in list
# The new skill should be visible in the skills list
# -----------------------------------------------------------------------------
Sleep 1s

# View the newly created skill
Enter
Sleep 2s

# Close detail
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
