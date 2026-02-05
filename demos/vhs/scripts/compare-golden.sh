#!/bin/bash
# Compare generated screenshots against golden file baselines
#
# Usage: ./demos/vhs/scripts/compare-golden.sh
#
# Prerequisites:
#   - ImageMagick (for compare command)
#   - Golden baselines in demos/vhs/golden/
#   - Generated screenshots in demos/vhs/output/screenshots/
#
# Exit codes:
#   0 - All screenshots match (or within threshold)
#   1 - Visual differences detected

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
VHS_DIR="$(dirname "$SCRIPT_DIR")"
GOLDEN_DIR="$VHS_DIR/golden"
GENERATED_DIR="$VHS_DIR/output/screenshots"
DIFF_DIR="$VHS_DIR/output/diff"

# Pixel difference threshold (allow minor anti-aliasing differences)
THRESHOLD=100

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check dependencies
if ! command -v compare &> /dev/null; then
    echo -e "${RED}Error: ImageMagick 'compare' command not found${NC}"
    echo "Install with: brew install imagemagick (macOS) or apt install imagemagick (Linux)"
    exit 1
fi

# Create diff directory
mkdir -p "$DIFF_DIR"

# Track results
PASS=0
FAIL=0
MISSING=0
NEW=0

echo "=== Golden File Comparison ==="
echo "Golden dir: $GOLDEN_DIR"
echo "Generated dir: $GENERATED_DIR"
echo "Threshold: $THRESHOLD pixels"
echo ""

# Check each golden file
for golden in "$GOLDEN_DIR"/*.png 2>/dev/null; do
    [ -e "$golden" ] || continue
    
    filename=$(basename "$golden")
    generated="$GENERATED_DIR/$filename"
    
    if [ ! -f "$generated" ]; then
        echo -e "${YELLOW}MISSING:${NC} $filename (no generated screenshot)"
        ((MISSING++))
        continue
    fi
    
    # Compare using ImageMagick
    diff_result=$(compare -metric AE "$golden" "$generated" "$DIFF_DIR/$filename" 2>&1) || true
    
    # Extract numeric value (compare outputs to stderr)
    diff_pixels=$(echo "$diff_result" | grep -oE '^[0-9]+' || echo "0")
    
    if [ "$diff_pixels" -gt "$THRESHOLD" ]; then
        echo -e "${RED}CHANGED:${NC} $filename (diff: $diff_pixels pixels)"
        ((FAIL++))
    else
        echo -e "${GREEN}OK:${NC} $filename"
        ((PASS++))
        # Clean up diff file if passed
        rm -f "$DIFF_DIR/$filename"
    fi
done

# Check for new screenshots not in golden
for generated in "$GENERATED_DIR"/*.png 2>/dev/null; do
    [ -e "$generated" ] || continue
    
    filename=$(basename "$generated")
    golden="$GOLDEN_DIR/$filename"
    
    if [ ! -f "$golden" ]; then
        echo -e "${YELLOW}NEW:${NC} $filename (no golden baseline)"
        ((NEW++))
    fi
done

# Summary
echo ""
echo "=== Summary ==="
echo -e "${GREEN}Passed:${NC} $PASS"
echo -e "${RED}Failed:${NC} $FAIL"
echo -e "${YELLOW}Missing:${NC} $MISSING"
echo -e "${YELLOW}New:${NC} $NEW"

if [ "$FAIL" -gt 0 ]; then
    echo ""
    echo -e "${RED}Visual regression detected!${NC}"
    echo "Diff images saved to: $DIFF_DIR"
    echo ""
    echo "To update golden files, run:"
    echo "  ./demos/vhs/scripts/update-golden.sh"
    exit 1
fi

if [ "$MISSING" -gt 0 ] || [ "$NEW" -gt 0 ]; then
    echo ""
    echo -e "${YELLOW}Warning: Golden file set is incomplete${NC}"
    echo "Run update-golden.sh to sync baseline files"
fi

echo ""
echo -e "${GREEN}All visual tests passed!${NC}"
exit 0
