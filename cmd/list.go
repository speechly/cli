package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
)

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List applications in the current context (project)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		projects, err := client.GetProject(ctx, &configv1.GetProjectRequest{})
		if err != nil {
			log.Fatalf("Getting projects failed: %s", err)
		}
		project := projects.Project[0]
		apps, err := client.ListApps(ctx, &configv1.ListAppsRequest{Project: project})
		if err != nil {
			log.Fatalf("Listing apps for project %s failed: %s", project, err)
		}
		cmd.Printf("List of applications in project %s:\n\n", project)
		if len(apps.Apps) > 0 {
			cmd.Printf("APP ID\t\t\t\t\tSTATUS\t\tNAME\n")
			for _, app := range apps.Apps {
				cmd.Println(fmt.Sprintf("%s\t%s\t%s", app.Id, app.Status, app.Name))
			}
		} else {
			cmd.Printf("No applications found.\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
