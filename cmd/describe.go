package cmd

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
)

func formatSeconds(seconds int32) string {
    m := seconds / 60
    s := seconds % 60
    if m == 0 {
	return fmt.Sprintf("%02ds", s)
    }
    return fmt.Sprintf("%02dm%02ds", m, s)
}

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Print details about an application",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		app, err := config_client.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to get app %s: %s", appId, err)
		}
		cmd.Printf("AppId:\t%s\n", app.App.Id)
		cmd.Printf("Name:\t%s\n", app.App.Name)
		cmd.Printf("Lang:\t%s\n", app.App.Language)
		if app.App.Status == configv1.App_STATUS_TRAINED {
			cmd.Printf("Status:\t%s\n", app.App.Status)
		}
		if app.App.Status == configv1.App_STATUS_FAILED {
			cmd.Printf("Status:\t%s\n", app.App.ErrorMsg)
		} else if app.App.Status == configv1.App_STATUS_TRAINING {
			if app.App.TrainingTimeSec > 0 {
				age := formatSeconds(app.App.TrainingTimeSec)
				cmd.Printf("Status:\t%s, age %s", app.App.Status, age)
				if app.App.EstimatedTrainingTimeSec > 0 {
					estim := (app.App.EstimatedTrainingTimeSec / 60) + 1
					cmd.Printf(", estimated about %02dm\n", estim)
				}
			}

			// if watch flag given, remain here and fetech app state in loop
			wait, _ := cmd.Flags().GetBool("watch")
			if wait {
				waitForDeploymentFinished(ctx, appId)
			}
		}
	},
}

func waitForDeploymentFinished(ctx context.Context, appId string) {
	time.Sleep(5 * time.Second)
	app, err := config_client.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
	if err != nil {
		log.Fatalf("Failed to get app %s: %s", appId, err)
	}

	for app.App.Status == configv1.App_STATUS_TRAINING {
		app, err = config_client.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to refresh app %s: %s", appId, err)
		}
		if app.App.Status == configv1.App_STATUS_TRAINING {
			r := "unknown"
			if app.App.EstimatedRemainingSec > 0 {
				r = fmt.Sprintf("%d seconds", app.App.EstimatedRemainingSec)
			}
			log.Println(fmt.Sprintf("Status:\t%s\testimated time remaining: %s", app.App.Status, r))
			time.Sleep(10 * time.Second)
		}
	}

	log.Println(fmt.Sprintf("Status:\t%s", app.App.Status))
}

func init() {
	rootCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("app", "a", "", "Application id to describe")
	describeCmd.MarkFlagRequired("app")
	describeCmd.Flags().BoolP("watch", "w", false, "If app status is training, wait until it is finished.")
}
