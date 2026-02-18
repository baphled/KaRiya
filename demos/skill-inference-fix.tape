# PR #172: Skill Inference Fix Demo
# Demonstrates the complete skill inference lifecycle:
# 1. Trigger inference on an event
# 2. See suggestions (BubbleTea appears because it's not linked to THIS event)
# 3. Accept the suggestion
# 4. Skill is saved and linked to the event
#
# Setup: One event ("Used BubbleTea for UI development") and one skill
# ("BubbleTea") that exists globally but is NOT linked to the event.
# The fix ensures inference suggests "BubbleTea" as new for this event.

Output demos/skill-inference-fix.gif

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

# Setup demo environment (hidden from recording)
Hide
Type "./demos/setup-demo.sh"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 500ms
Show

# Show what we're demonstrating
Type "# PR #172: Skill Inference Fix - event-skill linkage check"
Enter
Sleep 2s

# Launch KaRiya with demo environment
Type `export HOME=$(pwd)/demos/temp_demo_env && ./kariya --db $HOME/.kariya/events.db`
Enter

# Wait for main menu to render
Sleep 3s

# Navigate to Browse Timeline (second menu item)
Down
Sleep 500ms
Enter
Sleep 2s

# Select the event "Used BubbleTea for UI development"
Enter
Sleep 2s

# Open skills view for this event (press 's')
Type "s"
Sleep 2s

# Trigger skill inference with 'i' key
Type "i"
Sleep 4s

# Inference results appear - BubbleTea should be suggested
# because it's not linked to THIS event (even though it exists globally)
Sleep 2s

# Accept the "BubbleTea" skill suggestion with 'a' key
# This saves it and links it to the current event
Type "a"
Sleep 3s

# Confirmation feedback shows skill was created
# The modal auto-closes when all suggestions are processed
Sleep 2s

# Close any remaining feedback modal
Escape
Sleep 1s

# Close skills modal (back to event detail)
Escape
Sleep 1s

# Close event detail (back to timeline)
Escape
Sleep 1s

# ---------------------------------------------------------
# VERIFICATION STEP:
# Re-open everything to prove the skill was actually saved
# ---------------------------------------------------------

# Re-open the same event
Enter
Sleep 2s

# Re-open skills view
Type "s"
Sleep 3s

# Skill "BubbleTea" should now be listed in the skills list
# This confirms persistence worked!
Sleep 2s

# Exit application
Ctrl+C
Sleep 1s
