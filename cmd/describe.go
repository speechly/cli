package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Print details about an application",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		app, err := client.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to get app %s: %s", appId, err)
		}
		cmd.Printf("AppId:\t%s\n", app.App.Id)
		cmd.Printf("Name:\t%s\n", app.App.Name)
		cmd.Printf("Lang:\t%s\n", app.App.Language)
		cmd.Printf("Status:\t%s", app.App.Status)
		if app.App.Status == configv1.App_STATUS_FAILED {
			cmd.Printf("\t%s", app.App.ErrorMsg)
		} else if app.App.Status == configv1.App_STATUS_TRAINING {
			cmd.Printf("\testimated time remaining: ")
			if app.App.EstimatedRemainingSec > 0 {
				cmd.Printf("%d seconds", app.App.EstimatedRemainingSec)
			} else {
				cmd.Printf("unknown")
			}
		}
		cmd.Printf("\n")
		// if watch flag given, remain here and fetech app state in loop
		wait, _ := cmd.Flags().GetBool("watch")
		for wait && app.App.Status == configv1.App_STATUS_TRAINING {
			time.Sleep(10 * time.Second)
			app, err = client.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
			if err != nil {
				log.Fatalf("Failed to refresh app %s: %s", appId, err)
			}
			cmd.Printf("Status:\t%s", app.App.Status)
			if app.App.Status == configv1.App_STATUS_TRAINING {
				cmd.Printf("\testimated time remaining: ")
				if app.App.EstimatedRemainingSec > 0 {
					cmd.Printf("%d seconds", app.App.EstimatedRemainingSec)
				} else {
					cmd.Printf("unknown")
				}
			}
			cmd.Printf("\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("app", "a", "", "Application id to describe")
	describeCmd.MarkFlagRequired("app")
	describeCmd.Flags().BoolP("watch", "w", false, "If app status is training, wait until it is finished.")
}
