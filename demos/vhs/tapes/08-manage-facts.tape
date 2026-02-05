# KaRiya Manage Facts Demo
# Duration: ~25 seconds
# Shows fact management: browsing extracted facts, viewing details
#
# Facts are key achievements and data points extracted from career events.
# This demo shows navigating the facts list, viewing fact details,
# and the edit workflow.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/08-manage-facts.gif
Output demos/vhs/output/mp4s/08-manage-facts.mp4

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
# Import events silently (facts are extracted from events)
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Manage Facts (seventh option: Capture, Browse, Skills, CV, Configure, Bursts, Facts)
Down
Down
Down
Down
Down
Down
Enter

# Wait for facts screen to load
Sleep 2s

# Navigate through facts using vim keys
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "k"
Sleep 500ms

# View fact detail
Enter
Sleep 1500ms

# Close detail view
Escape
Sleep 500ms

# Try edit mode (press e)
Type "e"
Sleep 1500ms

# Cancel edit
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
