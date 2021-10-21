package cmd

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"github.com/speechly/cli/pkg/clients"
)

var annotateCmd = &cobra.Command{
	Use: "annotate",
	Example: `speechly annotate -a APP_ID --input input.txt
speechly annotate -a APP_ID --input input.txt > output.txt
speechly annotate -a APP_ID --reference-date 2021-01-20 --input input.txt > output.txt

To evaluate already deployed speechly app, you need a set of evaluation examples that users of your application might say.`,
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

func init() {
	rootCmd.AddCommand(annotateCmd)
	annotateCmd.Flags().StringP("app", "a", "", "app id of the application to evaluate.")
	annotateCmd.Flags().StringP("input", "i", "", "evaluation utterances, separated by newline.")
	annotateCmd.Flags().StringP("reference-date", "r", "", "reference date in ISO format, eg. YYYY-MM-DD. If no date provided, the current date is used.")
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
