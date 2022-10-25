package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	configv1 "github.com/speechly/api/go/speechly/config/v1"
	salv1 "github.com/speechly/api/go/speechly/sal/v1"
	wluv1 "github.com/speechly/api/go/speechly/slu/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

func waitForAppStatus(cmd *cobra.Command, configClient configv1.ConfigAPIClient, appId string, status configv1.App_Status) {
	ctx := cmd.Context()

	for {
		app, err := configClient.GetApp(ctx, &configv1.GetAppRequest{AppId: appId})
		if err != nil {
			log.Fatalf("Failed to refresh app %s: %s", appId, err)
		}
		cmd.Printf("Status:\t%s", app.App.Status)
		switch app.App.Status {
		case configv1.App_STATUS_NEW:
			cmd.Printf(", queued (%d jobs before this)", app.App.QueueSize)
		case configv1.App_STATUS_TRAINING:
			age := time.Duration(app.App.TrainingTimeSec) * time.Second
			est := time.Duration(app.App.EstimatedTrainingTimeSec) * time.Second
			cmd.Printf(", age %s, previous deployment took %s", age, est)
		case configv1.App_STATUS_FAILED:
			cmd.Println()
			cmd.Printf("Error: %s", app.App.ErrorMsg)
		}
		cmd.Println()

		if app.App.Status >= status {
			break
		}
		time.Sleep(10 * time.Second)
	}
}

func checkSoleAppArgument(cmd *cobra.Command, args []string) error {
	appId, err := cmd.Flags().GetString("app")
	if err != nil {
		log.Fatalf("Missing app ID: %s", err)
	}
	if appId == "" && len(args) < 1 {
		return fmt.Errorf("app_id must be given with flag --app or as the sole positional argument")
	}
	return nil
}

func readReferenceDate(cmd *cobra.Command) (time.Time, error) {
	refD := time.Now()
	refDS, err := cmd.Flags().GetString("reference-date")
	if err != nil {
		return time.Time{}, nil
		// log.Fatalf("reference-date is invalid: %s", err)
	}

	if len(refDS) > 0 {
		refD, err = time.Parse("2006-01-02", refDS)
		if err != nil {
			return time.Time{}, nil
		}
	}
	return refD, nil
}

func runThroughWLU(ctx context.Context, appID string, inputFile string, deAnnotate bool, refD time.Time) (*wluv1.TextsResponse, []string, error) {
	wluClient, err := clients.WLUClient(ctx)
	if err != nil {
		return nil, nil, err
	}

	data := readLines(inputFile)

	annotated := data
	transcripts := make([]string, len(data))
	for i, line := range data {
		transcripts[i] = removeAnnotations(line)
	}
	data = transcripts

	if deAnnotate {
		for _, line := range data {
			fmt.Println(line)
		}
		os.Exit(0)
	}

	wluRequests := make([]*wluv1.WLURequest, len(data))
	for i, line := range data {
		wluRequests[i] = &wluv1.WLURequest{
			Text:          line,
			ReferenceTime: timestamppb.New(refD),
		}
	}
	textsRequest := &wluv1.TextsRequest{
		AppId:    appID,
		Requests: wluRequests,
	}

	res, err := wluClient.Texts(ctx, textsRequest)
	return res, annotated, err
}

func removeAnnotations(line string) string {
	removeNormalizedPattern := regexp.MustCompile(`\|.+?]\(([^)]+)\)`)
	line = removeNormalizedPattern.ReplaceAllString(line, "]")

	removablePattern := regexp.MustCompile(`\*([^ ]+)(?: |$)|\(([^)]+)\)`)
	line = removablePattern.ReplaceAllString(line, "")

	entityValuePattern := regexp.MustCompile(`\[([^]]+)]`)
	return entityValuePattern.ReplaceAllStringFunc(line, func(s string) string {
		pipeIndex := strings.Index(s, "|")
		if pipeIndex == -1 {
			pipeIndex = len(s) - 1
		}
		return s[1:pipeIndex]
	})
}

func readLines(fn string) []string {
	if fn != "--" {
		file, err := os.Open(fn)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			err := file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
		return scanLines(file)
	} else {
		return scanLines(os.Stdin)
	}
}

func scanLines(file *os.File) []string {
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func evaluateAnnotatedUtterances(annotatedData []string, groundTruthData []string) {
	if len(annotatedData) != len(groundTruthData) {
		log.Fatalf(
			"Inputs should have same length, but input has %d items and ground-truths %d items.",
			len(annotatedData),
			len(groundTruthData),
		)
	}

	n := float64(len(annotatedData))
	hits := 0.0
	for i, aUtt := range annotatedData {
		gtUtt := groundTruthData[i]
		if strings.TrimSpace(aUtt) == strings.TrimSpace(gtUtt) {
			hits += 1.0
			continue
		}
		fmt.Printf("Line %d: Ground truth had:\n", i+1)
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

func wluResponsesToString(responses []*wluv1.WLUResponse) []string {
	results := make([]string, len(responses))
	for i, resp := range responses {
		segmentStrings := make([]string, len(resp.Segments))
		for j, segment := range resp.Segments {
			segmentStrings[j] = segment.AnnotatedText
		}
		results[i] = strings.Join(segmentStrings, " ")

	}
	return results
}
