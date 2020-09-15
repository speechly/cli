package cmd

import (
	"fmt"
	"io"
	"log"
	"text/tabwriter"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List applications in the current context (project)",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		projects, err := config_client.GetProject(ctx, &configv1.GetProjectRequest{})
		if err != nil {
			log.Fatalf("Getting projects failed: %s", err)
		}
		project := projects.Project[0]
		apps, err := config_client.ListApps(ctx, &configv1.ListAppsRequest{Project: project})
		if err != nil {
			log.Fatalf("Listing apps for project %s failed: %s", project, err)
		}
		cmd.Printf("List of applications in project %s:\n\n", project)
		if a := apps.GetApps(); len(a) > 0 {
			if err := printApps(cmd.OutOrStdout(), a...); err != nil {
				log.Fatalf("Error listing apps: %s", err)
			}
		} else {
			cmd.Printf("No applications found.\n")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func printApps(out io.Writer, apps ...*configv1.App) error {
	// Format in tab-separated columns with a tab stop of 8.
	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)

	fmt.Fprint(w, "APP ID\tSTATUS\tNAME\n")
	for _, app := range apps {
		fmt.Fprintf(w, "%s\t%s\t%s\n", app.GetId(), app.GetStatus(), app.GetName())
	}

	return w.Flush()
}
