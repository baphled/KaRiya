# KaRiya Capture Event Cancel Demo
# Duration: ~20 seconds
# Shows how to cancel event capture at various stages
#
# This demo shows:
# - Starting to fill a form
# - Cancelling with Escape
# - Returning to main menu without saving

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/14-capture-cancel.gif
Output demos/vhs/output/mp4s/14-capture-cancel.mp4

# Setup (hidden)
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
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Select Capture Event
Enter

# Wait for strategy selection
Sleep 1s

# Select Manual mode
Down
Enter

# Wait for event form
Sleep 1s

# Start typing an event description
Type@60ms "Started working on a new feature but changed my mind..."
Sleep 1s

# Cancel the form with Escape
Escape
Sleep 1s

# Back at strategy selection - cancel again
Escape
Sleep 1s

# Back at main menu - event was NOT saved
Sleep 1s

# Demonstrate we can start over
Enter
Sleep 1s

# Select Quick mode this time
Enter
Sleep 1s

# Type something briefly
Type@60ms "Quick test event"
Sleep 500ms

# Cancel again
Escape
Sleep 1s

# Cancel from strategy
Escape
Sleep 1s

# Back at main menu
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
