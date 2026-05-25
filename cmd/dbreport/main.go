package main

import (
	"os"

	"github.com/rswestmoreland/dbreport/internal/cli"
)

func main() {
	runner := cli.Runner{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	os.Exit(runner.Run(os.Args[1:]))
}
