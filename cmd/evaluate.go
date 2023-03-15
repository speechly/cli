package cmd

import (
	"fmt"
	"log"

	"github.com/speechly/nwalgo"
	"github.com/spf13/cobra"
)

var evaluateCmd = &cobra.Command{
	Use:   "evaluate [command]",
	Short: "Evaluate application model accuracy",
	Args:  cobra.NoArgs,
}

var nluCmd = &cobra.Command{
	Use:   "nlu",
	Short: "Evaluate the NLU accuracy of the given application model",
	Long:  "To run NLU evaluation, you need a set of ground truth annotations. Use the `annotate` command to get started.",
	Example: `speechly evaluate nlu <app_id> ground-truths.txt
speechly evaluate nlu <app_id> ground-truths.txt --reference-date 2021-01-20`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appID := args[0]
		refD, err := readReferenceDate(cmd)
		if err != nil {
			log.Fatalf("reading reference date flag failed: %v", err)
		}

		res, annotated, err := runThroughWLU(ctx, appID, args[1], false, refD)
		if err != nil {
			log.Fatalf("WLU failed: %v", err)
		}
		isRelaxed, err := cmd.Flags().GetBool("relax")
		if err != nil {
			log.Fatalf("WLU failed: %v", err)
		}

		evaluateAnnotatedUtterances(wluResponsesToString(res.Responses), annotated, isRelaxed)
	},
}

var asrCmd = &cobra.Command{
	Use:   "asr",
	Short: "Evaluate the ASR accuracy of the given application model",
	Long:  "To run ASR evaluation, you need a set of ground truth transcripts. Use the `transcribe` command to get started.",
	Example: `speechly evaluate asr <app_id> ground-truths.jsonl
speechly evaluate asr <app_id> ground-truths.jsonl --streaming`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appID := args[0]
		var ac []AudioCorpusItem
		useStreaming, err := cmd.Flags().GetBool("streaming")
		if err != nil {
			log.Fatalf("Reading streaming flag failed: %v", err)
		}
		if useStreaming {
			ac, err = transcribeWithStreamingAPI(ctx, appID, args[1], true)
		} else {
			ac, err = transcribeWithBatchAPI(ctx, appID, args[1], true)
		}
		if err != nil {
			log.Fatalf("Transcription failed: %v", err)
		}

		ed := EditDistance{}
		for _, aci := range ac {
			wd, err := wordDistance(aci.Transcript, aci.Hypothesis)
			if err != nil {
				log.Fatalf("Error in result generation: %v", err)
			}
			if wd.dist > 0 && wd.base > 0 {
				aln1, aln2, _ := nwalgo.Align(aci.Transcript, aci.Hypothesis, "*", 1, -1, -1)
				fmt.Printf("\nAudio: %s\n", aci.Audio)
				fmt.Printf("└─ Ground truth: %s\n", aln1)
				fmt.Printf("└─ Prediction:   %s\n", aln2)
			}
			ed = ed.Add(wd)
		}
		fmt.Printf("\nWord Error Rate (WER): %.2f (%.0d/%.0d)\n", ed.AsER(), ed.dist, ed.base)
	},
}

func init() {
	RootCmd.AddCommand(evaluateCmd)
	evaluateCmd.AddCommand(nluCmd)
	nluCmd.Flags().StringP("reference-date", "r", "", "Reference date in YYYY-MM-DD format, if not provided use current date.")
	nluCmd.Flags().Bool("relax", false, "Ignore normalized entity values and casing in matching.")

	evaluateCmd.AddCommand(asrCmd)
	asrCmd.Flags().Bool("streaming", false, "Use the Streaming API instead of the Batch API.")
}
