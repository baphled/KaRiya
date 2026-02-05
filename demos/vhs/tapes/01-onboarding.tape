# KaRiya Onboarding Demo
# Duration: ~30 seconds
# Shows first-time user experience with the 3-step onboarding wizard
#
# This demo shows what happens when a new user launches KaRiya for the first time.
# The onboarding wizard collects basic profile information.
#
# Uses empty-config.yaml which has no profile data, triggering onboarding.

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/01-onboarding.gif
Output demos/vhs/output/mp4s/01-onboarding.mp4

# Clean start for fresh onboarding (hidden setup)
Hide
Type `rm -rf /tmp/kariya-demo && mkdir -p /tmp/kariya-demo`
Enter
Sleep 300ms
# Copy empty config to trigger onboarding wizard
Type `cp demos/vhs/empty-config.yaml /tmp/kariya-demo/config.yaml`
Enter
Sleep 300ms
Type "cd /home/baphled/Projects/KoRiya"
Enter
Sleep 300ms
Show

# Launch KaRiya with empty config (triggers onboarding wizard)
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Sleep 300ms
Enter

# Wait for onboarding wizard Step 1
Sleep 2s

# Step 1: Enter name
Type@100ms "Alex Chen"
Sleep 500ms
Enter

# Wait for Step 2
Sleep 1s

# Step 2: Enter email and location
Type@100ms "alex.chen@techcorp.io"
Tab
Sleep 300ms
Type@100ms "London, UK"
Sleep 500ms
Enter

# Wait for Step 3
Sleep 1s

# Step 3: Professional details (optional but we'll fill them)
Type@100ms "Senior Software Engineer"
Tab
Sleep 300ms
Type@100ms "alexchendev"
Tab
Sleep 300ms
Type@100ms "https://alexchen.dev"
Sleep 500ms
Enter

# Wait for main menu to appear
Sleep 2s

# Navigate menu to show available options
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Down
Sleep 200ms
Up
Sleep 200ms
Up
Sleep 200ms
Up
Sleep 200ms
Up
Sleep 1s

# Exit gracefully
Ctrl+C
Sleep 500ms
