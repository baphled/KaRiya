# Default TUI Launch Demo
# Demonstrates running kariya with no subcommand launches the interactive TUI

Output demos/default-tui-launch.gif

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
Type "# Running kariya with no arguments launches the interactive TUI"
Enter
Sleep 2s

# Launch KaRiya with the isolated demo environment
Type `export HOME=$(pwd)/demos/temp_demo_env && ./kariya --db $HOME/.kariya/events.db`
Enter

# Wait for main menu to fully render (app runs DB migrations on startup)
Sleep 4s

# Show the main menu briefly
Sleep 3s

# Exit application
Ctrl+C
Sleep 1s
