package cmd

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/speechly/cli/pkg/upload"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Get the active configuration of the given app.",
	Long: `Fetches the currently stored configuration from the API. This command
does not check for validity of the stored configuration, but downloads the latest
version.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		outDir, _ := cmd.Flags().GetString("out")
		outDir, _ = filepath.Abs(outDir)

		log.Printf("Download current configuration to %s\n", outDir)

		client, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}

		if err := os.RemoveAll(outDir); err != nil {
			log.Fatalf("Could not clear the download directory %s: %s", outDir, err)
		}
		if err := os.MkdirAll(outDir, 0755); err != nil {
			log.Fatalf("Could not create the download directory %s: %s", outDir, err)
		}

		appId, _ := cmd.Flags().GetString("app")

		buf := []byte{}
		stream, err := client.DownloadCurrentTrainingData(ctx, &configv1.DownloadCurrentTrainingDataRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to get training data for %s: %s", appId, err)
		}
		ct := configv1.DownloadCurrentTrainingDataResponse_CONTENT_TYPE_UNSPECIFIED
		for {
			pkg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if ct == configv1.DownloadCurrentTrainingDataResponse_CONTENT_TYPE_UNSPECIFIED {
				ct = pkg.GetContentType()
			}
			if err != nil {
				log.Fatalf("Training data fetch failed: %s", err)
			}
			buf = append(buf, pkg.DataChunk...)
		}

		if ct == configv1.DownloadCurrentTrainingDataResponse_CONTENT_TYPE_TAR {
			if err := upload.ExtractTarToDir(outDir, bytes.NewReader(buf)); err != nil {
				log.Fatalf("Could not extract the configuration: %s", err)
			}
		} else {
			out := filepath.Join(outDir, "config.yaml")
			log.Printf("Writing file %s (%d bytes)\n", out, len(buf))
			if err := ioutil.WriteFile(out, buf, 0755); err != nil {
				log.Fatalf("Could not write configuration to %s: %s", out, err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("app", "a", "", "Which application's configuration to download")
	if err := downloadCmd.MarkFlagRequired("app"); err != nil {
		log.Fatalf("failed to init flags: %s", err)
	}
	downloadCmd.Flags().StringP("out", "o", "", "directory to write the training data in. All existing contents will be deleted.")
	if err := downloadCmd.MarkFlagRequired("out"); err != nil {
		log.Fatalf("failed to init flags: %s", err)
	}

}