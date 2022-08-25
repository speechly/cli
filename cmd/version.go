package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "development"
	commit  = ""
	date    = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println(fmt.Sprintf("Version: %s", version))
		if commit != "" {
			cmd.Println(fmt.Sprintf("Commit: %s", commit))
		}
		if date != "" {
			cmd.Println(fmt.Sprintf("Date: %s", date))
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
