package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	salv1 "github.com/speechly/api/go/speechly/sal/v1"
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
		inDir := args[0]
		absPath, _ := filepath.Abs(inDir)
		log.Printf("Project dir: %s\n", absPath)
		// create a tar package from files in memory
		uploadData := createTarFromDir(inDir)

		if len(uploadData.files) == 0 {
			log.Fatalf("No files found to compile!\n\nPlease ensure the files are named *.yaml or *.csv")
		}

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
			log.Println("Configuration validation failed")
			for _, message := range compileResult.Messages {
				var errorLevel string
				switch message.Level {
				case salv1.LineReference_LEVEL_NOTE:
					errorLevel = "NOTE"
				case salv1.LineReference_LEVEL_WARNING:
					errorLevel = "WARNING"
				case salv1.LineReference_LEVEL_ERROR:
					errorLevel = "ERROR"
				}
				if message.File != "" {
					log.Printf("%s:%d:%d:%s:%s\n", message.File, message.Line,
						message.Column, errorLevel, message.Message)
				} else {
					log.Printf("%s: %s", errorLevel, message.Message)
				}
			}
			os.Exit(1)
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
