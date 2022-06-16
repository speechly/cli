package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var transcribeCmd = &cobra.Command{
	Use:     "transcribe <app_id> <input_file>",
	Example: `speechly transcribe <app_id> <input_file>`,
	Short:   "Transcribe the given jsonlines file",
	Args:    cobra.RangeArgs(2, 2),
	Run: func(cmd *cobra.Command, args []string) {
		model, err := cmd.Flags().GetString("model")
		if err != nil {
			log.Fatalf("Error reading flags: %s", err)
		}
		appID := args[0]
		inputPath := args[1]
		if model != "" {
			err = transcribeOnDevice(strings.Split(model, ","), appID, inputPath)
			if err != nil {
				log.Fatalf("Error in On Device Transcription: %s", err)
			}
			return
		}
		log.Fatal("This version of the Speechly CLI tool does not support Cloud Transcription.")
	},
}

func init() {
	rootCmd.AddCommand(transcribeCmd)
	transcribeCmd.Flags().StringP("model", "m", "", "on device model file paths as a comma separated list")
}

type AudioCorpusItem struct {
	Audio      string `json:"audio"`
	Hypothesis string `json:"hypothesis"`
}
