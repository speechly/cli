package cmd

import (
	"log"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing application",
	Example: `speechly edit --app <app_id> --name <new_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")

		if name == "" {
			log.Println("Nothing to do.")
			return
		}
		appId, _ := cmd.Flags().GetString("app")

		ctx := cmd.Context()
		configClient, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}
		appRes, err := configClient.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to get app %s: %s", appId, err)
		}
		app := appRes.GetApp()

		app.Name = name

		_, err = configClient.UpdateApp(ctx, &configv1.UpdateAppRequest{App: app})
		if err != nil {
			log.Fatalf("Error editing application: %s", err)
		}
		cmd.Println("Updated application:")
		cmd.Printf("AppId:\t%s\n", app.Id)
		cmd.Printf("Name:\t%s\n", app.Name)
		cmd.Printf("Lang:\t%s\n", app.Language)
	},
}

func init() {
	editCmd.Flags().StringP("app", "a", "", "Application to edit")
	if err := editCmd.MarkFlagRequired("app"); err != nil {
		log.Fatalf("Internal error: %s", err)
	}
	editCmd.Flags().StringP("name", "n", "", "Application name")

	RootCmd.AddCommand(editCmd)
}
