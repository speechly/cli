package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"github.com/speechly/cli/pkg/clients"
)

var annotateCmd = &cobra.Command{
	Use: "annotate [<input file>] [<app id>]",
	Example: `speechly annotate -a APP_ID --input input.txt
speechly annotate -a APP_ID --input input.txt > output.txt
speechly annotate -a APP_ID --reference-date 2021-01-20 --input input.txt > output.txt

To evaluate already deployed Speechly app, you need a set of evaluation examples that users of your application might say.`,
	Short: "Create SAL annotations for a list of examples using Speechly.",
	Args:  cobra.RangeArgs(0, 1),
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

		wluClient, err := clients.WLUClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}

		refD := time.Now()
		refDS, err := cmd.Flags().GetString("reference-date")
		if err != nil {
			log.Fatalf("reference-date is invalid: %s", err)
		}

		if len(refDS) > 0 {
			refD, err = time.Parse("2006-01-02", refDS)
			if err != nil {
				log.Fatalf("reference-date is invalid: %s", err)
			}
		}

		data := readLines(inputFile)

		deAnnotate, err := cmd.Flags().GetBool("de-annotate")
		if err != nil {
			log.Fatalf("Missing de-annotated flag: %s", err)
		}

		annotated := data
		transcripts := make([]string, len(data))
		for i, line := range data {
			transcripts[i] = removeAnnotations(line)
		}
		data = transcripts

		if deAnnotate {
			for _, line := range data {
				fmt.Println(line)
			}
			os.Exit(0)
		}

		wluRequests := make([]*wluv1.WLURequest, len(data))
		for i, line := range data {
			wluRequests[i] = &wluv1.WLURequest{
				Text:          line,
				ReferenceTime: timestamppb.New(refD),
			}
		}
		textsRequest := &wluv1.TextsRequest{
			AppId:    appId,
			Requests: wluRequests,
		}

		res, err := wluClient.Texts(ctx, textsRequest)
		if err != nil {
			log.Fatal(err)
		}

		evaluate, err := cmd.Flags().GetBool("evaluate")
		if err != nil {
			log.Fatalf("Missing evaluate flag: %s", err)
		}

		if evaluate {
			EvaluateAnnotatedUtterances(wluResponsesToString(res.Responses), annotated)
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

		if err := printEvalResultCSV(outputWriter, res.Responses); err != nil {
			log.Fatalf("Error creating CSV: %s", err)
		}
	},
}

func removeAnnotations(line string) string {
	removablePattern := regexp.MustCompile(`\*([^ ]+)(?: |$)|\(([^)]+)\)`)
	line = removablePattern.ReplaceAllString(line, "")

	entityValuePattern := regexp.MustCompile(`\[([^]]+)]`)
	return entityValuePattern.ReplaceAllStringFunc(line, func(s string) string {
		pipeIndex := strings.Index(s, "|")
		if pipeIndex == -1 {
			pipeIndex = len(s) - 1
		}
		return s[1:pipeIndex]
	})
}

func readLines(fn string) []string {
	if fn != "--" {
		file, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		return scanLines(file)
	} else {
		return scanLines(os.Stdin)
	}
}

func scanLines(file *os.File) []string {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func init() {
	rootCmd.AddCommand(annotateCmd)
	annotateCmd.Flags().StringP("app", "a", "", "app id of the application to evaluate. Can alternatively be given as the last positional argument")
	annotateCmd.Flags().StringP("input", "i", "", "evaluation utterances, separated by newline, if not provided, read from stdin. Can alternatively be given as the first positional argument.")
	annotateCmd.Flags().StringP("output", "o", "", "where to store annotated utterances, if not provided, print to stdout.")
	annotateCmd.Flags().StringP("reference-date", "r", "", "reference date in YYYY-MM-DD format, if not provided use current date.")
	annotateCmd.Flags().BoolP("de-annotate", "d", false, "instead of adding annotation, remove annotations from output.")
	annotateCmd.Flags().BoolP("evaluate", "e", false, "print evaluation stats instead of the annotated output.")
}

func printEvalResultCSV(out io.Writer, items []*wluv1.WLUResponse) error {
	w := csv.NewWriter(out)
	for _, resp := range items {
		texts := make([]string, len(resp.Segments))
		for i, segment := range resp.Segments {
			texts[i] = segment.AnnotatedText
		}
		if err := w.Write([]string{strings.Join(texts, " ")}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func wluResponsesToString(responses []*wluv1.WLUResponse) []string {
	results := make([]string, len(responses))
	for i, resp := range responses {
		segmentStrings := make([]string, len(resp.Segments))
		for j, segment := range resp.Segments {
			segmentStrings[j] = segment.AnnotatedText
		}
		results[i] = strings.Join(segmentStrings, " ")

	}
	return results
}
