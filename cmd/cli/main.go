package main

import (
	"io"
	"os"

	"github.com/baphled/kariya/internal/cli/cmd"
)

var version = "dev"

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, out io.Writer, errOut io.Writer) int {
	rootCmd := cmd.NewRootCmd(version)
	rootCmd.SetOut(out)
	rootCmd.SetErr(errOut)
	rootCmd.SetArgs(args)

	if err := rootCmd.Execute(); err != nil {
		return 1
	}
	return 0
}
