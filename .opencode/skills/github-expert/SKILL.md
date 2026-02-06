---
name: github-expert
description: GitHub Actions, workflows, CLI, API, and repository management best practices
license: MIT
compatibility: opencode
metadata:
  audience: developers
  workflow: kariya
---

# GitHub Expert Skill

You are a GitHub expert proficient in GitHub Actions, workflows, CLI, API, and best practices for repository management.

## Overview

GitHub is more than version control. Master Actions, CLI, and API to automate everything.

---

## GitHub CLI (gh)

### Installation & Auth

```bash
# Install
brew install gh  # macOS
sudo apt install gh  # Ubuntu

# Authenticate
gh auth login
gh auth status
```

### Common Commands

```bash
# Repository
gh repo clone owner/repo
gh repo create my-repo --public
gh repo view
gh repo fork

# Issues
gh issue list
gh issue create --title "Bug" --body "Description"
gh issue view 123
gh issue close 123

# Pull Requests
gh pr list
gh pr create --title "Feature" --body "Description"
gh pr view 123
gh pr checkout 123
gh pr merge 123
gh pr review 123 --approve

# Workflows
gh workflow list
gh workflow run ci.yml
gh run list
gh run view 12345
gh run watch 12345

# Releases
gh release list
gh release create v1.0.0 --generate-notes
gh release download v1.0.0
```

### Advanced Usage

```bash
# Create PR with reviewers
gh pr create \
  --title "feat: add feature" \
  --body "## Summary\n- Added feature" \
  --reviewer user1,user2 \
  --label enhancement

# Merge with specific method
gh pr merge 123 --squash --delete-branch

# View PR checks
gh pr checks 123

# Run workflow with inputs
gh workflow run deploy.yml -f environment=staging

# API calls
gh api repos/{owner}/{repo}/issues
gh api repos/{owner}/{repo}/pulls --jq '.[].title'
```

---

## GitHub Actions

### Workflow Syntax

```yaml
name: CI

on:
  push:
    branches: [main]
    paths:
      - '**.go'
      - 'go.mod'
  pull_request:
    branches: [main]
  workflow_dispatch:  # Manual trigger
    inputs:
      environment:
        description: 'Target environment'
        required: true
        default: 'staging'
        type: choice
        options:
          - staging
          - production

permissions:
  contents: read
  pull-requests: write

env:
  GO_VERSION: '1.22'

jobs:
  build:
    runs-on: ubuntu-latest
    timeout-minutes: 10
    
    steps:
      - uses: actions/checkout@v4
      
      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}
          cache: true
      
      - name: Build
        run: go build ./...
```

### Job Dependencies

```yaml
jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: make lint

  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: make test

  deploy:
    needs: [lint, test]  # Wait for both
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - run: make deploy
```

### Conditional Execution

```yaml
steps:
  # Only on push to main
  - if: github.event_name == 'push' && github.ref == 'refs/heads/main'
    run: deploy

  # Only on PR
  - if: github.event_name == 'pull_request'
    run: test

  # Only if file changed
  - uses: dorny/paths-filter@v3
    id: changes
    with:
      filters: |
        src:
          - 'src/**'
  
  - if: steps.changes.outputs.src == 'true'
    run: build
```

### Secrets & Variables

```yaml
jobs:
  deploy:
    environment: production  # Use environment secrets
    steps:
      - name: Deploy
        env:
          API_KEY: ${{ secrets.API_KEY }}
          APP_VERSION: ${{ vars.APP_VERSION }}
        run: ./deploy.sh
```

### Matrix Strategy

```yaml
jobs:
  test:
    strategy:
      fail-fast: false
      matrix:
        os: [ubuntu-latest, macos-latest]
        go: ['1.21', '1.22']
        exclude:
          - os: macos-latest
            go: '1.21'
    runs-on: ${{ matrix.os }}
    steps:
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
      - run: go test ./...
```

### Reusable Workflows

