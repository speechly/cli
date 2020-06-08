package cmd

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	//"time"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
	homedir "github.com/mitchellh/go-homedir"
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
	client  configv1.ConfigAPIClient
	conf    Config
	sc      SpeechlyContext
	rootCmd = &cobra.Command{
		Use:   "speechly",
		Short: "Speechly API Client",
	}
)

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
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Println("Please create a configuration file first:")
		log.Println("")
		log.Println("\tspeechly config add --apikey APIKEY")
		log.Println("")
		//log.Fatalf("Could not read config file: %s", err)
	}
	if err := viper.Unmarshal(&conf); err != nil {
		log.Fatalf("Failed to unmarshal config file %s: %s", viper.ConfigFileUsed(), err)
	}
	for _, item := range conf.Contexts {
		if item.Name == conf.CurrentContext {
			sc = item
			break
		}
	}
	if sc == (SpeechlyContext{}) {
		log.Fatalf("Could not resolve selected context: %s", conf.CurrentContext)
	}
}

func Execute() error {

	//cs := viper.UnmarshalKey("contexts", &sc)

	md := metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", sc.Apikey))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	//ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	//defer cancel()

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
	}

	conn, err := grpc.DialContext(ctx, serverAddr, opts...)
	if err != nil {
		log.Fatalf("Connecting to host %s failed: %s", sc.Host, err)
	}
	defer conn.Close()

	client = configv1.NewConfigAPIClient(conn)

	return rootCmd.ExecuteContext(ctx)
}
