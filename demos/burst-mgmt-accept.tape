# Burst Management Accept Demo
# Demonstrates accepting a burst suggestion via the burst_management intent:
# 1. Five events are pre-seeded (same company, project, similar tags)
# 2. Navigate to Manage Bursts from the main menu
# 3. Press 's' to trigger burst suggestion detection
# 4. Press 'a' to accept the first burst suggestion
# 5. OBSERVE what happens (this is the key debugging step)

Output demos/burst-mgmt-accept.gif

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

# Set up a clean demo environment with 5 pre-seeded events (hidden from recording)
Hide
Type "./demos/setup-burst-mgmt-demo.sh"
Enter
Sleep 1s
Type "clear"
Enter
Sleep 300ms
Show

# Introduce what we're demonstrating
Type "# Burst Management: navigate → suggest bursts → accept → observe behaviour"
Enter
Sleep 1s

# Launch KaRiya with the isolated demo environment
Type `export HOME=$(pwd)/demos/temp_demo_env && ./kariya --db $HOME/.kariya/events.db`
Enter

# Wait for main menu to fully render (app runs DB migrations on startup)
Sleep 3s

# Navigate to Manage Bursts (index 5 in menu: 0=Capture, 1=Browse, 2=Skills, 3=CV, 4=Config, 5=Bursts)
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms

# Select Manage Bursts
Enter
Sleep 4s

# Now in burst_management intent - press 's' to suggest bursts (triggers detection)
Type "s"
Sleep 4s

# Detection should complete — suggestions modal shown (if any bursts detected)
# Let it sit for a moment so we can see what's displayed
Sleep 2s

# Press 'a' to accept the first suggestion
Type "a"
Sleep 3s

# OBSERVE what happens — this is the KEY DEBUGGING STEP
# The BDD scenario shows "Review Skill Suggestions" modal appears instead of success
# We want to visually confirm this behaviour
Sleep 3s

# Check if we need to interact further (escape or enter)
# Let's try pressing escape to see the state
Escape
Sleep 2s

# Let it sit so we can see the final state
Sleep 2s

# Exit application
Ctrl+C
Sleep 500ms
