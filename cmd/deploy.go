package cmd

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/speechly/cli/pkg/upload"
)

type DeployWriter struct {
	appId  string
	stream configv1.ConfigAPI_UploadTrainingDataClient
}

func (u DeployWriter) Write(data []byte) (n int, err error) {
	req := &configv1.UploadTrainingDataRequest{AppId: u.appId, DataChunk: data, ContentType: configv1.UploadTrainingDataRequest_CONTENT_TYPE_TAR}
	if err = u.stream.Send(req); err != nil {
		return 0, err
	}
	return len(data), nil
}

var deployCmd = &cobra.Command{
	Use: "deploy [<app_id>] <directory>",
	Example: `speechly deploy . -a <app_id>
speechly deploy <app_id> /usr/local/project/app`,
	Short: "Send the contents of a local directory to training",
	Long: `The contents of the directory given as argument is sent to the
API and validated. Then, a new model is trained and automatically deployed
as the active model for the application.`,
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
		inputDirectory := args[0]
		if appId == "" {
			appId = args[0]
			inputDirectory = args[1]
		}
		absPath, _ := filepath.Abs(inputDirectory)
		log.Printf("Project dir: %s\n", absPath)
		// create a tar package from files in memory
		uploadData := upload.CreateTarFromDir(inputDirectory)

		if len(uploadData.Files) == 0 {
			log.Fatalf("Nothing to deploy!\n\nPlease ensure the files are named *.yaml or *.csv")
		}

		messages, err := validateUploadData(ctx, appId, uploadData)
		if err != nil {
			log.Fatalf("Validate failed: %s", err)
		} else if len(messages) > 0 {
			printLineErrors(messages)
		}

		configClient, err := clients.ConfigClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}

		// open a stream for upload
		stream, err := configClient.UploadTrainingData(ctx)
		if err != nil {
			log.Fatalf("Failed to open deploy stream: %s", err)
		}

		// flush the tar from memory to the stream
		deployWriter := DeployWriter{appId, stream}
		n, err := uploadData.Buf.WriteTo(deployWriter)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}

		// Response from deploy is empty, ignore:
		_, err = stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Deploy failed: %s", err)
		}

		cmd.Printf("%d bytes uploaded, training and deployment proceeding.\n", n)

		// if watch flag given, wait for deployment to finish
		wait, _ := cmd.Flags().GetBool("watch")
		if wait {
			waitForAppStatus(cmd, configClient, appId, configv1.App_STATUS_TRAINED)
		}
	},
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringP("app", "a", "", "application to deploy the files to. Can alternatively be given as the first positional argument.")
	deployCmd.Flags().BoolP("watch", "w", false, "wait for training to be finished")
}
