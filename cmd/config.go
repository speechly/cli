package cmd

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/speechly/cli/pkg/clients"
)

var configCmd = &cobra.Command{
	Use:     "projects",
	Aliases: []string{"config"},
	Short:   "Manage API access to Speechly projects",
	Run: func(cmd *cobra.Command, args []string) {
		conf := clients.GetConfig(cmd.Context())
		cmd.Printf("Configuration file used: %s\n", viper.ConfigFileUsed())
		cmd.Printf("Current project: %s\n", conf.CurrentContext)
		cmd.Printf("Known projects:\n")
		for _, c := range conf.Contexts {
			cmd.Printf("- %s\n", c.Name)
		}
	},
}

var validName = regexp.MustCompile(`[^A-Za-z0-9]+`)

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Configure access to a pre-existing project",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if validName.MatchString(name) {
			return fmt.Errorf("Invalid name: %s", name)
		}
		conf := clients.GetConfig(cmd.Context())
		for _, c := range conf.Contexts {
			if c.Name == name {
				return fmt.Errorf("Context named %s already exists", name)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		conf := clients.GetConfig(cmd.Context())
		host, _ := cmd.Flags().GetString("host")
		apikey, _ := cmd.Flags().GetString("apikey")
		name, _ := cmd.Flags().GetString("name")
		viper.Set("contexts", append(conf.Contexts, clients.SpeechlyContext{Host: host, Apikey: apikey, Name: name}))
		viper.Set("current-context", name)
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Failed to write the config: %s", err)
		}
		cmd.Printf("Wrote configuration to file: %s\n", viper.ConfigFileUsed())
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove access to a project from configuration",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if err := ensureContextExists(cmd.Context(), name); err != nil {
			return err
		}
		if name == viper.Get("current-context") {
			return fmt.Errorf("Cannot remove active context: %s", name)
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
			log.Fatalf("Failed to write the config: %s", err)
		}
		cmd.Printf("Wrote configuration to file: %s\n", viper.ConfigFileUsed())
	},
}

var configUseCmd = &cobra.Command{
	Use:   "use",
	Short: "Select the default project used",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		return ensureContextExists(cmd.Context(), name)
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		viper.Set("current-context", name)
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Failed to write the config: %s", err)
		}
		cmd.Printf("Wrote configuration to file: %s\n", viper.ConfigFileUsed())
	},
}

func ensureContextExists(ctx context.Context, name string) error {
	conf := clients.GetConfig(ctx)
	for _, c := range conf.Contexts {
		if c.Name == name {
			return nil
		}
	}
	return fmt.Errorf("Context named %s not found in configuration", name)
}

func init() {
	configAddCmd.Flags().String("apikey", "", "API key, created in dashboard")
	if err := configAddCmd.MarkFlagRequired("apikey"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configAddCmd.Flags().String("name", "", "An unique name for the project. If not given the project name (or id) will be used.")
	if err := configAddCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configAddCmd.Flags().String("host", "api.speechly.com", "API address")
	configCmd.AddCommand(configAddCmd)

	configRemoveCmd.Flags().String("name", "", "The name for the project for which access is to be removed")
	if err := configRemoveCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configCmd.AddCommand(configRemoveCmd)

	configUseCmd.Flags().String("name", "", "An unique name for the project")
	if err := configUseCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configCmd.AddCommand(configUseCmd)

	rootCmd.AddCommand(configCmd)
}
