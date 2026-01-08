#!/bin/bash

# Install Git Hooks for AI Commit Attribution
#
# This script installs git hooks that enforce AI commit attribution rules

set -e

echo "================================================"
echo "🔧 Installing Git Hooks for AI Attribution"
echo "================================================"
echo ""

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if we're in a git repository
if [ ! -d ".git" ]; then
    echo "❌ Error: Not in a git repository"
    exit 1
fi

# Create hooks directory if it doesn't exist
mkdir -p .git/hooks

echo "Installing hooks..."
echo ""

# Install pre-commit hook
if [ -f ".git/hooks/pre-commit" ]; then
    echo -e "${YELLOW}⚠️  pre-commit hook already exists${NC}"
    echo -n "Overwrite? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Skipping pre-commit..."
    else
        cp .git-hooks/pre-commit .git/hooks/pre-commit
        chmod +x .git/hooks/pre-commit
        echo -e "${GREEN}✅ Installed pre-commit hook${NC}"
    fi
else
    cp .git-hooks/pre-commit .git/hooks/pre-commit
    chmod +x .git/hooks/pre-commit
    echo -e "${GREEN}✅ Installed pre-commit hook${NC}"
fi

echo ""

# Install prepare-commit-msg hook
if [ -f ".git/hooks/prepare-commit-msg" ]; then
    echo -e "${YELLOW}⚠️  prepare-commit-msg hook already exists${NC}"
    echo -n "Overwrite? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Skipping prepare-commit-msg..."
    else
        cp .git-hooks/prepare-commit-msg .git/hooks/prepare-commit-msg
        chmod +x .git/hooks/prepare-commit-msg
        echo -e "${GREEN}✅ Installed prepare-commit-msg hook${NC}"
    fi
else
    cp .git-hooks/prepare-commit-msg .git/hooks/prepare-commit-msg
    chmod +x .git/hooks/prepare-commit-msg
    echo -e "${GREEN}✅ Installed prepare-commit-msg hook${NC}"
fi

echo ""

# Install commit-msg hook
if [ -f ".git/hooks/commit-msg" ]; then
    echo -e "${YELLOW}⚠️  commit-msg hook already exists${NC}"
    echo -n "Overwrite? (y/N): "
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo "Skipping commit-msg..."
    else
        cp .git-hooks/commit-msg .git/hooks/commit-msg
        chmod +x .git/hooks/commit-msg
        echo -e "${GREEN}✅ Installed commit-msg hook${NC}"
    fi
else
    cp .git-hooks/commit-msg .git/hooks/commit-msg
    chmod +x .git/hooks/commit-msg
    echo -e "${GREEN}✅ Installed commit-msg hook${NC}"
fi

echo ""

# Install commit-msg-lint hook (optional, requires Node.js)
if command -v node &> /dev/null && [ -f ".git-hooks/commit-msg-lint" ]; then
    echo "Installing commitlint hook (optional)..."
    if [ -f ".git/hooks/commit-msg-lint" ]; then
        echo -e "${YELLOW}⚠️  commit-msg-lint hook already exists${NC}"
        echo -n "Overwrite? (y/N): "
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            cp .git-hooks/commit-msg-lint .git/hooks/commit-msg-lint
            chmod +x .git/hooks/commit-msg-lint
            echo -e "${GREEN}✅ Installed commit-msg-lint hook${NC}"
        fi
    else
        cp .git-hooks/commit-msg-lint .git/hooks/commit-msg-lint
        chmod +x .git/hooks/commit-msg-lint
        echo -e "${GREEN}✅ Installed commit-msg-lint hook${NC}"
    fi
    echo ""
else
    if ! command -v node &> /dev/null; then
        echo -e "${YELLOW}⚠️  Node.js not found. Skipping commitlint hook.${NC}"
        echo "   Install Node.js and run 'npm install' to enable commitlint."
    fi
    echo ""
fi

# Install Node.js dependencies if package.json exists
if [ -f "package.json" ] && command -v npm &> /dev/null; then
    echo "Installing Node.js dependencies for commitlint..."
    if npm ci 2>/dev/null || npm install 2>/dev/null; then
        echo -e "${GREEN}✅ Installed Node.js dependencies${NC}"
    else
        echo -e "${YELLOW}⚠️  Could not install Node.js dependencies${NC}"
    fi
    echo ""
fi

# Configure git commit template
echo "Setting up commit message template..."
if git config commit.template .gitmessage 2>/dev/null; then
    echo -e "${GREEN}✅ Configured commit template${NC}"
else
    echo -e "${YELLOW}⚠️  Could not configure commit template${NC}"
fi

echo ""
echo "================================================"
echo "✅ Installation Complete"
echo "================================================"
echo ""
echo "Installed hooks:"
echo "  - pre-commit: Code quality and TDD enforcement"
echo "  - prepare-commit-msg: Adds AI attribution reminder"
echo "  - commit-msg: Validates AI attribution format"
if command -v node &> /dev/null && [ -f ".git/hooks/commit-msg-lint" ]; then
    echo "  - commit-msg-lint: Validates conventional commits format"
fi
echo ""
echo "Configured:"
echo "  - Commit message template (.gitmessage)"
if [ -d "node_modules" ]; then
    echo "  - Commitlint (conventional commits validation)"
fi
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo "1. Read: docs/rules/AI_COMMIT_ATTRIBUTION.md"
echo "2. When committing AI-generated code, include:"
echo "   AI-Generated-By: <Assistant> (<Model>)"
echo "   Reviewed-By: <Your Name>"
echo "3. Follow conventional commits format:"
echo "   <type>(<scope>): <subject>"
echo ""
echo "To bypass hooks (for emergencies only):"
echo "  git commit --no-verify"
echo ""
echo "To uninstall hooks:"
echo "  rm .git/hooks/prepare-commit-msg"
echo "  rm .git/hooks/commit-msg"
if [ -f ".git/hooks/commit-msg-lint" ]; then
    echo "  rm .git/hooks/commit-msg-lint"
fi
echo ""

