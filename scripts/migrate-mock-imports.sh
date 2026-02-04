#!/bin/bash
set -e

if sed --version 2>/dev/null | grep -q 'GNU'; then
    SED_INPLACE=(sed -i)
else
    SED_INPLACE=(sed -i '')
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

BACKUP_DIR="$PROJECT_ROOT/.import-backup-$(date +%Y%m%d%H%M%S)"
mkdir -p "$BACKUP_DIR"
echo "Backing up test files to: $BACKUP_DIR"
echo ""

CHANGED=0

while read -r file; do
    relative_path="${file#$PROJECT_ROOT/}"

    if grep -q 'github.com/baphled/kariya/internal/testutil/mocks"' "$file" ||
       grep -q 'github.com/baphled/kariya/internal/repository/career/mocks"' "$file"; then

        mkdir -p "$BACKUP_DIR/$(dirname "$relative_path")"
        cp "$file" "$BACKUP_DIR/$relative_path"

        echo "Processing: $relative_path"
        CHANGED=$((CHANGED + 1))

        # Replace old top-level mocks import with repository layer
        "${SED_INPLACE[@]}" 's|"github.com/baphled/kariya/internal/testutil/mocks"|mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"|g' "$file"

        # Replace old repository/career/mocks import
        "${SED_INPLACE[@]}" 's|"github.com/baphled/kariya/internal/repository/career/mocks"|mockrepo "github.com/baphled/kariya/internal/testutil/mocks/repository"|g' "$file"

        # Update constructor calls: mocks.NewMock* -> mockrepo.NewMock*
        "${SED_INPLACE[@]}" 's|\bmocks\.NewMockEventRepository\b|mockrepo.NewMockEventRepository|g' "$file"
        "${SED_INPLACE[@]}" 's|\bmocks\.NewMockBurstRepository\b|mockrepo.NewMockBurstRepository|g' "$file"
        "${SED_INPLACE[@]}" 's|\bmocks\.NewMockFactRepository\b|mockrepo.NewMockFactRepository|g' "$file"
        "${SED_INPLACE[@]}" 's|\bmocks\.NewMockSkillRepository\b|mockrepo.NewMockSkillRepository|g' "$file"

        # Update recorder references
        "${SED_INPLACE[@]}" 's|\bmocks\.MockEventRepository\b|mockrepo.MockEventRepository|g' "$file"
        "${SED_INPLACE[@]}" 's|\bmocks\.MockBurstRepository\b|mockrepo.MockBurstRepository|g' "$file"
        "${SED_INPLACE[@]}" 's|\bmocks\.MockFactRepository\b|mockrepo.MockFactRepository|g' "$file"
        "${SED_INPLACE[@]}" 's|\bmocks\.MockSkillRepository\b|mockrepo.MockSkillRepository|g' "$file"

        # Replace hand-written mock references (from testutil/mocks package)
        "${SED_INPLACE[@]}" 's|\bmocks\.NewBurstServiceMock\b|mockintent.NewMockBurstService|g' "$file"
        "${SED_INPLACE[@]}" 's|\bmocks\.NewBurstRepositoryMock\b|mockrepo.NewMockBurstRepository|g' "$file"

        # Add mockintent import if mockintent references were introduced
        if grep -q 'mockintent\.' "$file" && ! grep -q 'mocks/intent"' "$file"; then
            "${SED_INPLACE[@]}" '/^import (/a\\tmockintent "github.com/baphled/kariya/internal/testutil/mocks/intent"' "$file"
        fi

        # Add mocksvc import if mocksvc references were introduced
        if grep -q 'mocksvc\.' "$file" && ! grep -q 'mocks/service"' "$file"; then
            "${SED_INPLACE[@]}" '/^import (/a\\tmocksvc "github.com/baphled/kariya/internal/testutil/mocks/service"' "$file"
        fi
    fi
done < <(find "$PROJECT_ROOT/internal" -name "*_test.go" -type f)

echo ""
echo "Migration complete! $CHANGED file(s) updated."
echo "Backup location: $BACKUP_DIR"
echo ""
echo "Next steps:"
echo "  1. Review changes: git diff"
echo "  2. Run tests: make test"
echo "  3. If issues arise, restore: cp -r $BACKUP_DIR/* $PROJECT_ROOT/"
