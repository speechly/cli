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
	Use: "create [<application name>]",
	Example: `speechly create "My App"
speechly create --name "My App" --output-dir /foo/bar
`,
	Short: "Create a new application in the current project",
	Long:  "Creates a new application in the current project and a config file in the current working directory.",
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

		outDir, _ := cmd.Flags().GetString("output-dir")

		if outDir == "" || outDir == "." {
			outDir = "."
		} else {
			outDir, _ = filepath.Abs(outDir)
			if _, err := os.Stat(outDir); os.IsNotExist(err) {
				if err := os.Mkdir(outDir, os.ModePerm); err != nil {
					log.Fatalf("Could not create the output directory %s: %s", outDir, err)
				}
			} else {
				log.Fatalf("Directory %s already exists", outDir)
			}
		}

		buf := []byte("templates: ''\nintents: []\nentities: []\n")
		outFile := filepath.Join(outDir, "config.yaml")
		log.Printf("Writing file %s (%d bytes)\n", outFile, len(buf))
		if err := os.WriteFile(outFile, buf, 0644); err != nil {
			log.Fatalf("Could not write configuration to %s: %s", outFile, err)
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

		cmd.Printf("Created an application in project \"%s\":\n\n", projectName)
		if err := printApps(cmd.OutOrStdout(), a); err != nil {
			log.Fatalf("Error listing app: %s", err)
		}
	},
}

func init() {
	createCmd.Flags().StringP("language", "l", "en-US", "Application language. Available options are 'en-US' and 'fi-FI'.")
	createCmd.Flags().StringP("name", "n", "", "Application name. Can be given as the sole positional argument.")
	createCmd.Flags().StringP("output-dir", "o", "", "Output directory for the config file.")
	RootCmd.AddCommand(createCmd)
}
