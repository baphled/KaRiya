# Neotest Setup for Ginkgo Tests with Automatic Coverage

## Overview

This project uses **neotest-go** for running Ginkgo tests in Neovim with automatic code coverage generation. This guide covers setup, configuration, usage, and troubleshooting.

## Why neotest-go for Ginkgo?

**Key Insight**: Ginkgo tests compile into standard Go test suites, so `go test` can run them natively!

```bash
# This works with Ginkgo tests
go test -v -coverprofile=coverage.out ./...
```

Since neotest-go uses `go test`, it can run Ginkgo tests AND generate coverage automatically.

## Configuration

### 1. Dependencies (`~/.config/nvim/lua/baphled/plugins.lua`)

```lua
"nvim-neotest/neotest-go",  -- Better coverage support
```

### 2. Adapter Configuration (`~/.config/nvim/lua/baphled/config/neotest.lua`)

```lua
require('neotest').setup({
  adapters = {
    require("neotest-go")({
      recursive_run = true,  -- Run tests in nested packages
      args = {
        "-coverprofile=coverage.out",  -- ✅ Automatic coverage!
        "-race"                        -- Race condition detection
      },
    }),
  },
  output = {
    open_on_run = true,  -- Automatically open output
    verbose = true,      -- Enable verbose output
  },
  diagnostic = {
    enabled = true,
  },
  log_level = vim.log.levels.DEBUG
})
```

### 3. Optional: nvim-coverage Integration

For visual coverage highlighting in Neovim:

```lua
require("coverage").setup({
  auto_reload = true,
  lang = {
    go = {
      coverage_file = "coverage.out",
    },
  },
})
```

## Benefits

### ✅ What You Gain

1. **Automatic Coverage Generation**
   - Every test run generates `coverage.out` at project root
   - No manual commands needed
   - Works with existing `make coverage` workflow

2. **Better Go Integration**
   - Uses standard `go test` command
   - Respects `go.mod` and Go tooling
   - Better maintained and documented

3. **Simpler Configuration**
   - Coverage flags passed via `args` parameter
   - No need for custom wrappers or scripts
   - Works out of the box

### ❌ Trade-offs

**Individual Ginkgo Spec Detection Lost**

- ~~nvim-ginkgo could detect individual `Describe`/`It`/`Context` blocks~~
- neotest-go only detects `Test*` functions (the Ginkgo suite entry point)

**Impact:**
- You can still run entire test files/packages ✅
- You just can't run individual Ginkgo specs from neotest ❌
- Use command line for granular Ginkgo testing: `ginkgo --focus="spec description"`

## Usage

### Running Tests with Coverage

All your neotest keybindings now generate coverage automatically:

| Keybinding | Action | Coverage Generated |
|------------|--------|-------------------|
| `<leader>tt` | Test current file | ✅ Yes |
| `<leader>ta` | Test all (suite) | ✅ Yes |
| `<leader>tn` | Test nearest | ✅ Yes |

### Viewing Coverage

After running tests, view coverage with:

```bash
# Generate HTML coverage report
make coverage

# View in browser
open coverage/index.html

# Or check coverage percentage
go tool cover -func=coverage.out | tail -1
```

### Example Workflow

1. **Open a test file:**
   ```bash
   nvim internal/domain/career/event_test.go
   ```

2. **Run tests:**
   - Press `<leader>tt` to test current file
   - Or `<leader>ta` to test entire suite

3. **Check coverage:**
   ```bash
   ls -lh coverage.out  # Should exist now!
   make coverage        # Generate HTML report
   ```

4. **View results:**
   - Check neotest output panel for test results
   - Open `coverage/index.html` for detailed coverage

## Verification

Run this verification script to confirm setup:

```bash
# Check configuration
grep "neotest-go" ~/.config/nvim/lua/baphled/plugins.lua
grep "coverprofile" ~/.config/nvim/lua/baphled/config/neotest.lua

# Test coverage generation manually
go test -coverprofile=coverage-test.out ./internal/domain/career
cat coverage-test.out | head -5
rm coverage-test.out
```

