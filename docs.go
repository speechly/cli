package main

import (
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra/doc"

	"github.com/speechly/cli/cmd"
)

func main() {
	cmd.RootCmd.DisableAutoGenTag = true
	dir := "docs"

	if err := doc.GenMarkdownTree(cmd.RootCmd, dir); err != nil {
		log.Fatal(err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var s string
	for _, f := range files {
		if strings.Contains(f.Name(), "speechly.md") || strings.Contains(f.Name(), "speechly_projects.md") {
			continue
		}

		bytes, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			log.Fatal(err)
		}

		links := regexp.MustCompile(`### SEE ALSO([\s\S]*?)\z`)
		outLinks := links.ReplaceAllString(string(bytes), "")

		headings := strings.ReplaceAll(outLinks, "### ", "##### ")

		commands := regexp.MustCompile(`## speechly ([a-z ]+)`)
		outCommands := commands.ReplaceAllString(headings, "### `$1`")

		s += outCommands
	}

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(s), 0644); err != nil {
		log.Fatal(err)
	}
}
