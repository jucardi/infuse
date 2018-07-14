package parser

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"bytes"
	"github.com/jucardi/infuse/templates"
	"github.com/jucardi/infuse/util/ioutils"
	"net/http"
	"strings"
	"github.com/jucardi/go-logger-lib/log"
)

// TemplateRequest encapsulates the required information to parse a template
type TemplateRequest struct {
	Filename      string
	String        string
	InputFile     string
	URL           string
	Output        string
	Definitions   []string
	SearchPattern string
}

func (t TemplateRequest) validate() error {
	count := 0
	if t.InputFile != "" {
		count++
	}
	if t.String != "" {
		count++
	}
	if t.URL != "" {
		count++
	}
	if count > 1 {
		return errors.New("only one input method allowed, specify either an input filename, a string or a URL")
	}
	if t.Filename == "" {
		return errors.New("template filename is required")
	}
	return nil
}

func (t TemplateRequest) load() ([]byte, error) {
	if err := t.validate(); err != nil {
		return nil, err
	}

	if t.InputFile == "" && t.String == "" && t.URL == "" {
		return []byte("{}"), nil
	}
	if t.InputFile != "" {
		return ioutil.ReadFile(t.InputFile)
	}
	if t.URL != "" {
		if resp, err := http.Get(t.URL); err != nil {
			return nil, fmt.Errorf("error occurred while fetching from URL, %+v", err)
		} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
			return nil, fmt.Errorf("unsuccessful response code (%d)", resp.StatusCode)
		} else {
			buf := new(bytes.Buffer)
			buf.ReadFrom(resp.Body)
			return buf.Bytes(), nil
		}
	}
	return []byte(t.String), nil
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

	filename := ""

	if req.Filename != "" {
		filename = req.Filename
	} else if req.URL != "" {
		filename = req.URL
	}

	split := strings.Split(strings.ToLower(filename), ".")
	fileType := split[len(split)-1]

	// Parse
	switch fileType {
	case "json":
		err = template.ParseJSON(writer, data)
	case "yaml":
		fallthrough
	case "yml":
		err = template.ParseYAML(writer, data)
	default:
		err = template.ParseMarshalled(writer, data)
	}

	if err != nil {
		return fmt.Errorf("failed to parse the template, %v", err)
	}

	return nil
}
