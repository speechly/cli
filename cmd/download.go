package cmd

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/speechly/cli/pkg/upload"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download the active configuration or model bundle of the given app",
	Long:  "Fetches the currently stored configuration or model bundle. This command does not check for validity of the stored configuration, but downloads the latest version.",
	Example: `speechly download <app_id> /path/to/config
speechly download --app <app_id> . --model tflite`,
	Args: cobra.RangeArgs(1, 2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		appId, _ := cmd.Flags().GetString("app")
		if appId == "" {
			if len(args) < 2 {
				return fmt.Errorf("app_id must be given with flag --app or as the first positional argument of two")
			}
		}

		model, _ := cmd.Flags().GetString("model")
		if !map[string]bool{"ort": true, "coreml": true, "tflite": true, "": true, "all": true}[model] {
			return fmt.Errorf("\"%s\" is not a valid option. Available options are: ort, tflite, coreml and all", model)
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		model, _ := cmd.Flags().GetString("model")
		outputDirectory := args[0]
		if appId == "" {
			appId = args[0]
			outputDirectory = args[1]
		}
		absPath, _ := filepath.Abs(outputDirectory)
		log.Printf("Downloading to %s\n", absPath)

		d, err := os.Open(absPath)
		if err != nil {
			log.Fatalf("Reading output directory failed: %s", err)
		}

		defer func() {
			_ = d.Close()
		}()
		if model == "" {
			downloadCurrentConfiguration(ctx, d, absPath, appId)
		} else if model == "all" {
			ok := downloadCurrentModel(ctx, absPath, appId, "ort")
			ok = downloadCurrentModel(ctx, absPath, appId, "coreml") || ok
			ok = downloadCurrentModel(ctx, absPath, appId, "tflite") || ok
			if !ok {
				log.Fatalf("This feature is available on Enterprise plans (https://speechly.com/pricing)")
			}
		} else {
			if !downloadCurrentModel(ctx, absPath, appId, model) {
				log.Fatalf("This feature is available on Enterprise plans (https://speechly.com/pricing)")
			}
		}

	},
}

func downloadCurrentConfiguration(ctx context.Context, d *os.File, absPath string, appId string) {
	client, err := clients.ConfigClient(ctx)
	if err != nil {
		log.Fatalf("Error connecting to API: %s", err)
	}

	var buf []byte
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
		if err := upload.ExtractTarToDir(absPath, bytes.NewReader(buf)); err != nil {
			log.Fatalf("Could not extract the configuration: %s", err)
		}
	} else {
		out := filepath.Join(absPath, "config.yaml")
		log.Printf("Writing file %s (%d bytes)\n", out, len(buf))
		if err := os.WriteFile(out, buf, 0755); err != nil {
			log.Fatalf("Could not write configuration to %s: %s", out, err)
		}
	}
}

func downloadCurrentModel(ctx context.Context, absPath string, appId string, model string) bool {
	client, err := clients.ModelClient(ctx)
	if err != nil {
		log.Fatalf("Error connecting to API: %s", err)
	}

	var ma configv1.DownloadModelRequest_ModelArchitecture
	if model == "ort" {
		ma = configv1.DownloadModelRequest_MODEL_ARCHITECTURE_ORT
	} else if model == "coreml" {
		ma = configv1.DownloadModelRequest_MODEL_ARCHITECTURE_COREML
	} else if model == "tflite" {
		ma = configv1.DownloadModelRequest_MODEL_ARCHITECTURE_TFLITE
	} else {
		ma = configv1.DownloadModelRequest_MODEL_ARCHITECTURE_INVALID
	}

	stream, err := client.DownloadModel(ctx, &configv1.DownloadModelRequest{AppId: appId, ModelArchitecture: ma})
	if err != nil {
		log.Fatalf("Failed to get training data for %s: %s", appId, err)
	}

	var (
		buf []byte
		out string
	)

	out = filepath.Join(absPath, fmt.Sprintf("%s.%s.bundle", appId, model))
	_, err = os.Stat(out)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("Failed to verify out file %s: %s", absPath, err)
		}

		for {
			pkg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if status.Code(err) == codes.PermissionDenied {
				return false
			}
			if err != nil {
				log.Fatalf("model fetch failed out file: %s", err)
			}
			buf = append(buf, pkg.Chunk...)
		}

		log.Printf("Writing file %s (%d bytes)\n", out, len(buf))
		if err := os.WriteFile(out, buf, 0644); err != nil {
			log.Fatalf("Could not write configuration to %s: %s", out, err)
		}

	} else {
		log.Printf("File %s exists, skipping\n", out)
	}

	return true
}

func init() {
	RootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringP("app", "a", "", "Application which configuration or model bundle to download. Can be given as the first positional argument.")
	downloadCmd.Flags().StringP("model", "m", "", "Model bundle machine learning framework. Available options are: ort, tflite, coreml and all. This feature is available on Enterprise plans (https://speechly.com/pricing)")
}
