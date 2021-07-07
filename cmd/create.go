package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
)

var validLangs = []string{"en-US", "fi-FI"}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new application in the current context (project)",
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("Missing name: %s", err)
		}

		lang, err := cmd.Flags().GetString("language")
		if err != nil {
			log.Fatalf("Missing language: %s", err)
		}

		if err := validateLang(lang); err != nil {
			log.Fatalf("Invalid language: %s", err)
		}

		ctx := cmd.Context()

		config_client, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}
		projects, err := config_client.GetProject(ctx, &configv1.GetProjectRequest{})
		if err != nil {
			log.Fatalf("Error fetching projects: %s", err)
		}

		if len(projects.GetProject()) < 1 {
			log.Fatal("Error fetching projects: no projects exist for the given token")
		}

		pid := projects.GetProject()[0]
		a := &configv1.App{
			Language: lang,
			Name:     name,
		}
		req := &configv1.CreateAppRequest{
			Project: pid,
			App:     a,
		}

		res, err := config_client.CreateApp(ctx, req)
		if err != nil {
			log.Fatalf("Error creating an app: %s", err)
		}

		// Cannot use the response here, because it only contains the id.
		a.Id = res.GetApp().GetId()

		cmd.Printf("Created an application in project %s:\n\n", pid)
		if err := printApps(cmd.OutOrStdout(), a); err != nil {
			log.Fatalf("Error listing app: %s", err)
		}
	},
}

func init() {
	createCmd.Flags().StringP("language", "l", "", "application language (current only 'en-US' and 'fi-FI' are supported)")
	if err := createCmd.MarkFlagRequired("language"); err != nil {
		log.Fatalf("Internal error: %s", err)
	}

	createCmd.Flags().StringP("name", "n", "", "application name")
	if err := createCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("Internal error: %s", err)
	}

	rootCmd.AddCommand(createCmd)
}

func validateLang(l string) error {
	for _, v := range validLangs {
		if v == l {
			return nil
		}
	}

	return fmt.Errorf("unsupported language: %s", l)
}
