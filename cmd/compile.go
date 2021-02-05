package cmd

import (
	"fmt"
	"io"
	"log"
	"text/tabwriter"
	"regexp"
	"sort"

	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use: "compile [directory]",
	Example: `speechly compile -a UUID_APP_ID .
speechly compile -a UUID_APP_ID /usr/local/project/app`,
	Short: "Compiles a sample of examples from the given configuration",
	Long: `The contents of the directory given as argument is sent to the
API and compiled. If suffcessful, a sample of examples are printed to stdout.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()
		appId, _ := cmd.Flags().GetString("app")
		uploadData := createAndValidateTar(args[0])

		// open a stream for upload
		stream, err := compile_client.Compile(ctx)
		if err != nil {
			log.Fatalf("Failed to open validate stream: %s", err)
		}

		// flush the tar from memory to the stream
		compileWriter := CompileWriter{appId, stream}
		_, err = uploadData.buf.WriteTo(compileWriter)
		if err != nil {
			log.Fatalf("Streaming file data failed: %s", err)
		}

		compileResult, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Validate failed: %s", err)
		}
		
		if len(compileResult.Messages) > 0 {
			printLineErrors(compileResult.Messages)
		} else {
			for _, message := range compileResult.Templates {
				log.Printf("%s", message)
			}
		}

		includeStats, _ := cmd.Flags().GetBool("stats")
		if includeStats {
			printStats(cmd.OutOrStdout(), compileResult.Templates)
		}
	},
}

func addRegexMatchesToMap(re *regexp.Regexp, bytes []byte, counts map[string]float32) {
	for _, val := range re.FindAll(bytes, -1) {
		key := string(val[1:(len(val) - 1)])
		counts[key] += 1
	}
}

func sumValues(counts map[string]float32) float32 {
	var sum float32
	for _,v := range counts {
		sum += v
	}
	return sum
}

func GetIntentAndEntityStats(examples []string) (map[string]float32, map[string]float32, map[string]float32) {
	intentsRe := regexp.MustCompile(`\*(.*?) `)
	entityValRe := regexp.MustCompile(`\[(.*?)\]`)
	entityTypeRe := regexp.MustCompile(`\((.*?)\)`)
	intents := make(map[string]float32, 0)
	entityTypes := make(map[string]float32, 0)
	entityValues := make(map[string]float32, 0)
	for _, example := range examples {
		bytes := []byte(example)
		addRegexMatchesToMap(intentsRe, bytes, intents)
		addRegexMatchesToMap(entityValRe, bytes, entityValues)
		addRegexMatchesToMap(entityTypeRe, bytes, entityTypes)
	}

	return intents, entityTypes, entityValues
}

func printLines(out io.Writer, name string, counts map[string]float32) error {
	total := sumValues(counts)
	// Sort keys by count
	keys := make([]string, 0)
	for k,_ := range counts {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return counts[keys[i]] > counts[keys[j]]
	})

	// Format in tab-separated columns with a tab stop of 8.
	w := tabwriter.NewWriter(out, 0, 8, 1, '\t', 0)
	fmt.Fprint(w, "\n")
	fmt.Fprint(w, name + "\tCOUNT\tPROPORTION\n")
	for _,k := range keys {
		v := counts[k]
		fmt.Fprintf(w, "%s\t%d\t%f\n", k, int(v), (v / total))
	}

	return w.Flush()
}

func printStats(out io.Writer, examples []string) {
	intents, entityTypes, entityValues := GetIntentAndEntityStats(examples)
	if err := printLines(out, "INTENT", intents); err != nil {
		log.Fatalf("Error listing intents: %s", err)
	}
	if err := printLines(out, "ENTITY TYPE", entityTypes); err != nil {
		log.Fatalf("Error listing entity types: %s", err)
	}
	if err := printLines(out, "ENTITY VALUE", entityValues); err != nil {
		log.Fatalf("Error listing entity values: %s", err)
	}
}

func init() {
	rootCmd.AddCommand(compileCmd)
	compileCmd.Flags().StringP("app", "a", "", "application to deploy the files to.")
	compileCmd.Flags().Bool("stats", false, "include intent and entity distributions to output.")
	compileCmd.MarkFlagRequired("app")
}
