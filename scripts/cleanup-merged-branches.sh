#!/bin/bash
# Script to clean up local branches that have been merged

set -e

# Colours for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Colour

# Get the default branch (usually main or master)
DEFAULT_BRANCH=$(git symbolic-ref refs/remotes/origin/HEAD 2>/dev/null | sed 's@^refs/remotes/origin/@@' || echo "main")

echo -e "${YELLOW}Fetching latest changes from remote...${NC}"
git fetch --prune

echo -e "\n${YELLOW}Finding merged branches...${NC}"

# Get list of local branches that have been merged into the default branch
# Exclude the default branch itself and current branch
MERGED_BRANCHES=$(git branch --merged "$DEFAULT_BRANCH" | \
    grep -v "^\*" | \
    grep -v "^  $DEFAULT_BRANCH$" | \
    grep -v "^  next$" | \
    sed 's/^[ *]*//')

if [ -z "$MERGED_BRANCHES" ]; then
    echo -e "${GREEN}No merged branches to clean up!${NC}"
    exit 0
fi

echo -e "\n${YELLOW}The following branches have been merged and will be deleted:${NC}"
echo "$MERGED_BRANCHES"

# Ask for confirmation unless --force flag is used
if [[ "$1" != "--force" && "$1" != "-f" ]]; then
    echo -e "\n${RED}Do you want to delete these branches? (y/N)${NC}"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        echo -e "${YELLOW}Cancelled.${NC}"
        exit 0
    fi
fi

# Delete the merged branches
echo -e "\n${YELLOW}Deleting merged branches...${NC}"
echo "$MERGED_BRANCHES" | while read -r branch; do
    if [ -n "$branch" ]; then
        git branch -d "$branch" && echo -e "${GREEN}Deleted: $branch${NC}" || echo -e "${RED}Failed to delete: $branch${NC}"
    fi
done

echo -e "\n${GREEN}Branch cleanup complete!${NC}"
