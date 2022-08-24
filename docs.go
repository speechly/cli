package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/cobra/doc"

	"github.com/speechly/cli/cmd"
)

func main() {
	dir := "docs"
	src := filepath.Join(dir, "speechly.md")
	dest := filepath.Join(dir, "README.md")

	if err := doc.GenMarkdownTree(cmd.RootCmd, dir); err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(src)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(dest, bytes, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
