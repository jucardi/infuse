package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/templates/helpers"
	"gopkg.in/jucardi/go-streams.v1/streams"
	"gopkg.in/jucardi/go-strings.v1/stringx"

	"github.com/jucardi/infuse/cmd/infuse/cli/parser"
	"github.com/jucardi/infuse/cmd/infuse/version"
	"github.com/jucardi/infuse/util/log"
	"github.com/spf13/cobra"
)

const (
	usage = `%s [template file] -i [JSON or YAML file] -u [URL to GET JSON or INPUT from] -s [JSON or YAML string] -o [output] -p [pattern] -d [template path 1] -d [template path 2]

  - All flags are optional
  - Max one input type allowed (-f, -s, -u)`
	long = `
Infuse - the templates CLI parser
    Version: V-%s
    Built: %s

Supports:
    - Go templates
    - Handlebars templates (coming soon)
`
)

var rootCmd = &cobra.Command{
	Use:              "infuse",
	Short:            "Parses a Golang template",
	Long:             fmt.Sprintf(long, version.Version, version.Built),
	PersistentPreRun: initCmd,
	Run:              parse,
}

// Execute starts the execution of the parse command.
func Execute() {
	rootCmd.Flags().StringArrayP("file", "f", nil, "INPUT: A JSON or YAML file to use as an input for the data to be parsed")
	rootCmd.Flags().StringP("string", "s", "", "INPUT: A JSON or YAML string representation")
	rootCmd.Flags().StringP("url", "u", "", "INPUT: A URL to HTTP GET a JSON or YAML file from. Useful to parse data from config servers")
	rootCmd.Flags().StringP("output", "o", "", "Set output file. If not specified, the resulting template will be printed to Stdout")
	rootCmd.Flags().StringP("pattern", "p", "", "Uses a search pattern to load definition files to be used in the 'templates' directive.")
	rootCmd.Flags().StringArrayP("definition", "d", []string{}, "Other templates to be loaded to be used in the 'templates' directive.")
	rootCmd.Flags().BoolP("listHelpers", "l", false, "Lists all registered helpers")

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func printUsage(cmd *cobra.Command) {
	cmd.Println(fmt.Sprintf(long, version.Version, version.Built))
	_ = cmd.Usage()
}

func initCmd(cmd *cobra.Command, _ []string) {
	FromCommand(cmd)
	cmd.Use = fmt.Sprintf(usage, cmd.Use)
}

func parse(cmd *cobra.Command, args []string) {
	if listHelpers, _ := cmd.Flags().GetBool("listHelpers"); listHelpers {
		printHelpers()
		os.Exit(0)
	}

	if !validate(args) {
		log.Error("Unexpected number of arguments")
		printUsage(cmd)
		os.Exit(-1)
	}

	filename := args[0]
	input, _ := cmd.Flags().GetStringArray("file")
	str, _ := cmd.Flags().GetString("string")
	url, _ := cmd.Flags().GetString("url")
	output, _ := cmd.Flags().GetString("output")
	definitions, _ := cmd.Flags().GetStringArray("definition")
	pattern, _ := cmd.Flags().GetString("pattern")

	request := parser.TemplateRequest{
		Filename:      filename,
		String:        str,
		URL:           url,
		Files:         input,
		Definitions:   definitions,
		SearchPattern: pattern,
		Output:        output,
	}

	if err := parser.Parse(request); err != nil {
		log.Errorf("%v", err)
		printUsage(cmd)
		os.Exit(-1)
	}
}

func validate(args []string) bool {
	return len(args) == 1
}

func printHelpers() {
	println("Available helpers:")
	for k, list := range helpersByCategory() {
		maxLenght := 0
		fmt.Printf("\n  %s\n", k)
		for _, h := range list {
			if len(h.Name) > maxLenght {
				maxLenght = len(h.Name)
			}
		}
		for _, h := range list {
			if h.Description == "" {
				fmt.Printf("   - %s\n", h.Name)
			} else {
				fmt.Printf("   - %s%s      > %s\n", h.Name, getSpaces(maxLenght-len(h.Name)), h.Description)
			}
		}
	}
	println()
}

func helpersByCategory() map[string][]*helpers.Helper {
	template := templates.Factory().New()
	ret := map[string][]*helpers.Helper{}

	for _, h := range template.Helpers() {
		list, _ := ret[h.Category]
		list = append(list, h)
		ret[h.Category] = list
	}
	for cat, list := range ret {
		ret[cat] = streams.
			From(list).
			OrderBy(func(i interface{}, j interface{}) int {
				x := i.(*helpers.Helper)
				y := j.(*helpers.Helper)

				return strings.Compare(x.Name, y.Name)
			}).
			ToArray().([]*helpers.Helper)
	}
	return ret
}
func getSpaces(count int) string {
	builder := stringx.Builder()
	for i := 0; i < count; i++ {
		builder.Append(" ")
	}
	return builder.Build()
}
