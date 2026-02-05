# KaRiya Manage Skills Demo
# Duration: ~25 seconds
# Shows skills management: browsing, viewing details, filtering
#
# This demo shows navigation through skills list, viewing skill details,
# and using the filter functionality.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/04-manage-skills.gif
Output demos/vhs/output/mp4s/04-manage-skills.mp4

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
# Import events silently (skills are inferred from events)
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Manage Skills (third option: Capture, Browse, Skills)
Down
Down
Enter

# Wait for skills screen to load
Sleep 2s

# Navigate through skills using vim keys
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "k"
Sleep 500ms

# View skill detail
Enter
Sleep 1500ms

# Close detail
Escape
Sleep 500ms

# Open filter modal
Type "f"
Sleep 1s

# Close filter (escape)
Escape
Sleep 500ms

# Open sort modal
Type "s"
Sleep 1s

# Close sort
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
