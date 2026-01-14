---
created: 2026-01-14T02:45
modified: 2026-01-14T02:45
---
# Editor & Tool Setup

This directory contains editor-specific setup guides and configurations for development tools used in the KaRiya project.

## 📁 Contents

### Neovim / Neotest Setup

- **[NEOTEST_SETUP.md](NEOTEST_SETUP.md)** - Complete Neotest setup guide for Ginkgo tests with automatic coverage generation

## 🚀 Quick Start

### For Neovim Users

Follow the [NEOTEST_SETUP.md](NEOTEST_SETUP.md) guide to:
1. Configure neotest-go adapter for Ginkgo tests
2. Enable automatic coverage generation
3. Set up keybindings and workflows
4. Troubleshoot common issues

### For Other Editors

This directory currently focuses on Neovim. For other editors:
- **VSCode**: Consider using the Go extension with built-in test support
- **IntelliJ/GoLand**: Built-in Go test runner with coverage
- **Vim**: Similar setup to Neovim with appropriate plugins

## 📖 Available Guides

### NEOTEST_SETUP.md

**Purpose**: Configure Neovim for running Ginkgo tests with automatic coverage

**Key Features**:
- Automatic `coverage.out` generation on every test run
- Race condition detection enabled
- Integration with neotest UI
- Coverage visualization support
- Troubleshooting guide

**Topics Covered**:
- Plugin installation and configuration
- Keybinding setup
- Usage examples
- Coverage viewing workflow
- Common issues and solutions

## 🎯 What You'll Learn

### Testing Workflow
- Run tests from within Neovim
- View test results in neotest panel
- Generate coverage automatically
- Visualize coverage in editor

### Configuration
- neotest-go adapter setup
- Coverage generation flags
- Test discovery settings
- Output configuration

### Troubleshooting
- Coverage not generated
- Tests not running
- Output not displaying
- Plugin conflicts

## 💡 Benefits

### Neotest-go Benefits
1. **Automatic Coverage**: Every test run generates `coverage.out`
2. **Standard Go Tools**: Uses `go test` natively
3. **Simple Configuration**: Minimal setup required
4. **Better Integration**: Works with Go tooling

### Trade-offs
- Individual Ginkgo spec detection not available
- Must run entire test files (not individual `It` blocks)
- Use CLI for focused Ginkgo testing

## 🔧 Prerequisites

### Required
- Neovim 0.7+ (0.9+ recommended)
- Go 1.20+
- Ginkgo v2.x
- neotest plugin
- neotest-go plugin

### Optional
- nvim-coverage (for coverage visualization)
- plenary.nvim (for async support)

## 📊 Configuration Examples

### Basic Setup
```lua
require('neotest').setup({
  adapters = {
    require("neotest-go")({
      args = { "-coverprofile=coverage.out" },
    }),
  },
})
```

### Advanced Setup
```lua
require('neotest').setup({
  adapters = {
    require("neotest-go")({
      recursive_run = true,
      args = {
        "-coverprofile=coverage.out",
        "-race",
        "-v",
      },
    }),
  },
  output = {
    open_on_run = true,
    verbose = true,
  },
})
```

## 🎨 Keybindings Reference

Common neotest keybindings (customize in your config):

| Key | Action |
|-----|--------|
| `<leader>tt` | Test current file |
| `<leader>ta` | Test all/suite |
| `<leader>tn` | Test nearest |
| `<leader>td` | Display test output |
| `<leader>ts` | Toggle test summary |

## 🔗 Related Documentation

### Project Documentation
- **[../../setup/](../../setup/)** - Project setup guides
- **[../../rules/go-guidelines.md](../../rules/go-guidelines.md)** - Go coding standards
- **[../../integration-test-strategy.md](../../integration-test-strategy.md)** - Testing strategy

### External Resources
- [neotest GitHub](https://github.com/nvim-neotest/neotest)
- [neotest-go GitHub](https://github.com/nvim-neotest/neotest-go)
- [Ginkgo Documentation](https://onsi.github.io/ginkgo/)
- [Go Testing Package](https://pkg.go.dev/testing)

## 🛠️ Alternative Testing Methods

### Command Line
```bash
# Full test suite with coverage
go test -coverprofile=coverage.out -race ./...

# Ginkgo-specific features
ginkgo --cover --race ./...

# Focused test
ginkgo --focus="specific test" ./path/to/package

# Using Makefile
make test
make coverage
```

### VSCode
1. Install Go extension
2. Use built-in test explorer
3. Run tests with code lens
4. View coverage in editor

### GoLand / IntelliJ
1. Use built-in test runner
2. Right-click test file → Run with Coverage
3. View coverage in gutter

## ✅ Verification

After setup, verify your configuration:

```bash
# Check plugins installed
:checkhealth neotest

# Check Go tools available
go version
ginkgo version

# Test coverage generation
cd /path/to/project
nvim internal/domain/career/event_test.go
# Press <leader>tt
ls -lh coverage.out  # Should exist!
```

## 📝 Adding Other Editors

To add setup guides for other editors:

1. Create new file: `EDITOR_NAME_SETUP.md`
2. Follow similar structure to NEOTEST_SETUP.md
3. Include:
   - Prerequisites
   - Installation steps
   - Configuration
   - Usage examples
   - Troubleshooting
4. Update this README with link

### Suggested Additions
- `VSCODE_SETUP.md` - VSCode Go extension setup
- `GOLAND_SETUP.md` - JetBrains IDE setup
- `VIM_SETUP.md` - Vim (non-Neovim) setup

## ℹ️ Getting Help

### Neotest Issues
1. Check `:checkhealth neotest`
2. Review [NEOTEST_SETUP.md](NEOTEST_SETUP.md) troubleshooting section
3. Check neotest logs: `:lua vim.notify(require('neotest').state)`

### General Testing Issues
1. Verify Go and Ginkgo installed correctly
2. Check project's test suite works from CLI
3. Review [../../integration-test-strategy.md](../../integration-test-strategy.md)

### Configuration Help
1. Review example configurations in [NEOTEST_SETUP.md](NEOTEST_SETUP.md)
2. Check your Neovim config syntax
3. Ensure all dependencies installed

---

**Last Updated**: 2025-12-23
**Directory**: `docs/tools/editor-setup/`
**Current Guides**: 1 (Neotest)
**Status**: Neovim setup complete, open for other editors

