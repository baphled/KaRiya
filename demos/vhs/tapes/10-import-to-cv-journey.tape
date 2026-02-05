# KaRiya Complete Journey: Import to CV
# Duration: ~90 seconds
# Shows the full workflow: Import CSV -> Browse timeline -> Generate CV
#
# This is the main demo showing the complete user journey from importing
# career data to generating a tailored CV.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/10-complete-journey.gif
Output demos/vhs/output/mp4s/10-complete-journey.mp4

# === PART 1: Show and Import Data ===

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

# Show CSV preview
Type "echo '=== Career Events to Import ===' && head -8 demos/vhs/sample-events.csv"
Enter
Sleep 3s

# Import the data
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review"
Enter
Sleep 4s

# === PART 2: Launch and Browse Timeline ===

Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Browse Timeline (second option)
Down
Enter

# Wait for timeline to load
Sleep 2s

# Scroll through events to show the imported data
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 500ms

# View an event detail
Enter
Sleep 2s

# Close detail view
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# === PART 3: Generate CV ===

# Navigate to Generate CV (fourth option)
Down
Down
Down
Enter

# Wait for CV wizard
Sleep 2s

# Step 1: Select Staff Engineer profile
Down
Sleep 300ms
Enter

# Wait for audience selection
Sleep 1s

# Step 2: Select Hiring Manager (default)
Enter

# Wait for technology focus
Sleep 1s

# Step 3: Technology focus (Language Agnostic)
Enter

# Wait for skills format
Sleep 1s

# Step 4: Skills format
Enter

# Wait for CV generation
Sleep 4s

# Scroll through CV preview
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 1s

# Show we're done - pause on the preview
Sleep 2s

# Exit
Ctrl+C
Sleep 500ms
