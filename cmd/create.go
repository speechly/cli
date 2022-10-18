package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
)

var createCmd = &cobra.Command{
	Use:   "create [<application name>]",
	Short: "Create a new application in the current project",
	Long:  "Creates a new application in the current project and a configuration file in the current working directory.",
	Args:  cobra.RangeArgs(0, 1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name == "" && len(args) == 0 {
			{
				return fmt.Errorf("name must be given either with flag --name or as the sole positional parameter")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("Missing name: %s", err)
		}
		if name == "" {
			name = args[0]
		}

		lang, err := cmd.Flags().GetString("language")
		if err != nil {
			log.Fatalf("Missing language: %s", err)
		}
		if lang == "" {
			lang = "en-US"
		}

		ctx := cmd.Context()

		configClient, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}
		projects, err := configClient.GetProject(ctx, &configv1.GetProjectRequest{})
		if err != nil {
			log.Fatalf("Error fetching projects: %s", err)
		}

		if len(projects.Project) < 1 {
			log.Fatal("Error fetching projects: no projects exist for the given token")
		}

		projectId := projects.Project[0]
		projectName := projects.ProjectNames[0]

		a := &configv1.App{
			Language: lang,
			Name:     name,
		}
		req := &configv1.CreateAppRequest{
			Project: projectId,
			App:     a,
		}

		res, err := configClient.CreateApp(ctx, req)
		if err != nil {
			log.Fatalf("Error creating an app: %s", err)
		}

		// Cannot use the response here, because it only contains the id.
		a.Id = res.GetApp().GetId()

		path, err := os.Getwd()
		if err != nil {
			log.Fatalf("Could not access current folder: %s", err)
		}
		buf := []byte(fmt.Sprintf("lang: %s\ntemplates: ''\nintents: []\nentities: []\n", lang))
		out := filepath.Join(path, "config.yaml")
		log.Printf("Writing file %s (%d bytes)\n", out, len(buf))
		if err := os.WriteFile(out, buf, 0644); err != nil {
			log.Fatalf("Could not write configuration to %s: %s", out, err)
		}

		cmd.Printf("Created an application in project \"%s\":\n\n", projectName)
		if err := printApps(cmd.OutOrStdout(), a); err != nil {
			log.Fatalf("Error listing app: %s", err)
		}
	},
}

func init() {
	createCmd.Flags().StringP("language", "l", "en-US", "Application language. Current only 'en-US' and 'fi-FI' are supported.")
	createCmd.Flags().StringP("name", "n", "", "Application name, can alternatively be given as the sole positional argument.")
	RootCmd.AddCommand(createCmd)
}
