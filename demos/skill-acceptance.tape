# Skill Acceptance Demo
# Demonstrates the enrichment review pipeline for skills:
# 1. Capture an event with technology keywords
# 2. Review screen shows inferred skills ('s' badge)
# 3. Open skill suggestion modal and accept all skills
# 4. Confirm review — skills are persisted and linked to the event

Output demos/skill-acceptance.gif

Set Shell "bash"
Set FontSize 14
Set FontFamily "JetBrains Mono"
Set Width 1200
Set Height 800
Set Padding 15
Set BorderRadius 8
Set WindowBar Colorful
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Framerate 30

# Set up a clean demo environment (hidden from recording)
Hide
Type "./demos/setup-skill-demo.sh"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 500ms
Show

# Introduce what we're demonstrating
Type "# Skill Acceptance: capture event → infer skills → accept → persist"
Enter
Sleep 2s

# Launch KaRiya with the isolated demo environment
Type `export HOME=$(pwd)/demos/temp_demo_env && ./kariya --db $HOME/.kariya/events.db`
Enter

# Wait for main menu to fully render
Sleep 4s

# Select Capture Event (first menu item — press Enter immediately)
Enter
Sleep 2s

# Strategy screen: Quick is the first option — select it
Enter
Sleep 2s

# Quick capture form: type an event with technology keywords (triggers skill inference)
Type "Built REST API in Go using Gin framework and PostgreSQL"
Sleep 500ms

# Submit the form (Enter)
Enter
Sleep 5s

# Review screen now shows inferred results
# The 's' badge in the footer indicates skills were inferred
# Press 's' to open the skill suggestion modal
Type "s"
Sleep 3s

# Skill suggestion modal is open
# Accept the first skill suggestion with 'a'
Type "a"
Sleep 2s

# Accept any additional skill suggestions
Type "a"
Sleep 2s

# Accept any further skill suggestions
Type "a"
Sleep 2s

# Modal auto-closes when all suggestions are processed
# Back on the review screen — press Enter to confirm and persist
Enter
Sleep 3s

# Review confirmed — skills are now persisted and linked to the event
Sleep 2s

# Exit application
Ctrl+C
Sleep 1s
