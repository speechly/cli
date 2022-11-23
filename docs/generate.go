package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"

	"github.com/speechly/cli/cmd"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Output directory must be given as argument\n")
		os.Exit(1)
	}
	outDir, _ := filepath.Abs(os.Args[1])

	rmBuf := new(bytes.Buffer)
	rmBuf.WriteString("# Available commands\n\n")

	for _, c := range cmd.RootCmd.Commands() {
		cBuf := new(bytes.Buffer)
		c.InitDefaultHelpFlag()

		name := c.Name()
		rmBuf.WriteString("#### [`" + name + "`](" + name + ".md)\n\n" + c.Short + "\n\n")
		cBuf.WriteString(header(c))

		if c.HasSubCommands() {
			cBuf.WriteString("### Subcommands\n\n")
			for _, sc := range c.Commands() {
				sc.InitDefaultHelpFlag()
				scBuf := new(bytes.Buffer)
				scName := sc.Name()
				link := fmt.Sprintf("[`%s %s`](%s_%s.md)", name, scName, name, scName)
				rmBuf.WriteString("#### " + link + "\n\n" + sc.Short + "\n\n")
				cBuf.WriteString("* " + link + " - " + sc.Short + "\n")

				scBuf.WriteString(header(sc))

				if sc.HasFlags() {
					scBuf.WriteString(flags(sc))
				}

				if sc.HasExample() {
					scBuf.WriteString(example(sc))
				}

				file := name + "_" + scName + ".md"
				createFile(file, outDir, scBuf.Bytes())
			}
			cBuf.WriteString("\n")
		}

		if c.HasFlags() {
			cBuf.WriteString(flags(c))
		}

		if c.HasExample() {
			cBuf.WriteString(example(c))
		}

		file := name + ".md"
		createFile(file, outDir, cBuf.Bytes())
	}

	createFile("README.md", outDir, rmBuf.Bytes())
}

func createFile(name string, dir string, buf []byte) {
	out := filepath.Join(dir, name)
	if err := os.WriteFile(out, buf, 0644); err != nil {
		log.Fatal(err)
	}
}

func header(c *cobra.Command) string {
	b := new(bytes.Buffer)

	name := ""
	if c.Parent() != c.Root() {
		name = c.Parent().Name() + " " + c.Name()
	} else {
		name = c.Name()
	}

	b.WriteString("# " + name + "\n\n")
	b.WriteString(c.Short + "\n\n")

	if c.Use != "" {
		usage := c.UseLine()
		b.WriteString("### Usage\n\n")
		b.WriteString("```\n" + usage + "\n```")
		b.WriteString("\n\n")
	}

	if c.Long != "" {
		b.WriteString(c.Long + "\n\n")
	}

	if c.Aliases != nil {
		a := fmt.Sprintf("%v", strings.Join(c.Aliases, ", "))
		b.WriteString(fmt.Sprintf("Command aliases: `%s`\n\n", a))
	}

	return b.String()
}

func flags(c *cobra.Command) string {
	b := new(bytes.Buffer)

	if c.HasFlags() {
		b.WriteString("### Flags\n\n")
		c.Flags().VisitAll(func(f *flag.Flag) {
			vt := f.Value.Type()
			if f.Shorthand == "" {
				b.WriteString("* `--" + f.Name + "` _(" + vt + ")_ - " + f.Usage + "\n")
			} else {
				b.WriteString("* `--" + f.Name + "` `-" + f.Shorthand + "` _(" + vt + ")_ - " + f.Usage + "\n")
			}
		})
	}

	return b.String()
}

func example(c *cobra.Command) string {
	b := new(bytes.Buffer)

	b.WriteString("\n")
	b.WriteString("### Examples\n\n")
	b.WriteString("```\n" + c.Example + "\n```")
	b.WriteString("\n")

	return b.String()
}
