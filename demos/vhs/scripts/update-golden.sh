#!/bin/bash
# Update golden file baselines from current screenshots
#
# Usage: ./demos/vhs/scripts/update-golden.sh
#
# This copies all current screenshots to the golden directory.
# Review the changes before committing!

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VHS_DIR="$(dirname "$SCRIPT_DIR")"
GOLDEN_DIR="$VHS_DIR/golden"
GENERATED_DIR="$VHS_DIR/output/screenshots"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo "=== Updating Golden File Baselines ==="
echo ""

# Check if screenshots exist
if [ ! -d "$GENERATED_DIR" ] || [ -z "$(ls -A "$GENERATED_DIR" 2>/dev/null)" ]; then
    echo -e "${YELLOW}No screenshots found in $GENERATED_DIR${NC}"
    echo ""
    echo "Generate screenshots first with:"
    echo "  vhs demos/vhs/golden-test.tape"
    exit 1
fi

# Create golden directory if needed
mkdir -p "$GOLDEN_DIR"

# Copy screenshots
count=0
for screenshot in "$GENERATED_DIR"/*.png; do
    [ -e "$screenshot" ] || continue
    
    filename=$(basename "$screenshot")
    cp "$screenshot" "$GOLDEN_DIR/$filename"
    echo -e "${GREEN}Updated:${NC} $filename"
    ((count++))
done

echo ""
echo "=== Summary ==="
echo -e "${GREEN}$count golden file(s) updated${NC}"
echo ""
echo "Next steps:"
echo "  1. Review the changes: git diff demos/vhs/golden/"
echo "  2. If correct, commit: git add demos/vhs/golden/ && git commit -m 'Update golden files'"