Expected output:
```
✅ neotest-go found in plugins.lua
✅ coverprofile flag found in neotest.lua
✅ Coverage file generated
```

## Troubleshooting

### Coverage Not Generated

**Issue:** Tests run but `coverage.out` not created

**Solutions:**
1. Restart Neovim: `:qa` then reopen
2. Sync plugins: `:Lazy sync`
3. Check test output: `<leader>td` (display test output)
4. Verify configuration: `grep coverprofile ~/.config/nvim/lua/baphled/config/neotest.lua`

### Tests Not Running

**Issue:** neotest can't find tests

**Solutions:**
1. Ensure file ends with `_test.go`
2. Ensure you're in the project root (where `go.mod` exists)
3. Check `:Neotest summary` for detected tests
4. Verify `go.mod` exists: `ls go.mod`

### No Output Displayed

**Issue:** Tests run but no output shows

**Solutions:**
1. Check `:messages` for hidden errors
2. Verify neotest configuration has `open_on_run = true`
3. Use `:Neotest output` to manually open output panel
4. Increase Neovim command line height:
   ```lua
   vim.o.cmdheight = 2
   ```

### Individual Specs Not Showing

**This is expected!** neotest-go doesn't parse Ginkgo specs.

**Workarounds:**
- Run entire test file with `<leader>tt`
- Use CLI for focused tests: `ginkgo --focus="specific spec"`
- Use ginkgo labels in test code for better organization

### Plugin Version Issues

**Issue:** Neotest-go not working correctly

**Solutions:**
1. Update plugins: `:Lazy sync`
2. Check plugin versions:
   - neotest: latest
   - neotest-go: latest
3. Verify Go and Ginkgo versions:
   ```bash
   go version  # Should be 1.20+
   ginkgo version  # Should be v2.x
   ```

### Debugging Test Discovery

If tests aren't discovered:

1. **Check Neovim messages:**
   ```vim
   :messages
   ```

2. **Verify test structure:**
   ```go
   var _ = Describe("MyTestSuite", func() {
       Context("Specific Scenario", func() {
           It("should do something specific", func() {
               // Test logic here
           })
       })
   })
   ```

3. **Check neotest log:**
   ```vim
   :lua vim.notify(vim.inspect(require('neotest').state))
   ```

## Alternative: Command-Line Coverage

If you prefer CLI-based testing with Ginkgo-specific features:

```bash
# Full coverage with Ginkgo
ginkgo --cover --coverprofile=coverage.out --race ./...

# Focused test with coverage
ginkgo --focus="CareerEvent Validation" --cover ./internal/domain/career

# Use your Makefile
make coverage  # Runs ginkgo and generates HTML
```

## Useful Neotest Commands

| Command | Description |
|---------|-------------|
| `:Neotest output` | Open test output panel |
| `:Neotest summary` | Show test summary tree |
| `:Neotest run` | Run nearest test |
| `:Neotest run file` | Run all tests in file |
| `:Neotest stop` | Stop running tests |
| `:messages` | Check for error messages |

## Related Files

- Configuration: `~/.config/nvim/lua/baphled/config/neotest.lua`
- Plugins: `~/.config/nvim/lua/baphled/plugins.lua`
- Keybindings: `~/.config/nvim/lua/baphled/config/which-key/neotest.lua`
- Project Makefile: `./Makefile` (has `make coverage` target)
- Ginkgo Config: `./ginkgo.yml` (still used by CLI `ginkgo` command)

## Summary

**Before:** Manual coverage setup with nvim-ginkgo
**After:** Automatic coverage with neotest-go

**Result:** Every test run generates `coverage.out` automatically! 🎉

The trade-off (losing individual spec detection) is worth it for automatic coverage generation, especially since you can still run entire test files and use CLI for focused testing.

## Next Steps

1. ✅ Configuration updated
2. ⏳ Restart Neovim (`:qa` then reopen)
3. ⏳ Run `:Lazy sync` to ensure plugins are synced
4. ⏳ Test with `<leader>tt` on a test file
5. ⏳ Verify `coverage.out` exists
6. ⏳ Run `make coverage` to view HTML report

**You're all set!** Coverage will now be generated automatically on every test run. 🚀

