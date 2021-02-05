package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use: "compile [directory]",
	Example: `speechly compile -a UUID_APP_ID .
speechly compile -a UUID_APP_ID /usr/local/project/app`,
	Short: "Compiles a sample of examples from the given configuration",
	Long: `The contents of the directory given as argument is sent to the
API and compiled. If suffcessful, a sample of examples are printed to stdout.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		uploadData := createAndValidateTar(args[0])

		// open a stream for upload
		stream, err := compile_client.Compile(ctx)
		if err != nil {
			log.Fatalf("Failed to open validate stream: %s", err)
		}

		// flush the tar from memory to the stream
		compileWriter := CompileWriter{appId, stream}
		_, err = uploadData.buf.WriteTo(compileWriter)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}

		compileResult, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Validate failed: %s", err)
		}
		
		if len(compileResult.Messages) > 0 {
			printLineErrors(compileResult.Messages)
		} else {
			for _, message := range compileResult.Templates {
				log.Printf("%s", message)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringP("app", "a", "", "application to deploy the files to.")
	compileCmd.MarkFlagRequired("app")
}
