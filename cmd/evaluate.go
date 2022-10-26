package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var evaluateCmd = &cobra.Command{
	Use:   "evaluate command [flags]",
	Short: "Evaluate application model accuracy.",
	Args:  cobra.NoArgs,
}

var nluCmd = &cobra.Command{
	Use:     "nlu <app_id> <input_file>",
	Example: `speechly evaluate nlu <app_id> annotated-utterances.txt`,
	Short:   "Evaluate the NLU accuracy of the given application model.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appID := args[0]
		res, annotated, err := runThroughWLU(ctx, appID, args[1], false, time.Now())
		if err != nil {
			log.Fatalf("WLU failed: %v", err)
		}

		evaluateAnnotatedUtterances(wluResponsesToString(res.Responses), annotated)
	},
}

var asrCmd = &cobra.Command{
	Use:     "asr <app_id> <input_file>",
	Example: `speechly evaluate asr <app_id> utterances.jsonlines`,
	Short:   "Evaluate the ASR accuracy of the given application model.",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appID := args[0]
		ac, err := transcribeWithBatchAPI(ctx, appID, args[1], true)
		if err != nil {
			log.Fatalf("Transcription failed: %v", err)
		}

		ed := EditDistance{}
		for _, aci := range ac {
			wd, err := wordDistance(aci.Transcript, aci.Hypothesis)
			if err != nil {
				log.Fatalf("Error in result generation: %v", err)
			}
			ed = ed.Add(wd)
		}
		fmt.Printf("Word Error Rate (WER):\n%.2f\t%d/%d\n", ed.AsER(), ed.dist, ed.base)
	},
}

func init() {
	RootCmd.AddCommand(evaluateCmd)
	evaluateCmd.AddCommand(nluCmd)
	nluCmd.Flags().StringP("reference-date", "r", "", "Reference date in YYYY-MM-DD format, if not provided use current date.")

	evaluateCmd.AddCommand(asrCmd)
	asrCmd.Flags().StringP("reference-date", "r", "", "Reference date in YYYY-MM-DD format, if not provided use current date.")
}
