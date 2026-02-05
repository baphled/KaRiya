# KaRiya Fact Extraction Flow Demo
# Duration: ~60 seconds
# Shows how facts are automatically extracted when confirming a burst
#
# This demo demonstrates the fact extraction workflow:
# 1. Import sample events (pre-populated data)
# 2. Navigate to Burst Management
# 3. Run burst detection to find career patterns
# 4. Accept a burst suggestion
# 5. See facts automatically extracted
# 6. View the extracted facts in Manage Facts
#
# Facts are competency evidence statements extracted from your career events.
# They're used to generate compelling CV content.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/15-fact-extraction-flow.gif
Output demos/vhs/output/mp4s/15-fact-extraction-flow.mp4

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
# Import sample events (these will be grouped into bursts)
Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
Enter
Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# ============================================
# PART 1: Detect bursts from imported events
# ============================================

# Navigate to Manage Bursts (6th option)
Down
Down
Down
Down
Down
Enter

# Wait for burst screen
Sleep 2s

# Trigger AI burst detection with 's'
Type "s"

# Wait for AI processing (analyzing events for patterns)
Sleep 5s

# Suggestion review modal appears
Sleep 2s

# Accept first suggested burst with 'a'
# This creates the burst AND extracts facts automatically
Type "a"

# Wait for burst creation and fact extraction
Sleep 4s

# Accept another burst suggestion
Type "a"

# Wait for fact extraction
Sleep 4s

# Cancel remaining suggestions
Escape
Sleep 1s

# ============================================
# PART 2: View extracted facts
# ============================================

# Return to main menu
Escape
Sleep 1s

# Navigate to Manage Facts (7th option)
Down
Down
Down
Down
Down
Down
Enter

# Wait for facts screen to load
Sleep 2s

# Navigate through facts
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms
Type "j"
Sleep 500ms

# View fact detail
Enter
Sleep 2s

# Close detail
Escape
Sleep 500ms

# Show more facts
Type "j"
Sleep 300ms
Type "j"
Sleep 500ms

# View another fact
Enter
Sleep 2s

# Close detail
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
