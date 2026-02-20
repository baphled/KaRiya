# Fact Acceptance Demo
# Demonstrates the enrichment review pipeline for facts:
# 1. Capture an event with a measurable achievement description
# 2. Review screen shows inferred facts ('f' badge)
# 3. Open fact suggestion modal and accept all facts
# 4. Confirm review — facts are persisted to the database

Output demos/fact-acceptance.gif

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
Type "./demos/setup-fact-demo.sh"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 500ms
Show

# Introduce what we're demonstrating
Type "# Fact Acceptance: capture event → infer facts → accept → persist"
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

# Quick capture form: type the event text (triggers fact inference)
Type "Reduced API latency by 40% through caching optimization"
Sleep 500ms

# Submit the form with Ctrl+S (quick submit — skips navigating to Submit button)
Ctrl+S
Sleep 8s

# Success modal auto-dismisses after countdown
# Enrichment review screen is now visible
# The 'f' badge in the footer indicates facts were inferred
# Press 'f' to open the fact suggestion modal
Type "f"
Sleep 4s

# Fact suggestion modal is open
# Accept the first (and any remaining) fact suggestion with 'a'
Type "a"
Sleep 2s

# Accept any additional fact suggestions
Type "a"
Sleep 2s

# Modal auto-closes when all suggestions are processed
# Back on the review screen — press Enter to confirm and persist
Enter
Sleep 3s

# Review confirmed — facts are now persisted to the database
Sleep 2s

# Exit application
Ctrl+C
Sleep 1s
