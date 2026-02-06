# Sad Path Demo Template
# ======================
# Demonstrates error handling and validation.
#
# Copy this file to demos/vhs/features/YOUR-FEATURE/sad-path.tape
#
# Usage:
#   1. Copy to demos/vhs/features/your-feature/sad-path.tape
#   2. Update the output paths and error scenarios
#   3. Run: vhs demos/vhs/features/your-feature/sad-path.tape
#
# Purpose:
#   - Show how errors are handled gracefully
#   - Demonstrate validation messages
#   - Prove recovery from error states
#   - Document user-friendly error experience

Source demos/vhs/config.tape

# Update these output paths for your feature
Output demos/vhs/features/YOUR-FEATURE/sad-path.gif
Output demos/vhs/features/YOUR-FEATURE/sad-path.mp4

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

# === SAD PATH SCENARIOS ===
#
# Scenario 1: Validation Error
# ----------------------------
# Trigger validation by submitting empty/invalid data
#
# Example:
# Down
# Enter
# Sleep 1s
# # Try to submit without required field
# Enter
# Sleep 1s
# Screenshot demos/vhs/features/YOUR-FEATURE/validation-error.png
#
# Scenario 2: Show Error Message
# ------------------------------
# Demonstrate clear error messaging
#
# Example:
# Type@100ms "invalid-data-format"
# Tab
# Enter
# Sleep 1s
# Screenshot demos/vhs/features/YOUR-FEATURE/error-message.png
#
# Scenario 3: Recovery from Error
# -------------------------------
# Show user can fix and retry
#
# Example:
# # Clear and enter valid data
# Ctrl+A
# Type@100ms "valid-data"
# Enter
# Sleep 2s
# Screenshot demos/vhs/features/YOUR-FEATURE/error-recovered.png

# Exit cleanly
Escape
Sleep 500ms
Ctrl+C
Sleep 500ms
