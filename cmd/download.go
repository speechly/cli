package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Get the active configuration.",
	Long: `Fetches the currently stored configuration from the API. This command
does not check for validity of the stored configuration, but downloads the latest
version.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		var out io.Writer
		outFile, _ := cmd.Flags().GetString("out")
		if outFile != "" {
			writer, err := os.Create(outFile)
			if err != nil {
				log.Fatalf("Could not open file for writing: %s", err)
			}
			defer writer.Close()
			out = writer
		} else {
			out = cmd.OutOrStdout()
		}

		appId, _ := cmd.Flags().GetString("app")
		stream, err := client.DownloadCurrentTrainingData(ctx, &configv1.DownloadCurrentTrainingDataRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to get training data for %s: %s", appId, err)
		}
		for {
			pkg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Training data fetch failed: %s", err)
			}
			if _, err = out.Write(pkg.DataChunk); err != nil {
				log.Fatalf("Writing data to output failed: %s", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("app", "a", "", "Which application's configuration to download")
	downloadCmd.MarkFlagRequired("app")
	downloadCmd.Flags().StringP("out", "o", "", "File to write the training data in. Will be overwritten.")

}
