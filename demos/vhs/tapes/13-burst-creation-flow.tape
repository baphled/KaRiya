# KaRiya Burst Creation Flow Demo
# Duration: ~120 seconds
# Shows complete workflow: creating related events -> AI burst detection -> fact extraction
#
# This demo demonstrates the full burst workflow:
# 1. Create several related career events (same project/timeframe)
# 2. Navigate to Burst Management
# 3. Trigger AI burst detection
# 4. Review and accept the suggested burst
# 5. Facts are automatically extracted
# 6. View the burst with its associated events and facts
#
# This showcases how KaRiya automatically identifies patterns in your
# career history and groups related work into meaningful "bursts".
#
# Note: Events use Project (not Company) to group related work together.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/13-burst-creation-flow.gif
Output demos/vhs/output/mp4s/13-burst-creation-flow.mp4

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

# ============================================
# PART 1: Create first event (Project Kickoff)
# ============================================

# Select Capture Event
Enter
Sleep 1s

# Select Manual mode (second option)
Down
Enter
Sleep 1s

# Field 1: Event Description (text area)
Type@50ms "Led architecture design for new microservices platform, defining service boundaries and communication patterns"
Tab
Sleep 300ms

# Field 2: Date (text input)
Type@80ms "2025-01-15"
Tab
Sleep 300ms

# Field 3: Company (skip - using project instead)
Tab
Sleep 300ms

# Field 4: Project (text input)
Type@80ms "Platform Modernization"
Tab
Sleep 300ms

# Field 5: Tags (multi-select)
Down
Space
Sleep 200ms
Down
Space
Sleep 200ms
Tab
Sleep 300ms

# Field 6: Categories (multi-select)
Down
Space
Sleep 200ms
Tab
Sleep 300ms

# Field 7: Submit
Enter

# Wait for save modal and success
Sleep 3s

# Review screen - confirm to complete
Enter
Sleep 2s

# ============================================
# PART 2: Create second event (Implementation)
# ============================================

# Capture Event again
Enter
Sleep 1s

# Manual mode
Down
Enter
Sleep 1s

# Field 1: Event Description
Type@50ms "Implemented core authentication service with JWT tokens and OAuth2 integration, handling 10k requests per second"
Tab
Sleep 300ms

# Field 2: Date
Type@80ms "2025-02-10"
Tab
Sleep 300ms

# Field 3: Company (skip - using project)
Tab
Sleep 300ms

# Field 4: Project
Type@80ms "Platform Modernization"
Tab
Sleep 300ms

# Field 5: Tags (multi-select)
Down
Space
Sleep 200ms
Down
Down
Space
Sleep 200ms
Tab
Sleep 300ms

# Field 6: Categories (multi-select)
Down
Space
Sleep 200ms
Tab
Sleep 300ms

# Field 7: Submit
Enter

# Wait for save modal and success
Sleep 3s

# Review screen - confirm
Enter
Sleep 2s

# ============================================
# PART 3: Create third event (Completion)
# ============================================

# Capture Event again
Enter
Sleep 1s

# Manual mode
Down
Enter
Sleep 1s

# Field 1: Event Description
Type@50ms "Successfully deployed microservices platform to production, reducing deployment time from 2 hours to 15 minutes"
Tab
Sleep 300ms

# Field 2: Date
Type@80ms "2025-03-20"
Tab
Sleep 300ms

# Field 3: Company (skip - using project)
Tab
Sleep 300ms

# Field 4: Project
Type@80ms "Platform Modernization"
Tab
Sleep 300ms

# Field 5: Tags (multi-select)
Down
Down
Space
Sleep 200ms
Down
Space
Sleep 200ms
Tab
Sleep 300ms

# Field 6: Categories (multi-select)
Down
Space
Sleep 200ms
Tab
Sleep 300ms

# Field 7: Submit
Enter

# Wait for save modal and success
Sleep 3s

# Review screen - confirm
Enter
Sleep 2s

# ============================================
# PART 4: Detect and review burst
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

# Wait for AI processing (loading modal)
Sleep 5s

# Suggestion review modal appears
Sleep 2s

# Accept the suggested burst with 'a'
# This will save the burst AND extract facts automatically
Type "a"

# Wait for burst creation and fact extraction
Sleep 4s

# Close any result modals
Escape
Sleep 1s

# View the new burst detail
Enter
Sleep 2s

# View events in this burst with 'v'
Type "v"
Sleep 2s

# Close events view
Escape
Sleep 500ms

# View facts extracted from this burst with 'f'
Type "f"
Sleep 2s

# Close facts view
Escape
Sleep 500ms

# Close detail
Escape
Sleep 500ms

# Return to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
