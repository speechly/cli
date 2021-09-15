package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"regexp"
	"sort"
	"text/tabwriter"

	"github.com/spf13/cobra"

	salv1 "github.com/speechly/api/go/speechly/sal/v1"
	"github.com/speechly/cli/pkg/clients"
	"github.com/speechly/cli/pkg/upload"
)

type CompileWriter struct {
	appId     string
	stream    salv1.Compiler_CompileClient
	batchSize int32
	seed      int32
}

func (u CompileWriter) Write(data []byte) (n int, err error) {
	contentType := salv1.AppSource_CONTENT_TYPE_TAR
	as := &salv1.AppSource{AppId: u.appId, DataChunk: data, ContentType: contentType}
	req := &salv1.CompileRequest{AppSource: as, BatchSize: u.batchSize, RandomSeed: u.seed}
	if err = u.stream.Send(req); err != nil {
		return 0, err
	}
	return len(data), nil
}

var sampleCmd = &cobra.Command{
	Use: "sample [directory]",
	Example: `speechly sample -a UUID_APP_ID .
speechly sample -a UUID_APP_ID /usr/local/project/app
speechly sample -a UUID_APP_ID /usr/local/project/app --stats`,
	Short: "Sample a set of examples from the given SAL configuration",
	Long: `The contents of the directory given as argument is sent to the
API and compiled. If configuration is valid, a set of examples are printed to stdout.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		batchSize, _ := cmd.Flags().GetInt("batch-size")
		if batchSize < 32 || batchSize > 10000 {
			log.Fatal("Batch size must be between 32 and 10000")
		}
		seed, _ := cmd.Flags().GetInt("seed")

		uploadData := upload.CreateTarFromDir(args[0])

		if len(uploadData.Files) == 0 {
			log.Fatal("No files to upload!\n\nPlease ensure the files are named *.yaml or *.csv")
		}

		// open a stream for upload
		compileClient, err := clients.CompileClient(ctx)
		if err != nil {
			log.Fatalf("Error connecting to API: %s", err)
		}
		stream, err := compileClient.Compile(ctx)
		if err != nil {
			log.Fatalf("Failed to open validate stream: %s", err)
		}

		// flush the tar from memory to the stream
		compileWriter := CompileWriter{appId, stream, int32(batchSize), int32(seed)}
		_, err = uploadData.Buf.WriteTo(compileWriter)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}
		log.Printf("Sampling %d examples \n", batchSize)

		compileResult, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Validate failed: %s", err)
		}

		if len(compileResult.Messages) > 0 {
			printLineErrors(compileResult.Messages)
		} else {
			simpleStats, _ := cmd.Flags().GetBool("stats")
			advancedStats, _ := cmd.Flags().GetBool("advanced-stats")
			limit, _ := cmd.Flags().GetInt("advanced-stats-limit")
			if simpleStats || advancedStats {
				printStats(cmd.OutOrStdout(), compileResult.Templates, simpleStats, advancedStats, int32(limit))
			} else {
				for _, message := range compileResult.Templates {
					fmt.Fprintf(cmd.OutOrStdout(), "%s\n", message)
				}
			}
		}
	},
}

type IntentEntityCounter struct {
	entityCounts      map[string]map[string]map[string]float32
	intentCounts      map[string]float32
	intentsRe         *regexp.Regexp
	entityRe          *regexp.Regexp
	utteranceCnt      float32
	entityToInt       map[string]int
	entityToType      map[string]string
	entityCooccurance map[string][][]bool
}

type Entity struct {
	entType string
	entVal  string
}

type Intent struct {
	name         string
	subUtterance []byte
}

type ResultRow struct {
	Name       string
	Count      int
	Distrib    float32
	Proportion float32
}

func (this *IntentEntityCounter) findEntities(subUtterance []byte) []Entity {
	result := make([]Entity, 0)
	for _, entity := range this.entityRe.FindAll(subUtterance, -1) {
		trippedEntity := entity[1 : len(entity)-1]
		splittedEntity := bytes.Split(trippedEntity, []byte(`](`))
		result = append(result, Entity{string(splittedEntity[1]), string(splittedEntity[0])})
	}
	return result
}

func (this *IntentEntityCounter) findIntents(utterance []byte) []Intent {
	result := make([]Intent, 0)
	intentPlaces := this.intentsRe.FindAllIndex(utterance, -1)
	intentStart := 0
	intentEnd := 0
	for i, intentNameStartEnd := range intentPlaces {
		if len(intentPlaces)-1 == i {
			// last intent ends in the end of utterance
			intentEnd = len(utterance)
		} else {
			// intent ends where next one starts
			intentEnd = intentPlaces[i+1][0]
		}
		intentNameStart := intentNameStartEnd[0] + 1
		intentNameEnd := intentNameStartEnd[1] - 1
		intentName := string(utterance[intentNameStart:intentNameEnd])
		subUtterance := utterance[intentStart:intentEnd]
		intentStart = intentEnd
		result = append(result, Intent{intentName, subUtterance})
	}
	return result
}

// CountSingle counts the occurance of individual intents and entities
func (this *IntentEntityCounter) CountSingle(utterance []byte) {
	for _, intent := range this.findIntents(utterance) {
		this.intentCounts[intent.name] += 1
		for _, ent := range this.findEntities(intent.subUtterance) {
			if _, ok := this.entityCounts[intent.name]; !ok {
				this.entityCounts[intent.name] = make(map[string]map[string]float32)
			}
			if _, ok := this.entityCounts[intent.name][ent.entType]; !ok {
				this.entityCounts[intent.name][ent.entType] = make(map[string]float32)
			}
			this.entityCounts[intent.name][ent.entType][ent.entVal] += 1
		}
		this.utteranceCnt += 1
	}
}

// CountEntityPairs counts the co-occurance of multiple entities in each intent
func (this *IntentEntityCounter) CountEntityPairs(utterance []byte) {
	if this.entityToInt == nil {
		entityToInt := make(map[string]int)
		entityToType := make(map[string]string)
		for _, entTypes := range this.entityCounts {
			for entType, entValues := range entTypes {
				for entVal := range entValues {
					if _, ok := entityToInt[entVal]; !ok {
						entityToInt[entVal] = len(entityToInt)
					}
					entityToType[entVal] = entType
				}
			}
		}
		this.entityToInt = entityToInt
		this.entityToType = entityToType
		this.entityCooccurance = make(map[string][][]bool)
	}
	for _, intent := range this.findIntents(utterance) {
		occurences := make([]bool, len(this.entityToInt))
		for _, ent := range this.findEntities(intent.subUtterance) {
			occurences[this.entityToInt[ent.entVal]] = true
		}
		this.entityCooccurance[intent.name] = append(this.entityCooccurance[intent.name], occurences)
	}
}

func normalize(rows []ResultRow, total float32) []ResultRow {
	result := make([]ResultRow, 0)
	for _, row := range rows {
		result = append(result, ResultRow{row.Name, row.Count, (row.Distrib / total), row.Proportion})
	}
	return result
}

func sortByCount(rows []ResultRow) {
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].Count == rows[j].Count {
			return rows[i].Name < rows[j].Name
		}
		return rows[i].Count > rows[j].Count
	})
}

func (this *IntentEntityCounter) generatePairs() [][]string {
	res := make([][]string, 0)
	for k1, v1 := range this.entityToInt {
		for k2, v2 := range this.entityToInt {
			if v1 < v2 && this.entityToType[k1] != this.entityToType[k2] {
				res = append(res, []string{k1, k2})
			}
		}
	}
	return res
}

func (this *IntentEntityCounter) GetIntentCounts() []ResultRow {
	result := make([]ResultRow, 0)
	var total float32
	for intentName, intentCnt := range this.intentCounts {
		total += intentCnt
		row := ResultRow{intentName, int(intentCnt), intentCnt, (intentCnt / this.utteranceCnt)}
		result = append(result, row)
	}
	sortByCount(result)
	return normalize(result, total)
}

func (this *IntentEntityCounter) GetEntityTypeCounts() []ResultRow {
	result := make([]ResultRow, 0)
	entityTypeCounts := make(map[string]float32)
	var total float32
	for _, entTypes := range this.entityCounts {
		for entType, entValues := range entTypes {
			for _, cnt := range entValues {
				entityTypeCounts[entType] += cnt
				total += cnt
			}
		}
	}
	for entityType, cnt := range entityTypeCounts {
		row := ResultRow{entityType, int(cnt), cnt, (cnt / this.utteranceCnt)}
		result = append(result, row)
	}
	sortByCount(result)
	return normalize(result, total)
}

func (this *IntentEntityCounter) GetEntityValueCounts() []ResultRow {
	result := make([]ResultRow, 0)
	entityValueCounts := make(map[string]float32)
	var total float32
	for _, entTypes := range this.entityCounts {
		for _, entValues := range entTypes {
			for entVal, cnt := range entValues {
				entityValueCounts[entVal] += cnt
				total += cnt
			}
		}
	}
	for entityVal, cnt := range entityValueCounts {
		row := ResultRow{entityVal, int(cnt), cnt, (cnt / this.utteranceCnt)}
		result = append(result, row)
	}
	sortByCount(result)
	return normalize(result, total)
}

func (this *IntentEntityCounter) GetIntentEntityTypeCounts() []ResultRow {
	result := make([]ResultRow, 0)
	var total float32
	for intentName, entTypes := range this.entityCounts {
		for entType, entValues := range entTypes {
			var entTypeCnt float32
			for _, cnt := range entValues {
				entTypeCnt += cnt
				total += cnt
			}
			name := intentName + "(" + entType + "=*)"
			row := ResultRow{name, int(entTypeCnt), entTypeCnt, (entTypeCnt / this.utteranceCnt)}
			result = append(result, row)
		}
	}
	sortByCount(result)
	return normalize(result, total)
}

func (this *IntentEntityCounter) GetIntentEntityValueCounts() []ResultRow {
	result := make([]ResultRow, 0)
	var total float32
	for intentName, entTypes := range this.entityCounts {
		for entType, entValues := range entTypes {
			for entValue, cnt := range entValues {
				total += cnt
				name := intentName + "(" + entType + "=" + entValue + ")"
				row := ResultRow{name, int(cnt), cnt, (cnt / this.utteranceCnt)}
				result = append(result, row)
			}
		}
	}
	sortByCount(result)
	return normalize(result, total)
}

func (this *IntentEntityCounter) GetIntentEntityValuePairCounts() []ResultRow {
	result := make([]ResultRow, 0)
	var total float32
	combinations := this.generatePairs()
	for intentName, occurences := range this.entityCooccurance {
		for _, comb := range combinations {
			ind1 := this.entityToInt[comb[0]]
			ind2 := this.entityToInt[comb[1]]
			var cnt float32
			for _, occurance := range occurences {
				if occurance[ind1] && occurance[ind2] {
					cnt += 1.0
					total += 1.0
				}
			}
			name := intentName + "(" + this.entityToType[comb[0]] + "=" + comb[0]
			name += "," + this.entityToType[comb[1]] + "=" + comb[1] + ")"
			row := ResultRow{name, int(cnt), cnt, (cnt / this.utteranceCnt)}
			result = append(result, row)
		}
	}
	sortByCount(result)
	return normalize(result, total)
}

func CreateCounter(examples []string, advanced bool) IntentEntityCounter {
	intentsRe := regexp.MustCompile(`\*(.*?) `)
	entityRe := regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	entityCounts := make(map[string]map[string]map[string]float32)
	intentCounts := make(map[string]float32)
	counter := IntentEntityCounter{entityCounts, intentCounts, intentsRe, entityRe, 0.0, nil, nil, nil}
	for _, example := range examples {
		counter.CountSingle([]byte(example))
	}
	if advanced {
		for _, example := range examples {
			counter.CountEntityPairs([]byte(example))
		}
	}
	return counter
}

func printLines(out io.Writer, name string, rows []ResultRow, lineLimit int32) {
	// Format in tab-separated columns with a tab stop of 8.
	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, name+"\tCOUNT\tDISTRIBUTION\tAVG PER UTTERANCE\n")
	length := int32(len(rows))
	if lineLimit == -1 || lineLimit > length {
		lineLimit = length
	}
	for _, row := range rows[:lineLimit] {
		fmt.Fprintf(w, "%s\t%d\t%f\t%f\n", row.Name, row.Count, row.Distrib, row.Proportion)
	}
	if err := w.Flush(); err != nil {
		log.Fatalf("When printing the section %s, error occured: %s", name, err)
	}
}

func printStats(out io.Writer, examples []string, normal bool, advanced bool, lineLimit int32) {
	counter := CreateCounter(examples, advanced)

	log.Printf("There was %d utterances in the sample of %d examples \n", int32(counter.utteranceCnt), len(examples))
	if normal {
		printLines(out, "INTENTS", counter.GetIntentCounts(), -1)
		printLines(out, "ENTITY TYPES", counter.GetEntityTypeCounts(), -1)
		printLines(out, "ENTITY VALUES", counter.GetEntityValueCounts(), -1)
	}
	if advanced {
		printLines(out, "ENTITY TYPES PER INTENT", counter.GetIntentEntityTypeCounts(), lineLimit)
		printLines(out, "ENTITY VALUES PER INTENT", counter.GetIntentEntityValueCounts(), lineLimit)
		printLines(out, "ENTITY VALUE PAIRS PER INTENT", counter.GetIntentEntityValuePairCounts(), lineLimit)
	}
}

func init() {
	rootCmd.AddCommand(sampleCmd)
	sampleCmd.Flags().StringP("app", "a", "", "application to deploy the files to.")
	sampleCmd.Flags().Int("batch-size", 100, "how many examples to return. Must be between 32 and 10000")
	sampleCmd.Flags().Int("seed", 0, "random seed to use when initializing the sampler.")

	sampleCmd.Flags().Bool("stats", false, "print intent and entity distributions to the output.")
	sampleCmd.Flags().Bool("advanced-stats", false, "print entity type, value and value pair distributions to the output.")
	sampleCmd.Flags().Int("advanced-stats-limit", 10, "line limit for advanced_stats. The lines are ordered by count.")
	sampleCmd.Flags().SortFlags = false
	if err := sampleCmd.MarkFlagRequired("app"); err != nil {
		log.Fatalf("failed to init flags: %v", err)
	}
}
