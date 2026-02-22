# PR #175: Version Injection Demo
# Demonstrates build-time version injection from VERSION file
# The binary gets its version from the VERSION file via ldflags

Output demos/version-injection.gif

Set Shell "bash"
Set FontSize 14
Set FontFamily "JetBrains Mono"
Set Width 900
Set Height 400
Set Padding 15
Set BorderRadius 8
Set WindowBar Colorful
Set Theme "Catppuccin Mocha"
Set TypingSpeed 50ms
Set Framerate 30

# Show the VERSION file
Type "cat VERSION"
Enter
Sleep 2s

# Build with version injection
Type "make build"
Enter
Sleep 3s

# Show the version output
Type "./kariya --version"
Enter
Sleep 3s
