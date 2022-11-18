package cmd

import (
	"fmt"
	"io"
	"log"
	"strings"
	"text/tabwriter"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List applications in the current project",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		configClient, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}

		projects, err := configClient.GetProject(ctx, &configv1.GetProjectRequest{})
		if err != nil {
			log.Fatalf("Getting projects failed: %s", err)
		}
		project := projects.Project[0]
		projectName := projects.ProjectNames[0]
		apps, err := configClient.ListApps(ctx, &configv1.ListAppsRequest{Project: project})
		if err != nil {
			log.Fatalf("Listing apps for project %s failed: %s", project, err)
		}
		cmd.Printf("Applications in project \"%s\" (%s):\n\n", projectName, project)
		if a := apps.GetApps(); len(a) > 0 {
			if err := printApps(cmd.OutOrStdout(), a...); err != nil {
				log.Fatalf("Error listing apps: %s", err)
			}
		} else {
			cmd.Printf("No applications found.\n")
		}

		// If the project name in settings is automatically generated, update it from server.
		conf := clients.GetConfig(cmd.Context())
		currentContextName := viper.Get("current-context")
		if projectName != currentContextName {
			idxToUpdate := -1
			for idx, c := range conf.Contexts {
				if c.Name == currentContextName {
					idxToUpdate = idx
					break
				}
			}
			if idxToUpdate >= 0 {
				if conf.Contexts[idxToUpdate].Name == conf.Contexts[idxToUpdate].RemoteName {
					conf.Contexts[idxToUpdate].Name = projectName
					viper.Set("current-context", projectName)
				}
				conf.Contexts[idxToUpdate].RemoteName = projectName
				viper.Set("contexts", conf.Contexts)
				if err := viper.WriteConfig(); err != nil {
					log.Fatalf("Failed to write settings: %s", err)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}

func printApps(out io.Writer, apps ...*configv1.App) error {
	// Format in tab-separated columns with a tab stop of 8.
	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	fmt.Fprint(w, "Name\tApp ID\tStatus\tDeployed at\n")
	for _, app := range apps {
		deployedAt := ""
		if app.GetDeployedAtTime() != nil {
			deployedAt = app.GetDeployedAtTime().AsTime().String()
		}
		status := app.GetStatus().String()
		if strings.Contains(status, "STATUS_") {
			c := cases.Title(language.English)
			status = c.String(strings.TrimPrefix(status, "STATUS_"))
		}
		fmt.Fprintf(w, "%-*.*s\t%s\t%s\t%s\n", 16, 32, app.GetName(), app.GetId(), status, deployedAt)
	}

	return w.Flush()
}
