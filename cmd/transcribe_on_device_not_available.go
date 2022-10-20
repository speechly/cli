//go:build !on_device
// +build !on_device

package cmd

import (
	"fmt"
	"log"
)

func transcribeOnDevice(appId string, filePath string) error {
	if appId != "" && filePath != "" {
		log.Printf("Using configuration in application: %s\nUploading file: %s …\n\n", appId, filePath)
		log.Printf("Transcription in the linguistic sense is the systematic representation of spoken language in written form. The source can either be utterances (speech or sign language) or preexisting text in another writing system.\n\n")
	}
	return fmt.Errorf("something went wrong")
}
