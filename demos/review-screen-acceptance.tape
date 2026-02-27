# ReviewScreen & PreviewScreen Acceptance Demo
# Demonstrates the new ReviewScreen wizard selections display:
# 1. Navigate to Generate CV from main menu
# 2. Complete the CV configuration wizard (defaults are fine)
# 3. Wait for CV generation to complete
# 4. ReviewScreen shows "YOUR SELECTIONS" and "GENERATED OUTPUT" sections
# 5. Navigate to PreviewScreen to view the generated CV content

Output demos/review-screen-acceptance.gif

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

# Set up a clean demo environment with seeded data (hidden from recording)
Hide
Type "./demos/setup-cv-demo.sh"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 500ms
Show

# Introduce what we're demonstrating
Type "# ReviewScreen: wizard selections display after CV generation"
Enter
Sleep 2s

# Launch KaRiya with the isolated demo environment
Type `export HOME=$(pwd)/demos/temp_demo_env && ./kariya --db $HOME/.kariya/events.db`
Enter

# Wait for main menu to fully render
Sleep 4s

# Navigate to Generate CV (index 3 in menu)
# Menu: 0=Capture Event, 1=Browse Timeline, 2=Manage Skills, 3=Generate CV
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms

# Select Generate CV
Enter
Sleep 3s

# CV Configuration Wizard is now visible
# Skip wizard with defaults (Ctrl+S uses pre-selected profile and audience)
Ctrl+S

# Wait for technology extraction and CV generation to complete
Sleep 12s

# ReviewScreen is now visible with:
# - Title: "CV Configuration Review"
# - Personal Details, CV Details, Summary sections
# Scroll down to show Generation Settings and Statistics sections
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms

Screenshot demos/review-screen-screenshot.png
Sleep 3s

# Scroll back to top before navigating to preview
Type "g"
Sleep 500ms

# Navigate to PreviewScreen to view the generated CV
Enter
Sleep 3s

# PreviewScreen is now visible with full CV content
# Scroll down to show CV sections
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms
Down
Sleep 500ms

# Take screenshot of PreviewScreen
Screenshot demos/preview-screen-screenshot.png
Sleep 2s

# Exit application
Ctrl+C
Sleep 1s
