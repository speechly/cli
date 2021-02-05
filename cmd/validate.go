package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use: "validate [directory]",
	Example: `speechly validate -a UUID_APP_ID .
speechly validate -a UUID_APP_ID /usr/local/project/app`,
	Short: "Validate the given configuration for syntax errors",
	Long: `The contents of the directory given as argument is sent to the
API and validated. Possible errors are printed to stdout.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		uploadData := createAndValidateTar(args[0])

		// open a stream for upload
		stream, err := compile_client.Validate(ctx)
		if err != nil {
			log.Fatalf("Failed to open validate stream: %s", err)
		}

		// flush the tar from memory to the stream
		validateWriter := ValidateWriter{appId, stream}
		_, err = uploadData.buf.WriteTo(validateWriter)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}

		validateResult, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Validate failed: %s", err)
		}
		if len(validateResult.Messages) > 0 {
			printLineErrors(validateResult.Messages)
		} else {
			log.Println("Configuration OK.")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("app", "a", "", "application to deploy the files to.")
	validateCmd.MarkFlagRequired("app")
}
