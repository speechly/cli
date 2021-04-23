package cmd

import (
	"fmt"
	"log"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Speechly API access configurations",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("Config file used: %s\n", viper.ConfigFileUsed())
		cmd.Printf("Current config: %s\n", conf.CurrentContext)
		cmd.Printf("Contexts:\n")
		for _, c := range conf.Contexts {
			cmd.Printf("- %s\n", c.Name)
		}

	},
}

var validName = regexp.MustCompile(`[^A-Za-z0-9]+`)

var configAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new configuration context",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if validName.MatchString(name) {
			return fmt.Errorf("Invalid name: %s", name)
		}
		for _, c := range conf.Contexts {
			if c.Name == name {
				return fmt.Errorf("Context named %s already exists", name)
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		host, _ := cmd.Flags().GetString("host")
		apikey, _ := cmd.Flags().GetString("apikey")
		name, _ := cmd.Flags().GetString("name")
		viper.Set("contexts", append(conf.Contexts, SpeechlyContext{Host: host, Apikey: apikey, Name: name}))
		viper.Set("current-context", name)
		if err := viper.WriteConfig(); err != nil {
			log.Fatalf("Failed to write the config: %s", err)
		}
		cmd.Printf("Wrote configuration to file: %s\n", viper.ConfigFileUsed())
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove a context from configuration",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		if err := ensureContextExists(name); err != nil {
			return err
		}
		if name == viper.Get("current-context") {
			return fmt.Errorf("Cannot remove active context: %s", name)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		cmd.Printf("Removing context: %s\n", name)
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
	Short: "Select the default context used",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		return ensureContextExists(name)
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

func ensureContextExists(name string) error {
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
	configAddCmd.Flags().String("name", "", "A short unique name for the context")
	if err := configAddCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configAddCmd.Flags().String("host", "api.speechly.com", "API address")
	configCmd.AddCommand(configAddCmd)

	configRemoveCmd.Flags().String("name", "", "The short name for context to be deleted")
	if err := configRemoveCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configCmd.AddCommand(configRemoveCmd)

	configUseCmd.Flags().String("name", "", "A short unique name for the context")
	if err := configUseCmd.MarkFlagRequired("name"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	configCmd.AddCommand(configUseCmd)

	rootCmd.AddCommand(configCmd)
}
