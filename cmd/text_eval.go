package cmd

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"github.com/speechly/cli/pkg/clients"
)

var textEvalCmd = &cobra.Command{
	Use: "text-eval",
	Example: `speechly text-eval
# To run this command, you need a csv file where each line corresponds utterances that
# users of your application might say to your Speechly app. The following command provides an example of input file format.
echo 'utterance 1\nsecond utt\nlast utterance' > input.csv
# Let's suppose you have these examples in input.csv.

speechly text-eval -a UUID_APP_ID --input input.csv > output.csv

# In order to see the accuracy of your application, you can use the following framework.
# 1) Run this command to get WLU (Written Language Understanding) results.
# 2) Go through the results and fix any potential errors. This gives you the (ground truth) reference.
# 3) Compute the accuracy based on the WLU results and the ground truth.
#
# So in practise, start by 1)
speechly text-eval -a UUID_APP_ID --input input.csv > output.csv
# Then step 2)
# Copy the content of output.csv to ground-truth.csv
cp output.csv ground-truth.csv
# Go through the lines in the ground-truth.csv and fix all errors in the annotations and save the file.
# Once done, proceed to 3).
# The following terminal command counts the proportion number of lines that differ between the output.csv and ground-truth.csv
` +
"hits=`diff -y --suppress-common-lines eval-output.csv ground-truth.csv | wc -l`; total=`cat ground-truth.csv | wc -l`; bc <<< \"scale=2; 1 - $hits / $total\"" +
`
# This is the accuracy of your current application. Closer the value 1, the better while 0 is the worst.
# If you modify and deploy your application, you can re-run the evaluation easily by using same 
# evaluation examples (input.csv) and ground truths (ground-truth.csv).
# Just repeat the steps 1 and 3, as described above to compute the accuracy of the latest version of the application.
`,
	Short: "Run WLU recognition for provided list of text utterances.",
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
			AppId: appId,
			Requests: wluRequests,
		}

		res, err := wluClient.Texts(ctx, textsRequest)
		if err != nil {
			log.Fatal(err)
		}
		if isatty.IsTerminal(os.Stdout.Fd()) {
			for _, resp := range res.Responses {
				texts := make([]string, len(res.Responses))
				for i, segment := range resp.Segments {
					texts[i] = segment.AnnotatedText
				}
				cmd.Println(strings.Join(texts, " "))
			}
		} else {
			if err := printEvalResultCSV(cmd.OutOrStdout(), res.Responses); err != nil {
				log.Fatalf("Error creating CSV: %s", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(textEvalCmd)
	textEvalCmd.Flags().StringP("app", "a", "", "application for WLU recognition.")
	textEvalCmd.Flags().StringP("input", "i", "", "csv file with your examples.")
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
