package cmd

import (
	"log"
	"os"
	"path/filepath"
	salv1 "github.com/speechly/api/go/speechly/sal/v1"
)

type ValidateWriter struct {
	appId  string
	stream salv1.Compiler_ValidateClient
}

type CompileWriter struct {
	appId  string
	stream salv1.Compiler_CompileClient
}

func createAppsource(appId string, data []byte) *salv1.AppSource {
	contentType := salv1.AppSource_CONTENT_TYPE_TAR
	return &salv1.AppSource{AppId: appId, DataChunk: data, ContentType: contentType}
}

func (u ValidateWriter) Write(data []byte) (n int, err error) {
	if err = u.stream.Send(createAppsource(u.appId, data)); err != nil {
		return 0, err
	}
	return len(data), nil
}

func (u CompileWriter) Write(data []byte) (n int, err error) {
	if err = u.stream.Send(createAppsource(u.appId, data)); err != nil {
		return 0, err
	}
	return len(data), nil
}

func printLineErrors(messages []*salv1.LineReference) {
	log.Println("Configuration validation failed")
	for _, message := range messages {
		var errorLevel string
		switch message.Level {
		case salv1.LineReference_LEVEL_NOTE:
			errorLevel = "NOTE"
		case salv1.LineReference_LEVEL_WARNING:
			errorLevel = "WARNING"
		case salv1.LineReference_LEVEL_ERROR:
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
}

func createAndValidateTar(inDir string) UploadData {
	absPath, _ := filepath.Abs(inDir)
	log.Printf("Project dir: %s\n", absPath)
	// create a tar package from files in memory
	uploadData := createTarFromDir(inDir)

	if len(uploadData.files) == 0 {
		log.Fatalf("No files found for validation!\n\nPlease ensure the files are named *.yaml or *.csv")
	}
	return uploadData
}