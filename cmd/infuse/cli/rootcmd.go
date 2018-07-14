package cli

import (
	"fmt"
	"os"

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
	rootCmd.Flags().StringP("file", "f", "", "INPUT: A JSON or YAML file to use as an input for the data to be parsed")
	rootCmd.Flags().StringP("string", "s", "", "INPUT: A JSON or YAML string representation")
	rootCmd.Flags().StringP("url", "u", "", "INPUT: A URL to HTTP GET a JSON or YAML file from. Useful to parse data from config servers")
	rootCmd.Flags().StringP("output", "o", "", "Set output file. If not specified, the resulting template will be printed to Stdout")
	rootCmd.Flags().StringP("pattern", "p", "", "Uses a search pattern to load definition files to be used in the 'templates' directive.")
	rootCmd.Flags().StringArrayP("definitions", "d", []string{}, "Other templates to be loaded to be used in the 'templates' directive.")

	rootCmd.Execute()
}

func printUsage(cmd *cobra.Command) {
	cmd.Println(fmt.Sprintf(long, version.Version, version.Built))
	cmd.Usage()
}

func initCmd(cmd *cobra.Command, args []string) {
	FromCommand(cmd)
	cmd.Use = fmt.Sprintf(usage, cmd.Use)
}

func parse(cmd *cobra.Command, args []string) {
	if !validate(args) {
		log.Error("Unexpected number of arguments")
		printUsage(cmd)
		os.Exit(-1)
	}

	filename := args[0]
	input, _ := cmd.Flags().GetString("file")
	str, _ := cmd.Flags().GetString("string")
	url, _ := cmd.Flags().GetString("url")
	output, _ := cmd.Flags().GetString("output")
	definitions, _ := cmd.Flags().GetStringArray("definitions")
	pattern, _ := cmd.Flags().GetString("pattern")

	request := parser.TemplateRequest{
		Filename:      filename,
		String:        str,
		URL:           url,
		InputFile:     input,
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
