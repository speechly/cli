package cmd

import (
	"bufio"
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/aebruno/nwalgo"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/schollz/progressbar/v3"
	configv1 "github.com/speechly/api/go/speechly/config/v1"
	salv1 "github.com/speechly/api/go/speechly/sal/v1"
	sluv1 "github.com/speechly/api/go/speechly/slu/v1"
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
		aln1, aln2, _ := nwalgo.Align(gtUtt, aUtt, 1, -1, -1)
		if strings.TrimSpace(aUtt) == strings.TrimSpace(gtUtt) {
			hits += 1.0
			continue
		}
		fmt.Printf("\nLine: %d\n", i+1)
		fmt.Printf("└─ Ground truth: %s\n", aln1)
		fmt.Printf("└─ Prediction:   %s\n", aln2)
	}
	fmt.Printf("\nAccuracy: %.2f (%.0f/%.0f)\n", hits/n, hits, n)
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

func readAudioCorpus(filename string) ([]AudioCorpusItem, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	ac := make([]AudioCorpusItem, 0)
	if strings.HasSuffix(filename, "wav") {
		return []AudioCorpusItem{{Audio: filename}}, nil
	}
	jd := json.NewDecoder(f)
	for {
		var aci AudioCorpusItem
		err := jd.Decode(&aci)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		ac = append(ac, aci)
	}
	return ac, nil
}

func readAudio(audioFilePath string, acItem AudioCorpusItem, callback func(buffer audio.IntBuffer, n int) error) error {
	file, err := os.Open(audioFilePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer func() {
		_ = file.Close()
	}()

	ad := wav.NewDecoder(file)
	ad.ReadInfo()
	if !ad.IsValidFile() {
		return fmt.Errorf("audio file is not valid")
	}

	afmt := ad.Format()

	if afmt.NumChannels != 1 || afmt.SampleRate != 16000 || ad.BitDepth != 16 {
		return fmt.Errorf("only audio with 1ch 16kHz 16bit PCM wav files are supported. The audio file is %dch %dHz %dbit",
			afmt.NumChannels, afmt.SampleRate, ad.BitDepth)
	}

	for {
		bfr := audio.IntBuffer{
			Format:         afmt,
			Data:           make([]int, 32768),
			SourceBitDepth: int(ad.BitDepth),
		}
		n, err := ad.PCMBuffer(&bfr)
		if err != nil {
			return fmt.Errorf("pcm buffer creation failed: %v", err)
		}

		if n == 0 {
			break
		}

		err = callback(bfr, n)
		if err != nil {
			return fmt.Errorf("processing read audio failed: %v", err)
		}
	}
	return nil
}

func transcribeWithBatchAPI(ctx context.Context, appID string, corpusPath string, requireGroundTruth bool) ([]AudioCorpusItem, error) {
	client, err := clients.BatchAPIClient(ctx)
	if err != nil {
		return nil, err
	}

	pending := make(map[string]AudioCorpusItem)

	ac, err := readAudioCorpus(corpusPath)
	if err != nil {
		return nil, err
	}
	bar := getBar("Uploading   ", "utt", len(ac))
	for _, aci := range ac {
		if requireGroundTruth && aci.Transcript == "" {
			barClearOnError(bar)
			return nil, fmt.Errorf("missing ground truth")
		}
		paStream, err := client.ProcessAudio(ctx)
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}

		audioFilePath := path.Join(path.Dir(corpusPath), aci.Audio)
		if corpusPath == aci.Audio {
			audioFilePath = corpusPath
		}

		err = readAudio(audioFilePath, aci, func(buffer audio.IntBuffer, n int) error {
			buffer16 := make([]uint16, len(buffer.Data))
			for i, x := range buffer.Data {
				buffer16[i] = uint16(x)
			}
			buf := new(bytes.Buffer)
			err = binary.Write(buf, binary.LittleEndian, buffer16)
			if err != nil {
				return fmt.Errorf("binary.Write: %v", err)
			}

			err = paStream.Send(&sluv1.ProcessAudioRequest{
				AppId: appID,
				Config: &sluv1.AudioConfiguration{
					Encoding:        sluv1.AudioConfiguration_ENCODING_LINEAR16,
					Channels:        1,
					SampleRateHertz: 16000,
				},
				Source: &sluv1.ProcessAudioRequest_Audio{Audio: buf.Bytes()},
			})
			if err != nil {
				return fmt.Errorf("sending %d process audio request failed: %w", buf.Len(), err)
			}
			return nil
		})
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}

		err = bar.Add(1)
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}

		paResp, err := paStream.CloseAndRecv()
		bID := paResp.GetOperation().GetId()
		pending[bID] = aci
	}

	err = bar.Close()
	if err != nil {
		return nil, err
	}

	inputSize := len(pending)
	var results []AudioCorpusItem

	bar = getBar("Transcribing", "utt", inputSize)
	for {
		for bID, aci := range pending {
			status, err := client.QueryStatus(ctx, &sluv1.QueryStatusRequest{Id: bID})
			if err != nil {
				barClearOnError(bar)
				return results, err
			}
			switch status.GetOperation().GetStatus() {
			case sluv1.Operation_STATUS_DONE:
				trs := status.GetOperation().GetTranscripts()
				words := make([]string, len(trs))
				for i, tr := range trs {
					words[i] = tr.Word
				}
				aci := AudioCorpusItem{
					Audio:      aci.Audio,
					Transcript: aci.Transcript,
					Hypothesis: strings.Join(words, " "),
				}
				results = append(results, aci)

				delete(pending, bID)
				err = bar.Add(1)
				if err != nil {
					barClearOnError(bar)
					return results, err
				}
			}
		}
		if len(pending) == 0 {
			break
		}
		time.Sleep(time.Second)
	}
	err = bar.Close()
	if err != nil {
		return results, err
	}

	return results, nil
}

