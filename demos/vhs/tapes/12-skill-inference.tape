# KaRiya AI Skill Inference Demo
# Duration: ~50 seconds
# Shows AI-powered skill inference from career events
#
# This demo demonstrates the AI feature that analyzes all career events
# to automatically detect and suggest skills based on technologies,
# tools, and methodologies mentioned in the user's career history.
#
# Workflow:
# 1. Navigate to Manage Skills
# 2. Press 'i' to trigger skill inference
# 3. AI analyzes all events for technology mentions
# 4. Review suggested skills
# 5. Accept (a), Reject (r), or View events (v)
# 6. Accepted skills are saved to your profile

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/12-skill-inference.gif
Output demos/vhs/output/mp4s/12-skill-inference.mp4

# Setup with pre-imported data (hidden)
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
# Import events (provides data for skill inference)
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Manage Skills (third option)
Down
Down
Enter

# Wait for skills screen to load
Sleep 2s

# Trigger AI skill inference with 'i' key
Type "i"

# Wait for AI processing (loading modal: "Analyzing all events for skills...")
Sleep 5s

# AI suggestion review modal appears
Sleep 2s

# Navigate through skill suggestions
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms

# Accept a skill suggestion with 'a'
# This skill will be added to your profile
Type "a"
Sleep 2s

# View events that suggest this skill with 'v'
# Shows which career events mention this technology
Type "v"
Sleep 2s

# Close events view
Escape
Sleep 500ms

# Accept another skill
Type "a"
Sleep 2s

# Accept one more
Type "a"
Sleep 2s

# Reject a skill suggestion with 'r'
# (Maybe this skill isn't accurate for your profile)
Type "r"
Sleep 1s

# Finish reviewing - cancel remaining suggestions
Escape
Sleep 1s

# Show the list now has new skills added
Sleep 1s

# Navigate to see newly added skills
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 500ms

# View a skill detail
Enter
Sleep 2s

# Close detail
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
