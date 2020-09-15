package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete an existing application",
	Run: func(cmd *cobra.Command, args []string) {
		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			log.Fatalf("Missing force flag: %s", err)
		}

		dry, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatalf("Missing dry-run flag: %s", err)
		}

		id, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatalf("Missing app ID: %s", err)
		}

		if !force && !confirm(fmt.Sprintf("Deleting app %s, are you sure?", id), cmd.OutOrStdout(), cmd.InOrStdin()) {
			cmd.Println("Deletion aborted.")
			return
		}

		if !dry {
			if _, err := config_client.DeleteApp(
				cmd.Context(),
				&configv1.DeleteAppRequest{
					AppId: id,
				},
			); err != nil {
				log.Fatalf("Error deleting the app: %s", err)
			}
		}

		cmd.Printf("Successfully deleted app %s.\n", id)
	},
}

func init() {
	deleteCmd.Flags().StringP("app", "a", "", "application ID to delete")
	if err := deleteCmd.MarkFlagRequired("app"); err != nil {
		log.Fatalf("Internal error: %s", err)
	}

	deleteCmd.Flags().BoolP("force", "f", false, "skip confirmation prompt")
	deleteCmd.Flags().BoolP("dry-run", "d", false, "don't perform the deletion")

	rootCmd.AddCommand(deleteCmd)
}

func confirm(prompt string, dst io.Writer, src io.Reader) bool {
	read := bufio.NewReader(src)

	for {
		fmt.Fprintf(dst, "%s [y/n]: ", prompt)

		r, err := read.ReadString('\n')
		if err != nil {
			return false
		}

		r = strings.ToLower(strings.TrimSpace(r))
		if r == "y" || r == "yes" {
			return true
		} else if r == "n" || r == "no" {
			return false
		}
	}
}
