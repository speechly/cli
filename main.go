package main

import (
	"os"

	"github.com/speechly/cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		// The cmd error handlers have reported this error already, no need to print it here.
		os.Exit(1)
	}
}
