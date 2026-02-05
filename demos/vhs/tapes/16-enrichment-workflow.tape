# KaRiya Enrichment Workflow Demo
# Duration: ~150 seconds
# Shows how inferred bursts, facts, and skills emerge from related events
#
# This demo demonstrates the complete enrichment workflow:
# 1. Create several career events that share a common theme (same project)
# 2. Run burst detection to identify the pattern
# 3. Accept the burst suggestion (facts are auto-extracted)
# 4. View the extracted facts in Manage Facts
# 5. Run skill inference to detect skills from the events
# 6. Review and accept skill suggestions
#
# This showcases how KaRiya automatically:
# - Groups related work into "bursts" (career themes)
# - Extracts factual achievements from your descriptions
# - Infers skills from technologies and methodologies mentioned

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/16-enrichment-workflow.gif
Output demos/vhs/output/mp4s/16-enrichment-workflow.mp4

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
# PART 1: Create first event (Design Phase)
# ============================================

# Select Capture Event
Enter
Sleep 1s

# Select Manual mode
Down
Enter
Sleep 1s

# Event Description - detailed achievement with metrics
Type@40ms "Designed cloud-native architecture using Kubernetes and microservices patterns, reducing infrastructure costs by 40% and improving deployment frequency from monthly to daily releases"
Tab
Sleep 300ms

# Date
Type@80ms "2025-01-10"
Tab
Sleep 300ms

# Company (skip - using project)
Tab
Sleep 300ms

# Project
Type@80ms "Cloud Migration Initiative"
Tab
Sleep 300ms

# Tags - select relevant ones
Down
Space
Sleep 150ms
Down
Space
Sleep 150ms
Tab
Sleep 300ms

# Categories
Down
Space
Sleep 150ms
Tab
Sleep 300ms

# Submit
Enter
Sleep 3s

# Confirm on review screen
Enter
Sleep 2s

# ============================================
# PART 2: Create second event (Implementation)
# ============================================

Enter
Sleep 1s

Down
Enter
Sleep 1s

# Event Description - implementation achievement
Type@40ms "Led team of 5 engineers implementing containerization strategy with Docker and Kubernetes, achieving 99.9% uptime and reducing mean time to recovery from 4 hours to 15 minutes"
Tab
Sleep 300ms

# Date
Type@80ms "2025-02-15"
Tab
Sleep 300ms

# Company (skip)
Tab
Sleep 300ms

# Project - same project to create burst pattern
Type@80ms "Cloud Migration Initiative"
Tab
Sleep 300ms

# Tags
Down
Space
Sleep 150ms
Down
Down
Space
Sleep 150ms
Tab
Sleep 300ms

# Categories
Down
Space
Sleep 150ms
Tab
Sleep 300ms

# Submit
Enter
Sleep 3s

# Confirm
Enter
Sleep 2s

# ============================================
# PART 3: Create third event (Completion)
# ============================================

Enter
Sleep 1s

Down
Enter
Sleep 1s

# Event Description - completion achievement
Type@40ms "Successfully migrated 50 legacy applications to cloud platform, mentoring junior developers on DevOps practices and establishing CI/CD pipelines that reduced deployment errors by 85%"
Tab
Sleep 300ms

# Date
Type@80ms "2025-03-20"
Tab
Sleep 300ms

# Company (skip)
Tab
Sleep 300ms

# Project - same project
Type@80ms "Cloud Migration Initiative"
Tab
Sleep 300ms

# Tags
Down
Down
Space
Sleep 150ms
Down
Space
Sleep 150ms
Tab
Sleep 300ms

# Categories
Down
Space
Sleep 150ms
Tab
Sleep 300ms

# Submit
Enter
Sleep 3s

# Confirm
Enter
Sleep 2s

# ============================================
# PART 4: Detect bursts from events
# ============================================

# Navigate to Manage Bursts (6th option)
Down
Down
Down
Down
Down
Enter
Sleep 2s

# Trigger burst detection
Type "s"
Sleep 5s

# Suggestion modal appears - shows detected burst
Sleep 2s

# Accept the burst (this also extracts facts automatically)
Type "a"
Sleep 4s

# Close result modal
Escape
Sleep 1s

# ============================================
# PART 5: View the burst and its facts
# ============================================

# View burst detail
Enter
Sleep 2s

# View events in this burst
Type "v"
Sleep 2s

# Close events
Escape
Sleep 500ms

# View facts extracted from this burst
Type "f"
Sleep 2s

# Close facts
Escape
Sleep 500ms

# Close detail
Escape
Sleep 500ms

# ============================================
# PART 6: View facts in Manage Facts
# ============================================

# Back to main menu
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
Sleep 2s

# Navigate through extracted facts
Type "j"
Sleep 300ms
Type "j"
Sleep 300ms

# View a fact detail
Enter
Sleep 2s

# Close detail
Escape
Sleep 500ms

# ============================================
# PART 7: Infer skills from events
# ============================================

# Back to main menu
Escape
Sleep 1s

# Navigate to Manage Skills (3rd option)
Down
Down
Enter
Sleep 2s

# Trigger skill inference
Type "i"
Sleep 5s

# Skill suggestions appear
Sleep 2s

# Accept a skill (e.g., Kubernetes)
Type "a"
Sleep 2s

# Accept another skill (e.g., Docker)
Type "a"
Sleep 2s

# Accept another skill
Type "a"
Sleep 2s

# View events for this skill
Type "v"
Sleep 2s

# Close events view
Escape
Sleep 500ms

# Finish skill review
Escape
Sleep 1s

# Show skills list with new skills
Sleep 1s

# Navigate to see skills
Type "j"
Sleep 300ms
Type "j"
Sleep 500ms

# ============================================
# PART 8: Return to main menu
# ============================================

# Back to main menu
Escape
Sleep 1s

# Exit
Ctrl+C
Sleep 500ms
