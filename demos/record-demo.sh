#!/bin/bash
# Demo recording script - handles config backup/restore around VHS

set -e

CONFIG_FILE="$HOME/.kariya/config.yaml"
BACKUP_FILE="$HOME/.kariya/config.yaml.demo-backup"

# Cleanup function to restore config
cleanup() {
    if [[ -f "$BACKUP_FILE" ]]; then
        echo "Restoring original config..."
        mv "$BACKUP_FILE" "$CONFIG_FILE"
        echo "Config restored."
    fi
}

# Set trap to restore config on exit (success or failure)
trap cleanup EXIT

# Backup existing config if it exists
if [[ -f "$CONFIG_FILE" ]]; then
    echo "Backing up existing config..."
    mv "$CONFIG_FILE" "$BACKUP_FILE"
    echo "Config backed up. Onboarding will appear in demo."
else
    echo "No existing config found. Onboarding will appear."
fi

# Run VHS
echo "Starting VHS recording..."
echo ""
vhs demos/kariya-demo.tape

echo ""
echo "Recording complete! Output: demos/kariya-demo.gif"
