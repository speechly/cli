package cmd

import (
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

func init() {
	RootCmd.AddCommand(evaluateCmd)
	evaluateCmd.AddCommand(nluCmd)

	nluCmd.Flags().StringP("reference-date", "r", "", "Reference date in YYYY-MM-DD format, if not provided use current date.")
}
