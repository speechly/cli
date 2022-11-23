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
	Example: `speechly delete <app_id>
speechly delete --app <app_id> --force`,
	Args: cobra.RangeArgs(0, 1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		appId, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatalf("Missing app ID: %s", err)
		}
		if appId == "" && len(args) < 1 {
			return fmt.Errorf("app_id must be given with flag --app or as the sole positional argument")
		}
		return nil
	},
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
		if id == "" {
			id = args[0]
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
		projectName := projects.ProjectNames[0]
		apps, err := configClient.ListApps(ctx, &configv1.ListAppsRequest{Project: project})
		if err != nil {
			log.Fatalf("Getting apps for project %s failed: %s", project, err)
		}

		if appList := apps.GetApps(); len(appList) > 0 {
			if !appIdInAppList(id, appList) {
				cmd.Printf("App ID '%s' does not exist.\n\nApplications in project \"%s\" (%s):\n\n", id, projectName, project)
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
	deleteCmd.Flags().StringP("app", "a", "", "Application to delete. Can be given as the sole positional argument.")
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt.")
	deleteCmd.Flags().BoolP("dry-run", "d", false, "Don't perform the deletion.")

	RootCmd.AddCommand(deleteCmd)
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
