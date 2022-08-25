package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/speechly/cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	cmd.RootCmd.DisableAutoGenTag = true

	if len(os.Args) < 2 {
		log.Fatal("Output directory must be given as argument\n")
		os.Exit(1)
	}

	outDir, _ := filepath.Abs(os.Args[1])
	src := filepath.Join(outDir, "speechly.md")
	dest := filepath.Join(outDir, "README.md")

	if err := doc.GenMarkdownTree(cmd.RootCmd, outDir); err != nil {
		log.Fatal(err)
	}

	bytes, err := os.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}

	err = os.WriteFile(dest, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
