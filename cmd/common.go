package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	salv1 "github.com/speechly/api/go/speechly/sal/v1"
	"github.com/spf13/cobra"
)

func printLineErrors(messages []*salv1.LineReference) {
	log.Println("Configuration validation failed")
	for _, message := range messages {
		var errorLevel string
		switch message.Level {
		case salv1.LineReference_LEVEL_NOTE:
			errorLevel = "NOTE"
		case salv1.LineReference_LEVEL_WARNING:
			errorLevel = "WARNING"
		case salv1.LineReference_LEVEL_ERROR:
			errorLevel = "ERROR"
		}
		if message.File != "" {
			log.Printf("%s:%d:%d:%s:%s\n", message.File, message.Line,
				message.Column, errorLevel, message.Message)
		} else {
			log.Printf("%s: %s", errorLevel, message.Message)
		}
	}
	os.Exit(1)
}

func waitForAppStatus(cmd *cobra.Command, configClient configv1.ConfigAPIClient, appId string, status configv1.App_Status) {
	ctx := cmd.Context()

	for {
		app, err := configClient.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to refresh app %s: %s", appId, err)
		}
		cmd.Printf("Status:\t%s", app.App.Status)
		switch app.App.Status {
		case configv1.App_STATUS_NEW:
			cmd.Printf(", queued (%d jobs before this)", app.App.QueueSize)
		case configv1.App_STATUS_TRAINING:
			age := time.Duration(app.App.TrainingTimeSec) * time.Second
			est := time.Duration(app.App.EstimatedTrainingTimeSec) * time.Second
			cmd.Printf(", age %s, estimated about %s", age, est)
		case configv1.App_STATUS_FAILED:
			cmd.Println()
			cmd.Printf("Error: %s", app.App.ErrorMsg)
		}
		cmd.Println()

		if app.App.Status >= status {
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func checkSoleAppArgument(cmd *cobra.Command, args []string) error {
	appId, err := cmd.Flags().GetString("app")
	if err != nil {
		log.Fatalf("Missing app ID: %s", err)
	}
	if appId == "" && len(args) < 1 {
		return fmt.Errorf("app_id must be given with flag --app or as the sole positional argument")
	}
	return nil
}
