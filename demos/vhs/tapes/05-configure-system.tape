# KaRiya Configure System Demo
# Duration: ~30 seconds
# Shows system configuration: domain selection, editing settings, reviewing changes
#
# This demo shows navigating configuration domains, editing settings,
# and the review/confirm workflow.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/05-configure-system.gif
Output demos/vhs/output/mp4s/05-configure-system.mp4

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

# Navigate to Configure System (fifth option: Capture, Browse, Skills, CV, Configure)
Down
Down
Down
Down
Enter

# Wait for configuration screen to load
Sleep 2s

# Navigate through configuration domains
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "k"
Sleep 500ms

# Select a domain to edit (Profile domain)
Type "j"
Sleep 200ms
Enter

# Wait for edit modal to load
Sleep 1500ms

# Navigate through form fields
Tab
Sleep 300ms
Tab
Sleep 300ms
Tab
Sleep 500ms

# Cancel edit and go back
Escape
Sleep 1s

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
