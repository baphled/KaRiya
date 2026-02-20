# KaRiya Golden File Test
# Generates screenshots for visual regression testing
#
# This tape captures key screens for golden file comparison.
# Run: vhs demos/vhs/golden-test.tape
# Compare: ./demos/vhs/scripts/compare-golden.sh

Source demos/vhs/config.tape

# Setup environment (hidden)
Hide
Type `mkdir -p /tmp/kariya-golden`
Enter
Sleep 300ms
Type `cp demos/vhs/minimal-config.yaml /tmp/kariya-golden/config.yaml`
Enter
Sleep 300ms
Type "cd /home/baphled/Projects/KaRiya"
Enter
Sleep 300ms
# Import events silently
Type `./kariya --config /tmp/kariya-golden/config.yaml --db /tmp/kariya-golden/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-golden/config.yaml --db /tmp/kariya-golden/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Screenshot: Main Menu
Screenshot demos/vhs/output/screenshots/main-menu.png

# Navigate to Capture Event
Enter
Sleep 1s

# Screenshot: Capture Strategy Selection
Screenshot demos/vhs/output/screenshots/capture-strategy.png

# Go back to menu
Escape
Sleep 500ms

# Navigate to Browse Timeline
Down
Enter
Sleep 2s

# Screenshot: Browse Timeline
Screenshot demos/vhs/output/screenshots/browse-timeline.png

# View event detail
Enter
Sleep 1s

# Screenshot: Event Detail
Screenshot demos/vhs/output/screenshots/event-detail.png

# Close and return to menu
Escape
Sleep 300ms
Escape
Sleep 500ms

# Navigate to Manage Skills (3rd option)
Down
Down
Enter
Sleep 2s

# Screenshot: Manage Skills
Screenshot demos/vhs/output/screenshots/manage-skills.png

# Return to menu
Escape
Sleep 500ms

# Navigate to Generate CV (4th option)
Down
Enter
Sleep 2s

# Screenshot: CV Wizard Profile
Screenshot demos/vhs/output/screenshots/cv-wizard-profile.png

# Select profile
Down
Enter
Sleep 1s

# Screenshot: CV Wizard Audience
Screenshot demos/vhs/output/screenshots/cv-wizard-audience.png

# Continue through wizard
Enter
Sleep 1s
Enter
Sleep 1s
Enter
Sleep 3s

# Screenshot: CV Preview
Screenshot demos/vhs/output/screenshots/cv-preview.png

# Cancel CV wizard and return to menu
Escape
Sleep 300ms
Escape
Sleep 500ms

# Navigate to Configure System (5th option)
Down
Down
Down
Down
Enter
Sleep 2s

# Screenshot: Configure System
Screenshot demos/vhs/output/screenshots/configure-system.png

# Return to menu
Escape
Sleep 500ms

# Navigate to Manage Bursts (6th option)
Down
Enter
Sleep 2s

# Screenshot: Manage Bursts
Screenshot demos/vhs/output/screenshots/manage-bursts.png

# Return to menu
Escape
Sleep 500ms

# Navigate to Manage Facts (7th option)
Down
Enter
Sleep 2s

# Screenshot: Manage Facts
Screenshot demos/vhs/output/screenshots/manage-facts.png

# Exit
Ctrl+C
Sleep 500ms