func getBar(desc string, unit string, inputSize int) *progressbar.ProgressBar {
	bar := progressbar.NewOptions(inputSize,
		// Default Options
		progressbar.OptionSetWriter(os.Stderr),
		progressbar.OptionSetWidth(10),
		progressbar.OptionThrottle(65*time.Millisecond),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSpinnerType(14),
		progressbar.OptionSetRenderBlankState(true),
		progressbar.OptionClearOnFinish(),
		// Custom Options
		progressbar.OptionSetWidth(20),
		progressbar.OptionSetDescription(desc),
		progressbar.OptionSetItsString(unit),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "█",
			SaucerHead:    "▌",
			SaucerPadding: "░",
			BarStart:      " ",
			BarEnd:        " ",
		}))
	return bar
}

func transcribeWithStreamingAPI(ctx context.Context, appID string, corpusPath string, requireGroundTruth bool) ([]AudioCorpusItem, error) {
	trIdx := 0
	var (
		results     []AudioCorpusItem
		audios      []string
		transcripts []string
	)

	ac, err := readAudioCorpus(corpusPath)
	if err != nil {
		return nil, err
	}

	bar := getBar("Transcribing", "utt", len(ac))
	for _, aci := range ac {
		client, err := clients.SLUClient(ctx)
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}

		stream, err := client.Stream(ctx)
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}

		done := make(chan error)
		words := make([]string, 0)

		go func() {
			for {
				res, err := stream.Recv()
				if err != nil {
					if err == io.EOF {
						done <- nil
					} else {
						done <- err
					}
					return
				}
				switch r := res.StreamingResponse.(type) {
				case *sluv1.SLUResponse_Started:
				case *sluv1.SLUResponse_Finished:
					aci := AudioCorpusItem{Audio: audios[trIdx], Transcript: transcripts[trIdx], Hypothesis: strings.Join(words, " ")}
					results = append(results, aci)
					trIdx++
					words = make([]string, 0)
				case *sluv1.SLUResponse_Transcript:
					words = append(words, r.Transcript.Word)
				case *sluv1.SLUResponse_Entity:
				case *sluv1.SLUResponse_Intent:
				}
			}
		}()

		err = stream.Send(&sluv1.SLURequest{StreamingRequest: &sluv1.SLURequest_Config{
			Config: &sluv1.SLUConfig{
				Encoding:        sluv1.SLUConfig_LINEAR16,
				Channels:        1,
				SampleRateHertz: 16000,
				LanguageCode:    "en-US",
			},
		}})

		audios = append(audios, aci.Audio)
		transcripts = append(transcripts, aci.Transcript)
		_ = stream.Send(&sluv1.SLURequest{StreamingRequest: &sluv1.SLURequest_Start{Start: &sluv1.SLUStart{
			AppId: appID,
		}}})

		audioFilePath := path.Join(path.Dir(corpusPath), aci.Audio)
		if corpusPath == aci.Audio {
			audioFilePath = corpusPath
		}

		err = readAudio(audioFilePath, aci, func(buffer audio.IntBuffer, n int) error {
			buffer16 := make([]uint16, len(buffer.Data))
			for i, x := range buffer.Data {
				buffer16[i] = uint16(x)
			}
			buf := new(bytes.Buffer)
			err = binary.Write(buf, binary.LittleEndian, buffer16)
			if err != nil {
				return err
			}

			_ = stream.Send(&sluv1.SLURequest{
				StreamingRequest: &sluv1.SLURequest_Audio{
					Audio: buf.Bytes(),
				},
			})
			return nil
		})
		_ = stream.Send(&sluv1.SLURequest{StreamingRequest: &sluv1.SLURequest_Stop{Stop: &sluv1.SLUStop{}}})
		_ = stream.CloseSend()
		if err != nil {
			barClearOnError(bar)
			return results, err
		}

		err = <-done
		if err != nil {
			barClearOnError(bar)
			return results, err
		}
		err = bar.Add(1)
		if err != nil {
			barClearOnError(bar)
			return nil, err
		}
	}
	err = bar.Close()
	if err != nil {
		return results, err
	}

	return results, nil
}
func barClearOnError(_ *progressbar.ProgressBar) {
	_, _ = fmt.Fprint(os.Stderr, "\n\n")
}
