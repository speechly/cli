package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/speechly/cli/pkg/clients"
)

var (
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

func failWithError(err error) {
	log.Fatalf("General failure; check your configuration with `speechly config`\n\nError: %v", err)
}

func Execute() error {
	ctx := clients.NewContext(failWithError)
	return rootCmd.ExecuteContext(ctx)
}
