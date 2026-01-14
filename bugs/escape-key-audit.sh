#!/bin/bash
# Escape Key Audit Script
# Identifies intents with potential escape key handling issues

echo "==================================="
echo "Escape Key Navigation Audit"
echo "==================================="
echo ""

# Function to check if HandleGlobalKeys is called BEFORE delegation
check_intent() {
    local file=$1
    local intent_name=$(basename "$file" | sed 's/_intent.go//')
    
    echo "Checking: $intent_name"
    echo "File: $file"
    
    # Look for Update(msg) delegation patterns
    local delegations=$(grep -n "\.Update(msg)" "$file" | grep -v "// " | head -5)
    
    if [ -n "$delegations" ]; then
        echo "  Found message delegations:"
        echo "$delegations" | while read line; do
            line_num=$(echo "$line" | cut -d: -f1)
            echo "    Line $line_num: $(echo "$line" | cut -d: -f2-)"
            
            # Check if HandleGlobalKeys appears BEFORE this line in the same function
            # Extract function name
            func_start=$(awk -v n="$line_num" 'NR < n && /^func/ {line=NR; text=$0} END {print line":"text}' "$file")
            func_line=$(echo "$func_start" | cut -d: -f1)
            
            if [ -n "$func_line" ]; then
                # Check if HandleGlobalKeys appears between function start and delegation
                global_keys=$(sed -n "${func_line},${line_num}p" "$file" | grep -n "HandleGlobalKeys")
                
                if [ -n "$global_keys" ]; then
                    global_line=$(echo "$global_keys" | head -1 | cut -d: -f1)
                    global_line=$((func_line + global_line - 1))
                    
                    if [ "$global_line" -lt "$line_num" ]; then
                        echo "      ✅ HandleGlobalKeys checked BEFORE delegation (line $global_line)"
                    else
                        echo "      ❌ HandleGlobalKeys checked AFTER delegation (line $global_line)"
                    fi
                else
                    echo "      ❌ NO HandleGlobalKeys check found in function"
                fi
            fi
        done
    else
        echo "  ✅ No message delegation found (or handles messages directly)"
    fi
    
    echo ""
}

# Check all intent files
for intent_file in internal/cli/intents/*_intent.go; do
    if [ -f "$intent_file" ]; then
        check_intent "$intent_file"
    fi
done

echo "==================================="
echo "Audit Complete"
echo "==================================="
