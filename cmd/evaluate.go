package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var evaluateCmd = &cobra.Command{
	Use:     "evaluate",
	Example: `speechly evaluate --input output.txt --ground-truth ground-truth.txt`,
	Short:   "Compute accuracy between annotated examples (given by 'speechly annotate') and ground truth.",
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		annotatedFn, err := cmd.Flags().GetString("input")
		if err != nil || len(annotatedFn) == 0 {
			log.Fatalf("Annotated file is invalid: %v", err)
		}

		groundTruthFn, err := cmd.Flags().GetString("ground-truth")
		if err != nil || len(groundTruthFn) == 0 {
			log.Fatalf("Ground-truth file is invalid: %v", err)
		}
		annotatedData := readLines(annotatedFn)
		groundTruthData := readLines(groundTruthFn)
		EvaluateAnnotatedUtterances(annotatedData, groundTruthData)
	},
}

func EvaluateAnnotatedUtterances(annotatedData []string, groundTruthData []string) {
	if len(annotatedData) != len(groundTruthData) {
		log.Fatalf(
			"Input files should have same length, but --annotated has %d lines and --ground-truth %d lines.",
			len(annotatedData),
			len(groundTruthData),
		)
	}

	n := float64(len(annotatedData))
	hits := 0.0
	for i, aUtt := range annotatedData {
		gtUtt := groundTruthData[i]
		if aUtt == gtUtt {
			hits += 1.0
			continue
		}
		fmt.Println("Ground truth had:")
		fmt.Println("  " + gtUtt)
		fmt.Println("but prediction was:")
		fmt.Println("  " + aUtt)
		fmt.Println()
	}
	fmt.Println("Matching rows out of total: ")
	fmt.Printf("%.0f / %.0f\n", hits, n)
	fmt.Println("Accuracy:")
	fmt.Printf("%.2f\n", hits/n)
}

func init() {
	rootCmd.AddCommand(evaluateCmd)
	evaluateCmd.Flags().StringP("input", "i", "", "SAL annotated utterances, as given by 'speechly annotate' command.")
	if err := evaluateCmd.MarkFlagRequired("input"); err != nil {
		log.Fatalf("Failed to init flags: %s", err)
	}
	evaluateCmd.Flags().StringP("ground-truth", "t", "", "manually verified ground-truths for annotated examples")
	if err := evaluateCmd.MarkFlagRequired("ground-truth"); err != nil {
		log.Fatalf("Failed to init flags: %s", err)
	}
}
