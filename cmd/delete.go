package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing application",
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalf("Missing force flag: %s", err)
		}

		dry, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatalf("Missing dry-run flag: %s", err)
		}

		id, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatalf("Missing app ID: %s", err)
		}

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
		apps, err := configClient.ListApps(ctx, &configv1.ListAppsRequest{Project: project})
		if err != nil {
			log.Fatalf("Getting apps for project %s failed: %s", project, err)
		}

		if appList := apps.GetApps(); len(appList) > 0 {
			if !appIdInAppList(id, appList) {
				cmd.Printf("App id '%s' does not exist. Your project has apps: \n", id)
				if err := printApps(cmd.OutOrStdout(), appList...); err != nil {
					log.Fatalf("Error listing app: %s", err)
				}
				return
			}
		} else {
			cmd.Println("No applications found.")
			return
		}

		if !force && !confirm(fmt.Sprintf("Deleting app %s, are you sure?", id), cmd.OutOrStdout(), cmd.InOrStdin()) {
			cmd.Println("Deletion aborted.")
			return
		}

		if !dry {
			if _, err := configClient.DeleteApp(
				ctx,
				&configv1.DeleteAppRequest{
					AppId: id,
				},
			); err != nil {
				log.Fatalf("Error deleting the app: %s", err)
			}
		}

		cmd.Printf("Successfully deleted app %s.\n", id)
	},
}

func init() {
	deleteCmd.Flags().StringP("app", "a", "", "application ID to delete")
	if err := deleteCmd.MarkFlagRequired("app"); err != nil {
		log.Fatalf("Internal error: %s", err)
	}

	deleteCmd.Flags().BoolP("force", "f", false, "skip confirmation prompt")
	deleteCmd.Flags().BoolP("dry-run", "d", false, "don't perform the deletion")

	rootCmd.AddCommand(deleteCmd)
}

func confirm(prompt string, dst io.Writer, src io.Reader) bool {
	read := bufio.NewReader(src)

	for {
		fmt.Fprintf(dst, "%s [y/n]: ", prompt)

		r, err := read.ReadString('\n')
		if err != nil {
			return false
		}

		r = strings.ToLower(strings.TrimSpace(r))
		if r == "y" || r == "yes" {
			return true
		} else if r == "n" || r == "no" {
			return false
		}
	}
}

func appIdInAppList(appId string, apps []*configv1.App) bool {
	for _, app := range apps {
		if app.GetId() == appId {
			return true
		}
	}
	return false
}
