package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var transcribeCmd = &cobra.Command{
	Use:     "transcribe <input_file>",
	Example: `speechly transcribe <input_file>`,
	Short:   "Transcribe the given jsonlines file",
	Args:    cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		model, err := cmd.Flags().GetString("model")
		if err != nil {
			log.Fatalf("Error reading flags: %s", err)
		}
		inputPath := args[0]
		if model != "" {
			err = transcribeOnDevice(model, inputPath)
			if err != nil {
				log.Fatalf("Error in On-device Transcription: %s", err)
			}
			return
		}
		log.Fatal("This version of the Speechly CLI tool does not support Cloud Transcription.")
	},
}

func init() {
	RootCmd.AddCommand(transcribeCmd)
	transcribeCmd.Flags().StringP("model", "m", "", "On-device model bundle file")
}

type AudioCorpusItem struct {
	Audio      string `json:"audio"`
	Hypothesis string `json:"hypothesis"`
}
