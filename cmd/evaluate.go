package cmd

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	"fmt"

	"github.com/spf13/cobra"

	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"github.com/speechly/cli/pkg/clients"
)

var evaluateCmd = &cobra.Command{
	Use: "evaluate",
	Example: `speechly evaluate run [flags]
speechly evaluate accuracy [flags]

To evaluate already deployed speechly app,
- check the appid of your app
- write down list of evaluation examples that users of your application might say

The examples should be written in a (single column) text/csv file, where each line corresponds one example.

Evaluation consists three steps
1) run 'speechly evaluate run' to annotate your evaluation examples. Check 'speechly evaluate run --help' for details.
2) create a ground truth reference based on the annotated examples.
3) compute accuracy between the annotated examples and ground truth. Check 'speechly evaluate accuracy --help' for details.

More information at docs.speechly.com
`,
	Short: "Evaluate a list of example utterances.",
}

var evaluateRunCmd = &cobra.Command{
	Use:     "run",
	Example: `speechly evaluate run -a APP_ID --input input.csv --annotated output.csv`,
	Short:   "Create SAL annotations for a list of examples using Speechly.",
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

		outputFile, err := cmd.Flags().GetString("annotated")
		if err != nil || len(outputFile) == 0 {
			log.Fatalf("Annotated file is invalid: %v", err)
		}

		file, err := os.Open(inputFile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		csvReader := csv.NewReader(file)
		data, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal(err)
		}

		wluRequests := make([]*wluv1.WLURequest, len(data))
		for i, line := range data {
			wluRequests[i] = &wluv1.WLURequest{
				Text: line[0],
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

		if err := printEvalResultCSV(outputFile, res.Responses); err != nil {
			log.Fatalf("Error creating CSV: %s", err)
		}
		fmt.Printf("Wrote the annotated examples to %v\n", outputFile)
	},
}

func readCsv(file io.Reader) [][]string {
	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	return data
}

type UtteranceComparator struct {
	entityRe        *regexp.Regexp
	postProcessedRe *regexp.Regexp
}

// compare two SAL utterance
func (this *UtteranceComparator) Equal(a string, b string) bool {
	aEntityIndexes := this.entityRe.FindAllIndex([]byte(a), -1)
	bEntityIndexes := this.entityRe.FindAllIndex([]byte(b), -1)

	if len(aEntityIndexes) != len(bEntityIndexes) {
		return false
	}
	if len(aEntityIndexes) == 0 {
		return a == b
	}

	aPtr := 0
	bPtr := 0
	for i, aIdx := range aEntityIndexes {
		bIdx := bEntityIndexes[i]
		if a[aPtr:aIdx[0]] != b[bPtr:bIdx[0]] {
			return false
		}
		aPtr = aIdx[1]
		bPtr = bIdx[1]
		aEntity := a[aIdx[0]:aIdx[1]]
		bEntity := b[bIdx[0]:bIdx[1]]

		aEntityOnlyValue := this.postProcessedRe.ReplaceAllString(aEntity, `$1`)
		bEntityOnlyValue := this.postProcessedRe.ReplaceAllString(bEntity, `$1`)
		if aEntityOnlyValue != bEntityOnlyValue {
			return false
		}
	}
	// check text after last entity
	aLastEntityEnd := aEntityIndexes[len(aEntityIndexes)-1][1]
	bLastEntityEnd := bEntityIndexes[len(bEntityIndexes)-1][1]
	return a[aLastEntityEnd:] == b[bLastEntityEnd:]
}

// pre-compile regexes for utterance comparison
func CreateComparator() UtteranceComparator {
	entityRe := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	postProcessedRe := regexp.MustCompile(`\[(.+?)\|(.*?)\]`)
	return UtteranceComparator{entityRe, postProcessedRe}
}

var evaluateAccuracyCmd = &cobra.Command{
	Use:     "accuracy",
	Example: `speechly evaluate accuracy --annotated output.csv --ground-truth ground-truth.csv`,
	Short:   "Compute accuracy between annotated examples (given by 'speechly evaluate run') and ground truth.",
	Run: func(cmd *cobra.Command, args []string) {
		annotatedFn, err := cmd.Flags().GetString("annotated")
		if err != nil || len(annotatedFn) == 0 {
			log.Fatalf("Annotated file is invalid: %v", err)
		}
		annotatedFile, err := os.Open(annotatedFn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := annotatedFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		groundTruthFn, err := cmd.Flags().GetString("ground-truth")
		if err != nil || len(groundTruthFn) == 0 {
			log.Fatalf("Ground-truth file is invalid: %v", err)
		}
		groundTruthFile, err := os.Open(groundTruthFn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := groundTruthFile.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		annotatedData := readCsv(annotatedFile)
		groundTruthData := readCsv(groundTruthFile)
		if len(annotatedData) != len(groundTruthData) {
			log.Fatalf(
				"Input csv files should have same length, but --annotated has %d lines and --ground-truth %d lines.",
				len(annotatedData),
				len(groundTruthData),
			)
		}

		n := float64(len(annotatedData))
		hits := 0.0
		comp := CreateComparator()
		for i, aLine := range annotatedData {
			gtUtt := groundTruthData[i][0]
			aUtt := aLine[0]
			if comp.Equal(aUtt, gtUtt) {
				hits += 1.0
				continue
			}
			fmt.Println("Expected:")
			fmt.Println("  " + gtUtt)
			fmt.Println("but had:")
			fmt.Println("  " + aUtt)
			fmt.Println()
		}
		fmt.Println("Accuracy:")
		fmt.Printf("%.2f\n", hits/n)
	},
}

func init() {
	rootCmd.AddCommand(evaluateCmd)
	evaluateCmd.AddCommand(evaluateRunCmd)
	evaluateCmd.AddCommand(evaluateAccuracyCmd)
	evaluateRunCmd.Flags().StringP("app", "a", "", "app id of the application to evaluate.")
	evaluateRunCmd.Flags().StringP("input", "i", "", "evaluation utterances, separated by newline.")
	evaluateRunCmd.Flags().StringP("annotated", "", "", "output location, where annotated examples will be written.")
	evaluateAccuracyCmd.Flags().StringP("annotated", "", "", "SAL annotated utterances, as given by 'speechly evaluate run' command.")
	evaluateAccuracyCmd.Flags().StringP("ground-truth", "", "", "manually verified ground-truths for annotated examples")
}

func printEvalResultCSV(outputFile string, items []*wluv1.WLUResponse) error {
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	w := csv.NewWriter(file)
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
