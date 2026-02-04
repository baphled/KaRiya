#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

BACKUP_DIR="$PROJECT_ROOT/.import-backup-$(date +%Y%m%d%H%M%S)"
mkdir -p "$BACKUP_DIR"
echo "Backing up test files to: $BACKUP_DIR"
echo ""

CHANGED=0

find "$PROJECT_ROOT/internal" -name "*_test.go" -type f | while read -r file; do
    relative_path="${file#$PROJECT_ROOT/}"

    if grep -q 'github.com/baphled/kariya/internal/testutil/mocks"' "$file" ||
       grep -q 'github.com/baphled/kariya/internal/repository/career/mocks"' "$file"; then

        mkdir -p "$BACKUP_DIR/$(dirname "$relative_path")"
        cp "$file" "$BACKUP_DIR/$relative_path"

        echo "Processing: $relative_path"
        CHANGED=$((CHANGED + 1))

        # Replace old top-level mocks import with repository layer
        sed -i 's|"github.com/baphled/kariya/internal/testutil/mocks"|mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"|g' "$file"

        # Replace old repository/career/mocks import
        sed -i 's|"github.com/baphled/kariya/internal/repository/career/mocks"|mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"|g' "$file"

        # Update constructor calls: mocks.NewMock* -> mockrepo.NewMock*
        sed -i 's|\bmocks\.NewMockEventRepository\b|mockrepo.NewMockEventRepository|g' "$file"
        sed -i 's|\bmocks\.NewMockBurstRepository\b|mockrepo.NewMockBurstRepository|g' "$file"
        sed -i 's|\bmocks\.NewMockFactRepository\b|mockrepo.NewMockFactRepository|g' "$file"
        sed -i 's|\bmocks\.NewMockSkillRepository\b|mockrepo.NewMockSkillRepository|g' "$file"

        # Update recorder references
        sed -i 's|\bmocks\.MockEventRepository\b|mockrepo.MockEventRepository|g' "$file"
        sed -i 's|\bmocks\.MockBurstRepository\b|mockrepo.MockBurstRepository|g' "$file"
        sed -i 's|\bmocks\.MockFactRepository\b|mockrepo.MockFactRepository|g' "$file"
        sed -i 's|\bmocks\.MockSkillRepository\b|mockrepo.MockSkillRepository|g' "$file"

        # Replace hand-written mock references (from testutil/mocks package)
        sed -i 's|\bmocks\.NewBurstServiceMock\b|mockintent.NewMockBurstService|g' "$file"
        sed -i 's|\bmocks\.NewBurstRepositoryMock\b|mockrepo.NewMockBurstRepository|g' "$file"
    fi
done

echo ""
echo "Migration complete!"
echo "Backup location: $BACKUP_DIR"
echo ""
echo "Next steps:"
echo "  1. Review changes: git diff"
echo "  2. Run tests: make test"
echo "  3. If issues arise, restore: cp -r $BACKUP_DIR/* $PROJECT_ROOT/"
