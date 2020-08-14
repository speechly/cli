package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	compilev1 "github.com/speechly/cli/gen/go/speechly/sal/v1"
)

type CompileWriter struct {
	appId  string
	stream compilev1.Compiler_ValidateClient
}

func (u CompileWriter) Write(data []byte) (n int, err error) {
	req := &compilev1.AppSource{AppId: u.appId, DataChunk: data}
	if err = u.stream.Send(req); err != nil {
		return 0, err
	}
	return len(data), nil
}

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
		inDir := args[0]
		absPath, _ := filepath.Abs(inDir)
		log.Printf("Project dir: %s\n", absPath)
		// create a tar package from files in memory
		uploadData := createTarFromDir(inDir)

		if len(uploadData.files) == 0 {
			log.Fatalf("No files found for validation!\n\nPlease ensure the files are named *.yaml or *.csv")
		}

		// open a stream for upload
		stream, err := compile_client.Validate(ctx)
		if err != nil {
			log.Fatalf("Failed to open validate stream: %s", err)
		}

		// flush the tar from memory to the stream
		validateWriter := CompileWriter{appId, stream}
		_, err = uploadData.buf.WriteTo(validateWriter)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}

		validateResult, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Validate failed: %s", err)
		}
		if len(validateResult.Messages) > 0 {
			log.Println("Configuration validation failed")
			for _, message := range validateResult.Messages {
				var errorLevel string
				switch message.Level {
				case compilev1.LineReference_LEVEL_NOTE:
					errorLevel = "NOTE"
				case compilev1.LineReference_LEVEL_WARNING:
					errorLevel = "WARNING"
				case compilev1.LineReference_LEVEL_ERROR:
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
			log.Println("Configuration OK.")
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("app", "a", "", "application to deploy the files to.")
	validateCmd.MarkFlagRequired("app")
}
