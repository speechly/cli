package cmd

import (
	"bytes"
	"log"
	"os"

	"github.com/spf13/cobra"

	salv1 "github.com/speechly/api/go/speechly/sal/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/speechly/cli/pkg/upload"
)

type ConvertWriter struct {
	stream   salv1.Compiler_ConvertClient
	format   salv1.ConvertRequest_InputFormat
	language string
}

func (u ConvertWriter) Write(data []byte) (n int, err error) {
	req := &salv1.ConvertRequest{InputFormat: u.format, Language: u.language, DataChunk: data}
	if err = u.stream.Send(req); err != nil {
		return 0, err
	}
	return len(data), nil
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Converts an Alexa Interaction Model in JSON format to a Speechly configuration",
	Example: `speechly convert my-alexa-skill.json
speechly convert --language en-US my-alexa-skill.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		language, _ := cmd.Flags().GetString("language")
		if len(language) == 0 {
			language = "en-US"
		}

		data, err := os.ReadFile(args[0])
		if err != nil {
			log.Fatalf("Error reading input from file %s", args[0])
		}

		client, err := clients.CompileClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}
		stream, err := client.Convert(ctx)
		if err != nil {
			log.Fatalf("Failed to open convert stream: %s", err)
		}

		convertWriter := ConvertWriter{stream, salv1.ConvertRequest_FORMAT_ALEXA, language}
		_, err = convertWriter.Write(data)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}
		log.Printf("Converting to Speechly configuration...")

		convertResult, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Conversion failed: %s", err)
		}

		if convertResult.Status == salv1.ConvertResult_CONVERT_FAILED {
			log.Fatalf("Conversion failed, message: %s\nAre you sure the input is an Alexa Interaction Model in JSON format?",
				convertResult.Warnings)
		}

		if convertResult.Status == salv1.ConvertResult_CONVERT_WARNINGS {
			log.Printf("Conversion done with warnings:\n%s", convertResult.Warnings)
		} else {
			// conversion must have been success!
			log.Printf("Conversion done!")
		}

		if err := upload.ExtractTarToDir(".", bytes.NewReader(convertResult.Result.DataChunk)); err != nil {
			log.Fatalf("Error when extracting configuration: %s", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(convertCmd)
	convertCmd.Flags().StringP("language", "l", "en-US", "Language of input")
}
