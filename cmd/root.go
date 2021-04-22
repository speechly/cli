package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	analyticsv1 "github.com/speechly/api/go/speechly/analytics/v1"
	configv1 "github.com/speechly/api/go/speechly/config/v1"
	salv1 "github.com/speechly/api/go/speechly/sal/v1"
)

type Config struct {
	CurrentContext string            `mapstructure:"current-context"`
	Contexts       []SpeechlyContext `mapstructure:"contexts"`
}

type SpeechlyContext struct {
	Name   string `mapstructure:"name"`
	Host   string `mapstructure:"host"`
	Apikey string `mapstructure:"apikey"`
}

var (
	config_client    configv1.ConfigAPIClient
	compile_client   salv1.CompilerClient
	analytics_client analyticsv1.AnalyticsAPIClient

	conf    Config
	sc      SpeechlyContext
	rootCmd = &cobra.Command{
		Use:   "speechly",
		Short: "Speechly API Client",
		Long:  logo,
	}
)

var logo = `
███████╗██████╗ ███████╗███████╗ ██████╗██╗  ██╗██╗  ██╗   ██╗
██╔════╝██╔══██╗██╔════╝██╔════╝██╔════╝██║  ██║██║  ╚██╗ ██╔╝
███████╗██████╔╝█████╗  █████╗  ██║     ███████║██║   ╚████╔╝
╚════██║██╔═══╝ ██╔══╝  ██╔══╝  ██║     ██╔══██║██║    ╚██╔╝
███████║██║     ███████╗███████╗╚██████╗██║  ██║███████╗██║
╚══════╝╚═╝     ╚══════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝

                      API Client
`

func init() {
	initConfig()
}

func initConfig() {
	log.SetFlags(0)

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
		if len(os.Args) < 2 || os.Args[1] != "config" {
			log.Print("Please create a configuration file first:\n\n")
			log.Printf("%s config add --apikey APIKEY --name NAME", os.Args[0])
			log.Println("")
			os.Exit(1)
		}
		// viper has a problem with non-existent config files, just touch the default:
		file, err := os.Create(filepath.Join(home, ".speechly.yaml"))
		if err != nil {
			log.Fatalf("Could not initialize speechly config file: %s", err)
		}
		if err := file.Close(); err != nil {
			log.Fatalf("Could not write speechly config file: %v", err)
		}
	} else {
		if err := viper.Unmarshal(&conf); err != nil {
			log.Fatalf("Failed to unmarshal config file %s: %s", viper.ConfigFileUsed(), err)
		}
	}
	for _, item := range conf.Contexts {
		if item.Name == conf.CurrentContext {
			sc = item
			break
		}
	}
}

func Execute() error {
	if sc == (SpeechlyContext{}) {
		return rootCmd.ExecuteContext(context.TODO())
	}

	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", sc.Apikey))
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	serverAddr := sc.Host
	opts := []grpc.DialOption{
		grpc.WithBlock(),
	}
	if strings.Contains(sc.Host, "speechly.com") {
		// Always use TLS for Speechly hosts
		serverAddr = serverAddr + ":443"
		creds := credentials.NewTLS(&tls.Config{
			ServerName: sc.Host,
		})
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	conn, err := grpc.DialContext(ctx, serverAddr, opts...)
	if err != nil {
		log.Fatalf("Connecting to host %s failed: %s", sc.Host, err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Fatalf("Could not close connection: %v", err)
		}
	}()

	config_client = configv1.NewConfigAPIClient(conn)
	compile_client = salv1.NewCompilerClient(conn)
	analytics_client = analyticsv1.NewAnalyticsAPIClient(conn)

	return rootCmd.ExecuteContext(ctx)
}
