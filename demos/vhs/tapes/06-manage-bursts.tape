# KaRiya Manage Bursts Demo
# Duration: ~25 seconds
# Shows burst management: browsing career bursts, viewing details
#
# Bursts are periods of focused career activity. This demo shows
# navigating the burst list, viewing burst details, and exploring
# related events.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/06-manage-bursts.gif
Output demos/vhs/output/mp4s/06-manage-bursts.mp4

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
# Import events silently (bursts are detected from events)
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Manage Bursts (sixth option: Capture, Browse, Skills, CV, Configure, Bursts)
Down
Down
Down
Down
Down
Enter

# Wait for burst screen to load
Sleep 2s

# Navigate through bursts using vim keys
Type "j"
Sleep 200ms
Type "j"
Sleep 200ms
Type "k"
Sleep 500ms

# View burst detail
Enter
Sleep 1500ms

# View events for this burst (if available)
Type "v"
Sleep 1500ms

# Close events view
Escape
Sleep 500ms

# Close detail modal
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
