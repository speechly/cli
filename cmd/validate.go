package cmd

import (
	"context"
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	salv1 "github.com/speechly/api/go/speechly/sal/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/speechly/cli/pkg/upload"
)

type ValidateWriter struct {
	appId  string
	stream salv1.Compiler_ValidateClient
}

func (u ValidateWriter) Write(data []byte) (n int, err error) {
	contentType := salv1.AppSource_CONTENT_TYPE_TAR
	as := &salv1.AppSource{AppId: u.appId, DataChunk: data, ContentType: contentType}
	if err = u.stream.Send(as); err != nil {
		return 0, err
	}
	return len(data), nil
}

var validateCmd = &cobra.Command{
	Use: "validate [<app_id>] <directory>",
	Example: `speechly validate -a UUID_APP_ID .
speechly validate UUID_APP_ID /usr/local/project/app`,
	Short: "Validate the given configuration for syntax errors",
	Long: `The contents of the directory given as argument is sent to the
API and validated. Possible errors are printed to stdout.`,
	Args: cobra.RangeArgs(1, 2),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		appId, _ := cmd.Flags().GetString("app")
		if appId == "" {
			if len(args) < 2 {
				return fmt.Errorf("app_id must be given with flag --app or as the first positional argument of two")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		inDir := args[0]
		if appId == "" {
			appId = args[0]
			inDir = args[1]
		}
		absPath, _ := filepath.Abs(inDir)
		log.Printf("Project dir: %s\n", absPath)
		// create a tar package from files in memory
		uploadData := upload.CreateTarFromDir(inDir)

		if len(uploadData.Files) == 0 {
			log.Fatalf("No files found for validation!\n\nPlease ensure the files are named *.yaml or *.csv")
		}

		messages, err := validateUploadData(ctx, appId, uploadData)
		if err != nil {
			log.Fatalf("Validate failed: %s", err)
		} else if len(messages) > 0 {
			printLineErrors(messages)
		}
		log.Println("Configuration OK.")
	},
}

func validateUploadData(ctx context.Context, appId string, ud upload.UploadData) ([]*salv1.LineReference, error) {
	compileClient, err := clients.CompileClient(ctx)
	if err != nil {
		return nil, err
	}

	// open a stream for upload
	stream, err := compileClient.Validate(ctx)
	if err != nil {
		return nil, err
	}

	// flush the tar from memory to the stream
	validateWriter := ValidateWriter{appId, stream}
	_, err = ud.Buf.WriteTo(validateWriter)
	if err != nil {
		return nil, err
	}

	validateResult, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}
	return validateResult.Messages, nil
}

func init() {
	rootCmd.AddCommand(validateCmd)
	validateCmd.Flags().StringP("app", "a", "", "application to validate the files for. Can alternatively be given as the first positional argument.")
}
