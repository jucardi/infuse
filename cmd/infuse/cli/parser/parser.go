package parser

import (
	"errors"
	"fmt"
	"github.com/jucardi/go-logger-lib/log"
	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/util/ioutils"
	"io"
	"os"
)

// TemplateRequest encapsulates the required information to parse a template
type TemplateRequest struct {
	Filename      string
	String        string
	Files         []string
	URL           string
	Output        string
	Definitions   []string
	SearchPattern string
}

func (t TemplateRequest) validate() error {
	if len(t.Files) == 0 && t.String != "" && t.URL != "" {
		return errors.New("only one input method allowed, specify either an input filename, a string or a URL")
	}
	if t.Filename == "" {
		return errors.New("template filename is required")
	}
	return nil
}

func (t TemplateRequest) load() (Data, error) {
	if err := t.validate(); err != nil {
		return nil, err
	}

	dataObj := Data{}

	if len(t.Files) > 0 {
		for _, file := range t.Files {
			if err := dataObj.LoadFile(file); err != nil {
				return nil, err
			}
		}
	}

	if t.URL != "" {
		if err := dataObj.LoadURL(t.URL); err != nil {
			return nil, err
		}
	}
	return dataObj, nil
}

// Parse parses the given template with the given information
func Parse(req TemplateRequest) error {
	log.SetLevel(log.DebugLevel)
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

	if err := template.Parse(writer, data.ToMap()); err != nil {
		return fmt.Errorf("failed to parse the template, %v", err)
	}

	return nil
}
