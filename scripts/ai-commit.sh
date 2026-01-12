#!/bin/bash

set -e

# ============================================================================
# AI Commit Helper
# ============================================================================
# Automates AI-attributed commits by:
# 1. Validating commit message format
# 2. Checking for staged changes
# 3. Adding AI attribution and human review trailers
# 4. Creating the commit
#
# Usage: make ai-commit MSG="feat(scope): description"
# ============================================================================

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Get commit message from first argument
COMMIT_MSG="$1"

# ============================================================================
# Step 1: Validate commit message provided
# ============================================================================

if [ -z "$COMMIT_MSG" ]; then
    echo -e "${RED}❌ ERROR: Commit message required${NC}"
    echo ""
    echo "Usage:"
    echo "  make ai-commit MSG=\"feat(scope): description\""
    echo ""
    echo "Examples:"
    echo "  make ai-commit MSG=\"feat(forms): add date validation helpers\""
    echo "  make ai-commit MSG=\"fix(tests): resolve race condition in burst tests\""
    echo "  make ai-commit MSG=\"docs(readme): update installation steps\""
    echo ""
    exit 1
fi

# ============================================================================
# Step 2: Check for staged changes
# ============================================================================

echo ""
echo -e "${BLUE}🔍 Checking for staged changes...${NC}"

if git diff --cached --quiet; then
    echo -e "${RED}❌ ERROR: No staged changes${NC}"
    echo ""
    echo "You must stage changes before committing:"
    echo "  git add -p <file>          # Stage specific hunks interactively"
    echo "  git add <file>             # Stage entire file"
    echo ""
    echo "Then try again:"
    echo "  make ai-commit MSG=\"${COMMIT_MSG}\""
    echo ""
    exit 1
fi

echo -e "${GREEN}✅ Staged changes detected${NC}"

# ============================================================================
# Step 3: Validate commit message format (via commitlint)
# ============================================================================

echo ""
echo -e "${BLUE}🔍 Validating commit message format...${NC}"

# Create temporary file for commit message
TEMP_MSG_FILE=$(mktemp)
echo "$COMMIT_MSG" > "$TEMP_MSG_FILE"

# Check if commitlint is available
if command -v npx &> /dev/null; then
    if ! npx commitlint --edit "$TEMP_MSG_FILE" 2>&1; then
        rm -f "$TEMP_MSG_FILE"
        echo ""
        echo -e "${RED}❌ Commit message does not follow conventional commit format${NC}"
        echo ""
        echo "Required format:"
        echo -e "${GREEN}  type(scope): subject${NC}"
        echo ""
        echo "Examples:"
        echo "  feat(cli): add new command"
        echo "  fix(tests): resolve race condition"
        echo "  docs(readme): update installation steps"
        echo ""
        echo "Allowed types: feat, fix, docs, style, refactor, test, chore, perf, ci, build, revert"
        echo ""
        echo "See: docs/rules/COMMIT_QUICK_REFERENCE.md"
        exit 1
    fi
    echo -e "${GREEN}✅ Commit message format valid${NC}"
else
    echo -e "${YELLOW}⚠️  commitlint not found - skipping format validation${NC}"
    echo "Run: npm install"
fi

rm -f "$TEMP_MSG_FILE"

# ============================================================================
# Step 4: Get AI agent and model information
# ============================================================================

# Use environment variables with sensible defaults
AGENT_NAME="${AI_AGENT:-OpenCode}"
MODEL_NAME="${AI_MODEL:-Claude Sonnet 4.5}"

# Get reviewer name from git config
REVIEWER_NAME=$(git config user.name)

if [ -z "$REVIEWER_NAME" ]; then
    echo -e "${YELLOW}⚠️  Warning: git user.name not set${NC}"
    echo "Set it with: git config user.name \"Your Name\""
    REVIEWER_NAME="Unknown"
fi

# ============================================================================
# Step 5: Create commit with AI attribution
# ============================================================================

echo ""
echo -e "${BLUE}🤖 Creating AI-attributed commit...${NC}"
echo ""
echo "Agent:    ${AGENT_NAME}"
echo "Model:    ${MODEL_NAME}"
echo "Reviewer: ${REVIEWER_NAME}"
echo ""

# Build full commit message with attribution
# Use temporary file to handle multi-line messages properly
COMMIT_MSG_FILE=$(mktemp)

# Write commit message to temp file
# This preserves newlines and formatting
cat > "$COMMIT_MSG_FILE" << EOF
${COMMIT_MSG}

AI-Generated-By: ${AGENT_NAME} (${MODEL_NAME})
Reviewed-By: ${REVIEWER_NAME}
EOF

# Create the commit using the temp file
if git commit -F "$COMMIT_MSG_FILE"; then
    echo ""
    echo -e "${GREEN}✅ Commit created successfully${NC}"
    echo ""
    echo "Commit message:"
    echo "─────────────────────────────────────────────"
    git log -1 --pretty=%B
    echo "─────────────────────────────────────────────"
    echo ""
    
    # Clean up temp file
    rm -f "$COMMIT_MSG_FILE"
else
    echo ""
    echo -e "${RED}❌ Commit failed${NC}"
    rm -f "$COMMIT_MSG_FILE"
    exit 1
fi

# ============================================================================
# Step 6: Summary
# ============================================================================

echo -e "${GREEN}✅ AI-attributed commit complete${NC}"
echo ""
echo "Next steps:"
echo "  git log -1                  # Review the commit"
echo "  make check-compliance       # Run compliance checks"
echo "  git push                    # Push to remote (when ready)"
echo ""