```yaml
# .github/workflows/reusable-test.yml
name: Reusable Test

on:
  workflow_call:
    inputs:
      go-version:
        required: true
        type: string
    secrets:
      token:
        required: true

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ inputs.go-version }}
      - run: go test ./...
```

```yaml
# .github/workflows/ci.yml
name: CI

on: [push]

jobs:
  call-test:
    uses: ./.github/workflows/reusable-test.yml
    with:
      go-version: '1.22'
    secrets:
      token: ${{ secrets.GITHUB_TOKEN }}
```

### Composite Actions

```yaml
# .github/actions/setup-go-build/action.yml
name: Setup Go Build
description: Setup Go and build

inputs:
  go-version:
    description: Go version
    required: true
    default: '1.22'

runs:
  using: composite
  steps:
    - uses: actions/setup-go@v5
      with:
        go-version: ${{ inputs.go-version }}
        cache: true
    
    - name: Download dependencies
      shell: bash
      run: go mod download
    
    - name: Build
      shell: bash
      run: go build ./...
```

---

## Branch Protection

### Settings

```yaml
# Via API or UI
branches:
  main:
    protection:
      required_status_checks:
        strict: true
        contexts:
          - lint
          - test
      required_pull_request_reviews:
        required_approving_review_count: 1
        dismiss_stale_reviews: true
      enforce_admins: true
      required_linear_history: true
```

---

## PR Templates

```markdown
<!-- .github/pull_request_template.md -->
## Summary

Brief description of changes.

## Type of Change

- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing

- [ ] Tests added/updated
- [ ] All tests pass locally

## Checklist

- [ ] Code follows style guidelines
- [ ] Self-reviewed
- [ ] Documentation updated
```

---

## GitHub API

### Using gh api

```bash
# Get repo info
gh api repos/{owner}/{repo}

# List PRs with jq
gh api repos/{owner}/{repo}/pulls --jq '.[].title'

# Create issue
gh api repos/{owner}/{repo}/issues -f title="Bug" -f body="Description"

# Add labels
gh api repos/{owner}/{repo}/issues/123/labels -f labels[]=bug

# Get PR reviews
gh api repos/{owner}/{repo}/pulls/123/reviews
```

### GraphQL

```bash
gh api graphql -f query='
  query {
    repository(owner: "owner", name: "repo") {
      pullRequests(last: 5, states: OPEN) {
        nodes {
          title
          number
          author { login }
        }
      }
    }
  }
'
```

---

## Dependabot

```yaml
# .github/dependabot.yml
version: 2
updates:
  - package-ecosystem: gomod
    directory: "/"
    schedule:
      interval: weekly
    commit-message:
      prefix: "deps"
    labels:
      - dependencies
    reviewers:
      - maintainer
    
  - package-ecosystem: github-actions
    directory: "/"
    schedule:
      interval: weekly
```

---

## Code Owners

```
# .github/CODEOWNERS
* @default-owner

/internal/cli/ @cli-team
/internal/domain/ @domain-team
*.md @docs-team

# Specific files
go.mod @maintainer
```

---

## Useful Patterns

### Auto-merge Dependabot

```yaml
name: Auto-merge Dependabot

on:
  pull_request:

permissions:
  contents: write
  pull-requests: write

jobs:
  auto-merge:
    if: github.actor == 'dependabot[bot]'
    runs-on: ubuntu-latest
    steps:
      - uses: dependabot/fetch-metadata@v2
        id: metadata
      
      - if: steps.metadata.outputs.update-type == 'version-update:semver-patch'
        run: gh pr merge --auto --squash "$PR_URL"
        env:
          PR_URL: ${{ github.event.pull_request.html_url }}
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### PR Comment on Build

```yaml
- name: Comment PR
  if: github.event_name == 'pull_request'
  uses: actions/github-script@v7
  with:
    script: |
      github.rest.issues.createComment({
        issue_number: context.issue.number,
        owner: context.repo.owner,
        repo: context.repo.repo,
        body: '✅ Build successful!'
      })
```

---

## Related Skills

- `devops` - CI/CD pipelines
- `automation` - Workflow automation
- `scripter` - Scripting in workflows
