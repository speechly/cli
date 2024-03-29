package cmd

import (
	"context"
	"fmt"
	"log"
	"regexp"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/speechly/cli/pkg/clients"
)

var configCmd = &cobra.Command{
	Use:     "projects [command]",
	Aliases: []string{"project", "config"},
	Short:   "Manage API access to Speechly projects",
	Args:    cobra.NoArgs,
}

var validName = regexp.MustCompile(`[^A-Za-z0-9]+`)

var configListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List known projects",
	Run: func(cmd *cobra.Command, args []string) {
		conf := clients.GetConfig(cmd.Context())
		cmd.Printf("Settings file used: %s\n", viper.ConfigFileUsed())
		cmd.Printf("Known projects:\n")
		for _, c := range conf.Contexts {
			prefix := "   "
			if c.Name == conf.CurrentContext {
				prefix = "✔  "
			}

			if c.Name == c.RemoteName {
				cmd.Printf("%s%s\n", prefix, c.Name)
			} else if c.RemoteName == "" {
				cmd.Printf("%s%s (name unknown)\n", prefix, c.Name)
			} else {
				cmd.Printf("%s%s (%s)\n", prefix, c.Name, c.RemoteName)
			}
		}
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add access to a pre-existing project",
	Example: `speechly projects add <api_token>
speechly projects add --apikey <api_token>`,
	Args: cobra.RangeArgs(0, 1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if name != "" && validName.MatchString(name) {
			return fmt.Errorf("invalid name: %s", name)
		}
		apikey, _ := cmd.Flags().GetString("apikey")
		if apikey == "" {
			if len(args) == 0 {
				return fmt.Errorf("apikey must be given either with --apikey flag or as the sole positional argument")
			}
		}

		conf := clients.GetConfig(cmd.Context())
		for _, c := range conf.Contexts {
			if c.Name == name {
				return fmt.Errorf("project with name %s already exists", name)
			}
			if c.Apikey == apikey {
				return fmt.Errorf("project with given apikey already exists")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf := clients.GetConfig(cmd.Context())
		host, _ := cmd.Flags().GetString("host")
		apikey, _ := cmd.Flags().GetString("apikey")
		if apikey == "" {
			apikey = args[0]
		}
		name, _ := cmd.Flags().GetString("name")
		isUserDefinedName := true
		if name == "" {
			name = apikey
			isUserDefinedName = false
		}
		previousContextName := viper.Get("current-context")
		viper.Set("contexts", append(conf.Contexts, clients.SpeechlyContext{Host: host, Apikey: apikey, Name: name, RemoteName: ""}))
		viper.Set("current-context", name)
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Failed to write settings: %s", err)
		}

		skipValidation, err := cmd.Flags().GetBool("skip-online-validation")
		if err != nil {
			log.Fatalf("Missing skip-online-validation flag: %s", err)
		}

		if !skipValidation {
			ctx := clients.NewContext(failWithError)
			configClient, err := clients.ConfigClient(ctx)
			if err != nil {
				log.Fatalf("Error connecting to API: %s", err)
			}

			projects, err := configClient.GetProject(ctx, &configv1.GetProjectRequest{})
			if err != nil {
				log.Fatalf("Verifying api token failed: %s", err)
			}
			projectName := projects.ProjectNames[0]
			viper.Set("current-context", previousContextName)
			for i, c := range conf.Contexts {
				if c.Name == name {
					conf.Contexts = append(conf.Contexts[:i], conf.Contexts[i+1:]...)
				}
			}
			viper.Set("contexts", conf.Contexts)
			actualName := projectName
			if isUserDefinedName {
				actualName = name
			} else {
				for _, c := range conf.Contexts {
					if actualName == c.Name {
						actualName = fmt.Sprintf("%s (%d)", projectName, len(conf.Contexts))
					}
				}
			}

			viper.Set("contexts", append(conf.Contexts, clients.SpeechlyContext{Host: host, Apikey: apikey, Name: actualName, RemoteName: projectName}))
			viper.Set("current-context", actualName)
			if err := viper.WriteConfig(); err != nil {
				log.Fatalf("Failed to write settings: %s", err)
			}
		}

		cmd.Printf("Wrote settings to file: %s\n", viper.ConfigFileUsed())
	},
}

var configRemoveCmd = &cobra.Command{
	Use:     "remove",
	Aliases: []string{"rm"},
	Short:   "Remove access to a project",
	Example: `speechly projects remove --name <project_name>`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if err := ensureContextExists(cmd.Context(), name); err != nil {
			return err
		}
		if name == viper.Get("current-context") {
			return fmt.Errorf("cannot remove active project: %s", name)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf := clients.GetConfig(cmd.Context())
		name, _ := cmd.Flags().GetString("name")
		cmd.Printf("Removing access to project: %s\n", name)
		for i, c := range conf.Contexts {
			if c.Name == name {
				conf.Contexts = append(conf.Contexts[:i], conf.Contexts[i+1:]...)
			}
		}
		viper.Set("contexts", conf.Contexts)
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Failed to write settings: %s", err)
		}
		cmd.Printf("Wrote settings to file: %s\n", viper.ConfigFileUsed())
	},
}

