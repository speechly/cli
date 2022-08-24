package upload

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

type UploadData struct {
	Files []string
	Buf   bytes.Buffer
}

func CreateTarFromDir(inDir string) UploadData {
	files, err := os.ReadDir(inDir)
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
			info, err := f.Info()
			if err != nil {
				log.Fatalf("Failed to read file: %s", err)
			}
			log.Printf("Adding %s (%d bytes)\n", f.Name(), info.Size())
			hdr := &tar.Header{
				Name: info.Name(),
				Mode: 0600,
				Size: info.Size(),
			}
			if err := tw.WriteHeader(hdr); err != nil {
				log.Fatalf("Failed to create a tar header: %s", err)
			}
			uploadFile := filepath.Join(inDir, f.Name())
			contents, err := os.ReadFile(uploadFile)
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

func ExtractTarToDir(outDir string, r io.Reader) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return err
		case header == nil:
			continue // skip empty files in tar
		}
		target := filepath.Join(outDir, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}
		case tar.TypeReg:
			fmt.Printf("Writing file %s (%d bytes)\n", target, header.Size)
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			_ = f.Close()
		}
	}
}
