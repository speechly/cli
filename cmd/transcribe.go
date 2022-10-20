package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var transcribeCmd = &cobra.Command{
	Use: "transcribe <input_file>",
	Example: `speechly transcribe <input_file> --model /path/to/model/bundle
speechly transcribe <input_file> --app <app_id>`,
	Short: "Transcribe the given file(s) using on-device or cloud transcription",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		model, err := cmd.Flags().GetString("model")
		if err != nil {
			log.Fatalf("Missing model bundle: %s", err)
		}

		inputPath := args[0]

		if model != "" {
			err = transcribeOnDevice(model, inputPath)
			if err != nil {
				log.Fatalf("Error in On-device Transcription: %s", err)
			}
			return
		}

		appId, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatalf("Missing app ID: %s", err)
		}

		if appId != "" {
			err = transcribeOnDevice(appId, inputPath)
			if err != nil {
				log.Fatalf("Error in cloud transcription: %s", err)
			}
			return
		}
	},
}

func init() {
	transcribeCmd.Flags().StringP("app", "a", "", "Application ID to use for cloud transcription")
	transcribeCmd.Flags().StringP("model", "m", "", "Model bundle file. This feature is available on Enterprise plans (https://speechly.com/pricing)")
	RootCmd.AddCommand(transcribeCmd)
}

type AudioCorpusItem struct {
	Audio      string `json:"audio"`
	Hypothesis string `json:"hypothesis"`
}
