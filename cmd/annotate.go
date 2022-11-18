package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"github.com/spf13/cobra"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Create SAL annotations for a list of examples using Speechly.",
	Long:  "To evaluate already deployed Speechly app, you need a set of evaluation examples that users of your application might say.",
	Example: `speechly annotate input.txt <app_id>
speechly annotate --app <app_id> --input input.txt > output.txt
speechly annotate --app <app_id> --reference-date 2021-01-20 --input input.txt > output.txt`,
	Args: cobra.RangeArgs(0, 1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		appId, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatalf("App ID is invalid: %s", err)
		}
		if appId == "" && len(args) < 1 {
			return fmt.Errorf("app_id must be given with flag --app or as the first positional argument of two")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		appId, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatalf("App ID is invalid: %s", err)
		}
		inputFile, err := cmd.Flags().GetString("input")
		if err != nil {
			log.Fatalf("Input file is invalid: %v", err)
		}
		if appId == "" && inputFile == "" {
			if len(args) == 2 {
				inputFile = args[0]
				appId = args[1]
			} else {
				inputFile = "--"
				appId = args[0]
			}
		}
		if inputFile == "" {
			if appId != "" && len(args) == 1 {
				inputFile = args[0]
			} else {
				inputFile = "--"
			}
		}

		if appId == "" {
			appId = args[0]
		}

		refD, err := readReferenceDate(cmd)
		if err != nil {
			log.Fatalf("Faild to get reference date: %s", err)
		}

		deAnnotate, err := cmd.Flags().GetBool("de-annotate")
		if err != nil {
			log.Fatalf("Missing de-annotated flag: %s", err)
		}

		res, annotated, err := runThroughWLU(ctx, appId, inputFile, deAnnotate, refD)
		if err != nil {
			log.Fatalf("WLU failed: %s", err)
		}

		evaluate, err := cmd.Flags().GetBool("evaluate")
		if err != nil {
			log.Fatalf("Missing evaluate flag: %s", err)
		}

		if evaluate {
			evaluateAnnotatedUtterances(wluResponsesToString(res.Responses), annotated)
			os.Exit(0)
		}

		var outputWriter io.Writer
		outputFile, err := cmd.Flags().GetString("output")
		if err == nil && len(outputFile) > 0 {
			outputFile, _ = filepath.Abs(outputFile)
			outputWriter, err = os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalf("output path is invalid: %s", err)
			}
		} else {
			outputWriter = cmd.OutOrStdout()
		}

		if err := printEvalResultTXT(outputWriter, res.Responses); err != nil {
			log.Fatalf("Error creating CSV: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(annotateCmd)
	annotateCmd.Flags().StringP("app", "a", "", "Application to evaluate. Can be given as the second positional argument.")
	annotateCmd.Flags().StringP("input", "i", "", "Evaluation utterances, separated by newline, if not provided, read from stdin. Can be given as the first positional argument.")
	annotateCmd.Flags().StringP("output", "o", "", "Where to store annotated utterances, if not provided, print to stdout.")
	annotateCmd.Flags().StringP("reference-date", "r", "", "Reference date in YYYY-MM-DD format, if not provided use current date.")
	annotateCmd.Flags().BoolP("de-annotate", "d", false, "Instead of adding annotation, remove annotations from output.")
	annotateCmd.Flags().BoolP("evaluate", "e", false, "Print evaluation stats instead of the annotated output.")
}

func printEvalResultTXT(out io.Writer, items []*wluv1.WLUResponse) error {
	for _, resp := range items {
		texts := make([]string, len(resp.Segments)+1)
		texts[len(resp.Segments)] = "\n"
		for i, segment := range resp.Segments {
			texts[i] = segment.AnnotatedText
		}
		if _, err := io.WriteString(out, strings.Join(texts, " ")); err != nil {
			return err
		}
	}
	return nil
}
