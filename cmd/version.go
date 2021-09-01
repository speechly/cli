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
		cmd.Println(fmt.Sprintf("version: %s", version))
		if commit != "" {
			cmd.Println(fmt.Sprintf("commit: %s", commit))
		}
		if date != "" {
			cmd.Println(fmt.Sprintf("date: %s", date))
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
