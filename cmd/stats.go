package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"text/tabwriter"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"

	analyticsv1 "github.com/speechly/api/go/speechly/analytics/v1"
	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Get utterance statistics for the current project or an application in it",
	Example: `speechly stats <app_id>
speechly stats --app <app_id>
speechly stats > output.csv
speechly stats --start-date 2021-03-01 --end-date 2021-04-01`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		appId, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatal("Error reading flags")
		}
		if appId == "" && len(args) == 1 {
			appId = args[0]
		}

		ctx := cmd.Context()
		configClient, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}
		analyticsClient, err := clients.AnalyticsClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}

		agg := analyticsv1.Aggregation_AGGREGATION_HOURLY
		req := &analyticsv1.UtteranceStatisticsRequest{
			Aggregation: agg,
		}
		if appId != "" {
			req.AppId = appId
		}
		startDate, err := cmd.Flags().GetString("start-date")
		if err != nil {
			log.Fatalf("start-date is invalid: %v", err)
		}
		req.StartDate = startDate
		endDate, err := cmd.Flags().GetString("end-date")
		if err != nil {
			log.Fatalf("end-date is invalid: %v", err)
		}
		req.EndDate = endDate
		export, err := cmd.Flags().GetBool("export")
		if err != nil {
			log.Fatalf("export flag is invalid: %v", err)
		}

		projects, err := configClient.GetProject(ctx, &configv1.GetProjectRequest{})
		if err != nil {
			log.Fatalf("Getting projects failed: %s", err)
		}
		projectId := projects.Project[0]

		res, err := analyticsClient.UtteranceStatistics(ctx, req)
		if err != nil {
			log.Fatalf("Getting statistics failed: %v", err)
		}

		if isatty.IsTerminal(os.Stdout.Fd()) && !export {
			cmd.Printf("Project ID: %s\n", projectId)
			if appId != "" {
				cmd.Printf("App ID: %s\n", appId)
			}
			cmd.Printf("Aggregation: %s\n", agg)
			cmd.Printf("Start time: %s\n", res.GetStartDate())
			cmd.Printf("End time: %s\n", res.GetEndDate())
			cmd.Printf("Total utterances: %d\n", res.GetTotalUtterances())
			cmd.Printf("Total duration: %d seconds\n", res.GetTotalDurationSeconds())
			if s := res.GetItems(); len(s) > 0 {
				if err := printAnalytics(cmd.OutOrStdout(), agg, s...); err != nil {
					log.Fatalf("Error printing statistics: %v", err)
				}
			}
		} else {
			if err := printCSV(cmd.OutOrStdout(), agg, res.GetItems()...); err != nil {
				log.Fatalf("Error creating CSV: %s", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(statsCmd)
	statsCmd.Flags().StringP("app", "a", "", "Application to get the statistics for. Can be given as the sole positional argument.")
	statsCmd.Flags().String("start-date", "", "Start date for statistics.")
	statsCmd.Flags().String("end-date", "", "End date for statistics, not included in results.")
	statsCmd.Flags().Bool("export", false, "Print report as CSV")
}

func printAnalytics(out io.Writer, agg analyticsv1.Aggregation, items ...*analyticsv1.UtteranceStatisticsPeriod) error {
	appId := ""
	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	for _, i := range items {
		if appId != i.GetAppId() {
			appId = i.GetAppId()
			fmt.Fprintf(w, "\n%s\n", appId)
			fmt.Fprint(w, "\tTIME\tUTTERANCE COUNT\tTOTAL AUDIO\tANNOTATED AUDIO\n")
		}
		fmt.Fprintf(w, "\t%s\t%d\t%d\t%d\n", formatDate(i.GetStartTime(), agg), i.GetCount(), i.GetUtterancesSeconds(), i.GetAnnotatedSeconds())
	}

	return w.Flush()
}

func printCSV(out io.Writer, agg analyticsv1.Aggregation, items ...*analyticsv1.UtteranceStatisticsPeriod) error {
	w := csv.NewWriter(out)
	if err := w.Write([]string{"APP", "START", "COUNT", "SECONDS", "ANNOTATED"}); err != nil {
		return err
	}
	for _, i := range items {
		if err := w.Write([]string{
			i.GetAppId(),
			formatDate(i.GetStartTime(), agg),
			strconv.Itoa(int(i.GetCount())),
			strconv.Itoa(int(i.GetUtterancesSeconds())),
			strconv.Itoa(int(i.GetAnnotatedSeconds())),
		}); err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}

func formatDate(ds string, agg analyticsv1.Aggregation) string {
	switch agg {
	case analyticsv1.Aggregation_AGGREGATION_DAILY:
		d, err := time.Parse(time.RFC3339, ds)
		if err != nil {
			log.Fatalf("invalid date: %s", err)
		}
		return d.Format("2006-01-02")
	default:
		return ds
	}
}
