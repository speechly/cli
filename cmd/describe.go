package cmd

import (
	"log"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
)

var describeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Print details about an application",
	Example: `speechly describe <app_id>
speechly describe --app <app_id>`,
	Args:    cobra.RangeArgs(0, 1),
	PreRunE: checkSoleAppArgument,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		wait, _ := cmd.Flags().GetBool("watch")
		if appId == "" {
			appId = args[0]
		}

		configClient, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}
		app, err := configClient.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to get app %s: %s", appId, err)
		}
		deployedAt := "Not available"
		if app.App.DeployedAtTime != nil {
			deployedAt = app.App.DeployedAtTime.AsTime().String()
		}
		cmd.Printf("App ID: %s\n", app.App.Id)
		cmd.Printf("Name: %s\n", app.App.Name)
		cmd.Printf("Language: %s\n", app.App.Language)

		waitFor := configv1.App_STATUS_UNSPECIFIED
		if wait {
			waitFor = configv1.App_STATUS_TRAINED
		}
		waitForAppStatus(cmd, configClient, appId, waitFor)
		cmd.Printf("Deployed at: %s\n", deployedAt)
	},
}

func init() {
	RootCmd.AddCommand(describeCmd)
	describeCmd.Flags().StringP("app", "a", "", "Application to describe. Can be given as the sole positional argument.")
	describeCmd.Flags().BoolP("watch", "w", false, "If app status is training, wait until it is finished.")
}
