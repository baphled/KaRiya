---
description: Create a new repository with proper patterns
agent: build
---

Create a new repository following KaRiya database patterns.

Load these skills:
- `db-operations` - Repository patterns and GORM usage
- `architecture` - Layer placement
- `security` - SQL injection prevention
- `clean-code` - Clean implementation

## Repository
$ARGUMENTS

## Process

1. **Define Interface** in `internal/repository/career/{name}_repository.go`
   - Add go:generate directive for mock
   - Define repository interface
   - Define filter struct
   - Define error variables

2. **Create GORM Model** in `internal/repository/models/career.go`
   - Add model struct with gorm tags
   - Add TableName() method
   - Add FromDomain() function
   - Add ToDomain() method

3. **Create Migration** in `internal/repository/career/migrations/`
   - Use sequential numbering (007_create_things.sql)
   - Use IF NOT EXISTS
   - Add appropriate indexes

4. **Implement SQL Repository** in `internal/repository/career/sql/`
   - Implement all interface methods
   - Use parameterized queries
   - Wrap errors appropriately

5. **Implement Memory Repository** in `internal/repository/career/memory/`
   - For testing purposes

6. **Generate Mocks**
   ```bash
   make generate-mocks
   ```

7. **Write Tests** for SQL repository

## Requirements

- Interface in career/, implementation in career/sql/
- All queries use parameters (never string concat)
- Proper error wrapping
- Memory implementation for testing
- >= 95% test coverage
