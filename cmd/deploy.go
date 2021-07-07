package cmd

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	"github.com/speechly/cli/pkg/clients"
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

type UploadData struct {
	files []string
	buf   bytes.Buffer
}

var deployCmd = &cobra.Command{
	Use: "deploy [directory]",
	Example: `speechly deploy . -a UUID_APP_ID
speechly deploy /usr/local/project/app -a UUID_APP_ID`,
	Short: "Send the contents of a local directory to training",
	Long: `The contents of the directory given as argument is sent to the
API and validated. Then, a new model is trained and automatically deployed
as the active model for the application.`,
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
			log.Fatalf("Nothing to deploy!\n\nPlease ensure the files are named *.yaml or *.csv")
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
		n, err := uploadData.buf.WriteTo(deployWriter)
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
			waitForDeploymentFinished(cmd, configClient, appId)
		}
	},
}

func createTarFromDir(inDir string) UploadData {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Fatalf("Could not read files from %s", inDir)
	}
	// only accept yaml and csv files in the tar package
	configFileMatch := regexp.MustCompile(`.*?(csv|yaml)$`)
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	uploadFiles := []string{}
	for _, f := range files {
		if configFileMatch.MatchString(f.Name()) {
			log.Printf("Adding %s (%d bytes)\n", f.Name(), f.Size())
			hdr := &tar.Header{
				Name: f.Name(),
				Mode: 0600,
				Size: f.Size(),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Fatalf("Failed to create a tar header: %s", err)
			}
			uploadFile := filepath.Join(inDir, f.Name())
			contents, err := ioutil.ReadFile(uploadFile)
			if err != nil {
				log.Fatalf("Failed to read file: %s", err)
			}
			if _, err := tw.Write(contents); err != nil {
				log.Fatalf("Failed to tar file: %s", err)
			}
			uploadFiles = append(uploadFiles, uploadFile)
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatalf("Package finalization failed: %s", err)
	}
	return UploadData{uploadFiles, buf}
}

func init() {
	rootCmd.AddCommand(deployCmd)
	deployCmd.Flags().StringP("app", "a", "", "application to deploy the files to.")
	if err := deployCmd.MarkFlagRequired("app"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
	deployCmd.Flags().BoolP("watch", "w", false, "wait for training to be finished")
}
