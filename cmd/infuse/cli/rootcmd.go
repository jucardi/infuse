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
	usage = "%s [template file] -j [json file]"
	long  = `
Infuse - the templates CLI parser
    Version: V-%s
    Built: %s

Supports:
    - Go templates
    - Handlebars templates
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
	rootCmd.Flags().StringP("json", "j", "", "A JSON file to use as an input for the data to be parsed")
	rootCmd.Flags().StringP("string", "s", "", "The JSON string to use as an input for the data to be parsed")
	rootCmd.Flags().StringP("output", "o", "", "Set output file. If not specified, the resulting template will be printed to Stdout")
	rootCmd.Flags().StringArrayP("definitions", "d", []string{}, "Other templates to be loaded to be used in the 'templates' directive.")
	rootCmd.Flags().StringP("pattern", "p", "", "Uses a search pattern to load definition files to be used in the 'templates' directive.")

	if err := rootCmd.Execute(); err != nil {
		printUsage(rootCmd)
	}
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
	input, _ := cmd.Flags().GetString("json")
	json, _ := cmd.Flags().GetString("string")
	output, _ := cmd.Flags().GetString("output")
	definitions, _ := cmd.Flags().GetStringArray("definitions")
	pattern, _ := cmd.Flags().GetString("pattern")

	request := parser.TemplateRequest{
		Filename:      filename,
		JSON:          json,
		InputFile:     input,
		Definitions:   definitions,
		SearchPattern: pattern,
		Output:        output,
	}

	if err := parser.Parse(request); err != nil {
		log.Panicf("%v", err)
	}
}

func validate(args []string) bool {
	return len(args) == 1

}
