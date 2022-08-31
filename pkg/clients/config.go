package clients

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

type Config struct {
	CurrentContext string            `mapstructure:"current-context"`
	Contexts       []SpeechlyContext `mapstructure:"contexts"`
}

type SpeechlyContext struct {
	Name       string `mapstructure:"name"`
	Host       string `mapstructure:"host"`
	Apikey     string `mapstructure:"apikey"`
	RemoteName string `mapstructure:"remotename"`
}

func (conf *Config) GetSpeechlyContext() *SpeechlyContext {
	if conf.CurrentContext == "" {
		return nil
	}
	var sc *SpeechlyContext
	for _, item := range conf.Contexts {
		if item.Name == conf.CurrentContext {
			sc = &item
			break
		}
	}
	if sc == nil {
		return nil
	}
	return sc
}

func getSpeechlyConfig() (*Config, error) {
	log.SetFlags(0)

	apikey := os.Getenv("SPEECHLY_APIKEY")
	if apikey != "" {
		host := os.Getenv("SPEECHLY_HOST")
		if host == "" {
			host = "api.speechly.com"
		}
		return &Config{
			CurrentContext: "default",
			Contexts: []SpeechlyContext{{
				Name:   "default",
				Host:   host,
				Apikey: apikey,
			}},
		}, nil
	}
	var conf Config

	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("Could not find $HOME: %s", err)
	}
	viper.AddConfigPath(home)
	viper.AddConfigPath(".")
	viper.SetConfigName(".speechly")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if len(os.Args) < 2 || (os.Args[1] != "config" && os.Args[1] != "projects") {
			log.Print("Please add a project first:\n\n")
			log.Print(" 1. Create a new API token in Speechly Dashboard (https://api.speechly.com/dashboard)\n")
			log.Printf(" 2. Copy the token and run: %s projects add --apikey <api_token>\n\n", os.Args[0])
			log.Print("Learn more: https://docs.speechly.com/dev-tools/command-line-tool")
			os.Exit(1)
		}
		// viper has a problem with non-existent config files, just touch the default:
		file, err := os.Create(filepath.Join(home, ".speechly.yaml"))
		if err != nil {
			return nil, fmt.Errorf("could not initialize Speechly settings file: %s", err)
		}
		if err := file.Close(); err != nil {
			return nil, fmt.Errorf("could not write Speechly settings file: %v", err)
		}
	} else {
		if err := viper.Unmarshal(&conf); err != nil {
			return nil, fmt.Errorf("failed to unmarshal Speechly settings file %s: %s", viper.ConfigFileUsed(), err)
		}
	}
	return &conf, nil
}
