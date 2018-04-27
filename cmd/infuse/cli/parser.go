package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/jucardi/go-infuse/cmd/infuse/version"
	"github.com/jucardi/go-infuse/templates"
	"github.com/jucardi/go-infuse/util/ioutils"
	"github.com/jucardi/go-infuse/util/log"
	"github.com/spf13/cobra"
)

const (
	usage = "%s [template file] [JSON string]"
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

	template := templates.Factory().New(args[0])
	var writer io.WriteCloser

	// Load template definitions.
	if definitions, _ := cmd.Flags().GetStringArray("definitions"); len(definitions) > 0 {
		if err := template.LoadFileDefinition(definitions...); err != nil {
			log.Panicf("failed to load definitions, %v", err)
		}
	}

	// Load definitions by search pattern
	if pattern, _ := cmd.Flags().GetString("pattern"); pattern != "" {
		if err := template.LoadFileDefinitionsByPattern(pattern); err != nil {
			log.Panicf("failed to load definitions, %v", err)
		}
	}

	// Load template
	if err := template.LoadFileTemplate(args[0]); err != nil {
		log.Panicf("failed to load template '%s', %v", args[0], err)
	}

	// Establish the io.Writer to use
	if output, _ := cmd.Flags().GetString("output"); output != "" {
		if w, err := ioutils.NewFileWriter(output); err != nil {
			log.Panicf("Unable to open file '%s', err", output, err)
		} else {
			writer = w
			defer w.Close()
		}
	} else {
		writer = os.Stdout
	}

	// Parse
	if err := template.ParseJSON(writer, args[1]); err != nil {
		log.Panicf("failed to parse the template, %v", err)
	}
}

func validate(args []string) bool {
	return len(args) == 2
}
