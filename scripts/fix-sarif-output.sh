#!/bin/bash
# Fix invalid SARIF output from gosec security scanner
#
# SARIF spec requires relationships to be objects with defined structure,
# but gosec outputs them as arrays. This script fixes the format for
# compatibility with GitHub's SARIF upload action.

set -euo pipefail

SARIF_FILE="${1:-gosec.sarif}"

if [ ! -f "$SARIF_FILE" ]; then
    echo "SARIF file not found: $SARIF_FILE"
    exit 1
fi

# Fix invalid relationships array in gosec SARIF output
# gosec outputs relationships as arrays but SARIF spec requires objects
jq 'walk(if type == "object" and has("relationships") then 
  .relationships = [] 
else . end)' "$SARIF_FILE" > "${SARIF_FILE}.tmp"

mv "${SARIF_FILE}.tmp" "$SARIF_FILE"

echo "✓ Fixed SARIF relationships format in $SARIF_FILE"
