package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/speechly/cli/pkg/clients"
)

var (
	RootCmd = &cobra.Command{
		Use:   "speechly",
		Short: "Speechly CLI",
		Long:  logoWithVersion(),
	}
)

var logo = `
███████╗██████╗ ███████╗███████╗ ██████╗██╗  ██╗██╗  ██╗   ██╗
██╔════╝██╔══██╗██╔════╝██╔════╝██╔════╝██║  ██║██║  ╚██╗ ██╔╝
███████╗██████╔╝█████╗  █████╗  ██║     ███████║██║   ╚████╔╝
╚════██║██╔═══╝ ██╔══╝  ██╔══╝  ██║     ██╔══██║██║    ╚██╔╝
███████║██║     ███████╗███████╗╚██████╗██║  ██║███████╗██║
╚══════╝╚═╝     ╚══════╝╚══════╝ ╚═════╝╚═╝  ╚═╝╚══════╝╚═╝
`

func logoWithVersion() string {
	if version == "development" || version == "latest" {
		return logo
	}
	return logo + "\n" + version
}

func failWithError(err error) {
	log.Fatalf("General failure; check your project settings with `speechly project`\n\nError: %v", err)
}

func Execute() error {
	ctx := clients.NewContext(failWithError)
	return RootCmd.ExecuteContext(ctx)
}
