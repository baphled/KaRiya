//go:build integration
// +build integration

package testutil

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	. "github.com/cucumber/godog"
)

func FeatureContext(s *Suite) {
	s.Step(`^I run "([^"]*)" with no arguments$`, iRunWithNoArguments)
	s.Step(`^I run "([^"]*) --help"$`, iRunWithHelp)
	s.Step(`^I run "([^"]*) version"$`, iRunWithVersion)
	s.Step(`^the TUI should launch and display the main menu$`, tuiShouldLaunch)
	s.Step(`^the CLI should display help text$`, cliShouldDisplayHelp)
	s.Step(`^the CLI should display the version$`, cliShouldDisplayVersion)
	s.Step(`^the TUI should not launch$`, tuiShouldNotLaunch)
}

var (
	lastOutput string
	lastErr    error
)

func iRunWithNoArguments(cmd string) error {
	c := exec.Command(cmd)
	out, err := c.CombinedOutput()
	lastOutput = string(out)
	lastErr = err
	return nil
}

func iRunWithHelp(cmd string) error {
	c := exec.Command(cmd, "--help")
	out, err := c.CombinedOutput()
	lastOutput = string(out)
	lastErr = err
	return nil
}

func iRunWithVersion(cmd string) error {
	c := exec.Command(cmd, "version")
	out, err := c.CombinedOutput()
	lastOutput = string(out)
	lastErr = err
	return nil
}

func tuiShouldLaunch() error {
	if !strings.Contains(lastOutput, "KaRiya") && !strings.Contains(lastOutput, "terminal user interface") {
		return context.Canceled // fail
	}
	return nil
}

func cliShouldDisplayHelp() error {
	if !strings.Contains(lastOutput, "Usage:") {
		return context.Canceled
	}
	return nil
}

func cliShouldDisplayVersion() error {
	if !strings.Contains(lastOutput, "version") {
		return context.Canceled
	}
	return nil
}

func tuiShouldNotLaunch() error {
	if strings.Contains(lastOutput, "KaRiya") && strings.Contains(lastOutput, "terminal user interface") {
		return context.Canceled
	}
	return nil
}
