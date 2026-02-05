# KaRiya Browse Timeline Demo
# Duration: ~20 seconds
# Shows browsing career events, viewing details, and filtering
#
# This demo shows navigation through the timeline, event detail view,
# and the filter functionality.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/03-browse-timeline.gif
Output demos/vhs/output/mp4s/03-browse-timeline.mp4

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
# Import events silently
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Browse Timeline (second option)
Down
Enter

# Wait for timeline to load
Sleep 2s

# Navigate through events using vim keys
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "k"
Sleep 200ms
Type "k"
Sleep 500ms

# View event detail
Enter
Sleep 2s

# Close detail
Escape
Sleep 500ms

# Open filter modal
Type "f"
Sleep 1s

# Close filter (escape)
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
