package cmd

import (
	"fmt"
	"log"

	analyticsv1 "github.com/speechly/api/go/speechly/analytics/v1"
	"github.com/spf13/cobra"

	"github.com/speechly/cli/pkg/clients"
)

var utterancesCmd = &cobra.Command{
	Use:   "utterances <app_id>",
	Short: "Get a sample of recent utterances.",
	Long:  `Fetches a sample of recent utterances and their SAL-annotated transcript.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId := args[0]

		client, err := clients.AnalyticsClient(ctx)
		if err != nil {
			log.Fatalf("Failed to init analytics client%s", err)
		}
		response, err := client.Utterances(ctx, &analyticsv1.UtterancesRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to fetch utterances data for %s: %s", appId, err)
		}
		for _, utt := range response.Utterances {
			fmt.Printf("%s\t%s\t%s\n", utt.Date, utt.Annotated, utt.Transcript)
		}

	},
}

func init() {
	RootCmd.AddCommand(utterancesCmd)
}
