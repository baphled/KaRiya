#!/bin/bash
# Comprehensive Escape Key Audit - All 10 Intents
# Checks for proper HandleGlobalKeys ordering in ALL update methods

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "========================================"
echo "Comprehensive Escape Key Audit"
echo "Checking ALL 10 intents for proper escape handling"
echo "========================================"
echo ""

INTENTS_DIR="internal/cli/intents"
ISSUES_FOUND=0
INTENTS_CHECKED=0

# Function to check a specific intent file
check_intent() {
    local file=$1
    local intent_name=$(basename "$file" _intent.go)
    
    INTENTS_CHECKED=$((INTENTS_CHECKED + 1))
    
    echo "----------------------------------------"
    echo "Intent: $intent_name"
    echo "File: $file"
    echo "----------------------------------------"
    
    # Find all update methods (update* functions)
    update_methods=$(grep -n "^func (.*) update[A-Z]" "$file" | cut -d: -f1,2 | sed 's/:/ - Line /')
    
    if [ -z "$update_methods" ]; then
        echo -e "${YELLOW}⚠️  No update methods found${NC}"
        echo ""
        return
    fi
    
    echo "Update methods found:"
    echo "$update_methods"
    echo ""
    
    # Check each update method for delegation patterns
    while IFS= read -r method_line; do
        line_num=$(echo "$method_line" | awk '{print $3}')
        method_name=$(sed -n "${line_num}p" "$file" | sed 's/func (.*) \(update[A-Za-z]*\).*/\1/')
        
        # Get the method body (from line to next method or end)
        next_method=$(grep -n "^func " "$file" | awk -F: -v line="$line_num" '$1 > line {print $1; exit}')
        if [ -z "$next_method" ]; then
            next_method=$(wc -l < "$file")
        fi
        
        # Extract method body
        method_body=$(sed -n "${line_num},${next_method}p" "$file")
        
        # Check for .Update(msg) delegation
        delegation_lines=$(echo "$method_body" | grep -n "\.Update(msg)" | head -5)
        
        if [ -n "$delegation_lines" ]; then
            # Check if HandleGlobalKeys appears BEFORE first delegation
            first_delegation_line=$(echo "$delegation_lines" | head -1 | cut -d: -f1)
            global_keys_line=$(echo "$method_body" | grep -n "HandleGlobalKeys" | head -1 | cut -d: -f1)
            
            if [ -z "$global_keys_line" ]; then
                echo -e "${RED}❌ $method_name (line $line_num): .Update(msg) delegation found, NO HandleGlobalKeys${NC}"
                echo "   Delegation at relative line: $first_delegation_line"
                ISSUES_FOUND=$((ISSUES_FOUND + 1))
            elif [ "$global_keys_line" -gt "$first_delegation_line" ]; then
                echo -e "${RED}❌ $method_name (line $line_num): HandleGlobalKeys AFTER delegation${NC}"
                echo "   Delegation at relative line: $first_delegation_line"
                echo "   HandleGlobalKeys at relative line: $global_keys_line"
                ISSUES_FOUND=$((ISSUES_FOUND + 1))
            else
                echo -e "${GREEN}✅ $method_name (line $line_num): HandleGlobalKeys BEFORE delegation${NC}"
            fi
        else
            # No delegation - check if HandleGlobalKeys is present anyway
            if echo "$method_body" | grep -q "HandleGlobalKeys"; then
                echo -e "${GREEN}✅ $method_name (line $line_num): Has HandleGlobalKeys, no delegation${NC}"
            else
                echo -e "${BLUE}ℹ️  $method_name (line $line_num): No delegation, no HandleGlobalKeys (may be leaf state)${NC}"
            fi
        fi
        
    done <<< "$update_methods"
    
    echo ""
}

# Check all 10 intent files
for intent_file in \
    "$INTENTS_DIR/browse_timeline_intent.go" \
    "$INTENTS_DIR/bulk_operations_intent.go" \
    "$INTENTS_DIR/burst_management_intent.go" \
    "$INTENTS_DIR/capture_event_intent.go" \
    "$INTENTS_DIR/configure_system_intent.go" \
    "$INTENTS_DIR/export_artifact_intent.go" \
    "$INTENTS_DIR/fact_management_intent.go" \
    "$INTENTS_DIR/generate_cv_intent.go" \
    "$INTENTS_DIR/import_wizard_intent.go" \
    "$INTENTS_DIR/metadata_editor_intent.go"
do
    if [ -f "$intent_file" ]; then
        check_intent "$intent_file"
    else
        echo "⚠️  File not found: $intent_file"
        echo ""
    fi
done

echo "========================================"
echo "Audit Complete"
echo "========================================"
echo "Intents checked: $INTENTS_CHECKED"
echo "Issues found: $ISSUES_FOUND"

if [ $ISSUES_FOUND -eq 0 ]; then
    echo -e "${GREEN}✅ ALL INTENTS PASS - No escape handling issues found!${NC}"
    exit 0
else
    echo -e "${RED}❌ ISSUES FOUND - $ISSUES_FOUND methods need fixing${NC}"
    exit 1
fi
