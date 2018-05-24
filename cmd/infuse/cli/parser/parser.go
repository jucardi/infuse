package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/util/ioutils"
)

// TemplateRequest encapsulates the required information to parse a template
type TemplateRequest struct {
	Filename      string
	JSON          string
	InputFile     string
	Output        string
	Definitions   []string
	SearchPattern string
}

func (t TemplateRequest) validate() error {
	if t.InputFile == "" && t.JSON == "" {
		return errors.New("an input file or a json string is required as a parsing data source")
	}
	if t.InputFile != "" && t.JSON != "" {
		return errors.New("specify either an input file or a json string but not both")
	}
	if t.Filename == "" {
		return errors.New("template filename is required")
	}
	return nil
}

func (t TemplateRequest) load() ([]byte, error) {
	if t.InputFile != "" {
		return ioutil.ReadFile(t.InputFile)
	}
	return []byte(t.JSON), nil
}

// Parse parses the given template with the given information
func Parse(req TemplateRequest) error {
	if err := req.validate(); err != nil {
		return err
	}

	data, err := req.load()

	if err != nil {
		return fmt.Errorf("unable to load data, %v", err)
	}

	template := templates.Factory().New(req.Filename)
	var writer io.WriteCloser

	// Load template
	if err := template.LoadFileTemplate(req.Filename); err != nil {
		return fmt.Errorf("failed to load template '%s', %v", req.Filename, err)
	}

	// Load template definitions.
	if len(req.Definitions) > 0 {
		if err := template.LoadFileDefinition(req.Definitions...); err != nil {
			return fmt.Errorf("failed to load definitions, %v", err)
		}
	}

	// Load definitions by search pattern
	if req.SearchPattern != "" {
		if err := template.LoadFileDefinitionsByPattern(req.SearchPattern); err != nil {
			return fmt.Errorf("failed to load definitions, %v", err)
		}
	}

	// Establish the io.Writer to use
	if req.Output != "" {
		if w, err := ioutils.NewFileWriter(req.Output); err != nil {
			return fmt.Errorf("unable to open file '%s', %v", req.Output, err)
		} else {
			writer = w
			defer w.Close()
		}
	} else {
		writer = os.Stdout
	}

	// Parse
	if err := template.ParseJSON(writer, data); err != nil {
		return fmt.Errorf("failed to parse the template, %v", err)
	}

	return nil
}
