package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

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

	filePrepender := func(name string) string {
		return ""
	}

	linkHandler := func(name string) string {
		base := strings.TrimPrefix(name, "speechly_")
		if strings.Contains(name, "speechly.md") {
			return "README.md"
		}
		return base
	}

	if err := doc.GenMarkdownTreeCustom(cmd.RootCmd, outDir, filePrepender, linkHandler); err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(outDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), "generate.go") {
			continue
		}

		file := filepath.Join(outDir, f.Name())
		bytes, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		headings := strings.ReplaceAll(string(bytes), "## ", "# ")
		seeAlso := strings.ReplaceAll(headings, "SEE ALSO", "See also")
		titles := strings.ReplaceAll(seeAlso, "# speechly ", "# ")
		links := strings.ReplaceAll(titles, "* [speechly ", "* [")
		if err := os.WriteFile(file, []byte(links), 0644); err != nil {
			log.Fatal(err)
		}

		if strings.Contains(f.Name(), "speechly.md") {
			if err := os.Rename(file, filepath.Join(outDir, "README.md")); err != nil {
				log.Fatal(err)
			}
		} else {
			if err := os.Rename(file, filepath.Join(outDir, strings.TrimPrefix(f.Name(), "speechly_"))); err != nil {
				log.Fatal(err)
			}
		}
	}
}
