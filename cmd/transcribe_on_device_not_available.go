//go:build !on_device
// +build !on_device

package cmd

import "fmt"

func transcribeOnDevice(bundlePath string, corpusPath string) error {
	return fmt.Errorf("this version of the Speechly CLI tool does not support on-device transcription")
}
