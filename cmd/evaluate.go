package cmd

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"fmt"

	"github.com/spf13/cobra"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"github.com/speechly/cli/pkg/clients"
)

var evaluateCmd = &cobra.Command{
	Use: "evaluate",
	Example: `speechly evaluate annotate [flags]
speechly evaluate accuracy [flags]

To evaluate already deployed speechly app,
- check the appid of your app
- write down list of evaluation examples that users of your application might say

The examples should be written in a text file, where each line corresponds one example.

Evaluation consists three steps
1) run 'speechly evaluate annotate' to annotate your evaluation examples. Check 'speechly evaluate annotate --help' for details.
2) compute accuracy between the annotated examples and ground truth. Check 'speechly evaluate accuracy --help' for details.

More information at docs.speechly.com
`,
	Short: "Evaluate a list of example utterances.",
}

var evaluateAnnotateCmd = &cobra.Command{
	Use: "annotate",
	Example: `speechly evaluate annotate -a APP_ID --input input.txt
speechly evaluate annotate -a APP_ID --input input.txt > output.txt
speechly evaluate annotate -a APP_ID --reference-date 2021-01-20 --input input.txt > output.txt`,
	Short: "Create SAL annotations for a list of examples using Speechly.",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		appId, err := cmd.Flags().GetString("app")
		if err != nil || len(appId) == 0 {
			log.Fatalf("App ID is invalid: %s", err)
		}

		wluClient, err := clients.WLUClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}

		inputFile, err := cmd.Flags().GetString("input")
		if err != nil || len(inputFile) == 0 {
			log.Fatalf("Input file is invalid: %v", err)
		}

		refDS, err := cmd.Flags().GetString("reference-date")
		if err != nil {
			log.Fatalf("reference-date is invalid: %s", err)
		}
		refD, err := time.Parse("2006-01-02", refDS)
		if err != nil {
			log.Fatalf("reference-date is invalid: %s", err)
		}

		data := readLines(inputFile)

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

		if err := printEvalResultCSV(cmd.OutOrStdout(), res.Responses); err != nil {
			log.Fatalf("Error creating CSV: %s", err)
		}
	},
}

func readLines(fn string) []string {
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
	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

var evaluateAccuracyCmd = &cobra.Command{
	Use:     "accuracy",
	Example: `speechly evaluate accuracy --input output.txt --ground-truth ground-truth.txt`,
	Short:   "Compute accuracy between annotated examples (given by 'speechly evaluate annotate') and ground truth.",
	Run: func(cmd *cobra.Command, args []string) {
		annotatedFn, err := cmd.Flags().GetString("input")
		if err != nil || len(annotatedFn) == 0 {
			log.Fatalf("Annotated file is invalid: %v", err)
		}

		groundTruthFn, err := cmd.Flags().GetString("ground-truth")
		if err != nil || len(groundTruthFn) == 0 {
			log.Fatalf("Ground-truth file is invalid: %v", err)
		}
		annotatedData := readLines(annotatedFn)
		groundTruthData := readLines(groundTruthFn)
		if len(annotatedData) != len(groundTruthData) {
			log.Fatalf(
				"Input files should have same length, but --annotated has %d lines and --ground-truth %d lines.",
				len(annotatedData),
				len(groundTruthData),
			)
		}

		n := float64(len(annotatedData))
		hits := 0.0
		for i, aUtt := range annotatedData {
			gtUtt := groundTruthData[i]
			if aUtt == gtUtt {
				hits += 1.0
				continue
			}
			fmt.Println("Ground truth had:")
			fmt.Println("  " + gtUtt)
			fmt.Println("but prediction was:")
			fmt.Println("  " + aUtt)
			fmt.Println()
		}
		fmt.Println("Matching rows out of total: ")
		fmt.Printf("%.0f / %.0f\n", hits, n)
		fmt.Println("Accuracy:")
		fmt.Printf("%.2f\n", hits/n)
	},
}

func init() {
	rootCmd.AddCommand(evaluateCmd)
	evaluateCmd.AddCommand(evaluateAnnotateCmd)
	evaluateCmd.AddCommand(evaluateAccuracyCmd)
	evaluateAnnotateCmd.Flags().StringP("app", "a", "", "app id of the application to evaluate.")
	evaluateAnnotateCmd.Flags().StringP("input", "i", "", "evaluation utterances, separated by newline.")
	evaluateAnnotateCmd.Flags().StringP("reference-date", "r", "", "reference date in ISO format, eg. YYYY-MM-DD.")
	evaluateAccuracyCmd.Flags().StringP("input", "", "", "SAL annotated utterances, as given by 'speechly evaluate annotate' command.")
	evaluateAccuracyCmd.Flags().StringP("ground-truth", "", "", "manually verified ground-truths for annotated examples")
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
