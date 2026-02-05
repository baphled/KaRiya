# KaRiya CV Generation Demo
# Duration: ~60 seconds
# Shows the complete CV wizard workflow with pre-seeded career data
#
# This demo focuses on the CV generation feature, showing the wizard
# that guides users through profile, audience, and format selection,
# then displays the generated CV with scrolling preview.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/07-generate-cv.gif
Output demos/vhs/output/mp4s/07-generate-cv.mp4

# Setup: Import sample career data (hidden)
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
# Import comprehensive events for better CV generation
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events-cv.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# Navigate to Generate CV (fourth option)
Down
Sleep 150ms
Down
Sleep 150ms
Down
Sleep 150ms
Enter

# Wait for CV wizard modal to appear
Sleep 2s

# Step 1: Profile Selection - Accept default (first profile)
# The wizard shows WHO step with Profile and Audience selection
Enter
Sleep 500ms

# Step 1 continued: Select Audience - use default "Hiring Manager"
Enter
Sleep 500ms

# Step 2: TECH - Technology Focus
# Since we don't have extracted technologies, this step may be skipped
# or we select "Language Agnostic" which is the default
Sleep 1s
Enter
Sleep 500ms

# Step 3: FORMAT - Skills formatting
# Select skills format and CV length
Enter
Sleep 500ms

# Confirm final settings
Enter

# Wait for technology extraction and CV generation
# The system will show progress indicators
Sleep 4s

# Now we should be at the Review screen showing CV metadata
Sleep 2s

# Press Enter or 'p' to view the full CV preview
Type "p"
Sleep 2s

# Now we're in the CV Preview screen - scroll through content
# The CV shows personal details, sections, and bullet points

# Scroll down through the CV content
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 500ms

# Continue scrolling to see more content
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 300ms
Down
Sleep 500ms

# Scroll back up to see the top
Type "g"
Sleep 1s

# Scroll to bottom to see all content
Type "G"
Sleep 1s

# Scroll back to top
Type "g"
Sleep 1s

# Optional: Show export option by pressing 'x'
Type "x"
Sleep 1s

# Select Markdown format (navigate down)
Down
Sleep 300ms
Enter

# Wait for export modal
Sleep 1s

# Select save to file
Enter

# Wait for export confirmation
Sleep 2s

# Press Enter to complete
Enter
Sleep 1s

# Exit cleanly
Ctrl+C
Sleep 500ms
