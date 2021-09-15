package upload

import (
	"archive/tar"
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
)

type UploadData struct {
	Files []string
	Buf   bytes.Buffer
}

func CreateTarFromDir(inDir string) UploadData {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		log.Fatalf("Could not read files from %s: %s", inDir, err)
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
