package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
)

var uploadCmd = &cobra.Command{
	Use:  "upload",
	Short: "Send a training data and configuration yaml to training",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId := args[0]
		inFile, _ := cmd.Flags().GetString("file")
		if inFile == "" {
			log.Fatalf("No file given for uploading")
		}
		reader, err := os.Open(inFile)
		if err != nil {
			log.Fatalf("Could not open file: %s: %s", inFile, err)
		}
		defer reader.Close()

		stream, err := client.UploadTrainingData(ctx)

		buffer := make([]byte, 32768)
		total := 0
		for {
			n, err := reader.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Could not read file %s: %s", inFile, err)
			}
			total += n
			req := &configv1.UploadTrainingDataRequest{AppId: appId, DataChunk: buffer[:n]}
			if err = stream.Send(req); err != nil {
				log.Fatalf("Uploading training data failed: %s", err)
			}
		}
		// Response from upload is empty, ignore:
		_, err = stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Upload failed: %s", err)
		}
		cmd.Printf("File %s (%d bytes) uploaded\n", inFile, total)
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringP("file", "f", "", "File to upload. Will start training.")
}
