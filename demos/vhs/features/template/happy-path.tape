# Happy Path Demo Template
# ==========================
# Demonstrates successful workflow completion.
#
# Copy this file to demos/vhs/features/YOUR-FEATURE/happy-path.tape
#
# Usage:
#   1. Copy to demos/vhs/features/your-feature/happy-path.tape
#   2. Update the output paths and workflow steps
#   3. Run: vhs demos/vhs/features/your-feature/happy-path.tape
#
# Purpose:
#   - Show the ideal user journey
#   - Demonstrate feature working correctly
#   - Capture success states and confirmations

Source demos/vhs/config.tape

# Update these output paths for your feature
Output demos/vhs/features/YOUR-FEATURE/happy-path.gif
Output demos/vhs/features/YOUR-FEATURE/happy-path.mp4

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
# Optional: Import sample data if needed for your feature
# Type `./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db --import demos/vhs/sample-events.csv --skip-import-review > /dev/null 2>&1`
# Enter
# Sleep 2s
Show

# Launch KaRiya
Type "./kariya --config /tmp/kariya-demo/config.yaml --db /tmp/kariya-demo/demo.db"
Enter

# Wait for main menu
Sleep 2s

# === HAPPY PATH STEPS ===
#
# Document each major step with:
# 1. Action (navigation, input, selection)
# 2. Brief pause to show result
# 3. Screenshot at key moments
#
# Example workflow:
#
# Step 1: Navigate to feature
# Down
# Down
# Enter
# Sleep 1s
# Screenshot demos/vhs/features/YOUR-FEATURE/step1-feature-entry.png
#
# Step 2: Fill in required data
# Type@100ms "Example input"
# Tab
# Type@100ms "More data"
# Sleep 500ms
# Screenshot demos/vhs/features/YOUR-FEATURE/step2-form-filled.png
#
# Step 3: Submit and see success
# Enter
# Sleep 2s
# Screenshot demos/vhs/features/YOUR-FEATURE/step3-success.png

# Exit cleanly
Escape
Sleep 500ms
Ctrl+C
Sleep 500ms
