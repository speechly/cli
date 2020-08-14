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
	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
	compilev1 "github.com/speechly/cli/gen/go/speechly/sal/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
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
	config_client  configv1.ConfigAPIClient
	compile_client compilev1.CompilerClient
	conf           Config
	sc             SpeechlyContext
	rootCmd        = &cobra.Command{
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
		if os.Args[1] != "config" {
			log.Println("Please create a configuration file first:")
			log.Println("")
			log.Println("\tspeechly config add --apikey APIKEY")
			log.Println("")
			os.Exit(1)
		}
		// viper has a problem with non-existent config files, just touch the default:
		file, err := os.Create(filepath.Join(home, ".speechly.yaml"))
		if err != nil {
			log.Fatalf("Could not initialize speechly config file: %s", err)
		}
		file.Close()
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
		return rootCmd.ExecuteContext(nil)
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
	defer conn.Close()

	config_client = configv1.NewConfigAPIClient(conn)
	compile_client = compilev1.NewCompilerClient(conn)

	return rootCmd.ExecuteContext(ctx)
}
