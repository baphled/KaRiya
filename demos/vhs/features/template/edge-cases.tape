# Edge Cases Demo Template
# ========================
# Demonstrates boundary conditions and special scenarios.
#
# Copy this file to demos/vhs/features/YOUR-FEATURE/edge-cases.tape
#
# Usage:
#   1. Copy to demos/vhs/features/your-feature/edge-cases.tape
#   2. Update the output paths and edge case scenarios
#   3. Run: vhs demos/vhs/features/your-feature/edge-cases.tape
#
# Purpose:
#   - Show cancel/back navigation works correctly
#   - Demonstrate empty state handling
#   - Prove interruption recovery
#   - Document unusual but valid scenarios

Source demos/vhs/config.tape

# Update these output paths for your feature
Output demos/vhs/features/YOUR-FEATURE/edge-cases.gif
Output demos/vhs/features/YOUR-FEATURE/edge-cases.mp4

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

# === EDGE CASE SCENARIOS ===
#
# Scenario 1: Cancel Workflow
# ---------------------------
# Start workflow, then cancel at various points
#
# Example:
# Down
# Enter
# Sleep 1s
# Screenshot demos/vhs/features/YOUR-FEATURE/cancel-step1.png
# # Cancel mid-workflow
# Escape
# Sleep 1s
# Screenshot demos/vhs/features/YOUR-FEATURE/cancel-confirmed.png
#
# Scenario 2: Back Navigation
# ---------------------------
# Navigate forward, then back, verify state preserved
#
# Example:
# Down
# Enter
# Sleep 500ms
# Type@100ms "Some input"
# Sleep 500ms
# Escape
# Sleep 500ms
# # Re-enter and verify input is preserved or cleared appropriately
# Enter
# Sleep 500ms
# Screenshot demos/vhs/features/YOUR-FEATURE/back-navigation.png
#
# Scenario 3: Empty State
# -----------------------
# Show behavior when no data exists
#
# Example:
# # Navigate to list view with no items
# Down
# Enter
# Sleep 1s
# Screenshot demos/vhs/features/YOUR-FEATURE/empty-state.png
#
# Scenario 4: Boundary Values
# ---------------------------
# Test with very long input, special characters, etc.
#
# Example:
# Type@50ms "This is a very long input that tests how the UI handles text that might overflow or wrap in unexpected ways"
# Sleep 1s
# Screenshot demos/vhs/features/YOUR-FEATURE/long-input.png
#
# Scenario 5: Rapid Navigation
# ----------------------------
# Quick key presses to test responsiveness
#
# Example:
# Down
# Down
# Down
# Up
# Up
# Enter
# Sleep 500ms
# Screenshot demos/vhs/features/YOUR-FEATURE/rapid-nav.png

# Exit cleanly
Escape
Sleep 500ms
Ctrl+C
Sleep 500ms
