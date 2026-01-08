#!/bin/bash

# Verify Git Hooks Installation Script
#
# Checks that all required git hooks are properly installed and executable

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

HOOKS_DIR=".git/hooks"
REQUIRED_HOOKS=("pre-commit" "commit-msg" "prepare-commit-msg")
MISSING_HOOKS=()
INVALID_HOOKS=()

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "🔍 Git Hooks Verification"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo ""

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo -e "${RED}❌ Error: Not in a git repository${NC}"
    exit 1
fi

# Check each required hook
for hook in "${REQUIRED_HOOKS[@]}"; do
    hook_path="$HOOKS_DIR/$hook"
    
    echo -n "Checking $hook: "
    
    if [ ! -f "$hook_path" ]; then
        echo -e "${RED}❌ Missing${NC}"
        MISSING_HOOKS+=("$hook")
    elif [ ! -x "$hook_path" ]; then
        echo -e "${YELLOW}⚠️  Not executable${NC}"
        INVALID_HOOKS+=("$hook")
    else
        # Verify hook has expected marker (basic validation)
        # We check for shebang to ensure it's a valid script
        if head -1 "$hook_path" | grep -q '^#!/'; then
            echo -e "${GREEN}✅ Installed${NC}"
        else
            echo -e "${YELLOW}⚠️  Invalid format${NC}"
            INVALID_HOOKS+=("$hook")
        fi
    fi
done

echo ""

# Report results
if [ ${#MISSING_HOOKS[@]} -eq 0 ] && [ ${#INVALID_HOOKS[@]} -eq 0 ]; then
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}✅ All git hooks are properly installed${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    exit 0
else
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${RED}❌ Git hooks installation incomplete${NC}"
    echo -e "${RED}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    if [ ${#MISSING_HOOKS[@]} -gt 0 ]; then
        echo "Missing hooks:"
        for hook in "${MISSING_HOOKS[@]}"; do
            echo "  - $hook"
        done
        echo ""
    fi
    
    if [ ${#INVALID_HOOKS[@]} -gt 0 ]; then
        echo "Invalid/non-executable hooks:"
        for hook in "${INVALID_HOOKS[@]}"; do
            echo "  - $hook"
        done
        echo ""
    fi
    
    echo -e "${BLUE}To install hooks, run:${NC}"
    echo "  make install-git-hooks"
    echo ""
    exit 1
fi
