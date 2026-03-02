# Configure System Demo
# Demonstrates the new unified modal overlay for system configuration:
# 1. Press ',' from the main menu to open the configure modal
# 2. Split-panel shows domain list (left) and settings form (right)
# 3. Navigate domains with j/k — form updates automatically
# 4. Press Esc to close without saving, Ctrl+S to save

Output demos/configure-acceptance.gif

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
Type "./demos/setup-demo.sh"
Enter
Sleep 2s
Type "clear"
Enter
Sleep 500ms
Show

# Introduce what we're demonstrating
Type "# Configure System: press ',' → navigate domains → preview settings → close"
Enter
Sleep 2s

# Launch KaRiya with the isolated demo environment
Type `export HOME=$(pwd)/demos/temp_demo_env && ./kariya --db $HOME/.kariya/events.db`
Enter

# Wait for main menu to fully render
Sleep 4s

# Open Configure System modal by pressing ','
Type ","
Sleep 3s

# Modal is open showing System domain settings by default
# Navigate to Profile domain
Type "j"
Sleep 2s

# Navigate to Export domain
Type "j"
Sleep 2s

# Navigate back to Profile domain
Type "k"
Sleep 2s

# Close the modal and return to main menu
Escape
Sleep 2s

# Exit application
Ctrl+C
Sleep 1s
