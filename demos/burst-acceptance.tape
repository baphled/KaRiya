# Burst Acceptance Demo
# Demonstrates the enrichment review pipeline for bursts:
# 1. Two events are pre-seeded (SuggestBursts requires ≥2 events)
# 2. Capture a third event related to the existing theme
# 3. Review screen shows inferred bursts ('b' badge)
# 4. Open burst suggestion modal and accept the burst
# 5. Confirm review — burst is persisted to the database

Output demos/burst-acceptance.gif

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

# Set up a clean demo environment with 2 pre-seeded events (hidden from recording)
Hide
Type "./demos/setup-burst-demo.sh"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 500ms
Show

# Introduce what we're demonstrating
Type "# Burst Acceptance: capture event → infer bursts → accept → persist"
Enter
Sleep 2s

# Launch KaRiya with the isolated demo environment
Type `export HOME=$(pwd)/demos/temp_demo_env && ./kariya --db $HOME/.kariya/events.db`
Enter

# Wait for main menu to fully render (app runs DB migrations on startup)
Sleep 4s

# Select Capture Event (first menu item — press Enter immediately)
Enter
Sleep 2s

# Strategy screen: Quick is the first option — select it
Enter
Sleep 2s

# Quick capture form: type a related event (triggers burst detection with pre-seeded events)
Type "Completed database migration reducing query time by 50%"
Sleep 500ms

# Submit the form with Ctrl+S (quick submit — skips navigating to Submit button)
Ctrl+S
Sleep 8s

# Success modal auto-dismisses after countdown
# Enrichment review screen is now visible
# The 'b' badge in the footer indicates bursts were detected
# Press 'b' to open the burst suggestion modal
Type "b"
Sleep 4s

# Burst suggestion modal is open
# Accept the burst suggestion with 'a'
Type "a"
Sleep 2s

# Modal auto-closes when the suggestion is processed
# Back on the review screen — press Enter to confirm and persist
Enter
Sleep 3s

# Review confirmed — burst is now persisted to the database
Sleep 2s

# Exit application
Ctrl+C
Sleep 1s
