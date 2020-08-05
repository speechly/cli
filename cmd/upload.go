package cmd

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"

	configv1 "github.com/speechly/cli/gen/go/speechly/config/v1"
)

type UploadWriter struct {
	appId  string
	stream configv1.ConfigAPI_UploadTrainingDataClient
}

func (u UploadWriter) Write(data []byte) (n int, err error) {
	req := &configv1.UploadTrainingDataRequest{AppId: u.appId, DataChunk: data, ContentType: configv1.UploadTrainingDataRequest_CONTENT_TYPE_TAR}
	if err = u.stream.Send(req); err != nil {
		return 0, err
	}
	return len(data), nil
}

var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "Send a training data and configuration yaml to training",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		var inDir string
		if len(args) > 0 {
			inDir = args[0]
		} else {
			inDir = "./"
		}
		absPath, _ := filepath.Abs(inDir)
		log.Printf("Project dir: %s\n", absPath)
		// create a tar package from files in memory
		buf := createTarFromDir(inDir)
		if buf.Len() == 0 {
			log.Fatalf("Nothing to upload.")
		}

		// open a stream for upload
		stream, err := client.UploadTrainingData(ctx)
		if err != nil {
			log.Fatalf("Failed to open upload stream: %s", err)
		}

		// flush the tar from memory to the stream
		uploadWriter := UploadWriter{appId, stream}
		n, err := buf.WriteTo(uploadWriter)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}

		// Response from upload is empty, ignore:
		_, err = stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Upload failed: %s", err)
		}

		cmd.Printf("%d bytes uploaded", n)
	},
}

func createTarFromDir(inDir string) bytes.Buffer {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Fatalf("Could not read files from %s", inDir)
	}
	// only accept yaml and csv files in the tar package
	configFileMatch := regexp.MustCompile(`.*?(csv|yaml)$`)
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
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

			contents, err := ioutil.ReadFile(filepath.Join(inDir, f.Name()))
			if err != nil {
				log.Fatalf("Failed to read file: %s", err)
			}
			if _, err := tw.Write(contents); err != nil {
				log.Fatalf("Failed to tar file: %s", err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		log.Fatalf("Package finalization failed: %s", err)
	}
	return buf
}

func init() {
	rootCmd.AddCommand(uploadCmd)
	uploadCmd.Flags().StringP("app", "a", "", "application id to upload the files to.")
	uploadCmd.MarkFlagRequired("app")
}