var configUseCmd = &cobra.Command{
	Use:     "use",
	Aliases: []string{"switch"},
	Short:   "Select the default project used",
	Example: `speechly projects use
speechly projects use --name <project_name>`,
	Run: func(cmd *cobra.Command, args []string) {
		previousContext := viper.Get("current-context")
		name, _ := cmd.Flags().GetString("name")
		if name == "" {
			conf := clients.GetConfig(cmd.Context())
			cmd.Printf("Select a project to use: (Enter number)\n")
			for ixd, c := range conf.Contexts {
				prefix := fmt.Sprintf("%d. ", ixd+1)
				if c.Name == previousContext {
					prefix = "✔  "
				}
				if c.Name == c.RemoteName {
					cmd.Printf("%s%s\n", prefix, c.Name)
				} else if c.RemoteName == "" {
					cmd.Printf("%s%s (name unknown)\n", prefix, c.Name)
				} else {
					cmd.Printf("%s%s (%s)\n", prefix, c.Name, c.RemoteName)
				}
			}

			cmd.Printf("Project: ")
			var i int
			_, err := fmt.Scanf("%d", &i)
			if err != nil {
				log.Fatalf("Invalid choice, not a number in range (1 – %d)", len(conf.Contexts))
			}
			if 1 > i || i > len(conf.Contexts) {
				log.Fatalf("Invalid choice, %d is not a number in range (1 – %d)", i, len(conf.Contexts))
			}
			name = conf.Contexts[i-1].Name
		}
		if err := ensureContextExists(cmd.Context(), name); err != nil {
			log.Fatalf("Unknown context %s", name)
		}

		viper.Set("current-context", name)
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Failed to write settings: %s", err)
		}
		cmd.Printf("Wrote settings to file: %s\n", viper.ConfigFileUsed())
	},
}

func ensureContextExists(ctx context.Context, name string) error {
	conf := clients.GetConfig(ctx)
	for _, c := range conf.Contexts {
		if c.Name == name {
			return nil
		}
	}
	return fmt.Errorf("project named %s is not known", name)
}

func init() {
	configCmd.AddCommand(configListCmd)

	configAddCmd.Flags().String("apikey", "", "API token, created in Speechly Dashboard. Can also be given as the sole positional argument.")
	configAddCmd.Flags().String("name", "", "An unique name for the project. If not given the project name configured in Speechly Dashboard will be used.")
	configAddCmd.Flags().String("host", "api.speechly.com", "API address")
	configAddCmd.Flags().Bool("skip-online-validation", false, "Skips validating the API token against the host.")
	configCmd.AddCommand(configAddCmd)

	configRemoveCmd.Flags().String("name", "", "The name for the project for which access is to be removed.")
	if err := configRemoveCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configCmd.AddCommand(configRemoveCmd)

	configUseCmd.Flags().String("name", "", "An unique name for the project.")
	configCmd.AddCommand(configUseCmd)

	RootCmd.AddCommand(configCmd)
}
