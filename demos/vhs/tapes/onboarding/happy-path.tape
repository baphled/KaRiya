# =============================================================================
# KaRiya Onboarding - Happy Path
# =============================================================================
# 
# WORKFLOW: First-time user completes the onboarding wizard
# DURATION: ~40 seconds
# 
# STEPS DEMONSTRATED:
# 1. Launch app with empty config (triggers onboarding)
# 2. Fill in Step 1: Basic profile (name, email, location)
# 3. Fill in Step 2: Professional info (title, social, website)
# 4. Fill in Step 3: Preferences (defaults)
# 5. Complete wizard and see main menu
#
# PREREQUISITES:
# - Empty config file (no profile data)
#
# KEYBOARD SHORTCUTS SHOWN:
# - Tab: Move to next field
# - Enter: Submit step / Confirm
# - Arrow keys: Navigate within fields
# =============================================================================

Source demos/vhs/config.tape

Output demos/vhs/output/gifs/onboarding-happy-path.gif
Output demos/vhs/output/mp4s/onboarding-happy-path.mp4

# -----------------------------------------------------------------------------
# SETUP: Create clean environment with empty config to trigger onboarding
# -----------------------------------------------------------------------------
Hide
Type `rm -rf /tmp/kariya-demo && mkdir -p /tmp/kariya-demo`
Enter
Sleep 300ms
# Use empty config - this triggers the onboarding wizard
Type `cp demos/vhs/empty-config.yaml /tmp/kariya-demo/config.yaml`
Enter
Sleep 300ms
Type "cd /home/baphled/Projects/KoRiya"
Enter
Sleep 300ms
Show

# -----------------------------------------------------------------------------
# STEP 1: Launch application - onboarding wizard appears automatically
# -----------------------------------------------------------------------------
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Sleep 300ms
Enter

# Wait for onboarding wizard to appear
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 2: Fill in basic profile information
# -----------------------------------------------------------------------------
# Field: Full Name
Type@100ms "Alex Chen"
Sleep 500ms
Enter
Sleep 1s

# Field: Email Address
Type@100ms "alex.chen@techcorp.io"
Tab
Sleep 300ms

# Field: Location
Type@100ms "London, UK"
Sleep 500ms
Enter
Sleep 1s

# -----------------------------------------------------------------------------
# STEP 3: Fill in professional information
# -----------------------------------------------------------------------------
# Field: Job Title
Type@100ms "Senior Software Engineer"
Tab
Sleep 300ms

# Field: Social Handle (optional)
Type@100ms "alexchendev"
Tab
Sleep 300ms

# Field: Website (optional)
Type@100ms "https://alexchen.dev"
Sleep 500ms
Enter

# Wait for wizard to complete
Sleep 2s

# -----------------------------------------------------------------------------
# STEP 4: Main menu appears - onboarding complete!
# -----------------------------------------------------------------------------
# Navigate menu to show it's working
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

# -----------------------------------------------------------------------------
# CLEANUP: Exit application
# -----------------------------------------------------------------------------
Ctrl+C
Sleep 500ms
