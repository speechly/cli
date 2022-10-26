package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
)

var transcribeCmd = &cobra.Command{
	Use: "transcribe <input_file>",
	Example: `speechly transcribe <input_file> --model /path/to/model/bundle
speechly transcribe <input_file> --app <app_id>`,
	Short: "Transcribe the given file(s) using on-device or cloud transcription",
	Args:  cobra.RangeArgs(1, 1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		model, err := cmd.Flags().GetString("model")
		if err != nil {
			log.Fatalf("Missing model bundle: %s", err)
		}

		inputPath := args[0]

		if model != "" {
			results, err := transcribeOnDevice(model, inputPath)
			printResults(results, inputPath, err == nil)
			if err != nil {
				log.Fatalf("Transcribing failed: %v", err)
			}
			return
		}

		appId, err := cmd.Flags().GetString("app")
		if err != nil {
			log.Fatalf("Missing app ID: %s", err)
		}

		if appId != "" {
			results, err := transcribeWithBatchAPI(ctx, appId, inputPath, false)
			printResults(results, inputPath, err == nil)
			if err != nil {
				log.Fatalf("Transcribing failed: %v", err)
			}
			return
		}
	},
}

func printResults(results []AudioCorpusItem, inputPath string, reportErrors bool) {
	for _, aci := range results {
		if strings.HasSuffix(inputPath, "wav") {
			fmt.Println(aci.Hypothesis)
		} else {
			b, err := json.Marshal(aci)
			if err != nil && reportErrors {
				log.Fatalf("Error in result generation: %v", err)
			}
			fmt.Println(string(b))
		}
	}
}

func init() {
	transcribeCmd.Flags().StringP("app", "a", "", "Application ID to use for cloud transcription")
	transcribeCmd.Flags().StringP("model", "m", "", "Model bundle file. This feature is available on Enterprise plans (https://speechly.com/pricing)")
	RootCmd.AddCommand(transcribeCmd)
}

type AudioCorpusItem struct {
	Audio      string `json:"audio"`
	Hypothesis string `json:"hypothesis,omitempty"`
	Transcript string `json:"transcript,omitempty"`
}

type AudioCorpusItemBatch struct {
	Audio   string `json:"audio"`
	BatchID string `json:"batch_id"`
}
