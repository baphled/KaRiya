# KaRiya AI Burst Detection Demo
# Duration: ~45 seconds
# Shows AI-powered burst detection and suggestion review workflow
#
# This demo demonstrates the AI feature that analyzes career events
# to automatically detect and suggest career "bursts" - periods of
# focused activity around specific technologies or projects.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/11-burst-detection.gif
Output demos/vhs/output/mp4s/11-burst-detection.mp4

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
# Import events silently (provides data for burst detection)
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Manage Bursts (sixth option)
Down
Down
Down
Down
Down
Enter

# Wait for burst screen to load
Sleep 2s

# Trigger AI burst detection with 's' key
Type "s"

# Wait for AI processing (loading modal appears)
Sleep 3s

# AI suggestion review modal should appear
# Navigate through suggestions
Sleep 1500ms
Type "j"
Sleep 500ms
Type "j"
Sleep 500ms

# Accept a suggestion with 'a'
Type "a"
Sleep 1500ms

# Accept another suggestion
Type "a"
Sleep 1500ms

# Reject a suggestion with 'r'
Type "r"
Sleep 1s

# Cancel review and return to list
Escape
Sleep 1s

# Show the list now has new bursts
Sleep 1s

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
